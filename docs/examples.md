# serve-grpc examples

## Docker examples

1. [Help]
1. [Using internal, transient SQLite database]
1. [Using Postgres database]
1. [Using custom Senzing license]
1. [Using TLS server-side authentication]
1. [Using TLS mutual authentication]
1. [Using bitnami/postgresql Docker container]

### Docker example - Help

This example shows environment variables and command-line arguments used to modify the behavior of the Senzing gRPC Server.

1. Show help to list environment variables that can be used in `docker run`'s `--env` parameter.
   Example:

    ```console
    docker run --rm senzing/serve-grpc --help
    ```

### Docker example - Using internal, transient SQLite database

This example shows the simplest use of the Senzing gRPC Server.

This usage has an SQLite database that is baked into the Docker container.
The container is mutable and the data in the database is lost when the container is terminated.

:warning: Use this technique for simple tests only.

1. Start the Senzing gRPC Server.
   Example:

    ```console
    docker run -it --publish 8261:8261 --rm senzing/serve-grpc
    ```

1. A quick test using [grpcurl].

    ```console
    grpcurl -plaintext -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

### Docker example - Using Postgres database

This example shows how to use a Postgres database with the Senzing gRPC Server.

The example brings up a Postgres Docker container.
If you already have a Postgres database:
    - Steps #1 and #2 may be skipped
    - The `SENZING_TOOLS_DATABASE_URL` value needs to reference your Postgres database
    - The `--network` argument is no longer needed

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

1. A quick test using [grpcurl].

    ```console
    grpcurl -plaintext -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

1. When the gRPC server is no longer needed, here's how to clean up.

    ```console
    docker kill my-postgres
    docker kill my-grpc-server
    docker network rm my-senzing-network
    ```

### Docker example - Using custom Senzing license

This example shows how to use your Senzing license key with the Senzing gRPC Server.

1. The Senzing engine come with a complementary license.
   To see this license, run
   Example:

    ```console
    docker run -it --publish 8261:8261 --rm senzing/serve-grpc
    ```

   In a separate terminal, view the license using [grpcurl].

    ```console
    grpcurl -plaintext -format text localhost:8261 szproduct.SzProduct.GetLicense
    ```

1. :pencil2: To use your custom license, a command-line argument may be used.
   To see your license, replace the value of `--license-string-base64` with your license key.
   *Note:* The example uses a non-working, expired license.
   Example:

    ```console
    docker run -it --publish 8261:8261 --rm senzing/serve-grpc --license-string-base64 AQAAADgCAAAAAAAAU2VuemluZyBUdXRvcmlhbAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARVZBTCAtIHN1cHBvcnRAc2VuemluZy5jb20AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADIwMjUtMDYtMTcAAAAAAAAAAAAARVZBTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNUQU5EQVJEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQKX8AAAEAAAAAAAAAMjAyNi0wNi0xNwAAAAAAAAAAAABZRUFSTFkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJ4NL/lYgnQcCMaMH6xHcSdR1Tjky7Ontd6s3VHQUmqQhCA5/9JVvd8O/V9qBetVWyDJ5XJKGIc7NgvYDJjd+etW68P6DCh3H9yY3WK+LMlvXWoBcpkEJS/Zx0zwpL/gMjii2GwjsTDUiYyg1xUDs/ScmA08B6gP44A4UQzqzWAGsVngGtcjsEgEyNvcVnEOR96xw3R5GBaagwE8C5RlMfyjijtSqNRo859VRo60zDaPIpLW9WCAD7FPxtVEESHSlnYCWIIp5Yv6Ba1v/i+4dLM6rtiVBC3Ukz5vjYZkj7wf3GdtcpGzEYfAWsaYT2rU75iymO1stw44i8Gguz8D8+xZJoeC8MmoXkmXfIz7ZKiFspUeer8kFcIeRu/qaXzwG10iut537P/YG6UbJh78lGlpWwdiuEf3I7kqaQiigCq/egjj2fB7E871PXitBFvYwhDEnGi6DYNbnOlaWN8QOZAEexES2uwYeFWXnnmreTf52Sxz66jPXNCGYcxgBXDmBRpZa3tUqpYsY+fxu/QyPC7nJIt29Peaa8Cxna4807BesCaSiYOZllSRTcZNkQK6LIzMTgTZwXVVL0uVcIk2GIhpC3C84KjBj8dDnBiyRMFs/ubWmrJQFSkMr2AOFFZwr4w7Zp0TNCoVNJ5zUjQBTDDd4RgNGozYHrHvWr1l3+NL
    ```

   In a separate terminal, view your license using [grpcurl].

    ```console
    grpcurl -plaintext -format text localhost:8261 szproduct.SzProduct.GetLicense
    ```

1. :pencil2: Alternatively, to use your custom license, an environment variable may be used.
   To see your license, replace the value of `SENZING_TOOLS_LICENSE_STRING_BASE64` with your license key.
    *Note:* The example uses a non-working, expired license.
   Example:

    ```console
    export SENZING_TOOLS_LICENSE_STRING_BASE64=AQAAADgCAAAAAAAAU2VuemluZyBUdXRvcmlhbAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARVZBTCAtIHN1cHBvcnRAc2VuemluZy5jb20AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADIwMjUtMDYtMTcAAAAAAAAAAAAARVZBTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFNUQU5EQVJEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQKX8AAAEAAAAAAAAAMjAyNi0wNi0xNwAAAAAAAAAAAABZRUFSTFkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJ4NL/lYgnQcCMaMH6xHcSdR1Tjky7Ontd6s3VHQUmqQhCA5/9JVvd8O/V9qBetVWyDJ5XJKGIc7NgvYDJjd+etW68P6DCh3H9yY3WK+LMlvXWoBcpkEJS/Zx0zwpL/gMjii2GwjsTDUiYyg1xUDs/ScmA08B6gP44A4UQzqzWAGsVngGtcjsEgEyNvcVnEOR96xw3R5GBaagwE8C5RlMfyjijtSqNRo859VRo60zDaPIpLW9WCAD7FPxtVEESHSlnYCWIIp5Yv6Ba1v/i+4dLM6rtiVBC3Ukz5vjYZkj7wf3GdtcpGzEYfAWsaYT2rU75iymO1stw44i8Gguz8D8+xZJoeC8MmoXkmXfIz7ZKiFspUeer8kFcIeRu/qaXzwG10iut537P/YG6UbJh78lGlpWwdiuEf3I7kqaQiigCq/egjj2fB7E871PXitBFvYwhDEnGi6DYNbnOlaWN8QOZAEexES2uwYeFWXnnmreTf52Sxz66jPXNCGYcxgBXDmBRpZa3tUqpYsY+fxu/QyPC7nJIt29Peaa8Cxna4807BesCaSiYOZllSRTcZNkQK6LIzMTgTZwXVVL0uVcIk2GIhpC3C84KjBj8dDnBiyRMFs/ubWmrJQFSkMr2AOFFZwr4w7Zp0TNCoVNJ5zUjQBTDDd4RgNGozYHrHvWr1l3+NL
    ```

    ```console
    docker run -it --env SENZING_TOOLS_LICENSE_STRING_BASE64 --publish 8261:8261 --rm senzing/serve-grpc
    ```

   In a separate terminal, view your license using [grpcurl].

    ```console
    grpcurl -plaintext -format text localhost:8261 szproduct.SzProduct.GetLicense
    ```

### Docker example - Using TLS server-side authentication

This example shows how to enable [server-side authentication] Transport Layer Security (TLS) in the Senzing gRPC Server.

1. :pencil2: To run this example, [git clone] the `senzing/serve-grpc` repository.
   Example:

    ```console
    export MY_SENZING_REPOSITORY=~/serve-grpc
    ```

    ```console
    git clone https://github.com/senzing-garage/serve-grpc.git ${MY_SENZING_REPOSITORY}
    ```

1. Run the Senzing gRPC Server container with `SENZING_TOOLS_SERVER_CERTIFICATE_FILE` and `SENZING_TOOLS_SERVER_KEY_FILE`
   environment variables set and `${MY_SENZING_REPOSITORY}` volume mounted.

    ```console
    docker run --env SENZING_TOOLS_SERVER_CERTIFICATE_FILE=/serve-grpc/testdata/certificates/server/certificate.pem --env SENZING_TOOLS_SERVER_KEY_FILE=/serve-grpc/testdata/certificates/server/private_key.pem --publish 8261:8261 --rm --volume ${MY_SENZING_REPOSITORY}:/serve-grpc senzing/serve-grpc
    ```

1. In a separate terminal, run a failing "no TLS" test using [grpcurl].

    ```console
    grpcurl -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

1. :pencil2: Run a successful test using [grpcurl].

    ```console
    export MY_SENZING_REPOSITORY=~/serve-grpc
    ```

    ```console
    grpcurl -authority www.senzing.com -cacert ${MY_SENZING_REPOSITORY}/testdata/certificates/certificate-authority/certificate.pem -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

### Docker example - Using TLS mutual authentication

This example shows how to enable [mutual authentication] Transport Layer Security (TLS) in the Senzing gRPC Server.

1. :pencil2: To run this example, [git clone] the `senzing/serve-grpc` repository.
   Example:

    ```console
    export MY_SENZING_REPOSITORY=~/serve-grpc
    ```

    ```console
    git clone https://github.com/senzing-garage/serve-grpc.git ${MY_SENZING_REPOSITORY}
    ```

1. Run the Senzing gRPC Server container with `SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE`,
   `SENZING_TOOLS_SERVER_CERTIFICATE_FILE` and `SENZING_TOOLS_SERVER_KEY_FILE` environment variables set
   and `${MY_SENZING_REPOSITORY}` volume mounted.

    ```console
    docker run --env SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE=/serve-grpc/testdata/certificates/certificate-authority/certificate.pem --env SENZING_TOOLS_SERVER_CERTIFICATE_FILE=/serve-grpc/testdata/certificates/server/certificate.pem --env SENZING_TOOLS_SERVER_KEY_FILE=/serve-grpc/testdata/certificates/server/private_key.pem --publish 8261:8261 --rm --volume ${MY_SENZING_REPOSITORY}:/serve-grpc senzing/serve-grpc
    ```

1. In a separate terminal, run a failing "no TLS" test using [grpcurl].

    ```console
    grpcurl -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

1. :pencil2: Run another failing "server-side TLS" test using [grpcurl].

    ```console
    export MY_SENZING_REPOSITORY=~/serve-grpc
    ```

    ```console
    grpcurl -authority www.senzing.com -cacert ${MY_SENZING_REPOSITORY}/testdata/certificates/certificate-authority/certificate.pem -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

1. :pencil2: Run a successful "mutual TLS" test using [grpcurl].

    ```console
    export MY_SENZING_REPOSITORY=~/serve-grpc
    ```

    ```console
    grpcurl -authority www.senzing.com -cacert ${MY_SENZING_REPOSITORY}/testdata/certificates/certificate-authority/certificate.pem -cert ${MY_SENZING_REPOSITORY}/testdata/certificates/client/certificate.pem -format text -key ${MY_SENZING_REPOSITORY}/testdata/certificates/client/private_key.pem localhost:8261 szproduct.SzProduct.GetVersion
    ```

### Docker example - Using bitnami/postgresql Docker container

1. Create a Docker network.

    ```console
    docker network create my-senzing-network --driver bridge
    ```

1. Bring up a PostgreSQL database using the [bitnami/postgresql] Docker image.

    ```console
    docker run -it --env POSTGRESQL_DATABASE=G2 --env POSTGRESQL_PASSWORD=my-password --name my-postgresql --network my-senzing-network --publish 5432:5432 --rm bitnami/postgresql
    ```

    This example does not persist data after the Docker container is terminated.
    For techniques on persisting data, see [bitnami/postgresql].

1. Using a separate terminal, populate the database with the Senzing schema and configuration.

    ```console
    docker run --env SENZING_TOOLS_DATABASE_URL=postgresql://postgres:my-password@my-postgresql:5432/G2/?sslmode=disable --network my-senzing-network --rm senzing/senzing-tools init-database
    ```

1. Run the gRPC server with the external Postgres database.

    ```console
    docker run -it --env SENZING_TOOLS_DATABASE_URL=postgresql://postgres:my-password@my-postgresql:5432/G2/?sslmode=disable --name my-grpc-server --network my-senzing-network --publish 8261:8261 --rm senzing/serve-grpc
    ```

   The gRPC service is available on port 8261.

1. A quick test using [grpcurl].

    ```console
    grpcurl -plaintext -format text localhost:8261 szproduct.SzProduct.GetVersion
    ```

1. When the gRPC server is no longer needed, here's how to clean up.

    ```console
    docker kill my-postgresql
    docker kill my-grpc-server
    docker network rm my-senzing-network
    ```

### Docker example - Using external SQLite database

:no_entry: This technique is not recommended.
It crashes on macOS and Windows and is unstable in Linux.
:no_entry:

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

## Command line examples using senzing-tools

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

[bitnami/postgresql]: https://hub.docker.com/r/bitnami/postgresql
[git clone]: https://git-scm.com/docs/git-clone
[grpcurl]: https://github.com/fullstorydev/grpcurl
[Help]: #docker-example---help
[mutual authentication]: https://en.wikipedia.org/wiki/Transport_Layer_Security#Client-authenticated_TLS_handshake
[postgresql]: https://hub.docker.com/_/postgres
[server-side authentication]: https://en.wikipedia.org/wiki/Transport_Layer_Security#Basic_TLS_handshake
[Using bitnami/postgresql Docker container]: #docker-example---using-bitnamipostgresql-docker-container
[Using custom Senzing license]: #docker-example---using-custom-senzing-license
[Using internal, transient SQLite database]: #docker-example---using-internal-transient-sqlite-database
[Using Postgres database]: #docker-example---using-postgres-database
[Using TLS mutual authentication]: #docker-example---using-tls-mutual-authentication
[Using TLS server-side authentication]: #docker-example---using-tls-server-side-authentication
