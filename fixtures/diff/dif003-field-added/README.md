# DIF003 field added

This fixture compares a base mapping without `customer_id` with a current
mapping where that field has been added.

Expected:

- `search-index-preflight diff --base base --current current`
- exit code `0` with the default `--fail-on error`
- one info `DIF003` finding
- exit code `1` when run with `--fail-on info`

The data is synthetic and public-safe.
