# Codex Implementation Plan

## Role

Codex CLI is the implementation agent.

Codex should implement SearchIndexPreflight in small, reviewable changes that follow the documented architecture and CLI contract.

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
- built-in rule: SIL003 dynamic template missing match_mapping_type
- internal diff foundation for field added/removed/type-changed changes
- internal diff-aware rule layer with DIF001 field-type-changed, DIF002 field-removed, and DIF003 field-added
- minimal public `diff --base --current` command
- public `rules list` metadata output for lint and diff rules
- public `explain <RULE_ID>` metadata output for one lint or diff rule

Not implemented:

- SIL004 and the rest of the rule catalog
- YAML
- config loading
- suppressions
- SARIF
- Markdown reporter
- baseline
- git-aware diff options
- diff rules beyond DIF001/DIF002/DIF003
- migration/versioning commands or validation
- cluster mode

Next expected sprint should not implement SIL004 by default. Do not reimplement parser, normalizer, corpus, traversal, static rule-runner, internal field-diff, internal DIF001/DIF002/DIF003, minimal public diff command foundations, `rules list` metadata output, or `explain` metadata output.

## Current Strategic Direction

SearchIndexPreflight is a preflight safety CLI for Elasticsearch/OpenSearch schema changes.

Do not implement SIL004 next unless explicitly approved. Do not add more heuristic static rules until the diff/preflight foundation starts.

The current diff foundation compares two normalized corpora and detects field added, field removed, and field type changed events. The public `diff` command currently emits DIF001 for field type changes, DIF002 for field removals, and DIF003 for field additions. The next code phase should choose one of these review-approved paths:

- harden minimal public diff behavior and fixtures
- harden the internal diff-rule layer further
- document or refine the future offline migration/versioning concept
- no oracle/engine-backed validation yet
- no cluster mode
- no cluster writes

Keep existing `lint` behavior working during and after the rename. Static checks SIL001-SIL003 remain the offline-fast subset of the future preflight product.

## Future Offline Migration/Versioning Direction

ADR 0004 and `docs/MIGRATION_VERSIONING_CONCEPT.md` describe planned future work only.

Do not implement migration/versioning commands unless explicitly requested. In particular, do not add `versions validate`, `migrations validate`, `migrations plan`, `migrate apply`, cluster writes, alias cutover execution, reindex execution, rollback execution, or migration state in Elasticsearch/OpenSearch.

If implementation is explicitly requested later, build on the current lint/diff foundations:

- validate ordered schema inputs or migration manifests offline
- run existing lint checks against each schema state
- run existing diff checks between consecutive states
- report preflight risks without mutating clusters

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

Create Go module, CLI root, `version`, `lint`, initial `rules list` command, initial `explain` command, CI workflow, and placeholder version package.

Files:

```text
go.mod
cmd/search-index-preflight/main.go
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

Load `search-index-preflight.yaml`, parse dialect, inputs, rules, severity threshold, suppressions, require reasons, warn on expiry, and apply suppressions.

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
