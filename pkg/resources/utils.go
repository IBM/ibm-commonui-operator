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
	"encoding/json"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("resource_utils")

const ReleaseName = "common-web-ui"
const Log4jsConfigMap = "common-web-ui-log4js"
const CommonConfigMap = "common-web-ui-config"
const DaemonSetName = "common-web-ui"
const ServiceName = "common-web-ui"
const APIIngress = "common-web-ui-api"
const CallbackIngress = "common-web-ui-callback"
const NavIngress = "common-web-ui"
const commonwebuiserviceCrType = "commonwebuiservice_cr"

const LegacyReleaseName = "platform-header"
const legacyheaderCrType = "legacyheader_cr"

const ChartName = "webui-nav"
const ChartVersion = "1.0.2"

var DeamonSetAnnotations = map[string]string{
	"scheduler.alpha.kubernetes.io/critical-pod": "",
	"productName":    "IBM Cloud Platform Common Services",
	"productID":      "068a62892a1e4db39641342e592daa25",
	"productVersion": "3.3.0",
	"productMetric":  "FREE",
}

var APIIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class": "ibm-icp-management",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;
        add_header Content-Security-Policy "default-src 'none'; font-src 'unsafe-inline' 'self'; script-src 'unsafe-inline' 'self' blob: cdn.segment.com fast.appcues.com; connect-src 'self' https://api.segment.io wss://api.appcues.net https://notify.bugsnag.com; img-src * data:; frame-src 'self' https://my.appcues.com; style-src 'unsafe-inline' 'self' https://fast.appcues.com; frame-ancestors 'self'";`,
}

var CallbackIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":         "ibm-icp-management",
	"icp.management.ibm.com/upstream-uri": "/auth/liberty/callback",
}

var CommonUIIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":      "ibm-icp-management",
	"icp.management.ibm.com/auth-type": "access-token",
	"icp.management.ibm.com/app-root":  "/common-nav?root=true",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;
        add_header Content-Security-Policy "default-src 'none'; font-src * 'unsafe-inline' 'self' data:; script-src 'unsafe-inline' 'self' blob: cdn.segment.com fast.appcues.com; connect-src 'self' https://api.segment.io wss://api.appcues.net https://notify.bugsnag.com; img-src * data:; frame-src 'self' https://my.appcues.com; style-src 'unsafe-inline' 'self' https://fast.appcues.com; frame-ancestors 'self'";`,
}

var CommonLegacyIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":      "ibm-icp-management",
	"icp.management.ibm.com/auth-type": "access-token",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;
        add_header Content-Security-Policy "default-src 'none'; font-src * 'unsafe-inline' 'self' data:; script-src 'unsafe-inline' 'self' blob: cdn.segment.com fast.appcues.com; connect-src 'self' https://api.segment.io wss://api.appcues.net https://notify.bugsnag.com; img-src * data:; frame-src 'self' https://my.appcues.com; style-src 'unsafe-inline' 'self' https://fast.appcues.com; frame-ancestors 'self'";`,
}

var Log4jsData = map[string]string{
	"log4js.json": `   {
		"appenders": {
		  "console": {
			"type": "console",
			"layout": {
			"type": "pattern",
			"pattern": "[%d] [%p] [webui-nav] [%c] %m"
			}
		  }
		},
		"categories": {
		  "default": { "appenders": ["console"], "level": "info" },
		  "request": { "appenders": ["console"], "level": "error" },
		  "socket.io": { "appenders": ["console"], "level": "error" },
		  "status": { "appenders": ["console"], "level": "info" },
		  "watcher": { "appenders": ["console"], "level": "debug" },
		  "service-watcher": { "appenders": ["console"], "level": "error" },
		  "session-poller": { "appenders": ["console"], "level": "error" },
		  "service-discovery": { "appenders": ["console"], "level": "info" },
		  "service-account": { "appenders": ["console"], "level": "info" },
		  "version": { "appenders": ["console"], "level": "error" },
		  "user-mgmt-client": { "appenders": ["console"], "level": "error" },
		  "oidc-client": { "appenders": ["console"], "level": "error" },
		  "server": { "appenders": ["console"], "level": "info" },
		  "auth": { "appenders": ["console"], "level": "error" },
		  "logout": { "appenders": ["console"], "level": "error" },
		  "app": { "appenders": ["console"], "level": "error" },
		  "userMgmt": { "appenders": ["console"], "level": "error" },
		  "catalog-client": { "appenders": ["console"], "level": "error" },
		  "template": { "appenders": ["console"], "level": "error" }
		}
	  }`,
}

// returns the labels associated with the resource being created
func LabelsForMetadata(deploymentName string) map[string]string {
	return map[string]string{"app.kubernetes.io/instance": deploymentName, "chart": ChartName, "version": ChartVersion,
		"heritage": "operator", "release": ReleaseName}
}

// returns the labels for selecting the resources belonging to the given metering CR name
func LabelsForSelector(deploymentName string, crType string, crName string) map[string]string {
	return map[string]string{"k8s-app": deploymentName, crType: crName}
}

// returns the labels associated with the Pod being created
func LabelsForPodMetadata(deploymentName string, crType string, crName string) map[string]string {
	podLabels := LabelsForMetadata(deploymentName)
	selectorLabels := LabelsForSelector(deploymentName, crType, crName)
	for key, value := range selectorLabels {
		podLabels[key] = value
	}
	return podLabels
}

func Log4jsConfigMapUI(instance *operatorsv1alpha1.CommonWebUI) *corev1.ConfigMap {
	reqLogger := log.WithValues("func", "log4jsConfigMap", "Name", instance.Name)
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(Log4jsConfigMap)
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Log4jsConfigMap,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Data: Log4jsData,
	}

	return configmap
}

// serviceForCommonWebUI returns a common web ui Service object
func ServiceForCommonWebUI(instance *operatorsv1alpha1.CommonWebUI) *corev1.Service {
	reqLogger := log.WithValues("func", "serviceForCommonWebUI", "instance.Name", instance.Name)
	metaLabels := LabelsForMetadata(ServiceName)
	metaLabels["kubernetes.io/cluster-service"] = "true"
	metaLabels["kubernetes.io/name"] = instance.Spec.CommonWebUIConfig.ServiceName
	metaLabels["app"] = instance.Spec.CommonWebUIConfig.ServiceName
	selectorLabels := LabelsForSelector(ServiceName, commonwebuiserviceCrType, instance.Name)

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
	return service
}

func APIIngressForCommonWebUI(instance *operatorsv1alpha1.CommonWebUI) *netv1.Ingress {
	reqLogger := log.WithValues("func", "apiIngressForCommonWebUI", "Ingress.Name", instance.Name)
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(APIIngress)
	Annotations := APIIngressAnnotations
	IngressPath := instance.Spec.CommonWebUIConfig.IngressPath
	APIIngressPath := IngressPath + "/api/"
	LogoutIngressPath := IngressPath + "/logout/"
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        APIIngress,
			Annotations: Annotations,
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
									Path: APIIngressPath,
									Backend: netv1.IngressBackend{
										ServiceName: instance.Spec.CommonWebUIConfig.ServiceName,
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
										},
									},
								},
								{
									Path: LogoutIngressPath,
									Backend: netv1.IngressBackend{
										ServiceName: instance.Spec.CommonWebUIConfig.ServiceName,
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
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
	return ingress

}

func CallbackIngressForCommonWebUI(instance *operatorsv1alpha1.CommonWebUI) *netv1.Ingress {
	reqLogger := log.WithValues("func", "callbackIngressForCommonWebUI", "Ingress.Name", instance.Name)
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(CallbackIngress)
	Annotations := CallbackIngressAnnotations
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        CallbackIngress,
			Annotations: Annotations,
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
									Path: "/auth/liberty/callback",
									Backend: netv1.IngressBackend{
										ServiceName: instance.Spec.CommonWebUIConfig.ServiceName,
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
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
	return ingress

}

func NavIngressForCommonWebUI(instance *operatorsv1alpha1.CommonWebUI) *netv1.Ingress {
	reqLogger := log.WithValues("func", "navIngressForCommonWebUI", "Ingress.Name", instance.Name)
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(NavIngress)
	Annotations := CommonUIIngressAnnotations
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        NavIngress,
			Annotations: Annotations,
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
									Path: instance.Spec.CommonWebUIConfig.IngressPath,
									Backend: netv1.IngressBackend{
										ServiceName: instance.Spec.CommonWebUIConfig.ServiceName,
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
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
	return ingress
}

func CommonConfigMapUI(instance *operatorsv1alpha1.LegacyHeader) *corev1.ConfigMap {
	reqLogger := log.WithValues("func", "commonConfigMap", "Name", instance.Name)
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(CommonConfigMap)
	data := map[string]interface{}{
		"ui-config.json": map[string]interface{}{
			"icpText": instance.Spec.LegacyConfig.LegacyLogoAltText,
			"loginDialog": map[string]interface{}{
				"enable":     false,
				"headerText": "Header text here",
				"dialogText": "You must set your dialog for this environment",
				"acceptText": "Your acceptance text here",
			},
			"login": map[string]interface{}{
				"path":   "/common-nav/api/graphics/logincloudpak.svg",
				"width":  "190px",
				"height": "47px",
			},
			"about": map[string]interface{}{
				"path": instance.Spec.LegacyConfig.LegacyLogoAltText,
				"text": instance.Spec.LegacyConfig.LegacyLogoAltText,
			},
			"header": map[string]interface{}{
				"path":   instance.Spec.LegacyConfig.LegacyLogoPath,
				"width":  instance.Spec.LegacyConfig.LegacyLogoWidth,
				"height": instance.Spec.LegacyConfig.LegacyLogoHeight,
			},
		},
	}
	jsonData, _ := json.Marshal(data["ui-config.json"])

	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CommonConfigMap,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Data: map[string]string{
			"uiconfig.json": string(jsonData),
		},
	}

	return configmap
}

// serviceForLegacyUI returns a legacy header Service object
func ServiceForLegacyUI(instance *operatorsv1alpha1.LegacyHeader) *corev1.Service {
	reqLogger := log.WithValues("func", "serviceForLegacyUI", "instance.Name", instance.Name)
	metaLabels := LabelsForMetadata(ServiceName)
	metaLabels["kubernetes.io/cluster-service"] = "true"
	metaLabels["kubernetes.io/name"] = instance.Spec.LegacyConfig.ServiceName
	metaLabels["app"] = instance.Spec.LegacyConfig.ServiceName
	selectorLabels := LabelsForSelector(ServiceName, legacyheaderCrType, instance.Name)

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
	return service
}

func IngressForLegacyUI(instance *operatorsv1alpha1.LegacyHeader) *netv1.Ingress {
	reqLogger := log.WithValues("func", "IngressForLegacyUI", "Ingress.Name", instance.Name)
	reqLogger.Info("CS??? Entry")
	metaLabels := LabelsForMetadata(NavIngress)
	Annotations := CommonLegacyIngressAnnotations
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        LegacyReleaseName,
			Annotations: Annotations,
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
									Path: instance.Spec.LegacyConfig.IngressPath,
									Backend: netv1.IngressBackend{
										ServiceName: instance.Spec.LegacyConfig.ServiceName,
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
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
	return ingress
}
