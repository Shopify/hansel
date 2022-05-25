# Development

* For tests, use standard Go tooling: `go test ./...`
* For lint, use [golangci-lint](https://golangci-lint.run/): `golangci-lint run ./...`
* For integration tests, use Docker to provide test environments: `docker build -f Dockerfile.test .`
* Hansel uses [goreleaser](https://goreleaser.com/), but this is not required for development.

# Releasing

Releases are built and published from GitHub Actions. Release versions follow semver.

To trigger a release:
1. Determine the appropriate version increment, as a rule of thumb:
   * If removing a CLI argument, increment the major version.
   * If adding a CLI argument, increment the minor version.
   * Else, increment the patch version.
1. Check the [release history](https://github.com/Shopify/hansel/releases) and increment to determine the next version. e.g. a patch increment to `v1.2.2` would be `v1.2.3`.
1. Create a tag with the next version, push the tag:
```bash
git checkout main
git pull origin main
git log -1 # double-check the expected commit
git tag -a v1.2.3 -m "v1.2.3 release"
git push --tags
```
4. Monitor the release process from [GitHub Actions](https://github.com/Shopify/hansel/actions/workflows/release.yaml).
