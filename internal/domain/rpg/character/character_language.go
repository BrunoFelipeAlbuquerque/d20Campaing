package character

import (
	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterlanguage "d20campaigngenerator/internal/domain/rpg/character/language"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
)

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

func NewBonusRacialCharacterLanguageFacts(
	selectedRace CharacterRace,
	abilityScores []CharacterAbilityScore,
	selectedLanguageIDs []characterlanguage.LanguageID,
) (CharacterLanguageFacts, bool) {
	race, ok := selectedRace.GetRace()
	if !ok {
		return characterLanguageFacts{}, false
	}

	bonusChoice, ok := race.GetBonusLanguageChoice()
	if !ok {
		return characterLanguageFacts{}, false
	}

	maximumChoices, ok := characterBonusLanguageChoiceLimit(abilityScores)
	if !ok || len(selectedLanguageIDs) != maximumChoices {
		return characterLanguageFacts{}, false
	}

	automaticLanguageIDs, ok := characterLanguageIDSet(race.GetAutomaticLanguages())
	if !ok {
		return characterLanguageFacts{}, false
	}

	languages, ok := bonusCharacterLanguagesFromIDs(
		selectedLanguageIDs,
		bonusChoice,
		automaticLanguageIDs,
	)
	if !ok {
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

func bonusCharacterLanguagesFromIDs(
	languageIDs []characterlanguage.LanguageID,
	bonusChoice characterrace.BonusLanguageChoice,
	automaticLanguageIDs map[characterlanguage.LanguageID]struct{},
) ([]CharacterLanguage, bool) {
	if len(languageIDs) == 0 {
		return nil, true
	}

	seen := make(map[characterlanguage.LanguageID]struct{}, len(languageIDs))
	languages := make([]CharacterLanguage, 0, len(languageIDs))

	for _, id := range languageIDs {
		if _, ok := seen[id]; ok {
			return nil, false
		}

		if _, ok := automaticLanguageIDs[id]; ok {
			return nil, false
		}

		if !bonusChoice.AllowsLanguageID(id) {
			return nil, false
		}

		language, ok := NewCharacterLanguage(id)
		if !ok {
			return nil, false
		}

		seen[id] = struct{}{}
		languages = append(languages, language)
	}

	return languages, true
}

func characterBonusLanguageChoiceLimit(
	abilityScores []CharacterAbilityScore,
) (int, bool) {
	scoreMap, ok := buildCharacterAbilityScoreMap(abilityScores)
	if !ok {
		return 0, false
	}

	intelligence, ok := scoreMap[ability.IntelligenceScore]
	if !ok {
		return 0, false
	}

	value, ok := ability.NewAbilityScoreValue(intelligence, true)
	if !ok {
		return 0, false
	}

	score, ok := ability.NewAbilityScore(ability.IntelligenceScore, value)
	if !ok {
		return 0, false
	}

	modifier, ok := score.GetModifier()
	if !ok {
		return 0, false
	}

	if modifier < 0 {
		return 0, true
	}

	return modifier, true
}

func characterLanguageIDSet(
	languageIDs []characterlanguage.LanguageID,
) (map[characterlanguage.LanguageID]struct{}, bool) {
	result := make(map[characterlanguage.LanguageID]struct{}, len(languageIDs))

	for _, id := range languageIDs {
		if _, ok := NewCharacterLanguage(id); !ok {
			return nil, false
		}

		if _, ok := result[id]; ok {
			return nil, false
		}

		result[id] = struct{}{}
	}

	return result, true
}
