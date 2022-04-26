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
	"context"
	"os"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Returns the labels associated with the resource being created
func LabelsForMetadata(name string) map[string]string {
	return map[string]string{"app.kubernetes.io/instance": "ibm-commonui-operator",
		"app.kubernetes.io/name": name, "app.kubernetes.io/managed-by": "ibm-commonui-operator", "intent": "projected"}
}

// Returns the labels for selecting the resources belonging to the given CR
func LabelsForSelector(name string, crType string, crName string) map[string]string {
	return map[string]string{"k8s-app": name, crType: crName}
}

// Returns the labels associated with the Pod being created
func LabelsForPodMetadata(name string, crType string, crName string) map[string]string {
	metaLabels := LabelsForMetadata(name)
	selectorLabels := LabelsForSelector(name, crType, crName)
	for key, value := range selectorLabels {
		metaLabels[key] = value
	}
	return metaLabels
}

// Constructs image IDs for operands: either <IMAGE_NAME>:<IMAGE_TAG> or <IMAGE_NAME>@<IMAGE_SHA>
func GetImageID(imageRegistry, imageName, defaultImageVersion, imagePostfix, envVarName string) string {
	var imageID string

	// Check if the env var exists. If it exists, use that image id; otherwise, use the default image version.
	imageValue := os.Getenv(envVarName)

	if len(imageValue) > 0 {
		imageID = imageValue
	} else {
		imageSuffix := ":" + defaultImageVersion
		if imagePostfix != "" {
			imageSuffix += imagePostfix
		}
		imageID = imageRegistry + "/" + imageName + imageSuffix
	}

	return imageID
}

// returns a bool after checking for a zen instance in cs namespace
func IsAdminHubOnZen(ctx context.Context, client client.Client, namespace string) bool {
	reqLogger := log.WithValues("func", "adminHubOnZen")
	reqLogger.Info("Checking zen optional install condition")

	zenDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zen-core",
			Namespace: namespace,
		},
	}
	getError := client.Get(ctx, types.NamespacedName{Name: "zen-core", Namespace: namespace}, zenDeployment)

	if getError == nil {
		reqLogger.Info("Got ZEN Deployment")
		return true
	}
	if errors.IsNotFound(getError) {
		reqLogger.Info("ZEN deployment not found")
	} else {
		reqLogger.Error(getError, "Error getting ZEN deployment")
	}
	return false
}

// returns kubernetes cluster type
func GetKubernetesClusterType(ctx context.Context, client client.Client, namespace string) bool {
	reqLogger := log.WithValues("func", "isCncf")
	reqLogger.Info("Checking kubernetes cluster type")

	ibmProjectK := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibm-cpp-config",
			Namespace: namespace,
		},
	}
	getError := client.Get(ctx, types.NamespacedName{Name: "ibm-cpp-config", Namespace: namespace}, ibmProjectK)

	if getError == nil {
		reqLogger.Info("Got ibm-cpp-config configmap")
		clusterType := ibmProjectK.Data["kubernetes_cluster_type"]
		reqLogger.Info("Kubernetes cluster type is " + clusterType)
		if clusterType == "cncf" {
			return true
		}
	} else {
		reqLogger.Error(getError, "error getting ibm-cpp-config configmap in cs namepace")
	}

	return false
}

// Returns the int64 representation of a resource string if properly formatted. Otherwise, returns the given default value.
func GetResourceLimitsWithDefault(valueStr string, defaultValue int64) int64 {
	value := defaultValue

	if valueStr != "" {
		limit, errLim := strconv.ParseInt(valueStr[0:len(valueStr)-1], 10, 64)
		if errLim == nil {
			value = limit
		}
	}

	return value
}

// Returns the int64 representation of a resource string if properly formatted. Otherwise, returns the given default value.
func GetResourceMemoryWithDefault(valueStr string, defaultValue int64) int64 {
	value := defaultValue

	if valueStr != "" {
		memory, errLim := strconv.ParseInt(valueStr[0:len(valueStr)-2], 10, 64)
		if errLim == nil {
			value = memory
		}
	}

	return value
}

// Returns the given string if is not empty. Otherwise, returns default string.
func GetStringWithDefault(str, defaultStr string) string {
	value := str

	if value == "" {
		value = defaultStr
	}

	return value
}

func PreserveKeyValue(key string, src, dest map[string]string) {
	if val, ok := src[key]; ok {
		dest[key] = val
	}
}

func ContainsString(strs []string, search string) bool {
	for _, item := range strs {
		if item == search {
			return true
		}
	}
	return false
}

func RemoveString(strs []string, search string) []string {
	result := []string{}
	for _, item := range strs {
		if item == search {
			continue
		}
		result = append(result, item)
	}
	return result
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}
