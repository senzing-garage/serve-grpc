# serve-grpc errors

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6011` - szconfig
1. `6012` - szconfigmgr
1. `6013` - szdiagnostic
1. `6014` - szengine
1. `6015` - szhasher
1. `6016` - szproduct
1. `6017` - szssadm

## Errors

## Common errors

### Postgresql

1. "Error: pq: SSL is not enabled on the server"
    1. The database URL needs the `sslmode` parameter.
       Example:

        ```console
        postgresql://username:password@postgres.example.com:5432/G2/?sslmode=disable
        ```

    1. [Connection String Parameters](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters)
