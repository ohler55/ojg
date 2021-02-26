# Changelog

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The structure and content of this file follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [1.7.1] - 2021-02-25
### Added
- Added HTMLUnsafe option to oj JSON writing to not encode &, <, and > to provide consistency
- Added HTMLSafe option to sen options to encode &, <, and > to provide consistency
### Fixed
- Fixed panic for `{"""":0}`. Now an error is returned.

## [1.7.0] - 2021-02-21
### Added
- Added support for a configuration file.
- Added ability to set colors when using the -c and -b option.
- Added ability to set HTML colors when using the -html option.

## [1.6.0] - 2021-02-19
### Added
- Added assembly plan package and cmd/oj option that allows assembling a new JSON from parsed data.
- Added sen.Parse() and sen.ParseReader() that use a new sen.DefaultParser
- Added the pretty package for prettier JSON layout.
- Added HTMLOptions for generating HTML color styled text.

## [1.5.0] - 2021-02-09
### Fixed
- Fixed reflection bug that occurred when a struct did not have the requested field.
### Added
- Added tab option for indentation.
### Changed
- Write operations now use panic and recovery internally for more
  robust error handling and for a very slight performance improvement.

## [1.4.1] - 2021-02-02
### Fixed
- The SEN parser and writer did not allow `\n` or `\t` in strings. It
  now does as would be expected from a friendly format.

## [1.4.0] - 2020-01-03
### Fixed
- JSONPath Slice end is now exclusive as called for in the Goessner description and the consensus.
- Nested array parsing bug fixed.

## [1.3.0] - 2020-10-28
### Added
- oj.Marshal added. The function fails if an un-encodeable value is encountered.
- UseTags option added for write and decompose Options.

## [1.2.1] - 2020-09-13
### Fixed
- Order is preserved when using JSONPath to follow wildcards, unions, and slices.

## [1.2.0] - 2020-07-20
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
