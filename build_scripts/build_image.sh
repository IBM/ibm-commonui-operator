#Define arch-specific image by appending -ARCH to the tag
export ARCH_IMAGE="$BUILD_IMAGE-$ARCH"

export NODE_OPTIONS="--max_old_space_size=4096 --openssl-legacy-provider"
env

echo "================================================="
echo "Installing Go                                    "
echo "================================================="
export GO_VERSION=1.23.11
wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
go install goimports

echo "================================================="
echo "BUILDING IMAGE                                   "
echo "================================================="
# Turn lint off for now - make install lint image
make build-image
miret=$?
if [ $miret -ne 0 ]; then
  echo "make build-image failed with code $miret"
  exit $miret
fi

echo "================================================="
echo "PUSHING IMAGE                                    "
echo "================================================="
echo "Pushing $ARCH_IMAGE"
#skopeo copy --dest-creds=$DOCKER_USER:$DOCKER_PASS --format=v2s2 docker-daemon:$ARCH_IMAGE docker://$ARCH_IMAGE
docker push $ARCH_IMAGE
dpret=$?
if [ $dpret -ne 0 ]; then
  echo "docker push failed with code $dpret"
  exit $dpret
fi

