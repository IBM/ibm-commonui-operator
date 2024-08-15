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
	"reflect"
	"strings"

	certmgr "github.com/ibm/ibm-cert-manager-operator/apis/cert-manager/v1"
	certmgrv1alpha1 "github.com/ibm/ibm-cert-manager-operator/apis/certmanager/v1alpha1"
	route "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
	im "github.com/IBM/ibm-commonui-operator/apis/operator/v1alpha1"
	res "github.com/IBM/ibm-commonui-operator/controllers/resources"
	"github.com/IBM/ibm-commonui-operator/version"
)

var log = logf.Log.WithName("controller_commonwebui")

// CommonWebUIReconciler reconciles a CommonWebUI object
type CommonWebUIReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
	IsCncf bool
}

const finalizerName = "commonui.operators.ibm.com"
const finalizerName1 = "commonui1.operators.ibm.com"

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
	reqLogger.Info("Reconciling CommonWebUI Controller")

	var err error

	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	// Fetch the CommonWebUIService CR instance
	instance := &operatorsv1alpha1.CommonWebUI{}

	//If the ibmcloud-cluster-info configmap has been updated then we need to reconcile routes
	//Since this isn't owned by our CR, we need to look our CR up
	if request.Name == "NON_OWNED_OBJECT_RECONCILE" {
		crList := &operatorsv1alpha1.CommonWebUIList{}
		err := r.Client.List(ctx, crList, client.InNamespace(instance.Namespace))
		if err != nil || len(crList.Items) == 0 {
			reqLogger.Error(err, "Cluster config configmap has changed, but unable to load list of CommonWebUI CRs")
			return ctrl.Result{}, err
		}
		instance = &crList.Items[0]
	} else {
		err = r.Client.Get(ctx, request.NamespacedName, instance)
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
	}

	reqLogger.Info("CommonWebUI instance version: " + instance.Spec.OperatorVersion)

	//Setup status update before returning
	defer func() {
		err := r.updateStatus(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Error updating current CR status")
		}
	}()

	// set a default Status value
	if len(instance.Status.Nodes) == 0 {
		instance.Status.Nodes = res.DefaultStatusForCR
		instance.Status.OperatorVersion = version.Version
		instance.Status.OperandVersion = version.Version
		err = r.Client.Status().Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set CommonWebUI default status")
			return ctrl.Result{}, err
		}
	}

	// Check to see if Zen instance exists in common services namespace
	isZen := false //ZEN DISABLED res.IsAdminHubOnZen(ctx, r.Client, instance.Namespace)

	// Check to see kubernetes cluster type is cncf
	isCncf := r.IsCncf

	// Check if the log4js configmap already exists. If not, create a new one.
	err = res.ReconcileLog4jsConfigMap(ctx, r.Client, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the common-web-ui-config configmap already exists. If not, create a new one.
	err = res.ReconcileCommonUIConfigConfigMap(ctx, r.Client, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = res.ReconcileServiceAccount(ctx, r.Client, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists. If not, create a new one.
	err = res.ReconcileDeployment(ctx, r.Client, instance, isZen, isCncf, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the service already exists. If not, create a new one.
	err = res.ReconcileService(ctx, r.Client, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile the required routes if this is not a cncf cluster
	if !isCncf {
		err = res.ReconcileRoutes(ctx, r.Client, instance, &needToRequeue)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Remove any legacy ingresses if they are found (this would only be on migration)
	// This was for cloudpak 3.0 work
	res.ReconcileRemoveIngresses(ctx, r.Client, instance, &needToRequeue)

	// For 1.15.0 operator version, check if v1alpha1 certs exits on upgrade and delete if so
	r.deleteCertsv1alpha1(ctx, instance)

	// Check if the certificates already exists. If not, create new v1 certs.
	err = res.ReconcileCertificates(ctx, r.Client, instance, &needToRequeue)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Update admin hub nav config, if it exists.
	err = res.ReconcileAdminHubNavConfig(ctx, r.Client, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Cleanup any remaining zen artifacts after removal of adminhub
	r.removeLegacyZenResources(ctx, instance)

	//Delete the operand request that may have been created by common ui prior to upgrade to cp 3.0
	//nolint
	res.DeleteGenericResource(ctx, "ibm-commonui-request", instance.Namespace, "operator.ibm.com", "v1alpha1", "operandrequests")

	r.removeLegacyFinalizers(ctx, instance)

	if needToRequeue {
		// One or more resources were created, so requeue the request
		reqLogger.Info("Requeuing the request")
		return ctrl.Result{Requeue: true}, nil
	}

	reqLogger.Info("COMMON UI CONTROLLER RECONCILE ALL DONE")
	return ctrl.Result{}, nil
}

func (r *CommonWebUIReconciler) removeLegacyZenResources(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI) {
	reqLogger := log.WithValues("func", "removeLegacyZenResources", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Removing legacy classic admin hub resources for zen")

	//Delete common ui bind info config map
	//nolint
	res.DeleteConfigMap(ctx, r.Client, "ibm-commonui-bindinfo-common-webui-ui-extensions", instance.Namespace)

	//Delete classic admin hub left nav menu item
	//nolint
	res.DeleteConfigMap(ctx, r.Client, res.ZenLeftNavExtensionsConfigMapName, instance.Namespace)

	//Delete zen adminhub card extensions
	//nolint
	res.DeleteConfigMap(ctx, r.Client, res.ZenCardExtensionsConfigMapName, instance.Namespace)

	//Delete zen adminhub quick nav extensions
	//nolint
	res.DeleteConfigMap(ctx, r.Client, res.ZenQuickNavExtensionsConfigMapName, instance.Namespace)

}

// Common UI 3.x added finalizers to the UI CR to manage OCP console links (and maybe redis)
// These are causing issues on 4.x upgrade because the finalizers still exist, however the code
// that would process them is long removed (this is because console links require cluster permissions
// and were essentially abandoned as objects in 4.x - customer must remove them if one exists)
func (r *CommonWebUIReconciler) removeLegacyFinalizers(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI) {
	reqLogger := log.WithValues("func", "removeLegacyFinalizers", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Checking for legacy finalizers for removal")

	hasFinalizer := res.ContainsString(instance.ObjectMeta.Finalizers, "commonui.operators.ibm.com")
	hasFinalizer1 := res.ContainsString(instance.ObjectMeta.Finalizers, "commonui1.operators.ibm.com")

	if hasFinalizer || hasFinalizer1 {
		if hasFinalizer {
			instance.ObjectMeta.Finalizers = res.RemoveString(instance.ObjectMeta.Finalizers, finalizerName)
			reqLogger.Info("Removing finalizer " + finalizerName)
		}
		if hasFinalizer1 {
			instance.ObjectMeta.Finalizers = res.RemoveString(instance.ObjectMeta.Finalizers, finalizerName1)
			reqLogger.Info("Removing finalizer " + finalizerName1)
		}

		if err := r.Client.Update(ctx, instance); err != nil {
			reqLogger.Error(err, "Failed to update after removing finalizer")
		}
	}
}

func (r *CommonWebUIReconciler) deleteCertsv1alpha1(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI) {
	reqLogger := log.WithValues("func", "deleteCertsv1alpha1", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	certificate := &certmgrv1alpha1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.UICertName,
			Namespace: instance.Namespace,
		},
	}
	err := r.Client.Get(ctx, types.NamespacedName{Name: res.UICertName, Namespace: instance.Namespace}, certificate)

	if err != nil {
		if !errors.IsNotFound(err) {
			reqLogger.Info("Unable to load v1alpha1 certificate - most likely this means the CRD doesn't exist and this can be ignored")
		}
		return
	}
	reqLogger.Info("Certificate common-web-ui-ca-cert found, checking api version..")
	reqLogger.Info("API version is: " + certificate.APIVersion)
	if certificate.APIVersion == res.Certv1alpha1APIVersion {
		reqLogger.Info("deleting cert: " + res.UICertName)
		err = r.Client.Delete(ctx, certificate)
		if err != nil {
			reqLogger.Error(err, "Failed to delete")
		} else {
			reqLogger.Info("Successfully deleted")
		}
	} else {
		reqLogger.Info("API version is NOT v1alpha1, returning..")
	}
}

func (r *CommonWebUIReconciler) updateStatus(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI) error {
	reqLogger := log.WithValues("func", "updateStatus", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Updating CommonWebUI status")

	updateServiceStatus := false
	updateNodeStatus := false

	//Check for updates to service status
	reqLogger.Info("Gather current service status")
	currentServiceStatus := res.GetCurrentServiceStatus(ctx, r.Client, instance, r.IsCncf)
	if !reflect.DeepEqual(currentServiceStatus, instance.Status.Service) {
		instance.Status.Service = currentServiceStatus
		updateServiceStatus = true
	}

	//Check for updates to node (pods) status
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(res.LabelsForSelector(res.DeploymentName, res.CommonWebUICRType, instance.Name)),
	}

	err := r.Client.List(ctx, podList, listOpts...)
	if err == nil {
		var podNames []string
		for _, pod := range podList.Items {
			podNames = append(podNames, pod.Name)
		}

		if !reflect.DeepEqual(podNames, instance.Status.Nodes) || instance.Status.OperatorVersion != version.Version ||
			instance.Status.OperandVersion != version.Version {
			instance.Status.Nodes = podNames
			instance.Status.OperatorVersion = version.Version
			instance.Status.OperandVersion = version.Version
			updateNodeStatus = true
		}
	} else {
		reqLogger.Error(err, "Failed to list pods - CR status will not be updated")
	}

	//Update any serivce status updates
	if updateServiceStatus || updateNodeStatus {
		reqLogger.Info("Updating status", "updateServiceStatus", updateServiceStatus, "updateNodeStatus", updateNodeStatus)
		err := r.Client.Status().Update(ctx, instance)
		if err != nil {
			return err
		}
	} else {
		reqLogger.Info("NO STATUS UPDATE REQUIRED - RECONCILE COMPLETE")
	}

	return nil
}

func clusterInfoCmPredicate() predicate.Predicate {
	namespaces := strings.Split(os.Getenv("WATCH_NAMESPACE"), ",")

	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew.GetName() == res.ClusterInfoConfigmapName && res.ContainsString(namespaces, e.ObjectNew.GetNamespace()) {
				return true
			}
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			if e.Object.GetName() == res.ClusterInfoConfigmapName && res.ContainsString(namespaces, e.Object.GetNamespace()) {
				return true
			}
			return false
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *CommonWebUIReconciler) SetupWithManager(mgr ctrl.Manager) error {

	//Skip routes when it is cncf
	if r.IsCncf {
		return ctrl.NewControllerManagedBy(mgr).
			For(&operatorsv1alpha1.CommonWebUI{}).
			Owns(&corev1.ConfigMap{}).
			Owns(&appsv1.Deployment{}).
			Owns(&corev1.Service{}).
			Owns(&corev1.Secret{}).
			Owns(&netv1.Ingress{}).
			Owns(&certmgr.Certificate{}).
			Owns(&corev1.ServiceAccount{}).
			Owns(&rbacv1.Role{}).
			Owns(&rbacv1.RoleBinding{}).
			Watches(&source.Kind{Type: &corev1.ConfigMap{}},
				handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
					return []ctrl.Request{
						{NamespacedName: types.NamespacedName{
							Name:      "NON_OWNED_OBJECT_RECONCILE",
							Namespace: a.GetNamespace(),
						}},
					}
				}), builder.WithPredicates(clusterInfoCmPredicate())).
			Complete(r)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsv1alpha1.CommonWebUI{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&netv1.Ingress{}).
		Owns(&certmgr.Certificate{}).
		Owns(&route.Route{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Watches(&source.Kind{Type: &corev1.ConfigMap{}},
			handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
				return []ctrl.Request{
					{NamespacedName: types.NamespacedName{
						Name:      "NON_OWNED_OBJECT_RECONCILE",
						Namespace: a.GetNamespace(),
					}},
				}
			}), builder.WithPredicates(clusterInfoCmPredicate())).
		Watches(&source.Kind{Type: &im.Authentication{}},
			handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
				return []ctrl.Request{
					{NamespacedName: types.NamespacedName{
						Name:      "NON_OWNED_OBJECT_RECONCILE",
						Namespace: a.GetNamespace(),
					}},
				}
			})).
		Complete(r)
}
