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
    export SENZING_OUTPUT_DIR=${GIT_REPOSITORY_DIR}/protobuf/${SENZING_COMPONENT}; \
    protoc \
        --proto_path=${GIT_REPOSITORY_DIR}/proto/ \
        --go_out=${SENZING_OUTPUT_DIR} \
        --go_opt=paths=source_relative \
        --go-grpc_out=${SENZING_OUTPUT_DIR} \
        --go-grpc_opt=paths=source_relative \
        ${GIT_REPOSITORY_DIR}/proto/${SENZING_COMPONENT}.proto; \
   done

    ```
