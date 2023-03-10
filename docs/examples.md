# servegrpc examples

## Command line examples

### Command line example - Enable only g2engine gRPC service

For security reasons, it may be that only certain gRPC services are started.
In this example, only the G2Engine gRPC is started.

1. Using command line options.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools servegrpc \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --enable-g2engine
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
            "RESOURCEPATH": "/opt/senzing/g2/resources",
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
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools servegrpc
    ```

1. For more information, visit
   [SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)

## Docker examples
