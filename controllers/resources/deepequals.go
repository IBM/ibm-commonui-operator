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
	"fmt"
	"reflect"

	certmgr "github.com/ibm/ibm-cert-manager-operator/apis/cert-manager/v1"
	route "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

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

	if !reflect.DeepEqual(oldPodTemplate.Spec.ImagePullSecrets, newPodTemplate.Spec.ImagePullSecrets) {
		logger.Info("Image pull secrets are not equal",
			"old", oldPodTemplate.Spec.ImagePullSecrets,
			"new", newPodTemplate.Spec.ImagePullSecrets)
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

func isContainerEnvEqual(oldContainer, newContainer corev1.Container, i int) bool {
	logger := log.WithValues("func", "isContainerEnvEqual")

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

				if !reflect.DeepEqual(oldContainer.SecurityContext, newContainer.SecurityContext) {
					logger.Info(containerType+" security context is not equal",
						"old", oldContainer.SecurityContext, "new", newContainer.SecurityContext)
					return false
				}

				//For some reason deep equals would not work on CPU values so went to a straight value comparison
				//for all types
				if !oldContainer.Resources.Limits.Cpu().Equal(*newContainer.Resources.Limits.Cpu()) {
					logger.Info(containerType+" resource cpu limits not equal", "container num", i,
						"old", oldContainer.Resources.Limits.Cpu(), "new", newContainer.Resources.Limits.Cpu())
					return false
				}

				if !reflect.DeepEqual(oldContainer.Resources.Limits["memory"], newContainer.Resources.Limits["memory"]) {
					logger.Info(containerType+" resource memory limits not equal", "container num", i,
						"old", oldContainer.Resources.Limits["memory"], "new", newContainer.Resources.Limits["memory"])
					return false
				}

				if !oldContainer.Resources.Limits.StorageEphemeral().Equal(*newContainer.Resources.Limits.StorageEphemeral()) {
					logger.Info(containerType+" resource ephemoral-storage limits not equal", "container num", i,
						"old", oldContainer.Resources.Limits["ephemeral-storage"], "new", newContainer.Resources.Limits["ephemeral-storage"])
					return false
				}

				//For some reason deep equals would not work on CPU values so went to a straight value comparison
				if !oldContainer.Resources.Requests.Cpu().Equal(*newContainer.Resources.Requests.Cpu()) {
					logger.Info(containerType+" resource cpu requests not equal", "container num", i,
						"old", oldContainer.Resources.Requests["cpu"], "new", newContainer.Resources.Requests["cpu"])
					return false
				}

				if !reflect.DeepEqual(oldContainer.Resources.Requests["memory"], newContainer.Resources.Requests["memory"]) {
					logger.Info(containerType+" resource memory requests not equal", "container num", i,
						"old", oldContainer.Resources.Requests["memory"], "new", newContainer.Resources.Requests["memory"])
					return false
				}

				if !oldContainer.Resources.Requests.StorageEphemeral().Equal(*newContainer.Resources.Requests.StorageEphemeral()) {
					logger.Info(containerType+" resource ephemoral-storage requests not equal", "container num", i,
						"old", oldContainer.Resources.Requests["ephemeral-storage"], "new", newContainer.Resources.Requests["ephemeral-storage"])
					return false
				}

				if !isContainerEnvEqual(oldContainer, newContainer, i) {
					return false
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
		if !reflect.DeepEqual(oldProbe.ProbeHandler, newProbe.ProbeHandler) {
			logger.Info(probeType+" probe Handler not equal",
				"old", fmt.Sprintf("%+v", oldProbe.ProbeHandler), "new", fmt.Sprintf("%+v", newProbe.ProbeHandler))
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

// Use DeepEqual to determine if 2 routes are equal.
// Check annotations and Spec.
// If there are any differences, return false. Otherwise, return true.
func IsRouteEqual(oldRoute, newRoute *route.Route) bool {
	logger := log.WithValues("func", "IsRouteEqual")

	if !reflect.DeepEqual(oldRoute.ObjectMeta.Name, newRoute.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldRoute.ObjectMeta.Name, "new", newRoute.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldRoute.ObjectMeta.Annotations, newRoute.ObjectMeta.Annotations) {
		logger.Info("Annotations not equal",
			"old", fmt.Sprintf("%v", oldRoute.ObjectMeta.Annotations),
			"new", fmt.Sprintf("%v", newRoute.ObjectMeta.Annotations))
		return false
	}

	if !reflect.DeepEqual(oldRoute.ObjectMeta.Labels, newRoute.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldRoute.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newRoute.ObjectMeta.Labels))
		return false
	}

	if !reflect.DeepEqual(oldRoute.Spec, newRoute.Spec) {
		//ugly, but don't print the CA to the log
		logger.Info("Specs not equal", "oldHost", oldRoute.Spec.Host, "newHost", newRoute.Spec.Host,
			"oldPath", oldRoute.Spec.Path, "newHost", newRoute.Spec.Path,
			"oldWildcardPolicy", oldRoute.Spec.WildcardPolicy, "newWildcardPolicy", newRoute.Spec.WildcardPolicy,
			"oldPort", fmt.Sprintf("%v", oldRoute.Spec.Port), "newPort", fmt.Sprintf("%v", newRoute.Spec.Port),
			"oldToService", fmt.Sprintf("%v", oldRoute.Spec.To), "newToService", fmt.Sprintf("%v", newRoute.Spec.To),
			"old.tls.termination", oldRoute.Spec.TLS.Termination, "new.tls.termination", newRoute.Spec.TLS.Termination,
			"old.tls.insecureEdgeTerminationPolicy", oldRoute.Spec.TLS.InsecureEdgeTerminationPolicy,
			"new.tls.insecureEdgeTerminationPolicy", newRoute.Spec.TLS.InsecureEdgeTerminationPolicy)
		return false
	}

	logger.Info("Routes are equal")

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

	//compare the specs, but allow renewBefore to be modified
	if !reflect.DeepEqual(oldCertificate.Spec.CommonName, newCertificate.Spec.CommonName) ||
		!reflect.DeepEqual(oldCertificate.Spec.DNSNames, newCertificate.Spec.DNSNames) ||
		!reflect.DeepEqual(oldCertificate.Spec.Duration, newCertificate.Spec.Duration) ||
		!reflect.DeepEqual(oldCertificate.Spec.IssuerRef, newCertificate.Spec.IssuerRef) ||
		!reflect.DeepEqual(oldCertificate.Spec.SecretName, newCertificate.Spec.SecretName) {
		// if !reflect.DeepEqual(oldCertificate.Spec, newCertificate.Spec) {
		logger.Info("Specs not equal (renewBefore is not checked)",
			"old", fmt.Sprintf("%v", oldCertificate.Spec),
			"new", fmt.Sprintf("%v", newCertificate.Spec))
		return false
	}

	logger.Info("Certificates are equal", "Certificate.Name", oldCertificate.ObjectMeta.Name)

	return true
}

// Use DeepEqual to determine if 2 service accounts are equal.
// Check metadata.
// If there are any differences, return false. Otherwise, return true.
func IsServiceAccountEqual(oldSA, newSA *corev1.ServiceAccount) bool {
	logger := log.WithValues("func", "IsServiceAccountEqual")

	if !reflect.DeepEqual(oldSA.ObjectMeta.Name, newSA.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldSA.ObjectMeta.Name, "new", newSA.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldSA.ObjectMeta.Labels, newSA.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldSA.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newSA.ObjectMeta.Labels))
		return false
	}

	logger.Info("Service accounts are equal")

	return true
}

// Use DeepEqual to determine if 2 roles are equal.
// Check metadata, labels, and rules.
// If there are any differences, return false. Otherwise, return true.
func IsRoleEqual(oldRole, newRole *rbacv1.Role) bool {
	logger := log.WithValues("func", "IsRoleEqual")

	if !reflect.DeepEqual(oldRole.ObjectMeta.Name, newRole.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldRole.ObjectMeta.Name, "new", newRole.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldRole.ObjectMeta.Labels, newRole.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldRole.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newRole.ObjectMeta.Labels))
		return false
	}

	if !reflect.DeepEqual(oldRole.Rules, newRole.Rules) {
		logger.Info("Rules not equal",
			"old", fmt.Sprintf("%v", oldRole.Rules),
			"new", fmt.Sprintf("%v", newRole.Rules))
		return false
	}

	logger.Info("Roles are equal")

	return true
}

// Use DeepEqual to determine if 2 role bindings are equal.
// Check metadata, labels, subjects, and role ref.
// If there are any differences, return false. Otherwise, return true.
func IsRoleBindingEqual(oldRoleBinding, newRoleBinding *rbacv1.RoleBinding) bool {
	logger := log.WithValues("func", "IsRoleBindingEqual")

	if !reflect.DeepEqual(oldRoleBinding.ObjectMeta.Name, newRoleBinding.ObjectMeta.Name) {
		logger.Info("Names not equal", "old", oldRoleBinding.ObjectMeta.Name, "new", newRoleBinding.ObjectMeta.Name)
		return false
	}

	if !reflect.DeepEqual(oldRoleBinding.ObjectMeta.Labels, newRoleBinding.ObjectMeta.Labels) {
		logger.Info("Labels not equal",
			"old", fmt.Sprintf("%v", oldRoleBinding.ObjectMeta.Labels),
			"new", fmt.Sprintf("%v", newRoleBinding.ObjectMeta.Labels))
		return false
	}

	if !reflect.DeepEqual(oldRoleBinding.Subjects, newRoleBinding.Subjects) {
		logger.Info("Rules not equal",
			"old", fmt.Sprintf("%v", oldRoleBinding.Subjects),
			"new", fmt.Sprintf("%v", newRoleBinding.Subjects))
		return false
	}

	if !reflect.DeepEqual(oldRoleBinding.RoleRef, newRoleBinding.RoleRef) {
		logger.Info("Rules not equal",
			"old", fmt.Sprintf("%v", oldRoleBinding.Subjects),
			"new", fmt.Sprintf("%v", newRoleBinding.Subjects))
		return false
	}

	logger.Info("Role bindings are equal")

	return true
}
