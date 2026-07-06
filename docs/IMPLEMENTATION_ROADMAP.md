# Implementation Roadmap

## Roadmap principles

- Documentation before code.
- Fixtures before or alongside rules.
- MVP must be small and sharp.
- No UI.
- No SaaS.
- No auto-fix.
- No live cluster mode in MVP.
- No writes to clusters ever.
- SARIF and GitHub Action are alpha/beta scope.
- Baseline is beta scope.
- Rule IDs must be stable.
- Heuristics must be conservative.

## Current Implementation Status

Completed foundations:

- M0 project skeleton
- M1 parse-only lint foundation
- M1.5 canonical model and traversal helpers
- M2 rule registry and runner foundation

Current CLI behavior:

- `lint` reports parse and normalization diagnostics only
- no real SIL rules are implemented
- rule execution is not wired into the CLI

Next:

- first real rule implementation, likely SIL001, with focused tests and fixtures

## Pre-alpha

Goal: prove parser, canonical model, rule engine, and fixture-driven development.

Scope:

- Go CLI skeleton
- JSON/JSONL parsing foundation
- YAML parsing
- mapping/index/component template models
- sample schema model
- rule registry
- 8-10 high-signal rules
- console and JSON output
- fixture harness
- golden report tests

Done criteria:

- `go test ./...` passes
- `search-index-lint lint fixtures/... --format json` works
- every implemented rule has fixtures
- invalid JSON/YAML returns clean error
- no network access
- no cluster code path
- no panics on malformed fixtures
- README accurately reflects implemented status

Out of scope: SARIF, GitHub Action, baseline, diff, cluster mode, auto-fix, plugin API, Docker image.

## MVP

Goal: deliver a small, usable offline CLI.

Scope:

- complete parser coverage
- directory input detection
- config file
- suppressions
- severity thresholds
- JSON report schema
- console UX
- 12-15 MVP rules
- fixtures
- `rules list`
- `explain`

Done criteria:

- CLI documented and tested
- config works
- suppressions require reason
- output deterministic
- at least 15 fixtures
- at least 12 implemented rules
- no cluster requirement
- docs match behavior

## Alpha

Goal: make SearchIndexLint useful in pull requests.

Scope:

- Markdown reporter
- SARIF reporter
- GitHub Action wrapper
- rule docs
- improved source locations
- sample doc inference improvements
- initial compatibility profile
- 15-20 rules

Done criteria:

- SARIF validates
- GitHub Action works on sample repo
- Markdown output readable
- at least 20 fixtures
- external users can run from README

## Beta

Goal: make adoption practical for legacy repositories.

Scope:

- baseline create/filtering
- diff old/new schema directories
- stable fingerprints
- dialect/version capability matrix
- compatibility fixtures
- fuzz tests
- benchmarks
- false-positive tuning
- 25-30 rules

Done criteria:

- baseline workflow documented
- diff supports basic mapping changes
- compatibility profiles documented
- no panics on fuzz corpus
- release artifacts available

## v1

Goal: stable public infra tool.

Scope:

- stable CLI contract
- stable JSON report schema
- stable baseline format
- stable rule ID policy
- signed releases/checksums
- GitHub Action v1
- complete docs
- complete contribution process
- 25+ documented rules

Done criteria:

- all default-on rules have fixtures
- all default-on rules have docs
- SemVer documented
- changelog complete
- README accurate
- SECURITY.md complete

## v1.1 ideas

- read-only cluster inspect
- read-only template simulate
- `_field_caps` integration
- drift detection
- Docker image hardening
- more OpenSearch compatibility rules

## v2 ideas

Only if v1 has adoption:

- custom rule API
- external rule packs
- query workload hints
- Terraform/Kubernetes integration
- editor/LSP integration
- migration advisory mode
