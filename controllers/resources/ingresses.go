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

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
)

// type DesiredStateGetter func(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) (*netv1.Ingress, error)

func ReconcileRemoveIngresses(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) {
	reqLogger := log.WithValues("func", "ReconcileRemoveIngresses")

	//No error checking as we will just make a best attempt to remove the legacy ingresses
	//Do not fail based on inability to delete the ingresses
	ingresses := []string{APIIngressName, CallbackIngressName, NavIngressName}
	for _, iname := range ingresses {
		err := DeleteIngress(ctx, client, iname, instance.Namespace, needToRequeue)
		if err != nil {
			reqLogger.Info("Failed to delete legacy ingress " + iname)
		}
	}
}

func DeleteIngress(ctx context.Context, client client.Client, ingressName string, ingressNS string, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "deleteIngress", "Name", ingressName, "Namespace", ingressNS)

	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: ingressNS,
		},
	}

	err := client.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, ingress)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		reqLogger.Error(err, "Failed to get legacy ingress")
		return err
	}

	// Delete ingress if found
	err = client.Delete(ctx, ingress)
	if err != nil {
		reqLogger.Error(err, "Failed to delete legacy ingress")
		return err
	}

	reqLogger.Info("Deleted legacy ingress")
	*needToRequeue = true
	return nil
}

func getDesiredAPIIngress(client client.Client, instance *operatorsv1alpha1.CommonWebUI, isCncf bool) (*netv1.Ingress, error) {
	reqLogger := log.WithValues("func", "getDesiredAPIIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	metaLabels := LabelsForMetadata(APIIngressName)

	ingressPath := instance.Spec.CommonWebUIConfig.IngressPath
	apiIngressPath := ingressPath + "/api/"
	logoutIngressPath := ingressPath + "/logout/"

	pathType := netv1.PathType("ImplementationSpecific")
	var ingress *netv1.Ingress

	if isCncf {
		ingress = &netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        APIIngressName,
				Annotations: APIIngressAnnotations,
				Labels:      metaLabels,
				Namespace:   instance.Namespace,
			},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     apiIngressPath,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: instance.Spec.CommonWebUIConfig.ServiceName,
												Port: netv1.ServiceBackendPort{
													Number: 3000,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	} else {
		ingress = &netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        APIIngressName,
				Annotations: APIIngressAnnotations,
				Labels:      metaLabels,
				Namespace:   instance.Namespace,
			},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     apiIngressPath,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: instance.Spec.CommonWebUIConfig.ServiceName,
												Port: netv1.ServiceBackendPort{
													Number: 3000,
												},
											},
										},
									},
									{
										Path:     logoutIngressPath,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: instance.Spec.CommonWebUIConfig.ServiceName,
												Port: netv1.ServiceBackendPort{
													Number: 3000,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	err := controllerutil.SetControllerReference(instance, ingress, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for API ingress")
		return nil, err
	}

	return ingress, nil
}

func ReconcileAPIIngress(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, isCncf bool, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileAPIIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling API ingress")

	desiredIngress, err := getDesiredAPIIngress(client, instance, isCncf)
	if err != nil {
		return err
	}

	return reconcileIngress(ctx, client, instance, APIIngressName, desiredIngress, needToRequeue)
}

//nolint
func getDesiredCallbackIngress(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*netv1.Ingress, error) {
	reqLogger := log.WithValues("func", "getDesiredCallbackIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	metaLabels := LabelsForMetadata(CallbackIngressName)
	pathType := netv1.PathType("ImplementationSpecific")

	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        CallbackIngressName,
			Annotations: CallbackIngressAnnotations,
			Labels:      metaLabels,
			Namespace:   instance.Namespace,
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/auth/liberty/callback",
									PathType: &pathType,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: instance.Spec.CommonWebUIConfig.ServiceName,
											Port: netv1.ServiceBackendPort{
												Number: 3000,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := controllerutil.SetControllerReference(instance, ingress, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for callback ingress")
		return nil, err
	}

	return ingress, nil
}

func ReconcileCallbackIngress(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileCallbackIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling callback ingress")

	desiredIngress, err := getDesiredCallbackIngress(client, instance)
	if err != nil {
		return err
	}

	return reconcileIngress(ctx, client, instance, CallbackIngressName, desiredIngress, needToRequeue)
}

func getDesiredNavIngress(client client.Client, instance *operatorsv1alpha1.CommonWebUI) (*netv1.Ingress, error) {
	reqLogger := log.WithValues("func", "getDesiredNavIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	metaLabels := LabelsForMetadata(NavIngressName)
	pathType := netv1.PathType("ImplementationSpecific")

	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        NavIngressName,
			Annotations: NavIngressAnnotations,
			Labels:      metaLabels,
			Namespace:   instance.Namespace,
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     instance.Spec.CommonWebUIConfig.IngressPath,
									PathType: &pathType,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: instance.Spec.CommonWebUIConfig.ServiceName,
											Port: netv1.ServiceBackendPort{
												Number: 3000,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := controllerutil.SetControllerReference(instance, ingress, client.Scheme())
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for common-nav ingress")
		return nil, err
	}

	return ingress, nil
}

func ReconcileNavIngress(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileNavIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)
	reqLogger.Info("Reconciling common-nav ingress")

	desiredIngress, err := getDesiredNavIngress(client, instance)
	if err != nil {
		return err
	}

	return reconcileIngress(ctx, client, instance, NavIngressName, desiredIngress, needToRequeue)
}

//nolint
func reconcileIngress(ctx context.Context, client client.Client, instance *operatorsv1alpha1.CommonWebUI, ingressName string, desiredIngress *netv1.Ingress, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileIngress", "instance.Name", instance.Name, "instance.Namespace", instance.Namespace)

	ingress := &netv1.Ingress{}

	err := client.Get(ctx, types.NamespacedName{Name: ingressName, Namespace: instance.Namespace}, ingress)
	if err != nil && !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get ingress", "Ingress.Namespace", instance.Namespace, "Ingress.Name", ingressName)
		return err
	}

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ingress", "Ingress.Namespace", desiredIngress.Namespace, "Ingress.Name", desiredIngress.Name)

		err = client.Create(ctx, desiredIngress)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				// Ingress already exists from a previous reconcile
				reqLogger.Info("Ingress already exists")
				*needToRequeue = true
			} else {
				// Failed to create a new ingress
				reqLogger.Info("Failed to create a new ingress", "Ingress.Namespace", desiredIngress.Namespace, "Ingress.Name", desiredIngress.Name)
				return err
			}
		} else {
			// Requeue after creating new ingress
			*needToRequeue = true
		}
	} else {
		// Determine if current ingress has changed
		reqLogger.Info("Comparing current and desired ingresses")

		if !IsIngressEqual(ingress, desiredIngress) {
			reqLogger.Info("Updating ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)

			ingress.ObjectMeta.Name = desiredIngress.ObjectMeta.Name
			ingress.ObjectMeta.Labels = desiredIngress.ObjectMeta.Labels
			ingress.ObjectMeta.Annotations = desiredIngress.ObjectMeta.Annotations
			ingress.Spec = desiredIngress.Spec

			err = client.Update(ctx, ingress)
			if err != nil {
				reqLogger.Error(err, "Failed to update ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
				return err
			}
		}
	}

	return nil
}
