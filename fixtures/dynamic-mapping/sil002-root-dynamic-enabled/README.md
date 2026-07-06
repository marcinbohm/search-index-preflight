# SIL002 Root Dynamic Enabled Fixtures

These fixtures exercise `SIL002` (`root-dynamic-enabled`), which detects mappings/templates where root-level dynamic mapping is explicitly enabled.

Why it matters:

- `dynamic: true` can allow unexpected fields to expand a mapping.
- Field growth can make index schemas harder to review and operate.
- The setting may still be intentional for flexible schemas, exploratory data, or controlled ingestion paths.

This rule is heuristic. It reports a warning and does not mean dynamic mapping is always wrong.

Current scope:

- checks only explicit root-level `dynamic: true`
- does not flag missing `dynamic`
- does not flag `dynamic: false`, `dynamic: strict`, or `dynamic: runtime`
- does not inspect child object dynamic settings
- does not estimate dynamic field expansion

Fixture cases:

- `mapping-root-dynamic-true.json`: emits one `SIL002` warning and exits `0` with the default `--fail-on error`.
- `mapping-root-dynamic-false.json`: emits no `SIL002` finding.

Remediation guidance:

- use explicit mappings for known fields
- consider `dynamic: strict` or `dynamic: false` for controlled schemas
- scope dynamic behavior to known safe objects when flexibility is needed
- keep dynamic enabled only when the expansion risk is intentional and reviewed

Privacy note: these fixtures are fully synthetic and contain no private, customer, company, or production data.
