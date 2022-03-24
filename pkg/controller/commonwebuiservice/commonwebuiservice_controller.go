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
package commonwebuiservice

import (
	"context"
	"encoding/json"
	"strconv"

	res "github.com/ibm/ibm-commonui-operator/pkg/resources"
	ver "github.com/ibm/ibm-commonui-operator/version"
	"k8s.io/apimachinery/pkg/util/intstr"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"
	certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
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

var log = logf.Log.WithName("controller_commonwebuiservice")

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

	// Watch for changes to secondary resource "Deployment" and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
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
func (r *ReconcileCommonWebUI) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUI")

	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	// Fetch the CommonWebUIService CR instance
	instance := &operatorsv1alpha1.CommonWebUI{}

	err := r.client.Get(ctx, request.NamespacedName, instance)
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
		instance.Status.Versions = operatorsv1alpha1.Versions{Reconciled: ver.Version}
		err = r.client.Status().Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set CommonWebUI default status")
			return reconcile.Result{}, err
		}
	}

	//Reconcile to see if Zen is enabled
	isZen := r.adminHubOnZen(ctx, instance.Namespace)

	//Check to see kubernetes cluster type
	isCncf := r.getKubernetesClusterType(ctx, instance.Namespace)

	// Check if the config maps already exist. If not, create a new one.
	err = r.reconcileConfigMaps(ctx, instance, res.Log4jsConfigMap, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the UI Deployment already exists, if not create a new one
	newDeployment, err := r.deploymentForUI(instance, isZen, isCncf)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = res.ReconcileDeployment(ctx, r.client, instance.Namespace, res.DeploymentName, newDeployment, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the common web ui Service already exist. If not, create a new one.
	newService, err := r.serviceForUI(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = res.ReconcileService(ctx, r.client, instance.Namespace, res.ServiceName, newService, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the common web ui Ingresses already exist. If not, create a new one.
	err = r.reconcileIngresses(ctx, instance, &needToRequeue, isCncf)
	if err != nil {
		return reconcile.Result{}, err
	}

	//Check if CR already exists. If not, create a new one
	err = r.reconcileCr(ctx, instance)
	if err != nil {
		reqLogger.Error(err, "Error creating custom resource")
	}

	// Check if the Certificates already exist, if not create new ones
	err = r.reconcileCertificates(ctx, instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.updateCustomResource(ctx, instance, res.CommonWebUICr)
	if err != nil {
		reqLogger.Error(err, "Failed updating navconfig CR")
	}

	err = r.updateCustomResource(ctx, instance, res.Cp4iCr)
	if err != nil {
		reqLogger.Error(err, "Failed updating icp4i navconfig CR")
	}

	// For 1.3.0 operator version check if daemonSet and navconfig crd exits on upgrade and delete if so
	r.deleteDaemonSet(ctx, instance)

	if needToRequeue {
		// one or more resources was created, so requeue the request
		reqLogger.Info("Requeue the request")
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Updating CommonWebUI staus")

	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(res.LabelsForSelector(res.DeploymentName, commonwebuiserviceCrType, instance.Name)),
	}
	if err = r.client.List(ctx, podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "CommonWebUI.Namespace", instance.Namespace, "CommonWebUI.Name", res.DeploymentName)
		return reconcile.Result{}, err
	}
	podNames := res.GetPodNames(podList.Items)

	//update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		instance.Status.Versions = operatorsv1alpha1.Versions{Reconciled: ver.Version}
		err := r.client.Status().Update(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update CommonWebUI status")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("CS??? all done")
	return reconcile.Result{}, nil
}

func (r *ReconcileCommonWebUI) adminHubOnZen(ctx context.Context, namespace string) bool {
	reqLogger := log.WithValues("func", "adminHubOnZen")
	reqLogger.Info("Checking zen optional install condition in commonui controller")

	zenDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zen-core",
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: "zen-core", Namespace: namespace}, zenDeployment)

	if getError == nil {
		reqLogger.Info("Got ZEN Deployment in commonui controller")
		return true
	}
	if errors.IsNotFound(getError) {
		reqLogger.Info("ZEN deployment not found in commonui controller")
	} else {
		reqLogger.Error(getError, "Error getting ZEN deployment  in commonui controller")
	}
	return false
}

func (r *ReconcileCommonWebUI) getKubernetesClusterType(ctx context.Context, namespace string) bool {
	reqLogger := log.WithValues("func", "isCncf")
	reqLogger.Info("Checking kubernetes cluster type")

	ibmProjectK := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibm-cpp-config",
			Namespace: namespace,
		},
	}
	getError := r.client.Get(ctx, types.NamespacedName{Name: "ibm-cpp-config", Namespace: namespace}, ibmProjectK)

	if getError == nil {
		reqLogger.Info("Got ibm project k config map")
		clusterType := ibmProjectK.Data["kubernetes_cluster_type"]
		if clusterType == "cncf" {
			reqLogger.Info("Kubernetes cluster type is " + clusterType)
			return true
		}
	}

	if errors.IsNotFound(getError) {
		reqLogger.Info("ibm project k config map not found in cs namepace")
	} else {
		reqLogger.Error(getError, "error getting ibm project k config map in cs namepace")
	}

	return false
}

func (r *ReconcileCommonWebUI) reconcileConfigMaps(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI, nameOfCM string, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileConfiMaps", "instance.Name", instance.Name)

	reqLogger.Info("checking log4js config map Service")
	// Check if the log4js config map already exists, if not create a new one
	currentConfigMap := &corev1.ConfigMap{}
	err := r.client.Get(ctx, types.NamespacedName{Name: nameOfCM, Namespace: instance.Namespace}, currentConfigMap)
	if err != nil && errors.IsNotFound(err) {
		// Define a new ConfigMap
		newConfigMap := &corev1.ConfigMap{}
		if nameOfCM == res.Log4jsConfigMap {
			newConfigMap = res.Log4jsConfigMapUI(instance)
		}
		err = controllerutil.SetControllerReference(instance, newConfigMap, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Failed to set owner for log4js config map", "Namespace", newConfigMap.Namespace,
				"Name", newConfigMap.Name)
			return err
		}

		reqLogger.Info("Creating a log4js config map", "Namespace", newConfigMap.Namespace, "Name", newConfigMap.Name)
		err = r.client.Create(ctx, newConfigMap)
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

func (r *ReconcileCommonWebUI) deploymentForUI(instance *operatorsv1alpha1.CommonWebUI, isZen bool, isCncf bool) (*appsv1.Deployment, error) {
	// CommonMainVolumeMounts will be added by the controller
	commonUIVolumeMounts := []corev1.VolumeMount{
		{
			Name:      res.Log4jsVolumeName,
			MountPath: "/etc/config",
		},
		{
			Name:      res.ClusterCaVolumeName,
			MountPath: "/opt/ibm/platform-header/certs",
		},
		{
			Name:      res.UICertVolumeName,
			MountPath: "/certs/common-web-ui",
		},
		{
			Name:      res.InternalTLSVolumeName,
			MountPath: "/etc/internal-tls",
		},
	}
	var commonVolume = []corev1.Volume{}
	reqLogger := log.WithValues("func", "newDeploymentForUI", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.DeploymentName)
	selectorLabels := res.LabelsForSelector(res.DeploymentName, commonwebuiserviceCrType, instance.Name)
	podLabels := res.LabelsForPodMetadata(res.DeploymentName, commonwebuiserviceCrType, instance.Name)
	Annotations := res.DeploymentAnnotations
	var replicas int32 = instance.Spec.Replicas
	var cpuLimits, cpuMemory, reqLimits, reqMemory int64
	var errLim error

	if replicas == 0 {
		replicas = 1
	}

	if instance.Spec.Resources.Limits.CPULimits != "" {
		limits := instance.Spec.Resources.Limits.CPULimits
		cpuLimits, errLim = strconv.ParseInt(limits[0:len(limits)-1], 10, 64)
		if errLim != nil {
			cpuLimits = 1000
		}
	} else {
		cpuLimits = 1000
	}

	if instance.Spec.Resources.Limits.CPUMemory != "" {
		memory := instance.Spec.Resources.Limits.CPUMemory
		cpuMemory, errLim = strconv.ParseInt(memory[0:len(memory)-2], 10, 64)
		if errLim != nil {
			cpuMemory = 512
		}
	} else {
		cpuMemory = 512
	}

	if instance.Spec.Resources.Requests.RequestLimits != "" {
		limits := instance.Spec.Resources.Requests.RequestLimits
		reqLimits, errLim = strconv.ParseInt(limits[0:len(limits)-1], 10, 64)
		if errLim != nil {
			reqLimits = 300
		}
	} else {
		reqLimits = 300
	}

	if instance.Spec.Resources.Requests.RequestMemory != "" {
		memory := instance.Spec.Resources.Requests.RequestMemory
		reqMemory, errLim = strconv.ParseInt(memory[0:len(memory)-2], 10, 64)
		if errLim != nil {
			reqMemory = 512
		}
	} else {
		reqMemory = 512
	}

	imageRegistry := instance.Spec.CommonWebUIConfig.ImageRegistry
	imageTag := instance.Spec.CommonWebUIConfig.ImageTag
	if imageRegistry == "" {
		imageRegistry = res.DefaultImageRegistry
	}
	if imageTag == "" {
		imageTag = res.DefaultImageTag
	}
	image := res.GetImageID(imageRegistry, res.DefaultImageName, imageTag, "", "COMMON_WEB_UI_IMAGE")
	reqLogger.Info("CS??? default Image=" + image)

	commonVolume = append(commonVolume, res.Log4jsVolume)
	commonVolumes := append(commonVolume, res.ClusterCaVolume)
	commonVolumes = append(commonVolumes, res.UICertVolume)
	allCommonVolumes := append(commonVolumes, res.InternalTLSVolume)

	commonwebuiContainer := res.CommonContainer
	commonwebuiContainer.Image = image
	commonwebuiContainer.Name = res.DaemonSetName
	commonwebuiContainer.Env[6].Value = instance.Spec.GlobalUIConfig.CloudPakVersion
	commonwebuiContainer.Env[8].Value = instance.Spec.GlobalUIConfig.DefaultAuth
	commonwebuiContainer.Env[9].Value = instance.Spec.GlobalUIConfig.EnterpriseLDAP
	commonwebuiContainer.Env[10].Value = instance.Spec.GlobalUIConfig.EnterpriseSAML
	commonwebuiContainer.Env[11].Value = instance.Spec.GlobalUIConfig.OSAuth
	commonwebuiContainer.Env[19].Value = instance.Spec.CommonWebUIConfig.LandingPage
	commonwebuiContainer.Resources.Limits["cpu"] = *resource.NewMilliQuantity(cpuLimits, resource.DecimalSI)
	commonwebuiContainer.Resources.Limits["memory"] = *resource.NewQuantity(cpuMemory*1024*1024, resource.BinarySI)
	commonwebuiContainer.Resources.Requests["cpu"] = *resource.NewMilliQuantity(reqLimits, resource.DecimalSI)
	commonwebuiContainer.Resources.Requests["memory"] = *resource.NewQuantity(reqMemory*1024*1024, resource.BinarySI)
	commonwebuiContainer.VolumeMounts = commonUIVolumeMounts

	if isZen {
		reqLogger.Info("Setting use zen to true in container def")
		commonwebuiContainer.Env[22].Value = "true"
	} else {
		reqLogger.Info("Setting use zen to false in container def")
		commonwebuiContainer.Env[22].Value = "false"
	}
	commonwebuiContainer.Env[23].Value = instance.Spec.Version

	if isCncf {
		reqLogger.Info("Setting cluster type env var to cncf")
		commonwebuiContainer.Env[24].Value = "cncf"
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.DeploymentName,
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
					Annotations: Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            res.GetServiceAccountName(),
					HostNetwork:                   false,
					HostPID:                       false,
					HostIPC:                       false,
					TerminationGracePeriodSeconds: &res.Seconds60,
					TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
						{
							MaxSkew:           1,
							TopologyKey:       "topology.kubernetes.io/zone",
							WhenUnsatisfiable: corev1.ScheduleAnyway,
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"k8s-app": res.DeploymentName,
								},
							},
						},
						{
							MaxSkew:           1,
							TopologyKey:       "topology.kubernetes.io/region",
							WhenUnsatisfiable: corev1.ScheduleAnyway,
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"k8s-app": res.DeploymentName,
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
												Values:   res.ArchitectureList,
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
													Values:   []string{res.DeploymentName},
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
					Volumes: allCommonVolumes,
					Containers: []corev1.Container{
						commonwebuiContainer,
					},
				},
			},
		},
	}
	// Set CommonUI instance as the owner and controller of the Deployment
	err := controllerutil.SetControllerReference(instance, deployment, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for UI Deployment")
		return nil, err
	}
	return deployment, nil
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
func (r *ReconcileCommonWebUI) reconcileIngresses(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool, isCncf bool) error {
	reqLogger := log.WithValues("func", "reconcileIngresses", "instance.Name", instance.Name)

	reqLogger.Info("checking  common web ui api ingress")
	// Define a new Ingress
	newAPIIngress := res.APIIngressForCommonWebUI(instance)
	if isCncf {
		newAPIIngress = res.APIIngressForCommonWebUICncf(instance)
	}
	// Set instance as the owner and controller of the ingress
	err := controllerutil.SetControllerReference(instance, newAPIIngress, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for api ingress")
		return nil
	}
	err = res.ReconcileIngress(ctx, r.client, instance.Namespace, res.APIIngress, newAPIIngress, needToRequeue)
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
	callbackErr = res.ReconcileIngress(ctx, r.client, instance.Namespace, res.CallbackIngress, newCallbackIngress, needToRequeue)
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
	navErr = res.ReconcileIngress(ctx, r.client, instance.Namespace, res.NavIngress, newNavIngress, needToRequeue)
	if navErr != nil {
		return err
	}
	reqLogger.Info("got common web ui nav Ingress")

	return nil
}

func (r *ReconcileCommonWebUI) reconcileCr(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI) error {
	reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	reqLogger.Info("RECONCILING CR")

	namespace := instance.Namespace
	var crTemplate map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	crTemplatesErr := json.Unmarshal([]byte(res.CrTemplates), &crTemplate)
	if crTemplatesErr != nil {
		reqLogger.Info("Failed to unmarshall crTemplates")
		return crTemplatesErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	name := unstruct.Object["metadata"].(map[string]interface{})["name"].(string)

	//Get CR and see if it exists
	getError := r.client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &unstruct)

	err1 := r.client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, instance)
	if err1 == nil {
		r.finalizerCr(ctx, instance, unstruct)
	}

	if getError != nil && !errors.IsNotFound(getError) {
		reqLogger.Error(getError, "Failed to get CR")
	} else {
		reqLogger.Info("Skipping CR creation")
	}

	return nil
}

func (r *ReconcileCommonWebUI) finalizerCr(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI, unstruct unstructured.Unstructured) {
	reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)

	finalizerName := "commonui.operators.ibm.com"
	finalizerName1 := "commonui1.operators.ibm.com"

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Add the finalizer to the metadata of the instance and update the object.
		if !containsString(instance.ObjectMeta.Finalizers, finalizerName) && !containsString(instance.ObjectMeta.Finalizers, finalizerName1) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName, finalizerName1)
			if err := r.client.Update(ctx, instance); err != nil {
				reqLogger.Error(err, "Failed to create finalizer")
			} else {
				reqLogger.Info("Created Finalizers")
			}
		}
	} else {
		// When the instance is being deleted. If finalizer is present
		if containsString(instance.ObjectMeta.Finalizers, finalizerName) {
			// Finalizer is present, so lets handle any external dependency - remove console link CR
			if err := r.client.Delete(ctx, &unstruct); err != nil {
				// if fails to delete the external dependency here, return with error
				reqLogger.Error(err, "Failed to delete Console Link CR")
			} else {
				reqLogger.Info("Deleted Console link CR")
			}

			// Remove our finalizer from the metadata of the object and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.client.Update(ctx, instance); err != nil {
				reqLogger.Error(err, "Failed to delete  Console link finalizer")
			} else {
				reqLogger.Info("Deleted Console link Finalizer")
			}
		} else if containsString(instance.ObjectMeta.Finalizers, finalizerName1) {
			// Finalizer is present, so lets handle any external dependency - remove console link CR
			if err := r.client.Delete(ctx, &unstruct); err != nil {
				// if fails to delete the external dependency here, return with error
				reqLogger.Error(err, "Failed to delete Redis CR")
			} else {
				reqLogger.Info("Deleted Redis CR")
			}

			// Remove our finalizer from the metadata of the object and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, finalizerName1)
			if err := r.client.Update(ctx, instance); err != nil {
				reqLogger.Error(err, "Failed to delete Redis finalizer")
			} else {
				reqLogger.Info("Deleted Redis Finalizer")
			}
		}
	}
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func (r *ReconcileCommonWebUI) reconcileCertificates(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI, needToRequeue *bool) error {
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
		err = res.ReconcileCertificate(ctx, r.client, instance.Namespace, certData.Name, newCertificate, needToRequeue)
		if err != nil {
			return err
		}
	}
	return nil
}

// delete the old common ui daemonset from an older version
func (r *ReconcileCommonWebUI) deleteDaemonSet(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI) {
	reqLogger := log.WithValues("func", "deleteDaemonSet", "instance.Name", instance.Name)
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.DaemonSetName,
			Namespace: res.DefaultNamespace,
		},
	}
	// check if the DaemonSet exists
	err := r.client.Get(ctx,
		types.NamespacedName{Name: res.DaemonSetName, Namespace: res.DefaultNamespace}, daemonSet)
	if err == nil {
		// DaemonSet found so delete it
		err := r.client.Delete(ctx, daemonSet)
		if err != nil {
			reqLogger.Error(err, "Failed to delete old common ui DaemonSet")
		} else {
			reqLogger.Info("Deleted old common ui DaemonSet")
		}
	} else if !errors.IsNotFound(err) {
		reqLogger.Error(err, "Failed to get old DaemonSet")
	}
}

func (r *ReconcileCommonWebUI) updateCustomResource(ctx context.Context, instance *operatorsv1alpha1.CommonWebUI, nameOfCR string) error {
	reqLogger := log.WithValues("Instance.Namespace", instance.Namespace, "Instance.Name", instance.Name)
	reqLogger.Info("UPDATE CUSTOM RESOURCE")
	namespace := instance.Namespace
	var crTemplate map[string]interface{}
	var jsonStringCr string
	// Unmarshal or Decode the JSON to the interface.
	if nameOfCR == res.CommonWebUICr {
		jsonStringCr = res.NavConfigCR
	} else {
		jsonStringCr = res.NavConfigCP4ICR
	}
	crTemplateErr := json.Unmarshal([]byte(jsonStringCr), &crTemplate)
	if crTemplateErr != nil {
		reqLogger.Info("Failed to unmarshall nav config cr")
		return crTemplateErr
	}
	var unstruct unstructured.Unstructured
	unstruct.Object = crTemplate
	getError := r.client.Get(ctx, types.NamespacedName{
		Name:      nameOfCR,
		Namespace: namespace,
	}, &unstruct)

	if getError == nil {
		reqLogger.Info("FOUND NAV CONFIG CR TRYING TO UPDATE")
		var currentTemplate map[string]interface{}
		crTemplateErr2 := json.Unmarshal([]byte(jsonStringCr), &currentTemplate)
		if crTemplateErr2 != nil {
			reqLogger.Info("Failed to unmarshall current nav config cr")
			return crTemplateErr2
		}
		var unstruct2 unstructured.Unstructured
		unstruct2.Object = currentTemplate
		navItems := unstruct2.Object["spec"].(map[string]interface{})["navItems"]
		var jsonData []byte
		jsonData, err := json.Marshal(navItems)
		if err != nil {
			reqLogger.Info("Failed to marshall navitems")
			return err
		}
		var updatedNavItems []map[string]interface{}
		//nolint
		navItemsErr := json.Unmarshal([]byte(jsonData), &updatedNavItems)
		if navItemsErr != nil {
			reqLogger.Info("Failed to unmarshall nav items array")
			return navItemsErr
		}
		for _, item := range updatedNavItems {
			if item["namespace"] != "" {
				item["namespace"] = namespace
			}
		}
		unstruct.Object["spec"].(map[string]interface{})["navItems"] = updatedNavItems

		if nameOfCR == res.CommonWebUICr {
			licenses := unstruct2.Object["spec"].(map[string]interface{})["about"].(map[string]interface{})["licenses"]
			unstruct.Object["spec"].(map[string]interface{})["about"].(map[string]interface{})["licenses"] = licenses
		}

		//Update the CR
		updateErr := r.client.Update(ctx, &unstruct)
		if updateErr == nil {
			reqLogger.Info("CLIENT UPDATED NAV CONFIG CR ")
		}
	}
	return nil
}
