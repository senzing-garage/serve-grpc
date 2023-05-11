# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

-

## [0.4.4] - 2023-05-11

### Changed in 0.4.4

- Update dependencies

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
