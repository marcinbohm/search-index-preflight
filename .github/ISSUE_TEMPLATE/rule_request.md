---
name: Rule request
about: Propose a new SearchIndexPreflight rule
title: "rule: "
labels: ["type: rule"]
assignees: ""
---

## Rule idea

Describe the rule.

## Problem it catches

What Elasticsearch/OpenSearch schema or template risk should this detect?

## Why it matters

Explain the operational impact.

## Applies to

- [ ] Elasticsearch
- [ ] OpenSearch
- [ ] Both
- [ ] Not sure

## Input required

- [ ] Mapping
- [ ] Index template
- [ ] Component template
- [ ] Dynamic template
- [ ] Settings
- [ ] Sample documents
- [ ] Multiple mappings/templates
- [ ] Cluster context; future only
- [ ] Not sure

## Determinism

- [ ] Deterministic
- [ ] Heuristic
- [ ] Cluster-context-required
- [ ] Not sure

## Suggested severity

- [ ] info
- [ ] warning
- [ ] error
- [ ] critical
- [ ] Not sure

## False-positive risk

- [ ] Low
- [ ] Medium
- [ ] High
- [ ] Not sure

## Bad input example

Use a minimal synthetic example only.

```json
{
}
```

## Suggested remediation

What should users do when this rule fires?

## References

Add links to public documentation only.

## Privacy check

- [ ] Examples are synthetic
- [ ] No production mappings
- [ ] No logs
- [ ] No customer data
- [ ] No credentials/tokens
- [ ] No internal service names
