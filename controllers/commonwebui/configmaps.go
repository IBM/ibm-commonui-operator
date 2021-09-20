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

	operatorv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"

	routesv1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *CommonWebUIReconciler) createConfigMap(ctx context.Context, cm *corev1.ConfigMap, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "createConfigMap", "instance.Name", instance.Name, "configmap.Name", cm.Name)
	
	err := controllerutil.SetControllerReference(instance, cm, r.Scheme)
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to set owner for configmap: %s", cm.Name))
		return err
	}

	reqLogger.Info(fmt.Sprintf("Creating configmap: %s", cm.Name))
	err = r.Client.Create(ctx, cm)
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to create configmap: %s", cm.Name))
		return err
	}

	// Created log4js configmap successfully, requiring a requeue
	*needToRequeue = true

	return nil
}

func (r *CommonWebUIReconciler) reconcileLog4jsConfigMap(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileLog4jsConfigMap", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling log4js configmap")

	cm := &corev1.ConfigMap{}

	// Check if the log4js configmap already exists, if not create a new one
	err := r.Client.Get(ctx, types.NamespacedName{Name: Log4jsConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			metaLabels := LabelsForMetadata(Log4jsConfigMapName)
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: Log4jsConfigMapName,
					Namespace: instance.Namespace,
					Labels: metaLabels,
				},
				Data: Log4jsConfigMapData,
			}
			
			err = r.createConfigMap(ctx, cm, instance, needToRequeue)
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

func (r *CommonWebUIReconciler) reconcileRedisConfigMap(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileRedisConfigMap", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling redis certs configmap")

	cm := &corev1.ConfigMap{}

	// Check if the redis certs configmap already exists, if not create a new one
	err := r.Client.Get(ctx, types.NamespacedName{Name: RedisConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			metaLabels := LabelsForMetadata(RedisConfigMapName)
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: RedisConfigMapName,
					Annotations: RedisConfigMapAnnotations,
					Namespace: instance.Namespace,
					Labels: metaLabels,
				},
			}

			err = r.createConfigMap(ctx, cm, instance, needToRequeue)
			if err != nil {
				return err
			}
		} else {
			reqLogger.Error(err, "Failed to get redis certs configmap")
			return err
		}
	}
	
	return nil
}

func (r *CommonWebUIReconciler) reconcileZenCardsConfigMap(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileZenCardsConfigMap", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling zen card extensions configmap")

	cm := &corev1.ConfigMap{}

	// Check if the zen card extensions configmap already exists, if not create a new one
	err := r.Client.Get(ctx, types.NamespacedName{Name: ZenCardsConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			metaLabels := LabelsForMetadata(ZenCardsConfigMapName)
			metaLabels["icpdata_addon"] = "true"
			metaLabels["icpdata_addon_version"] = "v" + instance.Spec.Version

			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: ZenCardsConfigMapName,
					Namespace: instance.Namespace,
					Labels: metaLabels,
				},
				Data: map[string]string{
					"nginx.conf": ZenNginxConfig,
					"extensions": ZenCardExtensions,
				},
			}

			err = r.createConfigMap(ctx, cm, instance, needToRequeue)
			if err != nil {
				return err
			}
		} else {
			reqLogger.Error(err, "Failed to get zen card extensions configmap")
			return err
		}
	}
	
	return nil
}

func (r *CommonWebUIReconciler) reconcileExtensionsConfigMap(ctx context.Context, instance *operatorv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileExtensionsConfigMap", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling extensions configmap")

	cm := &corev1.ConfigMap{}

	// Check if the extensions configmap already exists, if not create a new one
	err := r.Client.Get(ctx, types.NamespacedName{Name: ExtensionsConfigMapName, Namespace: instance.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			route := &routesv1.Route{}

			// Get the cp-console route and add it to the configmap
			err = r.Client.Get(ctx, types.NamespacedName{Name: ConsoleRouteName, Namespace: instance.Namespace}, route)
			if err != nil {
				reqLogger.Info("Failed to get route for cp-console, try again later")
				*needToRequeue = true
				return err
			}

			reqLogger.Info(fmt.Sprintf("Current cp-console route: %s", route.Spec.Host))

			metaLabels := LabelsForMetadata(ExtensionsConfigMapName)
			metaLabels["icpdata_addon"] = "true"

			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: ExtensionsConfigMapName,
					Namespace: instance.Namespace,
					Labels: metaLabels,
				},
				Data: map[string]string{
					"extensions": fmt.Sprintf(ExtensionsConfigMapData, route.Spec.Host),
				},
			}

			err = r.createConfigMap(ctx, cm, instance, needToRequeue)
			if err != nil {
				return err
			}
		} else {
			reqLogger.Error(err, "Failed to get extensions configmap")
			return err
		}
	}
	
	return nil
}
