# Getting started

This guide runs SearchIndexPreflight from source. Release binaries are not published yet.

## 1. Clone

```bash
git clone https://github.com/marcinbohm/search-index-preflight.git
cd search-index-preflight
```

Or install the current source-only pre-release with Go:

```bash
go install github.com/marcinbohm/search-index-preflight/cmd/search-index-preflight@v0.0.1-prealpha
search-index-preflight version
```

Expected:

```
SearchIndexPreflight version 0.0.1-prealpha
```

## 2. Run tests

```bash
go test ./...
```

## 3. Run a lint example

```bash
go run ./cmd/search-index-preflight lint \
  --mapping fixtures/dynamic-templates/sil003-missing-match-mapping-type/mapping-missing-match-mapping-type.json
```

Expected result: one `SIL003` warning about a dynamic template missing `match_mapping_type`.

## 4. Run a diff example

```bash
go run ./cmd/search-index-preflight diff \
  --base fixtures/diff/mixed-field-changes/base \
  --current fixtures/diff/mixed-field-changes/current
```

Expected result: one `DIF001` error, one `DIF002` warning, and one `DIF003` info finding. This command exits with status 1 because `DIF001` is an error finding.

## 5. Inspect available rules

```bash
go run ./cmd/search-index-preflight rules list
go run ./cmd/search-index-preflight rules list --format json
```

## 6. Explain a rule

```bash
go run ./cmd/search-index-preflight explain SIL001
go run ./cmd/search-index-preflight explain DIF003 --format json
```

## Next docs

- [CLI contract](CLI_CONTRACT.md)
- [Rule catalog](RULE_CATALOG.md)
- [Fixtures](FIXTURES.md)
- [Architecture](ARCHITECTURE.md)
