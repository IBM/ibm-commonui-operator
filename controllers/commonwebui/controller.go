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
	"fmt"
	"reflect"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
	certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
)

var log = logf.Log.WithName("controller_commonwebui")

// CommonWebUIReconciler reconciles a CommonWebUI object
type CommonWebUIReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *CommonWebUIReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUI")

	// Set a flag to requeue only once in the instance that multiple resources are created
	needToRequeue := false

	instance := &operatorv1alpha1.CommonWebUI{}

	// Fetch the CommonWebUIService CR instance
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

	// Set a default status value
	if len(instance.Status.Nodes) == 0 {
		instance.Status.Nodes = DefaultStatusForCR
		err = r.Client.Status().Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set CommonWebUI default status")
			return ctrl.Result{}, err
		}
	}

	adminHubOnZen := r.hasAdminHubOnZen(ctx, instance)

	// Check if the log4js configmap already exists. If not, create a new one.
	err = r.reconcileLog4jsConfigMap(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the redis certs configmap already exists. If not, create a new one.
	err = r.reconcileRedisConfigMap(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// If admin hub on zen is enabled, reconcile zen card extensions configmap.
	// Otherwise, reconcile zen left menu item extension.
	if adminHubOnZen {
		// Check if the zen card extensions configmap already exists. If not, create a new one.
		err = r.reconcileZenCardsConfigMap(ctx, instance, &needToRequeue)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// Check if the zen left menu extensions configmap already exists. If not, create a new one.
		err = r.reconcileExtensionsConfigMap(ctx, instance, &needToRequeue)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Duration(3) * time.Minute}, err
		}
	}

	// Check if the deployment already exists. If not, create a new one.
	err = r.reconcileDeployment(ctx, instance, adminHubOnZen, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the service already exists. If not, create a new one.
	err = r.reconcileService(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the API ingress already exists. If not, create a new one.
	err = r.reconcileAPIIngress(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the callback ingress already exists. If not, create a new one.
	err = r.reconcileCallbackIngress(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the common-nav ingress already exists. If not, create a new one.
	err = r.reconcileNavIngress(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the ConsoleLink CR already exists. If not, create a new one.
	err = r.reconcileConsoleLink(ctx, instance, adminHubOnZen, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the certificates already exists. If not, create new ones.
	err = r.reconcileCertificates(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Update admin hub nav config, if it exists.
	err = r.reconcileAdminHubNavConfig(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Update cp4i nav config, if it exists.
	err = r.reconcileCP4INavConfig(ctx, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// For 1.3.0 operator version, check if daemonset exits on upgrade and delete if so
	r.deleteDaemonSet(ctx, instance)

	if needToRequeue {
		// One or more resources were created, so requeue the request
		reqLogger.Info("Requeuing the request")
		return ctrl.Result{Requeue: true}, nil
	}

	err = r.updateStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *CommonWebUIReconciler) deleteDaemonSet(ctx context.Context, instance *operatorv1alpha1.CommonWebUI) {
	reqLogger := log.WithValues("func", "deleteDaemonSet", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: DaemonSetName,
			Namespace: DefaultNamespace,
		},
	}

	err := r.Client.Get(ctx, types.NamespacedName{Name: daemonSet.Name, Namespace: daemonSet.Namespace}, daemonSet)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("No CommonWebUI daemonset found")
	} else {
		// Delete daemonset if found
		err = r.Client.Delete(ctx, daemonSet)
		if err != nil {
			reqLogger.Error(err, "Failed to delete old CommonWebUI daemonset")
		} else {
			reqLogger.Info("Successfully deleted old CommonWebUI daemonset")
		}
	}
}

func (r *CommonWebUIReconciler) updateStatus(ctx context.Context, instance *operatorv1alpha1.CommonWebUI) error {
	reqLogger := log.WithValues("func", "updateStatus", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Updating CommonWebUI status")
	
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(LabelsForSelector(DeploymentName, CommonWebUICRType, instance.Name)),
	}

	err := r.Client.List(ctx, podList, listOpts...)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods")
		return err
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}

	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames

		err := r.Client.Status().Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Failed to updated CommonWebUI status")
			return err
		}
	}
	
	return nil
}

func (r *CommonWebUIReconciler) hasAdminHubOnZen(ctx context.Context, instance *operatorv1alpha1.CommonWebUI) bool {
	reqLogger := log.WithValues("func", "hasAdminHubOnZen", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Checking zen optional install condition")

	cm := &corev1.ConfigMap{}

	err := r.Client.Get(ctx, types.NamespacedName{Name: AdminHubOnZenConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		reqLogger.Info("Unable to find Zen optional install configmap")
	} else {
		reqLogger.Info(fmt.Sprintf("Zen optional install configmap found with value: %s", cm.Data["adminHubOnZen"]))
		return cm.Data["adminHubOnZen"] == "true"
	}

	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *CommonWebUIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.CommonWebUI{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Owns(&certmgr.Certificate{}).
		Complete(r)
}
