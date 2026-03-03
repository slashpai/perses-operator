# Perses Operator Release

In all commands below, replace `<remote>` with the name of your git remote pointing to `github.com/perses/perses-operator` (commonly `origin` or `upstream`).

## 1. Prepare your release

- Fetch the latest changes and tags:

  ```bash
  git fetch <remote> --tags
  ```

- Create a release branch named `release/v<major>.<minor>` from `main` (this can be done from the GitHub UI via the branch dropdown on the repository page).

  > ⚠️ Release candidates and patch releases happen in the same `release/v<major>.<minor>` branch. Do not create `release/<version>` for patch or release candidate releases.

- Create a working branch: `<yourname>/release-v<major>.<minor>.<patch>`.

- Update `VERSION` with the new version.

- Regenerate manifests and installer:

  ```bash
  make manifests
  make build-installer
  ```

- Generate changelog:

  ```bash
  make generate-changelog
  ```

  > ⚠️ Tags must be fetched before this step — the generator uses `git describe --tags` to find the previous release.

- Review `CHANGELOG.md`:
  - Entries should be ordered: `[FEATURE]`, `[ENHANCEMENT]`, `[BUGFIX]`, `[BREAKINGCHANGE]`, `[DOC]`.
  - Each entry should include a pull request number.
  - Categorize or remove any `[UNKNOWN]` entries.
  - `[BREAKINGCHANGE]` entries should include a brief migration note.

- Push and open a PR against the release branch. Share the link in `#perses-dev` on CNCF Slack.
- Once approved, merge into the release branch.

## 2. Create release tag

```bash
git checkout release/v<major>.<minor>
make tag
git push <remote> v<major>.<minor>
```

Pushing the tag triggers a GitHub Actions workflow that builds and publishes Docker images to Docker Hub and creates a GitHub release with changelog notes and binary tarballs.

## 3. Merge the release into `main`

Keep the release branch around briefly in case a patch release is needed. When ready, open a PR to merge it into `main`.

> ⚠️ Use **"merge commit"**, not "squash and merge" — squashing deletes the tagged commit and can cause problems.
