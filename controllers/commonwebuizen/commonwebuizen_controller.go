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

	res "github.com/IBM/ibm-commonui-operator/controllers/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

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
	reqLogger.Info("Reconciling CommonWebUIZen")

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
		err := res.ReconcileConfigMapsZen(ctx, r.Client, namespace, res.CommonConfigMapName)
		if err != nil {
			return ctrl.Result{}, err
		}
		if isCncf {
			// Create Zen card extensions for common ui on CNCF
			err = res.ReconcileConfigMapsZen(ctx, r.Client, namespace, res.ZenCardExtensionsConfigMapNameCncf)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			// Create Zen card extensions for common ui on openshift
			err = res.ReconcileConfigMapsZen(ctx, r.Client, namespace, res.ZenCardExtensionsConfigMapName)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		// Create Zen quick nav extensions for common ui
		err = res.ReconcileConfigMapsZen(ctx, r.Client, namespace, res.ZenQuickNavExtensionsConfigMapName)
		if err != nil {
			return ctrl.Result{}, err
		}

	}

	if needToRequeue {
		// one or more resources were created/updated, so requeue the request
		reqLogger.Info("Requeue the request")
		return ctrl.Result{Requeue: true}, nil
	}

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
	return ctrl.NewControllerManagedBy(mgr).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(zenDeploymentPredicate()).
		Watches(&source.Kind{Type: &corev1.ConfigMap{}}, handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
			return []ctrl.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "RECONCILE-ZEN-PRODUCT-CONFIGMAP",
					Namespace: a.GetNamespace(),
				}},
			}
		})).
		WithEventFilter(zenProductCmPredicate()).
		Complete(r)
}
