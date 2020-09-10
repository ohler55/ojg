# Changelog

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The structure and content of this file follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [1.2.1] - [2020-09-11]
### Fixed
- Order is preserved when using JSONPath to follow wildcards, unions, and slices.

## [1.2.0] - [2020-07-20]
### Added
- Parse Resuse option added to allow reusing maps on subsequent parses.
- In addition to callbacks, parsing multi-json documents can place elements on a `chan interface{}`.
### Changed
- A code refactoring resulting in a performance boost to Parsing and Validation.

## [1.1.4] - 2020-07-13
### Changed
- Validation speedup using a one switch statement and character maps.

## [1.1.3] - 2020-07-09
### Fixed
- Validator bug introduced in the speedup fixed.

## [1.1.2] - 2020-07-08
### Changed
- Performance improvement on validation and parsing.

## [1.1.1] - 2020-07-05
### Fixed
- Write bug that incorrectly wrote some UTF-8 sequences.

## [1.1.0] - 2020-07-04
### Added
- [Simple Encoding Notation](sen.md)
- Lazy input and out options to the `cmd/oj` command.

## [1.0.2] - 2020-07-01
### Added
- Filters will now iterate over Object members as well as Array members.

## [1.0.1] - 2020-06-23
### Added
- `cmd/oj` now correctly allows JSON as an argument in addition to reading from a file.

## [1.0.0] - 2020-06-22
### Added
- Initial release.
