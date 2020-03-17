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

// CS??? removed icp-serviceid-apikey-secret from CommonSecretCheckNames, CommonSecretCheckDirs,
// CS???   and CommonSecretCheckVolumeMounts
// Linter doesn't like "Secret" in string var names so use "Zecret"

package resources

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const DefaultImageRegistry = "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom-amd64"
const DefaultImageName = "common-web-ui"
const DefaultImageTag = "1.1.0"

const LegacyImageRegistry = "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom-amd64"
const LegacyImageName = "icp-platform-header"
const LegacyImageTag = "3.2.4"

var TrueVar = true
var FalseVar = false
var Replica1 int32 = 1
var Seconds60 int64 = 60

var cpu300 = resource.NewMilliQuantity(300, resource.DecimalSI)        // 300m
var memory256 = resource.NewQuantity(256*1024*1024, resource.BinarySI) // 256Mi

const Log4jsVolumeName = "log4js"
const ClusterCaVolumeName = "cluster-ca"

var Log4jsVolume = corev1.Volume{
	Name: Log4jsVolumeName,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "common-web-ui-log4js",
			},
			Items: []corev1.KeyToPath{
				{
					Key:  "log4js.json",
					Path: "log4js.json",
				},
			},
			// DefaultMode: &DefaultMode,
			Optional: &TrueVar,
		},
	},
}

var ClusterCaVolume = corev1.Volume{
	Name: ClusterCaVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName: "cluster-ca-cert",
			Items: []corev1.KeyToPath{
				{
					Key:  "tls.key",
					Path: "ca.key",
				},
				{
					Key:  "tls.crt",
					Path: "ca.crt",
				},
			},
			// DefaultMode: &DefaultMode,
			Optional: &TrueVar,
		},
	},
}

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
}

var CommonContainer = corev1.Container{
	Image: "common-web-ui",
	//CS??? Image: "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom-amd64/metering-data-manager:3.3.1",
	Name:            "common-web-ui",
	ImagePullPolicy: corev1.PullIfNotPresent,

	Resources: corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    *cpu300,
			corev1.ResourceMemory: *memory256},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    *cpu300,
			corev1.ResourceMemory: *memory256},
	},

	SecurityContext: &commonSecurityContext,

	ReadinessProbe: &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/readinessProbe",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 3000,
				},
				Scheme: "",
			},
		},
		InitialDelaySeconds: 100,
		TimeoutSeconds:      15,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	},

	LivenessProbe: &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/livenessProbe",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 3000,
				},
				Scheme: "",
			},
		},
		InitialDelaySeconds: 100,
		TimeoutSeconds:      5,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	},

	// CommonMainVolumeMounts will be added by the controller
	VolumeMounts: []corev1.VolumeMount{
		{
			Name:      Log4jsVolumeName,
			MountPath: "/etc/config",
		},
		{
			Name:      ClusterCaVolumeName,
			MountPath: "/opt/ibm/platform-header/certs",
		},
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
			Value: "https://icp-management-ingress:443/idprovider",
		},
		{
			Name:  "PLATFORM_AUTH_SERVICE_URL",
			Value: "https://icp-management-ingress:443/idauth",
		},
		{
			Name:  "NAV_PORT",
			Value: "8443",
		},
		{
			Name: "OAUTH2_CLIENT_REGISTRATION_SECRET",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "platform-oidc-credentials",
					},
					Key: "OAUTH2_CLIENT_REGISTRATION_SECRET",
				},
			},
		},
		{
			Name:  "CLOUDPAK_VERSION",
			Value: "1.0.0",
		},
		{
			Name:  "default_admin_user",
			Value: "admin",
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
			Name: "ROKS_ENABLED",
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "platform-auth-idp",
					},
					Key: "ROKS_ENABLED",
				},
			},
		},
		{
			Name: "WLP_CLIENT_ID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "platform-oidc-credentials",
					},
					Key: "WLP_CLIENT_ID",
				},
			},
		},
		{
			Name: "WLP_CLIENT_SECRET",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "platform-oidc-credentials",
					},
					Key: "WLP_CLIENT_SECRET",
				},
			},
		},
	},
}
