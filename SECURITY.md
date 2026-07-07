# Security Policy

## Supported versions

SearchIndexPreflight has no stable release yet.

| Version | Supported |
|---|---|
| pre-alpha | No security support guarantee |
| alpha | Best-effort |
| beta | Best-effort |
| v1.x | TBD |

## Reporting a vulnerability

Do not report security vulnerabilities through public GitHub issues.

Preferred reporting channel: TBD.

Until a private reporting channel is configured, open a minimal public issue saying that you need a private security contact, without including sensitive details.

Do not include production mappings, cluster URLs, credentials, tokens, customer data, internal service names, private logs, full production templates, or sensitive sample documents.

## Tool security model

SearchIndexPreflight is designed to be offline-first.

The default `lint` command must not:

- connect to clusters
- make network calls
- upload mappings
- upload sample documents
- send telemetry
- require production credentials
- write to clusters

Future cluster mode must be explicitly invoked, read-only, documented separately, safe for least-privilege credentials, and never part of default offline linting.

## Privacy considerations

Mappings, templates, and sample documents can reveal sensitive information, including internal service names, customer identifiers, tenant structure, index naming conventions, business events, security-relevant fields, infrastructure details, and logging conventions.

Do not paste confidential mappings or logs into public issues.

When reporting bugs, prefer minimal synthetic reproduction, redacted field names, fake index patterns, fake sample values, and small fixtures written from scratch.

## Handling sample documents

SearchIndexPreflight should avoid printing full sample values by default.

Findings should truncate long values, avoid printing secrets, include paths and types instead of raw data where possible, and allow verbose mode only when users explicitly request more detail.

## Dependency security

Before v1, maintainers should enable Dependabot/Renovate, CodeQL, gosec, OpenSSF Scorecard, release checksums, and signed release artifacts if practical.

## Security non-goals

SearchIndexPreflight is not a secret scanner, vulnerability scanner, cluster hardening scanner, access-control validator, compliance tool, or replacement for security review.

## Public issue warning

Public GitHub issues are public forever.

Before posting, remove confidential mappings, logs, credentials, customer data, internal names, and reduce the example to the smallest synthetic case.
