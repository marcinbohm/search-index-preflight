# SearchIndexPreflight Project Backlog

Generated: 2026-07-08
Updated: 2026-07-08
Status: working master backlog after `v0.0.1-prealpha`

This document is the working source of truth for planned SearchIndexPreflight work.
It is intentionally broader than the first GitHub Issues batch. Not every item should
be opened as an issue immediately.

## Product north star

SearchIndexPreflight is a PR-time preflight tool for Elasticsearch/OpenSearch schema
changes.

The core question is:

> Can this mapping/template/schema change break indexing, queries, aggregations,
> rollouts, compatibility, or future maintenance before it reaches production?

The project should not drift into a generic JSON linter. Static lint is useful, but
the highest-value direction is schema-change preflight: compare old/new schema states,
surface risks early, and make review easier without connecting CI to a production
cluster.

## Current implemented baseline

Released as source-only experimental pre-alpha: `v0.0.1-prealpha`.

Implemented commands:

- `lint`
- `diff`
- `rules list`
- `explain <RULE_ID>`
- `version`

Implemented static lint rules:

- `SIL001 total-fields-limit-risk`
- `SIL002 root-dynamic-enabled`
- `SIL003 dynamic-template-missing-match-mapping-type`

Implemented diff rules:

- `DIF001 field-type-changed`
- `DIF002 field-removed`
- `DIF003 field-added`

Current positioning:

- offline-first
- source-only pre-alpha release
- no binaries yet
- no cluster access in current commands
- no telemetry
- no writes to clusters
- no config loading yet
- no suppressions yet
- no baseline yet
- no YAML yet
- no SARIF/Markdown reporters yet
- no GitHub Action wrapper yet
- no doctor/cluster mode yet
- no migration/versioning implementation yet

---

# Backlog operating model

## GitHub organization model

Use GitHub Issues as the execution layer, but keep this document as the broader map.

Recommended structure:

- Milestones = release/phase targets
- Labels = module/type/priority/status
- Epic issues = normal issues titled `[EPIC] ...` with checklists linking child issues
- GitHub Project = kanban/roadmap view

## Recommended milestones

- `v0.0.2-prealpha` — post-release cleanup, repo operations, first public issue system
- `v0.0.3-prealpha` — diff/preflight hardening
- `v0.1.0-alpha` — Markdown/PR-time reporting
- `v0.2.0-alpha` — config, rule selection, early integration work
- `v0.3.0-beta` — suppressions, baseline, adoption in existing repos
- `future` — doctor, binaries, compatibility profiles, versioning/migration concepts
- `future/offline-versioning` — dedicated future milestone for offline version-chain work

## Recommended labels

Type labels:

- `type: epic`
- `type: feature`
- `type: rule`
- `type: docs`
- `type: test`
- `type: ci`
- `type: chore`
- `type: design`
- `type: release`

Area labels:

- `area: cli`
- `area: lint`
- `area: diff`
- `area: reports`
- `area: config`
- `area: suppressions`
- `area: baseline`
- `area: docs`
- `area: examples`
- `area: release`
- `area: github-action`
- `area: sarif`
- `area: doctor`
- `area: migration-versioning`
- `area: compatibility`
- `area: parser`
- `area: normalizer`
- `area: samples`

Priority labels:

- `priority: p0`
- `priority: p1`
- `priority: p2`
- `priority: p3`

Status labels:

- `status: ready`
- `status: needs-design`
- `status: blocked`
- `status: future`

Community labels:

- `good first issue`
- `help wanted`


---

# Governance and decision framework

This backlog is intentionally broad. It is not a commitment to implement every item soon.
The purpose is to prevent important ideas from getting lost while keeping active GitHub
Issues focused and executable.

## Backlog governance

Use the following separation:

- `docs/PROJECT_BACKLOG.md` = complete project map and memory.
- GitHub Issues = active execution layer for near-term work.
- Epic issues = grouped work packages with child issue checklists.
- Milestones = release or phase targets.
- Roadmap issue = public index linking the most important epics.

Rules:

- Do not open every backlog item as an issue immediately.
- Create issues only when the item is actionable or intentionally marked as a future candidate.
- Update this backlog when a strategic decision changes direction.
- Update GitHub issues when work moves from “concept” to “execution”.
- Keep current/future boundaries explicit in README and docs.

## Definition of Ready

A feature or rule issue is ready for implementation only when it has:

- clear user problem or failure mode;
- example input, preferably mapping/template/sample/diff fixture;
- expected output or finding shape;
- classification: `SIL`, `DIF`, report, config, doctor, compatibility profile, docs, or out of scope;
- severity and confidence proposal;
- determinism classification;
- false-positive risk assessment;
- source link or explicit rationale;
- acceptance criteria;
- known non-goals.

For rule issues, the minimum ready-state card should answer:

```text
Failure mode:
Source / evidence:
Offline detectability:
Rule family: SIL | DIF | doctor | profile | out-of-scope
Severity:
Confidence:
Determinism:
False-positive risk:
Example bad input:
Expected finding:
Remediation guidance:
Dialect/version dependency:
```

## Definition of Done

A public behavior change is done only when the relevant items are complete:

- implementation;
- unit tests;
- CLI tests where applicable;
- fixtures and expected output;
- docs update;
- `rules list` metadata update if a rule is added;
- `explain <RULE_ID>` coverage if a rule is added;
- README update if user-facing behavior changes;
- CLI contract update if command/flag/output behavior changes;
- rule catalog update if a rule changes;
- no stale “planned/not implemented” wording for implemented behavior;
- no overclaim about production readiness, config, YAML, SARIF, binaries, doctor, or migration/versioning;
- `go test ./...` and `go vet ./...` pass.

## Rule discovery decision record

Before adding many new rules, use short decision records to avoid rule sprawl.
Each candidate rule should record:

- what breaks in Elasticsearch/OpenSearch;
- whether the failure happens during indexing, query, aggregation, template rollout, alias/data stream operation, or future compatibility;
- whether the tool can detect it offline;
- whether it belongs in `SIL`, `DIF`, `doctor`, compatibility profile, config policy, or out of scope;
- how reliable the detection is;
- how noisy the rule may be;
- whether the behavior depends on Elasticsearch/OpenSearch version;
- what fixture proves the behavior;
- what remediation should be shown.

This avoids implementing rules only because they sound useful. The project should prefer
rules that are grounded in official engine behavior, documented limits, real rollout risks,
or common incident patterns.

## Rule prioritization scoring

When choosing between possible rules, score them roughly using:

```text
priority = user impact + likelihood + offline detectability + confidence - false-positive risk - implementation cost
```

Suggested dimensions:

- User impact: would this prevent a painful production incident or PR rollback?
- Likelihood: how often do teams make this mistake?
- Offline detectability: can we detect it without a cluster?
- Confidence: can we report it without guessing too much?
- False-positive risk: will users quickly ignore it?
- Implementation cost: how much model/parser/diff/report work is required?
- Documentation quality: can we cite or explain the behavior clearly?

Use this scoring especially when choosing between:

- `DIF004 field-options-changed`;
- dynamic setting diff;
- dynamic template diff;
- template priority/composition diff;
- `SIL007 dotted-field-collision`;
- `SIL009 sample-doc-mapping-conflict`;
- compatibility/profile rules.

## False-positive management

False-positive control is a product feature, not an afterthought.
For every heuristic rule, track:

- why it is heuristic;
- expected false-positive cases;
- default severity;
- whether it should ever fail CI by default;
- how to explain the finding without sounding certain when it is not;
- when the rule should be disabled by default;
- how user feedback should change severity or wording.

Guidelines:

- Deterministic `error` findings may fail default CI thresholds.
- Heuristic findings should usually be `warning` or `info`.
- If a warning is noisy, prefer sharper applicability over more suppressions.
- Suppressions and baseline are adoption tools, not a substitute for precise rules.

## Compatibility data lifecycle

Dialect/version profiles can become a major differentiator, but only if maintained carefully.
Plan the lifecycle before implementing version-specific behavior.

Questions to answer in the compatibility RFC:

- Where do profiles live: code, embedded JSON data, generated files, or docs-first tables?
- How is `latest` defined?
- How often is `latest` updated?
- How are Elasticsearch and OpenSearch divergences represented?
- How are uncertain or undocumented behaviors marked?
- How are profile changes tested?
- Can users override or extend profiles?
- How do reports show profile source and rule applicability?
- How are breaking profile-data changes released?

Do not scatter direct version checks across rules before this model exists.

## Report schema and public contract policy

JSON report output is a public contract once users integrate it with CI.
Define the policy before expanding machine-readable outputs.

Questions:

- Which outputs are versioned contracts?
- Does `rules list --format json` need its own schema version later?
- Does `explain --format json` need its own schema version later?
- When should `schema_version` change?
- What counts as a breaking report change?
- How are deprecated fields handled?
- Should a JSON Schema file be published?
- How do Markdown and SARIF relate to the JSON report model?

Until a v1 contract exists, docs should clearly say which formats are experimental.

## Data, privacy, and safety policy

Search schema files and sample documents can contain sensitive information.
The project should keep strong safety defaults:

- no telemetry;
- no upload of mappings, templates, or sample documents;
- no network calls during `lint` or `diff`;
- no cluster writes ever;
- no auto-fix;
- examples and fixtures must use synthetic data;
- issue templates should warn users not to paste production mappings, credentials, cluster URLs, private logs, customer data, or sensitive sample documents;
- future doctor/cluster features must be explicit, read-only, and safe by default.

This policy should remain visible in README, SECURITY, CONTRIBUTING, and issue templates.

## User journeys and personas

Use user journeys to guide UX and priority decisions.

### Platform/search engineer reviewing a schema PR

Wants:

- see whether a mapping/template change is risky;
- get concise PR-friendly output;
- understand remediation;
- avoid giving CI production cluster credentials.

Important features:

- `diff`;
- Markdown report;
- GitHub Action later;
- deterministic findings;
- good examples.

### Application developer changing a mapping

Wants:

- understand why the change is risky;
- know what to change before review;
- run locally without learning the whole tool.

Important features:

- quick demo;
- good remediation;
- `explain <RULE_ID>`;
- clean console output.

### SRE maintaining index templates

Wants:

- catch template rollout risks;
- compare old/new template sets;
- understand aliases/data streams/settings changes;
- eventually inspect cluster state read-only.

Important features:

- template-level diff;
- settings diff;
- alias diff;
- read-only `doctor` later.

### Team adopting tool in a legacy repo

Wants:

- avoid being blocked by existing findings;
- fail only on new risks;
- suppress known intentional findings with reasons.

Important features:

- baseline;
- suppressions;
- stable fingerprints;
- config.

### Open-source contributor adding a rule

Wants:

- clear rule authoring guide;
- fixture pattern;
- test expectations;
- rule metadata expectations;
- docs checklist.

Important features:

- contributor docs;
- good-first issues;
- rule request template.

## Release/channel policy

Use release channels to keep expectations honest.

### Pre-alpha

Purpose:

- source-only evaluation;
- synthetic examples;
- feedback on core idea and CLI behavior.

Not expected:

- production CI adoption;
- binaries;
- stable JSON schema;
- broad rule coverage.

### Alpha

Purpose:

- useful PR-time review flow;
- Markdown/SARIF direction;
- early GitHub Action or manual PR integration;
- more diff/preflight coverage.

### Beta

Purpose:

- adoption in existing repositories;
- baseline and suppressions;
- compatibility profiles start becoming useful;
- release artifacts may become important.

### v1

Purpose:

- stable CLI contract;
- stable report schema;
- documented SemVer policy;
- signed/checksummed releases;
- stable rule ID policy;
- mature contribution process.

## Documentation governance

Keep documentation layered:

- README = first 60 seconds and entry points;
- GETTING_STARTED = runnable first path;
- CLI_CONTRACT = exact behavior;
- RULE_CATALOG = rule truth;
- PROJECT_BACKLOG = broad memory and strategy;
- GitHub Issues = execution;
- ADRs = durable decisions;
- release notes = what changed in a specific release.

Whenever code changes public behavior, update the relevant docs in the same sprint.

---

# Rule discovery strategy: SIL vs DIF

## Core principle

We should not guess the final number of `SIL` and `DIF` rules from intuition alone.
The right answer to “what exactly and how many rules do we need?” should come from a
structured audit of Elasticsearch/OpenSearch failure modes.

## DIF rule source of truth

`DIF` rules should be derived primarily from real schema-evolution failure modes:
what can break when an existing mapping/template/schema changes between base and
current states.

Research inputs:

- Elasticsearch official documentation
- OpenSearch official documentation
- release notes and breaking changes
- mapping parameter docs
- index settings docs
- index template and component template docs
- data stream and alias behavior
- field capabilities and template simulation behavior
- real-world incident patterns from search/platform teams

Example failure-mode families:

- field type changes
- field removal
- field option changes without type changes
- analyzer/normalizer/search_analyzer changes
- `index` / `doc_values` / `store` / `norms` changes
- `ignore_above`, `null_value`, `coerce`, `format`, `ignore_malformed` changes
- dynamic mapping policy changes
- dynamic template matching and ordering changes
- index template priority / patterns / composed_of changes
- component template changes
- total fields / depth / nested limits changes
- aliases and write index changes
- data stream configuration changes
- runtime field changes
- dialect-specific unsupported mappings

## SIL rule source of truth

`SIL` rules are harder because they are not necessarily about a delta. They should be
based on:

- official engine limits
- mapping best practices
- operational risk patterns
- documented Elasticsearch/OpenSearch behavior
- conservative static heuristics
- sample-document consistency checks where available

Example static-risk families:

- mapping explosion / total fields / depth / nested limits
- broad dynamic mappings
- risky dynamic templates
- dotted field / object collisions
- text vs keyword mistakes
- fielddata on text
- analyzer/normalizer references
- template conflicts
- sample-document conflicts
- unsupported field types by dialect
- source/runtime/nested/object risk

## Research epic needed

Before adding many new rules, create a dedicated research/design epic:

> `[EPIC] Elasticsearch/OpenSearch failure-mode audit for predeploy schema safety`

Deliverables:

- failure-mode taxonomy
- mapping from failure mode -> possible `DIF` rule / `SIL` rule / doctor check / out of scope
- source links to official ES/OS docs
- severity model
- determinism classification
- false-positive risk classification
- fixtures needed
- prioritization by user value and implementation cost

## Dialect and version profile strategy

`SIL` and `DIF` rules should not assume that every Elasticsearch/OpenSearch version behaves identically. The default profile may target `latest`, but the tool should eventually support explicit dialect/version profiles, for example:

```text
elasticsearch/latest
elasticsearch/8.x
elasticsearch/8.13
opensearch/latest
opensearch/2.x
```

Why this matters:

- some mapping parameters, field types, defaults, limits, or settings differ between Elasticsearch and OpenSearch;
- some risky patterns are version-specific;
- some checks may be valid for latest Elasticsearch but wrong or incomplete for a specific OpenSearch version, or the other way around;
- users currently have to verify this manually against engine docs, which is exactly the kind of predeploy burden SearchIndexPreflight should reduce over time.

Future behavior should allow:

- defaulting to a conservative `latest` profile when no dialect/version is selected;
- selecting dialect/version through CLI/config, e.g. `--dialect elasticsearch --engine-version 8.x`;
- rule applicability by profile;
- profile-specific severity, remediation, supported mapping parameters, supported field types, and known limitations;
- config-level overrides or extensions for teams with internal engine constraints;
- clear reporting when a rule is skipped, downgraded, or adjusted because of the selected profile.

This must be designed before implementing version-specific behavior. Do not hardcode scattered version checks directly inside rules without a profile model.

---

# A. Project and repository operations

## A1. Update README after source-only release

The README should reflect that `v0.0.1-prealpha` exists as a source-only GitHub
pre-release, while still stating that there are no binary artifacts.

Tasks:

- Update install/run section after release.
- Mention source-only pre-release.
- Do not claim binaries.
- Keep pre-alpha warning.
- Keep first 30 seconds focused on value.

Acceptance criteria:

- README says source-only release exists.
- README says no binaries are published.
- README remains honest about pre-alpha status.

## A2. Verify and document `go install`

Verify:

```bash
go install github.com/marcinbohm/search-index-preflight/cmd/search-index-preflight@v0.0.1-prealpha
```

Tasks:

- Test from a clean environment.
- Confirm installed binary runs.
- Confirm `search-index-preflight version` reports `0.0.1-prealpha`.
- Add to README/GETTING_STARTED only if verified.

## A3. Configure GitHub repository metadata

Manual GitHub UI tasks:

- description: `Offline-first preflight CLI for Elasticsearch/OpenSearch schema changes`
- topics:
  - `elasticsearch`
  - `opensearch`
  - `search`
  - `schema`
  - `mapping`
  - `index-template`
  - `ci`
  - `preflight`
  - `linter`
  - `golang`
  - `devops`
  - `sre`
  - `schema-as-code`
- verify CI badge
- verify issue templates
- enable security advisories/private vulnerability reporting if available

## A4. Create GitHub labels, milestones, and project board

Tasks:

- Create labels from this document.
- Create milestones from this document.
- Create GitHub Project `SearchIndexPreflight Roadmap`.
- Add columns/statuses:
  - Backlog
  - Ready
  - In progress
  - Review
  - Done
  - Blocked

## A5. Create public roadmap issue index

Create issue:

```text
[ROADMAP] SearchIndexPreflight public roadmap
```

Contents:

- current release/status
- implemented rules
- next focus: diff/preflight core
- links to epic issues
- no dates promised

## A6. Maintain release process

Tasks:

- Keep `docs/RELEASE_CHECKLIST.md` current.
- Keep `CHANGELOG.md` current.
- Keep `docs/releases/` drafts current.
- Decide when to add binaries.
- Decide whether to use GoReleaser later.
- Add checksums/signing only in later release work.

---

# B. CLI contract and UX

## B1. CLI help consistency

Tasks:

- Review help output for all commands.
- Ensure examples match implemented behavior.
- Ensure planned flags are not shown as implemented.
- Add tests for help output where useful.

## B2. Argument parsing consistency

Tasks:

- Review `lint`, `diff`, `rules list`, `explain` flag parsing.
- Decide whether to support flags before/after positional args consistently.
- Do not break existing documented examples.

## B3. Exit code model

Tasks:

- Document exit code semantics clearly.
- Ensure error findings fail with default threshold.
- Ensure warning/info behavior is documented.
- Add tests for threshold behavior across lint/diff.

## B4. Stderr/stdout policy

Tasks:

- Successful reports go to stdout.
- Usage/errors go to stderr.
- `--output` behavior remains consistent.
- Add tests for output separation.

## B5. Future `check` alias

Current plan says `check` may become the preferred static-check command.

Tasks:

- RFC: keep `lint` only vs add `check` alias.
- If added, ensure no breaking change.
- Docs must explain alias clearly.

## B6. Future utility flags

Candidates:

- `--quiet`
- `--verbose`
- `--no-color`
- `--debug`

No implementation before design.

---

# C. Input discovery and parsing

## C1. JSON mapping/template parsing hardening

Tasks:

- Better malformed JSON diagnostics.
- Line/column diagnostics if feasible.
- JSON pointer accuracy.
- Raw mapping vs wrapped mapping coverage.
- Index template and component template examples.

## C2. Directory discovery hardening

Tasks:

- Deterministic file ordering.
- Include/exclude behavior.
- Hidden file policy.
- Symlink policy.
- Mixed valid/invalid files.

## C3. JSONL/NDJSON sample parsing hardening

Tasks:

- Large-file behavior.
- Max sample docs per file.
- Line-number diagnostics.
- Empty lines policy.
- Invalid JSONL line behavior.

## C4. YAML input RFC

YAML is planned, not implemented.

Tasks:

- Decide whether YAML should be supported.
- Pick parser dependency if yes.
- Define security and ambiguity policy.
- Add fixtures before implementation.

---

# D. Normalizer and model.Corpus

## D1. Mapping normalization completeness

Tasks:

- Properties.
- Multi-fields.
- Runtime fields.
- Object fields.
- Nested fields.
- Field aliases.
- JSON pointer correctness.

## D2. Index template normalization

Tasks:

- index template mappings.
- index patterns.
- priority.
- composed_of.
- data_stream.
- settings.
- aliases.

## D3. Component template normalization

Tasks:

- component template mappings.
- settings.
- aliases.
- composition relationship model.

## D4. Settings model

Tasks:

- analyzer definitions.
- normalizer definitions.
- total_fields limit.
- depth limits.
- nested limits.
- mapping options that affect lint/diff.

## D5. Alias model

Tasks:

- parse aliases.
- write index.
- filters.
- routing.
- diff support later.

---

# E. Static lint rules

## E1. Maintain implemented rules

### SIL001 total-fields-limit-risk

Future work:

- Read limit from settings when available.
- Count aliases and object mappers accurately where missing.
- Improve field count breakdown in output.
- Add more fixtures for runtime fields, aliases, nested objects, object mappers.
- Keep severity behavior clear: warning near threshold, error at/above limit.

### SIL002 root-dynamic-enabled

Future work:

- Child object dynamic settings later.
- Dynamic template interaction later.
- More examples around intentional dynamic usage.

### SIL003 dynamic-template-missing-match-mapping-type

Future work:

- Dynamic template ordering later.
- Type compatibility later.
- Better false-positive guidance.

## E2. Future static rule candidates

Candidate rules from catalog:

- `SIL004 overbroad-dynamic-template`
- `SIL005 dynamic-template-shadowing`
- `SIL006 path-match-object-collision-risk`
- `SIL007 dotted-field-collision`
- `SIL008 field-type-conflict`
- `SIL009 sample-doc-mapping-conflict`
- `SIL010 dynamic-date-numeric-detection-risk`
- `SIL011 likely-aggregatable-field-as-text`
- `SIL012 long-keyword-without-ignore_above`
- `SIL013 fielddata-true-on-text`
- `SIL014 missing-analyzer-normalizer-definition`
- `SIL015 template-priority-conflict`
- `SIL016 multi-field-expansion-risk`
- `SIL017 array-of-objects-object-mapping-risk`
- `SIL018 nested-limit-risk`
- `SIL019 keyword-likely-needs-normalizer`
- `SIL020 component-template-missing`
- `SIL021 component-template-override-conflict`
- `SIL022 legacy-composable-template-collision`
- `SIL023 data-stream-missing-timestamp`
- `SIL024 mixed-array-element-types`
- `SIL025 null-only-sample-field`
- `SIL026 mapping-depth-limit-risk`
- `SIL027 numeric-identifier-risk`
- `SIL028 runtime-fields-overuse-risk`
- `SIL029 unsupported-field-type-for-dialect`
- `SIL030 source-disabled-risk`

## E3. Static rule triage

Before implementing more SIL rules:

- Evaluate official ES/OS behavior.
- Rank by real failure risk.
- Prefer deterministic rules over broad heuristics.
- Avoid implementing `SIL004` just because it is next numerically.

---

# F. Diff / preflight core

This is the next highest-value implementation area.

## F1. Diff/preflight failure-mode RFC

Tasks:

- Audit ES/OS schema-change failure modes.
- Map failures to possible `DIF` rules.
- Select next 1-2 `DIF` rules.
- Define severity and determinism.
- Define fixtures.

## F2. Extend FieldSnapshot with mapping options

Potential options:

- `index`
- `doc_values`
- `store`
- `norms`
- `analyzer`
- `search_analyzer`
- `normalizer`
- `ignore_above`
- `ignore_malformed`
- `coerce`
- `format`
- `null_value`
- `similarity`

## F3. DIF004 field-options-changed

Potential rule:

```text
DIF004 field-options-changed
```

Purpose:

Detect important field mapping option changes where field type does not change.

Acceptance:

- registry metadata
- `rules list`
- `explain DIF004`
- console/JSON output
- fixtures
- docs

## F4. Dynamic setting diff

Potential rule:

```text
DIF005 dynamic-setting-changed
```

Examples:

- `strict -> true`
- `false -> true`
- `true -> strict`

## F5. Dynamic template diff

Potential rules:

- `DIF006 dynamic-template-added`
- `DIF007 dynamic-template-removed`
- `DIF008 dynamic-template-changed`

Need RFC before implementation because template order and shadowing matter.

## F6. Template-level diff

Potential checks:

- index_patterns changed
- priority changed
- composed_of changed
- data_stream changed
- settings changed
- aliases changed
- component template changed

## F7. Settings diff

Potential checks:

- total_fields limit changed
- depth limit changed
- nested limits changed
- analyzer/normalizer changed
- index sorting changed

## F8. Alias diff

Potential checks:

- alias added/removed
- write index changed
- filter changed
- routing changed

## F9. Rename detection policy

Tasks:

- Document current limitations.
- Decide if rename detection is needed.
- Consider resource identity/content hash.
- Do not implement early unless strong need.

## F10. Git refs support

Possible future features:

- compare current working tree to git ref
- `--base-ref`
- `--current-ref`
- local git only first

---

# G. Reports and output formats

## G1. Console UX polish

Tasks:

- summary formatting
- remediation formatting
- grouping by severity/rule/file
- optional color later
- quiet/verbose later

## G2. JSON report schema docs

Tasks:

- document schema
- versioning policy
- compatibility policy
- JSON schema file later

## G3. Markdown reporter RFC

Goal:

Design PR-friendly Markdown output.

Sections:

- summary
- findings by severity
- finding detail
- remediation
- affected files
- exit status

## G4. Implement `--format markdown`

Acceptance:

- `lint --format markdown`
- `diff --format markdown`
- `--output report.md`
- golden tests
- docs

## G5. SARIF reporter RFC

Tasks:

- SARIF rule metadata
- severity mapping
- JSON pointer to location mapping
- validator strategy

## G6. Implement `--format sarif`

Acceptance:

- SARIF validates
- GitHub code scanning compatible
- tests and fixtures
- no GitHub Action required yet

---

# H. GitHub Action integration

## H1. GitHub Action RFC

Inputs:

- mode: `lint|diff`
- schema path
- base path
- current path
- fail-on
- output format

Outputs:

- findings count
- report path
- exit code behavior

Safety:

- no cluster access
- no secrets required
- no writes

## H2. Experimental GitHub Action wrapper

Tasks:

- `action.yml`
- sample workflow
- docs
- pin version behavior
- no release automation initially

## H3. PR comment integration

Tasks:

- use Markdown reporter
- optional script example
- later GitHub Action output integration

---

# I. Config and rule selection

## I1. Config RFC

Questions:

- YAML vs JSON vs TOML
- explicit `--config` only or default search path
- unknown keys behavior
- config versioning
- future compatibility

## I2. Config loading MVP

Tasks:

- parse config
- validate config
- clear errors
- docs
- tests

## I3. Rule selection flags

Flags:

- `--only-rule`
- `--disable-rule`
- `--enable-rule`

Acceptance:

- works for implemented rules
- unknown rule error
- deterministic behavior
- docs and tests

## I4. Severity and threshold config

Tasks:

- fail-on config
- maybe rule severity override later
- document risks

## I5. Include/exclude config

Tasks:

- include globs
- exclude globs
- integration with discovery
- docs and tests

## I6. Dialect/version profile config

Future config should allow users to define or select the engine profile used by `SIL` and `DIF` rules.

Possible config shape, subject to RFC:

```yaml
engine:
  dialect: elasticsearch
  version: 8.x
  profile: latest
```

Potential extension/override shape, subject to RFC:

```yaml
profiles:
  custom-prod-es:
    dialect: elasticsearch
    version: 8.13
    overrides:
      rules:
        SIL029:
          enabled: true
```

Tasks:

- decide whether profile selection belongs in config, CLI flags, or both;
- define allowed dialect/version identifiers;
- define how rules declare applicability;
- define whether users may override rule applicability/severity;
- ensure reports include the selected profile;
- ensure default behavior remains conservative when no profile is configured.

---

# J. Suppressions and baseline


## J1. Suppressions RFC

Suppression should include:

- rule ID
- fingerprint/location
- reason
- optional expiration
- optional owner

## J2. Suppressions MVP

Tasks:

- parse suppressions
- apply suppressions
- report suppressed count
- JSON includes suppressed metadata
- docs

## J3. Baseline RFC

Baseline goal:

- allow adoption in existing repos without fixing everything immediately
- fail on new findings
- track stable fingerprints

## J4. Stable fingerprints

Inputs:

- rule ID
- normalized resource identity
- field path / JSON pointer
- message class

Need deterministic tests.

## J5. Baseline implementation

Later beta work.

---

# K. Compatibility, dialects, and version profiles

This is a major future product area. Many users have to manually check whether a mapping, setting, or template behavior is valid or risky for their exact Elasticsearch/OpenSearch version. SearchIndexPreflight should eventually reduce that burden.

## K1. Compatibility profile RFC

Scope:

- Elasticsearch versions
- OpenSearch versions
- `latest` profile semantics
- supported field types
- supported mapping parameters
- supported index settings
- analyzer/normalizer differences
- template/data stream differences
- deprecations and removals
- behavior differences that affect `SIL` and `DIF` rules

Questions:

- Should the default profile be `latest`, `unknown`, or `conservative-latest`?
- How should `latest` be updated between releases?
- Should profiles live in code, generated data files, or docs-first tables?
- How much of profile data should users be allowed to override?
- How do we report a rule that does not apply to the selected profile?

## K2. Dialect/version flags

Possible flags:

- `--dialect elasticsearch|opensearch`
- `--engine-version 8.x|2.x|latest`
- `--profile <name>`

Notes:

- Avoid overloading the existing `version` command.
- Prefer `--engine-version` or similar over bare `--version` if ambiguity is likely.
- CLI and config should resolve to one normalized engine profile.

## K3. Rule applicability by profile

Future rules should be able to declare applicability:

- applies to all profiles;
- applies only to Elasticsearch;
- applies only to OpenSearch;
- applies only below/above a version range;
- has different severity/remediation by profile;
- is skipped for unknown profile unless explicitly enabled.

This should apply to both:

- `SIL` static lint rules;
- `DIF` schema-change rules.

## K4. Profile-aware report output

Reports should eventually include:

- selected dialect;
- selected engine version/profile;
- profile source: CLI, config, default;
- rules skipped or adjusted by profile, if relevant;
- warnings when profile is unknown or defaulted to `latest`.

## K5. User-defined profile overrides

Future config may allow users to add or override profile behavior for internal constraints.

Examples:

- internal Elasticsearch fork with disabled field type;
- organization-specific mapping restrictions;
- intentionally stricter rule thresholds;
- pinned OpenSearch version where latest behavior would be misleading.

Any override system must be explicit and auditable in reports.

## K6. Unsupported field type checks

Possible rule:

- `SIL029 unsupported-field-type-for-dialect`

Only after compatibility model exists.

## K7. Version-specific diff risks

Some `DIF` rules may need profile-aware behavior. Examples:

- mapping parameter change that is safe in one engine version but breaking in another;
- field type available in latest Elasticsearch but unsupported in a target OpenSearch version;
- template/data stream behavior differences;
- default setting changes between engine versions.

Do not implement scattered one-off version checks. Add the profile model first.

---

# L. Sample document inference

## L1. Sample type inference

Tasks:

- infer primitive types
- objects
- arrays
- nulls
- mixed types
- date-like strings

## L2. SIL009 sample-doc-mapping-conflict

Goal:

Find sample documents that clearly conflict with supplied mapping.

## L3. SIL024 mixed-array-element-types

Detect mixed array element types from samples.

## L4. SIL025 null-only-sample-field

Info-level check for sample-only uncertainty.

---

# M. Doctor / read-only engine-backed validation

Future only.

## M1. Doctor safety RFC

Must define:

- read-only only
- explicit command only
- no default network calls
- credential handling
- log redaction
- no telemetry

## M2. `_field_caps` PoC

Tasks:

- mocked HTTP tests
- explicit URL input
- no writes
- compare with offline schema view

## M3. Template simulate PoC

Tasks:

- read-only simulate endpoint
- compare engine result with offline assumptions

## M4. Drift detection

Future:

- repo schema vs cluster schema
- read-only inspection

---

# N. Offline migration/versioning

Future only.

## N1. Command naming RFC

Candidates:

- `versions validate`
- `migrations validate`
- other term that does not imply apply/deploy

## N2. Input model RFC

Options:

- versioned directories
- migration manifests
- manifest + schema snapshots

## N3. Version chain validation

Checks:

- ordering
- duplicate IDs
- gaps
- missing base
- multi-index/multi-template support

## N4. Consecutive diff orchestration

Process:

- lint every version
- diff consecutive versions
- aggregate report

## N5. Migration/versioning report schema

Sections:

- chain summary
- per-step findings
- thresholds
- JSON/console
- Markdown later

## N6. Future fixtures

Only if clearly marked as not implemented.

---

# O. Architecture and internals

## O1. Package boundaries audit

Packages:

- input discovery
- parser
- normalizer
- model
- rules
- diff
- diffrules
- report
- cli

## O2. Error handling audit

Tasks:

- no panics
- wrapped errors
- deterministic diagnostics
- user-friendly messages

## O3. Performance baselines

Tasks:

- benchmark parser
- benchmark normalizer
- benchmark field collection
- large mapping fixture

## O4. Fuzz tests

Targets:

- parser
- normalizer
- diff engine

---

# P. Documentation and examples

## P1. README maintenance

Update after:

- releases
- new commands
- new rules
- new reporters
- install flow changes

## P2. Getting Started maintenance

Tasks:

- source release flow
- `go install` if verified
- examples current

## P3. Rule docs completeness

Every implemented rule should have:

- why it matters
- applicability
- bad input
- expected output
- remediation
- limitations
- fixture link

## P4. Example expansion

Future examples:

- field added
- total fields limit
- root dynamic enabled
- field option changed
- Markdown report example
- config example once implemented

## P5. Docs consistency checks

Run before each release:

- no stale stub wording
- no old config filename
- no overclaims
- no current/future confusion

---

# Q. CI and quality gates

## Q1. CI matrix

Future:

- multiple Go versions
- OS matrix
- possibly race tests later

## Q2. Test discipline

Tasks:

- unit tests
- CLI tests
- fixture tests
- golden outputs

## Q3. Coverage reporting

Future, not required now.

## Q4. Link checking

Future docs quality check.

## Q5. Golden update process

Document how expected outputs are updated.

---

# R. Distribution

## R1. Source-only release

Current approach.

## R2. Binary release RFC

Questions:

- OS/arch matrix
- checksums
- signing
- GoReleaser vs custom
- release automation

## R3. GoReleaser implementation

Future.

## R4. Homebrew tap

Future.

## R5. Docker/GHCR

Future, possibly unnecessary.

## R6. GitHub Packages

Not needed now.

---

# S. Community and contribution

## S1. Good first issues

Candidates:

- docs examples
- small fixtures
- README improvements
- rule catalog clarifications
- expected output updates

## S2. Contributor guide: add a rule

Document:

- metadata
- fixtures
- tests
- docs
- explain/rules list expectations

## S3. Contributor guide: add a diff rule

Document:

- diff engine changes
- diffrules package
- fixtures
- CLI tests
- docs

## S4. Rule request template polish

Ask for:

- risk description
- example mapping/template
- expected severity
- false positive concerns
- Elasticsearch/OpenSearch source if known

## S5. Security policy follow-up

- enable security advisories
- keep sensitive data out of issues

---

# T. Product and positioning

## T1. Positioning audit after release

Review:

- README first 30 seconds
- topics
- description
- examples
- release notes

## T2. Feedback collection

Ask users:

- Are current findings useful?
- Which rules are missing?
- Is output PR-friendly?
- What examples are unclear?

## T3. Competitive comparison

Possible docs:

- why not only `_simulate_index_template`
- why not only staging cluster
- why offline-first matters
- what this tool intentionally does not do

---

# First wave of GitHub Issues

Do not create every item above immediately.

Recommended first wave:

1. `[EPIC] Repository operations after v0.0.1-prealpha`
2. `[EPIC] Diff/preflight core`
3. `[EPIC] Reporting and PR review output`
4. `[EPIC] Config, rule selection, suppressions`
5. `[EPIC] Static lint rule backlog`
6. `[EPIC] Offline migration/versioning`
7. `Update README after v0.0.1-prealpha source release`
8. `Verify and document go install for the pre-alpha tag`
9. `Configure GitHub repository metadata`
10. `Create GitHub labels, milestones, and project board`
11. `Create public roadmap issue index`
12. `RFC: Elasticsearch/OpenSearch failure-mode audit for predeploy schema safety`
13. `RFC: Compatibility and dialect/version profile model`
14. `RFC: Choose next diff/preflight risk checks`
15. `Extend field snapshots with mapping options`
16. `Add fixtures for field option changes`
17. `RFC: Markdown report format for PR review`
18. `Convert future SIL rules into issue candidates`
19. `RFC: Minimal config contract`
20. `Create good-first-issue candidates`

# Recommended near-term execution order

## Sprint 1: GitHub organization sprint

- labels
- milestones
- project board
- epic issues
- public roadmap issue

## Sprint 2: Post-release docs/install sprint

- README after source-only release
- verify `go install`
- repository metadata
- first user path

## Sprint 3: Failure-mode and compatibility research sprint

- audit Elasticsearch/OpenSearch docs
- build failure-mode taxonomy
- map failures to `DIF`, `SIL`, doctor, config, out-of-scope
- identify version/dialect-specific behavior
- draft the compatibility/profile model before any version-specific rules

## Sprint 4: Diff RFC sprint

- choose next diff rules
- decide field options scope
- write fixtures plan

## Sprint 5: Diff implementation sprint

- extend `FieldSnapshot`
- implement first next diff rule
- add fixtures/docs/tests

## Sprint 6: Markdown reporter sprint

- design PR-friendly output
- implement `--format markdown`
- add examples

# Explicit non-priorities right now

Do not do these immediately:

- implement `SIL004` just because it is next numerically
- implement migration/versioning
- implement doctor/field_caps
- implement suppressions/baseline
- add Docker/Homebrew/GitHub Packages
- add GitHub Action before Markdown/SARIF/reporting work
- add release binaries before release RFC

