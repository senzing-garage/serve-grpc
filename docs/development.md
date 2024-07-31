# serve-grpc development

## Install Go

1. See Go's [Download and install].

## Install Senzing C library

Since the Senzing library is a prerequisite, it must be installed first.

1. Verify Senzing C shared objects, configuration, and SDK header files are installed.
    1. `/opt/senzing/g2/lib`
    1. `/opt/senzing/g2/sdk/c`
    1. `/etc/opt/senzing`

1. If not installed, see [How to Install Senzing for Go Development].

## Install Git repository

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing-garage
    export GIT_REPOSITORY=serve-grpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow
   steps in [clone-repository] to install the Git repository.

## Dependencies

1. A one-time command to install dependencies needed for `make` targets.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make dependencies-for-development

    ```

1. Install dependencies needed for [Go] code.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make dependencies

    ```

## Build

1. Build the binaries.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean build

    ```

1. The binaries will be found in the `${GIT_REPOSITORY_DIR}/target` directory.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

## Run

After running `make build`,
the binary built can be run.

1. Identify the database by setting the `SENZING_TOOLS_DATABASE_URL` environment variable.
    1. If using Sqlite database
        1. :pencil2: Create a new Sqlite database.
        Examples:

            ```console
            export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db
            export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
            senzing-tools init-database

            ```

    1. If using PostgreSQL database
        1. :pencil2: Identify an existing PostgreSQL database.
           Example:

            ```console
            export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/senzing-garage/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
            export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/G2/?sslmode=disable

            ```

        1. If needed, initialize PostgreSQL database.
           Example:

            ```console
            export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
            senzing-tools init-database

            ```

1. Set `LD_LIBRARY_PATH` and run the command.
   Examples:
    1. Linux

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        ${GIT_REPOSITORY_DIR}/target/linux-amd64/serve-grpc

        ```

    1. macOS

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        ${GIT_REPOSITORY_DIR}/target/darwin-amd64/serve-grpc

        ```

    1. Windows

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        ${GIT_REPOSITORY_DIR}/target/windows-amd64/serve-grpc

        ```

1. Clean up.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean

    ```

## Lint

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make lint

    ```

## Test

### Test using SQLite database

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup test

    ```

1. **Optional:** View the SQLite database.
   Example:

    ```console
    docker run \
        --env SQLITE_DATABASE=G2C.db \
        --interactive \
        --publish 9174:8080 \
        --rm \
        --tty \
        --volume /tmp/sqlite:/data \
        coleifer/sqlite-web

    ```

   Visit [localhost:9174].

### Test using Docker-compose stack with PostgreSql database

The following instructions show how to bring up a test stack to be used
in testing the `sz-sdk-go-core` packages.

1. Identify a directory to place docker-compose artifacts.
   The directory specified will be deleted and re-created.
   Example:

    ```console
    export SENZING_DEMO_DIR=~/my-senzing-demo

    ```

1. Bring up the docker-compose stack.
   Example:

    ```console
    export PGADMIN_DIR=${SENZING_DEMO_DIR}/pgadmin
    export POSTGRES_DIR=${SENZING_DEMO_DIR}/postgres
    export RABBITMQ_DIR=${SENZING_DEMO_DIR}/rabbitmq
    export SENZING_VAR_DIR=${SENZING_DEMO_DIR}/var
    export SENZING_UID=$(id -u)
    export SENZING_GID=$(id -g)

    rm -rf ${SENZING_DEMO_DIR:-/tmp/nowhere/for/safety}
    mkdir ${SENZING_DEMO_DIR}
    mkdir -p ${PGADMIN_DIR} ${POSTGRES_DIR} ${RABBITMQ_DIR} ${SENZING_VAR_DIR}
    chmod -R 777 ${SENZING_DEMO_DIR}

    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-versions-stable.sh \
        https://raw.githubusercontent.com/senzing-garage/knowledge-base/main/lists/docker-versions-stable.sh
    source ${SENZING_DEMO_DIR}/docker-versions-stable.sh
    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-compose.yaml \
        "https://raw.githubusercontent.com/senzing-garage/docker-compose-demo/main/resources/postgresql/docker-compose-postgresql.yaml"

    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose up

    ```

1. In a separate terminal window, set environment variables.
   Identify Database URL of database in docker-compose stack.
   Example:

    ```console
    export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/senzing-garage/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/G2/?sslmode=disable

    ```

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

    ```

1. **Optional:** View the PostgreSQL database.

   Visit [localhost:9171].
   For the initial login, review the instructions at the top of the web page.
   For server password information, see the `POSTGRESQL_POSTGRES_PASSWORD` value in `${SENZING_DEMO_DIR}/docker-compose.yaml`.
   Usually, it's "postgres".

1. Cleanup.

    ```console
    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose down

    cd ${GIT_REPOSITORY_DIR}
    make clean

    ```

### Test using bloomrpc

Using a (deprecated) BloomRPC client, test the Senzing gRPC Server.
For other gRPC tools, visit [Awesome gRPC].

1. Install the [bloomrpc] gRPC test client.
   1. Example for Ubuntu.

       1. Find [latest release].

       1. :pencil2: Install.
          Example:

           ```console
           wget https://github.com/bloomrpc/bloomrpc/releases/download/1.5.3/bloomrpc_1.5.3_amd64.deb
           sudo apt install ./bloomrpc_1.5.3_amd64.deb

           ```

1. Start the test server.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean run-serve-grpc

    ```

1. In a separate terminal, start the gRPC test client.
   Example:

    ```console
    bloomrpc

    ```

1. In `bloomrpc`:
    1. Near top-center, use the address of `0.0.0.0:8258` to reach the local gRPC server.
    1. In upper-left, click on plus sign ("+").
        1. Navigate to the ${GIT_REPOSITORY_DIR}/proto directory
        1. Choose one or more `.proto` files.
    1. In left-hand pane,
        1. Choose the `Init` message.
        1. Set the request values.
           Example:

            ```json
            {
              "moduleName": "Test of gRPC",
              "iniParams": "{\"PIPELINE\":{\"CONFIGPATH\":\"/etc/opt/senzing\",\"RESOURCEPATH\":\"/opt/senzing/g2/resources\",\"SUPPORTPATH\":\"/opt/senzing/data\"},\"SQL\":{\"CONNECTION\":\"sqlite3://na:na@/tmp/sqlite/G2C.db\"}}",
              "verboseLogging": 0
            }
            ```

        1. Near the center, click the green "play" button.
    1. The Senzing object is initialized and other messages can be tried.

## Coverage

Create a code coverage map.

1. Run Go tests.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean setup coverage

    ```

   A web-browser will show the results of the coverage.
   The goal is to have over 80% coverage.
   Anything less needs to be reflected in [testcoverage.yaml].

## Documentation

1. Start [godoc] documentation server.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean documentation

    ```

1. If a web page doesn't appear, visit [localhost:6060].
1. Senzing documentation will be in the "Third party" section.
   `github.com` > `senzing` > `go-cmdhelping`

1. When a versioned release is published with a `v0.0.0` format tag,
the reference can be found by clicking on the following badge at the top of the README.md page.
Example:

    [![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/template-go.svg)](https://pkg.go.dev/github.com/senzing-garage/template-go)

1. To stop the `godoc` server, run

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean

    ```

## Docker

1. Use make target to run a docker images that builds RPM and DEB files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make docker-build

    ```

1. Run docker container.
   Example:

    ```console
    docker run --rm senzing/serve-grpc

    ```

## Package

**Note:** This only packages the `serve-grpc` command.
It is only to be used in development and test.
The actual packaging is done in the [senzing-tools] repository.

### Package RPM and DEB files

1. Use make target to run a docker images that builds RPM and DEB files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make package

    ```

1. The results will be in the `${GIT_REPOSITORY_DIR}/target` directory.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

### Test DEB package on Ubuntu

1. Determine if `serve-grpc` is installed.
   Example:

    ```console
    apt list --installed | grep serve-grpc

    ```

1. :pencil2: Install `serve-grpc`.
   The `serve-grpc-...` filename will need modification.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo apt install ./serve-grpc-0.0.0.deb

    ```

1. :pencil2: Identify database.
   One option is to bring up PostgreSql as see in [Test using Docker-compose stack with PostgreSql database].
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db

    ```

1. :pencil2: Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    serve-grpc

    ```

1. Remove `serve-grpc` from system.
   Example:

    ```console
    sudo apt-get remove serve-grpc

    ```

### Test RPM package on Centos

1. Determine if `serve-grpc` is installed.
   Example:

    ```console
    yum list installed | grep serve-grpc

    ```

1. :pencil2: Install `serve-grpc`.
   The `serve-grpc-...` filename will need modification.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo yum install ./serve-grpc-0.0.0.rpm

    ```

1. :pencil2: Identify database.
   One option is to bring up PostgreSql as see in [Test using Docker-compose stack with PostgreSql database].
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db

    ```

1. Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    serve-grpc

    ```

1. Remove `serve-grpc` from system.
   Example:

    ```console
    sudo yum remove serve-grpc

    ```

[clone-repository]: https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/clone-repository.md
[Download and install]: https://go.dev/doc/install
[Go]: https://go.dev/
[godoc]: https://pkg.go.dev/golang.org/x/tools/cmd/godoc
[How to Install Senzing for Go Development]: https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md
[localhost:6060]: http://localhost:6060/pkg/github.com/senzing-garage/template-go/
[testcoverage.yaml]: ../.github/coverage/testcoverage.yaml
[localhost:9174]: http://localhost:9174
[localhost:9171]: http://localhost:9171
[Awesome gRPC]: https://github.com/grpc-ecosystem/awesome-grpc#tools
[bloomrpc]: https://github.com/bloomrpc/bloomrpc
[latest release]: https://github.com/bloomrpc/bloomrpc/releases
[senzing-tools]: https://github.com/senzing-garage/senzing-tools
[Test using Docker-compose stack with PostgreSql database]: #test-using-docker-compose-stack-with-postgresql-database
