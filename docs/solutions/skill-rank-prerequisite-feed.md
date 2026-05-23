---
title: "Skill Rank Prerequisite Feed"
source: "backlog"
status: "accepted"
tags: ["skill", "feat", "composition-boundary"]
created: "2026-05-23"
---

## Context

P8 feeds validated skill-rank allocations into feat prerequisite checks. Allocation facts use concrete skill identities, including grouped specializations such as `Craft (alchemy)`, while some feat prerequisites refer to a grouped family such as any `Craft` or `Profession` skill.

## Decision

Keep the allocation path concrete, and project grouped specializations into feat prerequisite state only at the `character` composition boundary.

For grouped family prerequisites, use the highest rank value from a single concrete specialization in that family. Do not sum multiple specializations into one family rank value.

## Reuse

Future skill or feat composition should consume validated `CharacterSkillRankAllocationFacts` instead of caller-provided rank side facts when the ranks originate from character creation.

If a later rule needs a chosen grouped skill, model that chosen concrete skill explicitly rather than treating the family ID as an allocated skill.

## Verification

`go test ./...` covers mounted feat prerequisites from allocated `Ride` ranks, `Master Craftsman` from concrete `Craft` ranks, rejection of split `Craft` ranks below the single-skill prerequisite, and rejection of malformed allocation facts before they enter feat prerequisite state.
