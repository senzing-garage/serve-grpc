# servegrpc

## :warning: WARNING: servegrpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`servegrpc` is a subcommand of the
[senzing-tools](https://github.com/Senzing/senzing-tools)
suite of tools.
This subcommand is a
[gRPC](https://grpc.io/)
server application that supports requests to the Senzing SDK via network access.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/servegrpc.svg)](https://pkg.go.dev/github.com/senzing/servegrpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/servegrpc)](https://goreportcard.com/report/github.com/senzing/servegrpc)

## Overview

`servegrpc` supports the
[Senzing Protocol Buffer definitions](https://github.com/Senzing/g2-sdk-proto).
Under the covers, the gRPC request is translated into a Senzing Go SDK API call using
[senzing/g2-sdk-go-base](https://github.com/Senzing/g2-sdk-go-base).
The response from the Senzing Go SDK API is returned to the gRPC client.

Senzing SDKs for accessing the gRPC server:

1. Go: [g2-sdk-go-grpc](https://github.com/Senzing/g2-sdk-go-grpc)

## Install

1. It is installed with the
   [senzing-tools](https://github.com/Senzing/senzing-tools)
   suite of tools.
   See senzing-tools [install](https://github.com/Senzing/senzing-tools#install).

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
senzing-tools servegrpc [flags]
```

1. For options and flags:
    1. Visit
       [hub.senzing.com/senzing-tools/senzing-tools_initdatabase.html](https://hub.senzing.com/senzing-tools/senzing-tools_initdatabase.html)
    1. Run:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        senzing-tools servegrpc --help
        ```

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools servegrpc --database-url postgresql://username:password@postgres.example.com:5432/G2
    ```

1. See [Parameters](#parameters) for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools servegrpc
    ```

1. See [Parameters](#parameters) for additional parameters.

### Using Docker

This usage shows how to initialze a database with a Docker container.

1. This usage has an SQLite database that is baked into the Docker container.
   The data in the database is lost when the container is terminated.
   Run `senzing/senzing-tools` container.
   Example:

    ```console
    docker run \
        --interactive \
        --publish 8258:8258 \
        --rm \
        --tty \
        senzing/senzing-tools servegrpc

    ```

1. This usage accepts a URL of an external database.

    1. :thinking: Identify the database URL.
       The example may not work in all cases.
       Example:

        ```console
        export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/Senzing/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
        export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/G2

        ```

    1. Run `senzing/servegrpc`.
       Example:

        ```console
        docker run \
            --env SENZING_TOOLS_DATABASE_URL \
            --interactive \
            --publish 8258:8258 \
            --rm \
            --tty \
        senzing/senzing-tools servegrpc

        ```

1. See [Parameters](#parameters) for additional parameters.

### Parameters

- **[SENZING_TOOLS_DATABASE_URL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url)**
- **[SENZING_TOOLS_ENABLE_G2CONFIG](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2config)**
- **[SENZING_TOOLS_ENABLE_G2CONFIGMGR](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2configmgr)**
- **[SENZING_TOOLS_ENABLE_G2DIAGNOSTIC](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2diagnostic)**
- **[SENZING_TOOLS_ENABLE_G2ENGINE](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2engine)**
- **[SENZING_TOOLS_ENABLE_G2PRODUCT](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2product)**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level)**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name)**
- **[SENZING_TOOLS_GRPC_PORT](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_grpc_port)**
- **[SENZING_TOOLS_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level)**

## References

- [Command reference](https://hub.senzing.com/senzing-tools/senzing-tools_servegrpc.html)
- [Errors](docs/errors.md)
- [Examples](docs/examples.md)
- [Development](docs/development.md)
