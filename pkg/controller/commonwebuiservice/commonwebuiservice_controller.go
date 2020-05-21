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
package commonwebuiservice

import (
	"context"
	"encoding/json"
	gorun "runtime"

	res "github.com/ibm/ibm-commonui-operator/pkg/resources"

	"k8s.io/apimachinery/pkg/util/intstr"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"
	certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	apiextv1beta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const commonwebuiserviceCrType = "commonwebuiservice_cr"

var commonVolume = []corev1.Volume{}

var log = logf.Log.WithName("controller_commonwebuiservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CommonWebUIService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCommonWebUI{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	reqLogger := log.WithValues("func", "add")

	// Create a new controller
	c, err := controller.New("commonwebui-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CommonWebUIService
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.CommonWebUI{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource ConfigMap and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CommonWebUI{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource "Daemonset" and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CommonWebUI{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Service" and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CommonWebUI{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Ingress" and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &netv1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CommonWebUI{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Certificate" and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &certmgr.Certificate{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CommonWebUI{},
	})
	if err != nil {
		// Log error instead of failing because "cert-manager" might not be installed
		reqLogger.Error(err, "Failed to watch Certificate")
	}

	return nil
}

// blank assignment to verify that ReconcileCommonWebUI implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCommonWebUI{}

// ReconcileCommonWebUI reconciles a CommonWebUIService object
type ReconcileCommonWebUI struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CommonWebUIService object and makes changes based on the state read
// and what is in the CommonWebUIService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a DaemonSet
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCommonWebUI) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUI")

	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	// Fetch the CommonWebUIService CR instance
	instance := &operatorsv1alpha1.CommonWebUI{}

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	opVersion := instance.Spec.OperatorVersion
	reqLogger.Info("got CommonWebUIService instance, version=" + opVersion)

	// set a default Status value
	if len(instance.Status.Nodes) == 0 {
		instance.Status.Nodes = res.DefaultStatusForCR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set CommonWebUI default status")
			return reconcile.Result{}, err
		}
	}
	// Check if the config maps already exist. If not, create a new one.
	err = r.reconcileConfigMaps(instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the DaemonSet already exists, if not create a new one
	newDaemonSet, err := r.newDaemonSetForCR(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = res.ReconcileDaemonSet(r.client, instance.Namespace, res.DaemonSetName, newDaemonSet, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the common web ui Service already exist. If not, create a new one.
	newService, err := r.serviceForUI(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = res.ReconcileService(r.client, instance.Namespace, res.ServiceName, newService, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the common web ui Ingresses already exist. If not, create a new one.
	err = r.reconcileIngresses(instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the Certificates already exist, if not create new ones
	err = r.reconcileCertificates(instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	if needToRequeue {
		// one or more resources was created, so requeue the request
		reqLogger.Info("Requeue the request")
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Updating CommonWebUI staus")

	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(res.LabelsForSelector(res.DaemonSetName, commonwebuiserviceCrType, instance.Name)),
	}
	if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "CommonWebUI.Namespace", instance.Namespace, "CommonWebUI.Name", res.DaemonSetName)
		return reconcile.Result{}, err
	}
	podNames := res.GetPodNames(podList.Items)

	//update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update CommonWebUI status")
			return reconcile.Result{}, err
		}
	}

	navConfigCRD := &apiextv1beta.CustomResourceDefinition{}
	recResult, err := r.handleCRD(instance, navConfigCRD)
	if err != nil {
		return recResult, err
	}

	err = r.reconcileCr(instance)
	if err != nil {
		reqLogger.Error(err, "Error creating custom resource")
	}

	reqLogger.Info("CS??? all done")
	return reconcile.Result{}, nil
}

func (r *ReconcileCommonWebUI) reconcileConfigMaps(instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileConfiMaps", "instance.Name", instance.Name)

	reqLogger.Info("checking log4js config map Service")
	// Check if the log4js config map already exists, if not create a new one
	currentConfigMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: res.Log4jsConfigMap, Namespace: instance.Namespace}, currentConfigMap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new ConfigMap
		newConfigMap := res.Log4jsConfigMapUI(instance)

		err = controllerutil.SetControllerReference(instance, newConfigMap, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Failed to set owner for log4js config map", "Namespace", newConfigMap.Namespace,
				"Name", newConfigMap.Name)
			return err
		}

		reqLogger.Info("Creating a log4js config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
		err = r.client.Create(context.TODO(), newConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to create a config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
			return err
		}
		// Service created successfully - return and requeue
		*needToRequeue = true
	} else if err != nil {
		reqLogger.Error(err, "Failed to get log 4js config map")
		return err
	}

	reqLogger.Info("got log4js config map")

	return nil

}

func (r *ReconcileCommonWebUI) newDaemonSetForCR(instance *operatorsv1alpha1.CommonWebUI) (*appsv1.DaemonSet, error) {
	reqLogger := log.WithValues("func", "newDaemonSetForCR", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.DaemonSetName)
	selectorLabels := res.LabelsForSelector(res.DaemonSetName, commonwebuiserviceCrType, instance.Name)
	podLabels := res.LabelsForPodMetadata(res.DaemonSetName, commonwebuiserviceCrType, instance.Name)
	Annotations := res.DeamonSetAnnotations
	imageRegistry := instance.Spec.CommonWebUIConfig.ImageRegistry
	imageTag := instance.Spec.CommonWebUIConfig.ImageTag
	if imageRegistry == "" {
		imageRegistry = res.DefaultImageRegistry
	}
	if imageTag == "" {
		imageTag = res.DefaultImageTag
	}
	image := res.GetImageID(imageRegistry, res.DefaultImageName, imageTag, "", "COMMON_WEB_UI_IMAGE_TAG_OR_SHA")
	reqLogger.Info("CS??? default Image=" + image)

	commonVolume := append(commonVolume, res.Log4jsVolume)
	commonVolumes := append(commonVolume, res.ClusterCaVolume)
	commonVolumes = append(commonVolumes, res.UICertVolume)

	commonwebuiContainer := res.CommonContainer
	commonwebuiContainer.Image = image
	commonwebuiContainer.Name = res.DaemonSetName
	commonwebuiContainer.Env[1].Value = instance.Spec.GlobalUIConfig.RouterURL
	commonwebuiContainer.Env[3].Value = instance.Spec.GlobalUIConfig.IdentityProviderURL
	commonwebuiContainer.Env[4].Value = instance.Spec.GlobalUIConfig.AuthServiceURL
	commonwebuiContainer.Env[7].Value = instance.Spec.GlobalUIConfig.CloudPakVersion
	commonwebuiContainer.Env[8].Value = instance.Spec.GlobalUIConfig.DefaultAdminUser
	commonwebuiContainer.Env[9].Value = instance.Spec.GlobalUIConfig.ClusterName
	commonwebuiContainer.Env[10].Value = instance.Spec.GlobalUIConfig.DefaultAuth
	commonwebuiContainer.Env[11].Value = instance.Spec.GlobalUIConfig.EnterpriseLDAP
	commonwebuiContainer.Env[12].Value = instance.Spec.GlobalUIConfig.EnterpriseSAML
	commonwebuiContainer.Env[13].Value = instance.Spec.GlobalUIConfig.OSAuth

	daemon := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.DaemonSetName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDaemonSet{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: res.GetServiceAccountName(),
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "beta.kubernetes.io/arch",
												Operator: corev1.NodeSelectorOpIn,
												Values:   []string{gorun.GOARCH},
											},
										},
									},
								},
							},
						},
					},
					Volumes:                       commonVolumes,
					TerminationGracePeriodSeconds: &res.Seconds60,
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
					Containers: []corev1.Container{
						commonwebuiContainer,
					},
				},
			},
		},
	}
	// Set Commonsvcsuiservice instance as the owner and controller of the DaemonSet
	err := controllerutil.SetControllerReference(instance, daemon, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for common ui DaemonSet")
		return nil, err
	}
	return daemon, nil
}

// Check if the Common web ui Service already exist. If not, create a new one.
// This function was created to reduce the cyclomatic complexity :)
func (r *ReconcileCommonWebUI) serviceForUI(instance *operatorsv1alpha1.CommonWebUI) (*corev1.Service, error) {
	reqLogger := log.WithValues("func", "serviceForCommonWebUI", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.ServiceName)
	metaLabels["kubernetes.io/cluster-service"] = "true"
	metaLabels["kubernetes.io/name"] = instance.Spec.CommonWebUIConfig.ServiceName
	metaLabels["app"] = instance.Spec.CommonWebUIConfig.ServiceName
	selectorLabels := res.LabelsForSelector(res.ServiceName, commonwebuiserviceCrType, instance.Name)

	reqLogger.Info("CS??? Entry")
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.CommonWebUIConfig.ServiceName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: instance.Spec.CommonWebUIConfig.ServiceName,
					Port: 3000,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 3000,
					},
				},
			},
			Selector: selectorLabels,
		},
	}
	// Set Commonsvcsuiservice instance as the owner and controller of the DaemonSet
	err := controllerutil.SetControllerReference(instance, service, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner service")
		return nil, err
	}
	return service, nil
}

// Check if the common web ui Ingresses already exist. If not, create a new one.
// This function was created to reduce the cyclomatic complexity :)
func (r *ReconcileCommonWebUI) reconcileIngresses(instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileIngresses", "instance.Name", instance.Name)

	reqLogger.Info("checking  common web ui api ingress")
	// Define a new Ingress
	newAPIIngress := res.APIIngressForCommonWebUI(instance)
	// Set instance as the owner and controller of the ingress
	err := controllerutil.SetControllerReference(instance, newAPIIngress, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for api ingress")
		return nil
	}
	err = res.ReconcileIngress(r.client, instance.Namespace, res.APIIngress, newAPIIngress, needToRequeue)
	if err != nil {
		return err
	}
	reqLogger.Info("got common web ui api Ingress, checking common web ui callback Ingress")

	// Define a new Ingress
	newCallbackIngress := res.CallbackIngressForCommonWebUI(instance)
	// Set instance as the owner and controller of the ingress
	callbackErr := controllerutil.SetControllerReference(instance, newCallbackIngress, r.scheme)
	if callbackErr != nil {
		reqLogger.Error(callbackErr, "Failed to set owner for callback ingress")
		return nil
	}
	callbackErr = res.ReconcileIngress(r.client, instance.Namespace, res.CallbackIngress, newCallbackIngress, needToRequeue)
	if callbackErr != nil {
		return err
	}
	reqLogger.Info("got common web ui callback Ingress, checking common web ui nav Ingress")

	// Define a new Ingress
	newNavIngress := res.NavIngressForCommonWebUI(instance)
	// Set instance as the owner and controller of the ingress
	navErr := controllerutil.SetControllerReference(instance, newNavIngress, r.scheme)
	if navErr != nil {
		reqLogger.Error(err, "Failed to set owner for Nav ingress")
		return nil
	}
	navErr = res.ReconcileIngress(r.client, instance.Namespace, res.NavIngress, newNavIngress, needToRequeue)
	if navErr != nil {
		return err
	}
	reqLogger.Info("got common web ui nav Ingress")

	return nil
}

func (r *ReconcileCommonWebUI) handleCRD(instance *operatorsv1alpha1.CommonWebUI, currentCRD *apiextv1beta.CustomResourceDefinition) (reconcile.Result, error) {
	reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	reqLogger.Info("ABOUT TO HANDLE THIS CRD")
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "navconfigurations.foundation.ibm.com", Namespace: ""}, currentCRD)
	if err != nil && errors.IsNotFound(err) {
		// Define CRD
		newCRD := r.newNavconfgCRD()
		reqLogger.Info("Creating a new CRD", "CRD.Namespace", instance.Namespace, "CRD.Name", "clients.oidc.security.ibm.com")
		err = r.client.Create(context.TODO(), newCRD)
		if err != nil {
			reqLogger.Error(err, "Failed to create new CRD", "CRD.Namespace", instance.Namespace, "CRD.Name", "clients.oidc.security.ibm.com")
			return reconcile.Result{}, err
		}
		// new CRD created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get CRD")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileCommonWebUI) newNavconfgCRD() *apiextv1beta.CustomResourceDefinition {
	newCRD := &apiextv1beta.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "navconfigurations.foundation.ibm.com",
			Labels:    res.LabelsForMetadata(res.DaemonSetName),
			Namespace: "ibm-common-services",
		},
		Spec: apiextv1beta.CustomResourceDefinitionSpec{
			Scope:   "Namespaced",
			Group:   "foundation.ibm.com",
			Version: "v1",
			Names: apiextv1beta.CustomResourceDefinitionNames{
				Kind:       "NavConfiguration",
				Singular:   "navconfiguration",
				Plural:     "navconfigurations",
				ShortNames: []string{"navconfig"},
			},
			Validation: &apiextv1beta.CustomResourceValidation{
				OpenAPIV3Schema: &apiextv1beta.JSONSchemaProps{
					Properties: res.GetNavConfigContent(),
				},
			},
		},
	}

	return newCRD
}

func (r *ReconcileCommonWebUI) reconcileCr(instance *operatorsv1alpha1.CommonWebUI) error {
	reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	reqLogger.Info("RECONCILING CR")

	namespace := instance.Namespace
	// Empty interface of type Array to hold the crs
	var crTemplates []map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	crTemplatesErr := json.Unmarshal([]byte(res.CrTemplates), &crTemplates)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	for _, crTemplate := range crTemplates {
		var unstruct unstructured.Unstructured
		unstruct.Object = crTemplate
		name := unstruct.Object["metadata"].(map[string]interface{})["name"].(string)

		getError := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}, &unstruct)

		if getError != nil && !errors.IsNotFound(getError) {
			reqLogger.Error(getError, "Failed to get the CR")
			continue
		} else if errors.IsNotFound(getError) {
			// Create Custom resource
			if createErr := r.createCustomResource(unstruct, name, namespace); createErr != nil {
				reqLogger.Error(createErr, "Failed to create CR")
				return createErr
			}
		} else {
			reqLogger.Info("Skipping CR creation")
		}
	}
	return nil
}

func (r *ReconcileCommonWebUI) createCustomResource(unstruct unstructured.Unstructured, name, namespace string) error {
	reqLogger := log.WithValues("CR namespace", namespace, "CR name", name)
	reqLogger.Info("creating a CR ", name)
	unstruct.Object["metadata"].(map[string]interface{})["namespace"] = namespace
	crCreateErr := r.client.Create(context.TODO(), &unstruct)
	if crCreateErr != nil && !errors.IsAlreadyExists(crCreateErr) {
		reqLogger.Error(crCreateErr, "Failed to Create the Custom Resource")
		return crCreateErr
	}
	return nil
}

func (r *ReconcileCommonWebUI) reconcileCertificates(instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileCertificates", "instance.Name", instance.Name)

	certificateList := []res.CertificateData{
		res.UICertificateData,
	}

	for _, certData := range certificateList {
		reqLogger.Info("Checking Certificate", "Certificate.Name", certData.Name)
		newCertificate := res.BuildCertificate(instance.Namespace, "", certData)
		// Set CommonWebUI instance as the owner and controller of the Certificate
		err := controllerutil.SetControllerReference(instance, newCertificate, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Failed to set owner for Certificate", "Certificate.Namespace", newCertificate.Namespace,
				"Certificate.Name", newCertificate.Name)
			return err
		}
		err = res.ReconcileCertificate(r.client, instance.Namespace, certData.Name, newCertificate, needToRequeue)
		if err != nil {
			return err
		}
	}
	return nil
}
