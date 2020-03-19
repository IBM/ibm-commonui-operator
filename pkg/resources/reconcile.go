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
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Check if a DaemonSet already exists. If not, create a new one.
func ReconcileDaemonSet(client client.Client, instanceNamespace string, daemonSetName string,
	newDaemonSet *appsv1.DaemonSet, needToRequeue *bool) error {
	logger := log.WithValues("func", "ReconcileDaemonSet")

	currentDaemonSet := &appsv1.DaemonSet{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: daemonSetName, Namespace: instanceNamespace}, currentDaemonSet)
	if err != nil && errors.IsNotFound(err) {
		// Create a new DaemonSet
		logger.Info("Creating a new DaemonSet", "DaemonSet.Namespace", newDaemonSet.Namespace, "DaemonSet.Name", newDaemonSet.Name)
		err = client.Create(context.TODO(), newDaemonSet)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue
			logger.Info(" DaemonSet already exists")
			*needToRequeue = true
		} else if err != nil {
			logger.Error(err, "Failed to create new DaemonSet", "DaemonSet.Namespace", newDaemonSet.Namespace,
				"DaemonSet.Name", newDaemonSet.Name)
			return err
		} else {
			// DaemonSet created successfully - return and requeue
			*needToRequeue = true
		}
	} else if err != nil {
		logger.Error(err, "Failed to get DaemonSet", "DaemonSet.Name", daemonSetName)
		return err
	} else {
		// Found DaemonSet, so determine if the resource has changed
		logger.Info("Comparing DaemonSets")
		if !IsDaemonSetEqual(currentDaemonSet, newDaemonSet) {
			logger.Info("Updating DaemonSet", "DaemonSet.Name", currentDaemonSet.Name)
			currentDaemonSet.ObjectMeta.Name = newDaemonSet.ObjectMeta.Name
			currentDaemonSet.ObjectMeta.Labels = newDaemonSet.ObjectMeta.Labels
			currentDaemonSet.Spec = newDaemonSet.Spec
			err = client.Update(context.TODO(), currentDaemonSet)
			if err != nil {
				logger.Error(err, "Failed to update DaemonSet",
					"DaemonSet.Namespace", currentDaemonSet.Namespace, "DaemonSet.Name", currentDaemonSet.Name)
				return err
			}
		}
	}
	return nil
}

// Check if a Service already exists. If not, create a new one.
func ReconcileService(client client.Client, instanceNamespace string, serviceName string,
	newService *corev1.Service, needToRequeue *bool) error {
	logger := log.WithValues("func", "ReconcileService")

	currentService := &corev1.Service{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: serviceName, Namespace: instanceNamespace}, currentService)
	if err != nil && errors.IsNotFound(err) {
		// Create a new Service
		logger.Info("Creating a new Service", "Service.Namespace", newService.Namespace, "Service.Name", newService.Name)
		err = client.Create(context.TODO(), newService)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue
			logger.Info(" Service already exists")
			*needToRequeue = true
		} else if err != nil {
			logger.Error(err, "Failed to create new Service", "Service.Namespace", newService.Namespace, "Service.Name", newService.Name)
			return err
		} else {
			// Service created successfully - return and requeue
			*needToRequeue = true
		}
	} else if err != nil {
		logger.Error(err, "Failed to get Service", "Service.Name", serviceName)
		return err
	} else {
		// Found service, so determine if the resource has changed
		logger.Info("Comparing Services")
		if !IsServiceEqual(currentService, newService) {
			logger.Info("Updating Service", "Service.Name", currentService.Name)
			// Can't copy the entire Spec because ClusterIP is immutable
			currentService.ObjectMeta.Name = newService.ObjectMeta.Name
			currentService.ObjectMeta.Labels = newService.ObjectMeta.Labels
			currentService.Spec.Ports = newService.Spec.Ports
			currentService.Spec.Selector = newService.Spec.Selector
			err = client.Update(context.TODO(), currentService)
			if err != nil {
				logger.Error(err, "Failed to update Service",
					"Service.Namespace", currentService.Namespace, "Service.Name", currentService.Name)
				return err
			}
		}
	}
	return nil
}

// Check if the Ingress already exists, if not create a new one.
func ReconcileIngress(client client.Client, instanceNamespace string, ingressName string,
	newIngress *netv1.Ingress, needToRequeue *bool) error {
	logger := log.WithValues("func", "ReconcileIngress")

	currentIngress := &netv1.Ingress{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: ingressName, Namespace: instanceNamespace}, currentIngress)
	if err != nil && errors.IsNotFound(err) {
		// Create a new Ingress
		logger.Info("Creating a new Ingress", "Ingress.Namespace", newIngress.Namespace, "Ingress.Name", newIngress.Name)
		err = client.Create(context.TODO(), newIngress)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue
			logger.Info("Ingress already exists")
			*needToRequeue = true
		} else if err != nil {
			logger.Error(err, "Failed to create new Ingress", "Ingress.Namespace", newIngress.Namespace,
				"Ingress.Name", newIngress.Name)
			return err
		} else {
			// Ingress created successfully - return and requeue
			*needToRequeue = true
		}
	} else if err != nil {
		logger.Error(err, "Failed to get Ingress", "Ingress.Name", ingressName)
		return err
	} else {
		// Found Ingress, so determine if the resource has changed
		logger.Info("Comparing Ingresses")
		if !IsIngressEqual(currentIngress, newIngress) {
			logger.Info("Updating Ingress", "Ingress.Name", currentIngress.Name)
			currentIngress.ObjectMeta.Name = newIngress.ObjectMeta.Name
			currentIngress.ObjectMeta.Labels = newIngress.ObjectMeta.Labels
			currentIngress.Spec = newIngress.Spec
			err = client.Update(context.TODO(), currentIngress)
			if err != nil {
				logger.Error(err, "Failed to update Ingress",
					"Ingress.Namespace", currentIngress.Namespace, "Ingress.Name", currentIngress.Name)
				return err
			}
		}
	}
	return nil
}

// Use DeepEqual to determine if 2 daemon sets are equal.
// Check labels, pod template labels, service account names, volumes,
// containers, init containers, image name, volume mounts, env vars, liveness, readiness.
// If there are any differences, return false. Otherwise, return true.
func IsDaemonSetEqual(oldDaemonSet, newDaemonSet *appsv1.DaemonSet) bool {
	logger := log.WithValues("func", "IsDaemonSetEqual")

	if !reflect.DeepEqual(oldDaemonSet.ObjectMeta.Name, newDaemonSet.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldDaemonSet.ObjectMeta.Name, "new", newDaemonSet.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldDaemonSet.ObjectMeta.Labels, newDaemonSet.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldDaemonSet.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newDaemonSet.ObjectMeta.Labels))
		return false
	}

	logger.Info("DaemonSets are equal", "DaemonSet.Name", oldDaemonSet.ObjectMeta.Name)

	return true
}

// Use DeepEqual to determine if 2 services are equal.
// Check ObjectMeta, Ports and Selector.
// If there are any differences, return false. Otherwise, return true.
func IsServiceEqual(oldService, newService *corev1.Service) bool {
	logger := log.WithValues("func", "IsServiceEqual")

	if !reflect.DeepEqual(oldService.ObjectMeta.Name, newService.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldService.ObjectMeta.Name, "new", newService.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldService.ObjectMeta.Labels, newService.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldService.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newService.ObjectMeta.Labels))
		return false
	}

	// Can't check the entire Spec because ClusterIP is immutable
	if !reflect.DeepEqual(oldService.Spec.Ports, newService.Spec.Ports) {
		logger.Info("Ports not equal",
			"old", fmt.Sprintf("%v", oldService.Spec.Ports),
			"new", fmt.Sprintf("%v", newService.Spec.Ports))
		return false
	}

	if !reflect.DeepEqual(oldService.Spec.Selector, newService.Spec.Selector) {
		logger.Info("Selectors not equal",
			"old", fmt.Sprintf("%v", oldService.Spec.Selector),
			"new", fmt.Sprintf("%v", newService.Spec.Selector))
		return false
	}

	logger.Info("Services are equal", "Service.Name", oldService.ObjectMeta.Name)

	return true
}

// Use DeepEqual to determine if 2 ingresses are equal.
// Check ObjectMeta and Spec.
// If there are any differences, return false. Otherwise, return true.
func IsIngressEqual(oldIngress, newIngress *netv1.Ingress) bool {
	logger := log.WithValues("func", "IsIngressEqual")

	if !reflect.DeepEqual(oldIngress.ObjectMeta.Name, newIngress.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldIngress.ObjectMeta.Name, "new", newIngress.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldIngress.ObjectMeta.Labels, newIngress.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldIngress.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newIngress.ObjectMeta.Labels))
		return false
	}

	if !reflect.DeepEqual(oldIngress.Spec, newIngress.Spec) {
		logger.Info("Specs not equal",
			"old", fmt.Sprintf("%v", oldIngress.Spec),
			"new", fmt.Sprintf("%v", newIngress.Spec))
		return false
	}

	logger.Info("Ingresses are equal", "Ingress.Name", oldIngress.ObjectMeta.Name)

	return true
}
