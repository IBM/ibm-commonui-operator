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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller_commonwebui")

var TrueVar = true
var FalseVar = false
var Seconds60 int64 = 60
var DefaultVolumeMode int32 = 420

var cpu300 = resource.NewMilliQuantity(300, resource.DecimalSI)        // 300m
var memory256 = resource.NewQuantity(256*1024*1024, resource.BinarySI) // 256Mi
var memory251 = resource.NewQuantity(251*1024*1024, resource.BinarySI) // 251Mi

const ReleaseName = "common-web-ui"
const DeploymentName = "common-web-ui"
const ServiceName = "common-web-ui"
const CommonWebUICRType = "commonwebuiservice_cr"
const ConsoleRouteName = "cp-console"

const ClusterInfoConfigmapName = "ibmcloud-cluster-info"
const PlatformAuthIdpConfigmapName = "platform-auth-idp"

const DaemonSetName = "common-web-ui"

const CertRestartLabel = "certmanager.k8s.io/time-restarted"
const NSSAnnotation = "nss.ibm.com/namespaceList"

const DefaultNamespace = "ibm-common-services"
const DefaultImageRegistry = "icr.io/cpopen/cpfs"
const DefaultImageName = "common-web-ui"
const DefaultImageTag = "1.2.1"

var DefaultStatusForCR = []string{"none"}

const Log4jsConfigMapName = "common-web-ui-log4js"

var Log4jsConfigMapData = map[string]string{
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

const CommonConfigMapName = "common-web-ui-config"
const ZenCardExtensionsConfigMapName = "common-web-ui-zen-card-extensions"
const ZenQuickNavExtensionsConfigMapName = "common-web-ui-zen-quicknav-extensions"
const ZenLeftNavExtensionsConfigMapName = "common-webui-ui-extensions"

const Log4jsVolumeName = "log4js"
const ClusterCaVolumeName = "cluster-ca"
const InternalTLSVolumeName = "internal-tls"
const IAMDataVolumeName = "iamdata"
const IAMAuthDataVolumeName = "iamadata"
const UICertVolumeName = "common-web-ui-certs"
const UICertSecretName = "common-web-ui-cert" + ""
const WebUIConfigVolumeName = "common-web-ui-config"
const ClusterInfoConfigVolumeName = "ibmcloud-cluster-info"
const PlatformAuthIdpConfigVolumeName = "platform-auth-idp"
const ZenProductInfoConfigVolumeName = "product-configmap"

var Log4jsVolume = corev1.Volume{
	Name: Log4jsVolumeName,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "common-web-ui-log4js",
			},
			Items: []corev1.KeyToPath{
				{
					Key:  "log4js.json",
					Path: "log4js.json",
				},
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var WebUIConfigVolume = corev1.Volume{
	Name: WebUIConfigVolumeName,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "common-web-ui-config",
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var ClusterInfoConfigVolume = corev1.Volume{
	Name: ClusterInfoConfigVolumeName,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "ibmcloud-cluster-info",
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var PlatformAuthIdpConfigVolume = corev1.Volume{
	Name: PlatformAuthIdpConfigVolumeName,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "platform-auth-idp",
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var ZenProductInfoConfigVolume = corev1.Volume{
	Name: ZenProductInfoConfigVolumeName,
	VolumeSource: corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "product-configmap",
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var ClusterCaVolume = corev1.Volume{
	Name: ClusterCaVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName: "cs-ca-certificate-secret",
			Items: []corev1.KeyToPath{
				{
					Key:  "tls.key",
					Path: "ca.key",
				},
				{
					Key:  "tls.crt",
					Path: "ca.crt",
				},
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}
var InternalTLSVolume = corev1.Volume{
	Name: InternalTLSVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName: "internal-tls",
			Items: []corev1.KeyToPath{
				{
					Key:  "tls.key",
					Path: "ca.key",
				},
				{
					Key:  "ca.crt",
					Path: "ca.crt",
				},
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}
var UICertVolume = corev1.Volume{
	Name: UICertVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName:  UICertSecretName,
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var IAMDataVolume = corev1.Volume{
	Name: IAMDataVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName: "platform-oidc-credentials",
			Items: []corev1.KeyToPath{
				{
					Key:  "WLP_CLIENT_SECRET",
					Path: "wlpcs",
				},
				{
					Key:  "WLP_CLIENT_ID",
					Path: "wlpcid",
				},
				{
					Key:  "OAUTH2_CLIENT_REGISTRATION_SECRET",
					Path: "oa2crs",
				},
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var IAMAuthDataVolume = corev1.Volume{
	Name: IAMAuthDataVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName: "platform-auth-idp-credentials",
			Items: []corev1.KeyToPath{
				{
					Key:  "admin_username",
					Path: "aun",
				},
			},
			Optional:    &TrueVar,
			DefaultMode: &DefaultVolumeMode,
		},
	},
}

var ArchitectureList = []string{
	"amd64",
	"ppc64le",
	"s390x",
}
var DeploymentAnnotations = map[string]string{
	"scheduler.alpha.kubernetes.io/critical-pod": "",
	"productName":   "IBM Cloud Platform Common Services",
	"productID":     "068a62892a1e4db39641342e592daa25",
	"productMetric": "FREE",
}

var CommonVolumeMounts = []corev1.VolumeMount{
	{
		Name:      Log4jsVolumeName,
		MountPath: "/etc/config",
	},
	{
		Name:      ClusterCaVolumeName,
		MountPath: "/opt/ibm/platform-header/certs",
	},
	{
		Name:      UICertVolumeName,
		MountPath: "/certs/common-web-ui",
	},
	{
		Name:      InternalTLSVolumeName,
		MountPath: "/etc/internal-tls",
	},
	{
		Name:      IAMDataVolumeName,
		MountPath: "/etc/iamdata",
	},
	{
		Name:      IAMAuthDataVolumeName,
		MountPath: "/etc/iamadata",
	},
	{
		Name:      WebUIConfigVolumeName,
		MountPath: "/etc/config/common-web-ui-config",
	},
	{
		Name:      ClusterInfoConfigVolumeName,
		MountPath: "/etc/config/ibmcloud-cluster-info",
	},
	{
		Name:      PlatformAuthIdpConfigVolumeName,
		MountPath: "/etc/config/platform-auth-idp",
	},
	{
		Name:      ZenProductInfoConfigVolumeName,
		MountPath: "/etc/config/product-configmap",
	},
}

const APIIngressName = "common-web-ui-api"

var APIIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":            "ibm-icp-management",
	"icp.management.ibm.com/secure-backends": "true",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;
        port_in_redirect off;`,
}

const CallbackIngressName = "common-web-ui-callback"

var CallbackIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":            "ibm-icp-management",
	"icp.management.ibm.com/upstream-uri":    "/auth/liberty/callback",
	"icp.management.ibm.com/secure-backends": "true",
}

const NavIngressName = "common-web-ui"

var NavIngressAnnotations = map[string]string{
	"kubernetes.io/ingress.class":            "ibm-icp-management",
	"icp.management.ibm.com/auth-type":       "access-token",
	"icp.management.ibm.com/secure-backends": "true",
	"icp.management.ibm.com/app-root":        "/common-nav?root=true",
	//nolint
	"icp.management.ibm.com/configuration-snippet": `
		add_header 'X-XSS-Protection' '1' always;`,
}

const DefaultClusterIssuer = "cs-ca-issuer"
const Certv1alpha1APIVersion = "certmanager.k8s.io/v1alpha1"
const UICertName = "common-web-ui-ca-cert"
const UICertCommonName = "common-web-ui"

// type CertificateData struct {
// 	Name      string
// 	Secret    string
// 	Common    string
// 	App       string
// 	Component string
// }

var UICertificateData = CertificateData{
	Name:      UICertName,
	Secret:    UICertSecretName,
	Common:    UICertCommonName,
	App:       "common-web-ui",
	Component: "common-web-ui",
}

const AdminHubNavConfigName = "common-web-ui-config"

// nolint
var AdminHubNavConfig = `
{
	"apiVersion": "foundation.ibm.com/v1",
	"kind": "NavConfiguration",
	"metadata": {
	  "labels": {
		"app.kubernetes.io/instance": "common-web-ui-config",
		"app.kubernetes.io/managed-by": "ibm-commonui-operator",
		"app.kubernetes.io/name": "ibm-commonui-operator",
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
		  "checker-qual, version 2.0.0, GPLv2",
		  "ocrad-bower, v1.0.0, GPL v3",
		  "Grafana, v7.5.12, AGPL",
		  "IBM Semeru Runtime Open Edition binaries, GPL v2",
		  "checker-compat-qual 2.5.5, GPL-2.0",
          "javax.annotation-api 1.3.1, GPL-2.0"
		],
		"logoUrl": "Administration",
		"logoAltText": "Administration"
	  },
	  "header": {
		"disabledItems": [
		  "createResource",
		  "catalog"
		],
		"docUrlMapping": "https://ibm.biz/cpcs_adminui",
		"logoAltText": "Administration",
		"logoHeight": "47px",
		"logoUrl": "/common-nav/graphics/ibm-cloudpack-logo.svg",
		"logoWidth": "190px"
	  },
	  "login": {
		"loginDialog": {
		  "acceptText": "Your acceptance text here",
		  "dialogText": "You must set your dialog for this environment",
		  "enable": false,
		  "headerText": "Header text here"
		},
		"logoAltText": "Cloud Pak",
		"logoHeight": "47px",
		"logoUrl": "/common-nav/api/graphics/logincloudpak.svg",
		"logoWidth": "190px"
	  },
	  "navItems": [
		{
			"id": "providers",
			"label": "Identity providers",
			"serviceId": "webui-nav",
			"url": "/common-nav/identity-access/realms"
		}
	  ]
	}
  }
`
