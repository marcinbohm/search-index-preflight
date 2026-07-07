# SearchIndexPreflight

Catch Elasticsearch/OpenSearch index schema and template risks before production.

SearchIndexPreflight is an offline-first preflight CLI for Elasticsearch and OpenSearch schema changes.

It is designed for teams that treat search schemas as code and want PR-time feedback before risky changes reach production.

## Status

Status: pre-alpha / foundation phase.

The current CLI is not production-ready. It can parse and normalize JSON mappings/templates and JSONL/NDJSON sample documents, run the first built-in static rules (`SIL001`, `SIL002`, `SIL003`), run the first diff rules (`DIF001`, `DIF002`, `DIF003`), and report diagnostics/findings. Rule coverage is intentionally very limited.

Current: `lint` static checks over supplied schema files plus a minimal experimental `diff` command for field type changes, removed fields, and added fields.  
Next strategic direction: deeper diff/preflight analysis for schema changes before merge or deployment.

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

SearchIndexPreflight moves part of that feedback earlier: into local development, CI, and pull requests.

## Who SearchIndexPreflight is for

SearchIndexPreflight is intended for:

- search infrastructure engineers
- platform engineers running Elasticsearch or OpenSearch clusters
- backend engineers shipping event, log, or search schemas
- SREs responsible for indexing reliability
- data platform teams maintaining schema-as-code repositories
- teams using GitOps, Terraform, CI/CD, or release pipelines for search schema changes

## Why existing APIs and tools are not enough

Elasticsearch and OpenSearch provide useful APIs such as template simulation, mapping retrieval, and field capability inspection. Those APIs are valuable, but they do not fully solve the pre-merge problem.

SearchIndexPreflight exists because:

- CI jobs should not require production cluster credentials.
- PR reviewers need offline feedback before templates are applied.
- Template simulation does not explain broader risk patterns by itself.
- Field capabilities are useful for existing indices, not proposed schema changes.
- Mapping generators create mappings; they do not evaluate rollout risk.
- Dashboards and cluster diagnostics are usually post-deploy tools.
- Teams need explainable rule IDs, suppressions, baselines, fixtures, and machine-readable reports.

SearchIndexPreflight is not a replacement for vendor APIs. Later versions may use read-only cluster APIs for drift detection and stronger validation. MVP remains offline.

## Installation

No release artifacts exist yet.

Planned future options:

```bash
brew install search-index-preflight
go install github.com/marcinbohm/search-index-preflight/cmd/search-index-preflight@latest
```

## Local Development

```bash
go test ./...
go vet ./...
go run ./cmd/search-index-preflight --help
go run ./cmd/search-index-preflight lint --mapping examples/basic/mapping.json
```

## Current Usage

Current pre-alpha behavior parses JSON mappings/templates and JSONL/NDJSON sample documents, normalizes supported schema shapes into internal models, runs `SIL001`, `SIL002`, and `SIL003` for static linting, and can run `DIF001`, `DIF002`, and `DIF003` through the experimental `diff` command.

```bash
search-index-preflight lint --mapping mapping.json
search-index-preflight lint --template index-template.json
search-index-preflight lint --component-template component-template.json
search-index-preflight lint --sample-docs samples.jsonl
search-index-preflight lint --sample-docs samples.ndjson
search-index-preflight lint --mapping mapping.json --sample-docs samples.jsonl
search-index-preflight lint ./schemas
search-index-preflight diff --base old-schemas/ --current new-schemas/
search-index-preflight diff --base fixtures/diff/dif001-field-type-changed/base --current fixtures/diff/dif001-field-type-changed/current
search-index-preflight diff --base fixtures/diff/dif002-field-removed/base --current fixtures/diff/dif002-field-removed/current
search-index-preflight diff --base fixtures/diff/dif003-field-added/base --current fixtures/diff/dif003-field-added/current
search-index-preflight diff --base fixtures/diff/mixed-field-changes/base --current fixtures/diff/mixed-field-changes/current
search-index-preflight rules list
search-index-preflight rules list --family diff --format json
search-index-preflight explain SIL001
search-index-preflight explain DIF003 --format json
```

Directory mode currently discovers only `.json`, `.jsonl`, and `.ndjson` files.

The experimental `diff` command currently emits `DIF001` field type changes, `DIF002` field removals, and `DIF003` field additions. `DIF002` is a warning and `DIF003` is info, so neither fails by default with `--fail-on error`; use `--fail-on warning` to fail on removed fields or `--fail-on info` to fail on added fields. Diff does not support git refs, PR comments, settings/alias diffs, dynamic template diffs, composed template analysis, sample document comparison, or cluster-backed validation.

Diff matching is intentionally simple: explicit file-vs-file inputs are compared as one logical resource even when filenames differ, while directory-vs-directory inputs are matched by relative path. File-vs-directory behavior is path-based and limited. Because rename detection is not implemented, renamed schema files may be reported as removed fields from the old relative path and added fields from the new relative path rather than matched as the same resource.

## Planned Direction

Planned future modes:

```bash
search-index-preflight check ./schemas      # planned future preferred static-check command
search-index-preflight diff --base old/ --current new/  # minimal experimental command exists
search-index-preflight doctor --url http://localhost:9200 --pattern "logs-*"  # planned later, read-only
```

`check` and `doctor` are not implemented today. The current static-check command remains:

```bash
search-index-preflight lint ./schemas
```

Future direction also includes offline migration/versioning validation for schema evolution. That future layer is planned as a helper/preflight workflow only: no cluster writes, no apply command, no alias cutover execution, no reindex execution, and no rollback execution.

## Example output

Current `SIL001` finding example:

```bash
search-index-preflight lint --mapping fixtures/mapping-limits/sil001-total-fields-limit/mapping-over-limit.json
```

```text
error SIL001: fixtures/mapping-limits/sil001-total-fields-limit/mapping-over-limit.json: Mapping has 1000 normalized fields, exceeding the default total fields limit of 1000.
  Remediation: Reduce explicit field count, restrict dynamic mappings, consider flattened/flat_object only when query semantics fit, split unrelated data into separate indices, or raise index.mapping.total_fields.limit only with operational review.
```

Planned future multi-rule output shape:

```text
SearchIndexPreflight 0.1.0

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
- `version`, `lint`, `diff`, `rules list`, and `explain` commands
- public rule metadata listing for lint and diff rules through `rules list`
- public rule metadata explanation for one lint or diff rule through `explain`
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
- second built-in rule: `SIL002` root dynamic enabled heuristic warning
- third built-in rule: `SIL003` dynamic template missing match mapping type heuristic warning
- experimental `diff --base <path> --current <path>` command
- first diff rule: `DIF001` field type changed
- second diff rule: `DIF002` field removed
- third diff rule: `DIF003` field added

Not implemented yet:

- SIL004 and the rest of the rule catalog
- deeper diff/preflight analysis beyond DIF001/DIF002/DIF003
- offline migration/versioning validation
- YAML parsing
- Markdown reporter
- SARIF reporter
- GitHub Action
- baseline mode
- git-aware diff mode
- cluster mode
- auto-fix
- config loading
- suppressions
- releases, Homebrew formula, and Docker image

Strategic note: further state-only heuristic rule expansion is paused unless explicitly approved. The next major product track is expected to be diff/preflight foundation.

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
- SIL002 root-dynamic-enabled (implemented, root-level explicit `dynamic: true` only)
- SIL003 dynamic-template-missing-match-mapping-type (implemented, missing `match_mapping_type` only)
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

SearchIndexPreflight v1 is not:

- a dashboard
- a SaaS product
- an OpenSearch Dashboards plugin
- a mapping generator
- an automatic fixer
- a cluster writer or automatic cluster remediator
- a replacement for Elasticsearch/OpenSearch APIs
- a replacement for load testing
- a replacement for staging clusters
- a replacement for operator judgment
- a query performance oracle
- an ingestion pipeline simulator
- a tool that requires production data
- a tool that writes to clusters

Live cluster mode is not part of the MVP. Future doctor mode, if added, must be read-only and explicitly separated from offline linting.

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
- built-in rules: `SIL001`, `SIL002`, `SIL003`
- parse, normalization, and rule finding reporting

No production release exists yet.

## Feedback wanted

SearchIndexPreflight is looking for feedback on:

- real Elasticsearch/OpenSearch schema failure cases
- rule false positives
- rule false negatives
- missing OpenSearch compatibility checks
- SARIF/GitHub Action workflow expectations
- fixture examples that can be shared publicly
- confusing remediation guidance

Please do not share confidential mappings, production logs, customer data, internal templates, or private cluster details in public issues.

## Safety warning

SearchIndexPreflight detects risk. It does not prove a schema is safe.

It does not replace staging validation, integration tests, load tests, rollout planning, cluster observability, experienced operator review, vendor documentation, or incident response judgment.

A clean current pre-alpha report means only that parsing and normalization completed without diagnostics and the currently implemented limited rule set did not report findings. It does not mean the full schema is safe.
