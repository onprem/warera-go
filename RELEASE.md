# Release process

`warera-go` is a Go library. A version is published to Go's module infrastructure when a semantic version tag is pushed; the repository does not need compiled release binaries. The project is currently pre-stable, so use `v0.x.y` versions and tell users that the API can change between releases. Go's module guidance describes `v0` as unstable and without compatibility guarantees.

## Versioning

- Use `v0.x.0` for a feature release or a release that includes public API changes.
- Use `v0.x.y` for a bug-fix-only release that does not change the public API.
- Start a new minor series when the current series has accumulated enough changes; do not imply the compatibility guarantees associated with `v1.x`.
- Choose `v1.0.0` only when we are ready to preserve backward compatibility for the exported API.

The first public preview should be `v0.1.0`, unless we decide to publish smaller incremental previews first. Do not create the tag until the release commit and its changelog have been reviewed.

## Changelog

Keep release notes in [`CHANGELOG.md`](CHANGELOG.md). Before a release, move the changes being released from `[Unreleased]` under a dated `## [0.x.y] - YYYY-MM-DD` heading, then create a fresh `[Unreleased]` section. Use the matching version section as the GitHub Release notes.

## Manual release steps

1. Start from an up-to-date `main` with a clean working tree. Confirm the CI workflow is green for the commit to release.
2. Review and update [`CHANGELOG.md`](CHANGELOG.md) as described above.
3. Push the release tag from the reviewed commit. For the initial preview, for example:

   ```bash
   git switch main
   git pull --ff-only
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

4. Confirm the tag workflow passes, then create a GitHub Release for that tag and use the matching changelog section as its release notes. Mark the release as a pre-release if it is explicitly an alpha or beta; a plain `v0.x.y` already communicates that the module is unstable.
5. Verify the version is available with `go list -m github.com/onprem/warera-go@v0.1.0` and that the README install command resolves to the tagged module.

Pushing a tag is the publication step for Go consumers. Do not rewrite or move a published tag. If a release needs correction, publish the next patch version.

## Future automation

The current project does not produce standalone binaries, so GoReleaser would mainly wrap GitHub Release creation and changelog generation; it does not decide the next module version or create a version-bump commit. Keep the first releases manual and reviewable. Once pull requests consistently use Conventional Commit titles, consider Release Please to open release pull requests with version bumps and changelog updates. Add GoReleaser only if we later want to publish actual build artifacts or have another concrete packaging need.
