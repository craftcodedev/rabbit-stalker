//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConditionSpec) DeepCopyInto(out *ConditionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConditionSpec.
func (in *ConditionSpec) DeepCopy() *ConditionSpec {
	if in == nil {
		return nil
	}
	out := new(ConditionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RabbitConnectionCredentialsSecretReferenceSpec) DeepCopyInto(out *RabbitConnectionCredentialsSecretReferenceSpec) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RabbitConnectionCredentialsSecretReferenceSpec.
func (in *RabbitConnectionCredentialsSecretReferenceSpec) DeepCopy() *RabbitConnectionCredentialsSecretReferenceSpec {
	if in == nil {
		return nil
	}
	out := new(RabbitConnectionCredentialsSecretReferenceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RabbitConnectionCredentialsSpec) DeepCopyInto(out *RabbitConnectionCredentialsSpec) {
	*out = *in
	out.Username = in.Username
	out.Password = in.Password
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RabbitConnectionCredentialsSpec.
func (in *RabbitConnectionCredentialsSpec) DeepCopy() *RabbitConnectionCredentialsSpec {
	if in == nil {
		return nil
	}
	out := new(RabbitConnectionCredentialsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RabbitConnectionSpec) DeepCopyInto(out *RabbitConnectionSpec) {
	*out = *in
	out.Credentials = in.Credentials
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RabbitConnectionSpec.
func (in *RabbitConnectionSpec) DeepCopy() *RabbitConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(RabbitConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretKeyReferenceSpec) DeepCopyInto(out *SecretKeyReferenceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretKeyReferenceSpec.
func (in *SecretKeyReferenceSpec) DeepCopy() *SecretKeyReferenceSpec {
	if in == nil {
		return nil
	}
	out := new(SecretKeyReferenceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SynchronizationSpec) DeepCopyInto(out *SynchronizationSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SynchronizationSpec.
func (in *SynchronizationSpec) DeepCopy() *SynchronizationSpec {
	if in == nil {
		return nil
	}
	out := new(SynchronizationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadAction) DeepCopyInto(out *WorkloadAction) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadAction.
func (in *WorkloadAction) DeepCopy() *WorkloadAction {
	if in == nil {
		return nil
	}
	out := new(WorkloadAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkloadAction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadActionList) DeepCopyInto(out *WorkloadActionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WorkloadAction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadActionList.
func (in *WorkloadActionList) DeepCopy() *WorkloadActionList {
	if in == nil {
		return nil
	}
	out := new(WorkloadActionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkloadActionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadActionSpec) DeepCopyInto(out *WorkloadActionSpec) {
	*out = *in
	out.Synchronization = in.Synchronization
	out.RabbitConnection = in.RabbitConnection
	out.Condition = in.Condition
	out.WorkloadRef = in.WorkloadRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadActionSpec.
func (in *WorkloadActionSpec) DeepCopy() *WorkloadActionSpec {
	if in == nil {
		return nil
	}
	out := new(WorkloadActionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadActionStatus) DeepCopyInto(out *WorkloadActionStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadActionStatus.
func (in *WorkloadActionStatus) DeepCopy() *WorkloadActionStatus {
	if in == nil {
		return nil
	}
	out := new(WorkloadActionStatus)
	in.DeepCopyInto(out)
	return out
}
