# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.19.4
ARG IMAGE_FPM_BUILDER=dockter/fpm:latest
ARG IMAGE_FINAL=alpine

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT 2022-12-22
LABEL Name="senzing/servegrpc-builder" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Build arguments.

ARG PROGRAM_NAME="unknown"
ARG BUILD_VERSION=0.0.0
ARG BUILD_ITERATION=0
ARG GO_PACKAGE_NAME="unknown"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/${GO_PACKAGE_NAME}

# Build go program.

WORKDIR ${GOPATH}/src/${GO_PACKAGE_NAME}
RUN make build

# --- Test go program ---------------------------------------------------------

# Run unit tests.

# RUN go get github.com/jstemmer/go-junit-report \
#  && mkdir -p /output/go-junit-report \
#  && go test -v ${GO_PACKAGE_NAME}/... | go-junit-report > /output/go-junit-report/test-report.xml

# Copy binaries to /output.

RUN mkdir -p /output \
      && cp -R ${GOPATH}/src/${GO_PACKAGE_NAME}/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: fpm_builder
#  - Reference: https://github.com/jordansissel/fpm/blob/master/Dockerfile
#  - FPM: https://fpm.readthedocs.io/en/latest/cli-reference.html
# -----------------------------------------------------------------------------

FROM ${IMAGE_FPM_BUILDER} as fpm_builder
ENV REFRESHED_AT 2022-12-22
LABEL Name="senzing/servegrpc-fpm-builder" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Use arguments from prior stage.

ARG PROGRAM_NAME
ARG BUILD_VERSION
ARG BUILD_ITERATION
ARG GO_PACKAGE_NAME

# Copy files from prior stage.

COPY --from=go_builder "/output/linux/*"       "/output/linux/"

# Create RPM package.

RUN fpm \
      --input-type dir \
      --output-type rpm \
      --name ${PROGRAM_NAME} \
      --package /output/${PROGRAM_NAME}-${BUILD_VERSION}.rpm \
      --version ${BUILD_VERSION} \
      --iteration ${BUILD_ITERATION} \
      /output/linux=/usr/bin

# Create DEB package.

RUN fpm \
      --deb-no-default-config-files \
      --input-type dir \
      --iteration ${BUILD_ITERATION} \
      --name ${PROGRAM_NAME} \
      --output-type deb \
      --package /output/${PROGRAM_NAME}-${BUILD_VERSION}.deb \
      --version ${BUILD_VERSION} \
      /output/linux/=/usr/bin

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT 2022-12-22
LABEL Name="senzing/servegrpc" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Use arguments from prior stage.

ARG PROGRAM_NAME

# Copy files from prior step.

COPY --from=fpm_builder "/output/linux/${PROGRAM_NAME}"    "/output/linux/${PROGRAM_NAME}"
COPY --from=fpm_builder "/output/go-servegrpc-*"           "/output/"

CMD ["/bin/bash"]
