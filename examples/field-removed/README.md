# Field removed

This example shows a schema change where `legacy_id` exists in the base mapping but is removed from the current mapping.

## Run

```bash
go run ./cmd/search-index-preflight diff \
  --base examples/field-removed/base \
  --current examples/field-removed/current
```

## Expected finding

The excerpt below shows the finding line. See `expected-output.txt` for the full expected stdout.

```text
warning DIF002: mapping.json#/properties/legacy_id: Field "legacy_id" was removed from the current schema.
```

This warning-only diff exits with status 0 under the default `--fail-on error` threshold.

## Why this matters

Removed fields can break producers, queries, dashboards, alerts, and downstream consumers that still depend on the field. Intentional removals should be coordinated with rollover, reindexing, or consumer migration plans.
