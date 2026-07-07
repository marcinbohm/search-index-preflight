# Contributing to SearchIndexPreflight

SearchIndexPreflight is an offline-first Elasticsearch/OpenSearch index schema risk linter. Contributions are welcome, but the project is strict about scope, fixtures, tests, and public-data safety.

## Project scope

In scope:

- offline CLI
- mappings
- index templates
- component templates
- dynamic templates
- sample documents
- rule engine
- reports
- fixtures
- documentation
- GitHub Action/SARIF after core CLI is stable

Out of scope before v1:

- dashboard
- SaaS
- auto-fix
- live cluster mode
- cluster write operations
- OpenSearch Dashboards plugin
- custom plugin API

## Local development

Requirements:

- Go version: see `go.mod`
- Make: optional, TBD
- golangci-lint: optional before alpha, recommended later

Clone:

```bash
git clone https://github.com/marcinbohm/search-index-preflight.git
cd search-index-preflight
```

Run tests:

```bash
go test ./...
```

Run locally:

```bash
go run ./cmd/search-index-preflight --help
go run ./cmd/search-index-preflight lint --mapping examples/basic/mapping.json
```

Format:

```bash
gofmt -w .
```

Vet:

```bash
go vet ./...
```

## How to add a rule

Real SIL rules are not implemented yet. Coordinate first-rule work with maintainers before adding rule behavior.

1. Open or create a rule request issue.
2. Confirm the rule fits the roadmap.
3. Pick a stable rule ID assigned by maintainers.
4. Add rule metadata.
5. Implement the rule.
6. Add unit tests.
7. Add positive fixture.
8. Add negative fixture.
9. Add expected JSON output.
10. Add suppression fixture if applicable.
11. Update `docs/RULE_CATALOG.md`.
12. Update `search-index-preflight rules list` metadata if generated manually.
13. Update README or CLI docs if behavior changes.

Rule implementation must not read files directly, format reports, call network APIs, depend on CLI flags directly, use production data, or panic on malformed input.

## Rule severity guidance

Use `error` for deterministic issues likely to break indexing/template application with clear fixes.  
Use `warning` for risky but context-dependent issues.  
Use `info` for advisory or high false-positive-risk findings.  
Use `critical` rarely.

## How to add a fixture

Create:

```text
fixtures/<category>/<case>/
  README.md
  fixture.yaml
  mapping.json              # optional
  index-template.json       # optional
  component-template.json   # optional
  samples.jsonl             # optional
  expected.json
```

Rules:

- use synthetic data only
- keep examples small
- isolate one main issue
- include expected output
- include remediation in README
- do not copy proprietary mappings
- do not include production logs
- do not include customer data

## Style guide

Go:

- use `gofmt`
- prefer small packages
- avoid global mutable state
- return errors with context
- avoid `panic` for user input
- use table-driven tests
- sort output deterministically
- do not use `map[string]any` directly in rules when a typed model exists

Docs:

- be direct
- avoid marketing language
- include examples
- mark unknowns as `TBD`
- do not invent references
- distinguish deterministic errors from heuristic advice

## PR process

Before opening a PR:

- run `go test ./...`
- run `gofmt`
- update docs
- add fixtures
- check for private data
- fill the PR template

PRs should be small.

## Issue labels

Recommended labels:

- `type: implementation`
- `type: docs`
- `type: fixture`
- `type: rule`
- `type: bug`
- `type: feature`
- `type: false-positive`
- `type: compatibility`
- `type: good-first-issue`
- `area: cli`
- `area: parser`
- `area: rules`
- `area: report`
- `area: fixtures`
- `area: docs`
- `stage: mvp`
- `stage: alpha`
- `stage: beta`
- `stage: v1`

## Security and privacy

Do not include production mappings, customer data, real logs, credentials, tokens, internal service names, or confidential index patterns.

Read `SECURITY.md` before reporting security issues or sharing examples.
