FROM docker-na-public.artifactory.swg-devops.com/hyc-cloud-private-edge-docker-local/build-images/ubi9-minimal:latest

ARG IMAGE_NAME
ARG IMAGE_DISPLAY_NAME
ARG IMAGE_NAME_ARCH
ARG IMAGE_MAINTAINER
ARG IMAGE_VENDOR
ARG IMAGE_VERSION
ARG IMAGE_RELEASE
ARG IMAGE_DESCRIPTION
ARG IMAGE_SUMMARY
ARG IMAGE_OPENSHIFT_TAGS
ARG SELF_METER_IMAGE_TAG
ARG VCS_REF
ARG VCS_URL

LABEL org.label-schema.vendor="$IMAGE_VENDOR" \
      org.label-schema.name="$IMAGE_NAME_ARCH" \
      org.label-schema.description="$IMAGE_DESCRIPTION" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.license="Apache 2.0" \
      org.label-schema.schema-version="1.0" \
      name="$IMAGE_NAME" \
      maintainer="$IMAGE_MAINTAINER" \
      vendor="$IMAGE_VENDOR" \
      version="$IMAGE_VERSION" \
      release="$IMAGE_RELEASE" \
      description="$IMAGE_DESCRIPTION" \
      summary="$IMAGE_SUMMARY" \
      io.k8s.display-name="$IMAGE_DISPLAY_NAME" \
      io.k8s.description="$IMAGE_DESCRIPTION" \
      io.openshift.tags="$IMAGE_OPENSHIFT_TAGS"

ENV BINARY=/usr/local/bin/ibm-commonui-operator \
  USER_UID=1001 \
  USER_NAME=ibm-commonui-operator

# install the binary
COPY build/_output/bin/ibm-commonui-operator ${BINARY}

# copy licenses
RUN mkdir /licenses
COPY LICENSE /licenses

ENTRYPOINT ["ibm-commonui-operator"]

USER ${USER_UID}
