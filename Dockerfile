# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.20.5@sha256:6b3fa4b908676231b50acbbc00e84d8cee9c6ce072b1175c0ff352c57d8a612f
ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.5.3

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2023-06-15
LABEL Name="senzing/serve-grpc-builder" \
      Maintainer="support@senzing.com" \
      Version="0.4.6"

# Build arguments.

ARG PROGRAM_NAME="unknown"
ARG BUILD_VERSION=0.0.0
ARG BUILD_ITERATION=0
ARG GO_PACKAGE_NAME="unknown"

# Copy remote files from DockerHub.

COPY --from=senzing/senzingapi-runtime:3.5.3  "/opt/senzing/g2/lib/"   "/opt/senzing/g2/lib/"
COPY --from=senzing/senzingapi-runtime:3.5.3  "/opt/senzing/g2/sdk/c/" "/opt/senzing/g2/sdk/c/"

# Copy local files from the Git repository.

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
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2023-06-15
LABEL Name="senzing/serve-grpc" \
      Maintainer="support@senzing.com" \
      Version="0.4.6"

# Copy files from repository.

COPY ./testdata/senzing-license/g2.lic /etc/opt/senzing/g2.lic
COPY ./testdata/sqlite/G2C.db          /tmp/sqlite/G2C.db

# Copy files from prior step.

COPY --from=go_builder "/output/linux/serve-grpc" "/app/serve-grpc"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/
ENV SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/serve-grpc"]
