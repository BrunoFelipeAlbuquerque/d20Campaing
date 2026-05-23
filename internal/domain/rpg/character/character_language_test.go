package character

import (
	"testing"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterlanguage "d20campaigngenerator/internal/domain/rpg/character/language"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
)

func TestNewAutomaticRacialCharacterLanguageFacts_ComposesCoreRaceAutomaticLanguages(t *testing.T) {
	selectedRace := mustNewCharacterRaceForLanguageTest(t, characterrace.DwarfRaceID)

	facts, ok := NewAutomaticRacialCharacterLanguageFacts(selectedRace)
	if !ok {
		t.Fatal("expected automatic racial language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		facts.GetLanguageIDs(),
		[]characterlanguage.LanguageID{
			characterlanguage.CommonLanguageID,
			characterlanguage.DwarvenLanguageID,
		},
	)

	if !facts.HasLanguage(characterlanguage.CommonLanguageID) {
		t.Fatal("expected automatic language facts to include Common")
	}

	if !facts.HasLanguage(characterlanguage.DwarvenLanguageID) {
		t.Fatal("expected automatic language facts to include Dwarven")
	}

	if facts.HasLanguage(characterlanguage.GiantLanguageID) {
		t.Fatal("expected automatic language facts not to include bonus languages")
	}
}

func TestNewAutomaticRacialCharacterLanguageFacts_ComposesMultiLanguageRace(t *testing.T) {
	selectedRace := mustNewCharacterRaceForLanguageTest(t, characterrace.GnomeRaceID)

	facts, ok := NewAutomaticRacialCharacterLanguageFacts(selectedRace)
	if !ok {
		t.Fatal("expected gnome automatic racial language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		facts.GetLanguageIDs(),
		[]characterlanguage.LanguageID{
			characterlanguage.CommonLanguageID,
			characterlanguage.GnomeLanguageID,
			characterlanguage.SylvanLanguageID,
		},
	)
}

func TestCharacterLanguageFacts_ExposesDefensiveLanguageCopies(t *testing.T) {
	selectedRace := mustNewCharacterRaceForLanguageTest(t, characterrace.ElfRaceID)

	facts, ok := NewAutomaticRacialCharacterLanguageFacts(selectedRace)
	if !ok {
		t.Fatal("expected elf automatic racial language facts to compose")
	}

	languages := facts.GetLanguages()
	if len(languages) != 2 {
		t.Fatalf("expected two language facts, got %d", len(languages))
	}

	languages[0] = mustNewCharacterLanguageForLanguageTest(t, characterlanguage.DwarvenLanguageID)

	if !facts.HasLanguage(characterlanguage.CommonLanguageID) {
		t.Fatal("expected copied language mutation not to remove Common")
	}

	if facts.HasLanguage(characterlanguage.DwarvenLanguageID) {
		t.Fatal("expected copied language mutation not to add Dwarven")
	}
}

func TestNewCharacterLanguage_RejectsUnknownLanguage(t *testing.T) {
	if _, ok := NewCharacterLanguage(characterlanguage.LanguageID("common")); ok {
		t.Fatal("expected unknown language casing to be rejected")
	}
}

func TestNewAutomaticRacialCharacterLanguageFacts_RejectsZeroValueOrUnknownRace(t *testing.T) {
	var zeroRace CharacterRace
	if _, ok := NewAutomaticRacialCharacterLanguageFacts(zeroRace); ok {
		t.Fatal("expected zero-value selected race to fail")
	}

	if _, ok := NewAutomaticRacialCharacterLanguageFacts(CharacterRace{id: characterrace.RaceID("android")}); ok {
		t.Fatal("expected unknown selected race to fail")
	}
}

func TestCharacterLanguagesFromIDs_DedupesAndRejectsMalformedLanguages(t *testing.T) {
	languages, ok := characterLanguagesFromIDs([]characterlanguage.LanguageID{
		characterlanguage.CommonLanguageID,
		characterlanguage.ElvenLanguageID,
		characterlanguage.CommonLanguageID,
	})
	if !ok {
		t.Fatal("expected duplicate language ids to dedupe")
	}

	if len(languages) != 2 {
		t.Fatalf("expected two deduped languages, got %d", len(languages))
	}

	if _, ok := characterLanguagesFromIDs([]characterlanguage.LanguageID{characterlanguage.LanguageID("Thieves' Cant")}); ok {
		t.Fatal("expected unknown language id to fail")
	}
}

func TestNewBonusRacialCharacterLanguageFacts_ComposesFixedListRaceFromIntelligence(t *testing.T) {
	selectedRace := mustNewCharacterRaceForLanguageTest(t, characterrace.ElfRaceID)

	facts, ok := NewBonusRacialCharacterLanguageFacts(
		selectedRace,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 14)},
		[]characterlanguage.LanguageID{
			characterlanguage.CelestialLanguageID,
			characterlanguage.DraconicLanguageID,
		},
	)
	if !ok {
		t.Fatal("expected fixed-list bonus language facts to compose")
	}

	assertCharacterLanguageIDs(
		t,
		facts.GetLanguageIDs(),
		[]characterlanguage.LanguageID{
			characterlanguage.CelestialLanguageID,
			characterlanguage.DraconicLanguageID,
		},
	)

	if facts.HasLanguage(characterlanguage.CommonLanguageID) || facts.HasLanguage(characterlanguage.ElvenLanguageID) {
		t.Fatal("expected bonus language facts not to include automatic languages")
	}
}

func TestNewBonusRacialCharacterLanguageFacts_ComposesAnyNonSecretRaceFromIntelligence(t *testing.T) {
	testCases := []struct {
		name       string
		raceID     characterrace.RaceID
		languageID characterlanguage.LanguageID
	}{
		{"human", characterrace.HumanRaceID, characterlanguage.InfernalLanguageID},
		{"half-elf", characterrace.HalfElfRaceID, characterlanguage.AquanLanguageID},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			selectedRace := mustNewCharacterRaceForLanguageTest(t, tc.raceID)

			facts, ok := NewBonusRacialCharacterLanguageFacts(
				selectedRace,
				[]CharacterAbilityScore{mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 12)},
				[]characterlanguage.LanguageID{tc.languageID},
			)
			if !ok {
				t.Fatal("expected any-non-secret bonus language facts to compose")
			}

			assertCharacterLanguageIDs(
				t,
				facts.GetLanguageIDs(),
				[]characterlanguage.LanguageID{tc.languageID},
			)
		})
	}
}

func TestNewBonusRacialCharacterLanguageFacts_AllowsNoChoicesWhenIntelligenceDoesNotGrantBonus(t *testing.T) {
	selectedRace := mustNewCharacterRaceForLanguageTest(t, characterrace.DwarfRaceID)

	facts, ok := NewBonusRacialCharacterLanguageFacts(
		selectedRace,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 9)},
		nil,
	)
	if !ok {
		t.Fatal("expected zero bonus language facts to compose when Intelligence grants none")
	}

	if len(facts.GetLanguageIDs()) != 0 {
		t.Fatalf("expected no bonus languages, got %v", facts.GetLanguageIDs())
	}
}

func TestNewBonusRacialCharacterLanguageFacts_RejectsMissingDuplicateUnknownOverBudgetOrDisallowedChoices(t *testing.T) {
	var zeroRace CharacterRace

	humanRace := mustNewCharacterRaceForLanguageTest(t, characterrace.HumanRaceID)
	dwarfRace := mustNewCharacterRaceForLanguageTest(t, characterrace.DwarfRaceID)
	unknownRace := CharacterRace{id: characterrace.RaceID("android")}
	intelligence12 := []CharacterAbilityScore{mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 12)}
	intelligence14 := []CharacterAbilityScore{mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 14)}

	testCases := []struct {
		name      string
		race      CharacterRace
		scores    []CharacterAbilityScore
		languages []characterlanguage.LanguageID
	}{
		{"zero race", zeroRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.InfernalLanguageID}},
		{"unknown race", unknownRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.InfernalLanguageID}},
		{"missing intelligence", humanRace, nil, []characterlanguage.LanguageID{characterlanguage.InfernalLanguageID}},
		{"missing choice", humanRace, intelligence12, nil},
		{"duplicate choices", humanRace, intelligence14, []characterlanguage.LanguageID{characterlanguage.InfernalLanguageID, characterlanguage.InfernalLanguageID}},
		{"unknown choice", humanRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.LanguageID("Thieves' Cant")}},
		{"over budget", humanRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.InfernalLanguageID, characterlanguage.AquanLanguageID}},
		{"disallowed fixed-list choice", dwarfRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.CelestialLanguageID}},
		{"secret any-non-secret choice", humanRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.DruidicLanguageID}},
		{"automatic language duplicate", humanRace, intelligence12, []characterlanguage.LanguageID{characterlanguage.CommonLanguageID}},
		{"no budget with choice", dwarfRace, []CharacterAbilityScore{mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 8)}, []characterlanguage.LanguageID{characterlanguage.GiantLanguageID}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewBonusRacialCharacterLanguageFacts(tc.race, tc.scores, tc.languages); ok {
				t.Fatal("expected bonus language facts to fail")
			}
		})
	}
}

func TestNewBonusRacialCharacterLanguageFacts_RejectsMalformedAbilityScores(t *testing.T) {
	selectedRace := mustNewCharacterRaceForLanguageTest(t, characterrace.HumanRaceID)

	if _, ok := NewBonusRacialCharacterLanguageFacts(
		selectedRace,
		[]CharacterAbilityScore{{id: ability.IntelligenceScore, score: -1}},
		[]characterlanguage.LanguageID{characterlanguage.InfernalLanguageID},
	); ok {
		t.Fatal("expected malformed Intelligence score to fail")
	}

	intelligence := mustNewCharacterAbilityScoreForLanguageTest(t, ability.IntelligenceScore, 12)
	if _, ok := NewBonusRacialCharacterLanguageFacts(
		selectedRace,
		[]CharacterAbilityScore{intelligence, intelligence},
		[]characterlanguage.LanguageID{characterlanguage.InfernalLanguageID},
	); ok {
		t.Fatal("expected duplicate Intelligence scores to fail")
	}
}

func assertCharacterLanguageIDs(
	t *testing.T,
	actual []characterlanguage.LanguageID,
	expected []characterlanguage.LanguageID,
) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected %d language ids, got %d: %v", len(expected), len(actual), actual)
	}

	for i, expectedID := range expected {
		if actual[i] != expectedID {
			t.Fatalf("expected language id at %d to be %q, got %q", i, expectedID, actual[i])
		}
	}
}

func mustNewCharacterRaceForLanguageTest(t *testing.T, id characterrace.RaceID) CharacterRace {
	t.Helper()

	selectedRace, ok := NewCharacterRace(id)
	if !ok {
		t.Fatalf("expected selected race %q to compose", id)
	}

	return selectedRace
}

func mustNewCharacterAbilityScoreForLanguageTest(
	t *testing.T,
	id ability.AbilityScoreID,
	score int,
) CharacterAbilityScore {
	t.Helper()

	abilityScore, ok := NewCharacterAbilityScore(id, score)
	if !ok {
		t.Fatalf("expected ability score %q %d to compose", id, score)
	}

	return abilityScore
}

func mustNewCharacterLanguageForLanguageTest(
	t *testing.T,
	id characterlanguage.LanguageID,
) CharacterLanguage {
	t.Helper()

	language, ok := NewCharacterLanguage(id)
	if !ok {
		t.Fatalf("expected language %q to compose", id)
	}

	return language
}
