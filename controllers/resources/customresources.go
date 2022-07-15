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

	routesv1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

func ReconcileConsoleLink(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, isZen bool, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileConsoleLink", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling ConsoleLink CR")

	var consoleLink map[string]interface{}
	var consoleLinkTemplate string

	if isZen {
		consoleLinkTemplate = ConsoleLinkTemplate2
	} else {
		consoleLinkTemplate = ConsoleLinkTemplate
	}

	err := json.Unmarshal([]byte(consoleLinkTemplate), &consoleLink)
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

	if getErr == nil {
		reqLogger.Info("CR already present, checking for updates")
		if unstruct.Object["spec"].(map[string]interface{})["applicationMenu"] == nil {
			reqLogger.Info("Console link CR missing attributes, trying to update")
			var currentTemplate map[string]interface{}
			crTemplateErr := json.Unmarshal([]byte(consoleLinkTemplate), &currentTemplate)
			if crTemplateErr != nil {
				reqLogger.Info("Failed to console link cr")
				return crTemplateErr
			}
			var unstruct2 unstructured.Unstructured
			unstruct2.Object = currentTemplate
			if updateErr := updateConsoleLink(ctx, client, unstruct, unstruct2, isZen, instance.Namespace); updateErr != nil {
				reqLogger.Error(updateErr, "Failed to update console link CR")
				return updateErr
			}
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
					if errors.IsNotFound(err) {
						reqLogger.Info("cpd route not found")
					} else {
						reqLogger.Error(err, "Failed to get route for cpd, try again later")
					}
				} else {
					reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
					//Will hold href for admin hub console link
					var href = "https://" + currentRoute.Spec.Host

					// Create Custom resource
					if createErr := createCustomResource(ctx, client, unstruct, name, href); createErr != nil {
						reqLogger.Error(createErr, "Failed to create CR")
						return createErr
					}
				}
				return err
			} else { //Get the cp-console route
				err := client.Get(ctx, types.NamespacedName{Name: "cp-console", Namespace: instance.Namespace}, currentRoute)
				if err != nil {
					if errors.IsNotFound(err) {
						reqLogger.Info("cp-console route not found")
					} else {
						reqLogger.Error(err, "Failed to get route for cp-console, try again later")
					}
				} else {
					reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
					//Will hold href for admin hub console link
					var href = "https://" + currentRoute.Spec.Host + "/common-nav/dashboard"

					// Create Custom resource
					if createErr := createCustomResource(ctx, client, unstruct, name, href); createErr != nil {
						reqLogger.Error(createErr, "Failed to create CR")
						return createErr
					}
				}
				return err
			}
		} else {
			reqLogger.Error(err, "Failed to get ConsoleLink CR", "ConsoleLink.Namespace", instance.Namespace, "ConsoleLink.Name", name)
			return err
		}
	}

	return nil
}

func createCustomResource(ctx context.Context, client client.Client, unstruct unstructured.Unstructured, name, href string) error {
	reqLogger := log.WithValues("func", "createCustomResource", "cr.Name", name)
	reqLogger.Info("Creating CR")

	unstruct.Object["spec"].(map[string]interface{})["href"] = href

	err := client.Create(ctx, &unstruct)
	if err != nil && !errors.IsAlreadyExists(err) {
		reqLogger.Error(err, "Failed to create CR")
		return err
	}

	return nil
}

//nolint
func updateConsoleLink(ctx context.Context, client client.Client, unstruct unstructured.Unstructured, unstruct2 unstructured.Unstructured, isZen bool, namespace string) error {
	reqLogger := log.WithValues("func", "updateCustomResource")
	reqLogger.Info("Updating console link cr")

	currentRoute := &routesv1.Route{}
	if isZen {
		err2 := client.Get(ctx, types.NamespacedName{Name: "cpd", Namespace: namespace}, currentRoute)
		if err2 != nil {
			reqLogger.Error(err2, "Failed to get route for cpd, try again later")
			return err2
		}
		reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
		var href = "https://" + currentRoute.Spec.Host
		unstruct2.Object["spec"].(map[string]interface{})["href"] = href
		unstruct2.Object["spec"].(map[string]interface{})["applicationMenu"].(map[string]interface{})["imageURL"] = href + "/common-nav/graphics/settings.svg"
	} else {
		err2 := client.Get(ctx, types.NamespacedName{Name: "cp-console", Namespace: namespace}, currentRoute)
		if err2 != nil {
			reqLogger.Error(err2, "Failed to get route for cp-console, try again later")
			return err2
		}
		reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
		var href = "https://" + currentRoute.Spec.Host
		unstruct2.Object["spec"].(map[string]interface{})["href"] = href + "/common-nav/dashboard"
		unstruct2.Object["spec"].(map[string]interface{})["applicationMenu"].(map[string]interface{})["imageURL"] = href + "/common-nav/graphics/settings.svg"
	}

	unstruct.Object["spec"] = unstruct2.Object["spec"]
	crUpdateErr := client.Update(ctx, &unstruct)
	if crUpdateErr != nil && !errors.IsAlreadyExists(crUpdateErr) {
		reqLogger.Error(crUpdateErr, "Failed to Create the Custom Resource")
		return crUpdateErr
	}
	return nil
}
