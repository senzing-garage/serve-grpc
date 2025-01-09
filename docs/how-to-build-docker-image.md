# How to build Docker image

## Senzing V4 Beta instructions

1. Email [sales@senzing.com] to request participation in Senzing V4 Beta.
1. Once you are signed up, you’ll get an email with a [URL for a Debian-based installation].
   You’ll use values this URL to create environment variables.
1. The first part of the URL, up to the final forward-slash,
   is the SENZING_APT_REPOSITORY_URL value.
   It starts with `https://senzing-` and ends with `.com`.
   It needs to be the value of the `SENZING_APT_REPOSITORY_URL` environment variable.
   Example:

    ```console
    export SENZING_APT_REPOSITORY_URL="https://senzing-xxxxxxxx.com"
    ```

1. The second part of the URL, everything after the final forward-slash,
   is the SENZING_APT_REPOSITORY_NAME value.
   It starts with `senzing` and ends with `.deb`.
   It needs to be the value of the `SENZING_APT_REPOSITORY_NAME` environment variable.

   Example:

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
[URL for a Debian-based installation]: ./response-email.png
