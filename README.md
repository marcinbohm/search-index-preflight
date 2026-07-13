# SearchIndexPreflight

[![status: pre-alpha](https://img.shields.io/badge/status-pre--alpha-orange)](#current-status)
[![CI](https://github.com/marcinbohm/search-index-preflight/actions/workflows/ci.yml/badge.svg)](https://github.com/marcinbohm/search-index-preflight/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/marcinbohm/search-index-preflight?include_prereleases)](https://github.com/marcinbohm/search-index-preflight/releases)
[![Go 1.22](https://img.shields.io/badge/go-1.22-blue)](go.mod)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)

Catch Elasticsearch/OpenSearch schema changes that are valid, but risky.

SearchIndexPreflight is an offline-first CLI for reviewing Elasticsearch/OpenSearch mappings, templates, and schema diffs before they reach production. It is aimed at PR-time feedback for teams that keep search schemas in Git.

## Why

Elasticsearch and OpenSearch mappings can be syntactically valid and still create rollout risk:

- field type changes that require reindexing or rollover planning
- removed fields that break producers, queries, dashboards, or alerts
- added fields that increase field counts or change downstream assumptions
- dynamic mappings and dynamic templates that expand more broadly than intended

SearchIndexPreflight moves part of that feedback into local development, CI, and pull requests without requiring cluster credentials.

## Current status

SearchIndexPreflight is experimental pre-alpha.

Use it for:

- evaluating the idea
- running local demos against synthetic fixtures
- reviewing early CLI behavior
- giving feedback on rule usefulness

Do not use it yet as:

- a required production CI gate
- a compliance tool
- a replacement for staging validation

## Quick demo

Run a mixed schema diff fixture:

```bash
go run ./cmd/search-index-preflight diff \
  --base fixtures/diff/mixed-field-changes/base \
  --current fixtures/diff/mixed-field-changes/current
```

Short output excerpt:

```text
error DIF001: mapping.json#/properties/status: Field "status" changed type from "keyword" to "long".
warning DIF002: mapping.json#/properties/legacy_id: Field "legacy_id" was removed from the current schema.
info DIF003: mapping.json#/properties/customer_id: Field "customer_id" was added to the current schema.
```

## Show me the failure

`DIF001` is an error, so the mixed fixture exits non-zero by default. That is intentional: a mapped field type change is often a schema-breaking change and usually needs a migration plan.

To see a static lint warning:

```bash
go run ./cmd/search-index-preflight lint \
  --mapping fixtures/dynamic-templates/sil003-missing-match-mapping-type/mapping-missing-match-mapping-type.json
```

```text
warning SIL003: ... Dynamic template "strings_as_keywords" does not declare match_mapping_type. It may apply more broadly than intended.
```

## Works today

- `lint` for offline static checks
- `diff --base <path> --current <path>` for schema diff checks
- `rules list` for public rule metadata
- `explain <RULE_ID>` for one-rule explanations
- console and JSON output
- JSON mappings/templates
- JSONL/NDJSON sample parsing
- static rules: `SIL001`, `SIL002`, `SIL003`
- diff rules: `DIF001`, `DIF002`, `DIF003`

## Planned, not implemented yet

- `check` alias / future preferred static-check command
- YAML input
- config loading
- suppressions
- baseline mode
- SARIF output
- Markdown output
- GitHub Action
- release binaries
- Homebrew / Docker packaging
- git refs and rename detection for `diff`
- read-only cluster `doctor`
- offline migration/versioning validation

The future migration/versioning layer is planned as offline preflight only. It will not apply migrations, write to clusters, execute alias cutovers, reindex, roll back, or store migration state in Elasticsearch/OpenSearch.

## Install / run from source

Latest pre-release: [`v0.0.1-prealpha`](https://github.com/marcinbohm/search-index-preflight/releases/tag/v0.0.1-prealpha).

This is currently a source-only pre-alpha release. Binary release artifacts are not published yet.

Requirements:

- Go 1.22 or newer

From a clone:

```bash
go test ./...
go run ./cmd/search-index-preflight --help
go run ./cmd/search-index-preflight lint --mapping examples/basic/mapping.json
```

Install the current source-only pre-release with Go:

```bash
go install github.com/marcinbohm/search-index-preflight/cmd/search-index-preflight@v0.0.1-prealpha
search-index-preflight version
```

Expected:

```
SearchIndexPreflight version 0.0.1-prealpha
```

## Commands

```bash
search-index-preflight lint --mapping mapping.json
search-index-preflight lint ./schemas
search-index-preflight diff --base old-schemas/ --current new-schemas/
search-index-preflight rules list
search-index-preflight rules list --family diff --format json
search-index-preflight explain SIL001
search-index-preflight explain DIF003 --format json
search-index-preflight version
```

`diff` currently matches directory inputs by relative path. Explicit file-vs-file inputs are compared as one logical resource even when filenames differ. Rename detection is not implemented, so renamed schema files may appear as removals from the old path and additions from the new path.

## Examples

- [Field type change](examples/field-type-change/README.md): emits `DIF001`.
- [Field removed](examples/field-removed/README.md): emits `DIF002`.
- [Dynamic template risk](examples/dynamic-template-risk/README.md): emits `SIL003`.
- [Diff fixtures](fixtures/diff/README.md): small public fixtures for `DIF001`, `DIF002`, `DIF003`, no-change, and mixed-change behavior.

## Documentation

- [Getting started](docs/GETTING_STARTED.md)
- [CLI contract](docs/CLI_CONTRACT.md)
- [Rule catalog](docs/RULE_CATALOG.md)
- [Fixtures](docs/FIXTURES.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Product brief](docs/PRODUCT_BRIEF.md)
- [Alpha readiness](docs/ALPHA_READINESS.md)
- [Release checklist](docs/RELEASE_CHECKLIST.md)
- [Offline migration/versioning concept](docs/MIGRATION_VERSIONING_CONCEPT.md)
- [Changelog](CHANGELOG.md)

Maintainer workflow notes live under [docs/maintainer](docs/maintainer/README.md).

## Safety model / non-goals

SearchIndexPreflight is offline-first by default.

It does not currently:

- connect to Elasticsearch/OpenSearch clusters
- make network calls during `lint` or `diff`
- upload mappings or sample documents
- send telemetry
- write to clusters
- auto-fix schemas

Future cluster-facing work must be explicitly invoked and read-only.

## Contributing

Contributions are welcome while the project is still shaping its alpha surface. Start with [CONTRIBUTING.md](CONTRIBUTING.md), use synthetic data only, and keep PRs small.

## License

Apache-2.0. See [LICENSE](LICENSE).
