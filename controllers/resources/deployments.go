//
// Copyright 2022 IBM Corporation
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
	"os"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

// nolint
func getDesiredDeployment(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, isZen bool, isCncf bool) (*appsv1.Deployment, error) {
	reqLogger := log.WithValues("func", "getDesiredDeployment", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	volumes := []corev1.Volume{}

	//Update any CR specified labels on deployment and container
	metaLabels := MergeMap(LabelsForMetadata(DeploymentName), instance.Spec.Labels)
	selectorLabels := LabelsForSelector(DeploymentName, CommonWebUICRType, instance.Name)
	podLabels := MergeMap(LabelsForPodMetadata(DeploymentName, CommonWebUICRType, instance.Name), instance.Spec.Labels)

	var err error

	replicas := instance.Spec.Replicas
	if replicas == 0 {
		replicas = 1
	}

	cpuLimits := GetResourceLimitsWithDefault(instance.Spec.Resources.Limits.CPULimits, 1000)
	cpuMemory := GetResourceMemoryWithDefault(instance.Spec.Resources.Limits.CPUMemory, 512)
	limEphemeral := GetResourceMemoryWithDefault(instance.Spec.Resources.Limits.EphemeralStorage, -1)
	reqLimits := GetResourceLimitsWithDefault(instance.Spec.Resources.Requests.RequestLimits, 300)
	reqMemory := GetResourceMemoryWithDefault(instance.Spec.Resources.Requests.RequestMemory, 512)
	reqEphemeral := GetResourceMemoryWithDefault(instance.Spec.Resources.Requests.EphemeralStorage, 251)

	imageRegistry := GetStringWithDefault(instance.Spec.CommonWebUIConfig.ImageRegistry, DefaultImageRegistry)
	imageTag := GetStringWithDefault(instance.Spec.CommonWebUIConfig.ImageTag, DefaultImageTag)
	image := GetImageID(imageRegistry, DefaultImageName, imageTag, "", "COMMON_WEB_UI_IMAGE")

	reqLogger.Info(fmt.Sprintf("Current image ID: %s", image))

	volumes = append(volumes, Log4jsVolume, ClusterCaVolume, UICertVolume, InternalTLSVolume, IAMDataVolume, IAMAuthDataVolume,
		WebUIConfigVolume, ClusterInfoConfigVolume, PlatformAuthIdpConfigVolume, ZenProductInfoConfigVolume)

	container := *CommonContainer.DeepCopy()
	container.Image = image
	container.Name = DeploymentName
	container.Env[6].Value = instance.Spec.GlobalUIConfig.CloudPakVersion
	container.Env[8].Value = instance.Spec.GlobalUIConfig.DefaultAuth
	container.Env[9].Value = instance.Spec.GlobalUIConfig.EnterpriseLDAP
	container.Env[10].Value = instance.Spec.GlobalUIConfig.EnterpriseSAML
	container.Env[11].Value = instance.Spec.GlobalUIConfig.OSAuth
	container.Env[19].Value = instance.Spec.CommonWebUIConfig.LandingPage
	container.Env[20].Value = os.Getenv("WATCH_NAMESPACE")
	container.Env[27].Value = strconv.FormatBool(instance.Spec.EnableInstanaMetricCollection)

	container.Resources.Limits["cpu"] = *resource.NewMilliQuantity(cpuLimits, resource.DecimalSI)
	container.Resources.Limits["memory"] = *resource.NewQuantity(cpuMemory*1024*1024, resource.BinarySI)
	container.Resources.Requests["cpu"] = *resource.NewMilliQuantity(reqLimits, resource.DecimalSI)
	container.Resources.Requests["memory"] = *resource.NewQuantity(reqMemory*1024*1024, resource.BinarySI)
	container.Resources.Requests["ephemeral-storage"] = *resource.NewQuantity(reqEphemeral*1024*1024, resource.BinarySI)
	//ephemeral-storage limit has no default, so only set it when it appears in the CR
	if limEphemeral > 0 {
		container.Resources.Limits["ephemeral-storage"] = *resource.NewQuantity(limEphemeral*1024*1024, resource.BinarySI)
	}
	container.VolumeMounts = CommonVolumeMounts

	if isZen {
		reqLogger.Info("Setting use zen to true in container def")
		container.Env[22].Value = "true"
	} else {
		reqLogger.Info("Setting use zen to false in container def")
		container.Env[22].Value = "false"
	}

	if isCncf {
		reqLogger.Info("Setting cluster type env var to cncf")
		container.Env[24].Value = "cncf"
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: DeploymentAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            "ibm-commonui-operand",
					HostNetwork:                   false,
					HostPID:                       false,
					HostIPC:                       false,
					TerminationGracePeriodSeconds: &Seconds60,
					TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
						{
							MaxSkew:           1,
							TopologyKey:       "topology.kubernetes.io/zone",
							WhenUnsatisfiable: corev1.ScheduleAnyway,
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"k8s-app": DeploymentName,
								},
							},
						},
						{
							MaxSkew:           1,
							TopologyKey:       "topology.kubernetes.io/region",
							WhenUnsatisfiable: corev1.ScheduleAnyway,
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"k8s-app": DeploymentName,
								},
							},
						},
					},
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "kubernetes.io/arch",
												Operator: corev1.NodeSelectorOpIn,
												Values:   ArchitectureList,
											},
										},
									},
								},
							},
						},
						PodAntiAffinity: &corev1.PodAntiAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: corev1.PodAffinityTerm{
										LabelSelector: &metav1.LabelSelector{
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "app.kubernetes.io/name",
													Operator: metav1.LabelSelectorOpIn,
													Values:   []string{DeploymentName},
												},
											},
										},
										TopologyKey: "topology.kubernetes.io/zone",
									},
								},
								{
									Weight: 100,
									PodAffinityTerm: corev1.PodAffinityTerm{
										LabelSelector: &metav1.LabelSelector{
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "app.kubernetes.io/name",
													Operator: metav1.LabelSelectorOpIn,
													Values:   []string{DeploymentName},
												},
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "dedicated",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
						{
							Key:      "CriticalAddonsOnly",
							Operator: corev1.TolerationOpExists,
						},
					},
					Volumes: volumes,
					Containers: []corev1.Container{
						container,
					},
				},
			},
		},
	}

	//If an image pull secret is set in the ENV, then set it into the pod spec
	ips, ipsExists := os.LookupEnv("IMAGE_PULL_SECRET")
	if ipsExists {
		reqLogger.Info(fmt.Sprintf("Setting image pull secret: %s", ips))
		deployment.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: ips},
		}
	}

	err = controllerutil.SetControllerReference(instance, deployment, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for deployment")
		return nil, err
	}

	return deployment, nil
}

// nolint
func ReconcileDeployment(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, isZen bool, isCncf bool, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileDeployment", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling deployment")

	deployment := &appsv1.Deployment{}

	desiredDeployment, desiredErr := getDesiredDeployment(ctx, client, instance, isZen, isCncf)
	if desiredErr != nil {
		return desiredErr
	}

	err := client.Get(ctx, types.NamespacedName{Name: DeploymentName, Namespace: instance.Namespace}, deployment)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new deployment", "Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)

		err = client.Create(ctx, desiredDeployment)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// Deployment already exists from a previous reconcile
				reqLogger.Info("Deployment already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new deployment
				reqLogger.Info("Failed to create a new deployment", "Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)
				return err
			}
		} else {
			// Requeue after creating new deployment
			*needToRequeue = true
		}
	} else if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get deployment", "Deployment.Namespace", instance.Namespace, "Deployment.Name", DeploymentName)
		return err
	} else {
		// Determine if current deployment has changed
		reqLogger.Info("Comparing current and desired deployments")

		// Preserve annotations added by NamespaceScope Operator
		PreserveKeyValue(NSSAnnotation, deployment.Spec.Template.ObjectMeta.Annotations, desiredDeployment.Spec.Template.ObjectMeta.Annotations)

		// Preserve labels added by cert-manager
		PreserveKeyValue(CertRestartLabel, deployment.ObjectMeta.Labels, desiredDeployment.ObjectMeta.Labels)
		PreserveKeyValue(CertRestartLabel, deployment.Spec.Template.ObjectMeta.Labels, desiredDeployment.Spec.Template.ObjectMeta.Labels)

		if !IsDeploymentEqual(deployment, desiredDeployment) {
			reqLogger.Info("Updating deployment", "Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)

			deployment.ObjectMeta.Name = desiredDeployment.ObjectMeta.Name
			deployment.ObjectMeta.Labels = desiredDeployment.ObjectMeta.Labels
			currentReplicas := *deployment.Spec.Replicas
			deployment.Spec = desiredDeployment.Spec

			if currentReplicas == 0 {
				// Since current deployment has been scaled to 0 replicas, do not use the default replica count in the new deployment.
				deployment.Spec.Replicas = &currentReplicas
			}

			err = client.Update(ctx, deployment)
			if err != nil {
				reqLogger.Error(err, "Failed to update deployment", "Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)
				return err
			}
		}
	}

	return nil
}
