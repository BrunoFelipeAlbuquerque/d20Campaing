---
title: "Core Language Identity Boundary"
source: "issue"
status: "accepted"
tags: ["language", "race", "composition-boundary"]
created: "2026-05-22"
---

## Context

Core races expose automatic and bonus-language metadata, but human and half-elf bonus choices require validating any known non-secret language.

## Decision

Keep core language identity and secret-language metadata in `character/language`.

`race` may alias language IDs and store race-specific lists, but it should delegate known-language validation to the language catalog. Any-non-secret checks should use the catalog secret flag instead of character-side string tables.

## Reuse

When composing character language facts, validate selected language IDs through `language` and race-specific eligibility through `BonusLanguageChoice.AllowsLanguageID`.

Do not add class-granted languages, campaign languages, or social/literacy behavior to the P7 race-language slice.

## Verification

`go test ./...` covers the core common language catalog, Druidic as secret, non-secret helper checks, race validation through the catalog, and fixed-list versus any-non-secret bonus-language eligibility.
