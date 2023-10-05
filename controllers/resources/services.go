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

package resources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

func getDesiredService(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*corev1.Service, error) {
	reqLogger := log.WithValues("func", "getDesiredService", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	metaLabels := LabelsForMetadata(ServiceName)
	metaLabels["kubernetes.io/cluster-service"] = "true"
	metaLabels["kubernetes.io/name"] = instance.Spec.CommonWebUIConfig.ServiceName
	metaLabels["app"] = instance.Spec.CommonWebUIConfig.ServiceName

	//Update any CR specified labels on the route
	if instance.Spec.Labels != nil {
		metaLabels = MergeMap(metaLabels, instance.Spec.Labels)
	}

	selectorLabels := LabelsForSelector(ServiceName, CommonWebUICRType, instance.Name)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.CommonWebUIConfig.ServiceName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     instance.Spec.CommonWebUIConfig.ServiceName,
					Port:     3000,
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 3000,
					},
				},
			},
			Selector: selectorLabels,
		},
	}

	err := controllerutil.SetControllerReference(instance, service, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for service")
		return nil, err
	}

	return service, nil
}

func ReconcileService(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileService", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling service")

	service := &corev1.Service{}

	desiredService, desiredErr := getDesiredService(client, instance)
	if desiredErr != nil {
		return desiredErr
	}

	err := client.Get(ctx, types.NamespacedName{Name: ServiceName, Namespace: instance.Namespace}, service)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new service", "Service.Namespace", desiredService.Namespace, "Service.Name", desiredService.Name)

		err = client.Create(ctx, desiredService)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// Service already exists from a previous reconcile
				reqLogger.Info("Service already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new service
				reqLogger.Info("Failed to create a new service", "Service.Namespace", desiredService.Namespace, "Service.Name", desiredService.Name)
				return err
			}
		} else {
			// Requeue after creating new service
			*needToRequeue = true
		}
	} else if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get service", "Service.Namespace", instance.Namespace, "Service.Name", ServiceName)
		return err
	} else {
		// Determine if current service has changed
		reqLogger.Info("Comparing current and desired services")

		if !IsServiceEqual(service, desiredService) {
			reqLogger.Info("Updating service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)

			service.ObjectMeta.Name = desiredService.ObjectMeta.Name
			service.ObjectMeta.Labels = desiredService.ObjectMeta.Labels
			service.Spec.Ports = desiredService.Spec.Ports
			service.Spec.Selector = desiredService.Spec.Selector

			err = client.Update(ctx, service)
			if err != nil {
				reqLogger.Error(err, "Failed to update service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
				return err
			}
		}
	}

	return nil
}
