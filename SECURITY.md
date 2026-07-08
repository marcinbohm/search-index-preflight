# Security Policy

## Supported versions

SearchIndexPreflight has no stable release yet.

| Version | Supported |
|---|---|
| pre-alpha | No security support guarantee |
| alpha | Best-effort |
| beta | Best-effort |
| v1.x | Planned stable support |

## Reporting a vulnerability

Please report suspected vulnerabilities by opening a GitHub security advisory if available for this repository.

If a GitHub security advisory is not available, open a minimal public issue asking for a private security contact. Do not include vulnerability details, exploit steps, production mappings, cluster URLs, credentials, tokens, or other sensitive information in that public issue.

Do not include customer data, internal service names, private logs, full production templates, sensitive sample documents, or other confidential material in public issues.

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

If GitHub security advisories are not available for this repository, security-sensitive reports may be sent to search-index-preflight@proton.me.

Do not include vulnerability details, credentials, production mappings, cluster URLs, private logs, internal service names, or customer data in public GitHub issues.
