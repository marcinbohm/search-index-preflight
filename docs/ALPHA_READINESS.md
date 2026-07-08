# Alpha readiness

This document tracks what is ready after `v0.0.1-prealpha` and what remains before the first usable alpha. It does not imply production readiness.

For the operational release verification flow, see the [Release checklist](RELEASE_CHECKLIST.md).

## Ready now

- README explains the PR-time schema preflight positioning.
- `lint` is implemented for `SIL001`, `SIL002`, and `SIL003`.
- `diff` is implemented for `DIF001`, `DIF002`, and `DIF003`.
- `rules list` and `explain <RULE_ID>` expose public rule metadata.
- Console and JSON output work for current commands.
- Public fixtures use synthetic data.
- Getting-started docs and practical examples exist.
- Future config example is labeled as planned/not implemented.
- `v0.0.1-prealpha` exists as a source-only GitHub pre-release.
- Release notes exist at `docs/releases/v0.0.1-prealpha.md`.
- Manual GitHub repository topics have been configured.

## Not ready yet

- No release binaries are published.
- No GitHub Action is published.
- No SARIF or Markdown reporter is implemented.
- No config loading, suppressions, or baseline mode is implemented.
- No YAML input is implemented.
- No git-ref diff or rename detection is implemented.
- No read-only cluster doctor mode is implemented.
- Offline migration/versioning is ADR/concept only.

## Completed for v0.0.1-prealpha

- [x] README demo commands are stable for the pre-release.
- [x] Examples and fixtures are synthetic and small.
- [x] Source-only release notes were published before binaries.
- [x] Issue templates and security reporting language were checked.
- [x] Repository topics were set manually in GitHub.
- [x] CI passed on the release commit.
- [x] `search-index-preflight.future.example.yaml` is clearly labeled as future-only until config loading exists.

## Before the first usable alpha

- [ ] Decide whether release binaries or package-manager installs are needed.
- [ ] Decide whether a GitHub Action wrapper is in scope.
- [ ] Add SARIF or Markdown output only if the reporting surface is designed.
- [ ] Add config loading, suppressions, baseline mode, or YAML only after the behavior is specified.
- [ ] Keep migration/versioning as offline preflight only unless a future ADR changes that direction.

## Release checklist

Use [Release checklist](RELEASE_CHECKLIST.md) as the source of truth for release verification.

- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] CLI smoke tests for `lint`, `diff`, `rules list`, `explain`, and `version`
- [ ] `CHANGELOG.md` updated
- [ ] README status still says pre-alpha
- [ ] No release artifacts are claimed unless they exist
- [ ] No cluster-writing functionality exists

## Visibility checklist

Recommended GitHub topics:

```text
elasticsearch
opensearch
search
schema
mapping
index-template
ci
preflight
linter
golang
devops
sre
schema-as-code
```
