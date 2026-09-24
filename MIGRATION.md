# Migrating from changelog-go v1 to v2

`changelog migrate` does the work. Read this first — it rewrites your config
file in place, and there are three changes it cannot make for you.

```bash
git checkout -b migrate-changelog-go
changelog migrate
changelog validate
git diff
```

`validate` parses everything the migration touched and reports what is wrong, so
run it before you commit.

## What migrate does

It converts `.changelog-go.yaml` into `.changelog-go.toml` and deletes the old
file, then rewrites every entry file and every `ReleaseInfo` from the v1 YAML
format into TOML.

These settings are carried across:

| v1 (YAML) | v2 (TOML) |
|---|---|
| `preRelease.detect` | `pre_release.detect` |
| `preRelease.deletePreRelease` | `pre_release.delete_pre_release` |
| `preRelease.foldPreReleases` | `pre_release.fold_pre_releases` |
| `changelog.entryPath` | `changelog_dir` |
| `changelog.changelog` | `output_path` |

These are dropped, with a warning:

- `entry.author` — v2 records an author when one is given, with no setting to
  turn it on or off.
- `changelog.customScheme` — v2 has no template customization yet. If you relied
  on a custom scheme, do not upgrade yet.

## What you have to do yourself

### 1. The command is still `changelog`, but the config file moved

`.changelog-go.yaml` is replaced by `.changelog-go.toml`. If anything in your CI
or tooling references the YAML file by name, update it.

### 2. Change types are now yours to define

v1 had six built-in types, numbered 0 to 6 in the entry files. v2 has no built-in
types at all: they live in your config under `[[entry.types]]`, and an entry
names one by ID.

The migration maps the old numbers onto the IDs it writes into your new config:

| v1 | v2 ID | Heading |
|---|---|---|
| 0 | `new_feat` | Added |
| 1 | `bug_fix` | Fixed |
| 2 | `feat_change` | Changed |
| 3 | `deprecate` | Deprecated |
| 4 | `rem_feat` | Removed |
| 5 | `security` | Security |
| 6 | `other` | Other |

Note the last one. v1's type 6 was "Other" and had no v2 equivalent, so the
migration adds an `other` type to your config with `effect = 'patch'`. If you
would rather those entries counted as a minor bump, change its `effect` after
migrating.

### 3. Every heading in CHANGELOG.md gets restyled

v2 renders versions through `version_prefix`, which defaults to `'v'`. On the
first regeneration every heading picks it up, whether or not the directory on
disk has the prefix:

```diff
-## 1.4.0 (2021-08-02)
+## v1.4.0 (2021-08-02)
```

Set `version_prefix = ''` if you want the old bare headings. Either way the
directories under `.changelogs/released/` keep the names they have — renaming
them would break any link pointing at them.

The regeneration happens on your next `changelog release`. Expect a large but
mechanical diff once.

## Things that changed in v2 regardless of migration

- **`new` no longer needs a terminal.** `changelog new -t bug_fix --title "..."`
  writes the entry and exits, which is what makes it usable from CI or a git
  hook. Without a TTY and without those flags it now fails instead of hanging.
- **Logs go to stderr**, so `$(changelog latest)` yields just the version.
- **Releasing nothing is an error**, exit 4. Pass `--allow-empty` for a version
  with no user-facing changes.
- **Unknown config keys are an error** naming the line, instead of being ignored.
  The same goes for entry files.
- **New commands**: `check`, `validate`, `show`, `next`, `latest`, `version`,
  `completion`.

## If migrate goes wrong

Everything it does is a file rewrite inside your repository, so
`git checkout -- .` puts it all back. That is the reason for the branch in the
first command above.

`changelog migrate --force` re-runs the migration on a project it thinks is
already migrated. It overwrites `.changelog-go.toml` with a freshly generated
one, so any hand-editing you did to it is lost — take a copy first.
