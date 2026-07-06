# SearchIndexLint

Catch Elasticsearch/OpenSearch index schema and template risks before production.

SearchIndexLint is an offline-first CLI for linting Elasticsearch and OpenSearch mappings, component templates, index templates, dynamic templates, and sample documents.

It is designed for teams that treat search schemas as code and want PR-time feedback before risky changes reach production.

## Status

Status: pre-alpha / foundation phase.

The current CLI is not production-ready. It can parse and normalize JSON mappings/templates and JSONL/NDJSON sample documents, run the first built-in rule (`SIL001`), and report diagnostics/findings. Rule coverage is intentionally very limited.

## Problem statement

Elasticsearch and OpenSearch mappings can be syntactically valid and operationally dangerous.

Common rollout failures include:

- uncontrolled field growth from dynamic mappings
- field type conflicts across indices or templates
- broad dynamic templates that silently map too much
- dotted field collisions such as `user.name` vs `user: { name: ... }`
- index template priority collisions
- `text` fields later used for sorting or aggregations
- long `keyword` values without `ignore_above`
- analyzer or normalizer references that are not defined
- sample documents that do not match declared mappings
- object vs nested modeling mistakes
- dialect differences between Elasticsearch and OpenSearch

These issues often show up after data starts flowing. Remediation may require rollover, reindexing, producer changes, query changes, DLQ drain, or manual operator intervention.

SearchIndexLint moves part of that feedback earlier: into local development, CI, and pull requests.

## Who SearchIndexLint is for

SearchIndexLint is intended for:

- search infrastructure engineers
- platform engineers running Elasticsearch or OpenSearch clusters
- backend engineers shipping event, log, or search schemas
- SREs responsible for indexing reliability
- data platform teams maintaining schema-as-code repositories
- teams using GitOps, Terraform, CI/CD, or release pipelines for search schema changes

## Why existing APIs and tools are not enough

Elasticsearch and OpenSearch provide useful APIs such as template simulation, mapping retrieval, and field capability inspection. Those APIs are valuable, but they do not fully solve the pre-merge problem.

SearchIndexLint exists because:

- CI jobs should not require production cluster credentials.
- PR reviewers need offline feedback before templates are applied.
- Template simulation does not explain broader risk patterns by itself.
- Field capabilities are useful for existing indices, not proposed schema changes.
- Mapping generators create mappings; they do not evaluate rollout risk.
- Dashboards and cluster diagnostics are usually post-deploy tools.
- Teams need explainable rule IDs, suppressions, baselines, fixtures, and machine-readable reports.

SearchIndexLint is not a replacement for vendor APIs. Later versions may use read-only cluster APIs for drift detection and stronger validation. MVP remains offline.

## Installation

No release artifacts exist yet.

Planned future options:

```bash
brew install search-index-lint
go install github.com/marcinbohm/search-index-lint/cmd/search-index-lint@latest
```

## Local Development

```bash
go test ./...
go vet ./...
go run ./cmd/search-index-lint --help
go run ./cmd/search-index-lint lint --mapping examples/basic/mapping.json
```

## Current Usage

Current pre-alpha behavior parses JSON mappings/templates and JSONL/NDJSON sample documents, normalizes supported schema shapes into internal models, runs `SIL001`, then reports parse/normalization diagnostics and rule findings.

```bash
search-index-lint lint --mapping mapping.json
search-index-lint lint --template index-template.json
search-index-lint lint --component-template component-template.json
search-index-lint lint --sample-docs samples.jsonl
search-index-lint lint --sample-docs samples.ndjson
search-index-lint lint --mapping mapping.json --sample-docs samples.jsonl
search-index-lint lint ./schemas
search-index-lint rules list
search-index-lint explain SIL001
```

Directory mode currently discovers only `.json`, `.jsonl`, and `.ndjson` files.

## Example output

Current `SIL001` finding example:

```bash
search-index-lint lint --mapping fixtures/mapping-limits/sil001-total-fields-limit/mapping-over-limit.json
```

```text
error SIL001: fixtures/mapping-limits/sil001-total-fields-limit/mapping-over-limit.json#/: Mapping has 1000 normalized fields, exceeding the default total fields limit of 1000.
  Remediation: Reduce explicit field count, restrict dynamic mappings, consider flattened/flat_object only when query semantics fit, split unrelated data into separate indices, or raise index.mapping.total_fields.limit only with operational review.
```

Planned future multi-rule output shape:

```text
SearchIndexLint 0.1.0

Dialect: elasticsearch 8.x
Scanned: 6 files
Findings: 1 error, 2 warnings

ERROR SIL015 template-priority-conflict
  schemas/logs-app.index-template.json#/priority

  Templates "logs-default" and "logs-app" both match "logs-app-*"
  with the same priority 100.

  Remediation:
  Give the more specific template a higher priority or narrow the index pattern.

WARNING SIL004 overbroad-dynamic-template
  schemas/logs-app.index-template.json#/template/mappings/dynamic_templates/0

WARNING SIL012 long-keyword-without-ignore-above
  schemas/events.mapping.json#/mappings/properties/payload
```

## Current implemented scope

Included now:

- Go CLI skeleton
- `version`, `lint`, `rules list`, and `explain` commands
- JSON mapping/template parsing
- JSONL/NDJSON sample document parsing with line-number diagnostics
- directory discovery for `.json`, `.jsonl`, and `.ndjson`
- parse diagnostics
- normalization diagnostics
- canonical mapping model foundation
- index template model foundation
- component template model foundation
- canonical corpus model
- normalized field traversal helpers
- explicit file loading and directory discovery
- `.local/` and default build/vendor directory ignores during discovery
- severity, confidence, finding, diagnostic, and summary models
- console and JSON reports
- rule registry skeleton
- rule runner foundation
- first built-in rule: `SIL001` total fields limit risk

Not implemented yet:

- SIL002 and the rest of the rule catalog
- YAML parsing
- Markdown reporter
- SARIF reporter
- GitHub Action
- baseline mode
- diff mode
- cluster mode
- auto-fix
- config loading
- suppressions
- releases, Homebrew formula, and Docker image

## Planned MVP scope

Included:

- offline CLI
- JSON and YAML input parsing
- JSONL sample document parsing
- mapping parser
- index template parser
- component template parser
- basic template composition approximation from supplied files
- canonical mapping model
- rule engine
- deterministic and heuristic finding model
- severity and confidence model
- console and JSON reports
- initial rule catalog focused on high-signal risks
- public fixture corpus
- golden file tests
- config file support
- local suppressions
- clear documentation

Initial MVP rules:

- SIL001 total-fields-limit-risk (implemented, limited default-threshold version)
- SIL002 root-dynamic-enabled
- SIL003 dynamic-template-missing-match-mapping-type
- SIL004 overbroad-dynamic-template
- SIL005 dynamic-template-shadowing
- SIL006 path-match-object-collision-risk
- SIL007 dotted-field-collision
- SIL008 field-type-conflict
- SIL009 sample-doc-mapping-conflict
- SIL010 dynamic-date-numeric-detection-risk
- SIL011 likely-aggregatable-field-as-text
- SIL012 long-keyword-without-ignore-above
- SIL013 fielddata-true-on-text
- SIL014 missing-analyzer-normalizer-definition
- SIL015 template-priority-conflict

## Non-goals

SearchIndexLint v1 is not:

- a dashboard
- a SaaS product
- an OpenSearch Dashboards plugin
- a mapping generator
- an automatic fixer
- a cluster doctor
- a replacement for Elasticsearch/OpenSearch APIs
- a replacement for load testing
- a replacement for staging clusters
- a replacement for operator judgment
- a query performance oracle
- an ingestion pipeline simulator
- a tool that requires production data
- a tool that writes to clusters

Live cluster mode is not part of the MVP. If added later, it must be read-only and explicitly separated from offline linting.

## Roadmap

### Pre-alpha

- Go CLI skeleton
- JSON/JSONL parsing
- YAML parsing
- canonical mapping/template model
- initial rule registry
- 8-10 high-signal rules
- console and JSON output
- fixture harness
- golden report tests

### Alpha

- 15-20 rules
- config file
- suppressions
- Markdown reporter
- SARIF reporter
- GitHub Action wrapper
- rule docs
- compatibility profile skeleton

### Beta

- baseline mode
- improved source locations
- richer sample document inference
- dialect/version capability matrix
- more compatibility fixtures
- better false-positive tuning
- release binaries

### v1

- stable CLI contract
- stable JSON report schema
- stable rule ID policy
- stable GitHub Action
- 25+ documented rules
- SemVer
- changelog
- signed release artifacts
- documented security policy
- documented contribution process

## Project Phase

Current phase:

- project skeleton
- input discovery
- JSON/JSONL parser foundations
- canonical model and normalizer foundations
- canonical corpus model
- normalized field traversal helpers
- rule runner foundation
- first built-in rule: `SIL001`
- parse, normalization, and rule finding reporting

No production release exists yet.

## Feedback wanted

SearchIndexLint is looking for feedback on:

- real Elasticsearch/OpenSearch schema failure cases
- rule false positives
- rule false negatives
- missing OpenSearch compatibility checks
- SARIF/GitHub Action workflow expectations
- fixture examples that can be shared publicly
- confusing remediation guidance

Please do not share confidential mappings, production logs, customer data, internal templates, or private cluster details in public issues.

## Safety warning

SearchIndexLint detects risk. It does not prove a schema is safe.

It does not replace staging validation, integration tests, load tests, rollout planning, cluster observability, experienced operator review, vendor documentation, or incident response judgment.

A clean current pre-alpha report means only that parsing and normalization completed without diagnostics and the currently implemented limited rule set did not report findings. It does not mean the full schema is safe.
