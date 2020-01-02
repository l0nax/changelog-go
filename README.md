# `changelog-go`

A Changelog Management Tool written in Go which is compatible with Linux,
Windows and Mac OS-X.

## Motivation

Think about that: You have a Project where you work with other 4-5 Peoples. How do you track changes so that they can
easily added to a Changelog File?

This Tool helps you to keep track of you Changelog (and changes) and its fully compatible with (eg) the Git Flow.
This Tool does **not** aim to replace Tools like _Jira_ moreover, it extends it.
Create Changelog Files easily with one Command (`changelog release <version>`).

Feature branches can contain Changelog Entries and can be merged whenever you want!
It does not break the functionality.

**And** this Tool does _not_ disturb any other Tool (eg IDE, CI/CD, Artifact, ...)


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

### Example workflow

```bash
## (1) Initialize git-flow
```
