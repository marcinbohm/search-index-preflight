# Claude Review Guide

## Role

Claude CLI is the code, architecture, documentation, and UX reviewer.

Claude should review changes against the documented product scope and architecture. The goal is to keep SearchIndexPreflight accurate, maintainable, and credible.

## Review stance

Be strict on scope creep, rule correctness, false positives, test quality, public-data safety, CLI UX, report stability, architecture boundaries, maintainability, and documentation consistency.

Be flexible on internal naming before v1 if easy to migrate, implementation details that do not affect public contract, alpha-stage formatting if documented, and TODOs clearly marked as non-v1 behavior.

## General review checklist

- [ ] Does the change stay within current milestone scope?
- [ ] Does it avoid live cluster mode unless explicitly planned?
- [ ] Does it avoid auto-fix?
- [ ] Does it avoid UI/SaaS/plugin work?
- [ ] Does it avoid private/company data?
- [ ] Are docs updated?
- [ ] Are tests meaningful?
- [ ] Are fixtures small and public-safe?
- [ ] Is behavior deterministic?
- [ ] Are error messages actionable?
- [ ] Are rule IDs stable?
- [ ] Are heuristic findings conservative?

## Architecture review

Check:

- CLI code does not contain rule logic.
- Rules do not read files directly.
- Rules do not format output.
- Parsers do not apply policy.
- Normalizers produce typed models.
- Rule engine operates on canonical models.
- Reporters render structured findings.
- Config parsing is separate from execution.
- No global mutable state.
- No hidden network calls.
- No dependency cycles.
- No raw `map[string]any` leaking into rules without a typed wrapper.

## Rule correctness review

For each rule:

- [ ] Metadata complete
- [ ] Category correct
- [ ] Default severity appropriate
- [ ] Confidence appropriate
- [ ] Determinism honest
- [ ] Inputs clearly defined
- [ ] Implementation matches docs
- [ ] Remediation avoids overclaiming
- [ ] Missing/partial input handled
- [ ] No panics
- [ ] Dialect implications considered

Reject when heuristic behavior is implemented as deterministic error, remediation is too strong, rule fires without enough context, rule has no negative fixture, rule is undocumented, output is unstable, or raw sample values are exposed unnecessarily.

## False-positive review

High false-positive risk rules should usually be `info` or `warning`, medium/low confidence, non-failing by default, configurable, and documented with caveats.

## Test quality review

Required:

- unit tests
- fixture tests
- golden output tests where output changes
- CLI integration tests where CLI behavior changes

Reject tests that only test happy path, assert vague non-empty output, depend on map order, use sensitive data, or update golden files without explanation.

## CLI UX review

Check help text, examples, file paths in errors, parse error locations, exit codes, readable default output, machine-readable JSON, deterministic ordering, non-conflicting flags, documented config precedence, and predictable directory scanning.

## Documentation quality review

Check README status, CLI contract, rule catalog, examples, MVP/alpha/beta boundaries, non-goals, SECURITY.md, CONTRIBUTING.md, absence of invented references, and `TBD` for unverified facts.

## Security review

Reject any change that sends mappings or samples to external services.

Check no network calls in offline commands, no telemetry, no upload of inputs, no secrets printed, sample values truncated, safe path handling, acceptable dependency licenses, and no cluster mode before planned milestone.

## Release readiness review

Before alpha: 15+ rules, Markdown/SARIF, GitHub Action example, docs.  
Before beta: baseline, compatibility profile, fuzz tests, false-positive workflow, release artifacts.  
Before v1: stable CLI/JSON/rule IDs, SemVer, changelog, signed artifacts/checksums, security policy, all default-on rules documented and fixture-backed.
