# Changelog

All notable changes to SearchIndexPreflight will be documented in this file.

## Unreleased

_No unreleased changes._

## v0.0.1-prealpha - 2026-07-07

Experimental pre-alpha source release. No release binaries are included.

### Added

- Offline `lint` command for static schema checks.
- Experimental `diff --base <path> --current <path>` command for schema diff checks.
- Static lint rules: `SIL001`, `SIL002`, `SIL003`.
- Diff rules: `DIF001`, `DIF002`, `DIF003`.
- `rules list` metadata command.
- `explain <RULE_ID>` metadata command.
- Console and JSON output.
- Public fixtures and practical examples.
- Getting started, alpha readiness, and release checklist docs.

### Changed

- Public documentation is positioned around PR-time preflight for Elasticsearch/OpenSearch schema changes.
- Future config example is explicitly marked as not implemented.

### Not included

- Production-ready CI gate guarantees.
- Release binaries.
- Config loading.
- Suppressions.
- Baseline mode.
- YAML input.
- SARIF/Markdown output.
- GitHub Action release integration.
- Cluster doctor.
- Migration/versioning implementation.
