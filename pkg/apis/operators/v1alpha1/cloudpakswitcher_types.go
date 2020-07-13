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

// CloudPakSwitcherSpec defines the desired state of CloudPakSwitcher
type CloudPakSwitcherSpec struct {
	CloudPakInfo    CloudPakInfo `json:"cloudPakInfo,omitempty"`
	OperatorVersion string       `json:"operatorVersion,omitempty"`
	Version         string       `json:"version,omitempty"`
}

type CloudPakInfo struct {
	LogoURL     string     `json:"logoURL,omitempty"`
	Label       string     `json:"label,omitempty"`
	Display     string     `json:"display,omitempty"`
	LandingPage string     `json:"landingPage,omitempty"`
	OtherLinks  OtherLinks `json:"otherLinks,omitempty"`
}

type OtherLinks struct {
	Label   string `json:"label,omitempty"`
	Display string `json:"display,omitempty"`
	URL     string `json:"url,omitempty"`
	Target  string `json:"target,omitempty"`
}

// CloudPakSwitcherStatus defines the observed state of CloudPakSwitcher
type CloudPakSwitcherStatus struct {
	// PodNames will hold the names of the commonwebui's
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudPakSwitcher is the Schema for the cloudpakswitchers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=cloudpakswitchers,scope=Namespaced
type CloudPakSwitcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudPakSwitcherSpec   `json:"spec,omitempty"`
	Status CloudPakSwitcherStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudPakSwitcherList contains a list of CloudPakSwitcher
type CloudPakSwitcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudPakSwitcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudPakSwitcher{}, &CloudPakSwitcherList{})
}
