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

	certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Check if a DaemonSet already exists. If not, create a new one.
func ReconcileDeployment(client client.Client, instanceNamespace string, deploymentName string,
	newDeployment *appsv1.Deployment, needToRequeue *bool) error {
	logger := log.WithValues("func", "ReconcileDeployment")

	currentDeployment := &appsv1.Deployment{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: deploymentName, Namespace: instanceNamespace}, currentDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Create a new deployment
		logger.Info("Creating a new Deployment", "Deployment.Namespace", newDeployment.Namespace, "Deployment.Name", newDeployment.Name)
		err = client.Create(context.TODO(), newDeployment)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue
			logger.Info("Deployment already exists")
			*needToRequeue = true
		} else if err != nil {
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", newDeployment.Namespace,
				"Deployment.Name", newDeployment.Name)
			return err
		} else {
			// Deployment created successfully - return and requeue
			*needToRequeue = true
		}
	} else if err != nil {
		logger.Error(err, "Failed to get Deployment", "Deployment.Name", deploymentName)
		return err
	} else {
		// Found deployment, so determine if the resource has changed
		logger.Info("Comparing Deployments")

		// Preserve cert-manager added labels in metadata
		if val, ok := currentDeployment.ObjectMeta.Labels[certRestartLabel]; ok {
			newDeployment.ObjectMeta.Labels[certRestartLabel] = val
		}

		// Preserve cert-manager added labels in spec
		if val, ok := currentDeployment.Spec.Template.ObjectMeta.Labels[certRestartLabel]; ok {
			newDeployment.Spec.Template.ObjectMeta.Labels[certRestartLabel] = val
		}

		if !IsDeploymentEqual(currentDeployment, newDeployment) {
			logger.Info("Updating Deployment", "Deployment.Name", currentDeployment.Name)
			currentDeployment.ObjectMeta.Name = newDeployment.ObjectMeta.Name
			currentDeployment.ObjectMeta.Labels = newDeployment.ObjectMeta.Labels
			currentReplicas := *currentDeployment.Spec.Replicas
			currentDeployment.Spec = newDeployment.Spec
			if currentReplicas == 0 {
				// since currentDeployment has been scaled to 0,
				// don't use the default replica count in newDeployment.
				currentDeployment.Spec.Replicas = &currentReplicas
			}
			err = client.Update(context.TODO(), currentDeployment)
			if err != nil {
				logger.Error(err, "Failed to update Deployment",
					"Deployment.Namespace", currentDeployment.Namespace, "Deployment.Name", currentDeployment.Name)
				return err
			}
		}
	}
	return nil
}

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
			currentIngress.ObjectMeta.Annotations = newIngress.ObjectMeta.Annotations
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

// Check if the Certificates already exist, if not create new ones.
func ReconcileCertificate(client client.Client, instanceNamespace, certificateName string,
	newCertificate *certmgr.Certificate, needToRequeue *bool) error {
	logger := log.WithValues("func", "ReconcileCertificate")

	currentCertificate := &certmgr.Certificate{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: certificateName, Namespace: instanceNamespace}, currentCertificate)
	if err != nil && errors.IsNotFound(err) {
		// Create a new Certificate
		logger.Info("Creating a new Certificate", "Certificate.Namespace", newCertificate.Namespace, "Certificate.Name", newCertificate.Name)
		err = client.Create(context.TODO(), newCertificate)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue
			logger.Info("Certificate already exists")
			*needToRequeue = true
		} else if err != nil {
			logger.Error(err, "Failed to create new Certificate", "Certificate.Namespace", newCertificate.Namespace,
				"Certificate.Name", newCertificate.Name)
			return err
		} else {
			// Certificate created successfully - return and requeue
			*needToRequeue = true
		}
	} else if err != nil {
		logger.Error(err, "Failed to get Certificate", "Certificate.Name", certificateName)
		return err
	} else {
		// Found Certificate, so determine if the resource has changed
		logger.Info("Comparing Certificates")
		if !IsCertificateEqual(currentCertificate, newCertificate) {
			logger.Info("Updating Certificate", "Certificate.Name", currentCertificate.Name)
			currentCertificate.ObjectMeta.Name = newCertificate.ObjectMeta.Name
			currentCertificate.ObjectMeta.Labels = newCertificate.ObjectMeta.Labels
			currentCertificate.Spec = newCertificate.Spec
			err = client.Update(context.TODO(), currentCertificate)
			if err != nil {
				logger.Error(err, "Failed to update Certificate", "Certificate.Namespace", currentCertificate.Namespace,
					"Certificate.Name", currentCertificate.Name)
				return err
			}
		}
	}
	return nil
}

// Use DeepEqual to determine if 2 deployments are equal.
// Check labels, replicas, pod template labels, service account names, volumes,
// containers, init containers, image name, volume mounts, env vars, liveness, readiness.
// If there are any differences, return false. Otherwise, return true.
// oldDeployment is the deployment that is currently running.
// newDeployment is what we expect the deployment to look like.
func IsDeploymentEqual(oldDeployment, newDeployment *appsv1.Deployment) bool {
	logger := log.WithValues("func", "IsDeploymentEqual")

	if !reflect.DeepEqual(oldDeployment.ObjectMeta.Name, newDeployment.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldDeployment.ObjectMeta.Name, "new", newDeployment.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldDeployment.ObjectMeta.Labels, newDeployment.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldDeployment.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newDeployment.ObjectMeta.Labels))
		return false
	}

	if *oldDeployment.Spec.Replicas != *newDeployment.Spec.Replicas {
		if *oldDeployment.Spec.Replicas == 0 {
			logger.Info("Allowing deployment to scale to 0", "name", oldDeployment.ObjectMeta.Name)
		} else {
			logger.Info("Replicas not equal", "old", oldDeployment.Spec.Replicas, "new", newDeployment.Spec.Replicas)
			return false
		}
	}

	oldPodTemplate := oldDeployment.Spec.Template
	newPodTemplate := newDeployment.Spec.Template
	if !isPodTemplateEqual(oldPodTemplate, newPodTemplate) {
		return false
	}

	logger.Info("Deployments are equal", "Deployment.Name", oldDeployment.ObjectMeta.Name)
	return true
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

	oldPodTemplate := oldDaemonSet.Spec.Template
	newPodTemplate := newDaemonSet.Spec.Template
	if !isPodTemplateEqual(oldPodTemplate, newPodTemplate) {
		return false
	}

	logger.Info("DaemonSets are equal", "DaemonSet.Name", oldDaemonSet.ObjectMeta.Name)

	return true
}

// Use DeepEqual to determine if 2 pod templates are equal.
// Check pod template labels, service account names, volumes,
// containers, init containers, image name, volume mounts, env vars, liveness, readiness.
// If there are any differences, return false. Otherwise, return true.
func isPodTemplateEqual(oldPodTemplate, newPodTemplate corev1.PodTemplateSpec) bool {
	logger := log.WithValues("func", "isPodTemplateEqual")

	if !reflect.DeepEqual(oldPodTemplate.ObjectMeta.Labels, newPodTemplate.ObjectMeta.Labels) {
		logger.Info("Pod labels not equal",
			"old", fmt.Sprintf("%v", oldPodTemplate.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newPodTemplate.ObjectMeta.Labels))
		return false
	}

	if !reflect.DeepEqual(oldPodTemplate.Spec.ServiceAccountName, newPodTemplate.Spec.ServiceAccountName) {
		logger.Info("Service account names not equal",
			"old", oldPodTemplate.Spec.ServiceAccountName,
			"new", newPodTemplate.Spec.ServiceAccountName)
		return false
	}

	oldVolumes := oldPodTemplate.Spec.Volumes
	newVolumes := newPodTemplate.Spec.Volumes
	if len(oldVolumes) == len(newVolumes) {
		if len(oldVolumes) > 0 {
			for i := range oldVolumes {
				oldVolume := oldVolumes[i]
				newVolume := newVolumes[i]
				if !reflect.DeepEqual(oldVolume.Name, newVolume.Name) {
					logger.Info("Pod volume names not equal", "volume num", i,
						"old", oldVolume.Name, "new", newVolume.Name)
					return false
				}
				if !reflect.DeepEqual(oldVolume.VolumeSource, newVolume.VolumeSource) {
					logger.Info("Pod volume sources not equal", "volume num", i,
						"old", fmt.Sprintf("%+v", oldVolume.VolumeSource), "new", fmt.Sprintf("%+v", newVolume.VolumeSource))
					return false
				}
			}
		}
	} else {
		logger.Info("Volume lengths not equal")
		return false
	}

	// check containers
	oldContainers := oldPodTemplate.Spec.Containers
	newContainers := newPodTemplate.Spec.Containers
	if !isContainerEqual(oldContainers, newContainers, false) {
		return false
	}

	// check init containers
	oldInitContainers := oldPodTemplate.Spec.InitContainers
	newInitContainers := newPodTemplate.Spec.InitContainers
	if !isContainerEqual(oldInitContainers, newInitContainers, true) {
		return false
	}

	logger.Info("Pod templates are equal")
	return true
}

// Use DeepEqual to determine if 2 container lists are equal.
// Check count, name, image name, image pull policy, env vars, volume mounts.
// If there are any differences, return false. Otherwise, return true.
// Set isInitContainer to true when checking init containers.
func isContainerEqual(oldContainers, newContainers []corev1.Container, isInitContainer bool) bool {
	logger := log.WithValues("func", "isContainerEqual")

	var containerType string
	if isInitContainer {
		containerType = "Init Container"
	} else {
		containerType = "Container"
	}

	if len(oldContainers) == len(newContainers) {
		if len(oldContainers) > 0 {
			for i := range oldContainers {
				oldContainer := oldContainers[i]
				newContainer := newContainers[i]
				logger.Info("Checking "+containerType, "old", oldContainer.Name)
				if !reflect.DeepEqual(oldContainer.Name, newContainer.Name) {
					logger.Info(containerType+" names not equal", "container num", i, "old", oldContainer.Name, "new", newContainer.Name)
					return false
				}

				if !reflect.DeepEqual(oldContainer.Image, newContainer.Image) {
					logger.Info(containerType+" images not equal", "container num", i, "old", oldContainer.Image, "new", newContainer.Image)
					return false
				}

				if !reflect.DeepEqual(oldContainer.ImagePullPolicy, newContainer.ImagePullPolicy) {
					logger.Info(containerType+" image pull policies not equal", "container num", i,
						"old", oldContainer.ImagePullPolicy, "new", newContainer.ImagePullPolicy)
					return false
				}

				if !reflect.DeepEqual(oldContainer.Resources.Limits["cpu"], newContainer.Resources.Limits["cpu"]) {
					logger.Info(containerType+" resource cpu limits not equal", "container num", i,
						"old", oldContainer.Resources.Limits["cpu"], "new", newContainer.Resources.Limits["cpu"])
					return false
				}

				if !reflect.DeepEqual(oldContainer.Resources.Limits["memory"], newContainer.Resources.Limits["memory"]) {
					logger.Info(containerType+" resource memory limits not equal", "container num", i,
						"old", oldContainer.Resources.Limits["memory"], "new", newContainer.Resources.Limits["memory"])
					return false
				}

				if !reflect.DeepEqual(oldContainer.Resources.Requests["cpu"], newContainer.Resources.Requests["cpu"]) {
					logger.Info(containerType+" resource cpu requests not equal", "container num", i,
						"old", oldContainer.Resources.Requests["cpu"], "new", newContainer.Resources.Requests["cpu"])
					return false
				}

				if !reflect.DeepEqual(oldContainer.Resources.Requests["memory"], newContainer.Resources.Requests["memory"]) {
					logger.Info(containerType+" resource memory requests not equal", "container num", i,
						"old", oldContainer.Resources.Requests["memory"], "new", newContainer.Resources.Requests["memory"])
					return false
				}

				oldEnvVars := oldContainer.Env
				newEnvVars := newContainer.Env
				if len(oldEnvVars) != len(newEnvVars) {
					logger.Info("Env var length not equal", "container num", i)
					return false
				} else if len(oldEnvVars) > 0 {
					for j := range oldEnvVars {
						oldEnvVar := oldEnvVars[j]
						newEnvVar := newEnvVars[j]
						if !reflect.DeepEqual(oldEnvVar.Name, newEnvVar.Name) {
							logger.Info("Env var names not equal", "container num", i, "old", oldEnvVar.Name, "new", newEnvVar.Name)
							return false
						}
						if !reflect.DeepEqual(oldEnvVar.Value, newEnvVar.Value) {
							logger.Info("Env var values not equal", "container num", i, "var", oldEnvVar.Name,
								"old", oldEnvVar.Value, "new", newEnvVar.Value)
							return false
						}
						if oldEnvVar.ValueFrom != nil && newEnvVar.ValueFrom != nil {
							if !reflect.DeepEqual(oldEnvVar.ValueFrom, newEnvVar.ValueFrom) {
								logger.Info("Env var ValueFrom not equal", "container num", i, "var", oldEnvVar.Name,
									"old", fmt.Sprintf("%+v", oldEnvVar.ValueFrom), "new", fmt.Sprintf("%+v", newEnvVar.ValueFrom))
								return false
							}
						} else if !(oldEnvVar.ValueFrom == nil && newEnvVar.ValueFrom == nil) {
							logger.Info("One of the env var's ValueFrom is nil", "container num", i, "var", oldEnvVar.Name)
							return false
						}
					}
				}

				oldVolumeMounts := oldContainer.VolumeMounts
				newVolumeMounts := newContainer.VolumeMounts
				if len(oldVolumeMounts) == len(newVolumeMounts) {
					if len(oldVolumeMounts) > 0 {
						for i := range oldVolumeMounts {
							oldVolumeMount := oldVolumeMounts[i]
							newVolumeMount := newVolumeMounts[i]
							if !reflect.DeepEqual(oldVolumeMount, newVolumeMount) {
								logger.Info("Volume mounts not equal", "mount num", i, "container num", i,
									"old", fmt.Sprintf("%+v", oldVolumeMount), "new", fmt.Sprintf("%+v", newVolumeMount))
								return false
							}
						}
					}
				} else {
					logger.Info("Volume mount lengths not equal", "container num", i)
					return false
				}

				if !isInitContainer {
					// check liveness and readiness probes
					oldLiveness := oldContainer.LivenessProbe
					newLiveness := newContainer.LivenessProbe
					if !isProbeEqual(oldLiveness, newLiveness, "Liveness") {
						return false
					}
					oldReadiness := oldContainer.ReadinessProbe
					newReadiness := newContainer.ReadinessProbe
					if !isProbeEqual(oldReadiness, newReadiness, "Readiness") {
						return false
					}
				}
			}
		}
	} else {
		logger.Info(containerType+" numbers not equal",
			"old", len(oldContainers), "new", len(newContainers))
		return false
	}
	return true
}

// Use DeepEqual to determine if 2 probes are equal.
// Check Handler, InitialDelaySeconds, TimeoutSeconds, PeriodSeconds.
// If there are any differences, return false. Otherwise, return true.
func isProbeEqual(oldProbe, newProbe *corev1.Probe, probeType string) bool {
	logger := log.WithValues("func", "isProbeEqual")
	logger.Info("Checking " + probeType + " probe")

	if oldProbe != nil && newProbe != nil {
		if !reflect.DeepEqual(oldProbe.Handler, newProbe.Handler) {
			logger.Info(probeType+" probe Handler not equal",
				"old", fmt.Sprintf("%+v", oldProbe.Handler), "new", fmt.Sprintf("%+v", newProbe.Handler))
			return false
		}

		if !reflect.DeepEqual(oldProbe.InitialDelaySeconds, newProbe.InitialDelaySeconds) {
			logger.Info(probeType+" probe Initial delay seconds not equal",
				"old", oldProbe.InitialDelaySeconds, "new", newProbe.InitialDelaySeconds)
			return false
		}

		if !reflect.DeepEqual(oldProbe.TimeoutSeconds, newProbe.TimeoutSeconds) {
			logger.Info(probeType+" probe Timeout seconds not equal",
				"old", oldProbe.TimeoutSeconds, "new", newProbe.TimeoutSeconds)
			return false
		}

		if !reflect.DeepEqual(oldProbe.PeriodSeconds, newProbe.PeriodSeconds) {
			logger.Info(probeType+" probe Period seconds not equal",
				"old", oldProbe.PeriodSeconds, "new", newProbe.PeriodSeconds)
			return false
		}
	} else if !(oldProbe == nil && newProbe == nil) {
		logger.Info("One "+probeType+" probe is nil",
			"old", fmt.Sprintf("%+v", oldProbe), "new", fmt.Sprintf("%+v", newProbe))
		return false
	}

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

	if !reflect.DeepEqual(oldIngress.ObjectMeta.Annotations, newIngress.ObjectMeta.Annotations) {
		logger.Info("Annotations not equal",
			"old", fmt.Sprintf("%v", oldIngress.ObjectMeta.Annotations),
			"new", fmt.Sprintf("%v", newIngress.ObjectMeta.Annotations))
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

// Use DeepEqual to determine if 2 certificates are equal.
// Check ObjectMeta and Spec.
// If there are any differences, return false. Otherwise, return true.
func IsCertificateEqual(oldCertificate, newCertificate *certmgr.Certificate) bool {
	logger := log.WithValues("func", "IsCertificateEqual")

	if !reflect.DeepEqual(oldCertificate.ObjectMeta.Name, newCertificate.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldCertificate.ObjectMeta.Name, "new", newCertificate.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldCertificate.ObjectMeta.Labels, newCertificate.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldCertificate.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newCertificate.ObjectMeta.Labels))
		return false
	}

	if !reflect.DeepEqual(oldCertificate.Spec, newCertificate.Spec) {
		logger.Info("Specs not equal",
			"old", fmt.Sprintf("%v", oldCertificate.Spec),
			"new", fmt.Sprintf("%v", newCertificate.Spec))
		return false
	}

	logger.Info("Certificates are equal", "Certificate.Name", oldCertificate.ObjectMeta.Name)

	return true
}
