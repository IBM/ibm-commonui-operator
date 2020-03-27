// +build !ignore_autogenerated

//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonWebUI) DeepCopyInto(out *CommonWebUI) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonWebUI.
func (in *CommonWebUI) DeepCopy() *CommonWebUI {
	if in == nil {
		return nil
	}
	out := new(CommonWebUI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CommonWebUI) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonWebUIConfig) DeepCopyInto(out *CommonWebUIConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonWebUIConfig.
func (in *CommonWebUIConfig) DeepCopy() *CommonWebUIConfig {
	if in == nil {
		return nil
	}
	out := new(CommonWebUIConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonWebUIList) DeepCopyInto(out *CommonWebUIList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CommonWebUI, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonWebUIList.
func (in *CommonWebUIList) DeepCopy() *CommonWebUIList {
	if in == nil {
		return nil
	}
	out := new(CommonWebUIList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CommonWebUIList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonWebUISpec) DeepCopyInto(out *CommonWebUISpec) {
	*out = *in
	out.CommonWebUIConfig = in.CommonWebUIConfig
	out.GlobalUIConfig = in.GlobalUIConfig
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonWebUISpec.
func (in *CommonWebUISpec) DeepCopy() *CommonWebUISpec {
	if in == nil {
		return nil
	}
	out := new(CommonWebUISpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonWebUIStatus) DeepCopyInto(out *CommonWebUIStatus) {
	*out = *in
	if in.PodNames != nil {
		in, out := &in.PodNames, &out.PodNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonWebUIStatus.
func (in *CommonWebUIStatus) DeepCopy() *CommonWebUIStatus {
	if in == nil {
		return nil
	}
	out := new(CommonWebUIStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalUIConfig) DeepCopyInto(out *GlobalUIConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalUIConfig.
func (in *GlobalUIConfig) DeepCopy() *GlobalUIConfig {
	if in == nil {
		return nil
	}
	out := new(GlobalUIConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LegacyConfig) DeepCopyInto(out *LegacyConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LegacyConfig.
func (in *LegacyConfig) DeepCopy() *LegacyConfig {
	if in == nil {
		return nil
	}
	out := new(LegacyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LegacyGlobalUIConfig) DeepCopyInto(out *LegacyGlobalUIConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LegacyGlobalUIConfig.
func (in *LegacyGlobalUIConfig) DeepCopy() *LegacyGlobalUIConfig {
	if in == nil {
		return nil
	}
	out := new(LegacyGlobalUIConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LegacyHeader) DeepCopyInto(out *LegacyHeader) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LegacyHeader.
func (in *LegacyHeader) DeepCopy() *LegacyHeader {
	if in == nil {
		return nil
	}
	out := new(LegacyHeader)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LegacyHeader) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LegacyHeaderList) DeepCopyInto(out *LegacyHeaderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LegacyHeader, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LegacyHeaderList.
func (in *LegacyHeaderList) DeepCopy() *LegacyHeaderList {
	if in == nil {
		return nil
	}
	out := new(LegacyHeaderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LegacyHeaderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LegacyHeaderSpec) DeepCopyInto(out *LegacyHeaderSpec) {
	*out = *in
	out.LegacyConfig = in.LegacyConfig
	out.LegacyGlobalUIConfig = in.LegacyGlobalUIConfig
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LegacyHeaderSpec.
func (in *LegacyHeaderSpec) DeepCopy() *LegacyHeaderSpec {
	if in == nil {
		return nil
	}
	out := new(LegacyHeaderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LegacyHeaderStatus) DeepCopyInto(out *LegacyHeaderStatus) {
	*out = *in
	if in.PodNames != nil {
		in, out := &in.PodNames, &out.PodNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LegacyHeaderStatus.
func (in *LegacyHeaderStatus) DeepCopy() *LegacyHeaderStatus {
	if in == nil {
		return nil
	}
	out := new(LegacyHeaderStatus)
	in.DeepCopyInto(out)
	return out
}
