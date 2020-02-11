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

// LegacyHeaderServiceSpec defines the desired state of LegacyHeaderService
// +k8s:openapi-gen=true
type LegacyHeaderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	LegacyConfig         LegacyConfig         `json:"legacyConfig,omitempty"`
	LegacyGlobalUIConfig LegacyGlobalUIConfig `json:"globalUIConfig,omitempty"`
	OperatorVersion      string               `json:"operatorVersion,omitempty"`
}

// LegacyConfig defines the desired state of LegacyConfig
// +k8s:openapi-gen=true
type LegacyConfig struct {
	ServiceName       string `json:"serviceName,omitempty"`
	ImageRegistry     string `json:"imageRegistry,omitempty"`
	ImageTag          string `json:"imageTag,omitempty"`
	CPULimits         string `json:"cpuLimits,omitempty"`
	CPUMemory         string `json:"cpuMemory,omitempty"`
	RequestLimits     string `json:"requestLimits,omitempty"`
	RequestMemory     string `json:"requestMemory,omitempty"`
	LegacyLogoPath    string `json:"legacyLogoPath,omitempty"`
	LegacyLogoWidth   string `json:"legacyLogoWidth,omitempty"`
	LegacyLogoHeight  string `json:"legacyLogoHeight,omitempty"`
	LegacySupportURL  string `json:"legacySupportURL,omitempty"`
	LegacyDocURL      string `json:"legacyDocURL,omitempty"`
	LegacyLogoAltText string `json:"legacyLogoAltText,omitempty"`
	IngressPath       string `json:"ingressPath,omitempty"`
}

// GlobalUIConfig defines the desired state of GlobalUIConfig
// +k8s:openapi-gen=true
type LegacyGlobalUIConfig struct {
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

// LegacyHeaderStatus defines the observed state of LegacyHeaderService
// +k8s:openapi-gen=true
type LegacyHeaderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LegacyHeader is the Schema for the legacyHeader API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=legacyheaders,scope=Namespaced
type LegacyHeader struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LegacyHeaderSpec   `json:"spec,omitempty"`
	Status LegacyHeaderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LegacyHeaderList contains a list of LegacyHeader
type LegacyHeaderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LegacyHeader `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LegacyHeader{}, &LegacyHeaderList{})
}
