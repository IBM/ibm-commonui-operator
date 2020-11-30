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
	"strings"

	operatorsv1alpha1 "github.com/ibm/ibm-commonui-operator/pkg/apis/operators/v1alpha1"
	certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os"

	apiextv1beta "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("resource_utils")

type CertificateData struct {
	Name      string
	Secret    string
	Common    string
	App       string
	Component string
}

const ReleaseName = "common-web-ui"
const Log4jsConfigMap = "common-web-ui-log4js"
const CommonConfigMap = "common-web-ui-config"
const DaemonSetName = "common-web-ui"
const DeploymentName = "common-web-ui"
const ServiceName = "common-web-ui"
const APIIngress = "common-web-ui-api"
const CallbackIngress = "common-web-ui-callback"
const NavIngress = "common-web-ui"

const LegacyReleaseName = "platform-header"

const ChartName = "webui-nav"
const ChartVersion = "1.0.2"

var DefaultStatusForCR = []string{"none"}

//GetImageID constructs image IDs for operands: either <IMAGE_NAME>:<IMAGE_TAG> or <IMAGE_NAME>@<IMAGE_SHA>
func GetImageID(imageRegistry, imageName, defaultImageVersion, imagePostfix, envVarName string) string {
	reqLogger := log.WithValues("Func", "GetImageID")

	var imageSuffix string

	//Check if the env var exists, if yes, check whether it's a SHA or tag and use accordingly; if no, use the default image version
	imageTagOrSHA := os.Getenv(envVarName)

	if len(imageTagOrSHA) > 0 {
		//check if it is a SHA or tag and prepend appropriately
		if strings.HasPrefix(imageTagOrSHA, "sha256:") {
			reqLogger.Info("Using SHA digest value from environment variable for image " + imageName)
			imageSuffix = "@" + imageTagOrSHA
		} else {
			reqLogger.Info("Using tag value from environment variable for image " + imageName)
			imageSuffix = ":" + imageTagOrSHA
			if imagePostfix != "" {
				imageSuffix += imagePostfix
			}
		}
	} else {
		//Use default value
		reqLogger.Info("Using default tag value for image " + imageName)
		imageSuffix = ":" + defaultImageVersion
		if imagePostfix != "" {
			imageSuffix += imagePostfix
		}
	}

	imageID := imageRegistry + "/" + imageName + imageSuffix

	reqLogger.Info("imageID: " + imageID)

	return imageID
}

var DeamonSetAnnotations = map[string]string{
	"scheduler.alpha.kubernetes.io/critical-pod": "",
	"productName":   "IBM Cloud Platform Common Services",
	"productID":     "068a62892a1e4db39641342e592daa25",
	"productMetric": "FREE",
}

var DeploymentAnnotations = map[string]string{
	"scheduler.alpha.kubernetes.io/critical-pod": "",
	"productName":   "IBM Cloud Platform Common Services",
	"productID":     "068a62892a1e4db39641342e592daa25",
	"productMetric": "FREE",
}

var APIIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":            "ibm-icp-management",
	"icp.management.ibm.com/secure-backends": "true",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;
        add_header Content-Security-Policy "default-src 'none'; font-src 'unsafe-inline' 'self'; script-src 'unsafe-inline' 'self' blob: cdn.segment.com fast.appcues.com; connect-src 'self' https://api.segment.io wss://api.appcues.net https://notify.bugsnag.com; img-src * data:; frame-src 'self' https://my.appcues.com; style-src 'unsafe-inline' 'self' https://fast.appcues.com; frame-ancestors 'self'";
        port_in_redirect off;`,
}

var CallbackIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":            "ibm-icp-management",
	"icp.management.ibm.com/upstream-uri":    "/auth/liberty/callback",
	"icp.management.ibm.com/secure-backends": "true",
}

var CommonUIIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":            "ibm-icp-management",
	"icp.management.ibm.com/auth-type":       "access-token",
	"icp.management.ibm.com/secure-backends": "true",
	"icp.management.ibm.com/app-root":        "/common-nav?root=true",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;
        add_header Content-Security-Policy "default-src 'none'; font-src * 'unsafe-inline' 'self' data:; script-src 'unsafe-inline' 'self' blob: cdn.segment.com fast.appcues.com; connect-src 'self' https://api.segment.io wss://api.appcues.net https://notify.bugsnag.com; img-src * data:; frame-src 'self' https://my.appcues.com; style-src 'unsafe-inline' 'self' https://fast.appcues.com; frame-ancestors 'self' https://*.multicloud-ibm.com";`,
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

var UICertificateData = CertificateData{
	Name:      UICertName,
	Secret:    UICertSecretName,
	Common:    UICertCommonName,
	App:       "common-web-ui",
	Component: "common-web-ui",
}

//nolint
var CrTemplates = `[
	{
		"apiVersion": "foundation.ibm.com/v1",
		"kind": "NavConfiguration",
		"metadata": {
		  "labels": {
			"default": "true",
			"name": "common-web-ui-config"
		  },
		  "name": "common-web-ui-config"
		},
		"spec": {
		  "about": {
			"copyright": "© 2018, 2020 IBM. All rights reserved.",
			"licenses": [
				"yq, version 3.3.0, MIT+GPL",
				"MongoDB, version 4.0.16 Community Edition, SSPL",
			  "Ansible: © 2017 Red Hat, Inc., http://www.redhat.com; © Henry Graham (hzgraham) \u003cHenry.Graham@mail.wvu.edu\u003e",
			  "calico-bird: © 1998–2008, Martin Mares \u003cmj@ucw.cz\u003e; © 1998–2000, Pavel Machek \u003cpavel@ucw.cz\u003e; © 1998–2008, Ondrej Filip \u003cfeela@network.cz\u003e; © 2009–2013,  CZ.NIC z.s.p.o.",
			  "chrony: © Richard P. Curnow  1997-2003, GPL v2",
			  "collectd, © 2017-2018, version 5.7.2, GPL v2, \u003chttps://github.com/collectd/collectd/tree/collectd-5.7.2\u003e",
			  "crudini: © Pádraig Brady \u003cP@draigBrady.com\u003e",
			  "Galera-3: © 2007–2014 Codership Oy \u003cinfo@codership.com\u003e",
			  "glusterfs: © 2010–2013+ James Shubin \u003chttps://ttboj.wordpress.com/\u003e",
			  "haproxy: © 2000–2013  Willy Tarreau \u003cw@1wt.eu\u003e",
			  "heketi v6.0.0: © 2015 The heketi Authors, GPL v2",
			  "heketi v8.0.0: © 2015 The heketi Authors, GPL v2",
			  "heketi v9.0.0: © 2015 The heketi Authors, LGPL v3",
			  "heketi-master/apps.app.go: © 2015 The heketi Authors",
			  "heketi-master/client/api/go-client/backup.go: © 2016 The heketi Authors",
			  "heketi-master/doc/man/heketi-cli.8: © 2016 The heketi Authors",
			  "heketi-master/extras/docker/gluster/gluster-setup.sh: © 2016 Red Hat, Inc. \u003chttp://www.redhat.com\u003e",
			  "ieee-data: © 2013 Luciano Bello \u003cluciano@debian.org\u003e",
			  "javax.mail: © 2017 Oracle and/or its affiliates. All rights reserved.",
			  "keepalived: © 2001-2017 Alexandre Cassen \u003cacassen@gmail.com\u003e",
			  "libonig2: © 2006–2008 Max Kellermann \u003cmax@duempel.org\u003e; © 2014–2015 Jörg Frings-Fürst \u003cdebian@jff-webhosting.net\u003e",
			  "libtomcrypt: © 2004 Sam Hocevar \u003csam@hocevar.net\u003e, GPL v2",
			  "mariadb-common: © 2018 MariaDB. All rights reserved.",
			  "mariaDB: © 2018 MariaDB. All rights reserved. \u003chttps://mariadb.com/\u003e",
			  "mariadb-server: © 2018 MariaDB. All rights reserved.",
			  "minitar: © 2004 Mauricio Julio Fernandez Pradier and Austin Ziegler",
			  "MongoDB: © 2007 Free Software Foundation, Inc. \u003chttp://fsf.org/\u003e",
			  "nvmi-cli: © 1989, 1991 Free Software Foundation, Inc., GPL v2",
			  "OpenJDK: © 2018 Oracle Corporation and/or its affiliates",
			  "openshift-mariadb-galera: © 2007 Free Software Foundation, Inc. \u003chttp://fsf.org/\u003e",
			  "percona-xtrabackup: © 2006–2018 Percona LLC.",
			  "pwgen: © Christian Thöing \u003cc.thoeing@web.de\u003e",
			  "rdoc: © 2001–2003 Dave Thomas, The Pragmatic Programmers",
			  "readline: © Chet Ramey \u003cchet.ramey@case.edu\u003e",
			  "John the Ripper password cracker: © 1996–2013 by Solar Designer \u003csolar@openwall.com\u003e",
			  "spdx-exceptions: © 2018 SPDX Workgroup a Linux Foundation Project. All rights reserved.",
			  "socat: © 2001–2010 Gerhard Rieger",
			  "sshpass: © 2006, 2008 Lingnu Open Source Consulting Ltd. \u003chttp://www.lingnu.com\u003e",
			  "timelimit: © 2001, 2007 - 2010  Peter Pentchev, GPL v2",
			  "ua-parser-js: © 2012-2018 Faisal Salman \u003cf@faisalman.com\u003e, GPL v2",
			  "ubuntu-cloud-keyring: © 2010 Michael Vogt \u003cmichael.vogt@canonical.com\u003e",
			  "unboundid-ldapsdk: © 2015 UnboundID. The LDAP SDK for Java is developed by UnboundID. \u003cinfo@unboundid.com\u003e",
			  "xmpp4r: © Lucas Nussbaum \u003clucas@lucas-nussbaum.net\u003e, Stephan Maka \u003cstephan@spaceboyz.net\u003e, and others.",
			  "module-assistant: © 2003-2008 Eduard Bloch \u003cblade@debian.org\u003e, version 0.11.8, GPL v2; © 2009 Cyril Brulebois \u003ckibi@debian.org\u003e, version 0.11.8, GPL v2; © 2013-2018 Andreas Beckmann \u003canbe@debian.org\u003e, version 0.11.8, GPL v2",
			  "module-init-tools: © 2011 ProFUSION embedded systems, version 22, GPL v2",
			  "thin: © 2017 Marc-Andre Cournoyer \u003cmacournoyer@gmail.com\u003e, version 1.7.2, GPL v2",
			  "gosu, © 1999-2014, version 1.1, GPL v3",
			  "mercurial (Python), © 2006-2018 ,version v4.5.3, GPL v2",
			  "garden-runc, © 2015-Present CloudFoundry.org Foundation, Inc. All Rights Reserved, version 1.17.0, GPLv2",
			  "libtomcrypt0, © 2003-2007 Tom St Denis \u003ctomstdenis@gmail.com\u003e, version 1.17-7, GPLv2",
			  "console-setup-min, © 1999,2000,2001,2002,2003,2006,2007,2008,2009,2010,2011 Anton Zinoviev, \u003canton@lml.bas.bg\u003e,version 1.108, GPLv2",
			  "dracut, © 2009 Harald Hoyer \u003charald@redhat.com\u003e, version 044+3-3, GPLv2",
			  "dracut-core, © 2009 Harald Hoyer \u003charald@redhat.com\u003e, version 044+3-3, GPLv2",
			  "g++, version 5.4.0-6ubuntu, GPL v2",
			  "libstdc++6, version 5.4.0-6ubuntu, GPL v3",
			  "libstdc++-5-dev, version 5.4.0-6ubuntu, GPL v3",
			  "docker-engine-selinux, version 3b5fac4, GPLv2",
			  "unorm, version 1.5.0, GPL v2",
			  "psmisc, version 22.20, GPL v2",
			  "lvm2-devel, version 2.0.2, GPL v2",
			  "nfs-utils, version 1.3, GPL v2",
			  "popt-static, version 1.13, GPL v2",
			  "sysvinit-tools, version 2.88, GPL v2",
			  "stunnel, version 5.53, GPL v2",
			  "stunnel, version 5.39, GPL v2",
			  "LVM2, version 2.02.180-10.el7_6.2, GPL v2",
			  "sysdig, version 2c43237, GPL",
			  "chisels, version 9722dbc, GPL",
			  "MongoDB, version 4.0.12, SSPL",
			  "ffi (Ruby Gem), 1.11.1, GPL",
			  "inotify-tools, v3.14, GPL v2",
			  "logrotate, v3.8.6, GPL v2",
			  "checker-qual, version 2.0.0, GPLv2"
			],
			"logoAltText": "IBM Cloud Pak Administration Hub",
			"logoUrl": "IBM Cloud Pak Administration Hub"
		  },
		  "header": {
			"logoAltText": "IBM Cloud Pak Administration Hub",
			"logoHeight": "47px",
			"logoUrl": "/common-nav/graphics/ibm-cloudpack-logo.svg",
			"logoWidth": "190px",
			"disabledItems":["createResource","catalog"],
			"docUrlMapping": "http://ibm.biz/cpcs_adminui"
		  },
		  "login": {
			"loginDialog": {
			  "acceptText": "Your acceptance text here",
			  "dialogText": "You must set your dialog for this environment",
			  "enable": false,
			  "headerText": "Header text here"
			},
			"logoAltText": "Cloud Pak Administration Hub",
			"logoHeight": "47px",
			"logoUrl": "/common-nav/api/graphics/logincloudpak.svg",
			"logoWidth": "190px"
		  },
		  "navItems": [
			{
			  "id": "home",
			  "label": "Home",
			  "url": "/common-nav/dashboard",
			  "iconUrl": "/common-nav/graphics/home.svg",
				"isAuthorized": [
					"ClusterAdministrator"
			  ]
			},
			{
			  "id": "id-access",
			  "label": "Identity and Access",
			  "serviceId": "webui-nav",
			  "url": "/common-nav/identity-access",
			  "iconUrl": "/common-nav/graphics/password.svg"
			},
			{
			  "detectionServiceName": true,
			  "id": "licensing",
			  "label": "Licensing",
			  "namespace": "ibm-common-services",
			  "serviceId": "ibm-license-service-reporter",
			  "serviceName": "ibm-license-service-reporter",
			  "url": "/license-service-reporter",
			  "iconUrl": "/common-nav/graphics/identification.svg"
			},
			{
			  "detectionServiceName": true,
			  "id": "metering",
			  "label": "Metering",
			  "namespace": "ibm-common-services",
			  "serviceId": "metering-ui",
			  "serviceName": "metering-ui",
			  "url": "/metering/dashboard?ace_config={ 'showClusterData': false }\u0026dashboard=cpi.icp.main",
			  "iconUrl": "/common-nav/graphics/meter--alt.svg"
			},
			{
			  "detectionServiceName": true,
			  "id": "monitoring",
			  "isAuthorized": [
				"Administrator",
				"ClusterAdministrator",
				"Operator"
			  ],
			  "label": "Monitoring",
			  "namespace": "ibm-common-services",
			  "serviceId": "monitoring-ui",
			  "serviceName": "ibm-monitoring-grafana",
			  "target": "_blank",
			  "url": "/grafana",
			  "iconUrl": "/common-nav/graphics/activity.svg"
			},
			{
			  "detectionServiceName": true,
			  "id": "logging",
			  "label": "Logging",
			  "namespace": "ibm-common-services",
			  "serviceId": "kibana",
			  "serviceName": "kibana",
			  "target": "_blank",
			  "url": "/kibana",
			  "iconUrl": "/common-nav/graphics/catalog.svg"
			}
		  ]
		}
	  },
	  {
		"apiVersion": "foundation.ibm.com/v1",
		"kind": "NavConfiguration",
		"metadata": {
		  "labels": {
			"name": "icp4i"
		  },
		  "name": "icp4i"
		},
		"spec": {
		  "header": {
			"disabledItems": [
			  "catalog",
			  "createResource",
			  "bookmark"
			],
			"logoAltText": "Cloud Pak for Integration",
			"logoUrl": "/common-nav/graphics/ibm-cloudpak-integration.svg"
		  },
		  "navItems": [
			{
			  "detectionServiceName": "true",
			  "id": "metering",
			  "label": "Metering",
			  "namespace": "ibm-common-services",
			  "serviceId": "metering-ui",
			  "serviceName": "metering-ui",
			  "url": "/metering/dashboard?ace_config={ 'showClusterData': false }\u0026dashboard=cpi.icp.main\u0026useNav=icp4i"
			},
			{
			  "detectionServiceName": "true",
			  "id": "monitoring",
			  "isAuthorized": [
				"Administrator",
				"ClusterAdministrator",
				"Operator"
			  ],
			  "label": "Monitoring",
			  "namespace": "ibm-common-services",
			  "serviceId": "monitoring-ui",
			  "serviceName": "ibm-monitoring-grafana",
			  "target": "_blank",
			  "url": "/grafana"
			},
			{
			  "detectionServiceName": true,
			  "id": "releases",
			  "label": "Helm Releases",
			  "namespace": "ibm-common-services",
			  "serviceId": "catalog-ui",
			  "serviceName": "catalog-ui",
			  "url": "/catalog/instances?useNav=icp4i"
			},
			{
			  "detectionServiceName": true,
			  "id": "repos",
			  "label": "Helm Repositories",
			  "namespace": "ibm-common-services",
			  "serviceId": "catalog-ui",
			  "serviceName": "catalog-ui",
			  "url": "/catalog/repositories?useNav=icp4i"
			},
			{
			  "id": "id-access",
			  "label": "Identity and Access",
			  "serviceId": "webui-nav",
			  "url": "/common-nav/identity-access?useNav=icp4i"
			},
			{
			  "detectionServiceName": "true",
			  "id": "logging",
			  "label": "Logging",
			  "namespace": "ibm-common-services",
			  "serviceId": "kibana",
			  "serviceName": "kibana",
			  "target": "_blank",
			  "url": "/kibana"
			},
			{
			  "detectionServiceName": true,
			  "id": "licensing",
			  "label": "Licensing",
			  "namespace": "ibm-common-services",
			  "serviceId": "ibm-license-service-reporter",
			  "serviceName": "ibm-license-service-reporter",
			  "url": "/license-service-reporter"
			}
		  ]
		}
	  }
	]`

// returns the labels associated with the resource being created
func LabelsForMetadata(deploymentName string) map[string]string {
	return map[string]string{"app.kubernetes.io/instance": "ibm-commonui-operator",
		"app.kubernetes.io/name": deploymentName, "app.kubernetes.io/managed-by": "ibm-commonui-operator", "intent": "projected"}
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

func BuildCertificate(instanceNamespace, instanceClusterIssuer string, certData CertificateData) *certmgr.Certificate {
	reqLogger := log.WithValues("func", "BuildCertificate")

	metaLabels := labelsForCertificateMeta(certData.App, certData.Component)
	var clusterIssuer string
	if instanceClusterIssuer != "" {
		reqLogger.Info("clusterIssuer=" + instanceClusterIssuer)
		clusterIssuer = instanceClusterIssuer
	} else {
		reqLogger.Info("clusterIssuer is blank, default=" + DefaultClusterIssuer)
		clusterIssuer = DefaultClusterIssuer
	}

	certificate := &certmgr.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      certData.Name,
			Labels:    metaLabels,
			Namespace: instanceNamespace,
		},
		Spec: certmgr.CertificateSpec{
			CommonName: certData.Common,
			SecretName: certData.Secret,
			IsCA:       false,
			DNSNames: []string{
				certData.Common,
				certData.Common + "." + instanceNamespace,
				certData.Common + "." + instanceNamespace + ".svc.cluster.local",
			},
			Organization: []string{"IBM"},
			IssuerRef: certmgr.ObjectReference{
				Name: clusterIssuer,
				Kind: certmgr.IssuerKind,
			},
		},
	}
	return certificate
}

func labelsForCertificateMeta(appName, componentName string) map[string]string {
	return map[string]string{
		"app":                          appName,
		"component":                    componentName,
		"release":                      ReleaseName,
		"app.kubernetes.io/instance":   "ibm-commonui-operator",
		"app.kubernetes.io/managed-by": "ibm-commonui-operator",
		"app.kubernetes.io/name":       UICertName,
	}
}

// GetPodNames returns the pod names of the array of pods passed in
func GetPodNames(pods []corev1.Pod) []string {
	reqLogger := log.WithValues("func", "GetPodNames")
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
		reqLogger.Info("pod name=" + pod.Name)
	}
	return podNames
}

// GetNavConfigContent returns all nav config crd content
func GetNavConfigContent() map[string]apiextv1beta.JSONSchemaProps {
	return map[string]apiextv1beta.JSONSchemaProps{
		"logoutRedirects": apiextv1beta.JSONSchemaProps{
			Description: `A list a URLs we make requests to logout the users of all applications within the cloudpack.`,
			Type:        "array",
			Items: &apiextv1beta.JSONSchemaPropsOrArray{
				Schema: &apiextv1beta.JSONSchemaProps{
					Type: "string",
				},
			},
		},

		"header": apiextv1beta.JSONSchemaProps{
			Type:        "object",
			Description: "Customized common web ui header items",
			Properties: map[string]apiextv1beta.JSONSchemaProps{
				"disabledItems": apiextv1beta.JSONSchemaProps{
					Type: "array",
					// nolint
					Description: "An array of header items that should be disabled when running within this CR context. Valid values are 'catalog', 'createResource', 'bookmark'",
					Items: &apiextv1beta.JSONSchemaPropsOrArray{
						Schema: &apiextv1beta.JSONSchemaProps{
							Type: "string",
						},
					},
				},
				"detectHeaderItems": apiextv1beta.JSONSchemaProps{
					Type: "array",
					// nolint
					Description: "An object that maps header items to service detection values, such as service name, label selector, and namespace. The only supported header item is 'search'.",
					AdditionalProperties: &apiextv1beta.JSONSchemaPropsOrBool{
						Schema: &apiextv1beta.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextv1beta.JSONSchemaProps{
								"detectionNamespace": apiextv1beta.JSONSchemaProps{
									Type: "string",
								},
								"detectionServiceName": apiextv1beta.JSONSchemaProps{
									Type: "string",
								},
								"detectionLabelSelector": apiextv1beta.JSONSchemaProps{
									Type: "string",
								},
								"isAuthorized": apiextv1beta.JSONSchemaProps{
									Type: "array",
									Items: &apiextv1beta.JSONSchemaPropsOrArray{
										Schema: &apiextv1beta.JSONSchemaProps{
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
				"logoUrl": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "The URL that provides the login page logo. Must be an unprotected URL.",
				},
				"logoWidth": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Width of the logo for the login page in pixels",
				},
				"logoHeight": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Height of the logo for the login page in pixels",
				},
				"docUrlMapping": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "URL of the Knowledge center page for the cloud pak",
				},
				"supportUrl": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "URL of the Support page for the cloud pak",
				},
				"gettingStartedUrl": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "URL of the Getting started page for the cloud pak",
				},
			},
		},
		"about": apiextv1beta.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextv1beta.JSONSchemaProps{
				"logoUrl": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "URL of the Logo on the About page for the cloud pak",
				},
				"licenses": apiextv1beta.JSONSchemaProps{
					Type:        "array",
					Description: "List of licenses we ship with the cloud pak",
					Items: &apiextv1beta.JSONSchemaPropsOrArray{
						Schema: &apiextv1beta.JSONSchemaProps{
							Type: "string",
						},
					},
				},
				"copyright": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Copyright string for the cloud pak",
				},
				"version": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Version of the cloud pak",
				},
				"edition": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Edition of the cloud pak",
				},
			},
		},
		"login": apiextv1beta.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextv1beta.JSONSchemaProps{
				"logoUrl": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "URL of the Logo on the About page for the cloud pak",
				},
				"logoAltText": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Alternate text of the shared header logo for cloud pak",
				},
				"loginDialog": apiextv1beta.JSONSchemaProps{
					Type:        "object",
					Description: "FISMA dialog contents can be modified here",
					Properties: map[string]apiextv1beta.JSONSchemaProps{
						"enabled": apiextv1beta.JSONSchemaProps{
							Type:        "boolean",
							Description: "This value is used to enable/disable the user acceptance dialog on the login page",
						},
						"dialogHeaderText": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Text that will display as the title of the user acceptance dialog on the login page",
						},
						"dialogText": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Text that will display as the content of the user acceptance dialog on the login page",
						},
						"acceptText": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Text that will display as the accept button text",
						},
					},
				},
				"logoWidth": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Width of the logo for the login page in pixels",
				},
				"logoHeight": apiextv1beta.JSONSchemaProps{
					Type:        "string",
					Description: "Height of the logo for the login page in pixels",
				},
			},
		},

		"navItems": apiextv1beta.JSONSchemaProps{
			Description: "Navigation items for the left hand nav within common ui header for the cloud pak",
			Type:        "array",
			Items: &apiextv1beta.JSONSchemaPropsOrArray{

				Schema: &apiextv1beta.JSONSchemaProps{

					Type: "object",
					Properties: map[string]apiextv1beta.JSONSchemaProps{
						"id": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "ID of the nav item, must be unique",
						},
						"label": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Displayed label of the nav item",
						},
						"url": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "URL of the nav item. It can either but an FQDN or a relative path based on the ingress of the cluster",
						},
						"target": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "name of the tab or _blank where the navigation item will launch within the window",
						},
						"iconUrl": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "URL of the icon that will display for the top level parents.",
						},
						"parentId": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "ID of the parent this child item will be nested under",
						},
						"namespace": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Namespace where the microservice associated with this item is running. Used with service detection",
						},
						"serviceName": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Name of the service running in the namespace above tied to the deployment/daemonset. Used for service detection",
						},
						"serviceId": apiextv1beta.JSONSchemaProps{
							Type: "string",
							// nolint
							Description: "Must be unique from a different microservice link. But the service id should remain the same for all links running on the same microservice for rendering purposes.",
						},
						"detectionServiceName": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "Informs the shared web console detection service to use the serviceName for auto discovery. Value should be true or false string",
						},
						"detectionLabelSelector": apiextv1beta.JSONSchemaProps{
							Type:        "string",
							Description: "The label selector for the microservice for detection.",
						},
						"isAuthorized": apiextv1beta.JSONSchemaProps{
							Type:        "array",
							Description: "The label selector for the microservice for detection.",
							Items: &apiextv1beta.JSONSchemaPropsOrArray{
								Schema: &apiextv1beta.JSONSchemaProps{
									Type: "string",
								},
							},
						},
					},
				},
			},
		}, // navitems
	}
}

// returns the service account name or default if it is not set in the environment
func GetServiceAccountName() string {

	sa := "ibm-commonui-operator"
	return sa
}
