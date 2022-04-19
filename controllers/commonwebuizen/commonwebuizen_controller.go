/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	res "github.com/IBM/ibm-commonui-operator/controllers/resources"
	version "github.com/IBM/ibm-commonui-operator/version"
)

var log = logf.Log.WithName("controller_commonwebuizen")

// CommonWebUIZenReconciler reconciles a CommonWebUIZen object
type CommonWebUIZenReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ibm.com,resources=commonwebuizens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ibm.com,resources=commonwebuizens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ibm.com,resources=commonwebuizens/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CommonWebUIZen object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CommonWebUIZenReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUIZen Controller")

	namespace := os.Getenv("WATCH_NAMESPACE")
	reqLogger.Info("In CommonWebUIZen Reconcile -- Common Services Pod namespace: " + namespace)
	reqLogger.Info("In CommonWebUIZen Reconcile -- Operator version: " + version.Version)

	//If the request is for the zen product-configmap, reconcile the adminhub values for zen
	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	if request.Name == "RECONCILE-ZEN-PRODUCT-CONFIGMAP" {
		reqLogger.Info("Change to zen product configmap " + res.ZenProductConfigMapName + " detected - reconciling common webui updates")
		// Check if the config maps already exist. If not, create a new one.
		err := res.ReconcileZenProductConfigMap(ctx, r.Client, request, &needToRequeue)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Check to see if Zen instance exists in common services namespace
	isZen := res.IsAdminHubOnZen(ctx, r.Client, namespace)

	// Check to see kubernetes cluster type is cncf
	isCncf := res.GetKubernetesClusterType(ctx, r.Client, namespace)

	if isZen {
		// Reconcile zen resources
		reqLogger.Info("Zen Instance found, deleting classic Admin hub resources and Reconciling Zen resources")

		// Delete classic adminhub resources
		deleteErr := r.deleteClassicAdminHubRes(ctx, namespace)
		if deleteErr != nil {
			reqLogger.Error(deleteErr, "Failed deleting classic admin hub resources")
		}

		// Create common-web-ui-config which contains the common ui app version needed for zen nls post
		err := res.ReconcileConfigMapsZen(ctx, r.Client, version.Version, namespace, res.CommonConfigMapName)
		if err != nil {
			return ctrl.Result{}, err
		}
		if isCncf {
			// Create Zen card extensions for common ui on CNCF
			err = res.ReconcileConfigMapsZen(ctx, r.Client, version.Version, namespace, res.ZenCardExtensionsConfigMapNameCncf)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			// Create Zen card extensions for common ui on openshift
			err = res.ReconcileConfigMapsZen(ctx, r.Client, version.Version, namespace, res.ZenCardExtensionsConfigMapName)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		// Create Zen quick nav extensions for common ui
		err = res.ReconcileConfigMapsZen(ctx, r.Client, version.Version, namespace, res.ZenQuickNavExtensionsConfigMapName)
		if err != nil {
			return ctrl.Result{}, err
		}

		// On upgrade, update zen extensions
		updateErr := r.updateZenResources(ctx, namespace, res.ZenQuickNavExtensionsConfigMapName)
		if updateErr != nil {
			reqLogger.Error(updateErr, "Failed updating zen card quick nav extensions")
			return ctrl.Result{}, updateErr
		}
		if isCncf {
			updateErr = r.updateZenResources(ctx, namespace, res.ZenCardExtensionsConfigMapNameCncf)
			if updateErr != nil {
				reqLogger.Error(updateErr, "Failed updating zen card extensions CNCF")
				return ctrl.Result{}, updateErr
			}
		} else {
			updateErr = r.updateZenResources(ctx, namespace, res.ZenCardExtensionsConfigMapName)
			if updateErr != nil {
				reqLogger.Error(updateErr, "Failed updating zen card extensions")
				return ctrl.Result{}, updateErr
			}
		}
		// Set env var USE_ZEN to true and update CLUSTER_TYPE
		updateErr = r.updateCommonUIDeployment(ctx, isZen, isCncf, namespace)
		if updateErr != nil {
			reqLogger.Error(updateErr, "Failed updating common ui deployment")
			return ctrl.Result{}, updateErr
		}

	} else {
		// Delete zen admin hub resources
		deleteErr := r.deleteZenAdminHubRes(ctx, namespace)
		if deleteErr != nil {
			reqLogger.Error(deleteErr, "Error deleting zen admin hub resources")
			return ctrl.Result{}, deleteErr
		}
		// Set env var USE_ZEN to false and update CLUSTER_TYPE
		updateErr := r.updateCommonUIDeployment(ctx, isZen, isCncf, namespace)
		if updateErr != nil {
			reqLogger.Error(updateErr, "Failed updating common ui deployment")
			return ctrl.Result{}, updateErr
		}
	}

	if needToRequeue {
		// one or more resources were created/updated, so requeue the request
		reqLogger.Info("Requeue the request")
		return ctrl.Result{Requeue: true}, nil
	}

	reqLogger.Info("COMMON UI ZEN CONTROLLER RECONCILE ALL DONE")
	return ctrl.Result{}, nil
}

func (r *CommonWebUIZenReconciler) deleteClassicAdminHubRes(ctx context.Context, namespace string) error {
	reqLogger := log.WithValues("func", "deleteClassicAdminHubRes")
	reqLogger.Info("Getting classic admin hub resources")

	reqLogger.Info("Checking to see if classic admin hub console link is present")
	var crTemplate map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	crTemplatesErr := json.Unmarshal([]byte(res.ConsoleLinkTemplate), &crTemplate)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	name := "admin-hub"

	//Get and delelte classic admin hub console link CR
	getCr := r.Client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &unstruct)

	if getCr != nil && !errors.IsNotFound(getCr) {
		reqLogger.Error(getCr, "Failed to get classic admin hub console link CR")
	} else if getCr == nil {
		reqLogger.Info("Got classic admin hub console link")
		err := r.Client.Delete(ctx, &unstruct)
		if err != nil {
			reqLogger.Error(err, "Failed to delete classic admin hub console link")
		} else {
			reqLogger.Info("Deleted classic admin hub console link")
		}
	}

	//Get and delete common ui bind info config map
	reqLogger.Info("Checking to see if Common UI bind info configmap exists")
	bindinfoConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibm-commonui-bindinfo-common-webui-ui-extensions",
			Namespace: namespace,
		},
	}
	getBindinfo := r.Client.Get(ctx, types.NamespacedName{Name: "ibm-commonui-bindinfo-common-webui-ui-extensions", Namespace: namespace}, bindinfoConfigMap)

	if getBindinfo == nil {
		reqLogger.Info("Got Common UI bind info")
		err := r.Client.Delete(ctx, bindinfoConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to delete Common UI bind info configmap")
		} else {
			reqLogger.Info("Deleted Common UI bind info configmap")
		}
	} else if !errors.IsNotFound(getBindinfo) {
		reqLogger.Error(getBindinfo, "Not found Common UI bind info configmap")
		return getBindinfo
	} else {
		reqLogger.Error(getBindinfo, "Failed to get Common UI bind info configmap")
	}

	//Get and delete classic admin hub left nav menu item
	reqLogger.Info("Checking to see if classic adminhub extensions configmap is present")
	extensionsConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.ZenLeftNavExtensionsConfigMapName,
			Namespace: namespace,
		},
	}
	getExtension := r.Client.Get(ctx, types.NamespacedName{Name: res.ZenLeftNavExtensionsConfigMapName, Namespace: namespace}, extensionsConfigMap)

	if getExtension == nil {
		reqLogger.Info("Got classic adminhub extensions configmap")
		err := r.Client.Delete(ctx, extensionsConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to delete classic adminhub extensions configmap")
		} else {
			reqLogger.Info("Deleted classic adminhub extensions cofigmap")
		}
	} else if !errors.IsNotFound(getExtension) {
		reqLogger.Error(getExtension, "Failed to get classic adminhub configmap")
	}

	return nil
}

func (r *CommonWebUIZenReconciler) shouldUpdateZenResources(ctx context.Context, nameOfCM string, namespace string) bool {
	reqLogger := log.WithValues("func", "shouldUpdateZenResources")
	reqLogger.Info("Checking zen upgrade condition")

	currentConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nameOfCM,
			Namespace: namespace,
		},
	}
	err := r.Client.Get(ctx, types.NamespacedName{Name: nameOfCM, Namespace: namespace}, currentConfigMap)

	if err == nil {
		reqLogger.Info("Comparing versions")
		oldVersion := currentConfigMap.Labels["icpdata_addon_version"]
		newVersion := "v" + version.Version
		reqLogger.Info("Old Version: " + oldVersion)
		reqLogger.Info("New Version: " + newVersion)
		if oldVersion != newVersion {
			return true
		}
	}

	return false
}

func (r *CommonWebUIZenReconciler) updateZenResources(ctx context.Context, namespace, nameOfCM string) error {
	reqLogger := log.WithValues("func", "updateZenResources")

	reqLogger.Info("checking if zen card extensions config map exists")

	if r.shouldUpdateZenResources(ctx, nameOfCM, namespace) {
		currentConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nameOfCM,
				Namespace: namespace,
			},
		}
		err := r.Client.Get(ctx, types.NamespacedName{Name: nameOfCM, Namespace: namespace}, currentConfigMap)

		if err == nil {
			reqLogger.Info("zen card extensions config map exists")

			var ExtensionsData map[string]string

			if nameOfCM == res.ZenQuickNavExtensionsConfigMapName {
				ExtensionsData = map[string]string{
					"extensions": res.ZenQuickNavExtensions,
				}
			} else if nameOfCM == res.ZenCardExtensionsConfigMapNameCncf {
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
			updateErr := r.Client.Update(ctx, currentConfigMap)
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

func (r *CommonWebUIZenReconciler) updateCommonUIDeployment(ctx context.Context, isZen bool, isCncf bool, namespace string) error {
	reqLogger := log.WithValues("func", "updateCommonUIDeployment")
	reqLogger.Info("Updating common ui deployment env variable")

	commonDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "common-web-ui",
			Namespace: namespace,
		},
	}
	getError := r.Client.Get(ctx, types.NamespacedName{Name: "common-web-ui", Namespace: namespace}, commonDeployment)

	if getError == nil {
		reqLogger.Info("Got Common UI deployment")
		env := commonDeployment.Spec.Template.Spec.Containers[0].Env[22].Value
		clusterType := "cncf"
		if isZen && env == "false" {
			reqLogger.Info("Setting use zen to true")
			commonDeployment.Spec.Template.Spec.Containers[0].Env[22].Value = "true"
			if isCncf {
				reqLogger.Info("Setting cluster type env var to cncf for zen case")
				envLength := len(commonDeployment.Spec.Template.Spec.Containers[0].Env)
				if envLength > 24 {
					clusterTypeEnvVar := commonDeployment.Spec.Template.Spec.Containers[0].Env[24].Value
					if clusterTypeEnvVar != clusterType {
						commonDeployment.Spec.Template.Spec.Containers[0].Env[24].Value = clusterType
					}
				}
			}
			updateErr := r.Client.Update(ctx, commonDeployment)
			if updateErr == nil {
				reqLogger.Info("Updated common ui deployment env variable")
			} else {
				reqLogger.Error(updateErr, "Could not update common ui deployment env variable")
				return updateErr
			}
		} else if !isZen && env == "true" {
			reqLogger.Info("Setting use zen to false")
			commonDeployment.Spec.Template.Spec.Containers[0].Env[22].Value = "false"
			updateErr := r.Client.Update(ctx, commonDeployment)
			if updateErr == nil {
				reqLogger.Info("Updated common ui deployment env variable")
			} else {
				reqLogger.Error(updateErr, "Could not update common ui deployment env variable")
				return updateErr
			}
		} else if !isZen && isCncf {
			reqLogger.Info("Setting cluster type env var to cncf for non zen case")
			envLength := len(commonDeployment.Spec.Template.Spec.Containers[0].Env)
			if envLength > 24 {
				commonDeployment.Spec.Template.Spec.Containers[0].Env[24].Value = clusterType
			}
			updateErr := r.Client.Update(ctx, commonDeployment)
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

func (r *CommonWebUIZenReconciler) deleteZenAdminHubRes(ctx context.Context, namespace string) error {
	reqLogger := log.WithValues("func", "deleteZenAdminHubRes")
	reqLogger.Info("Getting ZEN admin hub resources")
	reqLogger.Info("Checking to see if ZEN admin hub config maps are present")

	//Get and delete zen admin hub left nav menu item
	currentConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.ZenCardExtensionsConfigMapName,
			Namespace: namespace,
		},
	}
	getError := r.Client.Get(ctx, types.NamespacedName{Name: res.ZenCardExtensionsConfigMapName, Namespace: namespace}, currentConfigMap)

	if getError == nil {
		reqLogger.Info("Got ZEN admin hub config maps")
		err := r.Client.Delete(ctx, currentConfigMap)
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
	crTemplatesErr := json.Unmarshal([]byte(res.ConsoleLinkTemplate2), &crTemplate)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	name := "admin-hub-zen"

	//Get and delelte classic admin hub console link CR
	getError2 := r.Client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &unstruct)

	if getError2 != nil && !errors.IsNotFound(getError2) {
		reqLogger.Error(getError2, "Failed to get zen admin hub console link CR")
	} else if getError == nil {
		reqLogger.Info("Got zen admin hub console link")
		err := r.Client.Delete(ctx, &unstruct)
		if err != nil {
			reqLogger.Error(err, "Failed to delete zen admin hub console link")
		} else {
			reqLogger.Info("Deleted zen admin hub console link")
		}
	}

	currentConfigMap2 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.CommonConfigMapName,
			Namespace: namespace,
		},
	}
	getError3 := r.Client.Get(ctx, types.NamespacedName{Name: res.CommonConfigMapName, Namespace: namespace}, currentConfigMap2)

	if getError3 == nil {
		reqLogger.Info("Got common web ui config")
		err := r.Client.Delete(ctx, currentConfigMap2)
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

func zenDeploymentPredicate() predicate.Predicate {
	namespace := os.Getenv("WATCH_NAMESPACE")
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			if e.Object.GetName() == res.ZenDeploymentName && e.Object.GetNamespace() == namespace {
				return true
			}
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			if e.Object.GetName() == res.ZenDeploymentName && e.Object.GetNamespace() == namespace {
				return true
			}
			return false
		},
	}
}

func zenProductCmPredicate() predicate.Predicate {
	namespace := os.Getenv("WATCH_NAMESPACE")
	return predicate.Funcs{
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
}

// SetupWithManager sets up the controller with the Manager.
func (r *CommonWebUIZenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create a new controller
	c, err := controller.New("commonwebuizen-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}, zenDeploymentPredicate())
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}},
		handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
			return []ctrl.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "RECONCILE-ZEN-PRODUCT-CONFIGMAP",
					Namespace: a.GetNamespace(),
				}},
			}
		}),
		zenProductCmPredicate())
	if err != nil {
		return err
	}

	return nil
}
