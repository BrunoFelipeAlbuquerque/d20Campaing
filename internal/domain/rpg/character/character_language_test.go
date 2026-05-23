package character

import (
	"testing"

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
