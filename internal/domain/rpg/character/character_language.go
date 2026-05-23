package character

import characterlanguage "d20campaigngenerator/internal/domain/rpg/character/language"

type characterLanguage struct {
	id characterlanguage.LanguageID
}
type CharacterLanguage = characterLanguage

type characterLanguageFacts struct {
	languages []CharacterLanguage
}
type CharacterLanguageFacts = characterLanguageFacts

func NewCharacterLanguage(id characterlanguage.LanguageID) (CharacterLanguage, bool) {
	if !characterlanguage.IsKnownLanguageID(id) {
		return characterLanguage{}, false
	}

	return characterLanguage{id: id}, true
}

func NewAutomaticRacialCharacterLanguageFacts(
	selectedRace CharacterRace,
) (CharacterLanguageFacts, bool) {
	race, ok := selectedRace.GetRace()
	if !ok {
		return characterLanguageFacts{}, false
	}

	languages, ok := characterLanguagesFromIDs(race.GetAutomaticLanguages())
	if !ok || len(languages) == 0 {
		return characterLanguageFacts{}, false
	}

	return characterLanguageFacts{languages: languages}, true
}

func (l characterLanguage) GetLanguageID() characterlanguage.LanguageID {
	return l.id
}

func (f characterLanguageFacts) GetLanguages() []CharacterLanguage {
	return append([]CharacterLanguage(nil), f.languages...)
}

func (f characterLanguageFacts) GetLanguageIDs() []characterlanguage.LanguageID {
	languageIDs := make([]characterlanguage.LanguageID, 0, len(f.languages))

	for _, language := range f.languages {
		languageIDs = append(languageIDs, language.GetLanguageID())
	}

	return languageIDs
}

func (f characterLanguageFacts) HasLanguage(id characterlanguage.LanguageID) bool {
	if !characterlanguage.IsKnownLanguageID(id) {
		return false
	}

	for _, language := range f.languages {
		if language.GetLanguageID() == id {
			return true
		}
	}

	return false
}

func characterLanguagesFromIDs(
	languageIDs []characterlanguage.LanguageID,
) ([]CharacterLanguage, bool) {
	if len(languageIDs) == 0 {
		return nil, true
	}

	seen := make(map[characterlanguage.LanguageID]struct{}, len(languageIDs))
	deduped := make([]CharacterLanguage, 0, len(languageIDs))

	for _, id := range languageIDs {
		if _, ok := seen[id]; ok {
			continue
		}

		language, ok := NewCharacterLanguage(id)
		if !ok {
			return nil, false
		}

		seen[id] = struct{}{}
		deduped = append(deduped, language)
	}

	return deduped, true
}
