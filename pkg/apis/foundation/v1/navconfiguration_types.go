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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NavConfigurationSpec defines the desired state of NavConfiguration
type NavConfigurationSpec struct {
	LogoutRedirects []string   `json:"logoutRedirects,omitempty"`
	About           About      `json:"about,omitempty"`
	Header          Header     `json:"header,omitempty"`
	Login           Login      `json:"login,omitempty"`
	NavItems        []NavItems `json:"navItems,omitempty"`
	OperatorVersion string     `json:"operatorVersion,omitempty"`
	Version         string     `json:"version,omitempty"`
	License         License    `json:"license,omitempty"`
}

// NavConfigurationStatus defines the observed state of NavConfiguration
type NavConfigurationStatus struct {
	Versions Versions `json:"versions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NavConfiguration is the Schema for the navconfigurations API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=navconfigurations,scope=Namespaced
type NavConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NavConfigurationSpec   `json:"spec,omitempty"`
	Status NavConfigurationStatus `json:"status,omitempty"`
}

type License struct {
	Accept bool `json:"accept,omitempty"`
}

type Versions struct {
	Reconciled string `json:"reconciled,omitempty"`
}

type About struct {
	LogoURL   string   `json:"logoUrl,omitempty"`
	Licenses  []string `json:"licenses,omitempty"`
	Copyright string   `json:"copyright,omitempty"`
	Version   string   `json:"version,omitempty"`
	Edition   string   `json:"edition,omitempty"`
}

type Header struct {
	LogoURL           string            `json:"logoUrl,omitempty"`
	LogoWidth         string            `json:"logoWidth,omitempty"`
	LogoHeight        string            `json:"logoHeight,omitempty"`
	LogoAltText       string            `json:"logoAltText,omitempty"`
	DocURLMapping     string            `json:"docUrlMapping,omitempty"`
	DisabledItems     []string          `json:"disabledItems,omitempty"`
	DetectHeaderItems DetectHeaderItems `json:"detectHeaderItems,omitempty"`
}

type DetectHeaderItems struct {
	AdditionalProperties AdditionalProperties `json:"additionalProperties,omitempty"`
}

type AdditionalProperties struct {
	DetectionNamespace     string   `json:"detectionNamespace,omitempty"`
	DetectionServiceName   string   `json:"detectionServiceName,omitempty"`
	DetectionLabelSelector string   `json:"detectionLabelSelector,omitempty"`
	IsAuthorized           []string `json:"isAuthorized,omitempty"`
}

type Login struct {
	LogoAltText string      `json:"logoAltText,omitempty"`
	LogoURL     string      `json:"logoUrl,omitempty"`
	LogoWidth   string      `json:"logoWidth,omitempty"`
	LogoHeight  string      `json:"logoHeight,omitempty"`
	LoginDialog LoginDialog `json:"loginDialog,omitempty"`
}

type LoginDialog struct {
	Enable     bool   `json:"enable,omitempty"`
	HeaderText string `json:"headerText,omitempty"`
	DialogText string `json:"dialogText,omitempty"`
	AcceptText string `json:"acceptText,omitempty"`
}

type NavItems struct {
	ID                   string   `json:"id,omitempty"`
	Label                string   `json:"label,omitempty"`
	URL                  string   `json:"url,omitempty"`
	IconURL              string   `json:"iconUrl,omitempty"`
	Target               string   `json:"target,omitempty"`
	ParentID             string   `json:"parentId,omitempty"`
	Namespace            string   `json:"namespace,omitempty"`
	ServiceName          string   `json:"serviceName,omitempty"`
	ServiceID            string   `json:"serviceId,omitempty"`
	DetectionServiceName bool     `json:"detectionServiceName,omitempty"`
	IsAuthorized         []string `json:"isAuthorized,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NavConfigurationList contains a list of NavConfiguration
type NavConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NavConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NavConfiguration{}, &NavConfigurationList{})
}
