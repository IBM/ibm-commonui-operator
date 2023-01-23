//
// Copyright 2021 IBM Corporation
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
	"encoding/json"
	errorf "errors"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

func ReconcileAdminHubNavConfig(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI) error {
	reqLogger := log.WithValues("func", "reconcileAdminHubNavConfig", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling admin hub nav config")

	return reconcileNavConfig(ctx, client, instance, AdminHubNavConfigName, AdminHubNavConfig)
}

func reconcileNavConfig(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, name, config string) error {
	reqLogger := log.WithValues("func", "reconcileNavConfig", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	var template map[string]interface{}
	err := json.Unmarshal([]byte(config), &template)
	if err != nil {
		reqLogger.Info(fmt.Sprintf("Failed to unmarshal nav config: %s", name))
		return err
	}

	navConfig := &unstructured.Unstructured{
		Object: template,
	}

	desiredNavConfig := &unstructured.Unstructured{
		Object: CopyMap(template),
	}

	err = client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, navConfig)
	if err == nil {
		reqLogger.Info(fmt.Sprintf("Updating nav config: %s", name))

		// cast navItems interface to array of json objects
		navItemsValue, ok := desiredNavConfig.Object["spec"].(map[string]interface{})["navItems"]
		if !ok {
			msg := fmt.Sprintf("Failed to unmarshal navItems array for nav config: %s", name)
			reqLogger.Info(msg)
			return errorf.New(msg)
		}

		// get reflected value of navItems interface
		s := reflect.ValueOf(navItemsValue)

		// return error if navItems interface is not a slice
		if s.IsNil() || s.Kind() != reflect.Slice {
			msg := fmt.Sprintf("Invalid navItems array in nav config: %s", name)
			reqLogger.Info(msg)
			return errorf.New(msg)
		}

		// initialize navItems array based on length of value
		navItems := make([]map[string]interface{}, s.Len())

		// cast each navItem to a map of strings to interfaces
		for i := 0; i < s.Len(); i++ {
			navItems[i], ok = s.Index(i).Interface().(map[string]interface{})
			if !ok {
				msg := fmt.Sprintf("Failed to unmarshal navItem %d for nav config: %s", i, name)
				reqLogger.Info(msg)
				return errorf.New(msg)
			}
		}

		// Update namespace for all nav items
		for _, item := range navItems {
			if item["namespace"] != "" {
				item["namespace"] = instance.Namespace
			}
		}

		// Set nav items to array with updated namespaces
		navConfig.Object["spec"].(map[string]interface{})["navItems"] = navItems

		// Update with latest licenses
		if name == AdminHubNavConfigName {
			licenses := desiredNavConfig.Object["spec"].(map[string]interface{})["about"].(map[string]interface{})["licenses"]
			navConfig.Object["spec"].(map[string]interface{})["about"].(map[string]interface{})["licenses"] = licenses
		}

		// Update nav config
		err = client.Update(ctx, navConfig)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to update nav config: %s", name), "NavConfig.Namespace", instance.Namespace, "NavConfig.Name", name)
			return err
		}
	}

	return nil
}
