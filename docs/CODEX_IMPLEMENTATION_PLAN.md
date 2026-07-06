# Codex Implementation Plan

## Role

Codex CLI is the implementation agent.

Codex should implement SearchIndexLint in small, reviewable changes that follow the documented architecture and CLI contract.

Codex must not invent product scope. If a task is not in the roadmap, Codex should stop and ask for maintainer review.

## Current State Before Next Sprint

Implemented foundations:

- Go CLI skeleton
- parse-only `lint` command with console and JSON reports
- input discovery for `.json`, `.jsonl`, and `.ndjson`
- JSON and JSONL/NDJSON parsing
- canonical mapping/template model foundations
- canonical `model.Corpus`
- normalized field traversal helpers
- rule registry and runner foundation
- built-in rule: SIL001 total fields limit risk
- built-in rule: SIL002 root dynamic enabled

Not implemented:

- SIL003 and the rest of the rule catalog
- YAML
- config loading
- suppressions
- SARIF
- Markdown reporter
- baseline
- diff
- cluster mode

Next expected sprint should either harden implemented rule fixtures/report coverage or implement the next small rule. Do not reimplement parser, normalizer, corpus, traversal, or rule-runner foundations.

## Working rules

1. Do not start by adding many rules.
2. Build the skeleton first.
3. Keep changes small.
4. Every rule needs fixtures.
5. Every behavior needs tests.
6. Do not add live cluster mode.
7. Do not add auto-fix.
8. Do not add UI.
9. Do not add SaaS-related code.
10. Do not log sample document values unless needed and truncated.
11. Do not use company/private data.
12. Do not add dependencies without justification.
13. Do not change CLI contract silently.
14. Do not reuse rule IDs.
15. Do not make heuristic warnings fail CI by default.

## Implementation order

### Phase 0: repository setup

Create Go module, CLI root, `version`, `lint`, `rules list`, and `explain` stubs, CI workflow, and placeholder version package.

Files:

```text
go.mod
cmd/search-index-lint/main.go
internal/cli/root.go
internal/cli/lint.go
internal/cli/rules.go
internal/cli/explain.go
internal/version/version.go
.github/workflows/ci.yml
```

Stop for review after Phase 0.

### Phase 1: report and finding model

Define severity, confidence, determinism, finding, diagnostic, summary, JSON reporter, console reporter, and deterministic output tests.

### Phase 2: input discovery and parsers

Implement file source model, explicit file loading, directory walk, default excludes, JSON parser, YAML parser, JSONL parser, diagnostics, and invalid input tests.

### Phase 3: canonical models

Define mapping, field, dynamic template, index template, component template, sample schema, JSON pointer tracking, and source file tracking.

### Phase 4: normalizers

Detect raw/wrapped mappings, index/component templates, normalize property trees, multi-fields, runtime fields, dynamic templates, and sample schema.

### Phase 5: rule engine

Define rule interface, metadata, registry, context, corpus, rule filtering, and registry tests.

### Phase 6: first MVP rules

Implement in this order:

1. SIL001 total fields limit risk
2. SIL007 dotted field collision
3. SIL008 field type conflict
4. SIL014 analyzer/normalizer missing
5. SIL002 root dynamic enabled
6. SIL003 dynamic template missing match_mapping_type
7. SIL004 overbroad dynamic template
8. SIL005 dynamic template shadowing
9. SIL006 path_match object collision risk
10. SIL009 sample doc mapping conflict

Stop for review after every 3 rules.

### Phase 7: fixtures and golden harness

Define fixture metadata, discovery, golden JSON comparison, path normalization, deterministic finding order, and explicit golden update mechanism.

### Phase 8: config and suppressions

Load `search-index-lint.yaml`, parse dialect, inputs, rules, severity threshold, suppressions, require reasons, warn on expiry, and apply suppressions.

### Phase 9: CLI UX polish

Wire lint to engine, implement format/output/fail-on/quiet/verbose, rules list, explain, and exit codes.

## Naming conventions

Rule files:

```text
sil001_total_fields_limit.go
sil002_root_dynamic_enabled.go
sil003_dynamic_template_missing_match_mapping_type.go
```

Fixture names:

```text
fixtures/<category>/<short-kebab-case>/
```

## Coding standards

Required:

- gofmt
- go test
- table-driven tests where practical
- small packages
- explicit errors
- no panics for user input
- no global mutable state
- deterministic output ordering
- no hidden network calls
- no private data in tests

## What Codex must not do without approval

- add live cluster mode
- add write operations
- add auto-fix
- add SaaS code
- add dashboard/plugin code
- change language/stack
- add plugin API
- change rule IDs
- change default fail threshold
- make heuristic findings fail by default
- introduce telemetry
- upload samples anywhere
- add dependencies with unclear licenses
