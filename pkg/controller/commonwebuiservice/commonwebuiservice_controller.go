package commonwebuiservice

import (
	"context"
	gorun "runtime"

	res "github.com/ibm/metering-operator/pkg/resources"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	return &ReconcileCommonWebUIService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("commonwebuiservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CommonWebUIService
	err = c.Watch(&source.Kind{Type: &operatorsv1alpha1.CommonWebUIService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorsv1alpha1.CommonWebUIService{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource "Daemonset" and requeue the owner CommonWebUIService
	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.CommonWebUIService{},
	})
	if err != nil {
		return err
	}

	// // Watch for changes to secondary resource "Service" and requeue the owner CommonWebUIService
	// err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &operatorv1alpha1.CommonWebUIService{},
	// })
	// if err != nil {
	// 	return err
	// }

	// Watch for changes to secondary resource "Ingress" and requeue the owner CommonWebUIService
	// err = c.Watch(&source.Kind{Type: &netv1.Ingress{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &operatorv1alpha1.CommonWebUIService{},
	// })
	// if err != nil {
	// 	return err
	// }


	return nil
}

// blank assignment to verify that ReconcileCommonWebUIService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCommonWebUIService{}

// ReconcileCommonWebUIService reconciles a CommonWebUIService object
type ReconcileCommonWebUIService struct {
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
func (r *ReconcileCommonWebUIService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CommonWebUIService")

	// Fetch the CommonWebUIService CR instance
	instance := &operatorsv1alpha1.CommonWebUIService{}

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

	// Set CommonWebUIService instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// Check if the DaemonSet already exists, if not create a new one
	currentDaemonSet := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: res.DaemonSetName, Namespace: instance.Namespace}, currentDaemonSet)
	if err != nil && errors.IsNotFound(err) {
		// Define a new DaemonSet
		newDaemonSet := r.newDaemonSetForCR(instance)
		reqLogger.Info("Creating a new Rdr DaemonSet", "DaemonSet.Namespace", newDaemonSet.Namespace, "DaemonSet.Name", newDaemonSet.Name)
		err = r.client.Create(context.TODO(), newDaemonSet)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Rdr DaemonSet", "DaemonSet.Namespace", newDaemonSet.Namespace,
				"DaemonSet.Name", newDaemonSet.Name)
			return reconcile.Result{}, err
		}
		// DaemonSet created successfully - return and requeue
		needToRequeue = true
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Rdr DaemonSet")
		return reconcile.Result{}, err
	}
	

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func newDaemonSetForCR(instance *operatorsv1alpha1.CommonWebUIService) *appsv1.DaemonSet {
	reqLogger := log.WithValues("func", "daemonForReader", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.DaemonSetName)
	selectorLabels := res.LabelsForSelector(res.DaemonSetName, commonwebuiserviceCrType, instance.Name)
	podLabels := res.LabelsForPodMetadata(res.DaemonSetName, commonwebuiserviceCrType, instance.Name)

	var image string
	if instance.Spec.CommonUIConfig.ImageRegistry == "" {
		image = res.DefaultImageRegistry + "/" + res.DefaultImageName + ":" + res.DefaultImageTag
		reqLogger.Info("CS??? default rdrImage=" + image)
	} else {
		image = instance.Spec.CommonUIConfig.ImageRegistry + "/" + res.DefaultImageName + ":" + res.DefaultImageTag
		reqLogger.Info("CS??? rdrImage=" + image)
	}


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
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: corev1.PodSpec{
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
					Volumes: rdrVolumes,
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
						rdrMainContainer,
					},
				},
			},
	}
	// Set Commonsvcsuiserive instance as the owner and controller of the DaemonSet
	err := controllerutil.SetControllerReference(instance, daemon, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for Rdr DaemonSet")
		return nil
	}
	return daemon
}