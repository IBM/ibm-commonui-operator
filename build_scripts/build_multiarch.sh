echo "Pushing multi-arch image"

#make docker:manifest-tool
MANIFEST_TOOL ?= ./manifest-tool
MANIFEST_TOOL_VERSION ?= v1.0.3
MANIFEST_TOOL_OS ?= linux
MANIFEST_TOOL_ARCH ?= amd64
MANIFEST_TOOL_URL ?= https://github.com/estesp/manifest-tool/releases/download/$(MANIFEST_TOOL_VERSION)/manifest-tool-$(MANIFEST_TOOL_OS)-$(MANIFEST_TOOL_ARCH)

echo "Installing manifest-tool $(MANIFEST_TOOL_VERSION) ($(MANIFEST_TOOL_OS)-$(MANIFEST_TOOL_ARCH)) from $(MANIFEST_TOOL_URL)" && \
                curl '-#' -fL -o $(MANIFEST_TOOL) $(MANIFEST_TOOL_URL) && \
                chmod +x $(MANIFEST_TOOL) \
$(MANIFEST_TOOL) --version

## Push the manifest to a Docker registry
$(MANIFEST_TOOL) --debug push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template "$BUILD_IMAGE-ARCH" --target "$BUILD_IMAGE"
