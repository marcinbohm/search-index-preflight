# Quality Gates

## Quality philosophy

SearchIndexPreflight must be boring, deterministic, and trustworthy.

A noisy linter will be disabled. A misleading linter will damage trust. A hard-to-run linter will not be adopted.

## Test strategy

SearchIndexPreflight uses:

- unit tests
- parser tests
- normalizer tests
- rule tests
- fixture tests
- golden file tests
- CLI integration tests
- compatibility tests
- fuzz tests
- benchmark tests
- release tests

## Unit tests

Required for severity ordering, confidence handling, rule registry, AST traversal, field counting, pattern overlap, JSON pointer generation, config validation, suppression matching, baseline fingerprinting, and dialect capability lookup.

## Golden file tests

Required for JSON reports, Markdown reports, SARIF reports, fixture outputs, baseline outputs, and diff outputs.

Golden output must not include timestamps, absolute paths, random IDs, OS-specific separators, or map-order-dependent output.

Golden updates require review.

## Fixture tests

Every implemented rule must have positive fixture, negative fixture, expected JSON report, README, and synthetic data only.

Before v1, every default-on rule should also have suppression fixture, compatibility fixture where applicable, and reporter-specific golden fixture where needed.

## CLI integration tests

Cover:

- `search-index-preflight lint --mapping`
- `search-index-preflight lint --template`
- `search-index-preflight lint --component-template`
- `search-index-preflight lint --sample-docs`
- `search-index-preflight lint ./schemas`
- format selection
- output file writing
- fail thresholds
- config loading
- invalid usage
- invalid input parsing
- exit codes

## Compatibility tests

Cover Elasticsearch 8.x, Elasticsearch 7.17 best-effort where implemented, OpenSearch 2.x, unsupported field types, `flattened` vs `flat_object`, version-gated field types, and template behavior differences where modeled.

## Fuzz tests

Beta scope.

Targets:

- JSON parser wrapper
- YAML parser wrapper
- normalizer
- dynamic template parser
- sample document inference

Goals: no panics, no hangs, bounded memory, useful diagnostics.

## Benchmarks

Benchmark large mappings, many fields, many templates, dynamic templates, sample JSONL files, and report generation.

Initial target: lint a 10 MB schema directory under 2 seconds on a typical CI runner, excluding very large sample docs. This target is provisional.

## Linting and formatting

Required:

```bash
gofmt
go test ./...
go vet ./...
```

Recommended before alpha:

```bash
golangci-lint run
```

Recommended before beta:

```bash
gosec ./...
```

## Static analysis

Target tools:

- `go vet`
- `golangci-lint`
- `gosec`
- CodeQL
- OpenSSF Scorecard
- Dependabot or Renovate
- Trivy for Docker image when Docker image exists

## Coverage expectations

Pre-alpha: no strict global threshold, but all rule logic must have tests.  
Alpha/Beta/v1: 80%+ coverage for `internal/rules` and meaningful parser/normalizer coverage.

## Release checklist

### Pre-alpha

- [ ] README status says pre-alpha
- [ ] `go test ./...` passes
- [ ] first rules fixture-backed
- [ ] JSON report works
- [ ] no cluster mode
- [ ] no auto-fix

### Alpha

- [ ] 15+ rules
- [ ] Markdown reporter
- [ ] SARIF reporter
- [ ] GitHub Action preview
- [ ] rule docs
- [ ] compatibility profile skeleton
- [ ] CI on Linux/macOS/Windows
- [ ] release notes

### Beta

- [ ] baseline mode
- [ ] diff mode basic
- [ ] compatibility matrix
- [ ] fuzz tests
- [ ] benchmarks
- [ ] release binaries
- [ ] false-positive issue template
- [ ] changelog

### v1

- [ ] stable CLI contract
- [ ] stable JSON schema
- [ ] stable baseline schema
- [ ] stable rule ID policy
- [ ] all default-on rules documented
- [ ] all default-on rules fixture-backed
- [ ] SemVer policy
- [ ] signed artifacts
- [ ] checksums
- [ ] GitHub Action v1 tag
- [ ] SECURITY.md complete
- [ ] CONTRIBUTING.md complete
- [ ] changelog complete

## SemVer policy

After v1:

Patch releases: bug fixes, false-positive reductions, docs fixes, non-breaking report additions.  
Minor releases: new rules, new config options, new output fields, new compatibility profiles.  
Major releases: breaking CLI changes, breaking JSON schema changes, breaking baseline format changes.

## Rule versioning policy

Rule IDs are never reused. Rule severity changes must be documented. New heuristic rules should start as info/warning. Removed rules must remain documented as retired.

## PR blocking conditions

Block if tests fail, golden output changes without explanation, a new rule lacks fixtures, docs contradict behavior, private data appears, offline command makes network calls, heuristic rule fails CI by default without approval, CLI contract changes without docs, or SARIF output becomes invalid.
