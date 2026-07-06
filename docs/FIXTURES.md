# Fixtures

## Purpose

Fixtures are a core product asset, not only test data.

They provide repeatable rule tests, documentation examples, regression coverage for false positives, public learning material, reviewable behavior for contributors, and stable golden outputs.

Every rule must be backed by fixtures before it is considered done.

## Fixture principles

1. No company data.
2. No customer data.
3. No production logs.
4. No internal mappings.
5. No copied proprietary templates.
6. Small examples only.
7. Each fixture should show one primary problem.
8. Expected output must be deterministic.
9. Bad and good examples should both exist.
10. Remediation should be documented.
11. Fixture names should be boring and descriptive.
12. Every fixture should be safe to publish.

## Directory structure

```text
fixtures/
  mapping-explosion/
  field-conflict/
  dynamic-mapping/
  dynamic-template-risk/
  dotted-field-collision/
  text-keyword-misuse/
  keyword-too-long/
  analyzer-mismatch/
  template-priority-conflict/
  component-template-composition/
  nested-object-misuse/
  flattened-flat-object/
  sample-doc-conflict/
  compatibility/
  suppressions/
  baseline/
```

## Fixture case structure

```text
fixtures/<category>/<case>/
  README.md
  fixture.yaml
  search-index-lint.yaml          # optional
  mapping.json              # optional
  index-template.json       # optional
  component-template.json   # optional
  samples.jsonl             # optional
  expected.json
  expected.md               # optional
  expected.sarif.json       # optional
```

## Fixture metadata

```yaml
id: SIL015-template-priority-conflict
title: Same-priority overlapping index templates
rules:
  - SIL015
stage: MVP
dialects:
  - engine: elasticsearch
    version: "8.x"
  - engine: opensearch
    version: "2.x"
expected:
  critical: 0
  error: 1
  warning: 0
  info: 0
privacy:
  contains_production_data: false
  contains_company_data: false
```

## Public data rules

Allowed:

- synthetic service names
- synthetic IDs
- generic log messages
- tiny mappings written from scratch
- common field names such as `status`, `user_id`, `message`, `metadata`

Not allowed:

- internal service names
- internal team names
- internal index patterns
- customer names
- real URLs from logs
- real email addresses
- real tokens
- real trace IDs from production
- proprietary mappings
- copied vendor examples beyond minimal fair-use snippets

## Required fixture packs

| Pack | Primary rules |
|---|---|
| mapping-explosion | SIL001, SIL002, SIL016 |
| dynamic-template-risk | SIL003, SIL004, SIL005, SIL006 |
| dotted-field-collision | SIL007 |
| field-conflict | SIL008 |
| sample-doc-conflict | SIL009, SIL024, SIL025 |
| dynamic-mapping | SIL010 |
| text-keyword-misuse | SIL011, SIL012, SIL013, SIL027 |
| analyzer-mismatch | SIL014, SIL019 |
| template-priority-conflict | SIL015 |
| component-template-composition | SIL020, SIL021 |
| objects-nested | SIL017, SIL018 |
| compatibility | SIL029 |
| metadata | SIL030 |

## Current fixture packs

The first public fixture pack is:

```text
fixtures/mapping-limits/sil001-total-fields-limit/
```

It covers the implemented `SIL001` default-threshold behavior with synthetic near-limit and over-limit mappings plus expected JSON reports.

## Expected output

Every fixture must include `expected.json`.

Minimal expected report:

```json
{
  "schema_version": "0.1",
  "summary": {
    "files_scanned": 2,
    "findings_total": 1,
    "critical": 0,
    "error": 1,
    "warning": 0,
    "info": 0
  },
  "findings": [
    {
      "id": "SIL015",
      "severity": "error",
      "category": "templates",
      "file": "logs-app.index-template.json",
      "json_pointer": "/priority"
    }
  ]
}
```

Golden output rules:

- use relative paths
- sort findings deterministically
- exclude timestamps
- exclude absolute paths
- normalize OS-specific path separators
- stabilize fingerprints
- avoid nondeterministic map ordering

## Golden file testing

MVP:

- JSON reporter
- console summary
- rule findings

Alpha:

- Markdown reporter
- SARIF reporter

Beta:

- baseline filtering
- diff output
- compatibility matrix output

## Fixture review checklist

- Does the fixture isolate one main rule?
- Is all data synthetic?
- Is the expected finding deterministic?
- Does README explain why it matters?
- Is remediation practical?
- Is the fixture small enough to understand quickly?
- Does it avoid company/vendor-private assumptions?
- Does it test both bad and good behavior?
