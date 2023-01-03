# servegrpc

## :warning: WARNING: servegrpc is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing servegrpc...

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/servegrpc.svg)](https://pkg.go.dev/github.com/senzing/servegrpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/servegrpc)](https://goreportcard.com/report/github.com/senzing/servegrpc)
[![go-test.yaml](https://github.com/Senzing/servegrpc/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/servegrpc/actions/workflows/go-test.yaml)

## Overview

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
servegrpc [flags]
```

### Options

```console
      --enable-g2config       enable G2Config service [SENZING_TOOLS_ENABLE_G2CONFIG]
      --enable-g2configmgr    enable G2ConfigMgr service [SENZING_TOOLS_ENABLE_G2CONFIGMGR]
      --enable-g2diagnostic   enable G2Diagnostic service [SENZING_TOOLS_ENABLE_G2DIAGNOSTIC]
      --enable-g2engine       enable G2Config service [SENZING_TOOLS_ENABLE_G2ENGINE]
      --enable-g2product      enable G2Config service [SENZING_TOOLS_ENABLE_G2PRODUCT]
      --grpc-port int         port used to serve gRPC [SENZING_TOOLS_GRPC_PORT] (default 8258)
  -h, --help                  help for servegrpc
      --log-level string      log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [SENZING_TOOLS_LOG_LEVEL] (default "INFO")
```

## Development

### Build

1. Verify Senzing C SDK header files and shared objects are installed.
    1. `/opt/senzing/g2/lib`
    1. `/opt/senzing/g2/sdk/c`

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=servegrpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Build the binaries.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make build

     ```

1. The binaries will be found in ${GIT_REPOSITORY_DIR}/target.
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

### Test

1. Install the  [bloomrpc](https://github.com/bloomrpc/bloomrpc) gRPC test client.
   1. Example for Ubuntu.

       1. Find [latest release](https://github.com/bloomrpc/bloomrpc/releases).

       1. :pencil2: Install.
          Example:

           ```console
           wget https://github.com/bloomrpc/bloomrpc/releases/download/1.5.3/bloomrpc_1.5.3_amd64.deb
           sudo apt install ./bloomrpc_1.5.3_amd64.deb

           ```

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=servegrpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Create a test database.
   Example:

     ```console
     mkdir /tmp/sqlite
     cp ${GIT_REPOSITORY_DIR}/testdata/sqlite/G2C.db /tmp/sqlite/G2C.db

     ```

1. Start the test server.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make test-servegrpc

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

## Package

### Package RPM and DEB files

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

### Test DEB package on Ubuntu

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

1. Run command.
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

### Test RPM package on Centos

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

Make documents that are visible at
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
