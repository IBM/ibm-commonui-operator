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
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CommonWebUISpec defines the desired state of CommonWebUISpec
// +k8s:openapi-gen=true
type CommonWebUISpec struct {
	CommonWebUIConfig CommonWebUIConfig `json:"commonWebUIConfig,omitempty"`
	GlobalUIConfig    GlobalUIConfig    `json:"globalUIConfig,omitempty"`
	OperatorVersion   string            `json:"operatorVersion,omitempty"`
}

// CommonWebUIConfig defines the desired state of CommonWebUIConfig
// +k8s:openapi-gen=true
type CommonWebUIConfig struct {
	ServiceName   string `json:"serviceName,omitempty"`
	ImageRegistry string `json:"imageRegistry,omitempty"`
	ImageTag      string `json:"imageTag,omitempty"`
	CPULimits     string `json:"cpuLimits,omitempty"`
	CPUMemory     string `json:"cpuMemory,omitempty"`
	RequestLimits string `json:"requestLimits,omitempty"`
	RequestMemory string `json:"requestMemory,omitempty"`
	IngressPath   string `json:"ingressPath,omitempty"`
}

// GlobalUIConfig defines the desired state of GlobalUIConfig
// +k8s:openapi-gen=true
type GlobalUIConfig struct {
	PullSecret             string `json:"pullSecret,omitempty"`
	RouterURL              string `json:"routerURL,omitempty"`
	IdentityProviderURL    string `json:"identityProviderURL,omitempty"`
	AuthServiceURL         string `json:"authServiceURL,omitempty"`
	CloudPakVersion        string `json:"cloudPakVersion,omitempty"`
	DefaultAdminUser       string `json:"defaultAdminUser,omitempty"`
	RouterHTTPSPort        int32  `json:"routerHTTPSPort,omitempty"`
	ClusterName            string `json:"clusterName,omitempty"`
	SessionPollingInterval int32  `json:"sessionPollingInterval,omitempty"`
}

// CommonWebUIStatus defines the observed state of CommonWebUI
// +k8s:openapi-gen=true
type CommonWebUIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CommonWebUI is the Schema for the commonwebuis API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=commonwebuis,scope=Namespaced
type CommonWebUI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommonWebUISpec   `json:"spec,omitempty"`
	Status CommonWebUIStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CommonWebUIList contains a list of CommonWebUI
type CommonWebUIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommonWebUI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CommonWebUI{}, &CommonWebUIList{})
}
