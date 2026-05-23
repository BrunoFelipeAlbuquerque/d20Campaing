package race

import (
	"strings"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterlanguage "d20campaigngenerator/internal/domain/rpg/character/language"
)

type raceID string
type RaceID = raceID

type LanguageID = characterlanguage.LanguageID

type racialFeatureID string
type RacialFeatureID = racialFeatureID

type bonusLanguageChoice struct {
	languageIDs  []LanguageID
	anyNonSecret bool
}
type BonusLanguageChoice = bonusLanguageChoice

type abilityScoreModifier struct {
	scoreID  ability.AbilityScoreID
	modifier int
}
type AbilityScoreModifier = abilityScoreModifier

type race struct {
	id                             raceID
	size                           ability.Size
	baseSpeed                      int
	abilityScoreModifiers          []abilityScoreModifier
	selectableAbilityScoreModifier int
	automaticLanguages             []LanguageID
	bonusLanguageChoice            bonusLanguageChoice
	racialFeatures                 []racialFeatureID
}
type Race = race

func NewAbilityScoreModifier(scoreID ability.AbilityScoreID, modifier int) (AbilityScoreModifier, bool) {
	if scoreID.GetName() == "" || modifier == 0 {
		return abilityScoreModifier{}, false
	}

	return abilityScoreModifier{
		scoreID:  scoreID,
		modifier: modifier,
	}, true
}

func NewBonusLanguageChoice(
	languageIDs []LanguageID,
	anyNonSecret bool,
) (BonusLanguageChoice, bool) {
	dedupedLanguages, ok := dedupeLanguageIDs(languageIDs)
	if !ok {
		return bonusLanguageChoice{}, false
	}

	if len(dedupedLanguages) == 0 && !anyNonSecret {
		return bonusLanguageChoice{}, false
	}

	return bonusLanguageChoice{
		languageIDs:  dedupedLanguages,
		anyNonSecret: anyNonSecret,
	}, true
}

func NewRace(
	id RaceID,
	size ability.Size,
	baseSpeed int,
	abilityScoreModifiers []AbilityScoreModifier,
	selectableAbilityScoreModifier int,
	automaticLanguages []LanguageID,
	bonusLanguageChoice BonusLanguageChoice,
	racialFeatures []RacialFeatureID,
) (Race, bool) {
	if !isValidRaceID(id) || !isValidSize(size) || baseSpeed <= 0 || selectableAbilityScoreModifier < 0 {
		return race{}, false
	}

	dedupedModifiers, ok := dedupeAbilityScoreModifiers(abilityScoreModifiers)
	if !ok {
		return race{}, false
	}

	if selectableAbilityScoreModifier != 0 {
		if selectableAbilityScoreModifier != 2 || len(dedupedModifiers) != 0 {
			return race{}, false
		}
	}

	dedupedAutomaticLanguages, ok := dedupeLanguageIDs(automaticLanguages)
	if !ok {
		return race{}, false
	}

	if !isValidBonusLanguageChoice(bonusLanguageChoice) {
		return race{}, false
	}

	dedupedFeatures, ok := dedupeRacialFeatureIDs(racialFeatures)
	if !ok {
		return race{}, false
	}

	return race{
		id:                             id,
		size:                           size,
		baseSpeed:                      baseSpeed,
		abilityScoreModifiers:          dedupedModifiers,
		selectableAbilityScoreModifier: selectableAbilityScoreModifier,
		automaticLanguages:             dedupedAutomaticLanguages,
		bonusLanguageChoice:            cloneBonusLanguageChoice(bonusLanguageChoice),
		racialFeatures:                 dedupedFeatures,
	}, true
}

func (m abilityScoreModifier) GetScoreID() ability.AbilityScoreID {
	return m.scoreID
}

func (m abilityScoreModifier) GetModifier() int {
	return m.modifier
}

func (r race) GetID() RaceID {
	return r.id
}

func (r race) GetSize() ability.Size {
	return r.size
}

func (r race) GetBaseSpeed() int {
	return r.baseSpeed
}

func (r race) GetAbilityScoreModifiers() []AbilityScoreModifier {
	return append([]AbilityScoreModifier(nil), r.abilityScoreModifiers...)
}

func (r race) GetSelectableAbilityScoreModifier() (int, bool) {
	if r.selectableAbilityScoreModifier == 0 {
		return 0, false
	}

	return r.selectableAbilityScoreModifier, true
}

func (r race) GetAutomaticLanguages() []LanguageID {
	return append([]LanguageID(nil), r.automaticLanguages...)
}

func (r race) GetBonusLanguageChoice() (BonusLanguageChoice, bool) {
	if !hasBonusLanguageChoice(r.bonusLanguageChoice) {
		return bonusLanguageChoice{}, false
	}

	return cloneBonusLanguageChoice(r.bonusLanguageChoice), true
}

func (r race) GetRacialFeatures() []RacialFeatureID {
	return append([]RacialFeatureID(nil), r.racialFeatures...)
}

func (c bonusLanguageChoice) GetLanguageIDs() []LanguageID {
	return append([]LanguageID(nil), c.languageIDs...)
}

func (c bonusLanguageChoice) AllowsAnyNonSecret() bool {
	return c.anyNonSecret
}

func (c bonusLanguageChoice) AllowsLanguageID(id LanguageID) bool {
	if c.anyNonSecret && characterlanguage.IsNonSecretLanguageID(id) {
		return true
	}

	for _, languageID := range c.languageIDs {
		if languageID == id {
			return true
		}
	}

	return false
}

func (r race) HasFeature(featureID RacialFeatureID) bool {
	for _, current := range r.racialFeatures {
		if current == featureID {
			return true
		}
	}

	return false
}

func (r race) HasRacialFeature(featureID RacialFeatureID) bool {
	return r.HasFeature(featureID)
}

func isValidRaceID(value RaceID) bool {
	id := string(value)
	return id != "" && strings.TrimSpace(id) == id
}

func isValidSize(value ability.Size) bool {
	_, ok := value.GetModifier()
	return ok
}

func isValidLanguageID(value LanguageID) bool {
	return characterlanguage.IsKnownLanguageID(value)
}

func isValidRacialFeatureID(value RacialFeatureID) bool {
	_, ok := validRacialFeatureIDs[value]
	return ok
}

func isValidBonusLanguageChoice(value BonusLanguageChoice) bool {
	if !value.anyNonSecret && len(value.languageIDs) == 0 {
		return true
	}

	dedupedLanguages, ok := dedupeLanguageIDs(value.languageIDs)
	if !ok {
		return false
	}

	return len(dedupedLanguages) == len(value.languageIDs)
}

func hasBonusLanguageChoice(value BonusLanguageChoice) bool {
	return value.anyNonSecret || len(value.languageIDs) > 0
}

func cloneBonusLanguageChoice(value BonusLanguageChoice) BonusLanguageChoice {
	return bonusLanguageChoice{
		languageIDs:  append([]LanguageID(nil), value.languageIDs...),
		anyNonSecret: value.anyNonSecret,
	}
}

func dedupeAbilityScoreModifiers(modifiers []AbilityScoreModifier) ([]AbilityScoreModifier, bool) {
	if len(modifiers) == 0 {
		return nil, true
	}

	seen := make(map[ability.AbilityScoreID]struct{}, len(modifiers))
	deduped := make([]AbilityScoreModifier, 0, len(modifiers))

	for _, modifier := range modifiers {
		if modifier.scoreID.GetName() == "" || modifier.modifier == 0 {
			return nil, false
		}

		if _, ok := seen[modifier.scoreID]; ok {
			continue
		}

		seen[modifier.scoreID] = struct{}{}
		deduped = append(deduped, modifier)
	}

	return deduped, true
}

func dedupeLanguageIDs(languageIDs []LanguageID) ([]LanguageID, bool) {
	if len(languageIDs) == 0 {
		return nil, true
	}

	seen := make(map[LanguageID]struct{}, len(languageIDs))
	deduped := make([]LanguageID, 0, len(languageIDs))

	for _, languageID := range languageIDs {
		if !isValidLanguageID(languageID) {
			return nil, false
		}

		if _, ok := seen[languageID]; ok {
			continue
		}

		seen[languageID] = struct{}{}
		deduped = append(deduped, languageID)
	}

	return deduped, true
}

func dedupeRacialFeatureIDs(featureIDs []RacialFeatureID) ([]RacialFeatureID, bool) {
	if len(featureIDs) == 0 {
		return nil, true
	}

	seen := make(map[RacialFeatureID]struct{}, len(featureIDs))
	deduped := make([]RacialFeatureID, 0, len(featureIDs))

	for _, featureID := range featureIDs {
		if !isValidRacialFeatureID(featureID) {
			return nil, false
		}

		if _, ok := seen[featureID]; ok {
			continue
		}

		seen[featureID] = struct{}{}
		deduped = append(deduped, featureID)
	}

	return deduped, true
}
