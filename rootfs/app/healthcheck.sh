#!/usr/bin/env bash

# Determine port.

DEFAULT_GRPC_PORT=8261

# If SENZING_TOOLS_GRPC_PORT exists, use it's value.

if [[ -z "${SENZING_TOOLS_GRPC_PORT}" ]]; then
    GRPC_PORT=${DEFAULT_GRPC_PORT}
else
    GRPC_PORT=${SENZING_TOOLS_GRPC_PORT}
fi

# Health check.

wget --no-verbose --tries=1 --spider http://localhost:${GRPC_PORT}
exit $?
