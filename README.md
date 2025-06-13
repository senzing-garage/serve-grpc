# serve-grpc

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: serve-grpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

`serve-grpc` supports the [Senzing Protocol Buffer definitions].
Under the covers, the gRPC request is translated by the gRPC server into a Senzing Go SDK call using [senzing/sz-sdk-go-core].
The response from the Senzing Go SDK is returned to the gRPC client.

Senzing SDKs for accessing the gRPC server:

1. Python: [sz-sdk-python-grpc]
1. Java: [sz-sdk-java-grpc]
1. Go: [sz-sdk-go-grpc]

## Install

## Use

1. Docker container with internal, ephemeral database and
   gRPC accessible on port 8261.
   Example:

    ```console
    docker run -it -p 8261:8261 --rm senzing/serve-grpc
    ```

1. See [Parameters](#parameters) for additional parameters.

### Parameters

- **[SENZING_TOOLS_AVOID_SERVING]**
- **[SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE]**
- **[SENZING_TOOLS_CONFIG_PATH]**
- **[SENZING_TOOLS_DATABASE_URL]**
- **[SENZING_TOOLS_ENABLE_ALL]**
- **[SENZING_TOOLS_ENABLE_HTTP]**
- **[SENZING_TOOLS_ENABLE_SZCONFIG]**
- **[SENZING_TOOLS_ENABLE_SZCONFIGMANAGER]**
- **[SENZING_TOOLS_ENABLE_SZDIAGNOSTIC]**
- **[SENZING_TOOLS_ENABLE_SZENGINE]**
- **[SENZING_TOOLS_ENABLE_SZPRODUCT]**
- **[SENZING_TOOLS_ENGINE_INSTANCE_NAME]**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL]**
- **[SENZING_TOOLS_GRPC_PORT]**
- **[SENZING_TOOLS_HTTP_PORT]**
- **[SENZING_TOOLS_LICENSE_STRING_BASE64]**
- **[SENZING_TOOLS_LOG_LEVEL]**
- **[SENZING_TOOLS_OBSERVER_ORIGIN]**
- **[SENZING_TOOLS_OBSERVER_URL]**
- **[SENZING_TOOLS_RESOURCE_PATH]**
- **[SENZING_TOOLS_SENZING_DIRECTORY]**
- **[SENZING_TOOLS_SERVER_ADDRESS]**
- **[SENZING_TOOLS_SERVER_CERTIFICATE_FILE]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_MIN_TIME_IN_SECONDS]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_PERMIT_WITHOUT_STREAM]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_GRACE_IN_SECONDS]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_IN_SECONDS]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_IDLE_IN_SECONDS]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIME_IN_SECONDS]**
- **[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIMEOUT_IN_SECONDS]**
- **[SENZING_TOOLS_SERVER_KEY_FILE]**
- **[SENZING_TOOLS_SERVER_KEY_PASSPHRASE]**
- **[SENZING_TOOLS_SERVER_MAX_CONCURRENT_STREAMS]**
- **[SENZING_TOOLS_SERVER_MAX_HEADER_LIST_SIZE_IN_BYTES]**
- **[SENZING_TOOLS_SERVER_MAX_RECEIVE_MESSAGE_SIZE_IN_BYTES]**
- **[SENZING_TOOLS_SERVER_MAX_SEND_MESSAGE_SIZE_IN_BYTES]**
- **[SENZING_TOOLS_SERVER_READ_BUFFER_SIZE_IN_BYTES]**
- **[SENZING_TOOLS_SERVER_WRITE_BUFFER_SIZE_IN_BYTES]**
- **[SENZING_TOOLS_SUPPORT_PATH]**

## References

1. [Command reference]
1. [Development]
1. [Errors]
1. [Examples]

[Command reference]: https://hub.senzing.com/senzing-tools/senzing-tools_serve-grpc.html
[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/serve-grpc.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/serve-grpc
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/serve-grpc
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/serve-grpc/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/serve-grpc/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/serve-grpc/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/serve-grpc/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/serve-grpc/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/serve-grpc/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/serve-grpc/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/serve-grpc/actions/workflows/golangci-lint.yaml
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/serve-grpc/blob/main/LICENSE
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/serve-grpc
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Protocol Buffer definitions]: https://github.com/senzing-garage/sz-sdk-proto
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[SENZING_TOOLS_AVOID_SERVING]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_avoid_serving
[SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_client_ca_certificate_file
[SENZING_TOOLS_CONFIG_PATH]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_config_path
[SENZING_TOOLS_DATABASE_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url
[SENZING_TOOLS_ENABLE_ALL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_all
[SENZING_TOOLS_ENABLE_HTTP]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_http
[SENZING_TOOLS_ENABLE_SZCONFIG]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szconfig
[SENZING_TOOLS_ENABLE_SZCONFIGMANAGER]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szconfigmanager
[SENZING_TOOLS_ENABLE_SZDIAGNOSTIC]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szdiagnostic
[SENZING_TOOLS_ENABLE_SZENGINE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szengine
[SENZING_TOOLS_ENABLE_SZPRODUCT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szproduct
[SENZING_TOOLS_ENGINE_INSTANCE_NAME]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_instance_name
[SENZING_TOOLS_ENGINE_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level
[SENZING_TOOLS_GRPC_PORT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_grpc_port
[SENZING_TOOLS_HTTP_PORT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_http_port
[SENZING_TOOLS_LICENSE_STRING_BASE64]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_license_string_base64
[SENZING_TOOLS_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level
[SENZING_TOOLS_OBSERVER_ORIGIN]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_observer_origin
[SENZING_TOOLS_OBSERVER_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_observer_url
[SENZING_TOOLS_RESOURCE_PATH]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_resource_path
[SENZING_TOOLS_SENZING_DIRECTORY]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_senzing_directory
[SENZING_TOOLS_SERVER_ADDRESS]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_address
[SENZING_TOOLS_SERVER_CERTIFICATE_FILE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_certificate_file
[SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_MIN_TIME_IN_SECONDS]:https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_enforcement_policy_min_time_in_seconds
[SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_PERMIT_WITHOUT_STREAM]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_enforcement_policy_permit_without_stream
[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_GRACE_IN_SECONDS]:https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_server_parameter_max_connection_age_grace_in_seconds
[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_IN_SECONDS]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_server_parameter_max_connection_age_in_seconds
[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_IDLE_IN_SECONDS]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_server_parameter_max_connection_idle_in_seconds
[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIME_IN_SECONDS]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_server_parameter_time_in_seconds
[SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIMEOUT_IN_SECONDS]:https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_keepalive_server_parameter_timeout_in_seconds
[SENZING_TOOLS_SERVER_KEY_FILE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_key_file
[SENZING_TOOLS_SERVER_KEY_PASSPHRASE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_key_passphrase
[SENZING_TOOLS_SERVER_MAX_CONCURRENT_STREAMS]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_max_concurrent_streams
[SENZING_TOOLS_SERVER_MAX_HEADER_LIST_SIZE_IN_BYTES]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_max_header_list_size_in_bytes
[SENZING_TOOLS_SERVER_MAX_RECEIVE_MESSAGE_SIZE_IN_BYTES]:https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_max_receive_message_size_in_bytes
[SENZING_TOOLS_SERVER_MAX_SEND_MESSAGE_SIZE_IN_BYTES]:https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_max_send_message_size_in_bytes
[SENZING_TOOLS_SERVER_READ_BUFFER_SIZE_IN_BYTES]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_read_buffer_size_in_bytes
[SENZING_TOOLS_SERVER_WRITE_BUFFER_SIZE_IN_BYTES]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_server_write_buffer_size_in_bytes
[SENZING_TOOLS_SUPPORT_PATH]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_support_path
[Senzing]: https://senzing.com/
[senzing/sz-sdk-go-core]: https://github.com/senzing-garage/sz-sdk-go-core
[sz-sdk-go-grpc]: https://github.com/senzing-garage/sz-sdk-go-grpc
[sz-sdk-java-grpc]: https://github.com/senzing-garage/sz-sdk-java-grpc
[sz-sdk-python-grpc]: https://github.com/senzing-garage/sz-sdk-python-grpc
