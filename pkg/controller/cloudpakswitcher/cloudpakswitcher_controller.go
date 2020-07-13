//
// Copyright 2020 IBM Corporation
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

package cloudpakswitcher

import (
	"context"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_cloudpakswitcher")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CloudPakSwitcher Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCloudPakSwitcher{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cloudpakswitcher-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CloudPakSwitcher
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.CloudPakSwitcher{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CloudPakSwitcher
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CloudPakSwitcher{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCloudPakSwitcher implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCloudPakSwitcher{}

// ReconcileCloudPakSwitcher reconciles a CloudPakSwitcher object
type ReconcileCloudPakSwitcher struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CloudPakSwitcher object and makes changes based on the state read
// and what is in the CloudPakSwitcher.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCloudPakSwitcher) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CloudPakSwitcher")

	// Fetch the CloudPakSwitcher instance
	instance := &operatorsv1alpha1.CloudPakSwitcher{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	currentClusterRole := &rbacv1.ClusterRole{}
	recResult, err := r.handleClusterRole(instance, currentClusterRole)
	if err != nil {
		return recResult, err
	}

	return reconcile.Result{}, nil
}

//nolint
func (r *ReconcileCloudPakSwitcher) handleClusterRole(instance *operatorsv1alpha1.CloudPakSwitcher, currentClusterRole *rbacv1.ClusterRole) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "icp-cps-admin-aggregate", Namespace: ""}, currentClusterRole)
	if err != nil && errors.IsNotFound(err) {
		// Define admin cluster role
		adminClusterRole := r.adminClusterRoleForCloudPakSwitcher(instance)
		reqLogger.Info("Creating a new ClusterRole", "ClusterRole.Namespace", instance.Namespace, "ClusterRole.Name", "icp-cps-admin-aggregate")
		err = r.client.Create(context.TODO(), adminClusterRole)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ClusterRole", "ClusterRole.Namespace", instance.Namespace, "ClusterRole.Name", "icp-cps-admin-aggregate")
			return reconcile.Result{}, err
		}

	} else if err != nil {
		reqLogger.Error(err, "Failed to get ClusterRole")
		return reconcile.Result{}, err
	}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "icp-cps-operate-aggregate", Namespace: ""}, currentClusterRole)
	if err != nil && errors.IsNotFound(err) {
		// Define operator cluster role
		operatorClusterRole := r.operatorClusterRoleForCloudPakSwitcher(instance)
		reqLogger.Info("Creating a new ClusterRole", "ClusterRole.Namespace", instance.Namespace, "ClusterRole.Name", "icp-cps-operate-aggregate")
		err = r.client.Create(context.TODO(), operatorClusterRole)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ClusterRole", "ClusterRole.Namespace", instance.Namespace, "ClusterRole.Name", "icp-cps-operate-aggregate")
			return reconcile.Result{}, err
		}
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ClusterRole")
		return reconcile.Result{}, err
	}
	//admin roles created successfully
	return reconcile.Result{Requeue: true}, nil
}

//nolint
func (r *ReconcileCloudPakSwitcher) adminClusterRoleForCloudPakSwitcher(instance *operatorsv1alpha1.CloudPakSwitcher) *rbacv1.ClusterRole {
	// reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	adminClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "icp-cps-admin-aggregate",
			Labels: map[string]string{
				"kubernetes.io/bootstrapping":                  "rbac-defaults",
				"rbac.icp.com/aggregate-to-icp-admin":          "true",
				"rbac.authorization.k8s.io/aggregate-to-admin": "true",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "list"},
			},
		},
	}
	// Set OIDCClientWatcher instance as the owner and controller of the cluster role
	// err := controllerutil.SetControllerReference(instance, adminClusterRole, r.scheme)
	// if err != nil {
	// 	reqLogger.Error(err, "Failed to set owner for admin Cluster Role")
	// 	return nil
	// }
	return adminClusterRole
}

//nolint
func (r *ReconcileCloudPakSwitcher) operatorClusterRoleForCloudPakSwitcher(instance *operatorsv1alpha1.CloudPakSwitcher) *rbacv1.ClusterRole {
	// reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	operatorClusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "icp-cps-operate-aggregate",
			Labels: map[string]string{
				"kubernetes.io/bootstrapping":                 "rbac-defaults",
				"rbac.icp.com/aggregate-to-icp-operate":       "true",
				"rbac.authorization.k8s.io/aggregate-to-edit": "true",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "list"},
			},
		},
	}
	// Set OIDCClientWatcher instance as the owner and controller of the cluster role
	// err := controllerutil.SetControllerReference(instance, operatorClusterRole, r.scheme)
	// if err != nil {
	// 	reqLogger.Error(err, "Failed to set owner for operator Cluster Role")
	// 	return nil
	// }
	return operatorClusterRole
}
