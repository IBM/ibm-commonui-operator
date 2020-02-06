package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CommonWebUIServiceSpec defines the desired state of CommonWebUIServiceSpec
// +k8s:openapi-gen=true
type CommonWebUIServiceSpec struct {
	CommonUIConfig  CommonWebUIConfig       `json:"uiconfig,omitempty"`
	GlobalUIConfig  CommonWebUIGlobalConfig `json:"globalConfig,omitempty"`
	OperatorVersion string                  `json:"operatorVersion,omitempty"`
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

// CommonWebUIGlobalConfig defines the desired state of CommonWebUIGlobalConfig
// +k8s:openapi-gen=true
type CommonWebUIGlobalConfig struct {
	PullSecret               string `json:"pullSecret,omitempty"`
	RouterURL                string `json:"cfcRouterUrl,omitempty"`
	IdentityProviderURL      string `json:"identityProviderUrl,omitempty"`
	AuthServiceURL           string `json:"authServiceUrl,omitempty"`
	CloudPakVersion          string `json:"CLOUDPAK_VERSION,omitempty"`
	DefaultAdminUser         string `json:"default_admin_user,omitempty"`
	RouterHTTPSPort          int32 `json:"router_https_port,omitempty"`
	ClusterName              string `json:"cluster_name,omitempty"`
	SessionPollingInterval   int32 `json:"session_polling_interval,omitempty"`
}

// CommonWebUIServiceStatus defines the observed state of CommonWebUIService
// +k8s:openapi-gen=true
type CommonWebUIServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CommonWebUIService is the Schema for the commonwebuiservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=commonwebuiservices,scope=Namespaced
type CommonWebUIService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommonWebUIServiceSpec   `json:"spec,omitempty"`
	Status CommonWebUIServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CommonWebUIServiceList contains a list of CommonWebUIService
type CommonWebUIServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommonWebUIService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CommonWebUIService{}, &CommonWebUIServiceList{})
}
