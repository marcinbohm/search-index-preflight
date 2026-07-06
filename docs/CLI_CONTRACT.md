# CLI Contract

## Status

Draft contract.

Implemented behavior must match this contract unless an ADR updates it.

Current pre-alpha implementation supports `lint`, `version`, `rules list` as a stub, and `explain` as a stub.

Current `lint` behavior:

- accepts JSON mappings/templates and JSONL/NDJSON sample documents
- discovers `.json`, `.jsonl`, and `.ndjson` files in directory mode
- supports `--format console|json`
- reports parse and normalization diagnostics only

Planned but not implemented:

- real SIL rule findings
- YAML input
- Markdown output
- SARIF output
- config loading
- suppressions
- baseline
- diff
- cluster commands

## Command overview

```text
search-index-lint
  lint
  rules list
  explain
  baseline create
  diff
  version
```

Future:

```text
search-index-lint cluster inspect
search-index-lint cluster simulate
```

## Global flags

```text
--config <path>          Path to config file
--format <format>        Output format: console, json, markdown, sarif
--output <path>          Output file path; default stdout
--dialect <engine>       elasticsearch or opensearch
--version <version>      Dialect version, for example 8.x, 8.13, 2.x
--fail-on <severity>     info, warning, error, critical
--quiet                  Only print summary or errors
--verbose                Print detailed diagnostics
--no-color               Disable terminal colors
--debug                  Print debug diagnostics
--help                   Show help
```

## `search-index-lint lint`

Examples:

```bash
search-index-lint lint --mapping mapping.json
search-index-lint lint --template index-template.json
search-index-lint lint --component-template component-template.json
search-index-lint lint --sample-docs samples.jsonl
search-index-lint lint --mapping mapping.json --sample-docs samples.jsonl
search-index-lint lint ./schemas
search-index-lint lint ./schemas --format markdown              # planned
search-index-lint lint ./schemas --format sarif --output search-index-lint.sarif  # planned
```

Flags:

```text
--mapping <path>                 Mapping JSON file; YAML planned
--template <path>                Index template JSON file; YAML planned
--component-template <path>      Component template JSON file; YAML planned
--sample-docs <path>             JSONL sample documents
--config <path>                  Config file; planned
--format <format>                console, json; markdown and sarif planned
--output <path>                  Output file
--fail-on <severity>             Minimum severity that returns exit code 1
--baseline <path>                Baseline file; planned
--baseline-mode <mode>           hide_existing, report_existing, fail_on_new; planned
--disable-rule <id>              Disable rule; repeatable; planned
--enable-rule <id>               Enable rule; repeatable; planned
--only-rule <id>                 Run only selected rule; repeatable; planned
--include <glob>                 Include glob; repeatable
--exclude <glob>                 Exclude glob; repeatable
--max-sample-docs <n>            Limit sample docs loaded per file
--strict                         Treat warnings as errors
```

MVP required flags:

- `--mapping`
- `--template`
- `--component-template`
- `--sample-docs`
- `--format`
- `--output`
- `--fail-on`
- directory argument

Beta flags:

- `--baseline`
- `--baseline-mode`

## Input formats

### Mapping file

Raw mapping:

```json
{
  "dynamic": "strict",
  "properties": {
    "status": {
      "type": "keyword"
    }
  }
}
```

Wrapped mapping:

```json
{
  "mappings": {
    "dynamic": "strict",
    "properties": {
      "status": {
        "type": "keyword"
      }
    }
  }
}
```

### Index template file

```json
{
  "index_patterns": ["logs-*"],
  "priority": 200,
  "template": {
    "settings": {
      "index.mapping.total_fields.limit": 1000
    },
    "mappings": {
      "dynamic": "strict",
      "properties": {
        "@timestamp": {
          "type": "date"
        }
      }
    }
  },
  "composed_of": ["logs-common"]
}
```

### Component template file

```json
{
  "template": {
    "settings": {
      "analysis": {
        "normalizer": {
          "lowercase": {
            "type": "custom",
            "filter": ["lowercase"]
          }
        }
      }
    },
    "mappings": {
      "properties": {
        "service.name": {
          "type": "keyword",
          "normalizer": "lowercase"
        }
      }
    }
  }
}
```

### Sample documents

JSONL, one JSON object per line:

```jsonl
{"@timestamp":"2026-07-06T10:00:00Z","status":"ok","user":{"id":"u-123"}}
{"@timestamp":"2026-07-06T10:01:00Z","status":"error","user":{"id":"u-456"}}
```

## Directory input

```bash
search-index-lint lint ./schemas
```

Directory scan behavior:

- use config globs when config is present; planned
- currently include only `.json`, `.jsonl`, and `.ndjson`
- infer JSON document kind from top-level keys
- ignore hidden directories by default
- ignore `vendor`, `node_modules`, `.git`, `dist`, `build`, and `.local` by default
- report ambiguous files as diagnostics

## Output formats

- console: default human-readable output
- JSON: stable machine-readable output
- Markdown: PR comment or saved report; planned
- SARIF: GitHub code scanning, alpha scope; planned

JSON skeleton:

```json
{
  "schema_version": "0.1",
  "tool": {
    "name": "SearchIndexLint",
    "version": "0.1.0"
  },
  "dialect": {
    "engine": "elasticsearch",
    "version": "8.x"
  },
  "summary": {
    "files_scanned": 6,
    "findings_total": 3,
    "critical": 0,
    "error": 1,
    "warning": 2,
    "info": 0,
    "exit_code": 1
  },
  "findings": [],
  "diagnostics": []
}
```

## Config file

Default names:

- `search-index-lint.yaml`
- `search-index-lint.yml`
- `.search-index-lint.yaml`

Example config is provided in `search-index-lint.example.yaml`.

## Exit codes

| Code | Meaning |
|---:|---|
| 0 | Success; no findings at or above fail threshold |
| 1 | Findings at or above fail threshold |
| 2 | Invalid CLI usage |
| 3 | Parse/config/input error |
| 4 | Unsupported dialect/version |
| 5 | Cluster connection/auth error; future cluster mode only |
| 6 | Internal error |
| 7 | Baseline mismatch/corruption |

## Severity thresholds

Severity order:

```text
info < warning < error < critical
```

Default:

```text
--fail-on error
```

`--strict` is equivalent to:

```text
--fail-on warning
```

Heuristic low-confidence findings should not fail CI by default.

## Suppressions

Required fields:

- `rule`
- `reason`

Recommended fields:

- `file`
- `path`
- `owner`
- `expires`

Suppressions without reasons are invalid.

## Baseline mode

Beta scope.

Commands:

```bash
search-index-lint baseline create ./schemas --output search-index-lint.baseline.json
search-index-lint lint ./schemas --baseline search-index-lint.baseline.json
```

Modes:

- `hide_existing`
- `report_existing`
- `fail_on_new`

## Future cluster commands

Not MVP.

```bash
search-index-lint cluster inspect --url "$ES_URL" --index "logs-*"
search-index-lint cluster simulate --url "$ES_URL" --index-name "logs-app-2026.07.06"
```

Cluster commands must never write to a cluster.
