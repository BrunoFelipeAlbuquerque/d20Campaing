# ISSUES

## Rules for issue ordering

- issues are grouped by criticality
- `NEED` = blocking or correctness-critical
- `SHOULD` = important soon
- `CAN` = useful later

---

## NEED

- [X] Compose caster-level feat prerequisites before selected core feats silently fail:
  - `CasterLevelPrerequisite` is a validated prerequisite shape in the Feat domain
  - core item creation feats such as `Scribe Scroll`, `Brew Potion`, `Craft Wand`, and `Forge Ring` are seeded with caster-level prerequisites
  - core combat feats such as `Arcane Armor Training` and `Arcane Armor Mastery` are also seeded with caster-level prerequisites
  - `CharacterFeatPrerequisiteState` currently accepts ability scores, BAB, class levels, class features, skill ranks, and feats, but no caster-level fact
  - `CharacterFeatPrerequisiteState.SatisfiesPrerequisite` does not handle `CasterLevelPrerequisite`
  - valid core feats with caster-level prerequisites therefore cannot compose through `character`
  - fix the composition surface before adding another domain or expanding feat selection

- [X] Clarify fractional multiclass saving-throw math before Class composition builds on the wrong assumption:
  - the project uses Pathfinder Unchained fractional base bonuses for core class BAB and base save progression math
  - added local rule data under `docs/pf1/text/PF_Unchained_Fractional_Base_Bonuses.txt`
  - added a matching searchable chunk under `docs/pf1/chunks/PF_Unchained_Fractional_Base_Bonuses.jsonl`
  - under that variant, each good-save class contributes `1/2` per level and the `+2` good-save base bonus is added once per save type
  - `SavingThrow.AddClassLevel` and `TestSavingThrowAddClassLevel_DoesNotRepeatGoodSaveBaseBonusAcrossMulticlassing` therefore match the adopted variant
  - restoring Core Rulebook table-stacking behavior would be incorrect while this fractional variant remains the project rule

- [X] Define the Feat prerequisite model boundary before the Feat chassis locks in an underpowered shape:
  - the next backlog item asks for a `prerequisite model`, not just feat IDs and category flags
  - core feat prerequisites use several distinct inputs: ability scores, BAB, skill ranks, spellcasting/caster level, class level, class-specific access such as fighter bonus feats, and other feats
  - current `class` metadata exposes class identity, progressions, skills, proficiencies, and spellcasting kind, but not class-level prerequisite terms or class feature IDs
  - current `race` metadata represents human `Bonus Feat` and similar feature facts only as racial feature labels, not as typed feat-selection grants
  - if the Feat chassis models prerequisites as free-form strings or only a single flat ref type, later composition will need side tables or a redesign when `Compose feat prerequisites` arrives
  - before creating `feat`, decide the smallest core prerequisite value objects needed now and reject unsupported prerequisite shapes instead of accepting opaque text

- [X] Stop collapsing delayed-caster zero-slot unlock rows before Spell work builds on the wrong spell-availability model:
  - `normalizeSpellSlotsByClassLevel` trims trailing zeroes and returns `nil` for all-zero rows, so the domain cannot represent core `0 spells per day` breakpoints for newly unlocked paladin/ranger spell levels
  - `delayedFourthLevelCasterSpellSlotsByClassLevel` currently works around that gap by seeding level 4 as `{0, 1}`, which makes paladin/ranger level 4 resolve to one 1st-level spell slot instead of core bonus-only access
  - later delayed-caster breakpoints like paladin/ranger 7th, 10th, and 13th still collapse into the same `GetSpellSlots(...)=0` result as a spell level that is not unlocked yet, so callers cannot tell `bonus spells only` from `cannot cast this spell level`
  - current tests lock in that wrong surface by asserting paladin level 4 has one 1st-level slot and by treating absent spell levels as successful zero-slot lookups
  - the next Spell backlog work will otherwise miscompute delayed-caster spell access and bonus-spell eligibility

- [X] Validate exported modifier target and condition refs before `Modifier` construction accepts them:
  - `NewTargetRef` and `NewConditionRef` now exist as validated entry points
  - `NewModifier` still accepts zero-valued exported `TargetRef` and `ConditionRef` without checking them
  - this reopens the invalid-state path that the modifier chassis was supposed to close
  - future callers can still construct semantically broken modifiers even after following the new public surface

- [X] Split race language metadata before character language composition hardens around the wrong shape:
  - `Race` currently stores a single `racialLanguages []LanguageID` list
  - the seeds only model automatic starting languages and cannot represent core bonus-language choices
  - core races like dwarf, elf, gnome, half-elf, half-orc, halfling, and human therefore lose rule-important language selection metadata
  - future character language handling would otherwise need side tables or overloading of the current field

- [X] Stop treating grouped skills as family markers only before later skill modeling builds on the wrong unit:
  - `Skill` only records `grouped bool` for `Craft`, `Knowledge`, `Perform`, and `Profession`
  - there is no way to represent a concrete grouped skill entry such as `Knowledge (arcana)` or `Perform (sing)`
  - current tests explicitly reject specialized grouped entries as invalid IDs
  - later class-skill and character skill-rank work will otherwise have to either rank whole families or bolt on a second skill representation

- [X] Expose usable modifier target and condition value objects before consumers route around the domain:
  - `Target` and `Condition` are sealed by unexported marker methods
  - the package exports no concrete target or condition types and no constructors for them
  - outside `modifier` code can currently only pass `nil` into those slots
  - this will push future composition work toward ad-hoc side metadata instead of the declared modifier chassis

- [X] Remove the direct humanoid racial HD convenience path before Class work normalizes the wrong humanoid model:
  - `ResolvedCreatureRules.NewRacialHitDie` still constructs humanoid `d8` racial hit dice
  - `resolver_test.go` still locks in that humanoid construction succeeds while the class-rule flag is present
  - `character.NewRacialHitPoints` guards one entry point, but the lower-level `creaturetype` API still advertises the wrong convenience path
  - the next Class backlog work can easily model humanoids as having racial HD if this surface stays open

- [X] Restore core Strength carrying-capacity math before encumbrance depends on it:
  - `AbilityScore.GetCarryingCapacity` is currently backed by a custom metric table
  - the returned pound values are wrong for core PF1 load bands
  - example: `Strength 18` resolves to about `110.23 lb` light load instead of core `100 lb`
  - current tests lock in the non-core numbers

- [X] Stop shipping project HP house rules as default HP math:
  - `HitDie.GetAverageBaseHP` uses fixed averages (`d6=4`, `d8=5`, `d10=6`, `d12=7`)
  - construct HP uses the custom Large+ / `Titanic` size bonus table
  - `creaturetype` and `character` helpers already depend on those values
  - current HP totals are therefore not trustworthy as core PF1 outputs

- [X] Correct outsider breathing semantics in `creaturetype`:
  - `OutsiderType` currently gets `NoNeedToEatSleepBreathe`
  - standard PF1 outsiders still breathe
  - `NativeSubtype` currently compensates by adding `BreatheEatSleep`, which means the base outsider metadata is wrong

- [X] Validate modifier entries before they become shared composition inputs:
  - `Modifier` is a free-form exported struct with no validated constructor
  - unknown or mistyped `ModifierType` values silently create new stacking buckets
  - empty type has special resolution behavior instead of failing fast
  - this breaks the repo rule that invalid states should fail construction

- [X] Restore backlog alignment before adding more domains or composition:
  - the first unchecked backlog item is still Race query helpers
  - repo work has already moved into composition and off-backlog surfaces
  - `BACKLOG.md` must stay trustworthy as the project source of truth

- [X] Complete the missing Race query slice:
  - add `GetRaceByID`
  - expose the backlog query surface for seeded core races
  - keep returned race data defensive-copy safe

- [X] Fix core race seed correctness:
  - gnome is missing `Keen Senses`
  - half-elf is missing `Elven Immunities`
  - half-elf is missing `Keen Senses`
  - half-orc is missing `Orc Ferocity`
  - human, half-elf, and half-orc still represent the variable `+2` racial bonus as marker text instead of rules data

- [X] Stop duplicating structural race facts as ad-hoc racial features:
  - size already exists as `size`
  - speed already exists as `baseSpeed`
  - current seeds mix `Medium`, `Small`, `Normal Speed`, and `Slow Speed` into feature lists inconsistently

- [X] Correct the `Fine` size rules entry:
  - current space is `0 ft`
  - Pathfinder size data expects `0.5 ft` space

---

## SHOULD

- [ ] Resolve the P9 armor-check penalty dependency before composing skill check penalties:
  - independent audit on 2026-05-25 found a backlog ordering risk, not a current test failure
  - `BACKLOG.md` P9 asks to compose armor check penalty into skill totals using equipped armor and shield metadata
  - equipped armor and shield facts are not scheduled until P11 under `Core equipped inventory`
  - current `character` code models selected carryable equipment and carried weight, but not equipped armor or equipped shield facts
  - implementing the P9 armor-penalty item before equipped facts exist would require caller-provided side inputs, pull P11 work forward silently, or create a temporary modeling path
  - continue with the earlier P9 items, but before the armor-penalty item either introduce the equipped armor/shield facts explicitly or reshape the backlog dependency

- [X] Define the core language identity boundary before P7 bonus-language selection hardens around race-owned IDs:
  - independent audit before P7 found `BACKLOG.md`, `README.md`, and `docs/project-map.md` correctly point to P7 core character languages as the next path
  - `race.LanguageID` and `validLanguageIDs` currently live inside `internal/domain/rpg/character/race`
  - the current valid language set only covers languages needed by the seeded fixed core race lists: Common, Abyssal, Celestial, Draconic, Dwarven, Elven, Giant, Gnome, Gnoll, Goblin, Terran, Undercommon, Sylvan, Orc, and Halfling
  - local PF1 Linguistics text also lists common languages such as Aklo, Aquan, Auran, Ignan, and Infernal, and identifies Druidic as druids-only
  - local PF1 race text says humans and half-elves with high Intelligence can choose any language except secret languages such as Druidic
  - `race.BonusLanguageChoice` represents human and half-elf eligibility only as `AllowsAnyNonSecret()`, with no catalog or secret-language flag for `character` to validate concrete choices against
  - the first P7 automatic-language item is still safe, but the second P7 bonus-language item will otherwise either reject valid non-secret choices, accept arbitrary strings, or add side tables at the `character` boundary
  - added `internal/domain/rpg/character/language` as the core common language identity/query boundary with Druidic marked secret
  - `race.LanguageID` now aliases the language boundary and race language validation delegates to it
  - `BonusLanguageChoice.AllowsLanguageID` now lets later character composition validate fixed-list and any-non-secret choices without side tables or arbitrary strings

- [X] Refresh source-of-truth docs after P4 completion and new P5 audit path:
  - independent audit found `BACKLOG.md` now marks every P4 core feat selection context item complete
  - independent audit found the only unchecked backlog items before this audit were P9 far-future non-core items
  - refreshed `README.md` current status and near-term next steps to point at P5 core character ability composition
  - refreshed `docs/project-map.md` current goal and next backlog path to point at P5 core character ability composition
  - preserved backlog order and left P9 as far-future non-core work

- [X] Refresh README next-step guidance after the independent audit:
  - independent audit found `README.md` still says the next major tracked backlog area is core carryable equipment catalogs
  - `README.md` still says the next executable backlog item is the weapon query helper slice for seeded core simple weapons
  - `BACKLOG.md` and `docs/project-map.md` now show P4 core feat selection context is the next core-only path
  - stale README guidance can send future delivery back into already-completed P3 equipment work
  - update README status and near-term next steps without changing backlog order

- [X] Define the selected-weapon proficiency adapter boundary before composing P4 weapon feats:
  - independent audit found the next P4 item is valid but has a vocabulary mismatch risk
  - `class` stores weapon proficiency facts as `WeaponProficiencyID` labels such as `Simple Weapons`, `Martial Weapons`, and individual weapon names
  - `equipment` stores seeded weapons as `WeaponID` slugs plus `WeaponProficiencyCategory` values such as `Simple`, `Martial`, and `Exotic`
  - raw string matching between those surfaces would either fail valid core cases or couple domains accidentally
  - the first P4 backlog item should use an explicit `character`-boundary adapter for category and individual proficiency checks
  - cover category proficiency, individual proficiency, unknown or malformed selected weapons, and unsupported mappings
  - do not add attack rolls, damage rolls, feat bonus application, wielded slots, or combat behavior

- [X] Define the carryable item lookup boundary before seeding core weapons or armor:
  - `CharacterEquipment` currently accepts only `EquipmentID`
  - `CharacterEquipment.GetEquipment` resolves only through `equipment.GetEquipmentByID`
  - the seeded `Equipment` catalog currently covers the bounded adventuring-gear batch only
  - the new `Weapon` and `Armor` chassis have separate IDs while also carrying cost and weight facts
  - seeding core weapons or armor without a shared carryable lookup boundary would push inventory and carried weight toward duplicate side tables or special-case composition
  - decide the smallest query/composition surface that lets gear, weapons, armor, and shields participate in carried inventory without adding combat behavior

- [X] Restore an executable core backlog path before continuing delivery:
  - `BACKLOG.md` has completed P2 through the armor and shield chassis
  - the first unchecked backlog items are now P9 far-future non-core work
  - project rules say not to proceed into non-core content by default
  - add the next core-only milestone before asking agents to continue with "next step"
  - the likely next path is core equipment seed/query work after the carryable item boundary is clarified

- [X] Refresh README status after Feat and Equipment foundation work:
  - `README.md` still says the status matrix reflects delivery only through Class, Spellcasting Progression, and Spell
  - the status matrix omits the completed Feat and Equipment foundation surfaces
  - `README.md` still says the next major tracked backlog area is `feat`
  - `BACKLOG.md` and `docs/project-map.md` now show Feat is complete and the core equipment foundation path is complete through armor and shield chassis
  - stale README guidance can send later planning back toward already-completed work

- [X] Move non-core P2 behind the core equipment path:
  - `BACKLOG.md` now makes core equipment and inventory the next priority area after the completed character creation slice
  - non-core archetype, prestige class, race, feat, and spell work moved to far future
  - `docs/project-map.md` now points future work at equipment foundation instead of the completed vertical slice

- [X] Reshape the near-term backlog around a vertical character creation slice:
  - added a P1 milestone for the minimum level-1 core character creation slice
  - clarified the next executable backlog item as slice tests over existing composition adapters
  - added explicit out-of-scope boundaries for equipment, skill allocation, spell preparation, combat, full aggregate work, non-core sources, and folder reorganization
  - documented the intended vertical path in `docs/project-map.md`

- [X] Reduce repeated project context reads before continuing deeper character composition:
  - added `docs/project-map.md` as the compact package and workflow orientation guide
  - made `AGENTS.md` the single workflow source for delivery, planning, and audit modes
  - reduced `/internal/ai/agents/*` and `/internal/ai/skills/*` to short role reminders that point back to `AGENTS.md` and `docs/project-map.md`
  - linked the project map from `README.md`

- [X] Restore local PF1 extract/chunk tooling or update docs before rule-sensitive work depends on regeneration:
  - `Makefile` calls `docs/pf1/extract_rules.sh` and `docs/pf1/chunk_rules.go`
  - `docs/pf1/README.md` documents those commands as the supported extraction and chunk-generation workflow
  - restored `docs/pf1/extract_rules.sh` using local `pdftotext -layout`
  - restored `docs/pf1/chunk_rules.go` for deterministic local JSONL chunk generation
  - added `-dry-run` support so chunk generation can be validated without rewriting checked-in chunks

- [X] Align internal agent rule-source docs with `docs/pf1` before they misdirect rule lookups:
  - `AGENTS.md` says local PF1 rules are under `docs/pf1` and requires `rg` against `docs/pf1/chunks`
  - aligned `/internal/ai/agents/product-owner.md`, `/internal/ai/agents/senior-dev.md`, and `/internal/ai/agents/tech-lead.md`
  - aligned `/internal/ai/skills/codex.md`, `/internal/ai/skills/rules.md`, `/internal/ai/skills/architecture.md`, and `/internal/ai/skills/compound.md`
  - each now points to `docs/pf1`, local chunk search, and `docs/pf1/PFRPG_SRD.pdf`

- [X] Validate core feat prerequisite references against the seeded feat catalog before prerequisite composition:
  - `FeatPrerequisite`, `AnyFeatPrerequisite`, `SameSelectionFeatPrerequisite`, and `SpellSchoolFeatPrerequisite` still validate reusable prerequisite shape without coupling public constructors to the core catalog
  - core feat catalog initialization now validates every seeded feat-reference prerequisite against the seeded core feat ID set
  - missing referenced core feats now panic during package initialization with `missing referenced core feat seed`
  - added regression coverage for the seeded core catalog and misuse-boundary coverage for each feat-reference prerequisite shape

- [X] Expose read-only core feat catalog helpers before feat prerequisite composition reaches around package-private seeds:
  - core feat data now exists across `coreGeneralFeats`, `coreCombatFeats`, `coreCriticalFeats`, `coreItemCreationFeats`, and `coreMetamagicFeats`
  - all of those maps and order slices are package-private
  - outside `feat`, callers cannot ask for a seeded core feat by ID or enumerate the combined catalog
  - later `Compose feat prerequisites` will need to resolve feat prerequisites and selected feat ownership without duplicating scans or depending on test-only package internals
  - this should stay query-only with defensive-copy behavior and no new feat data

- [X] Tighten RaceID validation before race composition accepts malformed race identities:
  - `NewRace` currently accepts any non-empty `RaceID`, including values with surrounding whitespace such as `" human"`
  - `ClassID`, `FeatID`, and `SpellID` reject unnormalized IDs, but `RaceID` does not
  - `GetRaceByID` only resolves canonical core IDs, so malformed `Race` values can be constructed but will not round-trip through the seeded catalog
  - upcoming race composition should not have to guard against race values that the race domain itself could have rejected

- [X] Add public read-only query helpers for core spell data and spell list bindings before spell composition depends on package-private seeds:
  - `coreSpells` and `coreSpellListEntries` are complete enough to test, but outside `spell` cannot look up a spell by ID or ask for class spell-list entries
  - current class-list lookup coverage lives in test helpers only, so production composition would need to duplicate scans or reach around package boundaries
  - later `Compose spell list entry with Class spellcasting` needs a stable read surface such as spell lookup, class list lookup, and defensive-copy catalog helpers
  - this should stay query-only and must not add new spell data or composition logic

- [X] Expose core spellcasting progression lookup helpers before character spell-slot composition starts:
  - `coreSpellcastingProgressionTables` and `coreSpellcastingProgressionClassOrder` are package-private
  - outside `class`, callers can construct arbitrary progression tables, but cannot ask for the seeded core progression for bard, cleric, druid, paladin, ranger, sorcerer, or wizard
  - later character composition will otherwise duplicate progression access logic or depend on unexported test-only knowledge
  - add a read-only lookup/catalog surface with defensive copies and no new progression rules

- [X] Refresh README status now that Class, Spellcasting Progression, and Spell are complete:
  - `README.md` still says the status matrix reflects Class not being aligned yet
  - the matrix row for `class` still says `Exists: no`, `Core-correct now: not started`, and `Intentional limit: next backlog domain`
  - the backlog now marks Class, Spellcasting Progression, Spell data, and Spell tests complete
  - stale status docs can mislead Feat review by making already-delivered foundations look unavailable

- [X] Add a read-only core class catalog helper for consistency before composition expands class consumers:
  - `race` exposes both `GetRaceByID` and `GetRaces`
  - `skill` exposes both `GetSkillByID` and `GetSkills`
  - `class` exposes `GetClassByID`, but `coreClassOrder` remains package-private and there is no `GetClasses`
  - later character and feat prerequisite composition will need ordered class discovery without relying on private seed state

- [X] Move core spell-list seed consistency from tests into the seed-building path before spell data is consumed by other domains:
  - `mustBuildCoreSpellListEntries` validates entry shape and duplicate triples, but it does not prove each listed spell has seeded `Spell` data
  - `TestCoreSpellData_SeedsAllCoreSpells` catches that today, but production package initialization can still build a mismatched core spell list if tests are bypassed
  - once spell lookups become public, a spell-list entry pointing at missing spell data would become a runtime lookup inconsistency
  - keep the fix seed-local: do not make generic `SpellListEntry` construction depend on the core catalog

- [X] Remove or rename the legacy `GetRacialLanguages` alias before callers lock in the wrong race-language shape again:
  - `Race` now has `GetAutomaticLanguages` and `GetBonusLanguageChoice`
  - `GetRacialLanguages` still exists and currently returns only the automatic slice
  - the old name reads like "all racial languages" and hides the newly modeled bonus-language metadata
  - current tests still normalize that legacy entry point instead of pushing callers to the split model

- [X] Refresh stale README guidance before docs start reintroducing already-fixed mistakes:
  - `README.md` still says `HitDie` average HP uses the project's fixed averages
  - `README.md` still says the next major area is creature type and subtype
  - both statements are now false and can send future work and reviews down the wrong path
  - stale source-of-truth guidance increases the chance of fixing symptoms instead of following the current model

- [X] Expose a read-only skill query surface for the seeded core catalog:
  - `coreSkills` exists but is package-private
  - outside consumers cannot ask for a skill by ID or enumerate the seeded catalog
  - later class-skill and character work will have to reach around the seed instead of using it

- [X] Tighten Race chassis invariants around selectable ability-score bonuses:
  - `NewRace` accepts any non-negative `selectableAbilityScoreModifier`
  - it currently allows semantically impossible combinations, like fixed modifiers plus a selectable bonus on the same race
  - the current core seeds are fine, but the chassis still permits invalid states

- [X] Make the current repo status readable from the source-of-truth docs:
  - `BACKLOG.md` starts at Race / Modifier / Skill / Class
  - the repo also contains established `ability`, `creaturetype`, and `character` packages that are not reflected there
  - a reviewer currently has to reconcile code, README, and git history to understand what is foundational versus backlog-delivered

- [X] Clarify the humanoid misuse boundary before class and character composition expands:
  - `ResolvedCreatureRules.NewRacialHitDie` and `NewRacialHitPoints` happily build humanoid racial HD / HP
  - the contextual flag exists, but the convenience API still makes the wrong path look normal
  - this is easy to use incorrectly once Class arrives

- [X] Make `creaturetype` scope easier to read from the package surface:
  - exported type IDs suggest broad PF1 coverage
  - only a limited subtype set and a partial trait model are actually implemented
  - most of that boundary is discoverable only by reading tests and `profile_table.go`

- [X] Quarantine project house rules from upcoming core backlog work:
  - `TitanicSize`
  - custom construct bonus HP table
  - source-bucket caster levels
  - these should not blur what is core-only work vs project-specific rules

- [X] Replace free-form race language and feature strings with canonical IDs before expanding lookup usage:
  - current validation only rejects empty strings
  - typo and casing drift will leak into query helpers

- [X] Bring repository documentation back in sync with repo rules:
  - `README.md` emphasizes extensibility and homebrew
  - `AGENTS.md` and `BACKLOG.md` enforce backlog-first, core-first delivery

- [X] Make backlog status match repo reality:
  - some unchecked Race helper and test work already exists
  - repo packages like `modifier` are not represented in `BACKLOG.md`

---

## CAN

- [ ] Keep class-granted and class-substituted language rules out of the P7 race-language path until class feature language work is scheduled:
  - local PF1 chunks show class-specific language rules outside race metadata: clerics add Celestial, Abyssal, and Infernal to bonus language options; druids add Sylvan and know Druidic as a secret free language; wizards may substitute Draconic for one racial bonus language
  - P7 is scoped to race language metadata plus Intelligence-derived bonus language capacity
  - folding class-specific grants or substitutions into P7 would expand the slice beyond the current backlog item
  - leaving this visible prevents later character-language work from mistaking P7 race-language facts for complete PF1 language coverage

- [X] Add a compact status matrix for specialist review:
  - what exists
  - what is core-correct
  - what is intentionally partial
  - what is project-specific

- [X] Normalize the remaining vague PF1 seed labels where they are still less explicit than the rulebook:
  - example: elf `Magic` is less precise than `Elven Magic`

- [X] Add read-only discovery helpers for seeded catalogs if lookup pressure grows:
  - race has `GetRaceByID` but no catalog iterator
  - skill has no public discovery surface yet

- [X] Run a later audit pass over the already-built foundation domains after P0 Race, Skill, and Class are aligned

- [X] Collapse duplicated helper logic only if it becomes active maintenance pain:
  - example: ability modifier math is duplicated in more than one place
