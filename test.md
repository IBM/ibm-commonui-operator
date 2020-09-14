cat << EOF | oc apply -f -
apiVersion: operator.openshift.io/v1alpha1
kind: ImageContentSourcePolicy
metadata:
  name: ibm-cs-daily-mirror
spec:
  repositoryDigestMirrors:
  - mirrors:
    - hyc-cloud-private-daily-docker-local.artifactory.swg-devops.com/ibmcom
    source: quay.io/opencloudio
EOF

export ARTIFACTORY_USER=llcao@cn.ibm.com
export ARTIFACTORY_TOKEN=AKCp5cbwQnxszNViXEjCybVZpwvVfLLjJDPTjbJcxbCxxhYtPy26i6b7T2GbUdka1n9pyCWHJ
pull_secret=$(echo -n "$ARTIFACTORY_USER:$ARTIFACTORY_TOKEN" | base64 -w0)
oc get secret/pull-secret -n openshift-config -o jsonpath='{.data.\.dockerconfigjson}' | base64 -d | sed -e 's|:{|:{"hyc-cloud-private-daily-docker-local.artifactory.swg-devops.com":{"auth":"'$pull_secret'"\},|' > /tmp/dockerconfig.json
oc set data secret/pull-secret -n openshift-config --from-file=.dockerconfigjson=/tmp/dockerconfig.json

cat << EOF | oc apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: opencloud-operators
  namespace: openshift-marketplace
spec:
  displayName: IBMCS Operators
  publisher: IBM
  sourceType: grpc
  image: quay.io/morvencao/ibm-common-service-catalog:latest
  updateStrategy:
    registryPoll:
      interval: 45m
EOF

oc -n openshift-marketplace get catsrc
oc -n openshift-marketplace get pod

oc create ns common-service

Install ibm-common-services-operator from OCP Console UI in common-service namespace

cat << EOF | oc apply -f -
apiVersion: operator.ibm.com/v3
kind: CommonService
metadata:
  name: common-service
  namespace: ibm-common-services
spec:
  size: small
  services:
  - name: ibm-mongodb-operator
    spec:
      mongoDB:
        storageClass: rook-ceph-cephfs-internal
EOF

cat << EOF | oc apply -f -
apiVersion: operator.ibm.com/v1alpha1
kind: OperandRequest
metadata:
  name: common-service
  namespace: ibm-common-services
spec:
  requests:
    - operands:
        - name: ibm-cert-manager-operator
        - name: ibm-mongodb-operator
        - name: ibm-iam-operator
        - name: ibm-healthcheck-operator
        - name: ibm-management-ingress-operator
        - name: ibm-licensing-operator
        - name: ibm-commonui-operator
        - name: ibm-ingress-nginx-operator
        - name: ibm-auditlogging-operator
        - name: ibm-platform-api-operator
      registry: common-service
EOF

oc -n ibm-common-services delete operandrequest common-service
Uninstall the ibm-common-services-operator and ODLM from OCP Console
oc get sub -A
oc get csv -A

oc delete ns common-service ibm-common-services


oc -n openshift-marketplace get pod







ibm-iam-operand-privileged:

  jobs:
  - name: oidc-client-registration
    containers:
    - name: oidc-client-registration
      image: quay.io/opencloudio/icp-platform-auth
      repo: git@github.ibm.com:IBMPrivateCloud/platform-auth-service.git
      requiredPermissions: [] # doesn't require extra permission, just register the OIDC client ID and client secret
  depoyments:
  - name: auth-idp
    initContainers:
    - name: init-mongodb
      image: quay.io/opencloudio/icp-platform-auth
      repo: git@github.ibm.com:IBMPrivateCloud/platform-auth-service.git
      requiredPermissions: [] # doesn't require extra permission, just check if mongodb is ready or not
    containers:
    - name: icp-audit-service
      image: quay.io/opencloudio/icp-audit-service
      repo: git@github.ibm.com:IBMPrivateCloud/audit-sidecar-service.git
      requiredPermissions: [] # doesn't require extra permission? TBD with security squad
    - name: platform-auth-service
      image: quay.io/opencloudio/icp-platform-auth
      repo: git@github.ibm.com:IBMPrivateCloud/platform-auth-service.git
      requiredPermissions:
      - resources: ["secrets"]
        group: ""
        verbs: ["get", "patch"]
        clusterScope: false
    - name: platform-identity-provider
      image: quay.io/opencloudio/icp-identity-provider
      repo: git@github.ibm.com:IBMPrivateCloud/platform-identity-provider.git
      requiredPermissions:
      - resources: ["oauthaccesstokens", "oauthclients"]
        group: oauth.openshift.io
        verbs: ["create", "get", "delete", "update"]
        clusterScope: true
      - resources: ["users"]
        group: user.openshift.io
        verbs: ["create", "get"]
        clusterScope: true
    - name: platform-identity-manager
      image: quay.io/opencloudio/icp-identity-manager
      repo: git@github.ibm.com:IBMPrivateCloud/platform-identity-mgmt.git
      requiredPermissions:
      - resources: ["namespaces"]
        group: ""
        verbs: ["get", "delete"]
        clusterScope: true
      - resources: ["pod", "secrets", "services", "persistentvolumeclaims"]
        group: ""
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: v # from the code, it try to access these resources in cluster scope
      - resources: ["clusterrolebindings", "rolebindings"]
        group: rbac.authorization.k8s.io
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: true
      - resources: ["clusterrolebindings", "rolebindings"]
        group: authorization.openshift.io
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: true
      - resources: ["clusterservicebrokers", "instances", "clusterserviceclasses", "clusterserviceplans", "bindings"]
        group: servicecatalog.k8s.io
        verbs: ["get", "list", "watch"]
        clusterScope: true
      - resources: ["deployments"]
        group: extensions
        verbs: ["get", "list", "watch"]
        clusterScope: true
      - resources: ["users", "groups", "identities"]
        group: user.openshift.io
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: true
  - name: auth-pap
    containers:
    - name: icp-audit-service
      image: quay.io/opencloudio/icp-audit-service
      repo: git@github.ibm.com:IBMPrivateCloud/audit-sidecar-service.git
      requiredPermissions: [] # doesn't require extra permission? TBD with security squad
    - name: auth-pap
      image: quay.io/opencloudio/iam-policy-administration
      repo: git@github.ibm.com:IBMPrivateCloud/iam-policy-administration.git
      requiredPermissions: [] # doesn't require extra permission, just execute CURD from mongodb
  - name: auth-pdp
    containers:
    - name: icp-audit-service
      image: quay.io/opencloudio/icp-audit-service
      repo: git@github.ibm.com:IBMPrivateCloud/audit-sidecar-service.git
      requiredPermissions: [] # doesn't require extra permission? TBD with security squad
    - name: auth-pdp
      image: quay.io/opencloudio/iam-policy-decision
      repo: git@github.ibm.com:IBMPrivateCloud/iam-policy-decision.git
      requiredPermissions: [] # doesn't require extra permission, just execute CURD from idmgmt(TBD)



ibm-iam-operand-restricted:
  deployments:
  - name: oidcclient-watcher
    containers:
    - name: oidcclient-watcher
      image: quay.io/opencloudio/icp-oidcclient-watcher
      repo: git@github.ibm.com:IBMPrivateCloud/icp-oidcclient-watcher.git
      requiredPermissions:
      - resources: ["clients"]
        group: oidc.security.ibm.com
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: false
      - resources: ["secrets"]
        group: ""
        verbs: ["create", "get"]
        clusterScope: false
      - resources: ["oauthclients"]
        group: oauth.openshift.io
        verbs: ["create", "get", "update"]
        clusterScope: true
      - resources: ["clusterrolebindings", "rolebindings"]
        group: rbac.authorization.k8s.io
        verbs: ["create", "update"]
        clusterScope: true
      - resources: ["namespaces"]
        group: ""
        verbs: ["get", "update"]
        clusterScope: true
      - resources: ["servicemonitors"]
        group: monitoring.coreos.com
        verbs: ["get", "create"]
        clusterScope: false
  - name: iam-policy-controller
    containers:
    - name: iam-policy-controller
      image: quay.io/opencloudio/iam-policy-controller
      repo: git@github.ibm.com:IBMPrivateCloud/iam-policy-controller.git
      requiredPermissions:
      - resources: ["events"]
        group: ""
        verbs: ["create", "patch"]
        clusterScope: true
      - resources: ["deployments"]
        group: "apps"
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: false
      - resources: ["deployments/status"]
        group: "apps"
        verbs: ["get", "patch", "update"]
        clusterScope: false
      - resources: ["iampolicies"]
        group: "iam.policies.ibm.com"
        verbs: ["create", "delete", "get", "list", "patch", "update", "watch"]
        clusterScope: false
      - resources: ["iampolicies/status"]
        group: "iam.policies.ibm.com"
        verbs: ["get", "patch", "update"]
        clusterScope: false
      - resources: ["clusterrolebindings", "rolebindings"]
        group: rbac.authorization.k8s.io
        verbs: ["list", "watch"]
        clusterScope: true
      - resources: ["namespaces"]
        group: ""
        verbs: ["get", "list"]
        clusterScope: true
      - resources: ["pods"]
        group: ""
        verbs: ["get", "list"]
        clusterScope: false
  - name: secret-watcher
    containers:
    - name: secret-watcher
      image: quay.io/opencloudio/icp-secret-watcher
      repo: git@github.ibm.com:IBMPrivateCloud/icp-secret-watcher.git
      requiredPermissions:
      - resources: ["secrets"]
        group: ""
        verbs: ["get", "list", "watch"]
        clusterScope: true
  jobs:
  - name: security-onboarding
    containers:
    - name: security-onboarding
      image: quay.io/opencloudio/icp-iam-onboarding
      repo: git@github.ibm.com:IBMPrivateCloud/icp-iam-onboarding.git
      requiredPermissions:
      - resources: ["secrets"]
        group: ""
        verbs: ["create", "get", "update"]
        clusterScope: false
  - name: iam-onboarding
    containers:
    - name: iam-onboarding
      image: quay.io/opencloudio/icp-iam-onboarding
      repo: git@github.ibm.com:IBMPrivateCloud/icp-iam-onboarding.git
      requiredPermissions:
      - resources: ["secrets"]
        group: ""
        verbs: ["create", "get", "update"]
        clusterScope: false

