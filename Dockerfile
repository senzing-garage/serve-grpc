# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.25.5-bookworm@sha256:5117d68695f57faa6c2b3a49a6f3187ec1f66c75d5b080e4360bfe4c1ada398c
ARG IMAGE_FINAL=senzing/senzingsdk-runtime:4.1.0@sha256:e57d751dc0148bb8eeafedb7accf988413f50b54a7e46f25dfe4559d240063e5

# -----------------------------------------------------------------------------
# Stage: senzingsdk_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS senzingsdk_runtime

# -----------------------------------------------------------------------------
# Stage: builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_BUILDER} AS builder
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/go-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Run as "root" for system installation.

USER root

# Install packages via apt-get.

RUN apt-get update \
 && apt-get -y --no-install-recommends install \
        libsqlite3-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/serve-grpc

# Copy files from prior stage.

COPY --from=senzingsdk_runtime  "/opt/senzing/er/lib/"   "/opt/senzing/er/lib/"
COPY --from=senzingsdk_runtime  "/opt/senzing/er/sdk/c/" "/opt/senzing/er/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/

# Build go program.

WORKDIR ${GOPATH}/src/serve-grpc
RUN make build-with-libsqlite3

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/serve-grpc/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS final
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/serve-grpc" \
      Maintainer="support@senzing.com" \
      Version="0.9.1"
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD ["/app/healthcheck.sh"]
USER root

# Install packages via apt-get.

RUN apt-get update \
 && apt-get -y --no-install-recommends install \
        libsqlite3-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/serve-grpc /app/serve-grpc

# Run as non-root container

USER 1001

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/
ENV SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@nowhere/IN_MEMORY_DB?mode=memory&cache=shared
ENV SENZING_TOOLS_ENABLE_ALL=true

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/serve-grpc"]
