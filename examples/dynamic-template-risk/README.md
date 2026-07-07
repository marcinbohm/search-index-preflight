# Dynamic template risk

This example shows a dynamic template that maps fields to `keyword` but does not declare `match_mapping_type`.

## Run

```bash
go run ./cmd/search-index-preflight lint \
  --mapping examples/dynamic-template-risk/mapping.json
```

## Expected finding

The excerpt below shows the finding line. See `expected-output.txt` for the full expected stdout.

```text
warning SIL003: examples/dynamic-template-risk/mapping.json#/dynamic_templates/0/strings_as_keywords: Dynamic template "strings_as_keywords" does not declare match_mapping_type.
```

This warning-only example exits with status 0 under the default `--fail-on error` threshold.

## Why this matters

A dynamic template without `match_mapping_type` can apply more broadly than intended. That may be deliberate, but it should be reviewed because it can affect field growth, type compatibility, and query behavior.
