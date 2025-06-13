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
   Use this technique for simple tests only.
   Example:

    ```console
    docker run \
        --interactive \
        --publish 8261:8261 \
        --rm \
        --tty \
        senzing/serve-grpc
    ```

### Docker example - Using postgres Docker container

1. Create a Docker network.

    ```console
    docker network create my-senzing-network --driver bridge
    ```

1. Bring up a PostgreSQL database using the [postgresql] Docker image.

    ```console
    docker run --env POSTGRES_DB=G2 --env POSTGRES_PASSWORD=my-password --name my-postgres --network my-senzing-network --rm postgres
    ```

    This example does not persist data after the Docker container is terminated.
    For techniques on persisting data, see [postgresql].

1. Using a separate terminal, populate the database with the Senzing schema and configuration.

    ```console
    docker run --env SENZING_TOOLS_DATABASE_URL=postgresql://postgres:my-password@my-postgres:5432/G2/?sslmode=disable --network my-senzing-network --rm senzing/senzing-tools init-database
    ```

1. Run the gRPC server using the Postgres database.

    ```console
    docker run --env SENZING_TOOLS_DATABASE_URL=postgresql://postgres:my-password@my-postgres:5432/G2/?sslmode=disable --name my-grpc-server --network my-senzing-network --publish 8261:8261 --rm senzing/serve-grpc
    ```

   The gRPC service is available on port 8261.

1. When the gRPC server is no longer needed, here's how to clean up.

    ```console
    docker kill my-postgres
    docker kill my-grpc-server
    docker network rm my-senzing-network
    ```

### Docker example - Using bitnami/postgresql Docker container

1. Create a Docker network.

    ```console
    docker network create my-senzing-network --driver bridge
    ```

1. Bring up a PostgreSQL database using the [bitnami/postgresql] Docker image.

    ```console
    docker run \
        --env POSTGRESQL_DATABASE=G2 \
        --env POSTGRESQL_PASSWORD=my-password \
        --interactive \
        --name my-postgresql \
        --network my-senzing-network \
        --publish 5432:5432 \
        --rm \
        --tty \
        bitnami/postgresql:latest
    ```

    This example does not persist data after the Docker container is terminated.
    For techniques on persisting data, see [bitnami/postgresql].

1. Using a separate terminal, populate the database with the Senzing schema and configuration.

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://postgres:my-password@my-postgresql:5432/G2/?sslmode=disable \
        --network my-senzing-network \
        --rm \
        senzing/senzing-tools init-database
    ```

1. Run the gRPC server with the external Postgres database.

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://postgres:my-password@my-postgresql:5432/G2/?sslmode=disable \
        --interactive \
        --name my-grpc-server \
        --network my-senzing-network \
        --publish 8261:8261 \
        --rm \
        --tty \
        senzing/serve-grpc
    ```

   The gRPC service is available on port 8261.

1. When the gRPC server is no longer needed, here's how to clean up.

    ```console
    docker kill my-postgresql
    docker kill my-grpc-server
    docker network rm my-senzing-network
    ```

### Docker example - Using external SQLite database

:no_entry: This technique is not recommended.
It crashes on macOS and Windows and is unstable in Linux.

1. This usage creates an SQLite database that is outside the Docker container.
   The SQLite database may be reused across multiple `docker run` commands.
   Use this technique for simple tests only.
   Example:

   :pencil2: Specify a directory to store the database.

    ```console
    export MY_SENZING_DIRECTORY=~/my-senzing
    ```

   Create an empty SQLite database and populate it with the Senzing schema and configuration.
   Remember `SENZING_TOOLS_DATABASE_URL` references the SQLite file *inside* the Docker container.
   Example:

    ```console
    mkdir -p ${MY_SENZING_DIRECTORY}
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@nowhere/tmp/sqlite/G2C.db \
        --rm \
        --user $(id -u):$(id -g) \
        --volume ${MY_SENZING_DIRECTORY}:/tmp/sqlite \
        senzing/senzing-tools init-database
    ```

   Run the gRPC server with the SQLite database mounted.

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@nowhere/tmp/sqlite/G2C.db \
        --interactive \
        --publish 8261:8261 \
        --rm \
        --tty \
        --user $(id -u):$(id -g) \
        --volume ${MY_SENZING_DIRECTORY}:/tmp/sqlite \
        senzing/serve-grpc
    ```

[bitnami/postgresql]: https://hub.docker.com/r/bitnami/postgresql
[postgresql]: https://hub.docker.com/_/postgres
