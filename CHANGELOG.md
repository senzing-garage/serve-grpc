# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

-

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
