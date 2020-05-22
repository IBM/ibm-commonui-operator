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

package legacyheaderservice

import (
	"context"
	gorun "runtime"

	res "github.com/ibm/ibm-commonui-operator/pkg/resources"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"

	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const legacyheaderCrType = "legacyheader_cr"

var commonVolume = []corev1.Volume{}

var log = logf.Log.WithName("controller_legacyheader")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new LegacyHeaderService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLegacyHeader{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("legacyheader-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource LegacyHeaderService
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.LegacyHeader{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource ConfigMap and requeue the owner LegacyHeader
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.LegacyHeader{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource "Daemonset" and requeue the owner LegacyHeader
	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.LegacyHeader{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Service" and requeue the owner LegacyHeader
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.LegacyHeader{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Ingress" and requeue the owner
	err = c.Watch(&source.Kind{Type: &netv1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.LegacyHeader{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileLegacyHeader implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLegacyHeader{}

// ReconcileLegacyHeader reconciles a LegacyHeaderService object
type ReconcileLegacyHeader struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a LegacyHeader object and makes changes based on the state read
// and what is in the LegacyHeader.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileLegacyHeader) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling LegacyHeaderService")

	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	// Fetch the LegacyHeaderService instance
	instance := &operatorsv1alpha1.LegacyHeader{}

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
	reqLogger.Info("got LegacyHeaderService instance, version=" + opVersion)

	// set a default Status value
	if len(instance.Status.Nodes) == 0 {
		instance.Status.Nodes = res.DefaultStatusForCR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set LegacyHeader default status")
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
	err = res.ReconcileDaemonSet(r.client, instance.Namespace, res.LegacyReleaseName, newDaemonSet, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the platform header Service already exist. If not, create a new one.
	newService, err := r.serviceForUI(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = res.ReconcileService(r.client, instance.Namespace, res.LegacyReleaseName, newService, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Check if the platform header Ingress already exist. If not, create a new one.
	err = r.reconcileIngress(instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	if needToRequeue {
		// one or more resources was created, so requeue the request
		reqLogger.Info("Requeue the request")
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Updating LegacyHeader staus")

	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(res.LabelsForSelector(res.LegacyReleaseName, legacyheaderCrType, instance.Name)),
	}
	if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "LegacyHeader.Namespace", instance.Namespace, "LegacyHeader.Name", res.LegacyReleaseName)
		return reconcile.Result{}, err
	}
	podNames := res.GetPodNames(podList.Items)

	//update status.podNames if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update LegacyHeader status")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("got Services, checking Certificates")
	// Resources exists - don't requeue
	reqLogger.Info("CS??? all done")
	return reconcile.Result{}, nil

}

func (r *ReconcileLegacyHeader) reconcileConfigMaps(instance *operatorsv1alpha1.LegacyHeader, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileConfiMaps", "instance.Name", instance.Name)

	// Check if the common config map already exists, if not create a new one
	currentConfigMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: res.CommonConfigMap, Namespace: instance.Namespace}, currentConfigMap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new ConfigMap
		newConfigMap := res.CommonConfigMapUI(instance)

		err = controllerutil.SetControllerReference(instance, newConfigMap, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Failed to set owner for common config map", "Namespace", newConfigMap.Namespace,
				"Name", newConfigMap.Name)
			return err
		}

		reqLogger.Info("Creating a common config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
		err = r.client.Create(context.TODO(), newConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to create a config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
			return err
		}
		// Service created successfully - return and requeue
		*needToRequeue = true
	} else if err != nil {
		reqLogger.Error(err, "Failed to get common config map")
		return err
	}

	reqLogger.Info("got common config map")

	return nil

}

func (r *ReconcileLegacyHeader) newDaemonSetForCR(instance *operatorsv1alpha1.LegacyHeader) (*appsv1.DaemonSet, error) {
	reqLogger := log.WithValues("func", "newDaemonSetForCR", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.LegacyReleaseName)
	selectorLabels := res.LabelsForSelector(res.LegacyReleaseName, legacyheaderCrType, instance.Name)
	podLabels := res.LabelsForPodMetadata(res.LegacyReleaseName, legacyheaderCrType, instance.Name)
	Annotations := res.DeamonSetAnnotations
	imageRegistry := instance.Spec.LegacyConfig.ImageRegistry
	imageTag := instance.Spec.LegacyConfig.ImageTag
	if imageRegistry == "" {
		imageRegistry = res.DefaultImageRegistry
	}
	if imageTag == "" {
		imageTag = res.DefaultImageTag
	}
	image := res.GetImageID(imageRegistry, res.LegacyImageName, imageTag, "", "LEGACYHEADER_IMAGE_TAG_OR_SHA")
	reqLogger.Info("CS??? default Image=" + image)

	commonVolume := append(commonVolume, res.Log4jsVolume)
	commonVolumes := append(commonVolume, res.ClusterCaVolume)

	legacyContainer := res.CommonContainer
	legacyContainer.Image = image
	legacyContainer.Name = res.LegacyReleaseName
	legacyContainer.Env[1].Value = instance.Spec.LegacyGlobalUIConfig.RouterURL
	legacyContainer.Env[3].Value = instance.Spec.LegacyGlobalUIConfig.IdentityProviderURL
	legacyContainer.Env[4].Value = instance.Spec.LegacyGlobalUIConfig.AuthServiceURL
	legacyContainer.Env[7].Value = instance.Spec.LegacyGlobalUIConfig.CloudPakVersion
	legacyContainer.Env[8].Value = instance.Spec.LegacyGlobalUIConfig.DefaultAdminUser
	legacyContainer.Env[9].Value = instance.Spec.LegacyGlobalUIConfig.ClusterName

	daemon := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.LegacyReleaseName,
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
						legacyContainer,
					},
				},
			},
		},
	}
	// Set Commonsvcsuiservice instance as the owner and controller of the DaemonSet
	err := controllerutil.SetControllerReference(instance, daemon, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for legacy DaemonSet")
		return nil, err
	}
	return daemon, nil
}

// Check if the Common web ui Service already exist. If not, create a new one.
// This function was created to reduce the cyclomatic complexity :)
func (r *ReconcileLegacyHeader) serviceForUI(instance *operatorsv1alpha1.LegacyHeader) (*corev1.Service, error) {
	reqLogger := log.WithValues("func", "serviceForLegacyUI", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.LegacyReleaseName)
	metaLabels["kubernetes.io/cluster-service"] = "true"
	metaLabels["kubernetes.io/name"] = instance.Spec.LegacyConfig.ServiceName
	metaLabels["app"] = instance.Spec.LegacyConfig.ServiceName
	selectorLabels := res.LabelsForSelector(res.LegacyReleaseName, legacyheaderCrType, instance.Name)

	reqLogger.Info("CS??? Entry")
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.LegacyConfig.ServiceName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: instance.Spec.LegacyConfig.ServiceName,
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
	err := controllerutil.SetControllerReference(instance, service, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner service")
		return nil, err
	}
	return service, nil
}

// Check if the lagacy header Ingresses already exist. If not, create a new one.
// This function was created to reduce the cyclomatic complexity :)
func (r *ReconcileLegacyHeader) reconcileIngress(instance *operatorsv1alpha1.LegacyHeader, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileIngress", "instance.Name", instance.Name)
	// Define a new Ingress
	newNavIngress := res.IngressForLegacyUI(instance)
	// Set instance as the owner and controller of the ingress
	err := controllerutil.SetControllerReference(instance, newNavIngress, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for Nav ingress")
		return nil
	}
	err = res.ReconcileIngress(r.client, instance.Namespace, res.LegacyReleaseName, newNavIngress, needToRequeue)
	if err != nil {
		return err
	}
	reqLogger.Info("got legacy header Ingress")

	return nil
}
