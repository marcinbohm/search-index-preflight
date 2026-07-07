# Alpha readiness

This checklist tracks repository readiness for a future `v0.0.1-prealpha` release. It does not create a release and does not imply production readiness.

For the operational pre-release verification flow, see the [Release checklist](RELEASE_CHECKLIST.md).

## Ready now

- README explains the PR-time schema preflight positioning.
- `lint` is implemented for `SIL001`, `SIL002`, and `SIL003`.
- `diff` is implemented for `DIF001`, `DIF002`, and `DIF003`.
- `rules list` and `explain <RULE_ID>` expose public rule metadata.
- Console and JSON output work for current commands.
- Public fixtures use synthetic data.
- Getting-started docs and practical examples exist.
- Future config example is labeled as planned/not implemented.

## Not ready yet

- No release binaries are published.
- No GitHub Action is published.
- No SARIF or Markdown reporter is implemented.
- No config loading, suppressions, or baseline mode is implemented.
- No YAML input is implemented.
- No git-ref diff or rename detection is implemented.
- No read-only cluster doctor mode is implemented.
- Offline migration/versioning is ADR/concept only.

## Before v0.0.1-prealpha

- [ ] Confirm README demo commands stay stable.
- [ ] Confirm examples and fixtures are synthetic and small.
- [ ] Decide whether to publish source-only release notes before binaries.
- [ ] Verify issue templates and security reporting language.
- [ ] Set repository topics manually in GitHub.
- [ ] Confirm CI passes on the release commit.
- [ ] Keep `search-index-preflight.future.example.yaml` clearly labeled as future-only until config loading exists.

## Release checklist

Use [Release checklist](RELEASE_CHECKLIST.md) as the source of truth for pre-release verification.

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
