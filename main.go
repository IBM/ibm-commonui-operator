/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	runtimescheme "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	certmgr "github.com/ibm/ibm-cert-manager-operator/apis/cert-manager/v1"
	certmgrv1alpha1 "github.com/ibm/ibm-cert-manager-operator/apis/certmanager/v1alpha1"
	cmmeta "github.com/ibm/ibm-cert-manager-operator/apis/meta.cert-manager/v1"
	routesv1 "github.com/openshift/api/route/v1"

	"github.com/IBM/controller-filtered-cache/filteredcache"
	operatorsv1alpha1 "github.com/IBM/ibm-commonui-operator/api/v1alpha1"
	im "github.com/IBM/ibm-commonui-operator/apis/operator/v1alpha1"
	commonwebuicontrollers "github.com/IBM/ibm-commonui-operator/controllers/commonwebui"
	res "github.com/IBM/ibm-commonui-operator/controllers/resources"
	"github.com/IBM/ibm-commonui-operator/version"
	//+kubebuilder:scaffold:imports
)

const CommonServiceName string = "common-service"

var log = logf.Log.WithName("main")

var (
	scheme   = runtimescheme.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func printVersion() {
	log.Info(fmt.Sprintf("Operator Version: %s", version.Version))
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func init() {
	// add default kubernetes schemes to controller
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	// add cert manager scheme to controller
	utilruntime.Must(certmgr.AddToScheme(scheme))
	utilruntime.Must(certmgrv1alpha1.AddToScheme(scheme))

	// add cert manager scheme to controller
	utilruntime.Must(cmmeta.AddToScheme(scheme))

	// add openshift routes scheme to controller
	utilruntime.Must(routesv1.AddToScheme(scheme))

	// add common web ui scheme to controller
	utilruntime.Must(operatorsv1alpha1.AddToScheme(scheme))

	utilruntime.Must(im.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func getFilteredCache(namespaces []string) cache.NewCacheFunc {
	commonLabels := map[string]string{
		"app.kubernetes.io/instance":   "ibm-commonui-operator",
		"app.kubernetes.io/managed-by": "ibm-commonui-operator",
	}

	commonSelector := labels.SelectorFromSet(commonLabels).String()

	//We are kind of stuck here - there isn't an enhanced multinamespace cache, but we need
	//to watch multiple configmaps with different selectors (see below comments).  So for now
	//we will not limit the cache of configmaps
	//corev1.SchemeGroupVersion.WithKind("ConfigMap"): {
	//	LabelSelector: commonSelector,
	//},
	gvkLabelsMap := map[schema.GroupVersionKind]filteredcache.Selector{
		appsv1.SchemeGroupVersion.WithKind("Deployment"): {
			LabelSelector: commonSelector,
		},
		corev1.SchemeGroupVersion.WithKind("Service"): {
			LabelSelector: commonSelector,
		},
		corev1.SchemeGroupVersion.WithKind("Secret"): {
			FieldSelector: "metadata.name==" + res.UICertSecretName,
		},
	}

	return filteredcache.MultiNamespacedFilteredCacheBuilder(gvkLabelsMap, namespaces)
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
		o.TimeEncoder = zapcore.RFC3339TimeEncoder
	}))

	printVersion()

	watchNamespace, err := getWatchNamespace()
	if err != nil {
		setupLog.Error(err, "unable to get WatchNamespace, "+
			"the manager will watch and manage resources in all namespaces")
	}

	var ctrlOpt ctrl.Options
	if strings.Contains(watchNamespace, ",") {
		// Create MultiNamespacedCache with watched namespaces if the watch namespace string contains comma
		newCache := getFilteredCache(strings.Split(watchNamespace, ","))

		ctrlOpt = ctrl.Options{
			Scheme:                 scheme,
			MetricsBindAddress:     metricsAddr,
			Port:                   9443,
			HealthProbeBindAddress: probeAddr,
			LeaderElection:         enableLeaderElection,
			LeaderElectionID:       "cf857902.ibm.com",
			NewCache:               newCache,
		}
	} else {
		// Create manager option for watching all namespaces.
		ctrlOpt = ctrl.Options{
			Scheme:                 scheme,
			MetricsBindAddress:     metricsAddr,
			Port:                   9443,
			HealthProbeBindAddress: probeAddr,
			LeaderElection:         enableLeaderElection,
			LeaderElectionID:       "cf857902.ibm.com",
			Namespace:              watchNamespace, // namespaced-scope when the value is not empty
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrlOpt)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	//Determine if this is a cncf cluster, if it is, do not watch routes
	isCncf, err := isCncf(mgr)
	if err != nil {
		log.Error(err, "Unable to determine CNCF cluster, assuming NOT CNCF - routes will be managed")
	} else {
		log.Info("Cluster type determined", "isCncf", isCncf)
	}

	// Setup Scheme for all resources
	if err := clientgoscheme.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := operatorsv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup Scheme for cert-manager
	if err := certmgr.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err := cmmeta.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if !isCncf {
		//routes Scheme
		if err := routesv1.AddToScheme(mgr.GetScheme()); err != nil {
			log.Error(err, "")
			os.Exit(1)
		}
	}

	//rbac Scheme
	if err := rbacv1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	if err = (&commonwebuicontrollers.CommonWebUIReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		IsCncf: isCncf,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CommonWebUI")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// Returns the Namespace the operator should be watching for changes
func getWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	var watchNamespaceEnvVar = "WATCH_NAMESPACE"

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}

func isCncf(mgr manager.Manager) (iscncf bool, err error) {
	//We need to determine the cluster type during startup
	//so we will use direct API calls since they are only done once

	iscncf = false
	reqLogger := log.WithValues("func", "isCncf")
	reqLogger.Info("Checking kubernetes cluster type in ibm-cpp-config")

	//Try and locate the ibm-cpp-config configmap in any of the watched namespaces
	watchNamespace, err := getWatchNamespace()
	if err != nil {
		return
	}
	nsa := strings.Split(watchNamespace, ",")
	for _, ns := range nsa {
		ibmCppConfig := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ibm-cpp-config",
				Namespace: ns,
			},
		}

		err = mgr.GetAPIReader().Get(context.TODO(), types.NamespacedName{Name: "ibm-cpp-config", Namespace: ns}, ibmCppConfig)
		if err != nil {
			log.Error(err, "Unable to load ibm-cpp-config configmap", "watched namespace", ns)
		} else {
			clusterType := ibmCppConfig.Data["kubernetes_cluster_type"]
			reqLogger.Info("Got ibm-cpp-config configmap - Kubernetes cluster type is "+clusterType, "ibm-cpp-config namespace", ns)
			if clusterType == "cncf" {
				iscncf = true
			}
			return
		}
	}

	//If we get this far, then ibm-cpp-config was not found
	err = fmt.Errorf("Unable to load the ibm-cpp-config configmap from any of the watched namespaces")

	return
}

func getSharedServicesNamespaceFromCommonService(mgr manager.Manager) (namespace string, err error) {
	reqLogger := log.WithValues("func", "getSharedServicesNamespaceFromCommonService")
	reqLogger.Info("Getting shared services namespace from common service CR")

	var operatorNamespaceEnvVar = "OPERATOR_NAMESPACE"
	operatorNamespace, found := os.LookupEnv(operatorNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("failed to get operator namespace from %s ENV var", operatorNamespaceEnvVar)
	}

	key := types.NamespacedName{Name: CommonServiceName, Namespace: operatorNamespace}

	log.Info("key", "key", key)

	gvk := schema.GroupVersionKind{
		Group:   "operator.ibm.com",
		Version: "v3",
		Kind:    "CommonService",
	}

	unstrCS := &unstructured.Unstructured{}
	unstrCS.SetGroupVersionKind(gvk)

	err = mgr.GetAPIReader().Get(context.TODO(), key, unstrCS)
	if err != nil {
		log.Error(err, "Failed to get CommonService as unstructured object")
		return
	}

	spec, ok := unstrCS.Object["spec"].(map[string]interface{})
	if !ok {
		log.Error(nil, "Failed to convert CommonService spec into map[string]interface{}")
		err = fmt.Errorf(".spec of CommonService %s in namespace %s is not a map[string]interface{}", CommonServiceName, operatorNamespace)
		return
	}
	namespace, ok = spec["servicesNamespace"].(string)
	if !ok {
		log.Error(nil, "Failed to get string servicesNamespace from CommonService spec")
		err = fmt.Errorf(".spec.servicesNamespace of CommonService %s in namespace %s is not a string", CommonServiceName, operatorNamespace)
	}
	return
}
