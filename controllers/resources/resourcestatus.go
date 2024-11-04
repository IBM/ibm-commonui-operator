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

package resources

import (
	"context"

	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

const Ready = "Ready"
const NotReady = "NotReady"
const Unknown = "Unknown"

func getServiceStatus(ctx context.Context, k8sClient client.Client, namespacedName types.NamespacedName) (status v1alpha1.ManagedResourceStatus) {
	reqLogger := log.WithValues("func", "getServiceStatus", "namespacedName", namespacedName)

	status = v1alpha1.ManagedResourceStatus{
		ObjectName: namespacedName.Name,
		APIVersion: Unknown,
		Namespace:  namespacedName.Namespace,
		Kind:       "Service",
		Status:     NotReady,
	}
	service := &corev1.Service{}
	err := k8sClient.Get(ctx, namespacedName, service)

	if err != nil {
		if !errors.IsNotFound(err) {
			reqLogger.Error(err, "Error reading service for status update")
		}
		return
	}
	status.APIVersion = service.APIVersion
	status.Status = Ready
	return
}

func getDeploymentStatus(ctx context.Context, k8sClient client.Client, namespacedName types.NamespacedName) (status v1alpha1.ManagedResourceStatus) {
	reqLogger := log.WithValues("func", "getDeploymentStatus", "namespacedName", namespacedName)

	status = v1alpha1.ManagedResourceStatus{
		ObjectName: namespacedName.Name,
		APIVersion: Unknown,
		Namespace:  namespacedName.Namespace,
		Kind:       "Deployment",
		Status:     NotReady,
	}

	deployment := &appsv1.Deployment{}
	err := k8sClient.Get(ctx, namespacedName, deployment)

	if err != nil {
		if !errors.IsNotFound(err) {
			reqLogger.Error(err, "Error reading deployment for status update")
		}
		return
	}

	status.APIVersion = deployment.APIVersion

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionTrue {
			status.Status = Ready
			return
		}
	}
	return
}

func getRouteStatus(ctx context.Context, k8sClient client.Client, namespacedName types.NamespacedName) (status v1alpha1.ManagedResourceStatus) {
	reqLogger := log.WithValues("func", "getRouteStatus", "namespacedName", namespacedName)

	status = v1alpha1.ManagedResourceStatus{
		ObjectName: namespacedName.Name,
		APIVersion: Unknown,
		Namespace:  namespacedName.Namespace,
		Kind:       "Route",
		Status:     NotReady,
	}
	route := &routev1.Route{}
	err := k8sClient.Get(ctx, namespacedName, route)

	if err != nil {
		if !errors.IsNotFound(err) {
			reqLogger.Error(err, "Error reading route for status update")
		}
		return
	}

	status.APIVersion = route.APIVersion

	for _, routeIngress := range route.Status.Ingress {
		for _, condition := range routeIngress.Conditions {
			if condition.Type == routev1.RouteAdmitted && condition.Status != corev1.ConditionTrue {
				return
			}
		}
	}
	status.Status = Ready
	return
}

type statusRetrievalFunc func(context.Context, client.Client, []string, string) []v1alpha1.ManagedResourceStatus

func getAllServiceStatus(ctx context.Context, k8sClient client.Client, names []string, namespace string) (statuses []v1alpha1.ManagedResourceStatus) {
	reqLogger := log.WithValues("func", "getAllServiceStatus", "namespace", namespace)
	for _, name := range names {
		nsn := types.NamespacedName{Name: name, Namespace: namespace}
		statuses = append(statuses, getServiceStatus(ctx, k8sClient, nsn))
	}
	reqLogger.Info("New statuses", "statuses", statuses)
	return
}

func getAllDeploymentStatus(ctx context.Context, k8sClient client.Client, names []string, namespace string) (statuses []v1alpha1.ManagedResourceStatus) {
	reqLogger := log.WithValues("func", "getAllDeploymentStatus", "namespace", namespace)
	for _, name := range names {
		nsn := types.NamespacedName{Name: name, Namespace: namespace}
		statuses = append(statuses, getDeploymentStatus(ctx, k8sClient, nsn))
	}
	reqLogger.Info("New statuses", "statuses", statuses)
	return
}

func getAllRouteStatus(ctx context.Context, k8sClient client.Client, names []string, namespace string) (statuses []v1alpha1.ManagedResourceStatus) {
	reqLogger := log.WithValues("func", "getAllRouteStatus", "namespace", namespace)
	for _, name := range names {
		nsn := types.NamespacedName{Name: name, Namespace: namespace}
		statuses = append(statuses, getRouteStatus(ctx, k8sClient, nsn))
	}
	reqLogger.Info("New statuses", "statuses", statuses)
	return
}

func GetCurrentServiceStatus(ctx context.Context, k8sClient client.Client, instance *v1alpha1.CommonWebUI, isCncf bool) (status v1alpha1.ServiceStatus) {
	reqLogger := log.WithValues("func", "getCurrentServiceStatus", "namespace", instance.Namespace, "isCncf", isCncf)
	type statusRetrieval struct {
		names []string
		f     statusRetrievalFunc
	}

	//
	statusRetrievals := []statusRetrieval{
		{
			names: []string{
				"common-web-ui",
			},
			f: getAllServiceStatus,
		},
		{
			names: []string{
				"common-web-ui",
			},
			f: getAllDeploymentStatus,
		},
	}

	routeStatusRetrieval := statusRetrieval{
		names: []string{
			"cp-console",
		},
		f: getAllRouteStatus,
	}

	if !isCncf && !ZenFrontDoorEnabled(ctx, k8sClient, instance.Namespace) {
		statusRetrievals = append(statusRetrievals, routeStatusRetrieval)
	}

	status = v1alpha1.ServiceStatus{
		ObjectName:       instance.Name,
		Namespace:        instance.Namespace,
		APIVersion:       instance.APIVersion,
		Kind:             "CommonWebUI",
		ManagedResources: []v1alpha1.ManagedResourceStatus{},
		Status:           NotReady,
	}

	reqLogger.Info("Getting statuses")
	for _, getStatuses := range statusRetrievals {
		status.ManagedResources = append(status.ManagedResources, getStatuses.f(ctx, k8sClient, getStatuses.names, status.Namespace)...)
	}

	for _, managedResourceStatus := range status.ManagedResources {
		if managedResourceStatus.Status == NotReady {
			return
		}
	}
	status.Status = Ready
	return
}
