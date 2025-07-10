#Define arch-specific image by appending -ARCH to the tag
ARCH_IMAGE="$BUILD_IMAGE-$ARCH"

echo "Building image $ARCH_IMAGE"
#docker buildx build --platform linux/${ARCH} -f Dockerfile --load --build-arg VCS_REF=$VCS_REF --build-arg VCS_URL=$VCS_URL -t $ARCH_IMAGE .

echo "Pushing to $DOCKER_REGISTRY"
#docker push $ARCH_IMAGE
# Skopeo copy can be used instead if image is being pushed as oci
#skopeo copy --dest-creds=$DOCKER_USER:$DOCKER_PASS --format=v2s2 docker-daemon:$ARCH_IMAGE docker://$ARCH_IMAGE
