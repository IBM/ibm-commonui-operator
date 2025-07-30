# Set up Needed Variables
# NOTE - VARIABLES ONLY NEEDED FOR PRs ARE IN THE .pipeline-config-pr.yaml file
export ARTIFACTORY_USERNAME="$(get_env ARTIFACTORY_USERNAME)"
export ARTIFACTORY_TOKEN="$(get_env ARTIFACTORY_TOKEN)"
export DOCKER_REGISTRY="$(get_env DOCKER_REGISTRY)"
export DOCKER_USER="$(get_env DOCKER_USER)"
export DOCKER_PASS="$(get_env DOCKER_PASS)"
export GITHUB_TOKEN="$(get_env GITHUB_TOKEN)"
export GITHUB_USER="$(get_env GITHUB_USER)"
# Build Image Details
export IMAGE_NAME="ibm-commonui-operator"
export BUILD_IMAGE="$DOCKER_REGISTRY/$IMAGE_NAME:$BUILD_TAG"
# Add build harness to path outside of Makefile
export PATH="$PATH:$PWD/build-harness/vendor/"
# Configure Environment
echo -e "machine github.ibm.com\n  login $GITHUB_TOKEN" >> ~/.netrc
chmod 600 ~/.netrc
git config --global --add safe.directory $WORKSPACE/$(load_repo app-repo path)
# Output Paremeters
echo "Current branch : $GIT_BRANCH"
echo "Building commit $GIT_COMMIT"
echo "Using build tag $BUILD_TAG"

# Login is done via pipeline docker config env variable - also for operator
# there is no build harness currently
# Login to root artifactory (to cover both base images and build image)
# make init
# make docker:login DOCKER_REGISTRY=docker-na-public.artifactory.swg-devops.com
