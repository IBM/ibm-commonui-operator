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

// SwitcherItemSpec defines the desired state of SwitcherItem
type SwitcherItemSpec struct {
	CloudPakInfo    CloudPakInfo `json:"cloudPakInfo,omitempty"`
	OperatorVersion string       `json:"operatorVersion,omitempty"`
	Version         string       `json:"version,omitempty"`
	License         License      `json:"license"`
}

type CloudPakInfo struct {
	LogoURL     string `json:"logoURL,omitempty"`
	Label       string `json:"label,omitempty"`
	Display     string `json:"display,omitempty"`
	LandingPage string `json:"landingPage,omitempty"`
}

// SwitcherItemStatus defines the observed state of SwitcherItem
type SwitcherItemStatus struct {
	Versions Versions `json:"versions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwitcherItem is the Schema for the switcheritems API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=switcheritems,scope=Namespaced
type SwitcherItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitcherItemSpec   `json:"spec,omitempty"`
	Status SwitcherItemStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwitcherItemList contains a list of SwitcherItem
type SwitcherItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitcherItem `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitcherItem{}, &SwitcherItemList{})
}
