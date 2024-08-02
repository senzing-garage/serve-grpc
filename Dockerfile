# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.22.3-bullseye
ARG IMAGE_FINAL=senzing/senzingapi-runtime-staging:latest

# -----------------------------------------------------------------------------
# Stage: senzingapi_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as senzingapi_runtime

# -----------------------------------------------------------------------------
# Stage: builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_BUILDER} as builder
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/go-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Run as "root" for system installation.

USER root

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/serve-grpc

# Copy files from prior stage.

COPY --from=senzingapi_runtime  "/opt/senzing/g2/lib/"   "/opt/senzing/g2/lib/"
COPY --from=senzingapi_runtime  "/opt/senzing/g2/sdk/c/" "/opt/senzing/g2/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Build go program.

WORKDIR ${GOPATH}/src/serve-grpc
RUN make build

# Copy binaries to /output.

RUN mkdir -p /output \
      && cp -R ${GOPATH}/src/serve-grpc/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/serve-grpc" \
      Maintainer="support@senzing.com" \
      Version="0.6.0"
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/app/healthcheck.sh" ]
USER root

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder "/output/linux/serve-grpc" "/app/serve-grpc"

# Run as non-root container

USER 1001
COPY ./testdata/sqlite/G2C.db          /tmp/sqlite/G2C.db

# Runtime environment variables.

ENV SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@nowhere/tmp/sqlite/G2C.db

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/serve-grpc"]
