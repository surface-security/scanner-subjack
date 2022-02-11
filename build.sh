#!/bin/sh

# this build.sh is used in every scanner, as is.
# please DO NOT modify/customize it as it will make it harder to propagate changes across scanners
# it checks if there's a custom_build.sh and, if so, executes it instead.
# so use that if required (and only if required)

set -e

cd $(dirname $0)

if [ -e 'custom_build.sh' ]; then
    exec ./custom_build.sh
fi

if [ -e '.git' ]; then
    NAME=$(basename $(cat .git/config| grep '/scanners/' | tr -d ' ') | sed -e 's/\.git$//g')
else
    # fingers crossed name is the same as repo
    NAME=$(basename $(pwd))
fi

IMAGE=${IMAGE_PREFIX:-test}/${NAME}
TAG=${gitlabSourceBranch:-dev}

docker build -t ${IMAGE}:${TAG} .

if [ -n "${BUILD_NUMBER}" -o "$1" = "push" ]; then
    docker push ${IMAGE}:${TAG}
    if [ "${TAG}" = "master" ]; then
        if [ -n "${BUILD_NUMBER}" ]; then
            docker tag ${IMAGE}:master ${IMAGE}:${BUILD_NUMBER}
            docker push ${IMAGE}:${BUILD_NUMBER}
        fi
        docker tag ${IMAGE}:master ${IMAGE}:latest
        docker push ${IMAGE}:latest
    fi
fi
