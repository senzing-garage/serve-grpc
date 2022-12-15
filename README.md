# go-servegrpc

## Development

### Create protobuf directories

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=go-servegrpc
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

1. Start the test server.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make test-servegrpc
     ```

1. From the `senzing-99992001` message, copy the value of "SENZING_ENGINE_CONFIGURATION_JSON".
   It is in escaped JSON format.
   It will be used when working with the gRPC test client.

1. Start the gRPC test client.
   Example:

    ```console
    bloomrpc
    ```

1. In `bloomrpc`:
    1. In upper-left, click on plus sign ("+").
        1. Navigate to the ${GIT_REPOSITORY_DIR}/proto directory
        1. Choose a `.proto` file
    1. Near top-center, use the address of `0.0.0.0:50051`
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
