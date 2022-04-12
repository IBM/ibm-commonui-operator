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
	"encoding/json"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
	routesv1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileConsoleLink(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, isZen bool, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileConsoleLink", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling ConsoleLink CR")

	var consoleLink map[string]interface{}
	err := json.Unmarshal([]byte(ConsoleLinkTemplate), &consoleLink)
	if err != nil {
		reqLogger.Info("Failed to unmarshall ConsoleLink CR template")
		return err
	}

	var unstruct unstructured.Unstructured
	unstruct.Object = consoleLink
	name := unstruct.Object["metadata"].(map[string]interface{})["name"].(string)

	getErr := client.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: name}, &unstruct)

	hasFinalizer := ContainsString(instance.ObjectMeta.Finalizers, finalizerName)
	hasFinalizer1 := ContainsString(instance.ObjectMeta.Finalizers, finalizerName1)
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !hasFinalizer && !hasFinalizer1 {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName, finalizerName1)
			if err := client.Update(ctx, instance); err != nil {
				reqLogger.Error(err, "Failed to create finalizer")
			} else {
				reqLogger.Info("Created Finalizers")
			}
		}
	} else if hasFinalizer {
		// Finalizer is present, so lets handle any external dependency - remove console link CR
		if err := client.Delete(ctx, &unstruct); err != nil {
			// if fails to delete the external dependency here, return with error
			reqLogger.Error(err, "Failed to delete Console Link CR")
		} else {
			reqLogger.Info("Deleted Console link CR")
		}

		// Remove our finalizer from the metadata of the object and update it.
		instance.ObjectMeta.Finalizers = RemoveString(instance.ObjectMeta.Finalizers, finalizerName)

		if err := client.Update(ctx, instance); err != nil {
			reqLogger.Error(err, "Failed to delete  Console link finalizer")
		} else {
			reqLogger.Info("Deleted Console link Finalizer")
		}
	} else if hasFinalizer1 {
		// Finalizer is present, so lets handle any external dependency - remove console link CR
		if err := client.Delete(ctx, &unstruct); err != nil {
			// if fails to delete the external dependency here, return with error
			reqLogger.Error(err, "Failed to delete Redis CR")
		} else {
			reqLogger.Info("Deleted Redis CR")
		}

		// Remove our finalizer from the metadata of the object and update it.
		instance.ObjectMeta.Finalizers = RemoveString(instance.ObjectMeta.Finalizers, finalizerName1)
		if err := client.Update(ctx, instance); err != nil {
			reqLogger.Error(err, "Failed to delete Redis finalizer")
		} else {
			reqLogger.Info("Deleted Redis Finalizer")
		}
	}

	if getErr != nil {
		if errors.IsNotFound(getErr) {
			//If CR was not found, create it
			//Get the cpd route is zen is true
			currentRoute := &routesv1.Route{}
			if isZen {
				err := client.Get(ctx, types.NamespacedName{Name: "cpd", Namespace: instance.Namespace}, currentRoute)
				if err != nil {
					reqLogger.Error(err, "Failed to get route for cpd, try again later")
				}
				reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
				//Will hold href for admin hub console link
				var href = "https://" + currentRoute.Spec.Host

				// Create Custom resource
				if createErr := createCustomResource(ctx, client, unstruct, name, href); createErr != nil {
					reqLogger.Error(createErr, "Failed to create CR")
					return createErr
				}
			} else { //Get the cp-console route
				err := client.Get(ctx, types.NamespacedName{Name: "cp-console", Namespace: instance.Namespace}, currentRoute)
				if err != nil {
					reqLogger.Error(err, "Failed to get route for cp-console, try again later")
				}
				reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
				//Will hold href for admin hub console link
				var href = "https://" + currentRoute.Spec.Host + "/common-nav/dashboard"

				// Create Custom resource
				if createErr := createCustomResource(ctx, client, unstruct, name, href); createErr != nil {
					reqLogger.Error(createErr, "Failed to create CR")
					return createErr
				}
			}
		} else {
			reqLogger.Error(err, "Failed to get ConsoleLink CR", "ConsoleLink.Namespace", instance.Namespace, "ConsoleLink.Name", name)
			return err
		}
	}

	return nil
}

func createCustomResource(ctx context.Context, client client.Client, unstruct unstructured.Unstructured, name, href string) error {
	reqLogger := log.WithValues("func", "createCustomResource", "CR.Name", name)
	reqLogger.Info("Creating CR ", name)

	unstruct.Object["spec"].(map[string]interface{})["href"] = href

	err := client.Create(ctx, &unstruct)
	if err != nil && !errors.IsAlreadyExists(err) {
		reqLogger.Error(err, "Failed to create CR")
		return err
	}

	return nil
}