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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
	res "github.com/IBM/ibm-commonui-operator/controllers/resources"
)

var log = logf.Log.WithName("controller_commonwebui")

// CommonWebUIReconciler reconciles a CommonWebUI object
type CommonWebUIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operators.ibm.com,resources=commonwebuis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operators.ibm.com,resources=commonwebuis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operators.ibm.com,resources=commonwebuis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CommonWebUI object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CommonWebUIReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUI")

	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	// Fetch the CommonWebUIService CR instance
	instance := &operatorsv1alpha1.CommonWebUI{}

	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	reqLogger.Info("CommonWebUI instance version: " + instance.Spec.OperatorVersion)

	// set a default Status value
	if len(instance.Status.Nodes) == 0 {
		instance.Status.Nodes = res.DefaultStatusForCR
		// instance.Status.Versions = operatorsv1alpha1.Versions{Reconciled: ver.Version}
		err = r.Client.Status().Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set CommonWebUI default status")
			return ctrl.Result{}, err
		}
	}

	// Check to see if Zen instance exists in common services namespace
	isZen := res.IsAdminHubOnZen(ctx, r.Client, instance.Namespace)

	// Check to see kubernetes cluster type is cncf
	isCncf := res.GetKubernetesClusterType(ctx, r.Client, instance.Namespace)

	// Check if the log4js configmap already exists. If not, create a new one.
	err = res.ReconcileLog4jsConfigMap(ctx, r.Client, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists. If not, create a new one.
	err = res.ReconcileDeployment(ctx, r.Client, instance, isZen, isCncf, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CommonWebUIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsv1alpha1.CommonWebUI{}).
		Complete(r)
}
