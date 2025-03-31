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

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

const HPAName = "common-web-ui-hpa"

func getDesiredHorizontalPodAutoscaler(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	reqLogger := log.WithValues("func", "getDesiredHorizontalPodAutoscaler", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	//Determine min and max replicas
	//The formula is tshirt size of medium or > then use 2 for min and 2x+1 for max
	//If we have more than on replica, then assumed to be medium or greater for the purpose of calculating replicas
	var minReplicas int32 = 1
	var maxReplicas int32 = 3

	if instance.Spec.Replicas > 1 {
		minReplicas = 2
		maxReplicas = (instance.Spec.Replicas * 2) + 1
	}

	//Determine average utilization
	//The formula is when the limit is < 130% of request then use 90,
	//When the gap is greater than 130%, then use (limit * .7)/request * 100
	var averageUtilization int32 = 90
	request := int32(GetResourceMemoryWithDefault(instance.Spec.Resources.Requests.RequestMemory, 512))
	limit := int32(GetResourceMemoryWithDefault(instance.Spec.Resources.Limits.CPUMemory, 512))

	reqLogger.Info("computing average utilization", "request", request, "limit", limit, "base utilization (limit/request)*100", (float64(limit)/float64(request))*100)

	//When the gap between limit and request is > 130, bump averageUtilization
	if (float64(limit)/float64(request))*100 > 130 {
		averageUtilization = int32(float64(limit*70) / float64(request))
		reqLogger.Info("Setting large gap utilization", "averageUtilization", averageUtilization)
	}

	metaLabels := MergeMap(LabelsForMetadata(HPAName), instance.Spec.Labels)
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HPAName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       DeploymentName,
				APIVersion: "apps/v1",
			},
			MinReplicas: &minReplicas,
			MaxReplicas: maxReplicas,
			Metrics: []autoscalingv2.MetricSpec{
				{
					Type: autoscalingv2.ResourceMetricSourceType,
					Resource: &autoscalingv2.ResourceMetricSource{
						Name: "memory",
						Target: autoscalingv2.MetricTarget{
							Type:               "Utilization",
							AverageUtilization: &averageUtilization,
						},
					},
				},
			},
		},
	}

	err := controllerutil.SetControllerReference(instance, hpa, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for HPA")
		return nil, err
	}

	return hpa, nil
}

func ReconcileHorizontalPodAutoscaler(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileHorizontalPodAutoscaler", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling HPA")

	if !instance.Spec.AutoScaleConfig {
		//Horizontal pod autoscaling is disabled, delete the hpa if it exists
		reqLogger.Info("HPA disabled - delete HPA if it exists", "routeName", HPAName)

		hpa := &autoscalingv2.HorizontalPodAutoscaler{}
		err := client.Get(ctx, types.NamespacedName{Name: HPAName, Namespace: instance.Namespace}, hpa)
		if err != nil {
			if errors.IsNotFound(err) {
				reqLogger.Info("HPA not found - deletion is skipped", "Name", HPAName, "Namespace", instance.Namespace)
			} else {
				reqLogger.Error(err, "Unable to read the HPA for deletion - deletion skipped, but reconciliation will proceed")
			}
			return nil //Do not stop reconciliation if there was an error
		}
		err = client.Delete(ctx, hpa)
		if err != nil {
			reqLogger.Error(err, "Error deleting HPA - reconciliation will proceed")
		} else {
			reqLogger.Info("HPA deleted")
		}
		return nil //Do not stop reconciliation if there was an error
	}

	hpa := &autoscalingv2.HorizontalPodAutoscaler{}
	desiredHPA, desiredErr := getDesiredHorizontalPodAutoscaler(client, instance)
	if desiredErr != nil {
		return desiredErr
	}

	err := client.Get(ctx, types.NamespacedName{Name: HPAName, Namespace: instance.Namespace}, hpa)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new HPA", "HPA.Namespace", desiredHPA.Namespace, "HPA.Name", desiredHPA.Name)

		err = client.Create(ctx, desiredHPA)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// HPA already exists from a previous reconcile
				reqLogger.Info("HPA already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new HPA
				reqLogger.Info("Failed to create a new HPA", "HPA.Namespace", desiredHPA.Namespace, "HPA.Name", desiredHPA.Name)
				return err
			}
		} else {
			// Requeue after creating new HPA
			*needToRequeue = true
		}
	} else if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get HPA", "HPA.Namespace", instance.Namespace, "HPA.Name", HPAName)
		return err
	} else {
		// Determine if current HPA has changed
		reqLogger.Info("Comparing current and desired HPAs")

		if !IsHorizontalPodscalerEqual(hpa, desiredHPA) {
			reqLogger.Info("Updating HPA", "Namespace", hpa.Namespace, "Name", hpa.Name)
			hpa.ObjectMeta.Name = desiredHPA.ObjectMeta.Name
			hpa.ObjectMeta.Labels = desiredHPA.ObjectMeta.Labels
			hpa.Spec = desiredHPA.Spec

			err = client.Update(ctx, hpa)
			if err != nil {
				reqLogger.Error(err, "Failed to update HPA", "Namespace", hpa.Namespace, "Name", hpa.Name)
				return err
			}
		}
	}

	return nil
}
