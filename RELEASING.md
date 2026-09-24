# Releasing

For maintainers. CI does the release; the Homebrew formula needs to be pushed manually.

## 1. Cut the release

On the default branch, with everything merged:

```bash
changelog validate
changelog next auto                       # what the pending entries imply
changelog release "$(changelog next --remove-prefix auto)"
```
Review the result, then commit:

```bash
git add .changelogs CHANGELOG.md
git commit -m "Release <version>"
git push
```

## 2. Tag

```bash
git tag -a v<version>
git push origin v<version>
```

The tag must come *after* the release commit: CI builds the notes with
`changelog show "$CI_COMMIT_TAG"`, which fails if the tag has no matching
directory under `.changelogs/released/`.

The `release` job then runs GoReleaser, which builds the archives and creates
the GitLab release with those notes.

## 3. Update the Homebrew formula

This is not automated: `CI_JOB_TOKEN` cannot write to the repository.

Download the `dist/` artifact from the `release` job, then:

```bash
cp "$(find dist -name '*.rb')" changelog.rb
git add changelog.rb
git commit -m "Update Homebrew formula for v<version>"
git push
```

The formula pins the version and the SHA256 of every archive, so it is wrong
until this is done. Skipping it is how the v1 formula ended up stuck at 1.4.1
for three years.

## Pre-releases

Same steps with a pre-release version, e.g. `changelog release 2.1.0-rc.1`.
Two differences, both automatic:

- The entry files stay in `unreleased/`, so the final release of that version
  carries the same entries.
- GoReleaser marks it as a pre-release, and it is not published to Homebrew.
