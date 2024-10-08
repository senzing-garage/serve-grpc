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

`serve-grpc` is a command in the [senzing-tools] suite of tools.
This command is a [gRPC] server application that supports requests to the Senzing SDK via network access.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

`serve-grpc` supports the [Senzing Protocol Buffer definitions].
Under the covers, the gRPC request is translated by the gRPC server into a Senzing Go SDK API call using [senzing/sz-sdk-go-core].
The response from the Senzing Go SDK API is returned to the gRPC client.

Senzing SDKs for accessing the gRPC server:

1. Go: [sz-sdk-go-grpc]
1. Python: [sz-sdk-python-grpc]

A simple demonstration using `senzing-tools` and a SQLite database.

```console
export LD_LIBRARY_PATH=/opt/senzing/er/lib/
export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@nowhere/tmp/sqlite/G2C.db
senzing-tools init-database
senzing-tools serve-grpc --enable-all

```

## Install

1. The `serve-grpc` command is installed with the [senzing-tools] suite of tools.
   See [senzing-tools install](https://github.com/senzing-garage/senzing-tools#install).

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/er/lib/
senzing-tools serve-grpc [flags]
```

1. For options and flags:
    1. [Online documentation](https://hub.senzing.com/senzing-tools/senzing-tools_serve-grpc.html)
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/er/lib/
        senzing-tools serve-grpc --help
        ```

1. In addition to the following simple usage examples, there are additional [Examples](docs/examples.md).

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools serve-grpc \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --enable-all

    ```

1. Run `senzing-tools serve-grpc --help` or see [Parameters](#parameters) for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_ENABLE_ALL=true
    senzing-tools serve-grpc
    ```

1. Run `senzing-tools serve-grpc --help` or see [Parameters](#parameters) for additional parameters.

### Using Docker

This usage shows how to initialze a database with a Docker container.

1. This usage specifies a URL of an external database.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --env SENZING_TOOLS_ENABLE_ALL=true \
        --interactive \
        --publish 8258:8258 \
        --rm \
        --tty \
        senzing/senzing-tools serve-grpc

    ```

1. See [Parameters](#parameters) for additional parameters.

### Parameters

- **[SENZING_TOOLS_DATABASE_URL]**
- **[SENZING_TOOLS_ENABLE_ALL]**
- **[SENZING_TOOLS_ENABLE_G2CONFIG]**
- **[SENZING_TOOLS_ENABLE_G2CONFIGMGR]**
- **[SENZING_TOOLS_ENABLE_G2DIAGNOSTIC]**
- **[SENZING_TOOLS_ENABLE_G2ENGINE]**
- **[SENZING_TOOLS_ENABLE_G2PRODUCT]**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL]**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME]**
- **[SENZING_TOOLS_GRPC_PORT]**
- **[SENZING_TOOLS_LOG_LEVEL]**

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
[gRPC]: https://github.com/senzing-garage/knowledge-base/blob/main/WHATIS/grpc.md
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/serve-grpc/blob/main/LICENSE
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/serve-grpc
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Protocol Buffer definitions]: https://github.com/senzing-garage/sz-sdk-proto
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[SENZING_TOOLS_DATABASE_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url
[SENZING_TOOLS_ENABLE_ALL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_all
[SENZING_TOOLS_ENABLE_G2CONFIG]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2config
[SENZING_TOOLS_ENABLE_G2CONFIGMGR]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2configmgr
[SENZING_TOOLS_ENABLE_G2DIAGNOSTIC]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2diagnostic
[SENZING_TOOLS_ENABLE_G2ENGINE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2engine
[SENZING_TOOLS_ENABLE_G2PRODUCT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2product
[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json
[SENZING_TOOLS_ENGINE_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level
[SENZING_TOOLS_ENGINE_MODULE_NAME]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name
[SENZING_TOOLS_GRPC_PORT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_grpc_port
[SENZING_TOOLS_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level
[senzing-tools]: https://github.com/senzing-garage/senzing-tools
[Senzing]: https://senzing.com/
[senzing/sz-sdk-go-core]: https://github.com/senzing-garage/sz-sdk-go-core
[sz-sdk-go-grpc]: https://github.com/senzing-garage/sz-sdk-go-grpc
[sz-sdk-python-grpc]: https://github.com/senzing-garage/sz-sdk-python-grpc
