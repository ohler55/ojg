# Changelog

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The structure and content of this file follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [1.13.0] - 2022-03-05
### Added
- Added jp.Expr.Has() function.
- Added jp.Walk to walk data and provide a the path and value for each
  element.

## [1.12.14] - 2022-02-28
### Fixed
- `[]byte` are encoded according to the ojg.Options.

## [1.12.13] - 2022-02-23
### Fixed
- For JSONPath (jp) reflection Get returns `has` value correctly for zero field values.

## [1.12.12] - 2021-12-27
### Fixed
- JSONPath scripts (jp.Script or [?(@.foo == 123)]) is now thread safe.

## [1.12.11] - 2021-12-10
### Fixed
- Parser reuse was no resetting callback and channels. It does now.

## [1.12.10] - 2021-12-07
### Added
- Added a delete option to the oj application.

## [1.12.9] - 2021-10-31
### Fixed
- Stuttering extracted elements when using the `-x` options has been fixed.

## [1.12.8] - 2021-09-21
### Fixed
- Correct unicode character is now included in error messages.

## [1.12.7] - 2021-09-14
### Fixed
- Typo in maxEnd for 32 bit architecture fixed.
- json.Unmarshaler fields in a struct correctly unmarshal.

## [1.12.6] - 2021-09-12
### Fixed
- Due to limitation (a bug most likely) in the stardard math package
  math.MaxInt64 can not be used on 32 bit architectures. Changes were
  made to work around this limitation.

- Embedded (Anonymous) pointers types now encode correctly.

### Added
- Support for json.Unmarshaler interface added.

## [1.12.5] - 2021-08-17
# Changed
- Updated to use go 1.17.

## [1.12.4] - 2021-08-06
### Fixed
- Setting an element in an array that does not exist now creates the array is the Nth value is not negative.

## [1.12.3] - 2021-08-01
### Fixed
- Error message on failed recompose was fixed to display the correct error message.
- Marshal of a non-pointer that contains a json.Marshaller that is not a pointer no longer fails.

## [1.12.2] - 2021-07-28
### Fixed
- Structs with recursive lists no longer fail.

## [1.12.1] - 2021-07-23
### Fixed
- Applying filters to a non-simple list such as `[]*Sample` now supported as expected.

## [1.12.0] - 2021-07-03
### Added
- SEN format parsing now allows string to be delimited with the single quote character.
- SEN format parsing now allows strings to be concatenated with syntax like `["abc" + "def"]`.
- SEN format parsing now allows functions such as `ISODate("2021-06-28T10:11:12Z")` in SEN data.
### Changed
- When Pretty Align is true map members are now aligned.

## [1.11.1] - 2021-05-29
### Fixed
- Missing support for json.Marshaler and encoding.TextMarshaler added.

## [1.11.0] - 2021-05-23
### Fixed
- Struct with pointers to base types such as *float64 are fixed.
- Stack overflow when converting values to JSON which are a type alias
  of a builtin.
### Added
- Added `[]byte` converation option for decompose.
- Added MustXxx versions of multiple functions to allow a panic and recover code pattern.
### Changed
- oj.Unmarshal now emits float64 for all numbers instead of int64 for
  integers. The parse functions remain unchanged.

## [1.10.0] - 2021-04-22
### Fixed
- Multiple part json tags are now parsed correctly and the string
  options is supported in both decompose and compose.
### Added
- Tokenize callback parser added.

## [1.9.5] - 2021-04-04
### Fixed
- OmitNil now catches nil maps and slices more consistently.

## [1.9.4] - 2021-04-04
### Fixed
- Number parsing in the form of 2e-7 has been fixed.

## [1.9.3] - 2021-03-30
### Fixed
- Writer functions now decompose structs if possible instead of resorting to %v too quickly.

## [1.9.2] - 2021-03-24
### Fixed
- When parsing SEN format `\r` is now allowed in strings to support
  Windows line termination as it works in Linux and macOS.

## [1.9.1] - 2021-03-21
### Fixed
- oj.Unmarshal now supports the optional alt.Recomposer as documented.
- Recomposer handles time.Time recomposing like any other struct.
- Write writes time.Time to conform to other struct encoding.

## [1.9.0] - 2021-03-13
### Added
- The Recomposer is now more flexibly in regard to input types. It now
  allows json.Unmarshal() targets as well as the type create key
  approach.
- Added flag to alt.Options to determine whether embedded anonymous
  types whould be output as nested elements or flattened.
- Added oj.Unmarshal and sen.Unmarshal.

## [1.8.0] - 2021-03-05
### Added
- Added alignment option for pretty printing.
- Added alt.Diff() and alt.Compare().
- Added color option for encoded time.
- Add alt.Converter along with some built in converter for time and mongodb export maps.

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
