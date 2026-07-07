# ADR 0003: Rename to SearchIndexPreflight

## Status

Accepted.

## Date

2026-07-07

## Context

The project began as SearchIndexLint / `search-index-lint`: an offline-first linter for Elasticsearch/OpenSearch mappings, templates, dynamic templates, and sample documents.

ADR 0002 changed the strategic direction toward preflight safety for Elasticsearch/OpenSearch schema changes. Static linting remains useful, but it is now the offline-fast subset of a broader product direction centered on change safety.

The old name overemphasized linting as the whole product. The new direction needs a name that can cover static checks, future diff-based preflight analysis, and future read-only doctor mode.

## Decision

Rename the project, Go module, command path, binary references, docs, tests, and report tool name to SearchIndexPreflight / `search-index-preflight`.

The Go module path becomes:

```text
github.com/marcinbohm/search-index-preflight
```

The current implemented command remains:

```bash
search-index-preflight lint
```

Rule IDs remain unchanged:

```text
SIL001
SIL002
SIL003
```

No functional behavior changes are introduced by this rename.

## Consequences

- Go imports and module path are updated.
- The command directory is now `cmd/search-index-preflight`.
- User-facing CLI strings use SearchIndexPreflight / `search-index-preflight`.
- Public docs and tests use the new name.
- JSON report `tool.name` uses `SearchIndexPreflight`.
- The GitHub repository should be renamed or redirected outside this code change if it has not already been renamed.
- `check`, `diff`, and `doctor` remain future work.
- No new rules, diff behavior, doctor behavior, cluster access, config, suppressions, SARIF, Markdown reporter, GitHub Action wrapper, or dependencies are added.
