package character

import (
	"testing"

	"d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterequipment "d20campaigngenerator/internal/domain/rpg/character/equipment"
	characterfeat "d20campaigngenerator/internal/domain/rpg/character/feat"
	characterlanguage "d20campaigngenerator/internal/domain/rpg/character/language"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
	characterspell "d20campaigngenerator/internal/domain/rpg/character/spell"
)

func TestMinimumLevelOneCharacterCreationSlice_ComposesCoreRaceClassHPSpellcastingAndFeatPrerequisites(t *testing.T) {
	selectedRace := mustNewCharacterRaceForSliceTest(t, characterrace.HumanRaceID)
	race, ok := selectedRace.GetRace()
	if !ok {
		t.Fatal("expected selected race to resolve")
	}

	if race.GetSize() != ability.MediumSize {
		t.Fatalf("expected human size %q, got %q", ability.MediumSize, race.GetSize())
	}

	selectableModifier, ok := race.GetSelectableAbilityScoreModifier()
	if !ok || selectableModifier != 2 {
		t.Fatalf("expected human selectable ability modifier (2, true), got (%d, %t)", selectableModifier, ok)
	}

	selectedClass := mustNewCharacterClassForTest(t, characterclass.WizardClassID)
	class, ok := selectedClass.GetClass()
	if !ok {
		t.Fatal("expected selected class to resolve")
	}

	if class.GetHitDieType() != ability.D6HitDie {
		t.Fatalf("expected wizard hit die %q, got %q", ability.D6HitDie, class.GetHitDieType())
	}

	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1),
	}
	levelFacts, ok := NewCharacterLevelFacts(classLevels)
	if !ok {
		t.Fatal("expected wizard class level facts to compose")
	}

	if levelFacts.GetTotalCharacterLevel() != 1 {
		t.Fatalf("expected total character level 1, got %d", levelFacts.GetTotalCharacterLevel())
	}

	wizardLevel, ok := levelFacts.GetClassLevel(characterclass.WizardClassID)
	if !ok || wizardLevel != 1 {
		t.Fatalf("expected wizard class level (1, true), got (%d, %t)", wizardLevel, ok)
	}

	baseAttackBonusFacts, ok := NewCharacterBaseAttackBonusFacts(levelFacts.GetClassLevels())
	if !ok {
		t.Fatal("expected wizard base attack bonus facts to compose")
	}

	if baseAttackBonusFacts.GetActualValue().GetNumerator() != 1 ||
		baseAttackBonusFacts.GetActualValue().GetDenominator() != 2 ||
		baseAttackBonusFacts.GetValue() != 0 {
		t.Fatalf(
			"expected wizard BAB 1/2 and value 0, got %d/%d and value %d",
			baseAttackBonusFacts.GetActualValue().GetNumerator(),
			baseAttackBonusFacts.GetActualValue().GetDenominator(),
			baseAttackBonusFacts.GetValue(),
		)
	}

	baseSaveFacts, ok := NewCharacterBaseSavingThrowFacts(levelFacts.GetClassLevels())
	if !ok {
		t.Fatal("expected wizard base saving throw facts to compose")
	}

	assertSliceBaseSave(t, baseSaveFacts, ability.FortitudeSave, 1, 3, 0)
	assertSliceBaseSave(t, baseSaveFacts, ability.ReflexSave, 1, 3, 0)
	assertSliceBaseSave(t, baseSaveFacts, ability.WillSave, 5, 2, 2)

	hpFacts, ok := NewCharacterClassHitPointFacts(
		levelFacts.GetClassLevels(),
		characterclass.WizardClassID,
		14,
		nil,
	)
	if !ok {
		t.Fatal("expected wizard class HP facts to compose")
	}

	hp, ok := hpFacts.GetHitPoints()
	if !ok {
		t.Fatal("expected wizard class HP facts to expose hit points")
	}

	if hp.GetTotal() != 8 || hp.GetCurrent() != 8 {
		t.Fatalf("expected first-level wizard total/current HP 8, got total %d current %d", hp.GetTotal(), hp.GetCurrent())
	}

	d6Count, ok := hp.GetHitDie().GetDieCount(ability.D6HitDie)
	if !ok || d6Count != 1 {
		t.Fatalf("expected first-level wizard hit die (1, true), got (%d, %t)", d6Count, ok)
	}

	progression, ok := NewCharacterSpellcastingProgression(selectedClass)
	if !ok {
		t.Fatal("expected wizard spellcasting progression to compose")
	}

	cantrips, ok := progression.GetSpellSlots(1, 0)
	if !ok || cantrips != 3 {
		t.Fatalf("expected wizard level 1 cantrip slots (3, true), got (%d, %t)", cantrips, ok)
	}

	firstLevelSlots, ok := progression.GetSpellSlots(1, 1)
	if !ok || firstLevelSlots != 1 {
		t.Fatalf("expected wizard level 1 spell slots (1, true), got (%d, %t)", firstLevelSlots, ok)
	}

	if _, ok := progression.GetSpellSlots(1, 2); ok {
		t.Fatal("expected wizard level 1 spell level 2 slots to be unavailable")
	}

	prerequisiteState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		baseAttackBonusFacts.GetValue(),
		levelFacts.GetClassLevels(),
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ArcaneStrikeFeatID, prerequisiteState); !ok {
		t.Fatal("expected arcane strike to compose from level-1 wizard spellcasting")
	}

	if _, ok := NewCharacterFeat(characterfeat.PowerAttackFeatID, prerequisiteState); ok {
		t.Fatal("expected power attack to fail without strength and base attack prerequisites")
	}
}

func TestMinimumLevelOneCharacterCreationSlice_ComposesSelectedFeatContexts(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)
	fighterClassLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1),
	}
	fighterBaseAttackBonusFacts, ok := NewCharacterBaseAttackBonusFacts(fighterClassLevels)
	if !ok {
		t.Fatal("expected fighter base attack bonus facts to compose")
	}

	weaponFeatState := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
		t,
		nil,
		fighterBaseAttackBonusFacts.GetValue(),
		nil,
		fighterClassLevels,
		nil,
		nil,
		selectedWeapon,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.WeaponFocusFeatID, weaponFeatState); !ok {
		t.Fatal("expected weapon focus to compose from a level-1 fighter and selected dagger")
	}

	conjurationSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID)
	spellFocus := mustNewCharacterSelectedSpellSchoolFeatForTest(t, characterfeat.SpellFocusFeatID, conjurationSchool)
	wizardClassLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1),
	}
	spellSchoolFeatState := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		wizardClassLevels,
		nil,
		nil,
		conjurationSchool,
		[]CharacterSelectedSpellSchoolFeat{spellFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterSpellFocusFeatID, spellSchoolFeatState); !ok {
		t.Fatal("expected greater spell focus to compose from spell focus with the same selected school")
	}

	evocationSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.EvocationSchoolID)
	mismatchedSpellSchoolState := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		wizardClassLevels,
		nil,
		nil,
		evocationSchool,
		[]CharacterSelectedSpellSchoolFeat{spellFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterSpellFocusFeatID, mismatchedSpellSchoolState); ok {
		t.Fatal("expected greater spell focus to reject spell focus with a mismatched selected school")
	}
}

func TestMinimumLevelOneCharacterCreationSlice_ComposesRacialAbilityContexts(t *testing.T) {
	elfAbilityScores, ok := NewFixedRacialCharacterAbilityScores(
		mustNewCharacterRaceForSliceTest(t, characterrace.ElfRaceID),
		[]CharacterAbilityScore{
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.StrengthScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.DexterityScore, 11),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.ConstitutionScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.WisdomScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.CharismaScore, 10),
		},
	)
	if !ok {
		t.Fatal("expected elf fixed racial ability modifiers to compose")
	}

	dodgeState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		elfAbilityScores,
		0,
		nil,
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.DodgeFeatID, dodgeState); !ok {
		t.Fatal("expected dodge to compose from elf-adjusted dexterity")
	}

	humanAbilityScores, ok := NewSelectableRacialCharacterAbilityScores(
		mustNewCharacterRaceForSliceTest(t, characterrace.HumanRaceID),
		[]CharacterAbilityScore{
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.StrengthScore, 11),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.DexterityScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.ConstitutionScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.WisdomScore, 10),
			mustNewCharacterAbilityScoreForAbilityTest(t, ability.CharismaScore, 10),
		},
		[]CharacterSelectedAbilityScore{
			mustNewCharacterSelectedAbilityScoreForTest(t, ability.StrengthScore),
		},
	)
	if !ok {
		t.Fatal("expected human selectable racial ability modifier to compose")
	}

	fighterClassLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1),
	}
	fighterBaseAttackBonusFacts, ok := NewCharacterBaseAttackBonusFacts(fighterClassLevels)
	if !ok {
		t.Fatal("expected fighter base attack bonus facts to compose")
	}

	powerAttackState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		humanAbilityScores,
		fighterBaseAttackBonusFacts.GetValue(),
		fighterClassLevels,
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.PowerAttackFeatID, powerAttackState); !ok {
		t.Fatal("expected power attack to compose from human-adjusted strength")
	}

	backpack, ok := NewCharacterAdventuringGear(characterequipment.BackpackEmptyEquipmentID, 1)
	if !ok {
		t.Fatal("expected backpack adventuring gear to compose")
	}

	carriedWeight, ok := NewCharacterCarriedWeight(
		mustNewAbilityScoreFromCharacterScoresForSliceTest(t, humanAbilityScores, ability.StrengthScore),
		[]CharacterEquipment{backpack},
	)
	if !ok {
		t.Fatal("expected carried weight to compose from human-adjusted strength")
	}

	if carriedWeight.GetLoadCategory() != LightLoadCategory {
		t.Fatalf("expected backpack to stay in light load, got %q", carriedWeight.GetLoadCategory())
	}
}

func TestMinimumLevelOneCharacterCreationSlice_RacialAbilitySelectionsFailClosed(t *testing.T) {
	if _, ok := NewSelectableRacialCharacterAbilityScores(
		mustNewCharacterRaceForSliceTest(t, characterrace.HumanRaceID),
		mustNewBaseCharacterAbilityScoresForTest(t, 10),
		nil,
	); ok {
		t.Fatal("expected human selectable racial ability composition to fail without a selected ability")
	}

	if _, ok := NewSelectableRacialCharacterAbilityScores(
		mustNewCharacterRaceForSliceTest(t, characterrace.DwarfRaceID),
		mustNewBaseCharacterAbilityScoresForTest(t, 10),
		[]CharacterSelectedAbilityScore{
			mustNewCharacterSelectedAbilityScoreForTest(t, ability.StrengthScore),
		},
	); ok {
		t.Fatal("expected selected ability input to fail against a fixed-modifier race")
	}
}

func TestMinimumLevelOneCharacterCreationSlice_ComposesLanguageContexts(t *testing.T) {
	elfRace := mustNewCharacterRaceForSliceTest(t, characterrace.ElfRaceID)
	elfAutomaticLanguages, ok := NewAutomaticRacialCharacterLanguageFacts(elfRace)
	if !ok {
		t.Fatal("expected elf automatic racial language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		elfAutomaticLanguages.GetLanguageIDs(),
		[]characterlanguage.LanguageID{
			characterlanguage.CommonLanguageID,
			characterlanguage.ElvenLanguageID,
		},
	)

	elfBonusLanguages, ok := NewBonusRacialCharacterLanguageFacts(
		elfRace,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 14)},
		[]characterlanguage.LanguageID{
			characterlanguage.CelestialLanguageID,
			characterlanguage.DraconicLanguageID,
		},
	)
	if !ok {
		t.Fatal("expected elf fixed-list bonus language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		elfBonusLanguages.GetLanguageIDs(),
		[]characterlanguage.LanguageID{
			characterlanguage.CelestialLanguageID,
			characterlanguage.DraconicLanguageID,
		},
	)

	if elfAutomaticLanguages.HasLanguage(characterlanguage.CelestialLanguageID) {
		t.Fatal("expected automatic language facts not to include selected bonus languages")
	}

	humanRace := mustNewCharacterRaceForSliceTest(t, characterrace.HumanRaceID)
	humanAutomaticLanguages, ok := NewAutomaticRacialCharacterLanguageFacts(humanRace)
	if !ok {
		t.Fatal("expected human automatic racial language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		humanAutomaticLanguages.GetLanguageIDs(),
		[]characterlanguage.LanguageID{characterlanguage.CommonLanguageID},
	)

	humanBonusLanguages, ok := NewBonusRacialCharacterLanguageFacts(
		humanRace,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 12)},
		[]characterlanguage.LanguageID{characterlanguage.InfernalLanguageID},
	)
	if !ok {
		t.Fatal("expected human any-non-secret bonus language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		humanBonusLanguages.GetLanguageIDs(),
		[]characterlanguage.LanguageID{characterlanguage.InfernalLanguageID},
	)
}

func TestMinimumLevelOneCharacterCreationSlice_LanguageSelectionsFailClosed(t *testing.T) {
	humanRace := mustNewCharacterRaceForSliceTest(t, characterrace.HumanRaceID)
	dwarfRace := mustNewCharacterRaceForSliceTest(t, characterrace.DwarfRaceID)
	intelligence12 := []CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 12)}

	if _, ok := NewBonusRacialCharacterLanguageFacts(
		humanRace,
		intelligence12,
		[]characterlanguage.LanguageID{
			characterlanguage.InfernalLanguageID,
			characterlanguage.AquanLanguageID,
		},
	); ok {
		t.Fatal("expected over-budget bonus language selection to fail")
	}

	if _, ok := NewBonusRacialCharacterLanguageFacts(
		dwarfRace,
		intelligence12,
		[]characterlanguage.LanguageID{characterlanguage.CelestialLanguageID},
	); ok {
		t.Fatal("expected disallowed fixed-list language selection to fail")
	}

	if _, ok := NewBonusRacialCharacterLanguageFacts(
		humanRace,
		intelligence12,
		[]characterlanguage.LanguageID{characterlanguage.DruidicLanguageID},
	); ok {
		t.Fatal("expected secret language selection to fail")
	}
}

func TestMinimumLevelOneCharacterCreationSlice_InvalidSelectedInputsFailClosed(t *testing.T) {
	if _, ok := NewCharacterRace(characterrace.RaceID("android")); ok {
		t.Fatal("expected unknown race to fail")
	}

	if _, ok := NewCharacterClass(characterclass.ClassID("alchemist")); ok {
		t.Fatal("expected unknown class to fail")
	}

	fighterLevel := mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)
	if _, ok := NewCharacterLevelFacts([]CharacterClassLevel{fighterLevel, fighterLevel}); ok {
		t.Fatal("expected duplicate class levels to fail")
	}

	if _, ok := NewCharacterBaseAttackBonusFacts([]CharacterClassLevel{{classID: characterclass.FighterClassID, level: 0}}); ok {
		t.Fatal("expected malformed class level BAB facts to fail")
	}

	if _, ok := NewCharacterBaseSavingThrowFacts([]CharacterClassLevel{{classID: characterclass.ClassID("alchemist"), level: 1}}); ok {
		t.Fatal("expected non-core class level save facts to fail")
	}

	var zeroClass CharacterClass
	if _, ok := NewFirstLevelCharacterHitPoints(zeroClass, 14); ok {
		t.Fatal("expected zero-value class hit points to fail")
	}

	selectedClass := mustNewCharacterClassForTest(t, characterclass.WizardClassID)
	if _, ok := NewFirstLevelCharacterHitPoints(selectedClass, 0); ok {
		t.Fatal("expected invalid constitution hit points to fail")
	}

	fighter := mustNewCharacterClassForTest(t, characterclass.FighterClassID)
	if _, ok := NewCharacterSpellcastingProgression(fighter); ok {
		t.Fatal("expected non-spellcasting class progression to fail")
	}

	var prerequisiteState CharacterFeatPrerequisiteState
	if _, ok := NewCharacterFeat(characterfeat.ArcaneStrikeFeatID, prerequisiteState); ok {
		t.Fatal("expected zero-value feat prerequisite state to fail")
	}
}

func mustNewCharacterRaceForSliceTest(
	t *testing.T,
	id characterrace.RaceID,
) CharacterRace {
	t.Helper()

	race, ok := NewCharacterRace(id)
	if !ok {
		t.Fatalf("expected character race %q to compose", id)
	}

	return race
}

func mustNewAbilityScoreFromCharacterScoresForSliceTest(
	t *testing.T,
	scores []CharacterAbilityScore,
	id ability.AbilityScoreID,
) ability.AbilityScore {
	t.Helper()

	for _, score := range scores {
		if score.GetAbilityScoreID() != id {
			continue
		}

		value, ok := ability.NewAbilityScoreValue(score.GetScore(), true)
		if !ok {
			t.Fatalf("expected ability score value %d to compose", score.GetScore())
		}

		abilityScore, ok := ability.NewAbilityScore(id, value)
		if !ok {
			t.Fatalf("expected ability score %q to compose", id)
		}

		return abilityScore
	}

	t.Fatalf("expected composed ability score %q to exist", id)
	return ability.AbilityScore{}
}

func assertSliceBaseSave(
	t *testing.T,
	facts CharacterBaseSavingThrowFacts,
	id ability.SavingThrowID,
	expectedNumerator int,
	expectedDenominator int,
	expectedValue int,
) {
	t.Helper()

	save, ok := facts.GetSavingThrow(id)
	if !ok {
		t.Fatalf("expected save %q to resolve", id)
	}

	if save.GetActualValue().GetNumerator() != expectedNumerator ||
		save.GetActualValue().GetDenominator() != expectedDenominator ||
		save.GetValue() != expectedValue {
		t.Fatalf(
			"expected %q save %d/%d and value %d, got %d/%d and value %d",
			id,
			expectedNumerator,
			expectedDenominator,
			expectedValue,
			save.GetActualValue().GetNumerator(),
			save.GetActualValue().GetDenominator(),
			save.GetValue(),
		)
	}
}
