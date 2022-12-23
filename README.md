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

```
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
servegrpc [flags]
```

### Options

```
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

### Create protobuf directories

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=servegrpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Identify Senzing subcomponents.
   Example:

    ```console
    export SENZING_COMPONENTS=( \
      "g2config" \
      "g2configmgr" \
      "g2diagnostic" \
      "g2engine" \
      "g2hasher" \
      "g2product" \
      "g2ssadm" \
    )

    ```

1. Create files.
   Example:

    ```console
   for SENZING_COMPONENT in ${SENZING_COMPONENTS[@]}; \
   do \
     export SENZING_OUTPUT_DIR=${GIT_REPOSITORY_DIR}/protobuf/${SENZING_COMPONENT};
     mkdir -p ${SENZING_OUTPUT_DIR}
     protoc \
       --proto_path=${GIT_REPOSITORY_DIR}/proto/ \
       --go_out=${SENZING_OUTPUT_DIR} \
       --go_opt=paths=source_relative \
       --go-grpc_out=${SENZING_OUTPUT_DIR} \
       --go-grpc_opt=paths=source_relative \
       ${GIT_REPOSITORY_DIR}/proto/${SENZING_COMPONENT}.proto;
   done

    ```

### Test server

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

1. Start the test server.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make test-servegrpc

     ```

1. From the `senzing-62042001` message, copy the value of "SENZING_ENGINE_CONFIGURATION_JSON".
   It is in escaped JSON format.
   It will be used when working with the gRPC test client.

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

1. Use make target to run a docker images that builds RPM and DEB files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make package

    ```

   The results will be in the `${GIT_REPOSITORY_DIR}/target` directory.

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
