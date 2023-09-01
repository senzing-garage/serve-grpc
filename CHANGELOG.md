# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

-

## [0.4.15] - 2023-09-01

### Changed in 0.4.15

- Last version before SenzingAPI 3.8.0

## [0.4.14] - 2023-08-17

### Changed in 0.4.14

- In `go.mod` update to `go 1.21`
- In `Dockerfile` update to `golang:1.21.0-bullseye`
- Update dependencies
  - github.com/senzing/g2-sdk-go-base v0.2.4
  - github.com/senzing/go-cmdhelping v0.1.7
  - github.com/senzing/go-common v0.2.13

## [0.4.13] - 2023-08-08

### Changed in 0.4.13

- Refactor to `template-go`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.8
  - github.com/senzing/g2-sdk-go-base v0.2.3
  - github.com/senzing/go-cmdhelping v0.1.5
  - github.com/senzing/go-common v0.2.11
  - github.com/senzing/go-logging v1.3.2
  - github.com/senzing/go-observing v0.2.7

## [0.4.12] - 2023-08-03

### Changed in 0.4.12

- Update dependencies
  - github.com/senzing/go-cmdhelping v0.1.4
  - github.com/senzing/go-common v0.2.8
- Refactored to template for multi-platform

## [0.4.11] - 2023-07-25

### Changed in 0.4.11

- In `Dockerfile`, added `HEALTHCHECK`
- Switch default port to 8261

## [0.4.10] - 2023-07-25

### Changed in 0.4.10

- Update dependencies
  - github.com/senzing/go-cmdhelping v0.1.1
  - github.com/senzing/go-common v0.2.5

## [0.4.9] - 2023-07-17

### Changed in 0.4.9

- In `Dockerfile`, updated to `senzing/senzingapi-runtime:3.6.0`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.7
  - github.com/senzing/g2-sdk-go-base v0.2.2
  - github.com/senzing/go-common v0.2.3
  - github.com/senzing/go-logging v1.3.1
  - github.com/senzing/senzing-tools v0.3.1
  - google.golang.org/grpc v1.56.2

## [0.4.8] - 2023-06-16

### Changed in 0.4.8

- Support for `--enable-all`
- In `Dockerfile`, updated to `senzing/senzingapi-runtime:3.5.3`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.5
  - github.com/senzing/g2-sdk-go-base v0.2.1
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230608182106-25c8cdc02e3c
  - github.com/senzing/go-common v0.1.4
  - github.com/senzing/go-logging v1.2.6
  - github.com/senzing/go-observing v0.2.6
  - github.com/senzing/senzing-tools v0.2.9-0.20230613173043-18f1bd4cafdb
  - github.com/spf13/viper v1.16.0
  - github.com/stretchr/testify v1.8.4
  - google.golang.org/grpc v1.56.0

## [0.4.7] - 2023-05-26

### Changed in 0.4.7

- Modified Load() to match `g2-sdk-proto/go`
- In Dockerfile, update Senzing binaries to 3.5.2
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.4
  - github.com/senzing/g2-sdk-go-base v0.2.0
  - github.com/senzing/g2-sdk-proto/go v0.0.0-20230526140633-b44eb0f20e1b

## [0.4.6] - 2023-05-19

### Changed in 0.4.6

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.3
  - github.com/senzing/g2-sdk-go-base v0.1.11

## [0.4.5] - 2023-05-17

### Added in 0.4.5

- Support for gRPC Observer aggregator

## [0.4.4] - 2023-05-11

### Changed in 0.4.4

- Update dependencies
  - github.com/senzing/g2-sdk-go-base v0.1.10
  - github.com/senzing/go-common v0.1.3
  - github.com/senzing/go-logging v1.2.3

## [0.4.3] - 2023-05-10

### Changed in 0.4.3

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.2
  - github.com/senzing/g2-sdk-go-base v0.1.9
  - github.com/senzing/go-observing v0.2.2

## [0.4.2] - 2023-04-21

### Changed in 0.4.2

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.1
  - github.com/senzing/g2-sdk-go-base v0.1.8

## [0.4.1] - 2023-04-18

### Changed in 0.4.1

- Updated dependencies
- Migrated from `github.com/senzing/go-logging/logger` to `github.com/senzing/go-logging/logging`

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
