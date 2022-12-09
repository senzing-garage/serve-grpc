# go-servegrpc

## Development

### Create *protobuf directories

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=go-servegrpc
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"
    ```

1. :pencil2: Identify Senzing subcomponent.
   Example:

    ```console
    export SENZING_COMPONENT=g2diagnostic
    ```

1. Make output directory.
   Example:

    ```console
    export SENZING_OUTPUT_DIR=${GIT_REPOSITORY_DIR}/${SENZING_COMPONENT}protobuf
    mkdir -p ${SENZING_OUTPUT_DIR}
    ```

1. Copy the *.proto file to `${SENZING_OUTPUT_DIR}/${SENZING_COMPONENT}.proto`
   Files are in [g2-sdk-proto](https://github.com/Senzing/g2-sdk-proto)

1. Create protobuf files.
   Example:

    ```console
    protoc \
        --proto_path=${SENZING_OUTPUT_DIR} \
        --go_out=${SENZING_OUTPUT_DIR} \
        --go_opt=paths=source_relative \
        --go-grpc_out=${SENZING_OUTPUT_DIR} \
        --go-grpc_opt=paths=source_relative \
        ${SENZING_OUTPUT_DIR}/${SENZING_COMPONENT}.proto
    ```
