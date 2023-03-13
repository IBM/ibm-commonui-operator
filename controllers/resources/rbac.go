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
package resources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

const ServiceAccountName = "ibm-commonui-operand"
const OperandRoleName = "ibm-commonui-operand"
const OperandRoleBindingName = "ibm-commonui-operand"

func getDesiredServiceAccount(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*corev1.ServiceAccount, error) {
	reqLogger := log.WithValues("func", "getDesiredServiceAccount", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceAccountName,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/instance":   "ibm-commonui-operator",
				"app.kubernetes.io/name":       "ibm-commonui-operator",
				"app.kubernetes.io/managed-by": "ibm-commonui-operator",
			},
		},
	}

	err := controllerutil.SetControllerReference(instance, serviceAccount, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for service account")
		return nil, err
	}

	return serviceAccount, nil
}

func ReconcileServiceAccount(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileServiceAccount", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling service account")

	serviceAccount := &corev1.ServiceAccount{}

	desiredSA, desiredErr := getDesiredServiceAccount(client, instance)
	if desiredErr != nil {
		return desiredErr
	}

	err := client.Get(ctx, types.NamespacedName{Name: ServiceAccountName, Namespace: desiredSA.Namespace}, serviceAccount)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new service account", "SA.Namespace", desiredSA.Namespace, "SA.Name", desiredSA.Name)

		err = client.Create(ctx, desiredSA)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// SA already exists from a previous reconcile
				reqLogger.Info("Service account already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new SA
				reqLogger.Info("Failed to create a new service account", "SA.Namespace", desiredSA.Namespace, "SA.Name", desiredSA.Name)
				return err
			}
		} else {
			// Requeue after creating new service account
			*needToRequeue = true
		}
	} else if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get service account", "SA.Namespace", desiredSA.Namespace, "SA.Name", ServiceAccountName)
		return err
	} else {
		// Determine if current service account has changed
		reqLogger.Info("Comparing current and desired service accounts")

		if !IsServiceAccountEqual(serviceAccount, desiredSA) {
			reqLogger.Info("Updating service account", "SA.Namespace", desiredSA.Namespace, "SA.Name", desiredSA.Name)

			serviceAccount.ObjectMeta.Name = desiredSA.ObjectMeta.Name
			serviceAccount.ObjectMeta.Labels = desiredSA.ObjectMeta.Annotations

			err = client.Update(ctx, serviceAccount)
			if err != nil {
				reqLogger.Error(err, "Failed to update service account", "SA.Namespace", desiredSA.Namespace, "SA.Name", desiredSA.Name)
				return err
			}
		}
	}

	return nil
}

func getDesiredRole(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*rbacv1.Role, error) {
	reqLogger := log.WithValues("func", "getDesiredRole", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      OperandRoleName,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/instance":   "ibm-commonui-operator",
				"app.kubernetes.io/name":       "ibm-commonui-operator",
				"app.kubernetes.io/managed-by": "ibm-commonui-operator",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{"foundation.ibm.com"},
				Resources: []string{"navconfigurations"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{"operators.ibm.com"},
				Resources: []string{"switcheritems"},
				Verbs:     []string{"get", "list"},
			},
		},
	}

	err := controllerutil.SetControllerReference(instance, role, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for role")
		return nil, err
	}

	return role, nil
}

func ReconcileRole(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileRole", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling role")

	role := &rbacv1.Role{}

	desiredRole, desiredErr := getDesiredRole(client, instance)
	if desiredErr != nil {
		return desiredErr
	}

	err := client.Get(ctx, types.NamespacedName{Name: desiredRole.Name, Namespace: desiredRole.Namespace}, role)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new role", "Role.Namespace", desiredRole.Namespace, "Role.Name", desiredRole.Name)

		err = client.Create(ctx, desiredRole)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// role already exists from a previous reconcile
				reqLogger.Info("Role already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new role
				reqLogger.Info("Failed to create a new role", "Role.Namespace", desiredRole.Namespace, "Role.Name", desiredRole.Name)
				return err
			}
		} else {
			// Requeue after creating new role
			*needToRequeue = true
		}
	} else if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get role", "Role.Namespace", desiredRole.Namespace, "Role.Name", desiredRole.Name)
		return err
	} else {
		// Determine if current role has changed
		reqLogger.Info("Comparing current and desired role")

		if !IsRoleEqual(role, desiredRole) {
			reqLogger.Info("Updating role", "Role.Namespace", desiredRole.Namespace, "Role.Name", desiredRole.Name)

			role.ObjectMeta.Name = desiredRole.ObjectMeta.Name
			role.ObjectMeta.Labels = desiredRole.ObjectMeta.Annotations
			role.Rules = desiredRole.Rules

			err = client.Update(ctx, role)
			if err != nil {
				reqLogger.Error(err, "Failed to update role", "Role.Namespace", desiredRole.Namespace, "Role.Name", desiredRole.Name)
				return err
			}
		}
	}

	return nil
}

func getDesiredRoleBinding(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*rbacv1.RoleBinding, error) {
	reqLogger := log.WithValues("func", "getDesiredRoleBinding", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      OperandRoleBindingName,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/instance":   "ibm-commonui-operator",
				"app.kubernetes.io/name":       "ibm-commonui-operator",
				"app.kubernetes.io/managed-by": "ibm-commonui-operator",
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      ServiceAccountName,
				Namespace: instance.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     OperandRoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	err := controllerutil.SetControllerReference(instance, roleBinding, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for role binding")
		return nil, err
	}

	return roleBinding, nil
}

func ReconcileRoleBinding(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileRoleBinding", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling rolebinding")

	roleBinding := &rbacv1.RoleBinding{}

	desiredRoleBinding, desiredErr := getDesiredRoleBinding(client, instance)
	if desiredErr != nil {
		return desiredErr
	}

	err := client.Get(ctx, types.NamespacedName{Name: desiredRoleBinding.Name, Namespace: desiredRoleBinding.Namespace}, roleBinding)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new role binding", "RoleBinding.Namespace", desiredRoleBinding.Namespace, "RoleBinding.Name", desiredRoleBinding.Name)

		err = client.Create(ctx, desiredRoleBinding)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// role binding already exists from a previous reconcile
				reqLogger.Info("Role binding already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new role binding
				reqLogger.Info("Failed to create a new role binding", "RoleBinding.Namespace", desiredRoleBinding.Namespace, "RoleBinding.Name", desiredRoleBinding.Name)
				return err
			}
		} else {
			// Requeue after creating new role binding
			*needToRequeue = true
		}
	} else if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get role binding", "RoleBinding.Namespace", desiredRoleBinding.Namespace, "RoleBinding.Name", desiredRoleBinding.Name)
		return err
	} else {
		// Determine if current role binding has changed
		reqLogger.Info("Comparing current and desired role binding")

		if !IsRoleBindingEqual(roleBinding, desiredRoleBinding) {
			reqLogger.Info("Updating role binding", "RoleBinding.Namespace", desiredRoleBinding.Namespace, "RoleBinding.Name", desiredRoleBinding.Name)

			roleBinding.ObjectMeta.Name = desiredRoleBinding.ObjectMeta.Name
			roleBinding.ObjectMeta.Labels = desiredRoleBinding.ObjectMeta.Annotations
			roleBinding.Subjects = desiredRoleBinding.Subjects
			roleBinding.RoleRef = desiredRoleBinding.RoleRef

			err = client.Update(ctx, roleBinding)
			if err != nil {
				reqLogger.Error(err, "Failed to update role binding", "RoleBinding.Namespace", desiredRoleBinding.Namespace, "RoleBinding.Name", desiredRoleBinding.Name)
				return err
			}
		}
	}

	return nil
}
