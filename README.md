# ibm-commonui-operator

> **Important:** Do not install this operator directly. Only install this operator using the IBM Common Services Operator. For more information about installing this operator and other Common Services operators, see [Installer documentation](http://ibm.biz/cpcs_opinstall).
> If you are using this operator as part of an IBM Cloud Pak, see the documentation for that IBM Cloud Pak to learn more about how to install and use the operator service. For more information about IBM Cloud Paks, see [IBM Cloud Paks that use Common Services](http://ibm.biz/cpcs_cloudpaks).

You can use the ibm-commonui-operator to install the Common Web UI service for the IBM Cloud Platform Common Services and access the Common Web UI console. You can use the Common Web UI console to access information and features from other IBM Cloud Platform Common Services or IBM Cloud Paks that you install.

For more information about the available IBM Cloud Platform Common Services, see the [IBM Knowledge Center](http://ibm.biz/cpcsdocs).

## Supported platforms

Red Hat OpenShift Container Platform 4.2 or newer installed on one of the following platforms:

- Linux x86_64
- Linux on Power (ppc64le)
- Linux on IBM Z and LinuxONE

## Operator versions

- 1.2.0
- 1.1.0

## Prerequisites

The Common Web UI service has dependencies on other IBM Cloud Platform Common Services. Before you install this operator, you need to first install the operator dependencies and prerequisites:

- For the list of operator dependencies, see the IBM Knowledge Center [Common Services dependencies documentation](http://ibm.biz/cpcs_opdependencies).

- For the list of prerequisites for installing the operator, see the IBM Knowledge Center [Preparing to install services documentation](http://ibm.biz/cpcs_opinstprereq).

## Documentation

To install the operator with the IBM Common Services Operator follow the the installation and configuration instructions within the IBM Knowledge Center.

- If you are using the operator as part of an IBM Cloud Pak, see the documentation for that IBM Cloud Pak [IBM Cloud Paks that use Common Services](http://ibm.biz/cpcs_cloudpaks).
- If you are using the operator with an IBM Containerized Software, see the IBM Cloud Platform Common Services Knowledge Center [Installer documentation](http://ibm.biz/cpcs_opinstall).

## SecurityContextConstraints Requirements

The Common UI service supports running with the OpenShift Container Platform 4.3 default restricted Security Context Constraints (SCCs).

For more information about the OpenShift Container Platform Security Context Constraints, see [Managing Security Context Constraints.](https://docs.openshift.com/container-platform/4.3/authentication/managing-security-context-constraints.html)

## (Optional) Developer guide

If, as a developer, you are looking to build and test this operator to try out and learn more about the operator and its capabilities, you can use the following developer guide. This guide provides commands for a quick install and initial validation for running the operator. For information about accessing and using the console, see the IBM Cloud Platform Common Services Knowledge Center [Common Web UI documentation](http://ibm.biz/cpcs_opcwebui).

> **Important:** The following developer guide is provided as-is and only for trial and education purposes. IBM and IBM Support does not provide any support for the usage of the operator with this developer guide. For the official supported install and usage guide for the operator, see the the IBM Knowledge Center documentation for your IBM Cloud Pak or for IBM Cloud Platform Common Services.

### Quick start guide

Use the following quick start commands for building and testing the operator:

Prerequisite:

- Git
- Go programming version 1.12+
- Linting Tools:

      | Linting tool | Version |
      | ------------ | ------- |
      | [Haskell Dockerfile Linter (hadolint)](https://github.com/hadolint/hadolint#install) | [v1.17.2](https://github.com/hadolint/hadolint/releases/tag/v1.17.2) |
      | [ShellCheck](https://github.com/koalaman/shellcheck#installing) | [v0.7.0](https://github.com/koalaman/shellcheck/releases/tag/v0.7.0) |
      | [yamllint](https://github.com/adrienverge/yamllint#installation) | [v1.17.0](https://github.com/adrienverge/yamllint/releases/tag/v1.17.0)
      | [Helm client](https://helm.sh/docs/using_helm/#install-helm) | [v2.10.0](https://github.com/helm/helm/releases/tag/v2.10.0) |
      | [golangci-lint](https://github.com/golangci/golangci-lint#install) | [v1.18.0](https://github.com/golangci/golangci-lint/releases/tag/v1.18.0) |
      | [autopep8](https://github.com/hhatto/autopep8#installation) | [v1.4.4](https://github.com/hhatto/autopep8/releases/tag/v1.4.4) |
      | [Markdownlint (mdl)](https://github.com/markdownlint/markdownlint#installation) | [v0.5.0](https://github.com/markdownlint/markdownlint/releases/tag/v0.5.0) |
      | [awesome_bot](https://github.com/dkhamsing/awesome_bot#installation) | [1.19.1](https://github.com/dkhamsing/awesome_bot/releases/tag/1.19.1) |
      | [Sass-lint](https://github.com/sasstools/sass-lint#install) | [v1.13.1](https://github.com/sasstools/sass-lint/releases/tag/v1.13.1) |
      | [Tslint](https://github.com/palantir/tslint#installation--usage) | [v5.18.0](https://github.com/palantir/tslint/releases/tag/5.18.0)
      | [Prototool](https://github.com/uber/prototool/blob/dev/docs/install.md) | `7df3b95` |
      | [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) | `3792095` |

- Set up `GIT_HOST` to override the setting for your custom path:

  ```bash
  export GIT_HOST=github.com/<YOUR_GITHUB_ID>
  ```

- Run the `linter` and `test` before you build the binary:

  ```bash
  make check
  make test
  make build
  ```

- Build and push the docker image for local development:

  ```bash
  export IMG=<YOUR_CUSTOMIZED_IMAGE_NAME>
  export REGISTRY=<YOUR_CUSTOMIZED_IMAGE_REGISTRY>
  make build-push-images
  ```

  > **Note:** You need to log in the docker registry before you run the preceding command.

### Debugging guide

Use the following commands to debug the operator installation:

- Check the CSV installation status:

  ```bash
  oc get csv
  oc describe csv ibm-commonui-operator.v1.2.0
  ```

- Check the custom resource status:

  ```bash
  oc get commonwebuis.operators.ibm.com
  oc describe commonwebuis.operators.ibm.com example-commonwebui
  oc get commonwebuis.operators.ibm.com example-commonwebui -o yaml
  ```

  If there are nodes for the `commonwebuis` instances, the nodes are deployed successfully.
  Additionally, you can check the logs for each of deployed containers for any errors.

- Check the logs for a deployed container:

  ```bash
  oc logs <status.nodeName>
  ```

- Check the operator status and log:

  ```bash
  oc describe po -l name=ibm-commonui-operator
  oc logs -f $(oc get po -l name=ibm-commonui-operator -o name)
  ```

- Access the common ui via the route name
Use the following command to obtain the route for the common ui:

```bash
oc get route
```

Use the URL from the HOST/PORT column for the route named cp-console.

### End-to-End testing

For more instructions on how to run end-to-end testing with the Operand Deployment Lifecycle Manager, see [ODLM guide](https://github.com/IBM/operand-deployment-lifecycle-manager/blob/master/docs/dev/e2e.md#running-e2e-tests).
