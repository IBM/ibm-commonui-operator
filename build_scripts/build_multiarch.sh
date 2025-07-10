echo "Pushing multi-arch image"

make docker:manifest-tool
#manifest-tool --debug push from-args --platforms linux/amd64,linux/ppc64le,linux/s390x --template "$BUILD_IMAGE-ARCH" --target "$BUILD_IMAGE"
