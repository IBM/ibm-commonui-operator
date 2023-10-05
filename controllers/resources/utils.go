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
	"os"
	"strconv"
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

func MergeMap(in map[string]string, mergeMap map[string]string) map[string]string {
	if in == nil {
		in = make(map[string]string)
	}
	for k, v := range mergeMap {
		in[k] = v
	}
	return in
}
