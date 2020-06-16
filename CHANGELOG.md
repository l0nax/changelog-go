# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).


## 1.1.4 (2020-06-16)

### Fixed (3 changes)
- Fix incorrect representation of the 'Number of Changes' in the CHANGELOG.md
- Fix bug in 'update avail.' checker
- Fix that no stacktrace will be send to Sentry. It also improves how the messages are stacked and send.


## 1.1.3 (2020-04-02)


## 1.1.2 (2020-04-02)

### Added (1 change)
- Add 'update available' reminder

### Fixed (1 change)
- Trim 'v' from version string when calling `release` command ([#4](https://gitlab.com/l0nax/changelog-go/-/issues/4))


## 1.1.1 (2020-03-30)

### Added (1 change)
- Add 'update' command

### Fixed (1 change)
- Fix 'error on number conversion' on Windows platform

### Changed (1 change)
- Disable automatic updates [#2](https://gitlab.com/l0nax/changelog-go/-/issues/2)


## 1.1.0 (2020-03-12)

### Other (1 change)
- Add sentry crash reporter

### Added (1 change)
- Add auto-updater


## 1.0.3 (2020-01-02)


## 1.0.2 (2020-01-02)

### Fixed (1 change)
- Replace '_' with '_' in changelog entry titles


## 1.0.1 (2020-01-02)

### Fixed (3 changes)
- Fix error if 'unreleased' folder does not exists
- Fix version data problems with GoReleaser
- Fix 'panic: file does not exist'

### Changed (2 changes)
- Remove 'version' flag
- Add 'version' sub-command


## 1.0.0 (2020-01-02)

### Added (1 change)
- Add basic functionality

