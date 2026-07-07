# DIF001 field type changed

This fixture compares a base mapping where `status` is `keyword` with a current mapping where `status` is `long`.

Run:

```bash
search-index-preflight diff --base base --current current
```

Expected:

- exit code `1`
- one `DIF001` finding
- finding points to `mapping.json#/properties/status`
