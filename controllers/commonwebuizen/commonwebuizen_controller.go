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
	"os"

	res "github.com/IBM/ibm-commonui-operator/controllers/resources"
	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

	return ctrl.Result{}, nil
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
		Watches(&source.Kind{Type: &corev1.ConfigMap{}}, handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      "RECONCILE-ZEN-PRODUCT-CONFIGMAP",
					Namespace: a.GetNamespace(),
				}},
			}
		})).
		WithEventFilter(zenProductCmPredicate()).
		Complete(r)
}
