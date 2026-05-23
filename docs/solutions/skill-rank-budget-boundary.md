---
title: "Skill Rank Budget Boundary"
source: "backlog"
status: "accepted"
tags: ["skill", "class", "race", "composition-boundary"]
created: "2026-05-23"
---

## Context

P8 composes a character skill-rank budget from selected class levels, composed Intelligence, and currently modeled racial metadata.

PF1 class entries define skill ranks per level as class ranks plus Intelligence modifier, with at least 1 rank per level. Human `Skilled` adds one additional skill rank at first level and every later level.

## Decision

Keep skill-rank budget composition in `character`.

The budget fact should use seeded class skill-rank metadata, the composed Intelligence modifier, and the selected race metadata. It should apply the one-rank-per-level floor to class-derived ranks, then add supported racial skill-rank bonuses such as human `Skilled` based on total character level.

Do not include favored class bonuses, traits, magic item skill ranks, retroactive Intelligence rebuild behavior, final skill check totals, or armor penalties in this budget slice.

## Reuse

Later rank validation should consume the budget fact instead of recalculating class, Intelligence, and racial rank rules. Keep concrete allocation validation separate from budget composition.

If additional racial or class-feature rank bonuses are modeled later, add them through explicit metadata instead of string matching feature labels outside their owning domains.

## Verification

`go test ./...` covers single-class and multiclass budgets, human `Skilled`, low-Intelligence one-rank minimum behavior, malformed class levels, malformed Intelligence facts, duplicate inputs, and invalid race inputs.
