# d20campaigngenerator

Pathfinder 1e character-building and rules-domain project.

This repository is aiming for a rules-aware domain model, not just a loose bag of structs. The goal is to represent character data in a way that is:

- explicit
- validated
- core-first
- safe to compose as more Pathfinder systems are added

The project follows `BACKLOG.md` as the source of truth. Current delivery is backlog-first and core-only: finish the small validated foundation domains, then add core seed data, then add query logic, and only then expand composition under `character`.

For a compact package and workflow orientation before implementation, use [Project Map](docs/project-map.md). It is intentionally shorter than this README and should reduce repeated repo rediscovery during future tasks.

## What This Project Is Trying To Achieve

The current goal is to model Pathfinder 1e characters and creature logic in a way that stays maintainable as the core rules surface grows one slice at a time.

Pathfinder is full of systems that look simple on the surface but become hard once you combine:

- multiclassing
- racial rules
- monster rules
- exceptional cases
- later composition pressure

Because of that, this project prefers a domain-first approach:

- keep each rules concept in its own value object or small aggregate
- validate at construction time
- avoid magic numbers leaking all over the codebase
- avoid coupling rule containers to unrelated concepts too early
- preserve enough detail to explain where a number came from

This matters because many later systems will depend on these primitives:

- creature type and subtype
- classes
- race and ancestry
- spellcasting
- combat state
- monster advancement

If the early stat objects are sloppy, everything above them becomes brittle.

## Current Architectural Direction

The current design favors narrow, well-defined domain types under `internal/domain/rpg/character/ability`.

Each type should answer one specific rules question well.

Examples:

- `AbilityScore` answers what a score is, whether it is currently valid, and what modifier it produces
- `Alignment` answers which axes are present and how they should be named
- `BaseAttackBonus` answers both the exact fractional progression and the game-facing rounded value
- `SavingThrow` answers the same, while also handling the one-time good-save bonus correctly across multiclassing
- `CasterLevel` stores source-based caster level totals without deciding how classes map into them
- `HitPoints` is a ledger of raw HP totals and sources, not a combat-state engine
- `Size` is a rules table for size-derived modifiers and physical dimensions

The project is intentionally trying not to jump too quickly into giant “character” structs that know everything. That kind of design usually feels productive at first, then becomes painful once exceptions appear.

## Main Design Decisions

### Invalid Input Should Fail Construction

This is one of the most important project rules.

If a domain value is invalid, construction should fail instead of silently creating a fake default object.

In Go, this is a design philosophy and boundary goal, not a strict language-level guarantee.

The project minimizes invalid states through unexported fields, validated constructors, and narrow domain types, but zero values and partially invalid states can still exist inside the module.

In practice, that means constructors usually return:

- `(value, true)` when valid
- `(zero, false)` when invalid

Why this matters:

- invalid state is easier to catch near the source
- later code can trust constructed values more
- bugs do not hide behind harmless-looking zero values

### Prevent Issue Loops By Closing Wrong Paths

This project should avoid the loop of:

- finding issues
- fixing issues
- the fix introducing more issues

The main way to reduce that is:

- fix domain shape before fixing caller behavior
- close misleading convenience APIs when they encode the wrong model
- add misuse-boundary tests, not only happy-path tests
- encode the violated rule in constructors, resolution logic, or query surface

In practice, a good issue fix usually includes:

- one constructor or API guard
- one data or resolved-metadata correction when needed
- one regression test
- one misuse-boundary test

### Domain Objects Should Be Small And Purposeful

Types should not own concerns that belong elsewhere.

For example:

- `HitPoints` owns totals, sources, temporary HP, nonlethal damage, and recalculation
- `HitPoints` does not own the entire combat-state model

That means concepts like:

- disabled
- dying
- unconscious
- dead state transitions

should live in a separate combat-state domain later.

This keeps the HP ledger useful without turning it into a god object.

### Preserve Exact Rules Math Internally

When Pathfinder uses fractional systems, this project tries to preserve the exact math internally and only round when the rules actually require it.

That is why `BaseAttackBonus` and `SavingThrow` use exact rational values instead of floats.

This avoids:

- float drift
- awkward equality checks
- hidden rounding errors in multiclass calculations

### Keep Later Expansion From Driving Current Scope

A recurring design goal is to keep the current core backlog from forcing rewrites later.

If a later class, creature type, or spellcasting source is added, the project should not require rewriting the core stat containers already built from the backlog.

That is why some types are intentionally generic.

Example:

- `CasterLevel` stores `Arcane`, `Divine`, and `Primal` source totals already present in the repo
- it does not hardcode class names into itself
- classes are expected to contribute to those totals elsewhere when that work is actually scheduled

This keeps existing domains stable without changing the current core-first delivery order.

### Keep Exact Value And Display Value When The Rules Need Both

Some Pathfinder stats have:

- a real internal value
- a displayed or applied value

That is why:

- `BaseAttackBonus` stores exact rational progression and rounded integer BAB
- `SavingThrow` stores exact rational progression and rounded integer save

This keeps the math honest without losing the number that actually matters at the table.

### Keep A Ledger When Provenance Matters

For some systems, just knowing the final total is not enough.

`HitPoints` is the clearest example.

It keeps sources like:

- `Base Dice`
- `Constitution`
- `Charisma (Undead)`
- `Construct Size Bonus`
- temporary HP sources

That makes recalculation and debugging much easier than storing only one final integer.

## Current Scope And Project-Specific Differences

The active backlog is core-only. The repository also contains some project-specific rules and conventions that are documented separately so they do not become the default direction for upcoming backlog work.

These are documented in:

[PF1 Differences](docs/pf1-differences/README.md)

Current repo-specific differences include:

- caster level tracked by source instead of by individual class
- `Titanic` as an officialized project creature size
- project-authored dual imperial/metric creature size measurements

Those differences exist in the repository, but they are not the source of truth for current delivery priority. `BACKLOG.md` and `AGENTS.md` are.

## Current Status Matrix

This table is the current foundation-audit snapshot for specialist review of the existing repo surface. It reflects the delivered core foundation through Feat, the completed P3 carryable equipment catalog path, the completed P4 feat selection context path, the completed P5 core ability composition path, the completed P6 class-level derived-fact path, and the current P7 character language composition gap.

| Area | Exists | Core-correct now | Intentional limit | Project-specific note |
| --- | --- | --- | --- | --- |
| `ability` | yes | carrying capacity, hit point averages, explicit HP base dice, core construct HP table, and core size ladder behavior are aligned for the current surface | not a full combat-state engine | `Titanic`, dual imperial/metric size measurements, source-based caster levels |
| `creaturetype` | yes | supported base types and the currently supported subtype effects resolve structurally | partial subtype coverage and partial trait model by design | none beyond project-specific sizes if a caller uses them |
| `character` | yes | composition helpers for race, class, HP, spellcasting progression, spell-list entries, feat prerequisites, selected carryable equipment, carried weight, racial ability composition, character level facts, base attack bonus facts, base saving throw facts, class HP ledger facts, and a refreshed class-level derived creation slice | not a full character aggregate yet | none |
| `language` | yes | core common language identities are seeded with secret-language metadata | identity/query surface only; no character language selection yet | none |
| `race` | yes | 7 core races seeded with lookup helpers and automatic/bonus language metadata | core seed/query slice only | none |
| `modifier` | yes | validated modifier chassis, stacking resolution, and usable target/condition refs | not full downstream composition yet | none |
| `skill` | yes | core catalog seeded with public lookup and grouped specializations accepted by the domain | not skill-rank composition yet | none |
| `class` | yes | 11 core classes seeded with validated chassis, lookup helper, and spellcasting metadata | class feature details and character composition are not modeled yet | none |
| `spellcasting progression` | yes | core progression tables seeded for bard, cleric, druid, paladin, ranger, sorcerer, and wizard | table/query surface only; no character spell-slot composition yet | none |
| `spell` | yes | spell chassis, spell-list entries, core spell-list bindings, and core spell data are seeded with read-only query helpers | no spell preparation, casting, or character spellbook composition yet | none |
| `feat` | yes | feat chassis, typed prerequisites, core feat seeds, read-only catalog helpers, and character prerequisite composition are present | no full feat allocation engine, feat-slot accounting, retraining, or bonus-feat grant flow yet | none |
| `equipment` | yes | equipment, weapon, armor/shield chassis, dual imperial/metric weight and length values, adventuring-gear seed/query helpers, carryable item boundary, character inventory composition, carried-weight composition, and bounded core weapon, armor, and shield seed/query slices are present | no combat behavior, equipped slots, ammunition tracking, AC composition, or magic items | metric values are stored alongside PF1 imperial values for equipment measurements |

## Current Domain Snapshot

### `AbilityScore`

Represents one ability score and its validity.

Important notes:

- score validity is explicit
- score value and score availability are separate concepts
- modifier calculation depends on a valid visible score

### `Alignment`

Represents the order and morality axes.

Important notes:

- construction is validated
- true neutral is rendered as `Neutral`, not `Neutral Neutral`

### `BaseAttackBonus`

Represents BAB as both exact fraction and rounded final value.

Important notes:

- class progressions are stored as rational values
- rounding only happens when deriving the displayed BAB

### `SavingThrow`

Represents an individual save with exact and rounded values.

Important notes:

- good and poor progressions are exact fractions
- the one-time `+2` good-save base bonus is tracked so multiclassing does not repeat it incorrectly

### `CasterLevel`

Represents the repository's current source-based caster level totals.

Important notes:

- sources are currently `Arcane`, `Divine`, and `Primal`
- a source can be impossible for a character instead of merely zero
- this domain stores totals, not class mapping rules

### `HitDie`

Represents the hit-die composition of a character or creature.

Important notes:

- semantic validity matters, not just field presence
- zero-total hit dice are rejected
- average HP uses the current core-aligned averages for this surface

### `HitPoints`

Represents a raw HP ledger.

Important notes:

- keeps a source breakdown
- supports temporary HP
- supports nonlethal damage
- recalculates when the relevant underlying stat changes
- does not own full combat-state semantics

### `Size`

Represents size-based rules data.

Important notes:

- includes attack and AC modifier
- includes CMB and CMD modifier
- includes Fly and Stealth modifiers
- includes construct HP bonus
- includes space and reach by body shape
- includes typical height and weight ranges with imperial and metric values stored together
- includes the repo-specific `Titanic` size entry

## Glossary

### Valid

When this project says a value is `valid`, it means the value exists and is currently usable under the rules for that domain.

Example:

- an `AbilityScoreValue` can store a number but still be invalid for use

### Invalid

An invalid value is one that should not be constructed or accepted by a validated setter.

Examples:

- negative hit dice
- unknown saving throw ids
- malformed size body shapes

### Impossible

`Impossible` is stronger than zero.

It means a stat or source does not exist for that character or cannot currently be used at all.

This is used in `CasterLevel`.

Example:

- `0, valid` means the source exists and is currently zero
- `0, invalid` means the source is impossible or unavailable

### Source

`Source` is the project term for a spellcasting origin bucket in `CasterLevel`.

Current sources:

- `Arcane`
- `Divine`
- `Primal`

This was chosen because it is more general and future-friendly than tying caster level directly to class names.

### Exact Value

The mathematically exact internal value before the game-facing rounding is applied.

Examples:

- BAB of `3/4`
- save total of `10/3`

### Display Value

The rounded or directly used table-facing number.

Examples:

- BAB `0` from exact `3/4`
- Fortitude `3` from exact `10/3`

### Ledger

A ledger is a domain structure that preserves the pieces that make up a total instead of storing only the final result.

This is especially useful when values must be recalculated after related stats change.

### Body Shape

The movement/reach profile used by `Size` when determining space and natural reach.

Current supported shapes:

- `Tall`
- `Long`

Invalid body shapes are rejected.

## Near-Term Next Steps

The next major tracked backlog area is:

- P7 core character languages

That work will likely interact with:

- selected core race facts
- automatic language metadata
- bonus language metadata
- composed Intelligence facts
- race-specific bonus language eligibility

The next executable backlog item is automatic racial languages into character language facts at the `character` boundary. Keep it adapter-focused: no bonus language selection, campaign language behavior, social checks, language-dependent spell effects, or non-core content.

## Philosophy In One Sentence

Model Pathfinder rules as small validated domains first, then compose them into bigger systems only when the boundaries are clear.
