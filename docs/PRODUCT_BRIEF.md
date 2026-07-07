# Product Brief: SearchIndexPreflight

## Product vision

SearchIndexPreflight helps engineering teams catch Elasticsearch/OpenSearch schema and template rollout risks before they reach production.

SearchIndexPreflight is a preflight safety CLI for Elasticsearch/OpenSearch schema changes.

The current implementation is an offline-first CLI for linting:

- mappings
- component templates
- index templates
- dynamic templates
- sample documents

The next strategic focus is diff-first schema change analysis: comparing a base schema corpus with a proposed schema corpus and reporting whether the change is safe to merge or deploy.

The core idea is simple: search schema changes should be reviewable, testable, and enforceable in CI the same way application code and infrastructure code are.

SearchIndexPreflight should become a practical guardrail for search infrastructure reliability, not a demo project and not a generic JSON validator.

## Problem statement

Teams often ship Elasticsearch/OpenSearch schema changes through Git, Terraform, application deploys, CI/CD pipelines, or ad-hoc scripts.

Many dangerous changes are valid JSON and valid mapping/template syntax. They pass basic validation but still create operational risk:

- dynamic field growth
- mapping explosion
- template precedence mistakes
- field type conflicts
- sample payload incompatibility
- analyzer/normalizer mismatches
- bad `text` vs `keyword` choices
- overly broad dynamic templates
- object/nested modeling mistakes
- Elasticsearch/OpenSearch compatibility gaps

These problems frequently surface only after indexing begins. At that point, rollback can require new indices, reindexing, rollover, DLQ draining, producer fixes, query changes, and manual operator work.

## Target users

Primary users:

- search infrastructure engineers
- distributed systems engineers working on indexing platforms
- platform engineers owning shared Elasticsearch/OpenSearch clusters
- SREs responsible for search reliability
- backend engineers shipping event/log/search schemas
- data platform engineers maintaining schema repositories

Secondary users:

- teams migrating between Elasticsearch and OpenSearch
- teams adopting schema-as-code workflows
- teams introducing CI checks for search infrastructure
- consultants or maintainers reviewing search schema changes

## Primary use cases

### PR-time schema preflight

A developer changes an index template or mapping. SearchIndexPreflight runs in CI and reports static risks today. The minimal experimental diff mode can compare old and new schema states for field type changes, removed fields, and added fields before merge.

### Offline mapping/template risk review

A maintainer runs SearchIndexPreflight locally before applying templates to a cluster.

```bash
search-index-preflight lint ./schemas
search-index-preflight lint --template logs.index-template.json
```

### Schema change diff

A maintainer or CI workflow compares the current repository schema with a proposed schema directory.

```bash
search-index-preflight diff --base old/ --current new/
```

This is minimal experimental behavior. It currently emits `DIF001` field type changes, `DIF002` field removals, and `DIF003` field additions; broader diff analysis is still planned.

### Sample document compatibility checks

A team maintains sample payloads alongside mappings. SearchIndexPreflight verifies whether sample documents match declared schema expectations.

### Template precedence checks

A platform repository contains multiple index templates and component templates. SearchIndexPreflight detects likely precedence and composition issues before rollout.

### CI-friendly reports

SearchIndexPreflight currently emits console and JSON reports. Markdown and SARIF are planned future report formats.

## Secondary use cases

- public fixture corpus of Elasticsearch/OpenSearch schema pitfalls
- repository audit before stricter schema governance
- baseline for legacy schemas
- old/new schema directory comparison
- future read-only cluster validation / doctor mode
- education through rule explanations and remediation guidance

## Non-goals

SearchIndexPreflight is not a SaaS product, dashboard, OpenSearch Dashboards plugin, mapping generator, schema migration framework, automatic fixer, replacement for staging/load tests, or a tool that writes to clusters. Future doctor mode must be read-only.

## User stories

### Search infrastructure engineer

As a search infrastructure engineer, I want schema changes to fail CI when they introduce high-confidence mapping risks, so that production indexing incidents are less likely.

### Platform engineer

As a platform engineer, I want to lint all index templates in a repository, so that template priority and pattern collisions are caught before rollout.

### Backend engineer

As a backend engineer, I want to validate sample documents against mappings, so that producer serialization mistakes are found before deployment.

### SRE

As an SRE, I want suppressions and baselines, so that legacy issues do not block adoption while new risks still fail CI.

## Success metrics

| Phase | Metrics |
|---|---|
| Pre-alpha | CLI lints a schema directory; SIL001-SIL003; fixtures; deterministic JSON; no cluster required |
| Alpha | diff/preflight foundation; Markdown/PR report direction; public fixtures; first external feedback |
| Beta | baseline or SARIF if still needed; compatibility profiles; docs for default-on rules; CI on Linux/macOS/Windows |
| v1 | stable CLI/JSON/rule IDs; signed artifacts; action v1; 25+ documented rules |

## Adoption strategy

1. Documentation before code.
2. Fixture-first MVP.
3. GitHub-native alpha with SARIF.
4. Focused external feedback on false positives.
5. v1 contract freeze.

## Positioning

SearchIndexPreflight is offline-first, CI-first, schema-as-code oriented, rule-based, explainable, Elasticsearch/OpenSearch aware, conservative about heuristics, and designed for infrastructure teams. Change safety, especially future diff-based preflight checks, is the primary product focus.

> SearchIndexPreflight catches Elasticsearch/OpenSearch schema and template risks before they hit production.

Future positioning:

> SearchIndexPreflight answers whether an Elasticsearch/OpenSearch schema change is safe to merge or deploy.

## Risks

| Risk | Severity | Mitigation |
|---|---:|---|
| False positives | Critical | severity + confidence, suppressions, baseline, conservative defaults |
| Wrong remediation advice | High | no auto-fix, cautious wording, fixtures, review checklist |
| Scope creep | Critical | non-goals, staged roadmap, cluster mode after v1.1 |
| Elasticsearch/OpenSearch divergence | High | dialect model, compatibility matrix, fixtures |
| Weak adoption | Medium | GitHub Action, SARIF, clear README, fixture corpus |
| Privacy mistakes | High | SECURITY.md, synthetic fixtures, issue template warnings |

## Open questions

- Exact initial license: recommended Apache-2.0, maintainer to confirm.
- Exact initial Elasticsearch version matrix: TBD before alpha.
- Exact initial OpenSearch version matrix: TBD before alpha.
- Whether JSON Schema for `search-index-preflight.yaml` ships in alpha or beta: TBD.
