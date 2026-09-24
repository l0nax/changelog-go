# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v2.0.0 (2026-09-24)

### Fixed (13 changes)
- Fix `delete_pre_release` never being read
- Fix `fold_pre_releases` never taking effect
- Fix `init` not creating the changelog directories, which broke `release` on a fresh project
- Fix `latest` never loading the config, so the command always failed
- Fix `next auto` deriving the version from the previous release instead of the pending entries
- Fix `v`-prefixed versions sorting as 0.0.0, which misplaced them in the changelog
- Fix a leaked file handle when copying entry files
- Fix a missing blank line before the first version heading
- Fix migration overwriting the config with stock defaults, losing the `other` change type
- Fix panic on a `v`-prefixed latest version
- Fix paths being resolved against the working directory, so commands now work from a subdirectory
- Fix release timestamps being written with nanosecond precision
- Fix removing a change type from the config breaking every past release

### Changed (8 changes)
- Refuse to release with no pending entries unless `--allow-empty` is given
- Reject unknown keys in the config and in entry files, naming the offending line
- Rename the pre-release config keys to snake_case; the old camelCase spellings are now rejected with the replacement named
- Render a superseded pre-release collapsed instead of dropping its entries
- Replace the deprecated TOML library with pelletier/go-toml/v2
- Report the build version from the Go build info instead of linker stamps
- Send logs to stderr so that stdout carries only the requested value
- Switch to TOML as config language

### Added (15 changes)
- Add `--dry-run` and `--allow-empty` to `release`
- Add `--remove-prefix` to `latest` and `next`, and `--skip-prereleases` to `latest`
- Add `--type`, `--title`, `--body`, `--author` and `--dry-run` to `new`, so it runs without a terminal
- Add `check` command to enforce a changelog entry per merge request
- Add `latest` command
- Add `migrate` command to automatically migrate to the new v2 structure
- Add `next` command
- Add `release --auto`, deriving the version from the pending entries
- Add `show` command rendering a single release's notes
- Add `validate` command reporting every problem in one pass
- Add `version_prefix` config key, applied when rendering versions
- Add `version` command
- Add global `--config` and `--changelog-dir` overrides
- Add optional entry bodies, rendered underneath the title
- Add shell completions and a man page


## v2.0.0-rc.1 (2024-12-04)

<details>
<summary>This is a Pre-Release, Click to see details.</summary>

### Changed (1 change)
- Switch to TOML as config language

### Added (3 changes)
- Add `latest` command
- Add `migrate` command to automatically migrate to the new v2 structure
- Add `next` command

</details>


## v1.4.0 (2021-08-02)

### Added (1 change)
- Add `ci-get` command to render only latest release


## v1.3.2 (2021-07-20)

### Fixed (1 change)
- Fix bug where existing `ReleaseInfo` file of a release isn't processed by `changelog-go`


## v1.3.1 (2021-06-28)

### Fixed (1 change)
- Fix glibc problems with snap installation


## v1.3.0 (2021-06-26)

### Fixed (1 change)
- Enforce configuration presence only on `new` and `release` subcommands ([#11](https://gitlab.com/l0nax/changelog-go/-/issues/11))

### Deprecated (1 change)
- Deprecate 'gut' package

### Changed (3 changes)
- Compile application with go 1.16.x
- Search config file by going the file-system up stopping at '/' instead of relying on the git library
- Use go 1.16 "embed" package instead of 3rd-party tool


## v1.2.1 (2020-08-20)

### Fixed (1 change)
- Disable error message for update check


## v1.2.0 (2020-08-20)

### Changed (1 change)
- Improve update routine

### Added (3 changes)
- Add custom CHANGELOG.md file template
- Sort change-entries per Change-Type in resulting CHANGELOG.md
- Sort entry types in CHANGELOG.md


## v1.1.4 (2020-06-16)

### Fixed (5 changes)
- Fix bug in 'update avail.' checker
- Fix incorrect representation of the 'Number of Changes' in the CHANGELOG.md
- Fix runtime.boundError in 'release' command ([#8](https://gitlab.com/l0nax/changelog-go/-/issues/8))
- Fix runtime.errorString in CheckUpdate algo ([#9](https://gitlab.com/l0nax/changelog-go/-/issues/9))
- Fix that no stacktrace will be send to Sentry. It also improves how the messages


## v1.1.3 (2020-04-02)


## v1.1.2 (2020-04-02)

### Fixed (1 change)
- Trim 'v' from version string when calling `release` command ([#4](https://gitlab.com/l0nax/changelog-go/-/issues/4))

### Added (1 change)
- Add 'update available' reminder


## v1.1.1 (2020-03-30)

### Fixed (1 change)
- Fix 'error on number conversion' on Windows platform

### Changed (1 change)
- Disable automatic updates [#2](https://gitlab.com/l0nax/changelog-go/-/issues/2)

### Added (1 change)
- Add 'update' command


## v1.1.0 (2020-03-12)

### Changed (1 change)
- Add sentry crash reporter

### Added (1 change)
- Add auto-updater


## v1.0.3 (2020-01-02)


## v1.0.2 (2020-01-02)

### Fixed (1 change)
- Replace '_' with '_' in changelog entry titles


## v1.0.1 (2020-01-02)

### Fixed (3 changes)
- 'Fix ''panic: file does not exist'''
- Fix error if 'unreleased' folder does not exists
- Fix version data problems with GoReleaser

### Changed (2 changes)
- Add 'version' sub-command
- Remove 'version' flag


## v1.0.0 (2020-01-02)

### Added (1 change)
- Add basic functionality

