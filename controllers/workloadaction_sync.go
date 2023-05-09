package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"net/url"
	"reflect"
	"time"

	rabbitstalkerv1alpha1 "docplanner.com/rabbit-stalker/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tidwall/gjson"
)

const (
	// RabbitAdminApiQueueEndpoint TODO
	// The endpoint structure looks like this: /api/queues/vhost/name
	RabbitAdminApiQueueEndpoint = "/api/queues/%s/%s"

	// TODO
	HttpHeaderAccept   = "application/json"
	HttpRequestTimeout = 5 * time.Second

	// TODO
	AnnotationRestartedAt = "rabbit-stalker.docplanner.com/restartedAtt"

	// parseSyncTimeError error message for invalid value on 'synchronization' parameter
	parseSyncTimeError = "Can not parse the synchronization time from workloadAction: %s"

	WorkloadRestartNotSupportedErrorMessage   = "restarting is not supported"
	PausedWorkloadRestartErrorMessage         = "can't restart paused deployment (run rollout resume first)"
	WorkloadActionAnnotationPatchErrorMessage = "impossible to patch the annotations for the workload template: %s"
	UnauthorizedRequestErrorMessage           = "unauthorized request against the connection. Set the credentials for this server"
	InvalidJsonErrorMessage                   = "invalid json from the HTTP response"
	InvalidActionErrorMessage                 = "invalid action. supported: restart, delete"
	DeleteActionNotImplementedErrorMessage    = "deletion action is not implemented yet"

	// TODO
	ResourceKindDeployment  = "Deployment"
	ResourceKindStatefulSet = "StatefulSet"
	ResourceKindDaemonSet   = "DaemonSet"

	// TODO
	ActionDelete  = "delete"
	ActionRestart = "restart"
)

// GetSynchronizationTime return the spec.synchronization.time as duration, or default time on failures
func (r *WorkloadActionReconciler) GetSynchronizationTime(workloadActionManifest *rabbitstalkerv1alpha1.WorkloadAction) (synchronizationTime time.Duration, err error) {
	synchronizationTime, err = time.ParseDuration(workloadActionManifest.Spec.Synchronization.Time)
	if err != nil {
		err = NewErrorf(parseSyncTimeError, workloadActionManifest.Name)
		return synchronizationTime, err
	}

	return synchronizationTime, err
}

// GetSecretResource TODO
func (r *WorkloadActionReconciler) GetSecretResource(ctx context.Context, namespace string, name string) (secret *corev1.Secret, err error) {

	secretResource := corev1.Secret{}

	err = r.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &secretResource)

	secret = &secretResource
	return secret, err
}

// GetCredentialsResources TODO
func (r *WorkloadActionReconciler) GetCredentialsResources(ctx context.Context, workloadActionManifest *rabbitstalkerv1alpha1.WorkloadAction) (resources []*corev1.Secret, err error) {

	usernameSecret := &corev1.Secret{}
	passwordSecret := &corev1.Secret{}

	secretNamespace := workloadActionManifest.Namespace

	// Get the Secret with the username field inside
	secretName := workloadActionManifest.Spec.RabbitConnection.Credentials.Username.SecretRef.Name
	usernameSecret, usernameErr := r.GetSecretResource(ctx, secretNamespace, secretName)

	// Get the Secret with the password field inside
	secretName = workloadActionManifest.Spec.RabbitConnection.Credentials.Password.SecretRef.Name
	passwordSecret, passwordErr := r.GetSecretResource(ctx, secretNamespace, secretName)
	if usernameErr != nil || passwordErr != nil {
		return resources, err
	}
	resources = append(resources, usernameSecret)
	resources = append(resources, passwordSecret)

	return resources, err
}

// GetWorkloadResource TODO
func (r *WorkloadActionReconciler) GetWorkloadResource(ctx context.Context, workloadActionManifest *rabbitstalkerv1alpha1.WorkloadAction) (resource *unstructured.Unstructured, err error) {

	// Get the target manifest
	target := unstructured.Unstructured{}
	target.SetGroupVersionKind(workloadActionManifest.Spec.WorkloadRef.GroupVersionKind())

	err = r.Get(ctx, client.ObjectKey{
		Namespace: workloadActionManifest.Spec.WorkloadRef.Namespace,
		Name:      workloadActionManifest.Spec.WorkloadRef.Name,
	}, &target)

	resource = &target
	return resource, err
}

// SetWorkloadRestartAnnotation restart a workload by changing an annotation.
// This will trigger an automatic reconciliation on the workload in the same way done by kubectl
func (r *WorkloadActionReconciler) SetWorkloadRestartAnnotation(ctx context.Context, obj *unstructured.Unstructured) (err error) {

	var templateAnnotations map[string]interface{}

	resourceType := obj.GetObjectKind().GroupVersionKind()
	// TODO: Check the group version just in case future changes on Kubernetes APIs
	//groupVersion := resourceType.GroupVersion()

	// 1. Check allowed workload types
	kind := resourceType.Kind
	if kind != ResourceKindDeployment && kind != ResourceKindDaemonSet && kind != ResourceKindStatefulSet {
		return fmt.Errorf(WorkloadRestartNotSupportedErrorMessage)
	}

	// 2. Pay special attention on paused deployments
	if kind == ResourceKindDeployment {
		pausedValue, pausedFound, err := unstructured.NestedBool(obj.Object, "spec", "paused")
		if err != nil {
			return err
		}

		if pausedFound && pausedValue == true {
			err = errors.New(PausedWorkloadRestartErrorMessage)
			return err
		}
	}

	// 3. Modify template annotations (spec.template.metadata.annotations) to include AnnotationRestartedAt
	templateAnnotations, templateAnnotationsFound, err := unstructured.NestedMap(obj.Object, "spec", "template", "metadata", "annotations")
	if err != nil {
		return err
	}

	// TODO: Take care about annotations not being present
	if !templateAnnotationsFound || templateAnnotations == nil {
		// TODO: Craft them. May be they are not important?
	}
	templateAnnotations[AnnotationRestartedAt] = time.Now().Format(time.RFC3339)

	// 4. Actually update the workload object against Kubernetes API
	parsedTemplateAnnotations, err := json.Marshal(templateAnnotations)
	if err != nil {
		return err
	}

	patchBytes := []byte(fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":%s}}}}`, parsedTemplateAnnotations))

	err = r.Patch(ctx, obj, client.RawPatch(types.StrategicMergePatchType, patchBytes))
	if err != nil {
		err = errors.New(fmt.Sprintf(WorkloadActionAnnotationPatchErrorMessage, err))
	}

	return err
}

// WorkloadActionTarget call Kubernetes API to actually workloadAction the resource
func (r *WorkloadActionReconciler) reconcileWorkloadAction(ctx context.Context, workloadActionManifest *rabbitstalkerv1alpha1.WorkloadAction) (err error) {

	// 1. Look for the Secret resources only when credentials refs are set in WorkloadAction
	var credentialsObjects []*corev1.Secret

	if !reflect.ValueOf(workloadActionManifest.Spec.RabbitConnection.Credentials).IsZero() {
		credentialsObjects, err = r.GetCredentialsResources(ctx, workloadActionManifest)
		if err != nil {
			r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
				metav1.ConditionFalse,
				ConditionReasonCredentialsNotFound,
				ConditionReasonCredentialsNotFoundMessage,
			))
			return err
		}
	}

	// 2. Fill the credentials only when the Secret resources where extracted
	var username, password []byte
	var usernameFound, passwordFound bool

	if len(credentialsObjects) == 2 {
		username, usernameFound = credentialsObjects[0].Data[workloadActionManifest.Spec.RabbitConnection.Credentials.Username.SecretRef.Key]
		password, passwordFound = credentialsObjects[1].Data[workloadActionManifest.Spec.RabbitConnection.Credentials.Password.SecretRef.Key]
		if !usernameFound || !passwordFound {
			r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
				metav1.ConditionFalse,
				ConditionReasonCredentialsNotFound,
				ConditionReasonCredentialsNotFoundMessage,
			))
			return err
		}
	}

	// 3. Get the workload object
	targetObject, err := r.GetWorkloadResource(ctx, workloadActionManifest)
	if err != nil {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonWorkloadNotFound,
			ConditionReasonWorkloadNotFoundMessage,
		))
		return err
	}

	// 4. Get the params for the HTTP request
	urlString := workloadActionManifest.Spec.RabbitConnection.Url
	urlString += fmt.Sprintf(RabbitAdminApiQueueEndpoint,
		workloadActionManifest.Spec.RabbitConnection.Vhost,
		workloadActionManifest.Spec.RabbitConnection.Queue)

	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonHttpRequestExecutionFailed,
			ConditionReasonHttpRequestExecutionFailedMessage,
		))
		return err
	}

	// 5. Make the HTTP request to the RabbitMQ admin API
	// TODO: Move whole step 4) to a function
	httpClient := http.Client{Timeout: HttpRequestTimeout}

	request, err := http.NewRequest(http.MethodGet, parsedUrl.String(), http.NoBody)
	if err != nil {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonHttpRequestExecutionFailed,
			ConditionReasonHttpRequestExecutionFailedMessage,
		))
		return err
	}

	request.Header.Add("Accept", HttpHeaderAccept)
	request.SetBasicAuth(string(username), string(password))

	response, err := httpClient.Do(request)
	if err != nil {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonHttpRequestExecutionFailed,
			ConditionReasonHttpRequestExecutionFailedMessage,
		))
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonCredentialsNotValid,
			ConditionReasonCredentialsNotValidMessage,
		))
		err = errors.New(UnauthorizedRequestErrorMessage)
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonHttpRequestExecutionFailed,
			ConditionReasonHttpRequestExecutionFailedMessage,
		))
		return err
	}

	// 6. Evaluate the condition to execute an action
	// Ref: https://gjson.dev/
	if !gjson.Valid(string(responseBody)) {
		err = errors.New(InvalidJsonErrorMessage)
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonHttpResponseNotValid,
			ConditionReasonHttpResponseNotValidMessage,
		))
		return err
	}

	//
	parsedValue := gjson.Get(string(responseBody), workloadActionManifest.Spec.Condition.Key)
	if parsedValue.String() != workloadActionManifest.Spec.Condition.Value {
		return err
	}

	// 7. Condition is met. Execute the action
	switch workloadActionManifest.Spec.Action {
	case ActionRestart:
		err = r.SetWorkloadRestartAnnotation(ctx, targetObject)
	case ActionDelete:
		err = errors.New(DeleteActionNotImplementedErrorMessage)
	default:
		err = errors.New(InvalidActionErrorMessage)
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonInvalidAction,
			ConditionReasonInvalidActionMessage,
		))
		return err
	}

	if err != nil {
		r.UpdateWorkloadActionCondition(workloadActionManifest, r.NewWorkloadActionCondition(ConditionTypeWorkloadActionReady,
			metav1.ConditionFalse,
			ConditionReasonActionExecutionFailed,
			ConditionReasonActionExecutionFailedMessage,
		))
	}

	return err
}