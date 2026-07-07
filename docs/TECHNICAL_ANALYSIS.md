# Technical Analysis

## Scope

SearchIndexPreflight analyzes Elasticsearch/OpenSearch schema artifacts before they are applied to a cluster.

Initial artifact types:

- mappings
- component templates
- index templates
- dynamic templates
- settings relevant to mapping behavior
- sample documents in JSONL format

The MVP is offline-only. It must not require network access, cluster credentials, or production data.

## Mappings

Mappings define field types and mapping parameters for documents.

SearchIndexPreflight should normalize mappings into a canonical tree:

```text
Mapping
  dynamic
  date_detection
  numeric_detection
  properties
    field name
      type
      parameters
      fields
      properties
      dynamic
      enabled
      runtime
```

Rules should receive a typed normalized model with source file, JSON pointer, parent context, dialect/version capability context, and mapping parameters.

Rules should not operate directly on raw `map[string]any`.

## Field counting

Field count analysis should include concrete field mappings, object mappings, aliases, runtime fields, and multi-fields.

Rules compare count with:

- `index.mapping.total_fields.limit`, when supplied
- default limit profile for selected dialect/version
- configurable thresholds, e.g. warn at 80%

Raising the limit is not a safe default remediation. Reducing uncontrolled field growth should be preferred.

## Dynamic mapping

Dynamic mapping is convenient but often produces production schema drift.

Checks:

- root `dynamic: true`
- dynamic objects under high-risk paths
- missing explicit mapping for common sample fields
- dynamic date/numeric detection risks
- unknown fields in sample docs when strict mapping is expected
- sample docs that would likely create many new fields

SearchIndexPreflight cannot know production key cardinality offline. Findings should be framed as risks unless deterministic evidence exists.

## Dynamic templates

Dynamic templates are ordered. First match wins.

SearchIndexPreflight should parse and analyze:

- template name
- `match`, `unmatch`
- `path_match`, `path_unmatch`
- `match_mapping_type`
- `mapping`
- variables such as `{{name}}` and `{{dynamic_type}}`

Checks:

- missing `match_mapping_type`
- overbroad `match: "*"`
- overbroad `path_match: "*"`
- templates that may match object paths
- shadowed templates
- broad string-to-text+keyword expansion
- unsupported field types for selected dialect
- analyzer/normalizer not defined in supplied settings

## Index templates

SearchIndexPreflight should parse template name, index patterns, priority, composed components, template settings, template mappings, aliases, data stream declaration, and `_meta`.

Offline analysis can detect overlapping patterns, same-priority conflicts, suspicious priority changes, missing referenced component templates, component override conflicts, possible built-in pattern collisions, and data stream templates missing expected timestamp mapping.

Offline composition is approximate. The cluster remains authoritative.

## Component templates

SearchIndexPreflight should parse mappings, settings, aliases, version, and `_meta`.

Checks:

- missing referenced components
- conflicting fields between composed components
- override order risks
- analyzer/normalizer definitions used by fields
- settings conflicts

## Sample documents

Input format:

- JSONL
- one JSON object per line
- comments not allowed
- malformed lines produce parse errors

SearchIndexPreflight should infer observed types, null counts, max string length, arrays, array element types, object keys, and dotted field representations.

Checks:

- mapping type mismatch
- object vs scalar mismatch
- mixed array element types
- dotted field collision
- long keyword values
- numeric identifier inconsistency
- null-only field warning
- dynamic mapping would infer risky type
- sample field not declared in strict mapping

## What can be checked offline

Deterministic or near-deterministic:

- JSON/YAML syntax
- known mapping field type for selected dialect/version
- field count from supplied mapping/template
- total field limit threshold when settings are supplied
- explicit field type conflict across supplied files
- dotted field collision
- analyzer/normalizer reference missing from supplied settings
- index pattern overlap in supplied templates
- missing referenced component template
- unsupported field type for selected dialect/version
- malformed JSONL
- clear sample scalar/object conflicts

Heuristic:

- root dynamic mapping risk
- broad dynamic templates
- path matching risk
- text vs keyword misuse
- keyword too long
- object vs nested misuse
- flattened/flat_object semantic mismatch
- numeric ID as numeric type
- `index: false` on likely queried field
- `doc_values: false` on likely aggregatable field

## What requires a cluster

Future read-only cluster mode can improve:

- authoritative template simulation with `_simulate_index`
- actual cross-index field conflicts via `_field_caps`
- mapping drift
- current templates not present in repo
- cluster version and installed plugin capabilities
- analyzer behavior through `_analyze`
- live field count
- aliases/data streams resolution
- real runtime field usage and performance context

MVP must not depend on cluster checks.

## Elasticsearch vs OpenSearch differences

SearchIndexPreflight must model dialect explicitly:

```yaml
dialect:
  engine: elasticsearch
  version: "8.x"
```

or:

```yaml
dialect:
  engine: opensearch
  version: "2.x"
```

Divergence areas:

- supported field types
- version-specific field types
- `flattened` vs `flat_object`
- vector field mapping semantics
- plugin-dependent field types
- template API behavior
- data stream support
- runtime field support
- default/managed templates
- migration compatibility issues

## Version support recommendation

MVP:

- Elasticsearch 8.x target behavior
- Elasticsearch 7.17 best-effort compatibility
- OpenSearch 2.x target behavior
- OpenSearch `flat_object` only where supported

Alpha/Beta:

- expand compatibility fixtures
- add stricter capability matrix
- validate newer profiles if needed

v1:

- publish compatibility table
- fail fast on unsupported dialect/version combinations
- allow unknown version only in advisory mode unless configured otherwise

## Technical limitations

SearchIndexPreflight cannot fully determine production cardinality, query workload semantics, heap pressure, actual indexing throughput impact, plugin-provided field types, exact final template resolution without all templates/cluster simulation, or whether object should truly be nested without query requirements.

## False positives and false negatives

Expected false-positive areas include field naming heuristics, object vs nested recommendations, keyword length warnings, flattened/flat_object semantic warnings, numeric identifiers, and intentional dynamic mapping.

Expected false-negative areas include production payloads not represented in samples, missing external component templates, plugin-specific field types, query-driven modeling errors, and hidden producer serialization changes.

Mitigations:

- confidence field
- deterministic vs heuristic classification
- conservative default fail threshold
- suppressions with required reason
- baseline mode
- compatibility fixtures
