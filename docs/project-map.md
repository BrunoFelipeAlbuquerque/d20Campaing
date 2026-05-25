# Project Map

Use this file for quick orientation before reading package internals.

`AGENTS.md` is still the workflow source of truth. `BACKLOG.md` and `ISSUES.md` still decide task order.

---

## Current Goal

Build a core-first Pathfinder 1e character domain from small validated pieces.

Order of work:

1. domain/chassis
2. core seed data
3. resolution/query logic
4. character composition

As of this map, the foundation domains, P4 feat selection composition adapters, P5 core ability composition adapters, P6 class-level derived-fact adapters and slice refresh, P7 core character language composition with slice coverage, the P8 core skill rank allocation path, and the P9 core skill check total path are present. Check `BACKLOG.md` for exact unchecked work.
The next planned area is P10 core feat allocation and grants; later non-core work is intentionally deferred.

---

## Fast Context Path

For most tasks, read:

1. `ISSUES.md`, unless explicitly skipped
2. the active `BACKLOG.md` item or user request
3. this map
4. directly relevant source files

For rule-sensitive behavior, search local rules:

```bash
rg -i "<rule term>" docs/pf1/chunks
```

Do not bulk-read local rule text, solution notes, or internal role files unless needed.

---

## Package Map

| Path | Responsibility | Boundary |
| --- | --- | --- |
| `internal/domain/rpg/character/ability` | Primitive values and math: ability scores, BAB, saves, HP, size, speed, alignment, caster level | Must not import higher domains |
| `internal/domain/rpg/character/creaturetype` | Structural creature rule resolution: type, subtype, traits, contextual flags | Structural only; not a full character engine |
| `internal/domain/rpg/character/language` | Core common language identities and secret-language flags | Language identity/query only; character language selection composes in `character` |
| `internal/domain/rpg/character/race` | Core race chassis, core race seeds, race lookup/query helpers | Race facts only; character choices compose elsewhere |
| `internal/domain/rpg/character/skill` | Core skill chassis, grouped skill parsing, skill catalog lookup, and skill ability-key metadata | Skill identity and metadata only; ranks and totals compose in `character` |
| `internal/domain/rpg/character/class` | Core class chassis, class seeds, spellcasting progression tables, class feature/proficiency IDs | Class metadata only; final character stats compose elsewhere |
| `internal/domain/rpg/character/spell` | Spell chassis, core spell data, spell-list entries, class spell-list queries | Spell/list metadata only; preparation/casting compose later |
| `internal/domain/rpg/character/feat` | Feat chassis, typed prerequisites, core feat seeds, feat catalog lookup | Feat facts only; selected character feats compose elsewhere |
| `internal/domain/rpg/character/equipment` | Core equipment chassis, weapon chassis, armor/shield chassis, core adventuring-gear seed batch, and carryable item lookup | Equipment facts only; selection and carried inventory compose in `character` |
| `internal/domain/rpg/character` | Character composition boundary and thin adapters across domains | Only place for cross-domain character composition |
| `internal/domain/rpg/modifier` | Modifier refs, sources, entries, and stacking/penalty resolution | Shared modifier logic; not character-specific by itself |
| `internal/text` | Generic text helpers | No RPG rules |

---

## Composition Surface

Current character-boundary adapters:

- `character_race.go`: selected core race lookup
- `character_class.go`: selected core class lookup
- `character_class_level.go`: selected class levels and character level facts
- `character_base_attack_bonus.go`: base attack bonus facts from selected class levels
- `character_base_saving_throw.go`: base saving throw facts from selected class levels
- `character_class_hit_points.go`: first-level class HP from selected class hit die
- `character_class_hit_point_ledger.go`: explicit class HP ledger facts beyond first level
- `character_spellcasting_progression.go`: class spellcasting progression access
- `character_spell_list_entry.go`: class spell-list entry access
- `character_feat.go`: feat prerequisite state, validated skill-rank allocation feed, and selected feat validation
- `character_equipment.go`: selected carryable item lookup with quantity
- `character_carried_weight.go`: carried equipment weight and load category from Strength
- `character_race_ability.go`: fixed and selectable core racial ability composition
- `character_language.go`: automatic and bonus racial language facts
- `character_skill_rank_allocation.go`: concrete skill rank allocation facts
- `character_skill_rank_budget.go`: skill-rank budget facts from class levels, Intelligence, and supported racial metadata
- `character_skill_rank_allocation_facts.go`: selected skill rank validation against budget and character-level caps
- `character_skill_check_total.go`: static skill check totals from ranks, ability modifiers, class-skill metadata, and trained-only availability
- `racial_hit_points.go`: creature rules to racial HP bridge

Keep composition thin. If logic belongs to a lower domain, add it there only when the backlog item requires it.

---

## Next Backlog Path

The next core-only backlog path is P10 core feat allocation and grants.

The P9 skill check total path is complete. Continue by creating explicit character feat slot and feat grant facts before composing feat allocation behavior.

Check `BACKLOG.md` before starting any far-future non-core item.

The path should prove general feat slots, supported racial and class bonus-feat grants, and selected-feat validation against existing prerequisite contexts without adding feat effects or combat bonus application.

Out of scope for this path:

- retraining, feat replacement, or optional non-core feat systems
- combat bonus application
- item creation, crafting, or spell effects from feats
- traits, archetypes, prestige classes, or non-core feats
- full mutable character aggregate
- broad folder or package reorganization

---

## Common Decisions

- Invalid states should fail construction.
- Zero-value misuse should fail to resolve.
- Query helpers should return defensive copies when data can be mutated.
- Core seed data should stay core-only unless backlog explicitly says otherwise.
- Unsupported prerequisite or selection shapes should fail closed, not pass with guessed behavior.
- Do not use PDFs, text chunks, or rule tooling from domain code.

---

## Token-Saving Rule

Prefer one focused file read over a broad scan.

If package purpose is unclear, read this map first. If rule behavior is unclear, search `docs/pf1/chunks` for the exact rule term. If task priority is unclear, read `ISSUES.md` and the next unchecked `BACKLOG.md` item.
