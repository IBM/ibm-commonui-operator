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
type LegacyHeaderServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	LegacyUIConfig  CommonWebUILegacyConfig `json:"legacyuiconfig,omitempty"`
	GlobalUIConfig  CommonWebUIGlobalConfig `json:"globalConfig,omitempty"`
	OperatorVersion string                  `json:"operatorVersion,omitempty"`
}

// CommonWebUILegacyConfig defines the desired state of CommonWebUILegacyConfig
// +k8s:openapi-gen=true
type CommonWebUILegacyConfig struct {
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
	LegacySupportURL  string `json:"legacySupportUrl,omitempty"`
	LegacyDocURL      string `json:"legacyDocUrl,omitempty"`
	LegacyLogoAltText string `json:"legacyLogoAltText,omitempty"`
	IngressPath       string `json:"ingressPath,omitempty"`
}

// CommonWebUIGlobalConfig defines the desired state of CommonWebUIGlobalConfig
// +k8s:openapi-gen=true
type CommonWebUIGlobalConfig struct {
	PullSecret               string `json:"pullSecret,omitempty"`
	RouterURL                string `json:"cfcRouterUrl,omitempty"`
	IdentityProviderURL      string `json:identityProviderUrl,omitempty"`
	AuthServiceURL           string `json:"authServiceUrl,omitempty"`
	CloudPakVersion          string `json:"CLOUDPAK_VERSION,omitempty"`
	DefaultAdminUser         string `json:"default_admin_user,omitempty"`
	RouterHTTPSPort          int32 `json:"router_https_port,omitempty"`
	ClusterName              string `json:"cluster_name,omitempty"`
	SessionPollingInterval   int32 `json:"session_polling_interval,omitempty"`
}

// LegacyHeaderServiceStatus defines the observed state of LegacyHeaderService
// +k8s:openapi-gen=true
type LegacyHeaderServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LegacyHeaderService is the Schema for the legacyheaderservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=legacyheaderservices,scope=Namespaced
type LegacyHeaderService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LegacyHeaderServiceSpec   `json:"spec,omitempty"`
	Status LegacyHeaderServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LegacyHeaderServiceList contains a list of LegacyHeaderService
type LegacyHeaderServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LegacyHeaderService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LegacyHeaderService{}, &LegacyHeaderServiceList{})
}
