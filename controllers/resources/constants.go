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

var cpu300 = resource.NewMilliQuantity(300, resource.DecimalSI)        // 300m
var memory256 = resource.NewQuantity(256*1024*1024, resource.BinarySI) // 256Mi

const ReleaseName = "common-web-ui"
const DeploymentName = "common-web-ui"
const ServiceName = "common-web-ui"
const CommonWebUICRType = "commonwebuiservice_cr"
const ConsoleRouteName = "cp-console"

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

const AdminHubOnZenConfigMapName = "adminhub-on-zen-cm"
const CommonConfigMapName = "common-web-ui-config"
const ZenCardExtensionsConfigMapName = "common-web-ui-zen-card-extensions"
const ZenCardExtensionsConfigMapNameCncf = "common-web-ui-zen-card-extension-cncf"
const ZenQuickNavExtensionsConfigMapName = "common-web-ui-zen-quicknav-extensions"
const ZenWalkmeExtensionsConfigMapName = "common-web-ui-zen-walkme-extensions"

const ZenDeploymentName = "zen-core"
const ZenProductConfigMapName = "product-configmap"

var ZenPcmMap = map[string]string{
	"CLOUD_PAK_TYPE":           "admin",
	"CLOUD_PAK_URL":            "https://common-web-ui:3000/common-nav/zen/meta",
	"CLOUD_PAK_AUTH_URL":       "https://common-web-ui:3000/common-nav/zen/meta",
	"IBM_PRODUCT_NAME":         "IBM Cloud Pak | Administration",
	"IBM_DEFAULT_PRODUCT_NAME": "IBM Cloud Pak | Administration",
}

var ZenNginxConfig = `
	location /common-nav {
		access_by_lua_file /nginx_data/checkjwt.lua;
		set_by_lua $nsdomain 'return os.getenv("NS_DOMAIN")';
		proxy_set_header Host $host;
		proxy_set_header zen-namespace-domain $nsdomain;
		proxy_pass https://common-web-ui:3000;
		proxy_read_timeout 10m;
	}
`

var ZenQuickNavExtensions = `
[
      {
        "extension_point_id": "homepage_quick_navigation",
        "extension_name": "homepage_quick_navigation_id_providers",
        "display_name": "{{ .global_zen_homepage_nav_id_providers }}",
        "order_hint": 100,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "extension_type": "ootb",
          "reference": {
            "nav_item": "nav-id-providers"
          }
        },
        "details": {
          "label": "{{ .global_adminhub_id_providers }}",
          "nav_link": "/common-nav/zen/idproviders"
        }
      }
]
`

var ZenCardExtensions = `
[
	  {
        "extension_point_id": "left_menu_item",
        "extension_name": "nav-id-providers",
        "display_name": "{{ .global_adminhub_id_providers }}",
        "order_hint": 600,
        "match_permissions": "administrator",
        "meta": {},
        "details": {
			"parent_folder": "dap-header-administer",
			"href": "/common-nav/zen/idproviders"
        }
      },
      {
        "extension_point_id": "homepage_resource",
        "extension_name": "homepage_resource_documentation",
        "display_name": "{{ .global_zen_homepage_nav_documentation }}",
        "order_hint": 100,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {},
        "details": {
          "label": "{{ .global_adminhub_documentation }}",
          "nav_link": "https://ibm.biz/cpcs_adminui",
          "carbon_icon": "Document16"
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_memory_usage",
        "display_name": "{{ .global_zen_homepage_card_memory_usage }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 0,
            "row": 0
          }
        },
        "details": {
          "service_defined_id": "homepage_card_memory_usage",
          "title": "{{ .global_adminhub_memory_usage }}",
          "description": "{{ .global_zen_homepage_memory_usage_description }}",
          "drilldown_url": "",
          "template_type": "donut",
          "data_url": "/common-nav/zen/api/v1/memory_usage"
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_cluster_inventory",
        "display_name": "{{ .global_zen_homepage_card_inventory_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 1,
            "row": 0
          }
        },
        "details": {
		  "service_defined_id": "homepage_card_cluster_inventory",
		  "title": "{{ .global_adminhub_cluster_inventory }}",
          "description": "{{ .global_zen_homepage_card_inventory_description }}",
          "drilldown_url": "",
          "window_open_target": "ah_cluster_inventory",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/inventory"
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_deployments",
        "display_name": "{{ .global_zen_homepage_card_deployments_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 2,
            "row": 0
          }
        },
        "details": {
		  "service_defined_id": "homepage_card_deployments",
		  "title": "{{ .global_adminhub_deployments }}",
          "description": "{{ .global_zen_homepage_card_deployments_description }}",
          "drilldown_url": "",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/deployments",
          "empty_state": {
            "main_text": "{{ .global_adminhub_deployments_empty_main_text }}",
            "sub_text": "{{ .global_adminhub_deployments_empty_sub_text }}"
          }
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_monitoring_trends",
        "display_name": "{{ .global_zen_homepage_card_vulnerabilities }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 0,
            "row": 1
          }
        },
        "details": {
          "service_defined_id": "homepage_card_monitoring_trends",
          "title": "{{ .global_adminhub_monitoring_trends }}",
          "description": "{{ .global_zen_homepage_vulnerabilities_description }}",
          "drilldown_url": "",
          "template_type": "big_number",
          "data_url": "/common-nav/zen/api/v1/trends",
          "empty_state": {
            "main_text": "{{ .global_zen_homepage_card_vulnerabilities_empty_state_main_text }}",
            "sub_text": "{{ .global_zen_homepage_card_vulnerabilities_empty_state_sub_text }}"
          }
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_system_utility_status",
        "display_name": "{{ .global_common_core_homepage_card_system_utility_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 1,
            "row": 1
          }
        },
        "details": {
          "service_defined_id": "homepage_card_system_utility_status",
          "title": "{{ .global_adminhub_system_utility_status }}",
          "description": "{{ .global_common_core_homepage_card_recent_projects_description }}",
          "drilldown_url": "",
          "template_type": "condensed_list",
          "data_url": "/common-nav/zen/api/v1/system_utility_status",
          "empty_state": {
            "main_text": "{{ .global_adminhub_system_utility_status_empty_main_text }}",
            "sub_text": "{{ .global_adminhub_system_utility_status_empty_sub_text }}",
            "button_text": "",
            "button_url": ""
          }
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_workload_summary",
        "display_name": "{{ .global_zen_homepage_card_workload_summary }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 2,
            "row": 1
          }
        },
        "details": {
          "service_defined_id": "homepage_card_workload_summary",
          "title": "{{ .global_adminhub_workload_summary }}",
          "description": "{{ .global_zen_homepage_memory_workload_summary }}",
          "drilldown_url": "",
          "template_type": "multi_donut",
          "data_url": "/common-nav/zen/api/v1/workload-summary"
        }
      },
      {
       "extension_point_id": "homepage_card",
       "extension_name": "homepage_card_events",
       "display_name": "{{ .global_zen_homepage_card_events_name }}",
       "order_hint": 0,
       "match_permissions": "administrator",
       "match_instance_id": "",
       "match_instance_role": "",
       "meta": {
          "preferences": {
            "column": 0,
            "row": 2
          }
       },
       "details": {
         "service_defined_id": "homepage_card_events",
         "title": "{{ .global_adminhub_system_events }}",
         "description": "{{ .global_zen_homepage_card_events_description }}",
         "drilldown_url": "",
         "template_type": "text_list",
         "data_url": "/common-nav/zen/api/v1/events",
         "empty_state": {
			"main_text": "{{ .global_adminhub_system_events_empty_main_text }}",
			"sub_text": "{{ .global_adminhub_system_events_empty_sub_text }}"
         }
       }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_identity_and_users_access",
        "display_name": "{{ .global_zen_homepage_card_requests_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 1,
            "row": 2
          }
        },
        "details": {
          "service_defined_id": "homepage_card_identity_and_users_access",
          "title": "{{ .global_adminhub_identity_and_users_access }}",
          "description": "{{ .global_zen_homepage_card_requests_description }}",
          "drilldown_url": "",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/users",
          "empty_state": {
            "main_text": "{{ .global_zen_homepage_card_requests_details_empty_state_main_text }}",
            "sub_text": "{{ .global_zen_homepage_card_requests_details_empty_state_sub_text }}"
          }
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_license_products",
        "display_name": "{{ .global_zen_homepage_card_requests_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 2,
            "row": 2
          }
        },
        "details": {
          "service_defined_id": "homepage_card_license_products",
          "title": "{{ .global_adminhub_license_products }}",
          "description": "{{ .global_zen_homepage_card_requests_description }}",
          "drilldown_url": "",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/license_products",
          "empty_state": {
            "main_text": "{{ .global_zen_homepage_card_requests_details_empty_state_main_text }}",
            "sub_text": "{{ .global_zen_homepage_card_requests_details_empty_state_sub_text }}"
          }
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_diagnostics",
        "display_name": "{{ .global_common_core_homepage_card_diagnostics_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 1,
            "row": 3
          }
        },
        "details": {
          "service_defined_id": "homepage_card_diagnostics",
          "title": "Diagnostics",
          "description": "{{ .global_common_core_homepage_card_diagnostics_description }}",
          "template_type": "iframe",
          "source_url": "/common-nav/zen/api/v1/diagnostics",
          "refresh_rate": 10
        }
      }
    ]
`

var ZenCardExtensionsCncf = `
[
	  {
        "extension_point_id": "left_menu_item",
        "extension_name": "nav-id-providers",
        "display_name": "{{ .global_adminhub_id_providers }}",
        "order_hint": 600,
        "match_permissions": "administrator",
        "meta": {},
        "details": {
			"parent_folder": "dap-header-administer",
			"href": "/common-nav/zen/idproviders"
        }
      },
      {
        "extension_point_id": "homepage_resource",
        "extension_name": "homepage_resource_documentation",
        "display_name": "{{ .global_zen_homepage_nav_documentation }}",
        "order_hint": 100,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {},
        "details": {
          "label": "{{ .global_adminhub_documentation }}",
          "nav_link": "https://ibm.biz/cpcs_adminui",
          "carbon_icon": "Document16"
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_cluster_inventory",
        "display_name": "{{ .global_zen_homepage_card_inventory_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 1,
            "row": 0
          }
        },
        "details": {
		  "service_defined_id": "homepage_card_cluster_inventory",
		  "title": "{{ .global_adminhub_cluster_inventory }}",
          "description": "{{ .global_zen_homepage_card_inventory_description }}",
          "drilldown_url": "",
          "window_open_target": "ah_cluster_inventory",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/inventory"
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_deployments",
        "display_name": "{{ .global_zen_homepage_card_deployments_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 2,
            "row": 0
          }
        },
        "details": {
		  "service_defined_id": "homepage_card_deployments",
		  "title": "{{ .global_adminhub_deployments }}",
          "description": "{{ .global_zen_homepage_card_deployments_description }}",
          "drilldown_url": "",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/deployments",
          "empty_state": {
            "main_text": "{{ .global_adminhub_deployments_empty_main_text }}",
            "sub_text": "{{ .global_adminhub_deployments_empty_sub_text }}"
          }
        }
      },
      {
       "extension_point_id": "homepage_card",
       "extension_name": "homepage_card_events",
       "display_name": "{{ .global_zen_homepage_card_events_name }}",
       "order_hint": 0,
       "match_permissions": "administrator",
       "match_instance_id": "",
       "match_instance_role": "",
       "meta": {
          "preferences": {
            "column": 0,
            "row": 2
          }
       },
       "details": {
         "service_defined_id": "homepage_card_events",
         "title": "{{ .global_adminhub_system_events }}",
         "description": "{{ .global_zen_homepage_card_events_description }}",
         "drilldown_url": "",
         "template_type": "text_list",
         "data_url": "/common-nav/zen/api/v1/events",
         "empty_state": {
			"main_text": "{{ .global_adminhub_system_events_empty_main_text }}",
			"sub_text": "{{ .global_adminhub_system_events_empty_sub_text }}"
         }
       }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_identity_and_users_access",
        "display_name": "{{ .global_zen_homepage_card_requests_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 1,
            "row": 2
          }
        },
        "details": {
          "service_defined_id": "homepage_card_identity_and_users_access",
          "title": "{{ .global_adminhub_identity_and_users_access }}",
          "description": "{{ .global_zen_homepage_card_requests_description }}",
          "drilldown_url": "",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/users",
          "empty_state": {
            "main_text": "{{ .global_zen_homepage_card_requests_details_empty_state_main_text }}",
            "sub_text": "{{ .global_zen_homepage_card_requests_details_empty_state_sub_text }}"
          }
        }
      },
      {
        "extension_point_id": "homepage_card",
        "extension_name": "homepage_card_license_products",
        "display_name": "{{ .global_zen_homepage_card_requests_name }}",
        "order_hint": 0,
        "match_permissions": "administrator",
        "match_instance_id": "",
        "match_instance_role": "",
        "meta": {
          "preferences": {
            "column": 2,
            "row": 2
          }
        },
        "details": {
          "service_defined_id": "homepage_card_license_products",
          "title": "{{ .global_adminhub_license_products }}",
          "description": "{{ .global_zen_homepage_card_requests_description }}",
          "drilldown_url": "",
          "template_type": "number_list",
          "data_url": "/common-nav/zen/api/v1/license_products",
          "empty_state": {
            "main_text": "{{ .global_zen_homepage_card_requests_details_empty_state_main_text }}",
            "sub_text": "{{ .global_zen_homepage_card_requests_details_empty_state_sub_text }}"
          }
        }
      }
    ]
`

var ZenWalkmeExtensions = `
[
	{
		"extension_point_id":"generic_preferences",
		"extension_name":"guided_tours",
		"display_name":"Guided tours",
		"description":"",
		"match_permissions":"administrator",
		"meta":null,
		"details":{
		  "lite_tours_src": "/common-nav/walkme/walkme_760e1a0cad93453f8cc129ce436f336e_https.js"
		},
		"status":"enabled"
	},
	{
	  "extension_point_id": "zen_platform_customization",
	  "extension_name": "mypak_customization_tours",
	  "order_hint": 300,
	  "details": {
		"title": "{{.global_zen_platform_customization_tours_title}}",
		"description": "{{.global_zen_platform_customization_tours_description}}",
		"icon": "Crossroads20",
		"icon_alt": "{{.global_zen_platform_customization_tours_title}}",
		"nav_url": "/zen/#/guidedToursCustomization"
	  }
	}
]
`

const ZenLeftNavExtensionsConfigMapName = "common-webui-ui-extensions"

var ZenLeftNavExtensionsConfigMapData = `
[
	{
		"extension_point_id": "left_menu_item",
		"extension_name": "dap-admin-hub",
		"display_name": "Administration Hub",
		"order_hint": 100,
		"match_permissions": "administrator",
		"meta": {},
		"details": {
			"parent_folder": "dap-header-administer",
			"href": "https://%s/common-nav/dashboard",
			"target": "adminhub"
		}
	}
]`

const Log4jsVolumeName = "log4js"
const ClusterCaVolumeName = "cluster-ca"
const InternalTLSVolumeName = "internal-tls"
const UICertVolumeName = "common-web-ui-certs"
const UICertSecretName = "common-web-ui-cert" + ""

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
			Optional: &TrueVar,
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
			Optional: &TrueVar,
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
			Optional: &TrueVar,
		},
	},
}
var UICertVolume = corev1.Volume{
	Name: UICertVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName: UICertSecretName,
			Optional:   &TrueVar,
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

var ConsoleLinkTemplate = `{
	"apiVersion": "console.openshift.io/v1",
	"kind": "ConsoleLink",
	"metadata": {
		"name": "admin-hub"
	},
	"spec": {
		"applicationMenu": {
			"imageURL": "https://raw.githubusercontent.com/carbon-design-system/carbon/main/packages/icons/src/svg/32/cloud.svg",
			"section": "IBM Cloud Paks"
		},
		"href": "https://%s/common-nav/dashboard",
		"location": "ApplicationMenu",
		"text": "Administration"
	}
}`

var ConsoleLinkTemplate2 = `{
	"apiVersion": "console.openshift.io/v1",
	"kind": "ConsoleLink",
	"metadata": {
		"name": "admin-hub-zen"
	},
	"spec": {
		"applicationMenu": {
			"imageURL": "https://raw.githubusercontent.com/carbon-design-system/carbon/main/packages/icons/src/svg/32/cloud.svg",
			"section": "IBM Cloud Paks"
		},
		"href": "https://%s/common-nav/dashboard",
		"location": "ApplicationMenu",
		"text": "Administration"
	}
}`

const finalizerName = "commonui.operators.ibm.com"
const finalizerName1 = "commonui1.operators.ibm.com"

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
const CP4INavConfigName = "icp4i"

//nolint
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
		  "Grafana, v7.5.12, AGPL"
		],
		"logoUrl": "IBM Cloud Pak | Administration Hub",
		"logoAltText": "IBM Cloud Pak | Administration Hub"
	  },
	  "header": {
		"disabledItems": [
		  "createResource",
		  "catalog"
		],
		"docUrlMapping": "https://ibm.biz/cpcs_adminui",
		"logoAltText": "IBM Cloud Pak Administration Hub",
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
		  "id": "home",
		  "label": "Home",
		  "url": "/common-nav/dashboard",
		  "iconUrl": "/common-nav/graphics/home.svg",
		  "isAuthorized": [
			"ClusterAdministrator",
			"CloudPakAdministrator"
		  ]
		},
		{
			"id": "id-access",
			"label": "Identity and access",
			"serviceId": "webui-nav",
			"iconUrl": "/common-nav/graphics/password.svg"
		},
		{
			"id": "providers",
			"parentId": "id-access",
			"label": "Identity providers",
			"serviceId": "webui-nav",
			"url": "/common-nav/identity-access/realms"
		},
		{
			"id": "teams-ids",
			"parentId": "id-access",
			"label": "Teams and service IDs",
			"serviceId": "webui-nav",
			"url": "/common-nav/identity-access/teams"
		},
		{
		  "detectionServiceName": true,
		  "id": "licensing",
		  "label": "Licensing",
		  "serviceId": "ibm-license-service-reporter",
		  "url": "/license-service-reporter",
		  "iconUrl": "/common-nav/graphics/identification.svg",
		  "isAuthorized": [
			"ClusterAdministrator",
			"CloudPakAdministrator"
		  ]
		},
		{
		  "detectionServiceName": true,
		  "id": "metering",
		  "label": "Metering",
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
			"CloudPakAdministrator",
			"Operator"
		  ],
		  "label": "Monitoring",
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
		  "serviceId": "kibana",
		  "serviceName": "kibana",
		  "target": "_blank",
		  "url": "/kibana",
		  "iconUrl": "/common-nav/graphics/catalog.svg"
		}
	  ]
	}
  }
`

//nolint
var CP4INavConfig = `
{
	"apiVersion": "foundation.ibm.com/v1",
	"kind": "NavConfiguration",
	"metadata": {
	  "labels": {
		"app.kubernetes.io/instance": "icp4i",
		"app.kubernetes.io/managed-by": "ibm-commonui-operator",
		"app.kubernetes.io/name": "ibm-commonui-operator",
		"name": "icp4i"
	  },
	  "name": "icp4i"
	},
	"spec": {
	  "header": {
		"disabledItems": [
		  "createResource",
		  "catalog",
		  "bookmark"
		],
		"logoAltText": "Cloud Pak for Integration",
		"logoUrl": "/common-nav/graphics/ibm-cloudpak-integration.svg"
	  },
	  "navItems": [
		{
		  "detectionServiceName": true,
		  "id": "metering",
		  "label": "Metering",
		  "serviceId": "metering-ui",
		  "serviceName": "metering-ui",
		  "url": "/metering/dashboard?ace_config={ 'showClusterData': false }\u0026dashboard=cpi.icp.main"
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
		  "serviceId": "monitoring-ui",
		  "serviceName": "ibm-monitoring-grafana",
		  "target": "_blank",
		  "url": "/grafana"
		},
		{
			"id": "id-access",
			"label": "Identity and access",
			"serviceId": "webui-nav"
		},
		{
			"id": "providers",
			"parentId": "id-access",
			"label": "Identity providers",
			"serviceId": "webui-nav",
			"url": "/common-nav/identity-access/realms?useNav=icp4i"
		},
		{
			"id": "teams-ids",
			"parentId": "id-access",
			"label": "Teams and service IDs",
			"serviceId": "webui-nav",
			"url": "/common-nav/identity-access/teams?useNav=icp4i"
		},
		{
		  "detectionServiceName": true,
		  "id": "logging",
		  "label": "Logging",
		  "serviceId": "kibana",
		  "serviceName": "kibana",
		  "target": "_blank",
		  "url": "/kibana"
		},
		{
		  "detectionServiceName": true,
		  "id": "releases",
		  "label": "Helm Releases",
		  "serviceId": "catalog-ui",
		  "serviceName": "catalog-ui",
		  "url": "/catalog/instances?useNav=icp4i"
		},
		{
		  "detectionServiceName": true,
		  "id": "repos",
		  "label": "Helm Repositories",
		  "serviceId": "catalog-ui",
		  "serviceName": "catalog-ui",
		  "url": "/catalog/repositories?useNav=icp4i"
		},
		{
		  "detectionServiceName": true,
		  "id": "licensing",
		  "label": "Licensing",
		  "serviceId": "ibm-license-service-reporter",
		  "url": "/license-service-reporter",
		  "isAuthorized": [
			"ClusterAdministrator",
			"CloudPakAdministrator"
		  ]
		}
	  ]
	}
  }
`
