# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], [markdownlint],
and this project adheres to [Semantic Versioning].

## [Unreleased]

-

## [0.9.17] - 2025-07-09

### Added in 0.9.17

- Update dependencies

## [0.9.16] - 2025-07-02

### Added in 0.9.16

- Update dependencies

## [0.9.15] - 2025-06-18

### Added in 0.9.15

- Rename `SzConfigManager.GetConfigs` to `SzConfigManger.GetConfigRegistry`

## [0.9.14] - 2025-06-16

### Added in 0.9.14

- Added `SENZING_TOOLS_LICENSE_STRING_BASE64`

## [0.9.13] - 2025-06-13

### Changed in 0.9.13

- Update artifacts

## [0.9.12] - 2025-06-10

### Added in 0.9.12

- Command options for gRPC
  - `SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_MIN_TIME_IN_SECONDS`
  - `SENZING_TOOLS_SERVER_KEEPALIVE_ENFORCEMENT_POLICY_PERMIT_WITHOUT_STREAM`
  - `SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_GRACE_IN_SECONDS`
  - `SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_AGE_IN_SECONDS`
  - `SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_MAX_CONNECTION_IDLE_IN_SECONDS`
  - `SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIME_IN_SECONDS`
  - `SENZING_TOOLS_SERVER_KEEPALIVE_SERVER_PARAMETER_TIMEOUT_IN_SECONDS`
  - `SENZING_TOOLS_SERVER_MAX_CONCURRENT_STREAMS`
  - `SENZING_TOOLS_SERVER_MAX_HEADER_LIST_SIZE_IN_BYTES`
  - `SENZING_TOOLS_SERVER_MAX_RECEIVE_MESSAGE_SIZE_IN_BYTES`
  - `SENZING_TOOLS_SERVER_MAX_SEND_MESSAGE_SIZE_IN_BYTES`
  - `SENZING_TOOLS_SERVER_READ_BUFFER_SIZE_IN_BYTES`
  - `SENZING_TOOLS_SERVER_WRITE_BUFFER_SIZE_IN_BYTES`

## [0.9.11] - 2025-06-02

### Added in 0.9.11

- Improved error messages

## [0.9.10] - 2025-05-08

### Added in 0.9.10

- CORS

## [0.9.9] - 2025-05-07

### Added in 0.9.9

- gRPC over HTTP

## [0.9.8] - 2025-05-02

### Changed in 0.9.8

- Update dependencies

## [0.9.7] - 2025-04-18

### Changed in 0.9.7

- Migrate from `_Ctype_longlong` to `_Ctype_int64_t`

## [0.9.6] - 2025-04-15

### Added in 0.9.6

- `SzConfig.VerifyConfig`
- `SzEngine.WhySearch`

## [0.9.5] - 2025-04-09

### Changed in 0.9.5

- Update dependencies

## [0.9.4] - 2025-03-28

### Added in 0.9.4

- Update dependencies

## [0.9.3] - 2025-03-20

### Added in 0.9.3

- TLS pass phrase support

## [0.9.2] - 2025-03-14

### Added in 0.9.2

- Mutual TLS support

## [0.9.1] - 2025-02-27

### Added in 0.9.1

- Support `SENZING_PATH`

## [0.9.0] - 2025-02-23

### Added in 0.9.0

- Server-side TLS support

## [0.8.12] - 2025-02-12

### Changed in 0.8.12

- Update dependencies
- `.proto` field names

## [0.8.11] - 2025-02-10

### Changed in 0.8.11

- Update dependencies

## [0.8.10] - 2025-01-31

### Changed in 0.8.10

- Update dependencies

## [0.8.9] - 2024-12-17

### Changed in 0.8.9

- Update dependencies

## [0.8.8] - 2024-12-11

### Changed in 0.8.8

- Update dependencies

## [0.8.7] - 2024-11-14

### Changed in 0.8.7

- Support SQLite in-memory database

## [0.8.6] - 2024-10-30

### Changed in 0.8.6

- Update dependencies

## [0.8.5] - 2024-10-21

### Changed in 0.8.5

- Add ephemeral database to Dockerfile

## [0.8.4] - 2024-10-09

### Changed in 0.8.4

- Update dependencies

## [0.8.3] - 2024-10-01

### Added in 0.8.3

- Added `PreprocessRecord()`

## [0.8.2] - 2024-09-11

### Changed in 0.8.2

- Update dependencies

## [0.8.1] - 2024-08-27

### Changed in 0.8.1

- Modify method calls to match Senzing API 4.0.0-24237

## [0.8.0] - 2024-08-23

### Changed in 0.8.0

- Change from `g2` to `sz`/`er`

## [0.7.7] - 2024-08-12

### Changed in 0.7.7

- Updated `senzing/senzingapi-runtime-staging` to 4.0.0.24211

## [0.7.6] - 2024-08-05

### Changed in 0.7.6

- Fix permissions on Sqlite database file

## [0.7.5] - 2024-08-05

### Changed in 0.7.5

- Improve non-root container

## [0.7.4] - 2024-06-26

### Changed in 0.7.4

- Updated dependencies

## [0.7.3] - 2024-06-17

### Changed in 0.7.3

- Update methods to Senzing 4.0.0-24162
- From `GrpcServerImpl` to `BasicGrpcServer`

## [0.7.2] - 2024-05-08

### Added in 0.7.2

- `SzDiagnostic.GetFeature`
- `SzEngine.FindInterestingEntitiesByEntityId`
- `SzEngine.FindInterestingEntitiesByRecordId`
- `SzEngine.ProcessRedoRecord`

### Changed in 0.7.2

- Update `Dockerfile`

### Deleted in 0.7.2

- `SzEngine.GetRepositoryLastModifiedTime`

## [0.7.1] - 2024-04-26

### Changed in 0.7.1

- Moved from "g2" to "sz"

## [0.7.0] - 2024-03-01

### Changed in 0.7.0

- Updated dependencies
- Deleted methods not used in V4

## [0.6.1] - 2024-01-29

### Changed in 0.6.1

- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.9.0
  - github.com/senzing-garage/g2-sdk-go-base v0.5.0
  - github.com/senzing-garage/g2-sdk-proto/go v0.0.0-20240126210601-d02d3beb81d4
  - google.golang.org/grpc v1.61.0

## [0.6.0] - 2024-01-03

### Changed in 0.6.0

- Renamed module to `github.com/senzing-garage/serve-grpc`
- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/senzing-garage/go-cmdhelping v0.2.0
  - github.com/senzing-garage/go-common v0.4.0
  - github.com/senzing-garage/go-logging v1.4.0
  - github.com/senzing-garage/go-observing v0.3.0
  - github.com/senzing/g2-sdk-go v0.8.0
  - github.com/senzing/g2-sdk-go-base v0.4.0
  - github.com/spf13/viper v1.18.2
  - google.golang.org/grpc v1.60.1

## [0.5.5] - 2023-12-08

### Changed in 0.5.5

- In `Dockerfile` update to:
  - `golang:1.21.40-bullseye`
  - `senzing/senzingapi-runtime:3.8.0`
- Updated `testdata/senzing-license/g2.lic`
- Update dependencies
  - github.com/spf13/cobra v1.8.0
  - github.com/spf13/viper v1.18.1

## [0.5.4] - 2023-11-01

### Changed in 0.5.4

- Update dependencies
  - github.com/senzing-garage/g2-sdk-go-base v0.3.3
  - github.com/senzing-garage/go-common v0.3.2-0.20231018174900-c1895fb44c30

## [0.5.3] - 2023-10-23

### Changed in 0.5.3

- Update to [template.go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.7.4
  - github.com/senzing-garage/g2-sdk-go-base v0.3.2
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20231016131354-0d0fba649357
  - github.com/senzing-garage/go-cmdhelping v0.1.9
  - github.com/senzing-garage/go-common v0.3.1
  - github.com/senzing-garage/go-logging v1.3.3
  - github.com/senzing-garage/go-observing v0.2.8
  - google.golang.org/grpc v1.59.0

## [0.5.2] - 2023-10-13

### Changed in 0.5.2

- Changed from `int` to `int64` where required by the SenzingAPI
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.7.3
  - github.com/senzing-garage/g2-sdk-go-base v0.3.1
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20231013142630-30a869751ff0
  - google.golang.org/grpc v1.58.3

### Deleted in 0.5.2

- `g2product.ValidateLicenseFile`
- `g2product.ValidateLicenseStringBase64`

## [0.5.1] - 2023-10-02

### Changed in 0.5.1

- In `Dockerfile` update to:
  - `golang:1.21.0-bullseye@sha256:02f350d8452d3f9693a450586659ecdc6e40e9be8f8dfc6d402300d87223fdfa`
  - `senzing/senzingapi-runtime:staging` - until a release of Senzing with Go support
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go-base v0.3.0

## [0.5.0] - 2023-09-25

### Changed in 0.5.0

- Supports SenzingAPI 3.8.0
- Deprecated functions have been removed

## [0.4.15] - 2023-09-01

### Changed in 0.4.15

- Last version before SenzingAPI 3.8.0

## [0.4.14] - 2023-08-17

### Changed in 0.4.14

- In `go.mod` update to `go 1.21`
- In `Dockerfile` update to `golang:1.21.0-bullseye`
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go-base v0.2.4
  - github.com/senzing-garage/go-cmdhelping v0.1.7
  - github.com/senzing-garage/go-common v0.2.13

## [0.4.13] - 2023-08-08

### Changed in 0.4.13

- Refactor to `template-go`
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.8
  - github.com/senzing-garage/g2-sdk-go-base v0.2.3
  - github.com/senzing-garage/go-cmdhelping v0.1.5
  - github.com/senzing-garage/go-common v0.2.11
  - github.com/senzing-garage/go-logging v1.3.2
  - github.com/senzing-garage/go-observing v0.2.7

## [0.4.12] - 2023-08-03

### Changed in 0.4.12

- Update dependencies
  - github.com/senzing-garage/go-cmdhelping v0.1.4
  - github.com/senzing-garage/go-common v0.2.8
- Refactored to template for multi-platform

## [0.4.11] - 2023-07-25

### Changed in 0.4.11

- In `Dockerfile`, added `HEALTHCHECK`
- Switch default port to 8261

## [0.4.10] - 2023-07-25

### Changed in 0.4.10

- Update dependencies
  - github.com/senzing-garage/go-cmdhelping v0.1.1
  - github.com/senzing-garage/go-common v0.2.5

## [0.4.9] - 2023-07-17

### Changed in 0.4.9

- In `Dockerfile`, updated to `senzing/senzingapi-runtime:3.6.0`
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.7
  - github.com/senzing-garage/g2-sdk-go-base v0.2.2
  - github.com/senzing-garage/go-common v0.2.3
  - github.com/senzing-garage/go-logging v1.3.1
  - github.com/senzing-garage/senzing-tools v0.3.1
  - google.golang.org/grpc v1.56.2

## [0.4.8] - 2023-06-16

### Changed in 0.4.8

- Support for `--enable-all`
- In `Dockerfile`, updated to `senzing/senzingapi-runtime:3.5.3`
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.5
  - github.com/senzing-garage/g2-sdk-go-base v0.2.1
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230608182106-25c8cdc02e3c
  - github.com/senzing-garage/go-common v0.1.4
  - github.com/senzing-garage/go-logging v1.2.6
  - github.com/senzing-garage/go-observing v0.2.6
  - github.com/senzing-garage/senzing-tools v0.2.9-0.20230613173043-18f1bd4cafdb
  - github.com/spf13/viper v1.16.0
  - github.com/stretchr/testify v1.8.4
  - google.golang.org/grpc v1.56.0

## [0.4.7] - 2023-05-26

### Changed in 0.4.7

- Modified Load() to match `g2-sdk-proto/go`
- In Dockerfile, update Senzing binaries to 3.5.2
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.4
  - github.com/senzing-garage/g2-sdk-go-base v0.2.0
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230526140633-b44eb0f20e1b

## [0.4.6] - 2023-05-19

### Changed in 0.4.6

- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.3
  - github.com/senzing-garage/g2-sdk-go-base v0.1.11

## [0.4.5] - 2023-05-17

### Added in 0.4.5

- Support for gRPC Observer aggregator

## [0.4.4] - 2023-05-11

### Changed in 0.4.4

- Update dependencies
  - github.com/senzing-garage/g2-sdk-go-base v0.1.10
  - github.com/senzing-garage/go-common v0.1.3
  - github.com/senzing-garage/go-logging v1.2.3

## [0.4.3] - 2023-05-10

### Changed in 0.4.3

- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.2
  - github.com/senzing-garage/g2-sdk-go-base v0.1.9
  - github.com/senzing-garage/go-observing v0.2.2

## [0.4.2] - 2023-04-21

### Changed in 0.4.2

- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.6.1
  - github.com/senzing-garage/g2-sdk-go-base v0.1.8

## [0.4.1] - 2023-04-18

### Changed in 0.4.1

- Updated dependencies
- Migrated from `github.com/senzing-garage/go-logging/logger` to `github.com/senzing-garage/go-logging/logging`

## [0.4.0] - 2023-03-27

### Added in 0.4.0

- Repository name change from `servegrpc` to `serve-grpc`

## [0.3.9] - 2023-03-27

### Added in 0.3.9

- Added Stream methods
  - g2diagnosticserver.StreamEntityListBySize()
  - g2engineserver.StreamExportCSVEntityReport()
  - g2engineserver.StreamExportJSONEntityReport()
- Update dependencies
- Last versioned release before name change to serve-grpc

## [0.3.8] - 2023-03-14

### Changed in 0.3.8

- Update dependencies
- Standardize use of Viper/Cobra

## [0.3.7] - 2023-03-13

### Changed in 0.3.7

- Improved documentation

### Fixed in 0.3.7

- Fixed issue silent error when connecting to database.

## [0.3.6] - 2023-03-08

### Changed in 0.3.6

- Fixed issue with Cobra, Viper, and command parameters, again.

## [0.3.5] - 2023-03-07

### Changed in 0.3.5

- Fixed issue with Cobra, Viper, and command parameters

## [0.3.4] - 2023-03-03

### Added in 0.3.4

- Normalized input parameters

## [0.3.3] - 2023-02-16

### Added in 0.3.3

- Add a default Senzing configuration to the SQLite database, `/tmp/sqlite/G2C.db`

## [0.3.2] - 2023-02-16

### Added in 0.3.2

- A test SQLite database to the Docker image, `/tmp/sqlite/G2C.db`

## [0.3.1] - 2023-02-16

### Changed in 0.3.1

- Refactored to reduce cyclomatic complexities

## [0.3.0] - 2023-02-15

### Added to 0.3.0

- Using refactored g2-sdk-go
- Removed need for rootfs

## [0.2.0] - 2023-02-07

### Added to 0.2.0

- Added Observer support

## [0.1.1] - 2023-01-30

### Added to 0.1.1

- Initialize Senzing objects upon startup
- Disable Init() and Destroy()

## [0.1.0] - 2023-01-13

### Added to 0.1.0

- Initial functionality

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[markdownlint]: https://dlaa.me/markdownlint/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
