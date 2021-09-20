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

package commonwebui

import (
	"context"
	"encoding/json"
	errorf "errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	operatorv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

func (r *CommonWebUIReconciler) reconcileAdminHubNavConfig(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileAdminHubNavConfig", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling admin hub nav config")

	return r.reconcileNavConfig(ctx, instance, needToRequeue, AdminHubNavConfigName, AdminHubNavConfig)
}

func (r *CommonWebUIReconciler) reconcileCP4INavConfig(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileCP4INavConfig", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling cp4i nav config")

	return r.reconcileNavConfig(ctx, instance, needToRequeue, CP4INavConfigName, CP4INavConfig)
}

func (r *CommonWebUIReconciler) reconcileNavConfig(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool, name, config string) error {
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

	err = r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, navConfig)
	if err == nil {
		reqLogger.Info(fmt.Sprintf("Updating nav config: %s", name))
	
		// Cast navItems interface to array of json objects
		navItems, ok := desiredNavConfig.Object["spec"].(map[string]interface{})["navItems"].([]map[string]interface{})
		if !ok {
			msg := fmt.Sprintf("Failed to unmarshal navItems array for nav config: %s", name)
			reqLogger.Info(msg)
			return errorf.New(msg)
		}

		// Update namespace for all nav items
		for _, item := range navItems {
			if item["namespace"] != "" {
				item["namespace"] = instance.Namespace
			}
		}
		
		// Set nav items to array with updated namespaces
		navConfig.Object["spec"].(map[string]interface{})["navItems"] = navItems
	
		// Update nav config
		err = r.Client.Update(ctx, navConfig)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to update nav config: %s", name), "NavConfig.Namespace", instance.Namespace, "NavConfig.Name", name)
			return err
		}
	}

	return nil
}
