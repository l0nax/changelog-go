# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
## 1.4.0 (2021-08-02)

### Added (1 change)
- Add `ci-get` command to render only latest release


## 1.3.2 (2021-07-20)

### Fixed (1 change)
- Fix bug where existing `ReleaseInfo` file of a release isn't processed by `changelog-go`


## 1.3.1 (2021-06-28)

### Fixed (1 change)
- Fix glibc problems with snap installation


## 1.3.0 (2021-06-26)

### Fixed (1 change)
- Enforce configuration presence only on `new` and `release` subcommands ([#11](https://gitlab.com/l0nax/changelog-go/-/issues/11))

### Deprecated (1 change)
- Deprecate 'gut' package

### Changed (3 changes)
- Compile application with go 1.16.x
- Search config file by going the file-system up stopping at '/' instead of relying on the git library
- Use go 1.16 "embed" package instead of 3rd-party tool


## 1.2.1 (2020-08-20)

### Fixed (1 change)
- Disable error message for update check


## 1.2.0 (2020-08-20)

### Changed (1 change)
- Improve update routine

### Added (3 changes)
- Add custom CHANGELOG.md file template
- Sort change-entries per Change-Type in resulting CHANGELOG.md
- Sort entry types in CHANGELOG.md


## 1.1.4 (2020-06-16)

### Fixed (5 changes)
- Fix bug in 'update avail.' checker
- Fix incorrect representation of the 'Number of Changes' in the CHANGELOG.md
- Fix runtime.boundError in 'release' command ([#8](https://gitlab.com/l0nax/changelog-go/-/issues/8))
- Fix runtime.errorString in CheckUpdate algo ([#9](https://gitlab.com/l0nax/changelog-go/-/issues/9))
- Fix that no stacktrace will be send to Sentry. It also improves how the messages


## 1.1.3 (2020-04-02)


## 1.1.2 (2020-04-02)

### Fixed (1 change)
- Trim 'v' from version string when calling `release` command ([#4](https://gitlab.com/l0nax/changelog-go/-/issues/4))

### Added (1 change)
- Add 'update available' reminder


## 1.1.1 (2020-03-30)

### Fixed (1 change)
- Fix 'error on number conversion' on Windows platform

### Changed (1 change)
- Disable automatic updates [#2](https://gitlab.com/l0nax/changelog-go/-/issues/2)

### Added (1 change)
- Add 'update' command


## 1.1.0 (2020-03-12)

### Changed (1 change)
- Add sentry crash reporter

### Added (1 change)
- Add auto-updater


## 1.0.3 (2020-01-02)


## 1.0.2 (2020-01-02)

### Fixed (1 change)
- Replace '_' with '_' in changelog entry titles


## 1.0.1 (2020-01-02)

### Fixed (3 changes)
- 'Fix ''panic: file does not exist'''
- Fix error if 'unreleased' folder does not exists
- Fix version data problems with GoReleaser

### Changed (2 changes)
- Add 'version' sub-command
- Remove 'version' flag


## 1.0.0 (2020-01-02)

### Added (1 change)
- Add basic functionality


## v2.0.0-rc.1 (2024-12-04)

### Changed (1 change)
- Switch to TOML as config language

### Added (3 changes)
- Add `latest` command
- Add `migrate` command to automatically migrate to the new v2 structure
- Add `next` command

