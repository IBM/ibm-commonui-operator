# Copyright 2020 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This repo is build locally for dev/test by default;
# Override this variable in CI env.
BUILD_LOCALLY ?= 1

# Image URL to use all building/pushing image targets;
# Use your own docker registry and image name for dev/test by overridding the IMG and REGISTRY environment variable.
IMG ?= ibm-commonui-operator
REGISTRY ?= "hyc-cloud-private-integration-docker-local.artifactory.swg-devops.com/ibmcom"
REGISTRY_DEV ?= quay.io/ericabr
CSV_VERSION ?= $(VERSION)
NAMESPACE=ibm-common-services
LEGACY_TAG ?= 3.2.5
COMMON_TAG ?= 1.5.0
DASHBOARD_TAG ?= 1.1.1

# Set the registry and tag for the operand/operator images
OPERAND_REGISTRY ?= $(REGISTRY)
COMMON_WEB_UI_OPERAND_TAG ?= $(COMMON_TAG)
COMMON_WEB_UI_OPERAND_TAG_AMD ?= $(COMMON_TAG)-amd64
COMMON_WEB_UI_OPERAND_TAG_PPC ?= $(COMMON_TAG)-ppc64le
COMMON_WEB_UI_OPERAND_TAG_Z ?= $(COMMON_TAG)-s390x
LEGACY_HEADER_OPERAND_TAG ?= $(LEGACY_TAG)
LEGACY_HEADER_OPERAND_TAG_AMD ?= $(LEGACY_TAG)-amd64
LEGACY_HEADER_OPERAND_TAG_PPC ?= $(LEGACY_TAG)-ppc64le
LEGACY_HEADER_OPERAND_TAG_Z ?= $(LEGACY_TAG)-s390x
COMMONUI_OPERATOR_TAG ?= $(VERSION)
# Github host to use for checking the source tree;
# Override this variable ue with your own value if you're working on forked repo.
GIT_HOST ?= github.com/IBM

PWD := $(shell pwd)
BASE_DIR := $(shell basename $(PWD))

# Keep an existing GOPATH, make a private one if it is undefined
GOPATH_DEFAULT := $(PWD)/.go
export GOPATH ?= $(GOPATH_DEFAULT)
GOBIN_DEFAULT := $(GOPATH)/bin
export GOBIN ?= $(GOBIN_DEFAULT)
TESTARGS_DEFAULT := "-v"
export TESTARGS ?= $(TESTARGS_DEFAULT)
DEST := $(GOPATH)/src/$(GIT_HOST)/$(BASE_DIR)
VERSION ?= $(shell cat ./version/version.go | grep "Version =" | awk '{ print $$3}' | tr -d '"')

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
    TARGET_OS ?= linux
    XARGS_FLAGS="-r"
else ifeq ($(LOCAL_OS),Darwin)
    TARGET_OS ?= darwin
    XARGS_FLAGS=
else
    $(error "This system's OS $(LOCAL_OS) isn't recognized/supported")
endif

ARCH := $(shell uname -m)
LOCAL_ARCH := "amd64"
ifeq ($(ARCH),x86_64)
    LOCAL_ARCH="amd64"
else ifeq ($(ARCH),ppc64le)
    LOCAL_ARCH="ppc64le"
else ifeq ($(ARCH),s390x)
    LOCAL_ARCH="s390x"
else
    $(error "This system's ARCH $(ARCH) isn't recognized/supported")
endif

# Setup DOCKER_BUILD_OPTS after all includes complete
#Variables for redhat ubi certification required labels
IMAGE_NAME=$(IMG)
IMAGE_DISPLAY_NAME=IBM CommonUI Operator
IMAGE_MAINTAINER=ericabr@us.ibm.com
IMAGE_VENDOR=IBM
IMAGE_VERSION=$(VERSION)
IMAGE_DESCRIPTION=Operator used to install a service to display a common header and IAM pages in a kubernetes cluster
IMAGE_SUMMARY=$(IMAGE_DESCRIPTION)
IMAGE_OPENSHIFT_TAGS=ibm-common-ui
$(eval WORKING_CHANGES := $(shell git status --porcelain))
$(eval BUILD_DATE := $(shell date +%m/%d@%H:%M:%S))
$(eval GIT_COMMIT := $(shell git rev-parse --short HEAD))
$(eval VCS_REF := $(GIT_COMMIT))
IMAGE_RELEASE=$(VCS_REF)
GIT_REMOTE_URL = $(shell git config --get remote.origin.url)
$(eval DOCKER_BUILD_OPTS := --build-arg "IMAGE_NAME=$(IMAGE_NAME)" --build-arg "IMAGE_DISPLAY_NAME=$(IMAGE_DISPLAY_NAME)" --build-arg "IMAGE_MAINTAINER=$(IMAGE_MAINTAINER)" --build-arg "IMAGE_VENDOR=$(IMAGE_VENDOR)" --build-arg "IMAGE_VERSION=$(IMAGE_VERSION)" --build-arg "IMAGE_RELEASE=$(IMAGE_RELEASE)" --build-arg "IMAGE_DESCRIPTION=$(IMAGE_DESCRIPTION)" --build-arg "IMAGE_SUMMARY=$(IMAGE_SUMMARY)" --build-arg "IMAGE_OPENSHIFT_TAGS=$(IMAGE_OPENSHIFT_TAGS)" --build-arg "VCS_REF=$(VCS_REF)" --build-arg "VCS_URL=$(GIT_REMOTE_URL)" --build-arg "SELF_METER_IMAGE_TAG=$(SELF_METER_IMAGE_TAG)")


all: fmt check test coverage build images

ifneq ("$(realpath $(DEST))", "$(realpath $(PWD))")
    $(error Please run 'make' from $(DEST). Current directory is $(PWD))
endif

include common/Makefile.common.mk
############################################################
# SHA section
############################################################

.PHONY: get-all-operand-image-sha
get-all-operand-image-sha: get-common-web-ui-image-sha get-legacy-header-image-sha get-commonui-operator-image-sha
	@echo Got SHAs for all operand images

.PHONY: get-commonui-operator-image-sha
get-commonui-operator-image-sha:
	@echo Get SHA for ibm-commonui-operator:$(COMMONUI_OPERATOR_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) ibm-commonui-operator $(COMMONUI_OPERATOR_TAG) COMMONUI_OPERATOR_TAG_OR_SHA
	@echo Get SHA for ibm-commonui-operator:$(COMMONUI_OPERATOR_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) ibm-commonui-operator-amd64 $(COMMONUI_OPERATOR_TAG) COMMONUI_OPERATOR_TAG_OR_SHA
	@echo Get SHA for ibm-commonui-operator:$(COMMONUI_OPERATOR_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) ibm-commonui-operator-ppc64le $(COMMONUI_OPERATOR_TAG) COMMONUI_OPERATOR_TAG_OR_SHA
	@echo Get SHA for ibm-commonui-operator:$(COMMONUI_OPERATOR_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) ibm-commonui-operator-s390x $(COMMONUI_OPERATOR_TAG) COMMONUI_OPERATOR_TAG_OR_SHA


.PHONY: get-common-web-ui-image-sha
get-common-web-ui-image-sha:
	@echo Get SHA for common-web-ui:$(COMMON_WEB_UI_OPERAND_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) common-web-ui $(COMMON_WEB_UI_OPERAND_TAG) COMMON_WEB_UI_IMAGE_TAG_OR_SHA
	@echo Get SHA for common-web-ui:$(COMMON_WEB_UI_OPERAND_TAG_AMD)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) common-web-ui $(COMMON_WEB_UI_OPERAND_TAG_AMD) COMMON_WEB_UI_IMAGE_TAG_OR_SHA
	@echo Get SHA for common-web-ui:$(COMMON_WEB_UI_OPERAND_TAG_PPC)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) common-web-ui $(COMMON_WEB_UI_OPERAND_TAG_PPC) COMMON_WEB_UI_IMAGE_TAG_OR_SHA
	@echo Get SHA for common-web-ui:$(COMMON_WEB_UI_OPERAND_TAG_Z)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) common-web-ui $(COMMON_WEB_UI_OPERAND_TAG_Z) COMMON_WEB_UI_IMAGE_TAG_OR_SHA

.PHONY: get-legacy-header-image-sha
get-legacy-header-image-sha:
	@echo Get SHA for platform-header:$(LEGACY_HEADER_OPERAND_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) icp-platform-header $(LEGACY_HEADER_OPERAND_TAG) LEGACYHEADER_IMAGE_TAG_OR_SHA
	@echo Get SHA for platform-header:$(LEGACY_HEADER_OPERAND_TAG_AMD)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) icp-platform-header $(LEGACY_HEADER_OPERAND_TAG_AMD) LEGACYHEADER_IMAGE_TAG_OR_SHA
	@echo Get SHA for platform-header:$(LEGACY_HEADER_OPERAND_TAG_PPC)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) icp-platform-header $(LEGACY_HEADER_OPERAND_TAG_PPC) LEGACYHEADER_IMAGE_TAG_OR_SHA
	@echo Get SHA for platform-header:$(LEGACY_HEADER_OPERAND_TAG_Z)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) icp-platform-header $(LEGACY_HEADER_OPERAND_TAG_Z) LEGACYHEADER_IMAGE_TAG_OR_SHA

.PHONY: get-oper-image-sha
get-oper-image-sha:
	@echo Get SHA for ibm-commonui-operator:$(VERSION)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) ibm-commonui-operator $(CSV_VERSION) CSV_VERSION

.PHONY: get-dashboard-image-sha
get-dashboard-image-sha:
	@echo Get SHA for ibm-dashboard-data-collector:$(DASHBOARD_TAG)
	@common/scripts/get_image_sha_digest.sh $(OPERAND_REGISTRY) ibm-dashboard-data-collector $(DASHBOARD_TAG) DASHBOARD_DATA_COLL_IMAGE_TAG_OR_SHA

############################################################
# work section
############################################################
$(GOBIN):
	@echo "create gobin"
	@mkdir -p $(GOBIN)

work: $(GOBIN)

############################################################
# format section
############################################################

# All available format: format-go format-protos format-python
# Default value will run all formats, override these make target with your requirements:
#    eg: fmt: format-go format-protos
fmt: format-go format-protos format-python

############################################################
# check section
############################################################

check: lint

# All available linters: lint-dockerfiles lint-scripts lint-yaml lint-copyright-banner lint-go lint-python lint-helm lint-markdown lint-sass lint-typescript lint-protos
# Default value will run all linters, override these make target with your requirements:
#    eg: lint: lint-go lint-yaml
lint: lint-all

############################################################
# test section
############################################################

test:
	@go test ${TESTARGS} ./...

############################################################
# coverage section
############################################################

coverage:
	@common/scripts/codecov.sh ${BUILD_LOCALLY}

############################################################
# install operator sdk section
############################################################

install-operator-sdk: 
	@operator-sdk version 2> /dev/null ; if [ $$? -ne 0 ]; then ./common/scripts/install-operator-sdk.sh; fi

############################################################
# build section
############################################################

build: build-amd64 build-ppc64le build-s390x

build-amd64:
	@echo "Building the ${IMG} amd64 binary..."
	@GOARCH=amd64 common/scripts/gobuild.sh build/_output/bin/$(IMG) ./cmd/manager

build-ppc64le:
	@echo "Building the ${IMG} ppc64le binary..."
	@GOARCH=ppc64le common/scripts/gobuild.sh build/_output/bin/$(IMG)-ppc64le ./cmd/manager

build-s390x:
	@echo "Building the ${IMG} s390x binary..."
	@GOARCH=s390x common/scripts/gobuild.sh build/_output/bin/$(IMG)-s390x ./cmd/manager

local:
	@GOOS=darwin common/scripts/gobuild.sh build/_output/bin/$(IMG) ./cmd/manager

############################################################
# images section
############################################################

ifeq ($(BUILD_LOCALLY),0)
    export CONFIG_DOCKER_TARGET = config-docker
config-docker:
endif

image-dev: build-image-amd64
	docker tag $(REGISTRY)/$(IMG)-amd64:$(VERSION) $(REGISTRY_DEV)/$(IMG):$(VERSION)
	docker push $(REGISTRY_DEV)/$(IMG):$(VERSION)
build-image-amd64: build-amd64
	@docker build -t $(REGISTRY)/$(IMG)-amd64:$(VERSION) $(DOCKER_BUILD_OPTS) --build-arg "IMAGE_NAME_ARCH=$(IMAGE_NAME)-amd64" -f build/Dockerfile .

build-image-ppc64le: build-ppc64le
	@docker run --rm --privileged multiarch/qemu-user-static:register --reset
	@docker build -t $(REGISTRY)/$(IMG)-ppc64le:$(VERSION) $(DOCKER_BUILD_OPTS) --build-arg "IMAGE_NAME_ARCH=$(IMAGE_NAME)-ppc64le" -f build/Dockerfile.ppc64le .

build-image-s390x: build-s390x
	@docker run --rm --privileged multiarch/qemu-user-static:register --reset
	@docker build -t $(REGISTRY)/$(IMG)-s390x:$(VERSION) $(DOCKER_BUILD_OPTS) --build-arg "IMAGE_NAME_ARCH=$(IMAGE_NAME)-s390x" -f build/Dockerfile.s390x .

push-image-amd64: $(CONFIG_DOCKER_TARGET) build-image-amd64
	@docker push $(REGISTRY)/$(IMG)-amd64:$(VERSION)

push-image-ppc64le: $(CONFIG_DOCKER_TARGET) build-image-ppc64le
	@docker push $(REGISTRY)/$(IMG)-ppc64le:$(VERSION)

push-image-s390x: $(CONFIG_DOCKER_TARGET) build-image-s390x
	@docker push $(REGISTRY)/$(IMG)-s390x:$(VERSION)

############################################################
# multiarch-image section
############################################################

images: push-image-amd64 push-image-ppc64le push-image-s390x multiarch-image

multiarch-image:
	@curl -L -o /tmp/manifest-tool https://github.com/estesp/manifest-tool/releases/download/v1.0.0/manifest-tool-linux-amd64
	@chmod +x /tmp/manifest-tool
	/tmp/manifest-tool push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template $(REGISTRY)/$(IMG)-ARCH:$(VERSION) --target $(REGISTRY)/$(IMG) --ignore-missing
	/tmp/manifest-tool push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template $(REGISTRY)/$(IMG)-ARCH:$(VERSION) --target $(REGISTRY)/$(IMG):$(VERSION) --ignore-missing

############################################################
# clean section
############################################################
clean:
	rm -rf build/_output

############################################################
# CSV section
############################################################
csv: ## Push CSV package to the catalog
	@RELEASE=${CSV_VERSION} common/scripts/push-csv.sh

############################################################
# CSV section
############################################################
delete-csv: ## Delete CSV package to the catalog
	@RELEASE=${CSV_VERSION} common/scripts/delete-csv.sh

.PHONY: all work build check lint test coverage images multiarch-image

############################################################
# Install/uninstall
############################################################

install: ## Install all resources (CR/CRD's, RBCA and Operator)
	@echo ....... Set environment variables ......
	- export DEPLOY_DIR=deploy/crds
	- export WATCH_NAMESPACE=${NAMESPACE}
	@echo ....... Applying CRDS and Operator .......
	- for crd in $(shell ls deploy/crds/*_crd.yaml); do kubectl apply -f $${crd}; done
	@echo ....... Applying RBAC .......
	- kubectl apply -f deploy/service_account.yaml -n ${NAMESPACE}
	- kubectl apply -f deploy/role.yaml -n ${NAMESPACE}
	- kubectl apply -f deploy/role_binding.yaml -n ${NAMESPACE}
	@echo ....... Applying Operator .......
	- kubectl apply -f deploy/olm-catalog/${BASE_DIR}/${CSV_VERSION}/${BASE_DIR}.v${CSV_VERSION}.clusterserviceversion.yaml -n ${NAMESPACE}
	@echo ....... Creating the Instance .......
	- for cr in $(shell ls deploy/crds/*_cr.yaml); do kubectl apply -f $${cr} -n ${NAMESPACE}; done

uninstall: ## Uninstall all that all performed in the $ make install
	@echo ....... Uninstalling .......
	@echo ....... Deleting CR .......
	- for cr in $(shell ls deploy/crds/*_cr.yaml); do kubectl delete -f $${cr} -n ${NAMESPACE}; done
	@echo ....... Deleting Operator .......
	- kubectl delete -f deploy/olm-catalog/${BASE_DIR}/${CSV_VERSION}/${BASE_DIR}.v${CSV_VERSION}.clusterserviceversion.yaml -n ${NAMESPACE}
	@echo ....... Deleting CRDs.......
	- for crd in $(shell ls deploy/crds/*_crd.yaml); do kubectl delete -f $${crd}; done
	@echo ....... Deleting Rules and Service Account .......
	- kubectl delete -f deploy/role_binding.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/service_account.yaml -n ${NAMESPACE}
	- kubectl delete -f deploy/role.yaml -n ${NAMESPACE}