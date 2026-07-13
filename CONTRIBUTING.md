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

SIL001, SIL002, and SIL003 are implemented. Coordinate new rule work with maintainers before adding rule behavior, especially while diff/preflight foundation is the next priority.

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

The project organizes work using four main label families. Use the current labels
below when suggesting or helping with issues; maintainers may add labels as the
project evolves.

### Type

Describes the kind of work:

- `type: chore`
- `type: design`
- `type: docs`
- `type: epic`
- `type: feature`
- `type: feedback`
- `type: release`
- `type: rule`
- `type: test`

### Area

Identifies the part of the project affected:

- `area: cli`
- `area: config`
- `area: diff`
- `area: docs`
- `area: examples`
- `area: lint`
- `area: migration-versioning`
- `area: release`
- `area: reports`

### Priority

Indicates urgency or planned order:

- `priority: p0`
- `priority: p1`
- `priority: p2`
- `priority: p3`

### Status

Shows the current execution state:

- `status: needs-design`
- `status: ready`
- `status: blocked`
- `status: future`

### Community

These GitHub community labels are separate from the project taxonomy:

- `good first issue`
- `help wanted`

The bug report template also uses GitHub's standard `bug` label. It is not a
`type:` label.

## Security and privacy

Do not include production mappings, customer data, real logs, credentials, tokens, internal service names, or confidential index patterns.

Read `SECURITY.md` before reporting security issues or sharing examples.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

Code of conduct incidents may be reported to search-index-preflight@proton.me.
