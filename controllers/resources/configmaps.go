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
	"fmt"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func createConfigMap(ctx context.Context, client client.Client, cm *corev1.ConfigMap, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "createConfigMap", "instance.Name", instance.Name, "configmap.Name", cm.Name)

	err := controllerutil.SetControllerReference(instance, cm, client.Scheme())
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to set owner for configmap: %s", cm.Name))
		return err
	}

	reqLogger.Info(fmt.Sprintf("Creating configmap: %s", cm.Name))
	err = client.Create(ctx, cm)
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to create configmap: %s", cm.Name))
		return err
	}

	// Created log4js configmap successfully, requiring a requeue
	*needToRequeue = true

	return nil
}

func ReconcileLog4jsConfigMap(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileLog4jsConfigMap", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling log4js configmap")

	cm := &corev1.ConfigMap{}

	// Check if the log4js configmap already exists, if not create a new one
	err := client.Get(ctx, types.NamespacedName{Name: Log4jsConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			metaLabels := LabelsForMetadata(Log4jsConfigMapName)
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      Log4jsConfigMapName,
					Namespace: instance.Namespace,
					Labels:    metaLabels,
				},
				Data: Log4jsConfigMapData,
			}

			err = createConfigMap(ctx, client, cm, instance, needToRequeue)
			if err != nil {
				return err
			}
		} else {
			reqLogger.Error(err, "Failed to get log4js configmap")
			return err
		}
	}

	return nil
}

func ReconcileCommonUIConfigConfigMap(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileCommonUiConfigConfigMap", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling common-web-ui-config configmap")

	cm := &corev1.ConfigMap{}

	// Check if the log4js configmap already exists, if not create a new one
	err := client.Get(ctx, types.NamespacedName{Name: CommonConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      CommonConfigMapName,
					Namespace: instance.Namespace,
					Labels: map[string]string{"app.kubernetes.io/instance": "ibm-commonui-operator",
						"app.kubernetes.io/name": CommonConfigMapName, "app.kubernetes.io/managed-by": "ibm-commonui-operator"},
				},
			}

			err = createConfigMap(ctx, client, cm, instance, needToRequeue)
			if err != nil {
				return err
			}
		} else {
			reqLogger.Error(err, "Failed to get common-web-ui-config configmap")
			return err
		}
	}

	return nil
}

func CommonWebUIConfigMap(namespace string) *corev1.ConfigMap {
	reqLogger := log.WithValues("func", "CommonWebUIConfigMap")
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(CommonConfigMapName)
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CommonConfigMapName,
			Namespace: namespace,
			Labels:    metaLabels,
		},
	}

	return configmap
}

func DeleteConfigMap(ctx context.Context, client client.Client, name string, namespace string) error {
	reqLogger := log.WithValues("func", "deleteConfigmap", "name", name, "namespace", namespace)
	reqLogger.Info("Deleting configmap")

	//Get and delete common ui bind info config map
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	err := client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Configmap not found")
			return nil
		}
		reqLogger.Error(err, "Failed reading configmap")
		return err
	}

	// Delete configmap if found
	err = client.Delete(ctx, configMap)
	if err != nil {
		reqLogger.Error(err, "Failed to delete configmap")
		return err
	}

	reqLogger.Info("Deleted configmap")
	return nil

}
