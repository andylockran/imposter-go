# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [0.7.0] - 2025-01-08
### Added
- feat(wsdl): prepopulate responses based on WSDL.
- feat: logs start up time.
- feat: supports more legacy format fields.

### Changed
- docs: splits legacy and current schema files.
- refactor: generalises capture signature to support interceptors.
- refactor: improves failed handler logging.
- refactor: lambda init should be conditional on runtime.
- refactor: moves configuration parsing earlier in lifecycle.

### Fixed
- fix(soap): check all attributes of root element to determine WSDL version.

## [0.6.1] - 2025-01-06
### Changed
- build: removes unsupported install flag.

### Fixed
- fix: request header property should be 'requestHeaders'.

## [0.6.0] - 2025-01-06
### Added
- feat: adds anyOf expression matcher.

### Changed
- docs: adds example for security config.
- refactor: improves expression matcher config naming.
- refactor: switches transformed security interceptors to anyOf expressions.

### Fixed
- fix: improves SOAP version namespace validation.

## [0.5.0] - 2025-01-05
### Added
- feat: transforms security blocks into interceptors.

### Changed
- build: adds since config.
- test: adds delay tolerance to more unit test conditions.

## [0.4.1] - 2025-01-05
### Changed
- test: adds delay tolerance to unit test.

## [0.4.0] - 2025-01-05
### Added
- feat: adds 'evals' matcher.

### Changed
- test: improves coverage for legacy config.
- test: splits matcher tests.

## [0.3.0] - 2025-01-04
### Added
- feat: adds config file schema and validator.
- feat: adds support for legacy config format.

## [0.2.0] - 2025-01-03
### Added
- feat: supports system XML namespaces.

### Changed
- docs: fixes JSONPath and XPath example configs.
- docs: improves release template.

## [0.1.2] - 2025-01-03
### Changed
- build: aligns binary name and build tags.
- docs: fixes example path.
- docs: improves installation instructions.

## [0.1.1] - 2025-01-03
### Changed
- ci: goreleaser should set internal version.

## [0.1.0] - 2025-01-03
### Added
- feat: first release.
