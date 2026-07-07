# Rule Catalog


## Rule policy

Every rule has stable ID, stable name, category, description, explanation, applicability, required input, determinism classification, default severity, false-positive risk, bad input example, remediation guidance, reference placeholder, and stage.

Rule IDs must never be reused.

Severity and confidence are separate.

Deterministic findings may fail CI by default if severity is `error` or higher. Heuristic findings should be conservative by default.

## Stages

- MVP: required for first usable CLI
- Alpha: useful before public alpha
- Beta: useful before broad external testing
- v1: stable release target
- Future: not planned before v1

## Categories

- `mapping-limits`
- `dynamic-mapping`
- `dynamic-templates`
- `field-conflicts`
- `samples`
- `text-keyword`
- `analyzers`
- `templates`
- `objects-nested`
- `compatibility`
- `runtime`
- `metadata`

## Summary

`SIL001` through `SIL003` are currently implemented static checks. Further static rule expansion is paused while the diff/preflight foundation is introduced. `DIF001`, `DIF002`, and `DIF003` are implemented for the minimal public `diff` command, but they are not emitted by `lint`. Existing rule IDs remain stable and were not renamed during the SearchIndexPreflight transition.

| ID | Name | Category | Stage | Severity | Determinism | FP risk |
|---|---|---|---|---|---|---|
| SIL001 | total-fields-limit-risk | mapping-limits | MVP implemented | error when exceeded; warning near threshold | deterministic | low |
| SIL002 | root-dynamic-enabled | dynamic-mapping | MVP implemented | warning | heuristic | medium |
| SIL003 | dynamic-template-missing-match-mapping-type | dynamic-templates | MVP implemented | warning | heuristic | medium |
| SIL004 | overbroad-dynamic-template | dynamic-templates | MVP | warning | heuristic | medium |
| SIL005 | dynamic-template-shadowing | dynamic-templates | MVP | warning | deterministic in simple cases; heuristic for wildcards | medium |
| SIL006 | path-match-object-collision-risk | dynamic-templates | MVP | error when confirmed; warning otherwise | deterministic pattern detection; heuristic risk | medium |
| SIL007 | dotted-field-collision | field-conflicts | MVP | error | deterministic | low |
| SIL008 | field-type-conflict | field-conflicts | MVP | error | deterministic | low |
| SIL009 | sample-doc-mapping-conflict | samples | MVP | error | deterministic for clear conflicts; heuristic for coercible values | medium |
| SIL010 | dynamic-date-numeric-detection-risk | dynamic-mapping | MVP | warning | heuristic | medium |
| SIL011 | likely-aggregatable-field-as-text | text-keyword | MVP | warning | heuristic | high |
| SIL012 | long-keyword-without-ignore-above | text-keyword | MVP | warning | heuristic | medium |
| SIL013 | fielddata-true-on-text | text-keyword | MVP | error | deterministic | low |
| SIL014 | missing-analyzer-normalizer-definition | analyzers | MVP | error | deterministic if settings supplied | medium |
| SIL015 | template-priority-conflict | templates | MVP | error/warning | deterministic for overlap; heuristic for specificity | medium |
| SIL016 | multi-field-expansion-risk | mapping-limits | Alpha | warning | heuristic | medium |
| SIL017 | array-of-objects-object-mapping-risk | objects-nested | Alpha | warning | heuristic | high |
| SIL018 | nested-limit-risk | objects-nested | Alpha | warning | deterministic/heuristic | low |
| SIL019 | keyword-likely-needs-normalizer | analyzers | Alpha | info | heuristic | high |
| SIL020 | component-template-missing | templates | Alpha | error | deterministic relative to supplied corpus | medium |
| SIL021 | component-template-override-conflict | templates | Beta | warning | deterministic/heuristic | medium |
| SIL022 | legacy-composable-template-collision | templates | Beta | warning | heuristic | medium |
| SIL023 | data-stream-missing-timestamp | templates | Beta | error | deterministic when data_stream declared | low |
| SIL024 | mixed-array-element-types | samples | Beta | warning | deterministic for samples | medium |
| SIL025 | null-only-sample-field | samples | Beta | info | heuristic | medium |
| SIL026 | mapping-depth-limit-risk | mapping-limits | Beta | warning/error | deterministic | low |
| SIL027 | numeric-identifier-risk | text-keyword | Beta | info | heuristic | high |
| SIL028 | runtime-fields-overuse-risk | runtime | v1 | warning | heuristic | high |
| SIL029 | unsupported-field-type-for-dialect | compatibility | v1 | error | deterministic when capability matrix known | low |
| SIL030 | source-disabled-risk | metadata | v1 | warning | heuristic | medium |

---

## Diff rules

Diff rules operate on old/new normalized corpora. The minimal public `diff` command currently emits DIF001, DIF002, and DIF003.

| ID | Name | Category | Status | Severity | Determinism | Notes |
|---|---|---|---|---|---|---|
| DIF001 | field-type-changed | schema-diff | experimental implemented | error | deterministic | Emitted by `diff`; not emitted by `lint`. |
| DIF002 | field-removed | schema-diff | experimental implemented | warning | deterministic | Emitted by `diff`; not emitted by `lint`; does not fail default `--fail-on error`. |
| DIF003 | field-added | schema-diff | experimental implemented | info | deterministic | Emitted by `diff`; not emitted by `lint`; does not fail default `--fail-on error` or `--fail-on warning`. |

---

## SIL001: total-fields-limit-risk

ID: SIL001  
Name: total-fields-limit-risk  
Category: mapping-limits  
Description: Detect mappings/templates that exceed or approach `index.mapping.total_fields.limit`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping or template; settings when available  
Deterministic vs heuristic: deterministic  
Default severity: error when exceeded; warning near threshold  
False positive risk: low  
Stage: MVP  
Implementation status: implemented in pre-alpha with default limit `1000` and warning threshold `800`. Counts explicit normalized properties, multi-fields, and runtime fields per standalone mapping/template mapping. Does not estimate dynamic expansion, compose component templates, read cluster state, or read config yet.

Bad input:

```text
fixtures/mapping-limits/sil001-total-fields-limit/mapping-over-limit.json
```

Remediation:

Reduce field count, restrict dynamic mapping, split unrelated data, use flattened/flat_object only when query semantics fit.

References: TBD


---

## SIL002: root-dynamic-enabled

ID: SIL002  
Name: root-dynamic-enabled  
Category: dynamic-mapping  
Description: Detect root-level `dynamic: true`.  
Why it matters: Root `dynamic: true` can allow unexpected fields to expand a mapping. This may be intentional, but it should be reviewed because field growth can create operational risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping or template  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: MVP  
Implementation status: implemented in pre-alpha. Checks only explicit root-level `dynamic: true` on standalone mappings and template mappings. Does not inspect child object dynamic settings, estimate dynamic expansion, or evaluate dynamic templates.

Bad input:

```text
fixtures/dynamic-mapping/sil002-root-dynamic-enabled/mapping-root-dynamic-true.json
```

Remediation:

Use explicit mappings for known fields. Consider `dynamic: strict` or `dynamic: false` for controlled schemas, or scope dynamic behavior to known safe objects. Keep dynamic enabled only when the expansion risk is intentional and reviewed.

References: TBD


---

## SIL003: dynamic-template-missing-match-mapping-type

ID: SIL003  
Name: dynamic-template-missing-match-mapping-type  
Category: dynamic-templates  
Description: Detect dynamic templates without `match_mapping_type`.  
Why it matters: A dynamic template without `match_mapping_type` can match more inferred field types than intended. This may be intentional, but it can create mapping growth, type compatibility, or query behavior surprises.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: dynamic templates  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: MVP  
Implementation status: implemented in pre-alpha. Checks normalized dynamic templates for missing `match_mapping_type`. Does not validate type compatibility, estimate dynamic expansion, analyze ordering/shadowing, or compose component templates.

Bad input:

```text
fixtures/dynamic-templates/sil003-missing-match-mapping-type/mapping-missing-match-mapping-type.json
```

Remediation:

Add `match_mapping_type` when the template is intended for a specific detected field type, or document why broad matching is intentional.

References: TBD


---

## SIL004: overbroad-dynamic-template

ID: SIL004  
Name: overbroad-dynamic-template  
Category: dynamic-templates  
Description: Detect broad dynamic templates such as `match: "*"`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: dynamic templates  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL004.
```

Remediation:

Restrict to known paths and avoid arbitrary text+keyword expansion.

References: TBD


---

## SIL005: dynamic-template-shadowing

ID: SIL005  
Name: dynamic-template-shadowing  
Category: dynamic-templates  
Description: Detect templates that may be shadowed by earlier templates.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: dynamic templates  
Deterministic vs heuristic: deterministic in simple cases; heuristic for wildcards  
Default severity: warning  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL005.
```

Remediation:

Place specific templates before generic templates and document order.

References: TBD


---

## SIL006: path-match-object-collision-risk

ID: SIL006  
Name: path-match-object-collision-risk  
Category: dynamic-templates  
Description: Detect `path_match` patterns that may match object paths.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: dynamic templates; samples improve confidence  
Deterministic vs heuristic: deterministic pattern detection; heuristic risk  
Default severity: error when confirmed; warning otherwise  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL006.
```

Remediation:

Use `match`, narrow paths, or explicitly map objects.

References: TBD


---

## SIL007: dotted-field-collision

ID: SIL007  
Name: dotted-field-collision  
Category: field-conflicts  
Description: Detect dotted/object collisions like `user.name` and `user: {...}`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mappings and/or samples  
Deterministic vs heuristic: deterministic  
Default severity: error  
False positive risk: low  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL007.
```

Remediation:

Normalize producers and choose one representation.

References: TBD


---

## SIL008: field-type-conflict

ID: SIL008  
Name: field-type-conflict  
Category: field-conflicts  
Description: Detect same field path declared with incompatible types across files.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: two or more mappings/templates  
Deterministic vs heuristic: deterministic  
Default severity: error  
False positive risk: low  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL008.
```

Remediation:

Choose one type; version index patterns or reindex when needed.

References: TBD


---

## SIL009: sample-doc-mapping-conflict

ID: SIL009  
Name: sample-doc-mapping-conflict  
Category: samples  
Description: Detect sample values incompatible with mapping types.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template + sample docs  
Deterministic vs heuristic: deterministic for clear conflicts; heuristic for coercible values  
Default severity: error  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL009.
```

Remediation:

Fix producer serialization or change mapping intentionally.

References: TBD


---

## SIL010: dynamic-date-numeric-detection-risk

ID: SIL010  
Name: dynamic-date-numeric-detection-risk  
Category: dynamic-mapping  
Description: Detect dynamic areas where strings may be inferred as dates/numbers.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/settings; samples improve confidence  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL010.
```

Remediation:

Explicitly map identifiers as keyword; disable risky detection where appropriate.

References: TBD


---

## SIL011: likely-aggregatable-field-as-text

ID: SIL011  
Name: likely-aggregatable-field-as-text  
Category: text-keyword  
Description: Detect likely filter/sort/aggregation fields mapped only as `text`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: high  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL011.
```

Remediation:

Use `keyword` or add a keyword subfield.

References: TBD


---

## SIL012: long-keyword-without-ignore-above

ID: SIL012  
Name: long-keyword-without-ignore-above  
Category: text-keyword  
Description: Detect `keyword` fields with long sample values and no `ignore_above`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template; samples improve confidence  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL012.
```

Remediation:

Add `ignore_above`, map long content as `text`, or split payloads.

References: TBD


---

## SIL013: fielddata-true-on-text

ID: SIL013  
Name: fielddata-true-on-text  
Category: text-keyword  
Description: Detect `fielddata: true` on `text` fields.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template  
Deterministic vs heuristic: deterministic  
Default severity: error  
False positive risk: low  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL013.
```

Remediation:

Add a keyword subfield and aggregate/sort on keyword.

References: TBD


---

## SIL014: missing-analyzer-normalizer-definition

ID: SIL014  
Name: missing-analyzer-normalizer-definition  
Category: analyzers  
Description: Detect analyzer/normalizer references missing from supplied settings.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template + settings  
Deterministic vs heuristic: deterministic if settings supplied  
Default severity: error  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL014.
```

Remediation:

Define the analyzer/normalizer or mark external settings in config.

References: TBD


---

## SIL015: template-priority-conflict

ID: SIL015  
Name: template-priority-conflict  
Category: templates  
Description: Detect overlapping index template patterns with same/conflicting priorities.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: two or more index templates  
Deterministic vs heuristic: deterministic for overlap; heuristic for specificity  
Default severity: error/warning  
False positive risk: medium  
Stage: MVP  

Bad input:

```text
TBD minimal fixture example for SIL015.
```

Remediation:

Raise priority for the specific template or narrow patterns.

References: TBD


---

## SIL016: multi-field-expansion-risk

ID: SIL016  
Name: multi-field-expansion-risk  
Category: mapping-limits  
Description: Detect dynamic templates that create multi-fields broadly.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template/dynamic templates  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: Alpha  

Bad input:

```text
TBD minimal fixture example for SIL016.
```

Remediation:

Restrict templates or use explicit mappings.

References: TBD


---

## SIL017: array-of-objects-object-mapping-risk

ID: SIL017  
Name: array-of-objects-object-mapping-risk  
Category: objects-nested  
Description: Detect arrays of objects mapped as plain `object`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping + samples  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: high  
Stage: Alpha  

Bad input:

```text
TBD minimal fixture example for SIL017.
```

Remediation:

Use `nested` only if same-object correlation is required.

References: TBD


---

## SIL018: nested-limit-risk

ID: SIL018  
Name: nested-limit-risk  
Category: objects-nested  
Description: Detect nested fields/objects approaching configured limits.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/settings; samples improve confidence  
Deterministic vs heuristic: deterministic/heuristic  
Default severity: warning  
False positive risk: low  
Stage: Alpha  

Bad input:

```text
TBD minimal fixture example for SIL018.
```

Remediation:

Reduce nested usage, split docs, or denormalize safely.

References: TBD


---

## SIL019: keyword-likely-needs-normalizer

ID: SIL019  
Name: keyword-likely-needs-normalizer  
Category: analyzers  
Description: Detect keyword fields likely needing case-insensitive matching.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template  
Deterministic vs heuristic: heuristic  
Default severity: info  
False positive risk: high  
Stage: Alpha  

Bad input:

```text
TBD minimal fixture example for SIL019.
```

Remediation:

Add lowercase normalizer if needed; otherwise suppress with reason.

References: TBD


---

## SIL020: component-template-missing

ID: SIL020  
Name: component-template-missing  
Category: templates  
Description: Detect missing `composed_of` component templates.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: index + component templates  
Deterministic vs heuristic: deterministic relative to supplied corpus  
Default severity: error  
False positive risk: medium  
Stage: Alpha  

Bad input:

```text
TBD minimal fixture example for SIL020.
```

Remediation:

Add component file or declare known external.

References: TBD


---

## SIL021: component-template-override-conflict

ID: SIL021  
Name: component-template-override-conflict  
Category: templates  
Description: Detect composed components defining same field/setting differently.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: component + index templates  
Deterministic vs heuristic: deterministic/heuristic  
Default severity: warning  
False positive risk: medium  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL021.
```

Remediation:

Avoid conflicts or make override order explicit.

References: TBD


---

## SIL022: legacy-composable-template-collision

ID: SIL022  
Name: legacy-composable-template-collision  
Category: templates  
Description: Detect legacy and composable templates targeting same patterns.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: legacy + composable templates  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL022.
```

Remediation:

Migrate legacy templates or document coexistence.

References: TBD


---

## SIL023: data-stream-missing-timestamp

ID: SIL023  
Name: data-stream-missing-timestamp  
Category: templates  
Description: Detect data stream templates without explicit timestamp mapping.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: index template  
Deterministic vs heuristic: deterministic when data_stream declared  
Default severity: error  
False positive risk: low  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL023.
```

Remediation:

Add explicit timestamp date mapping.

References: TBD


---

## SIL024: mixed-array-element-types

ID: SIL024  
Name: mixed-array-element-types  
Category: samples  
Description: Detect arrays with inconsistent element types in samples.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: sample docs  
Deterministic vs heuristic: deterministic for samples  
Default severity: warning  
False positive risk: medium  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL024.
```

Remediation:

Normalize producer output or split fields.

References: TBD


---

## SIL025: null-only-sample-field

ID: SIL025  
Name: null-only-sample-field  
Category: samples  
Description: Detect fields that appear only as null in sample docs.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: sample docs  
Deterministic vs heuristic: heuristic  
Default severity: info  
False positive risk: medium  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL025.
```

Remediation:

Add representative non-null samples or explicit mapping.

References: TBD


---

## SIL026: mapping-depth-limit-risk

ID: SIL026  
Name: mapping-depth-limit-risk  
Category: mapping-limits  
Description: Detect deeply nested object paths approaching/exceeding depth limits.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/settings  
Deterministic vs heuristic: deterministic  
Default severity: warning/error  
False positive risk: low  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL026.
```

Remediation:

Flatten document structure or split entities.

References: TBD


---

## SIL027: numeric-identifier-risk

ID: SIL027  
Name: numeric-identifier-risk  
Category: text-keyword  
Description: Detect identifier-like fields mapped as numeric types.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template; samples improve confidence  
Deterministic vs heuristic: heuristic  
Default severity: info  
False positive risk: high  
Stage: Beta  

Bad input:

```text
TBD minimal fixture example for SIL027.
```

Remediation:

Use keyword for opaque IDs; keep numeric if numeric semantics are intended.

References: TBD


---

## SIL028: runtime-fields-overuse-risk

ID: SIL028  
Name: runtime-fields-overuse-risk  
Category: runtime  
Description: Detect many runtime fields or runtime fields under hot paths.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template; query hints improve confidence  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: high  
Stage: v1  

Bad input:

```text
TBD minimal fixture example for SIL028.
```

Remediation:

Index hot fields; keep runtime fields for rare/exploratory usage.

References: TBD


---

## SIL029: unsupported-field-type-for-dialect

ID: SIL029  
Name: unsupported-field-type-for-dialect  
Category: compatibility  
Description: Detect field types unsupported by configured engine/version.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template + dialect config  
Deterministic vs heuristic: deterministic when capability matrix known  
Default severity: error  
False positive risk: low  
Stage: v1  

Bad input:

```text
TBD minimal fixture example for SIL029.
```

Remediation:

Select correct dialect/version or use supported field type.

References: TBD


---

## SIL030: source-disabled-risk

ID: SIL030  
Name: source-disabled-risk  
Category: metadata  
Description: Detect mappings that disable `_source`.  
Why it matters: This pattern can create rollout, indexing, search correctness, compatibility, or operational reliability risk.  
Applies to: Elasticsearch/OpenSearch unless the selected dialect says otherwise.  
Input required: mapping/template  
Deterministic vs heuristic: heuristic  
Default severity: warning  
False positive risk: medium  
Stage: v1  

Bad input:

```text
TBD minimal fixture example for SIL030.
```

Remediation:

Keep `_source` enabled unless recovery/reindex strategy is documented.

References: TBD
