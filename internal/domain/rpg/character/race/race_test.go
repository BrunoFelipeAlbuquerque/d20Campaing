package race

import (
	"testing"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterlanguage "d20campaigngenerator/internal/domain/rpg/character/language"
)

func TestNewRace_ConstructsValidatedRaceChassis(t *testing.T) {
	strengthModifier, ok := NewAbilityScoreModifier(ability.StrengthScore, 2)
	if !ok {
		t.Fatal("expected strength modifier to be constructed")
	}

	constitutionModifier, ok := NewAbilityScoreModifier(ability.ConstitutionScore, -2)
	if !ok {
		t.Fatal("expected constitution modifier to be constructed")
	}

	race, ok := NewRace(
		RaceID("elf"),
		ability.MediumSize,
		30,
		[]AbilityScoreModifier{strengthModifier, constitutionModifier},
		0,
		[]LanguageID{CommonLanguageID, ElvenLanguageID},
		mustBonusLanguageChoiceForTest(t, []LanguageID{SylvanLanguageID}, false),
		[]RacialFeatureID{KeenSensesFeatureID, LowLightVisionFeatureID},
	)
	if !ok {
		t.Fatal("expected race chassis to be constructed")
	}

	if race.GetID() != RaceID("elf") {
		t.Fatalf("expected race id %q, got %q", RaceID("elf"), race.GetID())
	}

	if race.GetSize() != ability.MediumSize {
		t.Fatalf("expected race size %q, got %q", ability.MediumSize, race.GetSize())
	}

	if race.GetBaseSpeed() != 30 {
		t.Fatalf("expected base speed 30, got %d", race.GetBaseSpeed())
	}

	modifiers := race.GetAbilityScoreModifiers()
	if len(modifiers) != 2 {
		t.Fatalf("expected 2 ability score modifiers, got %d", len(modifiers))
	}

	if modifiers[0].GetScoreID() != ability.StrengthScore || modifiers[0].GetModifier() != 2 {
		t.Fatalf("expected first modifier to be STR +2, got (%q, %d)", modifiers[0].GetScoreID(), modifiers[0].GetModifier())
	}

	languages := race.GetAutomaticLanguages()
	if len(languages) != 2 || languages[0] != CommonLanguageID || languages[1] != ElvenLanguageID {
		t.Fatalf("expected automatic languages [Common Elven], got %v", languages)
	}

	bonusChoice, ok := race.GetBonusLanguageChoice()
	if !ok {
		t.Fatal("expected bonus language choice metadata to be present")
	}

	if len(bonusChoice.GetLanguageIDs()) != 1 || bonusChoice.GetLanguageIDs()[0] != SylvanLanguageID {
		t.Fatalf("expected bonus language choice [Sylvan], got %v", bonusChoice.GetLanguageIDs())
	}

	if !race.HasFeature(KeenSensesFeatureID) {
		t.Fatal("expected Keen Senses feature to be present")
	}

	if _, ok := race.GetSelectableAbilityScoreModifier(); ok {
		t.Fatal("expected fixed-modifier race chassis not to expose selectable ability score modifier metadata")
	}
}

func TestNewRace_StoresSelectableAbilityScoreModifierMetadata(t *testing.T) {
	race, ok := NewRace(
		RaceID("human"),
		ability.MediumSize,
		30,
		nil,
		2,
		[]LanguageID{CommonLanguageID},
		mustBonusLanguageChoiceForTest(t, nil, true),
		[]RacialFeatureID{BonusFeatFeatureID, SkilledFeatureID},
	)
	if !ok {
		t.Fatal("expected race chassis with selectable ability score modifier metadata to be constructed")
	}

	selectableModifier, ok := race.GetSelectableAbilityScoreModifier()
	if !ok {
		t.Fatal("expected selectable ability score modifier metadata to be present")
	}

	if selectableModifier != 2 {
		t.Fatalf("expected selectable ability score modifier 2, got %d", selectableModifier)
	}
}

func TestBonusLanguageChoice_AllowsFixedOrAnyNonSecretLanguages(t *testing.T) {
	fixedChoice := mustBonusLanguageChoiceForTest(t, []LanguageID{ElvenLanguageID}, false)

	if !fixedChoice.AllowsLanguageID(ElvenLanguageID) {
		t.Fatal("expected fixed-list bonus language choice to allow listed language")
	}

	if fixedChoice.AllowsLanguageID(characterlanguage.InfernalLanguageID) {
		t.Fatal("expected fixed-list bonus language choice to reject unlisted language")
	}

	anyNonSecretChoice := mustBonusLanguageChoiceForTest(t, nil, true)

	if !anyNonSecretChoice.AllowsLanguageID(characterlanguage.InfernalLanguageID) {
		t.Fatal("expected any-non-secret bonus language choice to allow known non-secret language")
	}

	if anyNonSecretChoice.AllowsLanguageID(characterlanguage.DruidicLanguageID) {
		t.Fatal("expected any-non-secret bonus language choice to reject secret language")
	}

	if anyNonSecretChoice.AllowsLanguageID(LanguageID("Thieves' Cant")) {
		t.Fatal("expected any-non-secret bonus language choice to reject unknown language")
	}
}

func TestNewRace_DedupesModifiersLanguagesAndFeatures(t *testing.T) {
	intelligenceModifier, ok := NewAbilityScoreModifier(ability.IntelligenceScore, 2)
	if !ok {
		t.Fatal("expected intelligence modifier to be constructed")
	}

	race, ok := NewRace(
		RaceID("gnome"),
		ability.SmallSize,
		20,
		[]AbilityScoreModifier{intelligenceModifier, intelligenceModifier},
		0,
		[]LanguageID{CommonLanguageID, GnomeLanguageID, CommonLanguageID},
		mustBonusLanguageChoiceForTest(t, []LanguageID{ElvenLanguageID, ElvenLanguageID}, false),
		[]RacialFeatureID{DefensiveTrainingFeatureID, DefensiveTrainingFeatureID, KeenSensesFeatureID},
	)
	if !ok {
		t.Fatal("expected race chassis to be constructed")
	}

	if len(race.GetAbilityScoreModifiers()) != 1 {
		t.Fatalf("expected deduped ability score modifier length 1, got %d", len(race.GetAbilityScoreModifiers()))
	}

	if len(race.GetAutomaticLanguages()) != 2 {
		t.Fatalf("expected deduped automatic languages length 2, got %d", len(race.GetAutomaticLanguages()))
	}

	bonusChoice, ok := race.GetBonusLanguageChoice()
	if !ok {
		t.Fatal("expected bonus language choice metadata to be present")
	}

	if len(bonusChoice.GetLanguageIDs()) != 1 {
		t.Fatalf("expected deduped bonus language choice length 1, got %d", len(bonusChoice.GetLanguageIDs()))
	}

	if len(race.GetRacialFeatures()) != 2 {
		t.Fatalf("expected deduped racial features length 2, got %d", len(race.GetRacialFeatures()))
	}
}

func TestNewRace_RejectsInvalidInputs(t *testing.T) {
	if _, ok := NewAbilityScoreModifier(ability.AbilityScoreID("LCK"), 2); ok {
		t.Fatal("expected invalid ability score id to be rejected")
	}

	if _, ok := NewAbilityScoreModifier(ability.StrengthScore, 0); ok {
		t.Fatal("expected zero ability score modifier to be rejected")
	}

	validModifier, ok := NewAbilityScoreModifier(ability.DexterityScore, 2)
	if !ok {
		t.Fatal("expected dexterity modifier to be constructed")
	}

	if _, ok := NewRace("", ability.MediumSize, 30, nil, 0, nil, bonusLanguageChoice{}, nil); ok {
		t.Fatal("expected empty race id to be rejected")
	}

	for _, id := range []RaceID{" human", "human ", "\thuman"} {
		if _, ok := NewRace(id, ability.MediumSize, 30, nil, 0, nil, bonusLanguageChoice{}, nil); ok {
			t.Fatalf("expected unnormalized race id %q to be rejected", id)
		}
	}

	if _, ok := NewRace(RaceID("human"), ability.Size("Gigantic"), 30, nil, 0, nil, bonusLanguageChoice{}, nil); ok {
		t.Fatal("expected invalid size to be rejected")
	}

	if _, ok := NewRace(RaceID("human"), ability.MediumSize, 0, nil, 0, nil, bonusLanguageChoice{}, nil); ok {
		t.Fatal("expected non-positive base speed to be rejected")
	}

	if _, ok := NewRace(RaceID("human"), ability.MediumSize, 30, nil, -2, nil, bonusLanguageChoice{}, nil); ok {
		t.Fatal("expected negative selectable ability score modifier to be rejected")
	}

	if _, ok := NewRace(RaceID("human"), ability.MediumSize, 30, nil, 1, nil, bonusLanguageChoice{}, nil); ok {
		t.Fatal("expected unsupported selectable ability score modifier to be rejected")
	}

	if _, ok := NewRace(
		RaceID("human"),
		ability.MediumSize,
		30,
		[]AbilityScoreModifier{{scoreID: ability.AbilityScoreID("LCK"), modifier: 2}},
		0,
		nil,
		bonusLanguageChoice{},
		nil,
	); ok {
		t.Fatal("expected invalid ability score modifier entry to be rejected")
	}

	if _, ok := NewRace(
		RaceID("human"),
		ability.MediumSize,
		30,
		[]AbilityScoreModifier{validModifier},
		2,
		nil,
		bonusLanguageChoice{},
		nil,
	); ok {
		t.Fatal("expected fixed and selectable ability score modifiers to be rejected together")
	}

	if _, ok := NewRace(
		RaceID("human"),
		ability.MediumSize,
		30,
		[]AbilityScoreModifier{validModifier},
		0,
		[]LanguageID{LanguageID("common")},
		bonusLanguageChoice{},
		nil,
	); ok {
		t.Fatal("expected unknown racial language to be rejected")
	}

	if _, ok := NewRace(
		RaceID("human"),
		ability.MediumSize,
		30,
		[]AbilityScoreModifier{validModifier},
		0,
		nil,
		BonusLanguageChoice{languageIDs: []LanguageID{LanguageID("unknown")}},
		nil,
	); ok {
		t.Fatal("expected invalid bonus language choice to be rejected")
	}

	if _, ok := NewRace(
		RaceID("human"),
		ability.MediumSize,
		30,
		[]AbilityScoreModifier{validModifier},
		0,
		nil,
		bonusLanguageChoice{},
		[]RacialFeatureID{RacialFeatureID("keen senses")},
	); ok {
		t.Fatal("expected unknown racial feature to be rejected")
	}
}

func TestRace_GettersReturnDefensiveCopies(t *testing.T) {
	dexterityModifier, ok := NewAbilityScoreModifier(ability.DexterityScore, 2)
	if !ok {
		t.Fatal("expected dexterity modifier to be constructed")
	}

	race, ok := NewRace(
		RaceID("halfling"),
		ability.SmallSize,
		20,
		[]AbilityScoreModifier{dexterityModifier},
		0,
		[]LanguageID{CommonLanguageID, HalflingLanguageID},
		mustBonusLanguageChoiceForTest(t, []LanguageID{ElvenLanguageID}, false),
		[]RacialFeatureID{SureFootedFeatureID},
	)
	if !ok {
		t.Fatal("expected race chassis to be constructed")
	}

	modifiers := race.GetAbilityScoreModifiers()
	languages := race.GetAutomaticLanguages()
	bonusChoice, _ := race.GetBonusLanguageChoice()
	features := race.GetRacialFeatures()

	modifiers[0] = AbilityScoreModifier{}
	languages[0] = "Changed"
	bonusChoice.languageIDs[0] = "Changed"
	features[0] = "Changed"

	if race.GetAbilityScoreModifiers()[0].GetScoreID() != ability.DexterityScore {
		t.Fatal("expected ability score modifiers getter to return a defensive copy")
	}

	if race.GetAutomaticLanguages()[0] != CommonLanguageID {
		t.Fatal("expected automatic languages getter to return a defensive copy")
	}

	storedBonusChoice, ok := race.GetBonusLanguageChoice()
	if !ok || storedBonusChoice.GetLanguageIDs()[0] != ElvenLanguageID {
		t.Fatal("expected bonus language choice getter to return a defensive copy")
	}

	if race.GetRacialFeatures()[0] != SureFootedFeatureID {
		t.Fatal("expected racial features getter to return a defensive copy")
	}
}

func mustBonusLanguageChoiceForTest(t *testing.T, languageIDs []LanguageID, anyNonSecret bool) BonusLanguageChoice {
	t.Helper()

	value, ok := NewBonusLanguageChoice(languageIDs, anyNonSecret)
	if !ok {
		t.Fatal("expected bonus language choice to be constructed")
	}

	return value
}
