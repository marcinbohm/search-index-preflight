# Release checklist

This checklist is for preparing a future `v0.0.1-prealpha` release.

It does not create a release, tag, or artifact. SearchIndexPreflight remains experimental pre-alpha until stated otherwise.

## Release goal

The first pre-alpha release should make it easy for external users to evaluate the idea, run examples, inspect rule metadata, and give feedback. It is not a production CI-gate release.

## Release type

Recommended first release: source/tag release first, with binaries only if a release process is explicitly added and verified.

Binaries, Homebrew, Docker images, and release automation are future options. Do not require or announce them until the repository has verified release tooling.

## Do not release if

- `go test ./...` fails.
- `go vet ./...` fails.
- README demo commands do not match current CLI behavior.
- Examples produce output different from `expected-output.txt`.
- README or docs imply production readiness.
- README or docs imply config loading, suppressions, baseline mode, YAML, SARIF, GitHub Action, release binaries, or cluster doctor are implemented when they are not.
- SECURITY.md contains a TBD reporting channel.
- `search-index-preflight.example.yaml` appears as a working config example.
- Release notes claim artifacts that do not exist.

## Pre-release verification

Run these commands before release:

```bash
git status --short
git status --short --ignored .local .idea

gofmt -w .
go test ./internal/cli
go test ./...
go vet ./...

go run ./cmd/search-index-preflight --help
go run ./cmd/search-index-preflight version
go run ./cmd/search-index-preflight lint --help
go run ./cmd/search-index-preflight lint --mapping examples/basic/mapping.json
go run ./cmd/search-index-preflight diff --help
go run ./cmd/search-index-preflight rules list
go run ./cmd/search-index-preflight explain SIL001
```

## Expected non-zero commands

Some demo commands intentionally return non-zero because they emit error findings. `go run` prints `exit status 1` when the program returns exit code 1.

```bash
go run ./cmd/search-index-preflight diff \
  --base examples/field-type-change/base \
  --current examples/field-type-change/current
```

Expected:

- emits `DIF001`
- exits non-zero because `DIF001` is error

```bash
go run ./cmd/search-index-preflight diff \
  --base fixtures/diff/mixed-field-changes/base \
  --current fixtures/diff/mixed-field-changes/current
```

Expected:

- emits `DIF001`, `DIF002`, `DIF003`
- exits non-zero because `DIF001` is error

## Expected zero commands with findings

```bash
go run ./cmd/search-index-preflight diff \
  --base examples/field-removed/base \
  --current examples/field-removed/current

go run ./cmd/search-index-preflight lint \
  --mapping examples/dynamic-template-risk/mapping.json
```

Expected:

- `field-removed` emits `DIF002` warning and exits 0 under default `--fail-on error`
- `dynamic-template-risk` emits `SIL003` warning and exits 0 under default `--fail-on error`

## Documentation verification

- README status still says experimental pre-alpha.
- README has a quick demo and a failure example.
- README separates Works today from Planned, not implemented yet.
- `docs/GETTING_STARTED.md` commands are current.
- `docs/CLI_CONTRACT.md` matches current command behavior.
- `docs/RULE_CATALOG.md` lists `SIL001`-`SIL003` and `DIF001`-`DIF003`.
- `search-index-preflight.future.example.yaml` remains clearly marked as future-only.
- `CHANGELOG.md` does not claim a release exists before it is created.
- Draft release notes exist at `docs/releases/v0.0.1-prealpha.md`.

## Example verification

For examples that should exit 0:

```bash
go run ./cmd/search-index-preflight diff \
  --base examples/field-removed/base \
  --current examples/field-removed/current \
  > /tmp/search-index-preflight-field-removed.txt

diff -u examples/field-removed/expected-output.txt /tmp/search-index-preflight-field-removed.txt
```

For examples that exit 1, capture stdout without failing the whole shell script:

```bash
set +e
go run ./cmd/search-index-preflight diff \
  --base examples/field-type-change/base \
  --current examples/field-type-change/current \
  > /tmp/search-index-preflight-field-type-change.txt
code=$?
set -e

test "$code" -ne 0
diff -u examples/field-type-change/expected-output.txt /tmp/search-index-preflight-field-type-change.txt
```

Repeat the same pattern for other examples that intentionally emit error findings.

## Security and repository hygiene

- SECURITY.md does not contain a TBD reporting channel.
- Public examples and fixtures contain synthetic data only.
- `.local/` and `.idea/` remain ignored and unstaged.
- No credentials, production mappings, real logs, internal service names, or customer data are present in docs, examples, fixtures, or release notes.
- No release artifact is claimed unless it exists.
- No tag is created until release notes are final.

## Stale wording checks

Run:

```bash
rg -n "explain.*stub|explain.*planned|explain.*not implemented|explain implementation is in progress|explain full.*incomplete" README.md docs internal --glob '!.git/**'
rg -n "rules list.*stub|rules list.*planned|rules list.*not implemented|rules.*stub" README.md docs internal --glob '!.git/**'
rg -n "SearchIndexLint|search-index-lint|github.com/marcinbohm/search-index-lint" . --glob '!.git/**' --glob '!.local/**' --glob '!.idea/**'
rg -n "TBD|Preferred reporting channel" README.md docs SECURITY.md examples --glob '!.git/**'
rg -n "search-index-preflight.example.yaml" README.md docs examples .github --glob '!.git/**'
```

Notes:

- Historical ADR old-name references are acceptable.
- Future rule catalog TBD entries may be acceptable if clearly future/planned.
- No active README/current docs should refer to `rules list` or `explain` as stub.
- The old config filename may appear in this checklist as a stale-search target only.

## GitHub manual checks

- CI is green on `main`.
- Repository description says something like: `Offline-first preflight CLI for Elasticsearch/OpenSearch schema changes`.
- Topics are set:
  - elasticsearch
  - opensearch
  - search
  - schema
  - mapping
  - index-template
  - ci
  - preflight
  - linter
  - golang
  - devops
  - sre
  - schema-as-code
- Security advisories/private vulnerability reporting are enabled if available.
- Issue templates render correctly.
- README badges render correctly.
- No release artifacts are claimed unless they exist.

## Changelog update

Before tagging:

- Move relevant items from `Unreleased` into `v0.0.1-prealpha`.
- Keep a new empty `Unreleased` section above it.
- Do not claim binaries unless they are attached.
- Mention that this is experimental pre-alpha.
- Review `docs/releases/v0.0.1-prealpha.md` before copying it into GitHub release notes.
- Confirm `go run ./cmd/search-index-preflight version` prints `SearchIndexPreflight version 0.0.1-prealpha`.

## Tag and release notes

Do not tag until verification passes and release notes are final.

### Release notes template

Title: `v0.0.1-prealpha`

Summary:
SearchIndexPreflight is an experimental offline-first CLI for reviewing Elasticsearch/OpenSearch mappings, templates, and schema diffs before they reach production.

Included:

- static lint rules: SIL001-SIL003
- experimental diff rules: DIF001-DIF003
- `rules list`
- `explain <RULE_ID>`
- console and JSON output
- examples and fixtures

Not included:

- production-ready CI gate guarantees
- release binaries, unless attached
- config loading
- suppressions
- baseline mode
- YAML input
- SARIF/Markdown output
- cluster doctor
- migration/versioning implementation

## After release

- Confirm the README release references still match the published release.
- Confirm the changelog has a fresh `Unreleased` section.
- Confirm no planned-only feature is described as implemented.
- Open follow-up issues for release artifacts or package managers only if maintainers want that next.

## Rollback / correction notes

If a release note overclaims behavior, publish a correction quickly and update README/docs.

If a tag is created by mistake, coordinate with maintainers before deleting or replacing it. Prefer a corrective follow-up release when possible.
