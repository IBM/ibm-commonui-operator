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

package resources

// import (

// 	corev1 "k8s.io/api/core/v1"
// 	netv1 "k8s.io/api/networking/v1beta1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/util/intstr"
// 	logf "sigs.k8s.io/controller-runtime/pkg/log"
// )

const ReleaseName = "common-web-ui"
const DaemonSetName = "common-web-ui"
const ChartName = "webui-nav"
const ChartVersion = "1.0.2"

// returns the labels associated with the resource being created
func LabelsForMetadata(deploymentName string) map[string]string {
	return map[string]string{"app": deploymentName, "chart": ChartName , "version": ChartVersion,
		"heritage": "operator", "release": ReleaseName}
}

// returns the labels for selecting the resources belonging to the given metering CR name
func LabelsForSelector(deploymentName string, crType string, crName string) map[string]string {
	return map[string]string{"k8s-app": deploymentName, crType: crName}
}

// returns the labels associated with the Pod being created
func LabelsForPodMetadata(deploymentName string, crType string, crName string) map[string]string {
	podLabels := LabelsForMetadata(deploymentName)
	selectorLabels := LabelsForSelector(deploymentName, crType, crName)
	for key, value := range selectorLabels {
		podLabels[key] = value
	}
	return podLabels
}