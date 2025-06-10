# serve-grpc examples

## Command line examples

### Command line example - Enable only szengine gRPC service

For security reasons, it may be that only certain gRPC services are started.
In this example, only the SzEngine gRPC is started.

1. Using command line options.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools serve-grpc \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --enable-szengine
    ```

### Command line example - using SENZING_TOOLS_ENGINE_CONFIGURATION_JSON environment variable

If using multiple databases or non-system locations of Senzing binaries,
`SENZING_TOOLS_ENGINE_CONFIGURATION_JSON` is used to configure the Senzing runtime engine.

1. :pencil2: Set the value of `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON`.
   Example:

    ```console
    export SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='{
        "PIPELINE": {
            "CONFIGPATH": "/etc/opt/senzing",
            "RESOURCEPATH": "/opt/senzing/er/resources",
            "SUPPORTPATH": "/opt/senzing/data"
        },
        "SQL": {
            "CONNECTION": "postgresql://username:password@host.example.com:G2/"
        }
    }'
    ```

1. Run the gRPC server.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools serve-grpc
    ```

1. For more information, visit
   [SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)

## Docker examples

### Docker example - help

1. Show help to list environment variables that can be used in `docker run`'s, `--env` parameter.
   Example:

    ```console
    docker run --rm senzing/serve-grpc --help
    ```

### Docker example - Using internal, transient SQLite database

1. This usage has an SQLite database that is baked into the Docker container.
   The container is mutable and the data in the database is lost when the container is terminated.
   Only use this technique for simple tests.
   Example:

    ```console
    docker run \
        --interactive \
        --publish 8261:8261 \
        --rm \
        --tty \
        senzing/serve-grpc
    ```

### Docker example - Using external SQLite database

1. This usage has an SQLite database that is baked into the Docker container.
   The container is mutable and the data in the database is lost when the container is terminated.
   Only use this technique for simple tests.
   Example:

   Specify a directory to store the database.

    ```console
    export MY_SENZING_DIRECTORY=~/my-senzing
    ```

   Create the directory and run the Docker container.

    ```console
    mkdir ${MY_SENZING_DIRECTORY}

    docker run \
        --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@nowhere/senzing/G2C.db \
        --interactive \
        --publish 8261:8261 \
        --rm \
        --tty \
        --volume ${MY_SENZING_DIRECTORY}:/senzing \
        senzing/serve-grpc
    ```
