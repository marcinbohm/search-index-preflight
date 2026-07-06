# SIL001 Total Fields Limit Fixtures

These fixtures exercise `SIL001` (`total-fields-limit-risk`), which detects explicit mappings/templates whose normalized field count approaches or exceeds the default total fields limit.

Current defaults:

- total fields limit: `1000`
- warning threshold: `800`

Counted fields:

- properties
- multi-fields
- runtime fields

Not counted or estimated yet:

- dynamic mapping expansion
- component template composition
- live cluster state
- configured `index.mapping.total_fields.limit`

Fixture cases:

- `mapping-near-limit.json`: exactly 800 keyword fields; emits one `warning` finding and exits `0` with the default `--fail-on error`.
- `mapping-over-limit.json`: exactly 1000 keyword fields; emits one `error` finding and exits `1` with the default `--fail-on error`.

Remediation guidance:

- reduce explicit field count
- restrict dynamic mappings
- consider `flattened` or `flat_object` only when query semantics fit
- split unrelated data into separate indices
- raise `index.mapping.total_fields.limit` only with operational review

Privacy note: these fixtures are fully synthetic and contain no private, customer, company, or production data.
