//go:build !ignore_autogenerated

/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CNIConfSpec) DeepCopyInto(out *CNIConfSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CNIConfSpec.
func (in *CNIConfSpec) DeepCopy() *CNIConfSpec {
	if in == nil {
		return nil
	}
	out := new(CNIConfSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUMemSpec) DeepCopyInto(out *CPUMemSpec) {
	*out = *in
	out.CPU = in.CPU.DeepCopy()
	out.Memory = in.Memory.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUMemSpec.
func (in *CPUMemSpec) DeepCopy() *CPUMemSpec {
	if in == nil {
		return nil
	}
	out := new(CPUMemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentSpec) DeepCopyInto(out *ComponentSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentSpec.
func (in *ComponentSpec) DeepCopy() *ComponentSpec {
	if in == nil {
		return nil
	}
	out := new(ComponentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Configuration) DeepCopyInto(out *Configuration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Configuration.
func (in *Configuration) DeepCopy() *Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Configuration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationList) DeepCopyInto(out *ConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Configuration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationList.
func (in *ConfigurationList) DeepCopy() *ConfigurationList {
	if in == nil {
		return nil
	}
	out := new(ConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationSpec) DeepCopyInto(out *ConfigurationSpec) {
	*out = *in
	in.Global.DeepCopyInto(&out.Global)
	out.Networking = in.Networking
	out.Component = in.Component
	out.IPv4 = in.IPv4
	out.IPv6 = in.IPv6
	out.DualStackSpec = in.DualStackSpec
	out.Performance = in.Performance
	out.Debug = in.Debug
	out.CNIConf = in.CNIConf
	out.KubeletConfig = in.KubeletConfig
	out.LogConfig = in.LogConfig
	out.HugePages = in.HugePages.DeepCopy()
	out.DPDKCPU = in.DPDKCPU.DeepCopy()
	out.DPDKMemory = in.DPDKMemory.DeepCopy()
	in.OVNCentral.DeepCopyInto(&out.OVNCentral)
	in.OVSOVN.DeepCopyInto(&out.OVSOVN)
	in.KubeOVNController.DeepCopyInto(&out.KubeOVNController)
	in.KubeOVNCNI.DeepCopyInto(&out.KubeOVNCNI)
	in.KubeOVNPinger.DeepCopyInto(&out.KubeOVNPinger)
	in.KubeOVNMonitor.DeepCopyInto(&out.KubeOVNMonitor)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationSpec.
func (in *ConfigurationSpec) DeepCopy() *ConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(ConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigurationStatus) DeepCopyInto(out *ConfigurationStatus) {
	*out = *in
	if in.MatchingNodeAddresses != nil {
		in, out := &in.MatchingNodeAddresses, &out.MatchingNodeAddresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigurationStatus.
func (in *ConfigurationStatus) DeepCopy() *ConfigurationStatus {
	if in == nil {
		return nil
	}
	out := new(ConfigurationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DebugSpec) DeepCopyInto(out *DebugSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DebugSpec.
func (in *DebugSpec) DeepCopy() *DebugSpec {
	if in == nil {
		return nil
	}
	out := new(DebugSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalSpec) DeepCopyInto(out *GlobalSpec) {
	*out = *in
	in.Registry.DeepCopyInto(&out.Registry)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalSpec.
func (in *GlobalSpec) DeepCopy() *GlobalSpec {
	if in == nil {
		return nil
	}
	out := new(GlobalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageDetails) DeepCopyInto(out *ImageDetails) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageDetails.
func (in *ImageDetails) DeepCopy() *ImageDetails {
	if in == nil {
		return nil
	}
	out := new(ImageDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletConfigSpec) DeepCopyInto(out *KubeletConfigSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletConfigSpec.
func (in *KubeletConfigSpec) DeepCopy() *KubeletConfigSpec {
	if in == nil {
		return nil
	}
	out := new(KubeletConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogConfigSpec) DeepCopyInto(out *LogConfigSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogConfigSpec.
func (in *LogConfigSpec) DeepCopy() *LogConfigSpec {
	if in == nil {
		return nil
	}
	out := new(LogConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkStackSpec) DeepCopyInto(out *NetworkStackSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkStackSpec.
func (in *NetworkStackSpec) DeepCopy() *NetworkStackSpec {
	if in == nil {
		return nil
	}
	out := new(NetworkStackSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkingSpec) DeepCopyInto(out *NetworkingSpec) {
	*out = *in
	out.Vlan = in.Vlan
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkingSpec.
func (in *NetworkingSpec) DeepCopy() *NetworkingSpec {
	if in == nil {
		return nil
	}
	out := new(NetworkingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PerformanceSpec) DeepCopyInto(out *PerformanceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PerformanceSpec.
func (in *PerformanceSpec) DeepCopy() *PerformanceSpec {
	if in == nil {
		return nil
	}
	out := new(PerformanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistrySpec) DeepCopyInto(out *RegistrySpec) {
	*out = *in
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Images = in.Images
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistrySpec.
func (in *RegistrySpec) DeepCopy() *RegistrySpec {
	if in == nil {
		return nil
	}
	out := new(RegistrySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSpec) DeepCopyInto(out *ResourceSpec) {
	*out = *in
	in.Requests.DeepCopyInto(&out.Requests)
	in.Limits.DeepCopyInto(&out.Limits)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSpec.
func (in *ResourceSpec) DeepCopy() *ResourceSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VlanSpec) DeepCopyInto(out *VlanSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VlanSpec.
func (in *VlanSpec) DeepCopy() *VlanSpec {
	if in == nil {
		return nil
	}
	out := new(VlanSpec)
	in.DeepCopyInto(out)
	return out
}
