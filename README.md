# `changelog-go`

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/changelog)
[![changelog](https://snapcraft.io//changelog/badge.svg)](https://snapcraft.io/changelog)

A Changelog Management Tool written in Go which is compatible with Linux,
Windows and Mac OS-X.

## Demo

[![asciicast](https://asciinema.org/a/291463.svg)](https://asciinema.org/a/291463)

## Motivation

Think about that: You have a Project where you work with other 4-5 Peoples. How do you track changes so that they can
easily added to a Changelog File?

This Tool helps you to keep track of you Changelog (and changes) and its fully
compatible with (eg.) the Git Flow. This Tool does **not** aim to replace Tools
like _Jira_ moreover, it extends it. Create Changelog Files easily with one
Command (`changelog release <version>`).

Feature branches can contain Changelog Entries and can be merged whenever you want!
It does not break the functionality.

**And** this Tool does _not_ disturb any other Tool (eg IDE, CI/CD, Artifact, ...)

## Features

- seamless integration into your Workflow and CI/CD environment
- works with Git, SVN, ... (every version control system)
- highly customizable
  * `CHANGELOG.md` output can be easily changed without writting any code
  * for internal and OpenSource Projects
- MIT license
- wide range of supported (and pre-compiled) Operating Systems
- keep track of changes made in your Project
- generates a changelog exactly the way you want it
- _free_ and _fast_ support


## Installation

There are many possibilities to install this Application.

### snap

```bash
~] snap install changelog
```

### script

You can also install `changelog-go` via a shell script:
```bash
~] curl https://gitlab.com/l0nax/changelog-go/raw/master/install.sh | bash
```

### RPM/ Deb

```bash
### (1) download the file from https://gitlab.com/l0nax/changelog-go/-/releases

### (2) install it via
~] rpm -iv PATH/TO/FILE
# or
~] dep -iv PATH/TO/FILE
```

### binary

1. Download the binary file from the [Release Page](https://gitlab.com/l0nax/changelog-go/-/releases)
2. Add the Path to the binary to your `$PATH` environment variable


## Usage

_NOTE: You don't have to be at the root path (eg. where your `.git` relies) to
generate/ create a Changelog(-Entry)_

### Initialize directory

1. Run `changelog init` and edit the config file (`.changelog-go.yaml`)

### Create a new change entry

1. Run `changelog new "<title>"` and then choose from the list which category
is the best fit.
2. Add and commit the new Changelog File

Example:
```bash   
## (1) create changelog entry
~] changelog new "Fix 'go pos'-Parser Bug"
Using config file: /root/.go/src/gitlab.com/l0nax/test/.changelog-go.yaml
[0] New Feature          (Added)
[1] Bug Fix              (Fixed)
[2] Feature change       (Changed)
[3] New deprecation      (Deprecated)
[4] Feature removal      (Removed)
[5] Security fix         (Security)
[6] Other                (Other)
>> 1

## (2) add and commit changelog Entry
~] git add .changelogs/unreleased/bugfix_P71128-114_bug-in-go-to-parser-gnAfCUyu
~] git commit -m "Add changelog entry"

## (3) push/ merge your changes or create (ie.) a Pull-Request
```

### Release a new version

Releasing a new Version is as simple as creating a new changelog entry:

1. Run `changelog release <version>`
2. Add and Commit the new `CHANGELOG.md` and changed files under `.changelogs`

Example:
```bash   
## (1) generate new changelog
~] changelog release 2.0.0

## (2) commit CHANGELOG.md and `.changelogs`
~] git add .changelogs CHANGELOG.md
~] git commit -m "Generating new CHANGELOG.md"

## (3) push/ merge your changes or create (ie.) a Pull-Request
```
