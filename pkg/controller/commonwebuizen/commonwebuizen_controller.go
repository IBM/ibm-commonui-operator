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

package commonwebuizen

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	// "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	res "github.com/ibm/ibm-commonui-operator/pkg/resources"
	version "github.com/ibm/ibm-commonui-operator/version"
	routesv1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var log = logf.Log.WithName("controller_commonwebuizen")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CommonWebUIZen Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCommonWebUIZen{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {

	reqLogger := log.WithValues("func", "add")

	// Create a new controller
	c, err := controller.New("commonwebuizen-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	namespace := os.Getenv("WATCH_NAMESPACE")
	reqLogger.Info("Namespace in ZEN Watch: " + namespace)

	zenp := predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			if e.Object.GetName() == "zen-core" && e.Object.GetNamespace() == namespace {
				return true
			}
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			if e.Object.GetName() == "zen-core" && e.Object.GetNamespace() == namespace {
				return true
			}
			return false
		},
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}, zenp)

	if err != nil {
		return err
	}

	//Add watch for configmaps - specifically the zen product-configmap - add adminhub fields when it
	//exists
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew.GetName() == res.ZenProductConfigMapName && e.ObjectNew.GetNamespace() == namespace {
				return true
			}
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			if e.Object.GetName() == res.ZenProductConfigMapName && e.Object.GetNamespace() == namespace {
				return true
			}
			return false
		},
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}},
		handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "RECONCILE-ZEN-PRODUCT-CONFIGMAP",
					Namespace: a.GetNamespace(),
				}},
			}
		}),
		p)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCommonWebUIZen implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCommonWebUIZen{}

// ReconcileCommonWebUIZen reconciles a CommonWebUIZen object
type ReconcileCommonWebUIZen struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CommonWebUIZen object and makes changes based on the state read
// and what is in the CommonWebUIZen.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCommonWebUIZen) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUIZen")

	namespace := os.Getenv("WATCH_NAMESPACE")
	reqLogger.Info("Namespace in ZEN Reconcile: " + namespace)

	reqLogger.Info("got CommonWebUIZen operator version=" + version.Version)

	//If the request is for the zen product-configmap, reconcile the adminhub values for zen
	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	if request.Name == "RECONCILE-ZEN-PRODUCT-CONFIGMAP" {
		reqLogger.Info("Change to zen product configmap " + res.ZenProductConfigMapName + " detected - reconciling common webui updates")
		// Check if the config maps already exist. If not, create a new one.
		err := r.reconcileZenProductConfigMap(ctx, request, &needToRequeue)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	if needToRequeue {
		// one or more resources were created/updated, so requeue the request
		reqLogger.Info("Requeue the request")
		return reconcile.Result{Requeue: true}, nil
	}

	//Check to see kubernetes cluster type
	isCncf := r.getKubernetesClusterType(ctx, namespace)

	//Reconcile to see if Zen is enabled
	isZen := r.adminHubOnZen(ctx, namespace)

	if isZen {
		// Reconcile zen resources
		reqLogger.Info("Zen enabled in Reconcile")

		deleteErr := r.deleteClassicAdminHubRes(ctx, namespace)
		if deleteErr != nil {
			reqLogger.Error(deleteErr, "Failed deleting classic admin hub resources")
		}
		// Create common-web-ui-config
		err := r.reconcileConfigMapsZen(ctx, namespace, res.CommonConfigMap)
		if err != nil {
			return reconcile.Result{}, err
		}
		if isCncf {
			err = r.reconcileConfigMapsZen(ctx, namespace, res.ZenCardExtensionsConfigMapCncf)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			err = r.reconcileConfigMapsZen(ctx, namespace, res.ZenCardExtensionsConfigMap)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		err = r.reconcileConfigMapsZen(ctx, namespace, res.ZenQuickNavExtensionsConfigMap)
		if err != nil {
			return reconcile.Result{}, err
		}
		//uncomment for walkme support err = r.reconcileConfigMapsZen(ctx, namespace, res.ZenWalkmeExtensionsConfigMap)
		//if err != nil {
		//	return reconcile.Result{}, err
		//}
		updateErr := r.updateZenResources(ctx, namespace, res.ZenQuickNavExtensionsConfigMap)
		if updateErr != nil {
			reqLogger.Error(updateErr, "Failed updating zen card quick nav extensions")
			return reconcile.Result{}, updateErr
		}
		if isCncf {
			updateErr = r.updateZenResources(ctx, namespace, res.ZenCardExtensionsConfigMapCncf)
			if updateErr != nil {
				reqLogger.Error(updateErr, "Failed updating zen card extensions CNCF")
				return reconcile.Result{}, updateErr
			}
		} else {
			updateErr = r.updateZenResources(ctx, namespace, res.ZenCardExtensionsConfigMap)
			if updateErr != nil {
				reqLogger.Error(updateErr, "Failed updating zen card extensions")
				return reconcile.Result{}, updateErr
			}
		}
		updateErr = r.updateCommonUIDeployment(ctx, isZen, isCncf, namespace)
		if updateErr != nil {
			reqLogger.Error(updateErr, "Failed updating common ui deployment")
			return reconcile.Result{}, updateErr
		}
		deleteErr = r.deleteBindInfo(ctx, namespace)
		if deleteErr != nil {
			reqLogger.Error(deleteErr, "Failed deleting Common UI bind Info")
			return reconcile.Result{}, deleteErr
		}

		if !isCncf {
			err = r.reconcileCrZen(ctx, namespace, "admin-hub-zen", res.CrTemplates2, isZen)
			if err != nil {
				reqLogger.Error(err, "Error creating console link cr for zen")
				return reconcile.Result{}, err
			}
		}

	} else {
		err := r.deleteZenAdminHubRes(ctx, namespace)
		if err != nil {
			reqLogger.Error(err, "Error deleting zen admin hub resources")
			return reconcile.Result{}, err
		}
		updateErr := r.updateCommonUIDeployment(ctx, isZen, isCncf, namespace)
		if updateErr != nil {
			reqLogger.Error(updateErr, "Failed updating common ui deployment")
			return reconcile.Result{}, updateErr
		}
		if !isCncf {
			err = r.reconcileCrZen(ctx, namespace, "admin-hub", res.CrTemplates, isZen)
			if err != nil {
				reqLogger.Error(err, "Error creating console link cr")
				return reconcile.Result{}, err
			}
		}
	}

	reqLogger.Info("ZEN CONTROLLER??? all done")
	return reconcile.Result{}, nil
}

func (r *ReconcileCommonWebUIZen) adminHubOnZen(ctx context.Context, namespace string) bool {
	reqLogger := log.WithValues("func", "adminHubOnZen")
	reqLogger.Info("Checking zen optional install condition in commonui-zen controller")

	zenDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zen-core",
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: "zen-core", Namespace: namespace}, zenDeployment)

	if getError == nil {
		reqLogger.Info("Got ZEN Deployment in commonui-zen controller")
		return true
	}

	if errors.IsNotFound(getError) {
		reqLogger.Info("ZEN deployment not found in commonui-zen controller")
	} else {
		reqLogger.Error(getError, "Error getting ZEN deployment in commonui-zen controller")
	}

	return false
}

func (r *ReconcileCommonWebUIZen) reconcileConfigMapsZen(ctx context.Context, namespace, nameOfCM string) error {
	reqLogger := log.WithValues("func", "reconcileConfigMapsZen")

	reqLogger.Info("Checking if config map: " + nameOfCM + " exists")
	// Check if the config map already exists, if not create a new one
	currentConfigMap := &corev1.ConfigMap{}
	err := r.client.Get(ctx, types.NamespacedName{Name: nameOfCM, Namespace: namespace}, currentConfigMap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new ConfigMap
		newConfigMap := &corev1.ConfigMap{}
		if nameOfCM == res.ZenCardExtensionsConfigMap {
			reqLogger.Info("Creating zen card extensions config map")
			var ExtensionsData = map[string]string{
				"nginx.conf": res.ZenNginxConfig,
				"extensions": res.ZenCardExtensions,
			}
			newConfigMap = res.ZenCardExtensionsConfigMapUI(res.ZenCardExtensionsConfigMap, namespace, version.Version, ExtensionsData)
		} else if nameOfCM == res.ZenCardExtensionsConfigMapCncf {
			reqLogger.Info("Creating zen card extensions config map for CNCF")
			var ExtensionsData = map[string]string{
				"nginx.conf": res.ZenNginxConfig,
				"extensions": res.ZenCardExtensionsCncf,
			}
			newConfigMap = res.ZenCardExtensionsConfigMapUI(res.ZenCardExtensionsConfigMapCncf, namespace, version.Version, ExtensionsData)
		} else if nameOfCM == res.ZenQuickNavExtensionsConfigMap {
			reqLogger.Info("Creating zen quick nav extensions config map")
			var ExtensionsData = map[string]string{
				"extensions": res.ZenQuickNavExtensions,
			}
			newConfigMap = res.ZenCardExtensionsConfigMapUI(res.ZenQuickNavExtensionsConfigMap, namespace, version.Version, ExtensionsData)
		} else if nameOfCM == res.ZenWalkmeExtensionsConfigMap {
			reqLogger.Info("Creating zen walkme extensions config map")
			var ExtensionsData = map[string]string{
				"extensions": res.ZenWalkmeExtensions,
			}
			newConfigMap = res.ZenCardExtensionsConfigMapUI(res.ZenWalkmeExtensionsConfigMap, namespace, version.Version, ExtensionsData)
		} else if nameOfCM == res.CommonConfigMap {
			reqLogger.Info("Creating common-web-ui-config config map")
			newConfigMap = res.CommonWebUIConfigMap(namespace)
		} else if nameOfCM == res.ExtensionsConfigMap {
			currentRoute := &routesv1.Route{}
			//Get the cp-console route and add it to the configmap below
			err2 := r.client.Get(ctx, types.NamespacedName{Name: "cp-console", Namespace: namespace}, currentRoute)
			if err2 != nil {
				reqLogger.Error(err2, "Failed to get route for cp-console, try again later")
				return err2
			}
			reqLogger.Info("Current route is: " + currentRoute.Spec.Host)

			var ExtensionsData = map[string]string{
				"extensions": strings.Replace(res.Extensions, "/common-nav/dashboard", "https://"+currentRoute.Spec.Host+"/common-nav/dashboard", 1),
			}

			newConfigMap = res.ExtensionsConfigMapUI(namespace, ExtensionsData)

		}

		reqLogger.Info("Creating a config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
		err = r.client.Create(ctx, newConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to create a config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
			return err
		}
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Zen Config map")
		return err
	}

	reqLogger.Info("got Zen config map")

	return nil

}

func (r *ReconcileCommonWebUIZen) reconcileCrZen(ctx context.Context, namespace string, crName string, template string, isZen bool) error {
	reqLogger := log.WithValues("func", "reconcileCrZen")
	reqLogger.Info("RECONCILING CR ZEN")

	var crTemplate map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	crTemplatesErr := json.Unmarshal([]byte(template), &crTemplate)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	name := crName

	//Get CR and see if it exists
	getError := r.client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &unstruct)

	if getError != nil && !errors.IsNotFound(getError) {
		reqLogger.Error(getError, "Failed to get CR")
	} else if errors.IsNotFound(getError) {
		//If CR was not found, create it
		//Get the cpd route is zen is true
		currentRoute := &routesv1.Route{}
		if isZen {
			err2 := r.client.Get(ctx, types.NamespacedName{Name: "cpd", Namespace: namespace}, currentRoute)
			if err2 != nil {
				reqLogger.Error(err2, "Failed to get route for cpd, try again later")
				return err2
			}
			reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
			//Will hold href for admin hub console link
			var href = "https://" + currentRoute.Spec.Host

			// Create Custom resource
			if createErr := r.createCustomResource(ctx, unstruct, name, href); createErr != nil {
				reqLogger.Error(createErr, "Failed to create CR")
				return createErr
			}
		} else {
			err2 := r.client.Get(ctx, types.NamespacedName{Name: "cp-console", Namespace: namespace}, currentRoute)
			if err2 != nil {
				reqLogger.Error(err2, "Failed to get route for cp-console, try again later")
				return err2
			}
			reqLogger.Info("Current route is: " + currentRoute.Spec.Host)
			//Will hold href for admin hub console link
			var href = "https://" + currentRoute.Spec.Host + "/common-nav/dashboard"

			// Create Custom resource
			if createErr := r.createCustomResource(ctx, unstruct, name, href); createErr != nil {
				reqLogger.Error(createErr, "Failed to create CR")
				return createErr
			}
		}

	} else {
		reqLogger.Info("Skipping CR creation")
	}

	return nil
}

func (r *ReconcileCommonWebUIZen) createCustomResource(ctx context.Context, unstruct unstructured.Unstructured, name, href string) error {
	reqLogger := log.WithValues("CR name", name)
	reqLogger.Info("creating a CR ", name)

	unstruct.Object["spec"].(map[string]interface{})["href"] = href
	crCreateErr := r.client.Create(ctx, &unstruct)
	if crCreateErr != nil && !errors.IsAlreadyExists(crCreateErr) {
		reqLogger.Error(crCreateErr, "Failed to Create the Custom Resource")
		return crCreateErr
	}
	return nil
}

func (r *ReconcileCommonWebUIZen) deleteClassicAdminHubRes(ctx context.Context, namespace string) error {
	reqLogger := log.WithValues("func", "deleteClassicAdminHubRes")
	reqLogger.Info("Getting classic admin hub resources")

	reqLogger.Info("Checking to see if classic admin hub console link is present")
	var crTemplate map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	crTemplatesErr := json.Unmarshal([]byte(res.CrTemplates), &crTemplate)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	name := "admin-hub"

	//Get and delelte classic admin hub console link CR
	getError := r.client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &unstruct)

	if getError != nil && !errors.IsNotFound(getError) {
		reqLogger.Error(getError, "Failed to get classic admin hub console link CR")
	} else if getError == nil {
		reqLogger.Info("Got classic admin hub console link")
		err := r.client.Delete(ctx, &unstruct)
		if err != nil {
			reqLogger.Error(err, "Failed to delete classic admin hub console link")
		} else {
			reqLogger.Info("Deleted classic admin hub console link")
		}
	}

	//Get and delete classic admin hub left nav menu item
	reqLogger.Info("Checking to see if classic admin hub config map is present")
	currentConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.ExtensionsConfigMap,
			Namespace: namespace,
		},
	}
	getError2 := r.client.Get(ctx, types.NamespacedName{Name: res.ExtensionsConfigMap, Namespace: namespace}, currentConfigMap)

	if getError2 == nil {
		reqLogger.Info("Got classic admin hub config map")
		err := r.client.Delete(ctx, currentConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to delete classic admin hub config map")
		} else {
			reqLogger.Info("Deleted classic admin hub cofig map")
		}
	} else if !errors.IsNotFound(getError2) {
		reqLogger.Error(getError2, "Failed to get classic admin hub config map")
	}

	return nil
}

func (r *ReconcileCommonWebUIZen) deleteBindInfo(ctx context.Context, namespace string) error {
	reqLogger := log.WithValues("func", "deleteBindInfo")
	//Get and delete classic admin hub left nav menu item
	reqLogger.Info("Checking to see if Common UI bind info exists")
	currentConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibm-commonui-bindinfo-common-webui-ui-extensions",
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: "ibm-commonui-bindinfo-common-webui-ui-extensions", Namespace: namespace}, currentConfigMap)

	if getError == nil {
		reqLogger.Info("Got Common UI bind info")
		err := r.client.Delete(ctx, currentConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to delete Common UI bind info")
		} else {
			reqLogger.Info("Deleted Common UI bind info")
		}
	} else if !errors.IsNotFound(getError) {
		reqLogger.Error(getError, "Not found Common UI bind info")
		return getError
	} else {
		reqLogger.Error(getError, "Failed to get Common UI bind info")
	}

	return nil
}

func (r *ReconcileCommonWebUIZen) shouldUpdateZenResources(ctx context.Context, nameOfCM string, namespace string) bool {
	reqLogger := log.WithValues("func", "shouldUpdateZenResources")
	reqLogger.Info("Checking zen upgrade condition")

	currentConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nameOfCM,
			Namespace: namespace,
		},
	}
	err := r.client.Get(ctx, types.NamespacedName{Name: nameOfCM, Namespace: namespace}, currentConfigMap)

	if err == nil {
		reqLogger.Info("Comparing versions")
		currentVersion := currentConfigMap.Labels["icpdata_addon_version"]
		newVersion := "v" + version.Version
		reqLogger.Info("Old Version: " + currentVersion)
		reqLogger.Info("New Version: " + newVersion)
		if currentVersion != newVersion {
			return true
		}
	}

	return false
}

func (r *ReconcileCommonWebUIZen) updateZenResources(ctx context.Context, namespace, nameOfCM string) error {
	reqLogger := log.WithValues("func", "updateZenResources")

	reqLogger.Info("checking if zen card extensions config map exists")

	if r.shouldUpdateZenResources(ctx, nameOfCM, namespace) {
		currentConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nameOfCM,
				Namespace: namespace,
			},
		}
		err := r.client.Get(ctx, types.NamespacedName{Name: nameOfCM, Namespace: namespace}, currentConfigMap)

		if err == nil {
			reqLogger.Info("zen card extensions config map exists")

			var ExtensionsData map[string]string

			if nameOfCM == res.ZenQuickNavExtensionsConfigMap {
				ExtensionsData = map[string]string{
					"extensions": res.ZenQuickNavExtensions,
				}
			} else if nameOfCM == res.ZenCardExtensionsConfigMapCncf {
				ExtensionsData = map[string]string{
					"nginx.conf": res.ZenNginxConfig,
					"extensions": res.ZenCardExtensionsCncf,
				}
			} else {
				ExtensionsData = map[string]string{
					"nginx.conf": res.ZenNginxConfig,
					"extensions": res.ZenCardExtensions,
				}
			}

			currentConfigMap.Labels["icpdata_addon_version"] = version.Version
			currentConfigMap.Data = ExtensionsData

			reqLogger.Info("Updating zen card extensions CM")
			updateErr := r.client.Update(ctx, currentConfigMap)
			if updateErr == nil {
				reqLogger.Info("Card extensions updated")
			} else {
				reqLogger.Error(updateErr, "Could not update card extensions")
				return updateErr
			}
		} else {
			return err
		}
	}

	return nil
}

func (r *ReconcileCommonWebUIZen) deleteZenAdminHubRes(ctx context.Context, namespace string) error {
	reqLogger := log.WithValues("func", "deleteZenAdminHubRes")
	reqLogger.Info("Getting ZEN admin hub resources")
	//Get and delete classic admin hub left nav menu item
	reqLogger.Info("Checking to see if ZEN admin hub config maps are present")

	currentConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.ZenCardExtensionsConfigMap,
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: res.ZenCardExtensionsConfigMap, Namespace: namespace}, currentConfigMap)

	if getError == nil {
		reqLogger.Info("Got ZEN admin hub config maps")
		err := r.client.Delete(ctx, currentConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to delete ZEN admin hub config maps")
		} else {
			reqLogger.Info("Deleted ZEN admin hub config maps")
		}
	} else if !errors.IsNotFound(getError) {
		reqLogger.Error(getError, "Failed to get ZEN admin hub config maps")
	}

	reqLogger.Info("Checking to see if zen admin hub console link is present")
	var crTemplate map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	crTemplatesErr := json.Unmarshal([]byte(res.CrTemplates2), &crTemplate)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	name := "admin-hub-zen"

	//Get and delelte classic admin hub console link CR
	getError2 := r.client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &unstruct)

	if getError2 != nil && !errors.IsNotFound(getError2) {
		reqLogger.Error(getError2, "Failed to get zen admin hub console link CR")
	} else if getError == nil {
		reqLogger.Info("Got zen admin hub console link")
		err := r.client.Delete(ctx, &unstruct)
		if err != nil {
			reqLogger.Error(err, "Failed to delete zen admin hub console link")
		} else {
			reqLogger.Info("Deleted zen admin hub console link")
		}
	}

	currentConfigMap2 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.CommonConfigMap,
			Namespace: namespace,
		},
	}
	getError3 := r.client.Get(ctx, types.NamespacedName{Name: res.CommonConfigMap, Namespace: namespace}, currentConfigMap2)

	if getError3 == nil {
		reqLogger.Info("Got common web ui config")
		err := r.client.Delete(ctx, currentConfigMap2)
		if err != nil {
			reqLogger.Error(err, "Failed to delete common web ui config")
		} else {
			reqLogger.Info("Deleted common web ui config")
		}
	} else if !errors.IsNotFound(getError3) {
		reqLogger.Error(getError3, "Failed to get common web ui config")
	}

	return nil
}

func (r *ReconcileCommonWebUIZen) updateCommonUIDeployment(ctx context.Context, isZen bool, isCncf bool, namespace string) error {
	reqLogger := log.WithValues("func", "updateCommonUIDeployment")
	reqLogger.Info("Updating common ui deployment env variable")

	commonDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "common-web-ui",
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: "common-web-ui", Namespace: namespace}, commonDeployment)

	if getError == nil {
		reqLogger.Info("Got Common UI deployment")
		env := commonDeployment.Spec.Template.Spec.Containers[0].Env[22].Value
		clusterTypeEnvVar := commonDeployment.Spec.Template.Spec.Containers[0].Env[24].Value
		clusterType := "cncf"
		if isZen && env == "false" {
			reqLogger.Info("Setting use zen to true")
			commonDeployment.Spec.Template.Spec.Containers[0].Env[22].Value = "true"
			if isCncf && clusterTypeEnvVar != clusterType {
				reqLogger.Info("Setting cluster type env var to cncf for zen case")
				commonDeployment.Spec.Template.Spec.Containers[0].Env[24].Value = clusterType
			}
			updateErr := r.client.Update(ctx, commonDeployment)
			if updateErr == nil {
				reqLogger.Info("Updated common ui deployment env variable")
			} else {
				reqLogger.Error(updateErr, "Could not update common ui deployment env variable")
				return updateErr
			}
		} else if !isZen && env == "true" {
			reqLogger.Info("Setting use zen to false")
			commonDeployment.Spec.Template.Spec.Containers[0].Env[22].Value = "false"
			updateErr := r.client.Update(ctx, commonDeployment)
			if updateErr == nil {
				reqLogger.Info("Updated common ui deployment env variable")
			} else {
				reqLogger.Error(updateErr, "Could not update common ui deployment env variable")
				return updateErr
			}
		} else if !isZen && isCncf {
			reqLogger.Info("Setting cluster type env var to cncf for non zen case")
			commonDeployment.Spec.Template.Spec.Containers[0].Env[24].Value = clusterType
			updateErr := r.client.Update(ctx, commonDeployment)
			if updateErr == nil {
				reqLogger.Info("Updated common ui deployment with cluster type")
			} else {
				reqLogger.Error(updateErr, "Could not update common ui with cluster type")
				return updateErr
			}
		}
	} else {
		reqLogger.Info("Failed to get Common UI deployment")
		return getError
	}
	return nil
}

//reconcileZenProductConfigMap will add the required adminhub values to the zen product-configmap
//if the configmap isn't found, it is assumed to be deleted and nothing is done
//returns an error or nil if successful
func (r *ReconcileCommonWebUIZen) reconcileZenProductConfigMap(ctx context.Context, request reconcile.Request, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileConfiMaps", "request", request)

	productConfigMap := &corev1.ConfigMap{}
	err := r.client.Get(ctx, types.NamespacedName{Name: res.ZenProductConfigMapName, Namespace: request.Namespace}, productConfigMap)
	if err != nil {
		if errors.IsNotFound(err) {
			//The zen product-configmap doesn't exist, if this is a zen adminhub instance then there is nothing to do, however
			//in this case the product-configmap will already exist so in theory we should never get here in that condition
			reqLogger.Info("Zen product configmap " + res.ZenProductConfigMapName + " was not found - assumed to be deleted")
			return nil
		}

		reqLogger.Error(err, "Failed to get zen product configmap "+res.ZenProductConfigMapName)
		return err
	}

	//If the requested fields are missing or not equal, then update the configmap
	update := false

	if productConfigMap.Data == nil {
		productConfigMap.Data = map[string]string{}
	}

	for key, element := range res.ZenPcmMap {
		zelement, ok := productConfigMap.Data[key]
		if !ok || element != zelement {
			productConfigMap.Data[key] = element
			update = true
			log.Info("updated zen product configmap "+res.ZenProductConfigMapName, "key", key, "value", element)
		}
	}

	if update {
		updateErr := r.client.Update(ctx, productConfigMap)
		if updateErr != nil {
			reqLogger.Info("Could not update zen product configmap "+res.ZenProductConfigMapName, updateErr)
			return updateErr
		}
		reqLogger.Info("Zen product configmap " + res.ZenProductConfigMapName + " updated with meta values")
		*needToRequeue = true
	} else {
		reqLogger.Info("No updates required for commonwebui in zen product configmap " + res.ZenProductConfigMapName)
	}

	return nil

}

func (r *ReconcileCommonWebUIZen) getKubernetesClusterType(ctx context.Context, namespace string) bool {
	reqLogger := log.WithValues("func", "isCncf")
	reqLogger.Info("Checking kubernetes cluster type")

	ibmProjectK := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibm-cpp-config",
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: "ibm-cpp-config", Namespace: namespace}, ibmProjectK)

	if getError == nil {
		reqLogger.Info("Got ibm project k config map")
		clusterType := ibmProjectK.Data["kubernetes_cluster_type"]
		if clusterType == "cncf" {
			reqLogger.Info("Kubernetes cluster type is " + clusterType)
			return true
		}
	}

	if errors.IsNotFound(getError) {
		reqLogger.Info("ibm project k config map not found in cs namepace")
	} else {
		reqLogger.Error(getError, "error getting ibm project k config map in cs namepace")
	}

	return false
}
