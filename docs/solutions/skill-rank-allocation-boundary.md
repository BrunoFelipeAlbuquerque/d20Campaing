---
title: "Skill Rank Allocation Boundary"
source: "backlog"
status: "accepted"
tags: ["skill", "character", "composition-boundary"]
created: "2026-05-23"
---

## Context

P8 introduces character skill-rank allocation facts after the `skill` domain learned grouped skill specializations such as `Knowledge (arcana)` and `Craft (alchemy)`.

Existing feat prerequisite inputs may still use grouped family IDs such as `Craft` or `Profession` because some core feats ask for ranks in any member of a family.

## Decision

Keep skill-rank allocation facts stricter than feat prerequisite rank facts.

A character allocation must point to one concrete skill identity. Ungrouped core skills like `Ride` are concrete as seeded. Grouped skills require a specialization such as `Knowledge (arcana)`; allocating ranks directly to `Knowledge`, `Craft`, `Perform`, or `Profession` should fail.

Do not change the existing feat prerequisite input surface in the allocation chassis slice. Conversion from concrete allocations into prerequisite facts belongs to the later P8 feed-allocated-ranks item.

## Reuse

When later composing skill-rank budgets, level caps, or feat prerequisite inputs, keep the allocation path concrete and preserve grouped-family references only where the target rule explicitly asks for any member of that family.

Do not add final skill check totals, armor check penalties, favored class bonuses, traits, or skill-use behavior to the allocation fact slice.

## Verification

`go test ./...` covers core skill allocation, grouped specializations for Craft, Knowledge, Perform, and Profession, invalid ranks, unknown skills, grouped-family misuse, malformed grouped IDs, and zero-value resolution failure.
