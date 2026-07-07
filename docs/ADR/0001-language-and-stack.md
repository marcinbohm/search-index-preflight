# ADR 0001: Language and Stack

## Status

Accepted.

## Date

2026-07-06

## Context

SearchIndexPreflight is an offline-first CLI and future GitHub Action for linting Elasticsearch/OpenSearch mappings, templates, dynamic templates, and sample documents.

The tool needs to be easy to run in CI, easy to distribute as a single binary, fast enough for schema repositories, maintainable by OSS contributors, suitable for infrastructure tooling, capable of producing JSON/Markdown/SARIF, cross-platform, practical for GitHub Actions, and safe by default.

## Decision

SearchIndexPreflight will be implemented in Go.

Initial stack:

- Go
- Cobra or a similarly mature CLI framework
- standard `encoding/json`
- YAML parser library, exact dependency TBD
- internal canonical model for mappings/templates
- golden file tests
- GoReleaser or equivalent release tooling later
- GitHub Actions CI

Recommended license: Apache-2.0, pending maintainer confirmation.

## Options considered

| Language | Pros | Cons | Decision |
|---|---|---|---|
| Go | strong CLI fit, static binaries, fast startup, simple releases, common infra-tooling language | less expressive than Rust, model code can be verbose | selected |
| Rust | performance, type safety, good binaries | higher contributor friction, slower early iteration | rejected for MVP |
| Java | ES ecosystem familiarity, mature libraries | heavier runtime, weaker single-binary story | rejected for MVP |
| Kotlin | concise JVM, good modeling | JVM distribution burden, build complexity | rejected for MVP |
| Python | fastest prototyping, easy JSON/YAML | packaging friction, slower, weaker binary story | rejected for MVP |
| TypeScript | GitHub Action ecosystem, JSON ergonomics | runtime dependency, weaker infra CLI signal | rejected for core CLI |

## Rationale

Go is selected because SearchIndexPreflight should behave like a serious infrastructure CLI: single binary, quick startup, deterministic behavior, simple CI integration, low installation burden, straightforward release process, readable codebase, and practical contributor experience.

The hard parts are domain modeling, rule quality, fixture coverage, false-positive control, report UX, compatibility handling, and documentation. Go is sufficient for those and avoids unnecessary complexity.

## Consequences

Positive:

- easy release automation
- easy local usage
- easy GitHub Action wrapper
- predictable CI behavior
- strong fit for open-source infra tooling
- maintainable codebase

Negative:

- model code may be verbose
- error handling requires discipline
- contributors must avoid untyped map-based rule logic
- no Rust-level compile-time guarantees
- plugin/custom rule API may be harder later

## Guardrails

Implementation must define a typed internal model, avoid rule logic over raw `map[string]any`, keep CLI and rule engine separate, use table-driven tests, use fixture/golden tests, return structured errors, avoid global mutable state, and avoid network access in offline commands.

## Revisit criteria

Revisit only if Go prevents required functionality, performance targets cannot be met, contributor feedback strongly favors another ecosystem, project direction changes toward plugin APIs or embedded libraries, or maintainership changes materially.

No revisit is planned before v1.
