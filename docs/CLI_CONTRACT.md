# CLI Contract

## Status

Draft contract.

Implemented behavior must match this contract unless an ADR updates it.

Current pre-alpha implementation supports `lint`, minimal experimental `diff`, `version`, `rules list` as a stub, and `explain` as a stub.

Current `lint` behavior:

- accepts JSON mappings/templates and JSONL/NDJSON sample documents
- discovers `.json`, `.jsonl`, and `.ndjson` files in directory mode
- supports `--format console|json`
- reports parse and normalization diagnostics
- runs the built-in rule registry with `SIL001`, `SIL002`, and `SIL003`
- emits findings for `SIL001` total fields limit risk
- emits heuristic warning findings for `SIL002` root dynamic enabled
- emits heuristic warning findings for `SIL003` dynamic templates missing `match_mapping_type`

Current `diff` behavior:

- accepts `--base <path>` and `--current <path>`
- paths may be JSON files or directories containing `.json`, `.jsonl`, and `.ndjson`
- parses and normalizes both inputs
- compares normalized corpora with `internal/diff`
- emits `DIF001` field type changed findings
- supports `--format console|json`, `--output`, and `--fail-on`

Planned but not implemented:

- SIL004 and the rest of the rule catalog
- YAML input
- Markdown output
- SARIF output
- config loading
- suppressions
- baseline
- git-aware diff options
- diff rules beyond `DIF001`
- cluster commands

## Current vs Planned Commands

| Command | Status | Notes |
|---|---|---|
| `search-index-preflight lint` | Current; future compatibility alias | Static checks over supplied mappings/templates/sample docs. |
| `search-index-preflight diff` | Current experimental | Minimal old/new schema comparison; currently emits only `DIF001`. |
| `search-index-preflight version` | Current | Prints version information. |
| `search-index-preflight rules list` | Current stub | Command exists; full rule listing UX is not complete. |
| `search-index-preflight explain` | Current stub | Command exists; full rule explanation UX is not complete. |
| `search-index-preflight check` | Planned | Future preferred name for static checks. |
| `search-index-preflight doctor` | Planned later | Future read-only cluster inspection mode. |

The project, Go module, and binary are now named `search-index-preflight`. `check` and `doctor` remain planned.

No command may perform cluster write operations.

## Command overview

```text
search-index-preflight
  lint
  diff
  rules list
  explain
  version
```

Future:

```text
search-index-preflight check ./schemas
search-index-preflight doctor --url http://localhost:9200 --pattern "logs-*"
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

## `search-index-preflight lint`

Examples:

```bash
search-index-preflight lint --mapping mapping.json
search-index-preflight lint --template index-template.json
search-index-preflight lint --component-template component-template.json
search-index-preflight lint --sample-docs samples.jsonl
search-index-preflight lint --mapping mapping.json --sample-docs samples.jsonl
search-index-preflight lint ./schemas
search-index-preflight lint ./schemas --format markdown              # planned
search-index-preflight lint ./schemas --format sarif --output search-index-preflight.sarif  # planned
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

## `search-index-preflight diff`

Examples:

```bash
search-index-preflight diff --base old/ --current new/
search-index-preflight diff --base old/mapping.json --current new/mapping.json --format json
search-index-preflight diff --base fixtures/diff/dif001-field-type-changed/base --current fixtures/diff/dif001-field-type-changed/current
```

Flags:

```text
--base <path>       Base schema file or directory
--current <path>    Current schema file or directory
--format <format>   console or json
--output <path>     Output file
--fail-on <severity> Minimum severity that returns exit code 1
```

Current limitations:

- emits only `DIF001 field-type-changed`
- explicit file-vs-file inputs are compared as one logical resource, even when filenames differ
- directory-vs-directory inputs are matched by relative path
- file-vs-directory inputs are path-based and limited
- no rename detection
- no git refs or `--base origin/main`
- no Markdown or SARIF output
- no settings, aliases, dynamic template, template priority, composed template, sample document, or cluster-backed comparison
- no doctor/oracle/engine-backed validation

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
search-index-preflight lint ./schemas
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
    "name": "SearchIndexPreflight",
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

- `search-index-preflight.yaml`
- `search-index-preflight.yml`
- `.search-index-preflight.yaml`

Example config is provided in `search-index-preflight.example.yaml`.

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
search-index-preflight baseline create ./schemas --output search-index-preflight.baseline.json
search-index-preflight lint ./schemas --baseline search-index-preflight.baseline.json
```

Modes:

- `hide_existing`
- `report_existing`
- `fail_on_new`

## Future doctor command

Not MVP.

```bash
search-index-preflight doctor --url "$ES_URL" --pattern "logs-*"
```

Doctor mode must be read-only and must never write to a cluster.
