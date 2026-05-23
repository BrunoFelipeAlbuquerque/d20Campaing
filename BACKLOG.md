# BACKLOG

## Rules for backlog items

- each item must be implementable in one PR
- each item must have tests
- no item may require redesign unless explicitly stated
- each item must be explicitly one of:
  - domain/chassis
  - core seed data
  - resolution/query logic

---

## Established foundation already present in repo

These packages already exist and are part of the current project baseline. They are not pending backlog items in this file:

- `internal/domain/rpg/character/ability`
  - validated foundational stat and chassis domains already established
- `internal/domain/rpg/character/creaturetype`
  - structural creature rule resolution already established
  - intentionally partial to current project scope
- `internal/domain/rpg/character`
  - current composition boundary already established

The remaining planned backlog below starts at Race / Modifier / Skill / Class because those are the remaining source-of-truth delivery items after that existing foundation.

---

## Tooling and documentation

- [x] Formalize local Pathfinder rules lookup workflow
  - Add `docs/pf1/README.md`
  - Add `pf1-extract`, `pf1-chunk`, and `pf1-search` Makefile targets
  - Update `AGENTS.md` with Pathfinder rules lookup policy
  - Ensure rule-sensitive tasks require local `rg` lookup before implementation

---

## Completed foundation alignment

- [X] Align Size measurement storage with Equipment dual-unit measurements (domain/chassis):
  - length values store both feet and meters
  - weight values store both pounds and kilograms
  - existing imperial and metric getters remain available
  - tests cover stored dual-unit values

---

## P0 — Core foundation

### Race

- [X] Create the Race domain chassis:
  - RaceID
  - validated constructor
  - size
  - base speed
  - ability score modifiers
  - racial languages
  - racial features container (no APG traits)

- [X] Seed the 7 core races only:
  - dwarf
  - elf
  - gnome
  - half-elf
  - half-orc
  - halfling
  - human

- [X] Add race query helpers:
  - GetRaceByID
  - HasFeature
  - defensive-copy getters

- [X] Add race tests:
  - valid IDs
  - invalid construction
  - core lookup
  - feature presence

---

### Modifier

- [X] Create the Modifier domain chassis:
  - ModifierType
  - ModifierSource
  - Modifier entry with target and condition slots
  - circumstance source registry

- [X] Add modifier tests:
  - source normalization and validation
  - circumstance registry behavior
  - stacking and penalty resolution

---

### Skill

- [X] Create the Skill domain chassis:
  - SkillID
  - trained-only flag
  - armor-check-penalty flag

- [X] Model grouped skills:
  - Craft
  - Knowledge
  - Perform
  - Profession

- [X] Seed the core skill catalog only

- [X] Add skill tests:
  - valid IDs
  - grouped skill handling
  - invalid rejection

---

### Class

- [X] Create the Class domain chassis:
  - ClassID
  - hit die type
  - BAB progression
  - save progression metadata
  - skill ranks per level
  - class skills
  - weapon/armor proficiency metadata
  - spellcasting kind
  - key ability metadata

- [X] Seed the 11 core classes only:
  - barbarian
  - bard
  - cleric
  - druid
  - fighter
  - monk
  - paladin
  - ranger
  - rogue
  - sorcerer
  - wizard

- [X] Add class tests:
  - valid IDs
  - progression correctness
  - lookup correctness

---

### Spellcasting progression

- [X] Create core spellcasting progression tables:
  - class → spell slots per level

- [X] Seed progression for:
  - bard
  - cleric
  - druid
  - paladin
  - ranger
  - sorcerer
  - wizard

- [X] Add progression tests:
  - known breakpoints
  - caster vs non-caster validation

---

### Spell

- [X] Create the Spell domain chassis:
  - SpellID
  - school
  - descriptors
  - components
  - casting time
  - range
  - target/effect
  - duration
  - saving throw
  - spell resistance

- [X] Create Spell List Entry:
  - spell id
  - class id
  - spell level

- [X] Seed core spell list bindings:
  - class ↔ spell level mapping

- [X] Seed spell data (batch 1):
  - all cantrips/orisons

- [X] Seed spell data (batch 2):
  - levels 1–3

- [X] Seed spell data (batch 3):
  - levels 4–6

- [X] Seed spell data (batch 4):
  - levels 7–9

- [X] Add spell tests:
  - valid construction
  - class list lookup
  - invalid rejection

---

### Feat

- [X] Create the Feat domain chassis:
  - FeatID
  - category
  - prerequisite model
  - fighter bonus feat flag
  - metamagic flag
  - item creation flag

- [X] Seed core general feats

- [X] Seed core combat feats

- [X] Seed core critical feats

- [X] Seed core item creation feats

- [X] Seed core metamagic feats

- [X] Add feat tests:
  - prerequisite validation
  - category correctness
  - invalid rejection

---

## P1 — Core composition

### Minimum level-1 core character creation slice

Near-term goal: prove the existing core domains compose into one small, reviewable character creation path before adding more systems.

The slice should answer:

- can a core race and core class be selected and resolved through `character`?
- can HP use the existing race/class/ability foundations without a broad character aggregate?
- can caster spell slots resolve for a caster class?
- can chosen feats be accepted or rejected from existing prerequisite state?

Do not expand this slice into:

- equipment, wealth, or encumbrance
- skill-rank allocation rules
- spell preparation, spellbooks, or known-spell selection
- combat state
- a full mutable character aggregate unless the slice proves it is required
- non-core sources
- broad folder or package reorganization

- [X] Compose Race with character boundary (no redesign)

- [X] Compose Class with character boundary (no redesign)

- [X] Compose spellcasting progression with Class

- [X] Compose spell list entry with Class spellcasting

- [X] Compose feat prerequisites with:
  - ability scores
  - BAB
  - class features
  - skill ranks
  - other feats

- [X] Add minimum level-1 character creation slice tests (resolution/query logic):
  - race + class + HP through existing character adapters
  - caster spell slots through existing class spellcasting progression
  - feat prerequisite checks through existing feat prerequisite state
  - invalid selected inputs fail closed

---

## P2 — Core equipment and inventory

### Equipment foundation

Near-term goal: model core equipment facts in small PR-sized pieces before character inventory, carrying load, or combat-facing equipment expands.

The path should answer:

- can a core item be constructed with validated ID, cost, weight, and category?
- can a bounded core adventuring-gear batch be looked up defensively?
- can a character selection reference seeded equipment without owning equipment facts?

Do not expand this path into:

- magic items
- wealth generation, shops, or economy rules
- crafting
- combat resolution
- non-core equipment
- archetypes, prestige classes, or other later non-core work

- [X] Create the Equipment domain chassis (domain/chassis):
  - EquipmentID
  - display name
  - equipment category
  - cost
  - weight
  - validation
  - tests
  - no seed catalog, inventory, encumbrance, weapon/armor mechanics, or magic items

- [X] Seed core adventuring gear batch 1 (core seed data):
  - Backpack (empty)
  - Bedroll
  - Flint and steel
  - Pouch, belt (empty)
  - Rations, trail (per day)
  - Rope, hemp (50 ft.)
  - Torch
  - Waterskin
  - use Core Rulebook cost and weight
  - no special item rules beyond cost, weight, and category

- [X] Add equipment query helpers (resolution/query logic):
  - GetEquipmentByID
  - catalog enumeration
  - defensive-copy getters

- [X] Compose selected equipment with character boundary (resolution/query logic):
  - selected core equipment ID
  - quantity
  - seeded lookup
  - invalid or unknown equipment fails closed
  - no purchasing, containers, encumbrance, weapon/armor combat effects, or magic items

- [X] Compose carried weight with existing Strength carrying capacity (resolution/query logic):
  - inventory total weight
  - load category from total weight and existing Strength carrying-capacity math
  - invalid or unknown carried equipment fails closed
  - no armor penalties, speed recalculation, containers, or combat effects

- [X] Create the Weapon domain chassis (domain/chassis):
  - WeaponID
  - proficiency category
  - weapon category
  - damage profile
  - critical profile
  - range increment
  - cost
  - weight
  - validation
  - tests
  - no attack roll, damage roll, proficiency application, or combat resolution

- [X] Create the Armor and Shield domain chassis (domain/chassis):
  - ArmorID
  - armor or shield category
  - armor bonus
  - maximum Dexterity bonus
  - armor check penalty
  - arcane spell failure
  - speed impact metadata
  - cost
  - weight
  - validation
  - tests
  - no AC composition, proficiency application, or combat resolution

- [X] Add dual imperial and metric equipment measurements (domain/chassis):
  - equipment weight stores both ounces and kilograms
  - weapon range increments store both feet and meters
  - armor speed impacts store both feet and meters
  - constructors support imperial and metric entry points
  - existing PF1 imperial getters remain available
  - tests cover both unit systems and seeded storage
  - no combat behavior, encumbrance behavior changes, UI formatting, or new seed data

---

## P3 — Core carryable equipment catalogs

### Weapon, armor, and shield seed/query path

Near-term goal: expand the core carryable catalog in bounded batches now that inventory can address gear, weapons, armor, and shields through one carryable lookup.

The path should answer:

- can core weapon facts be seeded and queried without combat resolution?
- can core armor and shield facts be seeded and queried without AC or proficiency resolution?
- can seeded weapons, armor, and shields resolve through the carryable lookup used by character inventory?

Do not expand this path into:

- attack rolls or damage rolls
- AC composition
- weapon or armor proficiency application
- wielded slots, equipped slots, or container contents
- ammunition tracking
- masterwork items, special materials, or magic items
- wealth generation, shops, economy rules, or crafting
- non-core sources

- [X] Seed core simple weapons batch 1 (core seed data):
  - unarmed strike
  - gauntlet
  - dagger
  - light mace
  - sickle
  - club
  - heavy mace
  - morningstar
  - shortspear
  - longspear
  - quarterstaff
  - spear
  - crossbow, heavy
  - crossbow, light
  - dart
  - javelin
  - sling
  - use Core Rulebook cost, weight, damage, critical, range, category, and proficiency
  - no attack roll, damage roll, ammunition tracking, proficiency application, special materials, masterwork, or magic versions

- [X] Add weapon query helpers (resolution/query logic):
  - GetWeaponByID
  - catalog enumeration for seeded core weapons
  - defensive-copy behavior
  - unknown weapon IDs fail closed
  - no new weapon seed data

- [X] Resolve seeded weapons through carryable item lookup (resolution/query logic):
  - `GetCarryableItemByRef` resolves seeded `Weapon` refs
  - display name, cost, and weight come from the weapon seed
  - unseeded or malformed weapon refs fail closed
  - no attack roll, damage roll, ammunition tracking, wielded slots, or combat behavior

- [X] Seed core light armor and shields batch 1 (core seed data):
  - padded armor
  - leather armor
  - studded leather armor
  - chain shirt
  - buckler
  - light wooden shield
  - light steel shield
  - heavy wooden shield
  - heavy steel shield
  - tower shield
  - use Core Rulebook cost, weight, armor bonus, max Dexterity, armor check penalty, arcane spell failure, and speed impact metadata
  - no AC composition, proficiency application, armor speed recalculation, special materials, masterwork, or magic versions

- [X] Add armor and shield query helpers (resolution/query logic):
  - GetArmorByID
  - catalog enumeration for seeded core armor and shields
  - defensive-copy behavior
  - unknown armor IDs fail closed
  - no new armor or shield seed data

- [X] Resolve seeded armor and shields through carryable item lookup (resolution/query logic):
  - `GetCarryableItemByRef` resolves seeded `Armor` refs
  - display name, cost, and weight come from the armor or shield seed
  - unseeded or malformed armor refs fail closed
  - no AC composition, proficiency application, equipped slots, or combat behavior

- [X] Compose carried weapons, armor, and shields with character inventory (resolution/query logic):
  - selected weapon, armor, and shield carryable refs compose through `character`
  - carried weight uses the shared carryable lookup for seeded weapons, armor, and shields
  - invalid, malformed, or unknown carryable refs fail closed
  - no equipped slots, attack rolls, AC composition, proficiency application, or combat behavior

---

## P4 — Core feat selection context

### Selection-sensitive core feat prerequisites

Audit snapshot:

- `go test ./...` passes
- `ISSUES.md` has no open issue markers
- the remaining executable core gap is selection-sensitive feat prerequisite composition
- current `character` feat composition intentionally fails closed for:
  - selected weapon proficiency prerequisites
  - same-selection feat prerequisites
  - spell-school feat prerequisites
  - selected familiar eligibility prerequisites

Near-term goal: model the smallest explicit character selection facts needed for seeded core feats with selection-sensitive prerequisites.

The path should answer:

- can a selected core weapon satisfy selected-weapon prerequisites without combat behavior?
- can selected feat ownership preserve the selected weapon or spell school needed by same-selection prerequisites?
- can core spell-school feat chains compose without spell DC or casting behavior?
- can explicit familiar eligibility compose without familiar stat blocks?

Do not expand this path into:

- attack rolls or damage rolls
- feat bonus application
- weapon proficiency bonus application
- spell DCs, school powers, or spell preparation
- familiar stat blocks, animal companion rules, or companion advancement
- feat-slot accounting, retraining, or a full mutable character aggregate
- non-core feats or archetypes

- [X] Compose selected weapon proficiency feat prerequisites (resolution/query logic):
  - add explicit selected core weapon context to `CharacterFeatPrerequisiteState`
  - selected weapon IDs resolve through the seeded weapon catalog
  - selected weapon proficiency checks use existing core class weapon proficiency metadata
  - use an explicit `character` adapter between `equipment.WeaponID` / `WeaponProficiencyCategory` and `class.WeaponProficiencyID`
  - do not compare raw weapon IDs, display names, or proficiency labels as the rule boundary
  - cover category proficiency, individual proficiency, unknown selected weapons, malformed selected weapon facts, and unsupported mappings
  - unknown, malformed, or unsupported selected weapon facts fail closed
  - no attack rolls, damage rolls, feat bonuses, or combat behavior

- [X] Compose same-selection weapon feat prerequisites (resolution/query logic):
  - model selected feat ownership with a selected core weapon key
  - same-selection prerequisites compare the requested feat and selected weapon key
  - missing, duplicate, mismatched, or unknown selected feat facts fail closed
  - no feat-slot accounting, retraining, attack bonuses, or weapon bonus application

- [X] Compose spell-school feat prerequisites (resolution/query logic):
  - model selected feat ownership with a selected core spell school key
  - `Spell Focus` can satisfy seeded same-school prerequisites such as `Greater Spell Focus`
  - `Spell Focus (conjuration)` can satisfy `Augment Summoning`
  - invalid spell schools or mismatched selected schools fail closed
  - no spell DCs, prepared spells, spell slots, or school power behavior

- [X] Compose selected familiar eligibility prerequisites (resolution/query logic):
  - add explicit familiar eligibility context to `CharacterFeatPrerequisiteState`
  - `Improved Familiar` requires both seeded familiar access and familiar eligibility facts
  - missing or zero-value eligibility facts fail closed
  - no familiar stat blocks, familiar selection catalog, or companion advancement

- [X] Refresh minimum level-1 character creation slice for selected feat contexts (resolution/query logic):
  - keep the slice test-only and adapter-focused
  - include one accepted selected-weapon feat path
  - include one accepted spell-school feat path
  - include fail-closed mismatched selection coverage
  - no full character aggregate or non-core content

---

## P5 — Core character ability composition

### Race ability modifiers and selected ability facts

Near-term goal: compose explicit level-1 ability score facts with core race ability metadata before broader character aggregation, skill allocation, or combat-facing derived stats.

The path should answer:

- can fixed core racial ability modifiers compose onto explicit base ability scores?
- can human, half-elf, and half-orc selectable `+2` ability metadata choose exactly one core ability score?
- can composed ability scores feed existing character feat prerequisite and carrying-load adapters without duplicate side math?

Do not expand this path into:

- point buy, rolling methods, elite arrays, or NPC arrays
- aging effects, level-up ability increases, magic bonuses, or temporary ability damage
- skill-rank allocation, saving throw totals, attack rolls, armor class, or combat state
- feat-slot accounting, favored class bonuses, traits, archetypes, or non-core sources
- a full mutable character aggregate

- [X] Compose fixed racial ability modifiers with base ability scores (resolution/query logic):
  - accept explicit base ability score facts for all six core abilities
  - apply fixed core race ability modifiers from the selected race
  - reject missing, duplicate, malformed, or non-core ability score facts
  - reject selectable-modifier races in this fixed-only path
  - expose composed scores as existing `CharacterAbilityScore` facts
  - no point buy, rolling, age, level-up, magic, temporary ability damage, or full aggregate behavior

- [X] Compose selectable racial ability modifiers with base ability scores (resolution/query logic):
  - support human, half-elf, and half-orc selectable `+2` metadata
  - require exactly one selected core ability score for the selectable modifier
  - reject missing, duplicate, unknown, or non-core selected abilities
  - reject fixed-modifier races in this selectable-only path
  - expose composed scores as existing `CharacterAbilityScore` facts
  - no point buy, rolling, age, level-up, magic, temporary ability damage, or full aggregate behavior

- [X] Refresh minimum level-1 character creation slice for racial ability composition (resolution/query logic):
  - keep the slice test-only and adapter-focused
  - include one accepted fixed-race ability modifier path
  - include one accepted selectable-race ability modifier path
  - include fail-closed missing or mismatched selected ability coverage
  - demonstrate composed ability facts feeding an existing feat prerequisite or carrying-load adapter
  - no full character aggregate, skill allocation, combat behavior, or non-core content

---

## P6 — Core class-level totals and base derived facts

### Class levels, character level, BAB, saves, and HP

Independent audit snapshot:

- `go test ./...` passes
- P5 core character ability composition is complete
- `class` exposes BAB progression, save progressions, skill ranks per level, class skills, proficiencies, and spellcasting metadata
- `character` currently accepts class-level facts for feat prerequisites, but does not compose class levels into reusable character-level totals
- local PF1 rules identify character level, BAB, saving throws, feats, skill ranks, and combat statistics as separate character-building concerns

Near-term goal: make class levels produce reusable base character facts before skills, feat slots, equipment use, or combat-facing totals depend on ad-hoc inputs.

The path should answer:

- can selected core class levels produce character level and per-class level facts?
- can class levels produce base attack bonus using the project's fractional BAB variant?
- can class levels produce base saving throws using the project's fractional save variant?
- can class levels produce a class HP ledger from explicit per-level HP choices and Constitution facts?

Do not expand this path into:

- XP thresholds, advancement automation, retraining, or favored class bonuses
- feat-slot allocation
- skill-rank allocation
- armor class, attack rolls, damage rolls, or combat maneuvers
- spell preparation, known spells, or spell DCs
- a full mutable character aggregate
- non-core classes or sources

- [X] Compose selected class levels into character level facts (resolution/query logic):
  - accept explicit selected core class level facts
  - reject missing, duplicate, malformed, zero, or non-core class levels
  - expose total character level and per-class levels as reusable character facts
  - no XP thresholds, advancement automation, favored class bonuses, or full aggregate behavior

- [X] Compose class levels into base attack bonus facts (resolution/query logic):
  - use seeded core class BAB progressions
  - use the project's adopted fractional BAB math
  - expose exact and rounded BAB facts for downstream feat prerequisites and attack snapshots
  - reject unknown classes, invalid levels, or duplicate class entries
  - no attack rolls, iterative attack actions, weapon behavior, or combat state

- [X] Compose class levels into base saving throw facts (resolution/query logic):
  - use seeded core class save progressions
  - use the project's adopted fractional saving throw math
  - expose Fortitude, Reflex, and Will base save facts
  - reject unknown classes, invalid levels, or duplicate class entries
  - no ability modifiers, resistance bonuses, conditions, or combat state

- [X] Compose class hit point ledger beyond first level (resolution/query logic):
  - accept explicit HP choices or rolls for each class level after first
  - first class level still uses maximum class hit die plus Constitution modifier
  - later levels use explicit validated HP entries plus Constitution modifier
  - reject missing, duplicate, out-of-range, or non-core class HP entries
  - no random rolling, favored class bonuses, temporary HP, healing, or death-state behavior

- [X] Refresh minimum level-1 character creation slice for class-level derived facts (resolution/query logic):
  - keep the slice test-only and adapter-focused
  - demonstrate class levels feeding character level, BAB, and base saves
  - demonstrate class levels still feeding feat prerequisites without duplicate side inputs
  - include fail-closed duplicate or malformed class-level coverage
  - no full character aggregate, skills, feat slots, combat behavior, or non-core content

---

## P7 — Core character languages

### Automatic and bonus language composition

Near-term goal: compose race language metadata with Intelligence-derived bonus language capacity before broader character identity snapshots.

The path should answer:

- can automatic racial languages become character language facts?
- can bonus language choices respect race-specific lists and human/half-elf any-non-secret metadata?
- can Intelligence modifier limit bonus language count without allowing duplicates or unknown languages?

Do not expand this path into:

- literacy, secret languages beyond current core metadata, or campaign languages
- language-dependent spell effects
- skill checks, social interaction, or exploration behavior
- non-core races or languages
- a full mutable character aggregate

- [X] Compose automatic racial languages into character language facts (resolution/query logic):
  - use selected core race automatic language metadata
  - reject zero-value or unknown races
  - expose deduped character language facts
  - no bonus language selection or campaign language behavior

- [ ] Compose bonus language selections from race and Intelligence (resolution/query logic):
  - use selected race bonus-language metadata
  - derive maximum bonus language choices from composed Intelligence modifier
  - support human and half-elf any non-secret language metadata
  - reject missing, duplicate, unknown, over-budget, or disallowed choices
  - no secret languages, campaign-specific languages, or non-core content

- [ ] Refresh minimum level-1 character creation slice for languages (resolution/query logic):
  - keep the slice test-only and adapter-focused
  - include one fixed-list bonus-language race
  - include one any-non-secret bonus-language race
  - include fail-closed over-budget or disallowed language coverage
  - no full character aggregate or social/skill behavior

---

## P8 — Core skill rank allocation

### Rank budgets and selected ranks

Near-term goal: allocate skill ranks from class, level, Intelligence, and supported racial metadata before computing final skill check totals.

The path should answer:

- can class and Intelligence facts produce a skill-rank budget?
- can selected ranks validate against character level caps and grouped skill identities?
- can allocated ranks feed existing feat prerequisite checks without caller-provided side facts?

Do not expand this path into:

- final skill check totals
- armor check penalties
- aid another, taking 10/20, or skill-use DCs
- retraining or favored class bonuses
- non-core skills or traits
- a full mutable character aggregate

- [ ] Create character skill rank allocation facts (domain/chassis):
  - selected concrete skill identity
  - allocated ranks
  - validation for core skill and grouped skill identities
  - no final skill totals, armor penalties, or skill-use rules

- [ ] Compose skill-rank budget from class levels and Intelligence (resolution/query logic):
  - use seeded core class skill ranks per level
  - use composed Intelligence modifier
  - support currently modeled core racial rank metadata such as human `Skilled`
  - reject invalid class levels, invalid Intelligence facts, or negative budgets
  - no favored class bonuses, traits, or retroactive rebuild behavior

- [ ] Validate selected skill ranks against budget and level cap (resolution/query logic):
  - enforce total allocated ranks not exceeding budget
  - enforce rank cap from character level
  - reject duplicate, malformed, unknown, or over-cap skill allocations
  - support concrete grouped skills already accepted by the `skill` domain
  - no final check totals or armor penalties

- [ ] Feed allocated skill ranks into feat prerequisite state (resolution/query logic):
  - convert validated rank allocations into existing `CharacterSkillRanks` facts
  - demonstrate a seeded core feat with skill-rank prerequisites
  - reject caller-provided malformed rank facts through the allocation path
  - no feat allocation engine or skill check totals

---

## P9 — Core skill check totals

### Ability keys, class skill bonus, and armor penalties

Near-term goal: compute supported static skill check totals from allocated ranks and existing character facts.

The path should answer:

- can each core skill expose the ability score used for its check?
- can trained-only, class-skill, and rank metadata produce static skill totals?
- can armor check penalty apply only to applicable skills from equipped armor or shields?

Do not expand this path into:

- skill-use DCs, opposed checks, take 10/20, or aid another
- situational modifiers, conditions, spells, or magic items
- movement or combat action resolution
- non-core skills or traits

- [ ] Add core skill ability-key metadata (domain/chassis):
  - each seeded core skill exposes its ability score key
  - grouped skills share the correct family ability key
  - validation rejects unknown ability keys
  - no rank allocation or check totals in the skill domain

- [ ] Compose static skill check totals from ranks and ability scores (resolution/query logic):
  - use allocated ranks and composed ability facts
  - apply core class-skill bonus when at least one rank is present
  - enforce trained-only availability
  - reject missing ability facts, invalid ranks, or unknown skills
  - no armor penalties, situational modifiers, or roll behavior

- [ ] Compose armor check penalty into applicable skill totals (resolution/query logic):
  - use equipped armor and shield metadata once equipped inventory exists
  - apply penalties only to skills marked for armor check penalty
  - reject malformed equipped armor or shield facts
  - no combat actions, speed changes, or magic equipment

---

## P10 — Core feat allocation and grants

### Feat slots, bonus feats, and validated selections

Near-term goal: move from prerequisite checking to an explicit feat allocation surface.

The path should answer:

- can general feat slots be derived from character level?
- can modeled racial and class bonus-feat grants create typed feat slots?
- can selected feats validate prerequisites and required selections through existing character contexts?

Do not expand this path into:

- retraining, feat replacement, or optional non-core feat systems
- combat bonus application
- item creation, crafting, or spell effects from feats
- traits, archetypes, prestige classes, or non-core feats

- [ ] Create character feat slot and feat grant facts (domain/chassis):
  - general feat slots
  - racial bonus feat slots
  - class bonus feat slots with category constraints when supported
  - validation for slot identity and source
  - no allocation engine or feat effects

- [ ] Compose general feat slots from character level (resolution/query logic):
  - use character level facts from class-level composition
  - create core general feat slots at supported level thresholds
  - reject invalid, zero, or missing character level facts
  - no retraining, mythic feats, or non-core rules

- [ ] Compose supported racial and class bonus feat grants (resolution/query logic):
  - support human bonus feat metadata
  - support currently modeled class bonus-feat metadata such as fighter bonus feats
  - reject unsupported grant sources instead of accepting free-form slots
  - no archetypes, favored class bonuses, or non-core grants

- [ ] Validate selected feats against slots and prerequisites (resolution/query logic):
  - consume selected feat slot facts
  - use existing prerequisite and selected-context adapters
  - reject duplicate feats, missing slots, invalid selections, or unsatisfied prerequisites
  - no feat bonus application or combat behavior

---

## P11 — Core equipped inventory

### Equipped armor, shields, and wielded weapons

Near-term goal: distinguish carried items from equipped items before AC, speed, attack, or damage calculations.

The path should answer:

- can carried armor, shields, and weapons be equipped through validated character facts?
- can equipment slots prevent duplicate or impossible equipment states?
- can equipped facts expose the metadata needed by later defensive and attack snapshots?

Do not expand this path into:

- armor class totals
- attack rolls or damage rolls
- magic item body slots
- ammunition tracking
- containers, shops, or wealth generation
- non-core equipment

- [ ] Create equipped armor and shield facts (domain/chassis):
  - equipped armor ref
  - equipped shield ref
  - validation through the shared carryable lookup
  - reject unknown, malformed, duplicate, or non-armor refs
  - no AC totals, speed changes, or proficiency penalties

- [ ] Create wielded weapon facts (domain/chassis):
  - selected carried weapon ref
  - supported wield mode metadata where needed by core weapons
  - reject unknown, malformed, duplicate, or non-weapon refs
  - no attack rolls, damage rolls, ammunition tracking, or handedness edge cases beyond current seed needs

- [ ] Compose equipped items from carried inventory (resolution/query logic):
  - equipped refs must exist in selected carried inventory
  - reject equipping items not carried or carried in non-positive quantity
  - expose equipped armor, shield, and weapon facts for downstream snapshots
  - no combat behavior or magic item body slots

---

## P12 — Core defensive stats and movement

### AC, saves, initiative, and speed snapshots

Near-term goal: compose supported static defensive stats from existing ability, class, race, and equipped-item facts.

The path should answer:

- can AC use armor, shield, Dexterity, size, and maximum Dexterity metadata?
- can saving throw totals use base saves and ability modifiers?
- can initiative and speed use supported static facts without combat state?

Do not expand this path into:

- touch attacks, combat maneuvers, surprise, or condition effects beyond explicitly modeled static facts
- spell effects, magic items, or situational modifiers
- attack rolls, damage rolls, or combat actions
- non-core content

- [ ] Compose Armor Class snapshot from equipped armor and shield (resolution/query logic):
  - base 10 plus supported armor, shield, Dexterity, and size facts
  - enforce maximum Dexterity bonus from equipped armor
  - expose normal, touch, and flat-footed AC where supported by available facts
  - reject malformed equipped facts or missing ability/size facts
  - no combat maneuvers, conditions, or magic bonuses

- [ ] Compose saving throw totals from base saves and abilities (resolution/query logic):
  - use base saving throw facts from class-level composition
  - add Constitution, Dexterity, and Wisdom modifiers for Fortitude, Reflex, and Will
  - reject missing or invalid base save or ability facts
  - no resistance bonuses, conditions, spells, or magic items

- [ ] Compose initiative and static movement speed (resolution/query logic):
  - initiative uses Dexterity modifier and supported static facts
  - movement uses race base speed, armor speed impact, and carried-load category
  - reject missing speed, ability, armor, or load facts when required
  - no terrain, tactical movement, conditions, or combat state

---

## P13 — Core weapon attack and damage snapshots

### Supported static weapon math

Near-term goal: expose basic static attack and damage facts from BAB, ability, size, proficiency, and wielded weapon metadata.

The path should answer:

- can a wielded weapon produce a supported attack bonus snapshot?
- can weapon damage expose the correct core damage dice and ability modifier contribution?
- can unsupported weapon or wield cases fail closed instead of guessing?

Do not expand this path into:

- combat actions, hit resolution, iterative full attacks, critical confirmation, or damage application
- ammunition tracking
- two-weapon fighting
- special materials, masterwork, magic weapons, or situational modifiers
- non-core weapons

- [ ] Compose basic weapon attack bonus snapshot (resolution/query logic):
  - use BAB, ability modifier, size modifier, weapon proficiency, and weapon category metadata
  - support one wielded core weapon at a time
  - reject unsupported proficiency mappings, malformed wielded facts, or missing BAB/ability/size facts
  - no roll resolution, iterative attacks, or combat actions

- [ ] Compose basic weapon damage snapshot (resolution/query logic):
  - use seeded weapon damage profile for character size where supported
  - add supported ability modifier contribution for simple melee/ranged cases
  - reject unsupported double-weapon, ammunition, or special-case damage paths until scheduled
  - no damage application, critical hits, or resistance handling

- [ ] Refresh minimum character slice for equipped combat snapshots (resolution/query logic):
  - keep the slice test-only and adapter-focused
  - demonstrate equipped weapon facts feeding static attack and damage snapshots
  - include fail-closed unsupported or mismatched equipment coverage
  - no combat turn engine or roll resolution

---

## P14 — Core spellcasting capacity and caster facts

### Caster level, slots, and key ability constraints

Near-term goal: make class spellcasting metadata usable as character spellcasting capacity before modeling prepared or known spell selections.

The path should answer:

- can class levels produce caster-level facts for supported spellcasting sources?
- can spell slots combine class progression and key ability bonus-spell rules?
- can spell availability reject unsupported class, level, or ability combinations?

Do not expand this path into:

- prepared spells, spellbooks, spells known, domains, schools, or bloodlines
- concentration checks, spell DCs, casting actions, or spell effects
- magic items or non-core classes

- [ ] Compose caster level from selected core class levels (resolution/query logic):
  - map supported core spellcasting classes into existing caster source facts
  - reject non-spellcasting classes and unsupported source mappings where appropriate
  - expose caster-level facts for feat prerequisites and spellcasting snapshots
  - no spell preparation, spell effects, or multiclass caster-level exceptions beyond current scope

- [ ] Compose spell slots with key ability support (resolution/query logic):
  - use seeded class spellcasting progression tables
  - use composed key ability score facts for bonus-spell eligibility
  - distinguish unavailable spell levels from zero-slot unlocked levels
  - reject invalid class levels, missing key ability facts, or unsupported spellcasting profiles
  - no prepared spells, spells known, domains, schools, or spell DCs

- [ ] Refresh minimum character slice for spellcasting capacity (resolution/query logic):
  - keep the slice test-only and adapter-focused
  - demonstrate caster level and spell-slot facts for one prepared caster and one spontaneous caster
  - include fail-closed non-caster or low-key-ability coverage
  - no spell selection, spellbook, or casting behavior

---

## P15 — Core spell selection and spell DC snapshots

### Prepared spells, known spells, spellbooks, and DCs

Near-term goal: model character spell choices separately from spell effects.

The path should answer:

- can prepared casters select prepared spells from valid class spell lists?
- can spontaneous casters select known spells from valid class spell lists?
- can wizard spellbook ownership and spell DC snapshots be represented without casting behavior?

Do not expand this path into:

- spell effects, targeting, durations in play, concentration checks, or combat casting
- domains, bloodlines, arcane schools, or familiars beyond already modeled eligibility facts
- scrolls, wands, potions, or magic item casting
- non-core spells or classes

- [ ] Create prepared spell selection facts (domain/chassis):
  - class spell list entry
  - prepared spell level and slot identity
  - validation for seeded core spell/list bindings
  - no spell effects, casting action, or prepared-slot consumption

- [ ] Create spontaneous known spell selection facts (domain/chassis):
  - class spell list entry
  - known spell level
  - validation for seeded core spell/list bindings
  - no spontaneous casting action or daily resource spending

- [ ] Create wizard spellbook ownership facts (domain/chassis):
  - owned wizard spell list entries
  - validation against seeded wizard spell bindings
  - reject unknown or non-wizard spell entries
  - no scroll copying costs, spell research, or spell preparation behavior

- [ ] Compose spell DC snapshots from selected spells and key ability (resolution/query logic):
  - DC uses spell level and key ability modifier
  - reject missing key ability facts or invalid selected spell facts
  - no school powers, metamagic adjustment, spell effects, or saving throw resolution

---

## P16 — Core character sheet snapshot

### Immutable aggregate for supported core facts

Near-term goal: create the first broad character snapshot only after the supporting composition paths exist.

The path should answer:

- can supported race, class, ability, language, skill, feat, equipment, defensive, attack, and spell facts be assembled once?
- can the snapshot expose read-only facts without moving domain logic out of lower packages?
- can invalid or unsupported partial states fail closed at the aggregate boundary?

Do not expand this path into:

- mutable editing workflows
- persistence, serialization, UI, imports, or exports
- non-core sources
- combat turn state or action resolution
- campaign-specific house rules beyond already documented project differences

- [ ] Create immutable core character sheet snapshot (domain/chassis):
  - selected race
  - selected class levels
  - composed ability facts
  - composed languages
  - allocated skills
  - allocated feats
  - selected carried and equipped items
  - supported derived defensive, attack, and spellcasting facts
  - no persistence, UI, mutable aggregate behavior, or non-core sources

- [ ] Add supported snapshot query helpers (resolution/query logic):
  - expose read-only facts grouped by character-sheet area
  - return defensive copies for slice-backed data
  - reject unsupported facts instead of manufacturing defaults
  - no recalculation side effects or editing workflow

- [ ] Add end-to-end core level-1 character sheet slice (resolution/query logic):
  - build one complete supported level-1 core character from existing domains
  - demonstrate race, class, abilities, languages, skills, feats, equipment, derived stats, and spellcasting capacity where supported
  - include fail-closed missing required area coverage
  - no non-core content, combat turn engine, or persistence

---

## P17 — Core advancement and resource state

### Later core-only lifecycle work

Near-term goal: capture level advancement and mutable daily/resource state after the immutable sheet snapshot is trustworthy.

The path should answer:

- can a supported character advance a level without bypassing validation?
- can daily resources and damage state be represented without corrupting immutable build facts?
- can later play-state systems depend on the character sheet snapshot without owning build rules?

Do not expand this path into:

- non-core classes, archetypes, prestige classes, or mythic rules
- campaign persistence, UI, imports, or exports
- full combat automation
- magic item economy or crafting systems

- [ ] Compose supported level-up transition inputs (resolution/query logic):
  - accept next class level selection
  - require updated HP, skill, feat, and spell choices when applicable
  - reject incomplete or unsupported advancement states
  - no XP table automation or retraining

- [ ] Create daily resource state facts for supported class and spellcasting resources (domain/chassis):
  - spell slots used
  - supported per-day class resources where metadata exists
  - validation for non-negative and in-cap usage
  - no class feature effects or combat action engine

- [ ] Separate damage/play-state facts from build snapshot facts (domain/chassis):
  - current HP and nonlethal damage use existing HP ledger behavior
  - status conditions remain explicit play-state facts
  - reject states that contradict the immutable character sheet snapshot
  - no full combat automation or condition effect engine

---

## P18 — Core favored class and level-choice ledgers

### Character creation choices that affect later totals

Near-term goal: represent per-character and per-level choices that PF1 expects before class features and full snapshots can be trusted.

The path should answer:

- can favored class selection be represented separately from class identity?
- can per-level favored class bonuses be tracked without changing class seed data?
- can later skill/HP/sheet totals consume level-choice ledgers without guessing?

Do not expand this path into:

- alternate favored class bonuses from non-core books
- retraining, rebuild rules, or campaign-specific house rules
- full level-up workflow UI or persistence
- archetypes, traits, or prestige classes

- [ ] Create favored class selection facts (domain/chassis):
  - selected core class reference
  - validation against existing selected class levels
  - no alternate racial favored class bonuses or non-core classes

- [ ] Create favored class bonus ledger facts (domain/chassis):
  - per-level HP or skill-rank bonus choice
  - validation that bonuses apply only to favored class levels
  - reject duplicate or out-of-range level entries
  - no retraining or alternate favored class bonus tables

- [ ] Compose favored class bonuses into supported HP and skill budgets (resolution/query logic):
  - HP bonus feeds the HP ledger path
  - skill bonus feeds skill-rank budget
  - reject unresolved favored class or level-choice state
  - no aggregate character editor or persistence

---

## P19 — Core racial feature effects

### Typed composition for already seeded racial features

Near-term goal: turn supported core racial feature metadata into explicit character facts only where the repo has enough data to do so correctly.

The path should answer:

- can seeded racial feature IDs produce typed character facts instead of loose strings?
- can static racial skill, save, senses, and weapon familiarity facts be exposed without a general effects engine?
- can unsupported racial features fail closed instead of silently disappearing?

Do not expand this path into:

- spell-like abilities, natural attacks, or combat automation
- non-core races or alternate racial traits
- broad modifier engine redesign
- creaturetype ownership of character-specific choices

- [ ] Create racial senses and movement facts for supported core race features (resolution/query logic):
  - low-light vision and darkvision facts where seeded
  - base land speed facts where seeded
  - validation against known race feature IDs
  - no perception simulation or lighting engine

- [ ] Create racial skill and save modifier facts (resolution/query logic):
  - supported static skill bonuses
  - supported static saving throw bonuses
  - source labels that preserve race and feature identity
  - no contextual adjudication beyond explicit seeded metadata

- [ ] Create racial weapon familiarity facts (resolution/query logic):
  - supported weapon familiarity groups
  - validation against seeded core weapon metadata where available
  - no attack roll or proficiency expansion outside supported facts

---

## P20 — Core class feature chassis

### Class feature grants before feature effects

Near-term goal: make class feature availability explicit by class level before modeling individual feature resources or effects.

The path should answer:

- can a class level produce a validated list of core class feature grants?
- can supported feature IDs be queried by character level without encoding effects yet?
- can unsupported feature effects remain visible as unsupported rather than guessed?

Do not expand this path into:

- full class feature behavior
- archetypes, alternate class features, prestige classes, or traits
- non-core class options
- combat actions, spell effects, or UI workflows

- [ ] Add class feature grant progression metadata for core classes (core seed data):
  - feature IDs by class and level
  - validation that IDs are non-empty and class-scoped
  - no effect behavior, choices, or non-core features

- [ ] Compose class feature grants from selected class levels (resolution/query logic):
  - aggregate feature grants by selected class levels
  - expose class, level, and feature identity
  - reject unknown class feature IDs
  - no resource pools or effect resolution

- [ ] Add unsupported-feature visibility for the character sheet path (resolution/query logic):
  - distinguish supported facts from known-but-unimplemented feature grants
  - keep unsupported feature IDs inspectable
  - no fabricated defaults for missing behavior

---

## P21 — Core class feature choices and resources

### Validated choices that many core classes require

Near-term goal: capture the first tier of class-specific choices and daily resources without building a full effects engine.

The path should answer:

- can class features that require a choice store that choice with validation?
- can per-day or per-round resource pools be represented from existing class levels?
- can later snapshots show class choices without owning class-specific behavior?

Do not expand this path into:

- combat action resolution
- spell effects, summoned creatures, or domain power automation
- archetypes or non-core options
- full companion stat advancement

- [ ] Create domain, school, bloodline, and favored enemy choice facts (domain/chassis):
  - validate choice identity against core class feature eligibility
  - preserve selected class and level source
  - no power effects or spell-list expansion beyond explicit later items

- [ ] Create rage, bardic performance, channel energy, and smite resource facts (domain/chassis):
  - resource maximum from class level where core rules define it
  - current usage validation for non-negative in-cap values
  - no action timing, targets, or combat effect automation

- [ ] Compose supported class feature choices into the character sheet snapshot (resolution/query logic):
  - expose selected feature choices and resource caps
  - reject missing required choices for supported classes
  - no broad class-feature effect engine

---

## P22 — Core companions and familiars

### Companion ownership before full play automation

Near-term goal: model core animal companion, familiar, and paladin mount ownership enough for a character sheet to represent the relationship correctly.

The path should answer:

- can a class feature require and validate a companion or familiar selection?
- can supported companion/familiar base choices be queried from core data?
- can companion snapshots remain separate from the owning character snapshot?

Do not expand this path into:

- eidolons, cohorts, leadership, or non-core companions
- animal trick training workflows
- full combat AI or action automation
- magic item sharing or campaign persistence

- [ ] Seed core familiar and animal companion base choices (core seed data):
  - core legal choices only
  - identity, size, movement, and base stat references needed for later snapshots
  - no non-core companions or template variants

- [ ] Create companion and familiar selection facts (domain/chassis):
  - owner class feature reference
  - selected companion/familiar reference
  - validation against eligibility
  - no advancement calculations yet

- [ ] Create minimal companion/familiar sheet snapshots (resolution/query logic):
  - base choice facts
  - owner level reference
  - known unsupported advancement areas reported explicitly
  - no combat turn automation or persistence

---

## P23 — Core wealth and purchase ledgers

### Money and ownership before broader equipment use

Near-term goal: represent starting wealth, currency, purchases, and equipment ownership as validated ledgers rather than loose inventory assumptions.

The path should answer:

- can a core character track currency in PF denominations?
- can starting wealth by class be represented without a marketplace engine?
- can equipment ownership and carried/equipped state validate against available purchased items?

Do not expand this path into:

- magic item economy, crafting, selling, settlement markets, or downtime
- randomized starting wealth rolls unless explicitly added by a later item
- non-core equipment
- encumbrance redesign

- [ ] Create currency and price value objects (domain/chassis):
  - copper, silver, gold, and platinum denominations
  - normalized comparison and addition behavior
  - reject negative money values

- [ ] Seed core class starting wealth metadata (core seed data):
  - average starting gold by core class
  - validation against existing class IDs
  - no randomized roll workflow

- [ ] Create purchase and ownership ledger facts (domain/chassis):
  - item reference, quantity, and price paid
  - total spend validation against available currency
  - carried/equipped state must reference owned items
  - no shop, sale, crafting, or magic item support

---

## P24 — Expanded core mundane equipment

### Fill the mundane catalog after ownership exists

Near-term goal: broaden core equipment coverage in bounded batches so character inventory, proficiencies, AC, and attacks can use real core equipment data.

The path should answer:

- can the core mundane equipment catalog cover common character-building choices?
- can weapon, armor, shield, ammunition, and adventuring gear seeds stay bounded and reviewable?
- can query helpers remain read-only and fail closed for unknown equipment IDs?

Do not expand this path into:

- magic items, alchemical item effects, vehicles, siege engines, or services automation
- non-core equipment
- combat action resolution
- economy simulation

- [ ] Add remaining simple, martial, and exotic core weapon seed batches (core seed data):
  - weapon group, proficiency category, damage, critical, range, weight, and length where available
  - dual imperial/metric values following the equipment gold standard
  - no attack behavior or magic weapons

- [ ] Add remaining core armor, shield, and ammunition seed batches (core seed data):
  - armor bonus, shield bonus, max Dex, armor check penalty, speed impact, weight, and cost
  - dual imperial/metric values where physical measurements exist
  - no special material, enhancement, or magic behavior

- [ ] Add bounded adventuring gear and tool seed batches (core seed data):
  - only items needed for ordinary core character sheets first
  - carryable item bridge coverage
  - no item effects unless explicitly represented as unsupported metadata

---

## P25 — Character modifier integration

### Shared modifier resolution for character snapshots

Near-term goal: connect existing modifier primitives to supported character facts so static bonuses are not hand-summed in every adapter.

The path should answer:

- can race, class, feat, equipment, and condition bonuses share one static modifier path?
- can stacking rules be applied at the character snapshot boundary?
- can unsupported contexts remain explicit rather than silently ignored?

Do not expand this path into:

- a general rules scripting engine
- dynamic encounter state or turn automation
- spell effect resolution
- non-core content

- [ ] Map supported static racial and equipment bonuses into modifier entries (resolution/query logic):
  - preserve source identity and bonus type
  - validate modifier target refs
  - no contextual or temporary effects

- [ ] Map supported feat and class-feature static bonuses into modifier entries (resolution/query logic):
  - only seeded feats/features with explicit supported behavior
  - fail closed for unsupported bonus shapes
  - no broad effect engine

- [ ] Resolve static character modifier totals for sheet snapshots (resolution/query logic):
  - apply stacking and penalty behavior through `modifier`
  - expose included and excluded sources for debugging
  - no combat action or spell-effect automation

---

## P26 — Core conditions and play-state modifiers

### Explicit temporary state before combat automation

Near-term goal: represent common PF1 character conditions as validated facts with limited static modifier output.

The path should answer:

- can a character snapshot accept explicit play-state conditions?
- can common static penalties and bonuses flow through modifier resolution?
- can unsupported condition behavior remain visible without being approximated?

Do not expand this path into:

- initiative order, rounds, durations, or turn engine
- full condition interaction automation
- spell or attack resolution
- campaign persistence

- [ ] Create core condition fact values (domain/chassis):
  - condition identity and optional source label
  - validation for supported core condition IDs
  - no duration or encounter timeline

- [ ] Add static condition modifier mappings for common conditions (resolution/query logic):
  - supported static AC, attack, save, skill, and ability-check effects
  - explicit unsupported notes for complex conditions
  - no movement/action denial automation

- [ ] Compose condition play-state into character snapshots (resolution/query logic):
  - apply supported modifier entries
  - keep build facts separate from play-state facts
  - reject unknown condition IDs

---

## P27 — Core combat stat snapshots

### Static combat-facing totals without combat resolution

Near-term goal: expose the combat numbers a PF character sheet needs after base facts, equipment, modifiers, and conditions exist.

The path should answer:

- can AC variants, initiative, CMB, CMD, and attack summaries be produced from validated facts?
- can equipped weapon and armor state influence static combat snapshots?
- can the snapshot stop before rolling attacks or resolving actions?

Do not expand this path into:

- combat turn engine, targeting, hit confirmation, or damage application
- spell attacks or spell effects
- monster/NPC automation
- non-core combat options

- [ ] Compose AC, touch AC, and flat-footed AC snapshots (resolution/query logic):
  - armor, shield, Dex, size, natural armor where supported, and modifier entries
  - condition-aware static variants where supported
  - no attack resolution or cover/concealment adjudication

- [ ] Compose CMB, CMD, and initiative snapshots (resolution/query logic):
  - BAB, ability modifiers, size, feats, equipment, and supported modifiers
  - validate missing required facts
  - no maneuver action resolution

- [ ] Compose iterative weapon attack and damage summaries (resolution/query logic):
  - BAB-derived iterative attacks
  - selected wielded weapon, ability modifier, proficiency, and size support
  - no hit rolls, critical confirmation, ammunition depletion, or damage application

---

## P28 — Core spellcasting choices and preparation lifecycle

### Spellcaster-specific choices after basic capacity exists

Near-term goal: represent the class choices and preparation state that make core prepared and spontaneous casters usable on a sheet.

The path should answer:

- can domain, school, bloodline, and spellbook choices affect supported spell availability?
- can prepared slots and spontaneous known spells validate against class-specific limits?
- can spent spell slots remain play-state separate from build facts?

Do not expand this path into:

- spell effects, targeting, concentration, counterspelling, or metamagic behavior
- magic items, scroll copying workflow, or spell research
- non-core spells, domains, bloodlines, or schools
- encounter automation

- [ ] Compose domain, school, bloodline, and opposition-school spell availability facts (resolution/query logic):
  - use existing class feature choices
  - validate against core spell/list data where seeded
  - no power effects or non-core options

- [ ] Validate prepared spell and spontaneous known spell limits by class level (resolution/query logic):
  - prepared slots by spell level
  - spontaneous known spell counts where core data supports them
  - reject missing key ability or invalid list entries
  - no casting action behavior

- [ ] Create spell preparation and slot-use play-state facts (domain/chassis):
  - prepared spell identity by slot
  - used slot counts or used prepared slots
  - separate build choices from daily play state
  - no spell effects or rest automation yet

---

## P29 — Core rest and reset lifecycle

### Daily reset boundaries

Near-term goal: define the smallest lifecycle rules needed to reset supported daily resources without turning the project into a campaign manager.

The path should answer:

- can daily resources reset from a validated character snapshot?
- can prepared spell state and class resources be refreshed consistently?
- can mutable play-state be updated without mutating build facts?

Do not expand this path into:

- calendar, travel, downtime, or campaign persistence
- encounter timelines or duration tracking
- healing services, crafting, or spell effect automation
- non-core resource systems

- [ ] Create daily reset input facts (domain/chassis):
  - reset source snapshot
  - selected prepared-spell changes where applicable
  - validation that reset applies only to supported resources
  - no calendar or rest-interruption rules

- [ ] Reset supported class and spell resources (resolution/query logic):
  - clear used daily class resources
  - refresh prepared/used spell slot state
  - preserve build choices and unsupported resources explicitly
  - no combat or duration automation

- [ ] Add lifecycle regression slice for one martial and one caster character (resolution/query logic):
  - prove build facts remain immutable
  - prove mutable state resets only supported resources
  - no UI, persistence, or non-core content

---

## P30 — Core audit and source-of-truth refresh gate

### Refresh before any non-core expansion

Near-term goal: force a complete source-of-truth pass after the core character runway has enough implemented surface to make a new audit useful.

The path should answer:

- do `BACKLOG.md`, `README.md`, `docs/project-map.md`, code, tests, and local PF1 rules still agree?
- are missing core PF character areas captured before non-core work resumes?
- are unsupported behaviors visible as explicit gaps instead of accidental omissions?

Do not expand this path into:

- implementing new behavior during the audit item
- reordering existing backlog work without explicit user direction
- non-core sources
- broad architecture redesign

- [ ] Run independent core character audit using `docs/independent-audit-protocol.md` (resolution/query logic):
  - compare implemented behavior against current backlog and project map
  - verify local PF1 rule references for implemented rule-sensitive behavior
  - record only actionable gaps as `ISSUES.md` or new backlog items when explicitly requested
  - no unrelated refactors or feature work

- [ ] Refresh README current status and project map after the audit (resolution/query logic):
  - update current implemented surface
  - update next backlog path
  - remove stale phase references
  - no source reorganization

---

## P99 — Far future (non-core)

- [ ] Archetype / Alternate Class Feature

- [ ] Prestige Classes

- [ ] Non-core races

- [ ] Non-core feats

- [ ] Non-core spells
