# changelog-go

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/l0nax/changelog-go)](https://goreportcard.com/report/gitlab.com/l0nax/changelog-go)

Keep a changelog without merge conflicts.

Every change gets its own small file under `.changelogs/unreleased/`, committed on
the branch that made the change. Two branches touching the changelog add two
different files, so they never conflict. On release, `changelog` collects them
into a `CHANGELOG.md` in the [Keep a Changelog](https://keepachangelog.com/)
format and files them under the version they shipped in.

The command is `changelog`. Linux, macOS and Windows.

## Install

### Release page

Download the archive for your platform from the
[releases page](https://gitlab.com/l0nax/changelog-go/-/releases), unpack it and
put `changelog` on your `$PATH`.

### Homebrew

```bash
brew tap l0nax/changelog-go https://gitlab.com/l0nax/changelog-go
brew install changelog
```

The formula is updated by hand after a release, so it may trail the latest
version by a little.

## Getting started

```bash
# once per project
changelog init

# once per change, on the branch that made it
changelog new -t bug_fix --title "Fix the sprocket"
git add .changelogs/unreleased && git commit -m "Add changelog entry"

# once per release
changelog release 1.4.0
git add .changelogs CHANGELOG.md && git commit -m "Release 1.4.0"
```

`init` writes `.changelog-go.toml` and creates the `.changelogs/` layout. Every
command finds that config by walking up from the working directory, so you can
run them from anywhere inside the project.

## Commands

| Command | What it does |
|---|---|
| `changelog init` | Create the config and directory layout |
| `changelog new` | Write a changelog entry |
| `changelog release <version>` | Release a version and regenerate the changelog file |
| `changelog latest` | Print the highest released version |
| `changelog next auto` | Print the version the pending entries imply |
| `changelog show <version>` | Print one release's notes, e.g. for a release page |
| `changelog check` | Verify the branch adds an entry — for CI |
| `changelog validate` | Parse everything and report what is wrong |
| `changelog migrate` | Migrate a v1 project — see [MIGRATION.md](MIGRATION.md) |
| `changelog version` | Print the version of `changelog` itself |
| `changelog completion <bash\|zsh>` | Print the shell completion script |

Run `changelog <command> --help` for the flags.

### Writing entries

With `--type` and `--title` the command writes the file and exits, so it works
in CI, a git hook or an editor plugin:

```bash
changelog new -t bug_fix --title "Fix the sprocket"
```

Leave either out and it prompts, unless stdin is not a terminal — then it fails
rather than hanging. `--interactive=false` forces that behaviour, `--dry-run`
prints the entry instead of writing it. An entry can carry a longer `--body`,
rendered underneath the title, and an `--author`.

### Releasing

`changelog release <version>` moves the pending entries into
`.changelogs/released/<version>/` and rewrites the changelog file. It refuses to
release nothing, so a wrong branch fails loudly; pass `--allow-empty` if a
version really has no user-facing changes. `--dry-run` prints the notes and
changes nothing.

`--auto` takes the version from the pending entries instead: a `rem_feat` entry
makes it a major bump, `new_feat` a minor one, anything else a patch. `next auto`
prints the same version without releasing it.

```bash
changelog release --auto
```

Pre-releases are detected from the version itself. A pre-release keeps its entry
files, so the final release of the same version carries the same entries; the
superseded pre-release stays in the changelog with its body folded into a
`<details>` block.

### Enforcing entries in CI

`changelog check` exits non-zero when a branch adds no entry, which is how you
keep the changelog from rotting:

```yaml
changelog:
  script:
    - changelog validate
    - changelog check --compare-with "origin/$CI_MERGE_REQUEST_TARGET_BRANCH_NAME"
```

It compares against the merge base, so a branch does not fail because the base
moved on. A branch that changes the changelog file itself is a release branch
rather than a change needing its own entry, so it is skipped. Use `--staged` in
a pre-commit hook.

For a release page, `show` prints one release's notes:

```bash
changelog show --no-heading "$(changelog latest)"
```

### Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Failure without a more specific code; `check` found no entry; `validate` found problems |
| 2 | Wrong invocation: a missing argument, or an unknown change type ID |
| 3 | The requested version does not exist, or nothing has been released |
| 4 | Nothing to do: releasing with no entries and without `--allow-empty` |

## Configuration

`.changelog-go.toml` lives at the project root. Paths in it are relative to that
file, not to the working directory.

```toml
changelog_dir = '.changelogs'
output_path = 'CHANGELOG.md'
version = '2'

# Prepended to every rendered version, whatever the directory on disk is
# called. Set it to '' for bare versions.
version_prefix = 'v'

[entry]
  [[entry.types]]
  id = 'bug_fix'
  group_title = 'Fixed'      # the heading in CHANGELOG.md
  title = 'Bug Fixed'        # the label in the picker
  effect = 'patch'           # major, minor or patch: what `next auto` does
  # hidden = true            # retire a type without breaking past releases

[pre_release]
detect = true                # derive the pre-release flag from the version
fold_pre_releases = true     # collapse a superseded pre-release into <details>
delete_pre_release = false   # or delete it outright; takes precedence

[check]
base_branch = ''             # what `check` compares against; empty = detect
```

Unknown keys are an error naming the offending line, so a typo never fails
silently. The same goes for entry files.

To stop using a change type, set `hidden = true` on it rather than deleting it:
that keeps it out of the picker while past releases still render under their
real title. Deleting it outright leaves old entries rendering under their raw
ID, with a warning.

## How it looks on disk

```
.changelog-go.toml
CHANGELOG.md
.changelogs/
├── released/
│   ├── 1.3.0/
│   │   ├── ReleaseInfo          # version, date, pre-release flag
│   │   └── develop-hIwUTCob     # one file per change
│   └── 1.4.0/
└── unreleased/
    └── 1733309257759            # waiting for the next release
```

By default nothing is thrown away: the entries that made up a release stay next
to it, which is what lets the changelog be regenerated from scratch at any time.
(`delete_pre_release` is the one exception, and it is off unless you ask for it.)

## Releasing

Maintainers: see [RELEASING.md](RELEASING.md).

## Upgrading from v1

See [MIGRATION.md](MIGRATION.md). Read it before running `changelog migrate` —
it rewrites your config.

## License

MIT. See [LICENSE](LICENSE).
