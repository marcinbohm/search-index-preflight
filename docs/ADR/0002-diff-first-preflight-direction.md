# ADR 0002: Diff-First Preflight Direction

## Status

Accepted.

## Date

2026-07-07

## Context

At the time of this ADR, the project was still named SearchIndexLint and had a working foundation: input discovery, JSON and JSONL/NDJSON parsing, mapping/template normalization, `model.Corpus`, a rule registry and runner, console/JSON reports, and the implemented static rules `SIL001`, `SIL002`, and `SIL003`.

Those static checks are useful, but they are not enough as the sole product direction. Pull request workflows usually need to understand a proposed schema change, not only whether the current checked-out schema has risky patterns. A mapping or template may be acceptable in isolation but risky as a change because it narrows a type, removes a field, alters dynamic behavior, changes template coverage, or diverges from existing indexed data expectations.

The project is therefore evolving from a state-first mapping linter into a preflight safety CLI for Elasticsearch/OpenSearch schema changes.

## Decision

The project direction is preflight checks for Elasticsearch/OpenSearch schema changes.

Current static `lint` behavior remains valid and should continue to work. It becomes the offline-fast static check subset of the broader preflight product.

The next major product capability should be `diff`: comparing a base schema corpus with a proposed schema corpus, deriving semantic changes, and running diff-aware rules that are useful in pull requests and deployment gates.

Read-only `doctor` mode is planned later. It may inspect cluster state, mappings, templates, field capabilities, versions, or template simulations, but it must remain read-only.

Future oracle or engine-backed validation is allowed as a possible direction, but it is not MVP code and should not be introduced before the diff foundation is established.

## Consequences

- The next milestone should focus on diff/preflight foundation.
- Avoid adding many more heuristic static rules before diff exists.
- Keep static rules as the offline-fast subset; do not delete or rewrite them.
- Reports must remain useful for pull request comments and CI.
- Rule IDs must remain stable.
- Existing `SIL001` through `SIL003` behavior and fixtures remain part of the project.
- No cluster write operations are allowed.
- No SaaS, UI, dashboard, or telemetry direction is introduced.
- The planned SearchIndexPreflight rename should be handled as a separate dedicated change, recorded later in ADR 0003.

## Non-Goals

This ADR does not implement diff mode, doctor mode, cluster access, oracle validation, config, suppressions, SARIF, Markdown reporting, a GitHub Action wrapper, or any repository/module rename. The rename itself is recorded separately in ADR 0003.
