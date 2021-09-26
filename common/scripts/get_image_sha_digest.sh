#!/bin/bash
#
# Copyright 2020 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# Get the SHA from an operand image and put it in operator.yaml and the CSV file.
# Do "docker login" before running this script.
# Run this script from the parent dir by typing "scripts/get-image-sha.sh"
# array of valid TYPE values
declare -A TYPE_LIST
TYPE_LIST[common-web-ui]=2
TYPE_LIST[platform-header]=2
# check the input parms
TYPE=$2
NAME=$1
TAG=$3
if [[ $TAG == "" ]]
then
   echo "Missing parm. Need image type, image name and image tag"
   echo "Examples:"
   echo "   common-web-ui quay.io/opencloudio/common-web-ui 1.2.1"
   echo "   platform-header quay.io/opencloudio/common-web-ui 3.2.4"
   exit 1
fi
# check the TYPE value
if ! [[ ${TYPE_LIST[$TYPE]} ]]
then
   echo "$TYPE is not valid. Must be DM, UI, MCMUI, or REPORT"
   exit 1
fi
# pull the image
IMAGE="$NAME/$TYPE:$TAG"
echo "Pulling image $IMAGE"
docker pull "$IMAGE"
# get the SHA for the image
DIGEST="$(docker images --digests "$NAME/$TYPE" | grep "$TAG" | awk 'FNR==1{print $3}')"
# DIGEST should look like this: sha256:10a844ffaf7733176e927e6c4faa04c2bc4410cf4d4ef61b9ae5240aa62d1456
if [[ $DIGEST != sha256* ]]
then
    echo "Cannot find SHA (sha256:nnnnnnnnnnnn) in digest: $DIGEST"
    exit 1
fi
SHA=$DIGEST
echo "SHA=$SHA"