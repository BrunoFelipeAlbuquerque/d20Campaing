package language

import "testing"

func TestCoreLanguages_SeedsCoreCommonLanguages(t *testing.T) {
	testCases := []struct {
		id     LanguageID
		secret bool
	}{
		{AbyssalLanguageID, false},
		{AkloLanguageID, false},
		{AquanLanguageID, false},
		{AuranLanguageID, false},
		{CelestialLanguageID, false},
		{CommonLanguageID, false},
		{DraconicLanguageID, false},
		{DruidicLanguageID, true},
		{DwarvenLanguageID, false},
		{ElvenLanguageID, false},
		{GiantLanguageID, false},
		{GnomeLanguageID, false},
		{GoblinLanguageID, false},
		{GnollLanguageID, false},
		{HalflingLanguageID, false},
		{IgnanLanguageID, false},
		{InfernalLanguageID, false},
		{OrcLanguageID, false},
		{SylvanLanguageID, false},
		{TerranLanguageID, false},
		{UndercommonLanguageID, false},
	}

	if len(coreLanguages) != len(testCases) {
		t.Fatalf("expected %d core languages, got %d", len(testCases), len(coreLanguages))
	}

	for _, tc := range testCases {
		language, ok := GetLanguageByID(tc.id)
		if !ok {
			t.Fatalf("expected language %q to resolve", tc.id)
		}

		if language.GetID() != tc.id {
			t.Fatalf("expected language id %q, got %q", tc.id, language.GetID())
		}

		if language.IsSecret() != tc.secret {
			t.Fatalf("expected language %q secret=%t, got %t", tc.id, tc.secret, language.IsSecret())
		}
	}
}

func TestCoreLanguages_QueryHelpersDistinguishKnownSecretAndNonSecretLanguages(t *testing.T) {
	if !IsKnownLanguageID(InfernalLanguageID) {
		t.Fatal("expected Infernal to be a known language")
	}

	if !IsNonSecretLanguageID(InfernalLanguageID) {
		t.Fatal("expected Infernal to be a known non-secret language")
	}

	if IsSecretLanguageID(InfernalLanguageID) {
		t.Fatal("expected Infernal not to be secret")
	}

	if !IsKnownLanguageID(DruidicLanguageID) {
		t.Fatal("expected Druidic to be a known language")
	}

	if !IsSecretLanguageID(DruidicLanguageID) {
		t.Fatal("expected Druidic to be secret")
	}

	if IsNonSecretLanguageID(DruidicLanguageID) {
		t.Fatal("expected Druidic not to be a non-secret language")
	}

	if IsKnownLanguageID(LanguageID("common")) ||
		IsSecretLanguageID(LanguageID("common")) ||
		IsNonSecretLanguageID(LanguageID("common")) {
		t.Fatal("expected unknown language casing not to resolve")
	}
}

func TestGetLanguageByID_RejectsUnknownLanguage(t *testing.T) {
	if _, ok := GetLanguageByID(LanguageID("Thieves' Cant")); ok {
		t.Fatal("expected unknown language lookup to fail")
	}
}

func TestGetLanguages_ReturnsCoreLanguagesInOrder(t *testing.T) {
	languages := GetLanguages()
	if len(languages) != len(coreLanguageOrder) {
		t.Fatalf("expected %d queried languages, got %d", len(coreLanguageOrder), len(languages))
	}

	for i, expectedID := range coreLanguageOrder {
		if languages[i].GetID() != expectedID {
			t.Fatalf("expected language at index %d to be %q, got %q", i, expectedID, languages[i].GetID())
		}
	}
}

func TestGetLanguages_ReturnsDetachedCopies(t *testing.T) {
	first := GetLanguages()
	second := GetLanguages()

	first[0].id = "Changed"
	first[0].secret = true

	if second[0].GetID() != AbyssalLanguageID {
		t.Fatalf("expected stored language id to remain %q, got %q", AbyssalLanguageID, second[0].GetID())
	}

	if second[0].IsSecret() {
		t.Fatal("expected stored Abyssal language not to become secret")
	}
}
