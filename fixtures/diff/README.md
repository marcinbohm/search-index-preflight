# Diff fixtures

These fixtures exercise the minimal public `search-index-preflight diff` command.

Current diff coverage is intentionally small:

- `dif001-field-type-changed/` emits one `DIF001` finding.
- `no-changes/` emits no diagnostics or findings.

The diff command currently matches directory inputs by relative path and emits only `DIF001`.
