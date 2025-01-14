# How to build Docker image

## Senzing V4 Beta instructions

1. Email [sales@senzing.com] to request participation in Senzing V4 Beta.
1. Once you are signed up, you’ll get an email with [docker build information].
1. Use the information to set environment variables.

   Examples:

    ```console
    export SENZING_APT_REPOSITORY_URL="https://senzing-xxxxxxxx.com"
    ```

    ```console
    export SENZING_APT_REPOSITORY_NAME="senzingxxxx-xxxxxxxx.deb"
    ```

1. Build the first Docker image which contains the Senzing binaries.

    ```console
    docker build --build-arg SENZING_APT_REPOSITORY_NAME --build-arg SENZING_APT_REPOSITORY_URL --tag senzing/senzingsdk-runtime-beta:latest https://github.com/senzing/senzingsdk-runtime.git#main
    ```

1. Build the second Docker image which adds the gRPC server to the Senzing binaries.

    ```console
    docker build --tag senzing/serve-grpc:latest https://github.com/senzing-garage/serve-grpc.git#main
    ```

[sales@senzing.com]: mailto:sales@senzing.com
[docker build information]: ./response-email.png
