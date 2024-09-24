/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CommonWebUIConfig struct {
	ServiceName   string `json:"serviceName,omitempty"`
	ImageRegistry string `json:"imageRegistry,omitempty"`
	ImageTag      string `json:"imageTag,omitempty"`
	CPULimits     string `json:"cpuLimits,omitempty"`
	CPUMemory     string `json:"cpuMemory,omitempty"`
	RequestLimits string `json:"requestLimits,omitempty"`
	RequestMemory string `json:"requestMemory,omitempty"`
	IngressPath   string `json:"ingressPath,omitempty"`
	LandingPage   string `json:"landingPage,omitempty"`
}

type GlobalUIConfig struct {
	PullSecret             string `json:"pullSecret,omitempty"`
	CloudPakVersion        string `json:"cloudPakVersion,omitempty"`
	DefaultAdminUser       string `json:"defaultAdminUser,omitempty"`
	DefaultAuth            string `json:"defaultAuth,omitempty"`
	OSAuth                 string `json:"osAuth,omitempty"`
	EnterpriseLDAP         string `json:"enterpriseLDAP,omitempty"`
	EnterpriseSAML         string `json:"enterpriseSAML,omitempty"`
	SessionPollingInterval int32  `json:"sessionPollingInterval,omitempty"`
}

type Resources struct {
	Requests Requests `json:"requests,omitempty"`
	Limits   Limits   `json:"limits,omitempty"`
}

type Requests struct {
	RequestLimits    string `json:"cpu,omitempty"`
	RequestMemory    string `json:"memory,omitempty"`
	EphemeralStorage string `json:"ephemeral-storage,omitempty"`
}

type Limits struct {
	CPULimits        string `json:"cpu,omitempty"`
	CPUMemory        string `json:"memory,omitempty"`
	EphemeralStorage string `json:"ephemeral-storage,omitempty"`
}

// CommonWebUISpec defines the desired state of CommonWebUI
type CommonWebUISpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CommonWebUIConfig             CommonWebUIConfig `json:"commonWebUIConfig,omitempty"`
	GlobalUIConfig                GlobalUIConfig    `json:"globalUIConfig,omitempty"`
	OperatorVersion               string            `json:"operatorVersion,omitempty"`
	Version                       string            `json:"version,omitempty"`
	Replicas                      int32             `json:"replicas,omitempty"`
	Resources                     Resources         `json:"resources,omitempty"`
	Labels                        map[string]string `json:"labels,omitempty"`
	EnableInstanaMetricCollection bool              `json:"enableInstanaMetricCollection,omitempty"`
	// License           License           `json:"license,omitempty"`
}

// CommonWebUIStatus defines the observed state of CommonWebUI
type CommonWebUIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// PodNames will hold the names of the commonwebui's
	Nodes []string `json:"nodes"`
	// Versions Versions `json:"versions,omitempty"`
	Service         ServiceStatus `json:"service,omitempty"`
	OperatorVersion string        `json:"operatorVersion,omitempty"`
	OperandVersion  string        `json:"operandVersion,omitempty"`
}

// ServiceStatus struct
type ServiceStatus struct {
	ObjectName       string                  `json:"objectName,omitempty"`
	APIVersion       string                  `json:"apiVersion,omitempty"`
	Namespace        string                  `json:"namespace,omitempty"`
	Kind             string                  `json:"kind,omitempty"`
	Status           string                  `json:"status,omitempty"`
	ManagedResources []ManagedResourceStatus `json:"managedResources,omitempty"`
}

type ManagedResourceStatus struct {
	ObjectName string `json:"objectName,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Status     string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CommonWebUI is the Schema for the commonwebuis API
type CommonWebUI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommonWebUISpec   `json:"spec,omitempty"`
	Status CommonWebUIStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CommonWebUIList contains a list of CommonWebUI
type CommonWebUIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommonWebUI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CommonWebUI{}, &CommonWebUIList{})
}
