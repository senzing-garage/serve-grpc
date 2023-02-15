# servegrpc

## :warning: WARNING: servegrpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `servegrpc` is a
[gRPC](https://grpc.io/)
server application that supports requests to the Senzing SDK via network access.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/servegrpc.svg)](https://pkg.go.dev/github.com/senzing/servegrpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/servegrpc)](https://goreportcard.com/report/github.com/senzing/servegrpc)

## Overview

The Senzing `servegrpc` supports the
[Senzing Protocol Buffer definitions](https://github.com/Senzing/g2-sdk-proto).
Under the covers, the gRPC request is translated into a Senzing Go SDK API call using
[senzing/g2-sdk-go-base](https://github.com/Senzing/g2-sdk-go-base).
The response from the Senzing Go SDK API is returned to the gRPC client.

Other implementations of the
[g2-sdk-go](https://github.com/Senzing/g2-sdk-go)
interface include:

- [g2-sdk-go-base](https://github.com/Senzing/g2-sdk-go-base) - for
  calling Senzing SDK APIs natively
- [g2-sdk-go-mock](https://github.com/Senzing/g2-sdk-go-mock) - for
  unit testing calls to the Senzing Go SDK
- [go-sdk-abstract-factory](https://github.com/Senzing/go-sdk-abstract-factory) - An
  [abstract factory pattern](https://en.wikipedia.org/wiki/Abstract_factory_pattern)
  for switching among implementations

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
servegrpc [flags]
```

For options and flags, see
[hub.senzing.com/servegrpc](https://hub.senzing.com/servegrpc/) or run:

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
servegrpc --help
```

## Docker

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
        senzing/servegrpc

## Development

### Install Go

1. See Go's [Download and install](https://go.dev/doc/install)

### Install Senzing C library

Since the Senzing library is a prerequisite, it must be installed first.

1. Verify Senzing C shared objects, configuration, and SDK header files are installed.
    1. `/opt/senzing/g2/lib`
    1. `/opt/senzing/g2/sdk/c`
    1. `/etc/opt/senzing`

1. If not installed, see
   [How to Install Senzing for Go Development](https://github.com/Senzing/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md).

### Install Git repository

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=servegrpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in
   [clone-repository](https://github.com/Senzing/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.

### Build

1. Build the binaries.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make build

    ```

1. The binaries will be found in the `${GIT_REPOSITORY_DIR}/target` directory.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

1. Clean up.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean

    ```

### Test using SQLite database

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

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

   Visit [localhost:9174](http://localhost:9174).

### Test using Docker-compose stack with PostgreSql database

The following instructions show how to bring up a test stack to be used
in testing the `g2-sdk-go-base` packages.

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
        https://raw.githubusercontent.com/Senzing/knowledge-base/main/lists/docker-versions-stable.sh
    source ${SENZING_DEMO_DIR}/docker-versions-stable.sh
    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-compose.yaml \
        "https://raw.githubusercontent.com/Senzing/docker-compose-demo/main/resources/postgresql/docker-compose-postgresql.yaml"

    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose up

    ```

1. In a separate terminal window, set environment variables.
   Identify Database URL of database in docker-compose stack.
   Example:

    ```console
    export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/Senzing/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/G2

    ```

1. Run tests.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make clean test

    ```

1. **Optional:** View the PostgreSQL database.

   Visit [localhost:9171](http://localhost:9171).
   For the initial login, review the instructions at the top of the web page.
   For server password information, see the `POSTGRESQL_POSTGRES_PASSWORD` value in `${SENZING_DEMO_DIR}/docker-compose.yaml`.
   Usually, it's "postgres".

1. Cleanup.

    ```console
    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose down

    ```

### Test using bloomrpc

Using a (deprecated) BloomRPC client, test the Senzing gRPC Server.
For other gRPC tools, visit
[Awesome gRPC](https://github.com/grpc-ecosystem/awesome-grpc#tools).

1. Install the  [bloomrpc](https://github.com/bloomrpc/bloomrpc) gRPC test client.
   1. Example for Ubuntu.

       1. Find [latest release](https://github.com/bloomrpc/bloomrpc/releases).

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
    make clean run-servegrpc

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

### Package

#### Package RPM and DEB files

1. :warning: **FIXME:**
   This won't work automatically until
   `/opt/senzing/g2/sdk/c/*.h` and `/opt/senzing/g2/lib/`
   files can be copied from an existing Docker image.
1. :thinking: *Work-around:*
   Copy files from `/opt/senzing/g2/lib` into the repository.
   Example:

    ```console
    cp /opt/senzing/g2/lib/* ${GIT_REPOSITORY_DIR}/rootfs/opt/senzing/g2/lib/

    ```

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

#### Test DEB package on Ubuntu

1. Determine if `servegrpc` is installed.
   Example:

    ```console
    apt list --installed | grep servegrpc

    ```

1. :pencil2: Install `servegrpc`.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo apt install ./servegrpc-0.0.0.deb

    ```

1. :pencil2: Identify database.
   One option is to bring up PostgreSql as see in
   [Test using Docker-compose stack with PostgreSql database](#test-using-docker-compose-stack-with-postgresql-database).
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db

    ```

1. :pencil2: Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    servegrpc
    ```

1. Remove `servegrpc` from system.
   Example:

    ```console
    sudo apt-get remove servegrpc

    ```

#### Test RPM package on Centos

1. Determine if `servegrpc` is installed.
   Example:

    ```console
    yum list installed | grep servegrpc

    ```

1. :pencil2: Install `servegrpc`.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo yum install ./servegrpc-0.0.0.rpm

    ```

1. Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    servegrpc

    ```

1. Remove `servegrpc` from system.
   Example:

    ```console
    sudo yum remove servegrpc

    ```

### Make documents

Make documents visible at
[hub.senzing.com/servegrpc](https://hub.senzing.com/servegrpc).

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=servegrpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Make documents.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    servegrpc docs --dir ${GIT_REPOSITORY_DIR}/docs

    ```

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6011` - g2config
1. `6012` - g2configmgr
1. `6013` - g2diagnostic
1. `6014` - g2engine
1. `6015` - g2hasher
1. `6016` - g2product
1. `6017` - g2ssadm
