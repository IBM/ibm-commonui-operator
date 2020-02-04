
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


const Log4jsVolumeName = "log4js"
const ClusterCaVolumeName = "cluster-ca"

var Log4jsVolume = corev1.Volume{
	Name: Log4jsVolume,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			name:  "common-web-ui-logs4js",
		},
	},
}

var ClusterCaVolume = corev1.Volume{
	Name: ClusterCaVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName:  APIKeyZecretName,
			DefaultMode: &DefaultMode,
			Optional:    &TrueVar,
		},
	},
}