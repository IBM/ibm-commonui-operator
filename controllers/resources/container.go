//
// Copyright 2022 IBM Corporation
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

package resources

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var commonSecurityContext = corev1.SecurityContext{
	AllowPrivilegeEscalation: &FalseVar,
	Privileged:               &FalseVar,
	ReadOnlyRootFilesystem:   &TrueVar,
	RunAsNonRoot:             &TrueVar,
	Capabilities: &corev1.Capabilities{
		Drop: []corev1.Capability{
			"ALL",
		},
	},
	SeccompProfile: &corev1.SeccompProfile{
		Type: corev1.SeccompProfileTypeRuntimeDefault,
	},
}

var CommonContainer = corev1.Container{
	Image:           "common-web-ui",
	Name:            "common-web-ui",
	ImagePullPolicy: corev1.PullIfNotPresent,

	Resources: corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    *cpu300,
			corev1.ResourceMemory: *memory256,
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:              *cpu300,
			corev1.ResourceMemory:           *memory256,
			corev1.ResourceEphemeralStorage: *memory251,
		},
	},

	SecurityContext: &commonSecurityContext,

	ReadinessProbe: &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/readinessProbe",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 3000,
				},
				Scheme: corev1.URISchemeHTTPS,
			},
		},
		InitialDelaySeconds: 30,
		TimeoutSeconds:      15,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	},

	LivenessProbe: &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/livenessProbe",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 3000,
				},
				Scheme: corev1.URISchemeHTTPS,
			},
		},
		InitialDelaySeconds: 30,
		TimeoutSeconds:      5,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	},

	// CommonEnvVars will be added by the controller
	Env: []corev1.EnvVar{
		{
			Name:  "contextPath",
			Value: "/common-nav",
		},
		{
			Name:  "cfcRouterUrl",
			Value: "https://icp-management-ingress:443",
		},
		{
			Name:  "NODE_EXTRA_CA_CERTS",
			Value: " /opt/ibm/platform-header/certs/ca.crt",
		},
		{
			Name:  "PLATFORM_IDENTITY_PROVIDER_URL",
			Value: "https://platform-identity-provider:4300",
		},
		{
			Name:  "PLATFORM_AUTH_SERVICE_URL",
			Value: "https://platform-auth-service:9443",
		},
		{
			Name:  "NAV_PORT",
			Value: "8443",
		},
		{
			Name:  "CLOUDPAK_VERSION",
			Value: "1.0.0",
		},
		{
			Name:  "CLUSTER_NAME",
			Value: "mycluster",
		},
		{
			Name:  "defaultAuth",
			Value: "",
		},
		{
			Name:  "enterpriseLDAP",
			Value: "",
		},
		{
			Name:  "enterpriseSAML",
			Value: "",
		},
		{
			Name:  "osAuth",
			Value: "",
		},
		{
			Name:  "SESSION_POLLING_INTERVAL",
			Value: "300",
		},
		{
			Name:  "PREFERRED_LOGIN",
			Value: "deprecated",
		},
		{
			Name:  "ROKS_ENABLED",
			Value: "deprecated",
		},
		{
			Name:  "USE_HTTPS",
			Value: "true",
		},
		{
			Name:  "UI_SSL_CA",
			Value: "/certs/common-web-ui/ca.crt",
		},
		{
			Name:  "UI_SSL_CERT",
			Value: "/certs/common-web-ui/tls.crt",
		},
		{
			Name:  "UI_SSL_KEY",
			Value: "/certs/common-web-ui/tls.key",
		},
		{
			Name:  "LANDING_PAGE",
			Value: "",
		},
		{
			Name:  "WATCH_NAMESPACE",
			Value: "",
		},
		{
			Name: "POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.namespace",
				},
			},
		},
		{
			Name:  "USE_ZEN",
			Value: "false",
		},
		{
			Name:  "APP_VERSION",
			Value: "",
		},
		{
			Name:  "CLUSTER_TYPE",
			Value: "unknown",
		},
		{
			Name:  "OSAUTH_ENABLED",
			Value: "deprecated",
		},
		{
			Name: "INSTANA_AGENT_HOST",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "status.hostIP",
				},
			},
		},
		{
			Name:  "INSTANA_AGENT_ENABLED",
			Value: "false",
		},
	},
}
