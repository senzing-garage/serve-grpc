# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

# FIXME: ARG IMAGE_SENZINGAPI_RUNTIME=senzing/senzingapi-runtime:3.7.1
ARG IMAGE_SENZINGAPI_RUNTIME=senzing/senzingapi-runtime:staging

ARG IMAGE_GO_BUILDER=golang:1.21.0-bullseye@sha256:02f350d8452d3f9693a450586659ecdc6e40e9be8f8dfc6d402300d87223fdfa

# FIXME: ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.7.1
ARG IMAGE_FINAL=senzing/senzingapi-runtime:staging

# -----------------------------------------------------------------------------
# Stage: senzingapi_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_SENZINGAPI_RUNTIME} as senzingapi_runtime

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2023-09-29
LABEL Name="senzing/serve-grpc-builder" \
      Maintainer="support@senzing.com" \
      Version="0.5.1"

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
ENV REFRESHED_AT=2023-09-29
LABEL Name="senzing/serve-grpc" \
      Maintainer="support@senzing.com" \
      Version="0.5.1"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY ./testdata/senzing-license/g2.lic /etc/opt/senzing/g2.lic
COPY ./testdata/sqlite/G2C.db          /tmp/sqlite/G2C.db

# Copy files from prior stage.

COPY --from=go_builder "/output/linux-amd64/serve-grpc" "/app/serve-grpc"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/
ENV SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db

# Runtime execution.

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/app/healthcheck.sh" ]

WORKDIR /app
ENTRYPOINT ["/app/serve-grpc"]
