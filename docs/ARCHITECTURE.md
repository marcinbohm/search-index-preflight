# Architecture

## Design principles

SearchIndexLint should be offline-first, deterministic where possible, explicit about heuristics, safe for CI, easy to run locally, useful without production cluster access, fixture-driven, dialect-aware, explainable through stable rule IDs, and conservative about failing builds.

## High-level architecture

```text
search-index-lint CLI
  |
  +-- input discovery
  |     +-- explicit files
  |     +-- directories
  |     +-- globs from config
  |
  +-- parsers
  |     +-- JSON
  |     +-- YAML
  |     +-- JSONL / NDJSON sample docs
  |
  +-- normalizers
  |     +-- mapping normalization
  |     +-- index template normalization
  |     +-- component template normalization
  |     +-- dynamic template normalization
  |     +-- sample schema inference
  |
  +-- dialect capability layer
  |     +-- Elasticsearch profiles
  |     +-- OpenSearch profiles
  |
  +-- composition layer
  |     +-- component template merge approximation
  |     +-- index template pattern overlap
  |     +-- priority analysis
  |
  +-- rule engine
  |     +-- rule registry
  |     +-- deterministic rules
  |     +-- heuristic rules
  |     +-- suppressions
  |     +-- baseline filtering
  |
  +-- reporters
        +-- console
        +-- JSON
        +-- Markdown
        +-- SARIF
```

## Current Implemented Architecture

Current pre-alpha code implements this foundation path:

```text
input discovery -> parser -> normalizer -> model.Corpus -> rule runner foundation -> reports
```

Implemented foundations:

- input discovery for explicit files and directories
- JSON parser for mappings/templates
- JSONL/NDJSON parser for sample documents
- mapping, index template, and component template normalization
- `model.Corpus` as the canonical shared corpus
- normalized field traversal helpers in `internal/model`
- rule registry and runner foundation
- console and JSON diagnostic reports

Current CLI behavior:

- `lint` reports parse and normalization diagnostics only
- rule runner exists but is not wired into `lint` for real findings
- no real SIL rules are implemented
- YAML, Markdown, SARIF, baseline, diff, config, suppressions, and cluster mode are planned future work

## Module architecture

```text
cmd/search-index-lint/
internal/
  cli/
  config/
  input/
  parser/
  model/
  dialect/
  compose/
  rules/
  report/
  suppress/
  baseline/
  version/
  cluster/      # future read-only mode only
fixtures/
tests/
docs/
examples/
action/
```

## CLI layer

Responsibilities:

- parse commands and flags
- load config
- validate user input
- collect input files
- select dialect/version
- call core engine
- render reports
- set exit codes

The CLI layer must not contain rule logic.

## Core engine

The core engine coordinates input discovery, parsing, normalization, composition approximation, sample schema inference, rule execution, suppression filtering, baseline filtering, and report model generation.

The engine should return structured results, not formatted text.

## Parsers

Parsers must eventually handle JSON/YAML mappings, JSON/YAML index templates, JSON/YAML component templates, and JSONL sample documents. The current implementation supports JSON mappings/templates and JSONL/NDJSON sample documents only.

Parser output should preserve source file path, raw document kind, JSON pointer, line/column where practical, and parse diagnostics.

Malformed inputs should produce clean errors, not panics.

## Normalizers

Normalizers convert raw documents into canonical models.

Supported document kinds:

- raw mapping
- index template
- component template
- settings + mappings bundle
- sample document collection
- config

## Canonical model

Core model entities:

```text
Mapping
Field
DynamicTemplate
IndexTemplate
ComponentTemplate
TemplateComposition
SampleDocumentSet
SampleField
AnalyzerReference
NormalizerReference
DialectProfile
Finding
RuleMetadata
```

Field model should include full path, parent path, field name, declared type, inferred type when applicable, mapping parameters, multi-fields, child properties, raw JSON pointer, source file, source location, and dialect support status.

## Rule engine

Rules should be small, testable, and metadata-driven.

```go
type Rule interface {
    Metadata() Metadata
    Check(ctx Context, corpus model.Corpus) ([]Finding, error)
}
```

Rules must not read files directly, write output, inspect CLI flags directly, panic on malformed input, depend on global mutable state, or make network calls.

## Severity model

Severity:

- `info`
- `warning`
- `error`
- `critical`

Confidence:

- `low`
- `medium`
- `high`

Determinism:

- `deterministic`
- `heuristic`
- `cluster-context-required`

Default CI behavior:

- fail on `error` and `critical`
- do not fail on `warning` or `info`
- do not fail on low-confidence heuristic findings unless configured

## Reporters

Console is default human output. JSON is stable machine output. Markdown is useful for PR comments. SARIF is alpha scope for GitHub code scanning.

Reports should be produced from structured findings, not from rule-specific text formatting.

## Config

Default config file names:

- `search-index-lint.yaml`
- `search-index-lint.yml`
- `.search-index-lint.yaml`

Config controls dialect/version, input globs, rules, rule config, severity threshold, suppressions, baseline, known external templates, and output defaults.

## Suppressions

Suppressions require rule ID and reason. File/path/owner/expiry are recommended. Expired suppressions should produce warnings.

## Baseline

Baseline is beta scope. It allows adoption in repositories with legacy issues and should fail only on new findings.

Fingerprint ingredients:

- rule ID
- normalized file path
- JSON pointer
- normalized finding key
- dialect

Do not include line numbers.

## Future GitHub Action

The GitHub Action should wrap the CLI and not add hidden behavior.

## Future cluster mode

Cluster mode is future work and must be read-only. It may fetch mappings/templates, run template simulation, run field capabilities, and fetch cluster version. It must never write to clusters.

## Fixture architecture

Every rule must be fixture-backed:

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

## Decisions

- Go as implementation language.
- Offline-first MVP.
- Rule IDs are stable and never reused.
- No auto-fix in MVP.
- Deterministic and heuristic rules are classified separately.
- Report model is internal first, rendered second.
- Fixtures are product assets.
