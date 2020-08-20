# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).


## 1.2.0 (2020-08-20)

### Added (3 changes)
- Add custom CHANGELOG.md file template
- Sort change-entries per Change-Type in resulting CHANGELOG.md
- Sort entry types in CHANGELOG.md

### Other (1 change)
- Improve update routine


## 1.1.4 (2020-06-16)

### Fixed (5 changes)
- Fix bug in 'update avail.' checker
- Fix incorrect representation of the 'Number of Changes' in the CHANGELOG.md
- Fix runtime.boundError in 'release' command ([#8](https://gitlab.com/l0nax/changelog-go/-/issues/8))
- Fix runtime.errorString in CheckUpdate algo ([#9](https://gitlab.com/l0nax/changelog-go/-/issues/9))
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

### Changed (1 change)
- Disable automatic updates [#2](https://gitlab.com/l0nax/changelog-go/-/issues/2)

### Fixed (1 change)
- Fix 'error on number conversion' on Windows platform


## 1.1.0 (2020-03-12)

### Added (1 change)
- Add auto-updater

### Other (1 change)
- Add sentry crash reporter


## 1.0.3 (2020-01-02)


## 1.0.2 (2020-01-02)

### Fixed (1 change)
- Replace '_' with '_' in changelog entry titles


## 1.0.1 (2020-01-02)

### Changed (2 changes)
- Add 'version' sub-command
- Remove 'version' flag

### Fixed (3 changes)
- Fix 'panic: file does not exist'
- Fix error if 'unreleased' folder does not exists
- Fix version data problems with GoReleaser


## 1.0.0 (2020-01-02)

### Added (1 change)
- Add basic functionality

