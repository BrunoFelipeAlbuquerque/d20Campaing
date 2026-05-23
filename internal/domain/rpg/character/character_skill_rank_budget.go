package character

import (
	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
)

const minimumCoreSkillRanksPerLevel = 1

type characterSkillRankBudgetFacts struct {
	totalSkillRanks      int
	classSkillRanks      int
	racialBonusRanks     int
	intelligenceModifier int
	totalCharacterLevel  int
}
type CharacterSkillRankBudgetFacts = characterSkillRankBudgetFacts

func NewCharacterSkillRankBudgetFacts(
	classLevels []CharacterClassLevel,
	abilityScores []CharacterAbilityScore,
	selectedRace CharacterRace,
) (CharacterSkillRankBudgetFacts, bool) {
	levelFacts, ok := NewCharacterLevelFacts(classLevels)
	if !ok {
		return characterSkillRankBudgetFacts{}, false
	}

	intelligenceModifier, ok := characterIntelligenceModifierForSkillRanks(abilityScores)
	if !ok {
		return characterSkillRankBudgetFacts{}, false
	}

	racialBonusRanks, ok := characterRacialSkillRankBonus(selectedRace, levelFacts.GetTotalCharacterLevel())
	if !ok {
		return characterSkillRankBudgetFacts{}, false
	}

	classSkillRanks, ok := characterClassSkillRankBudget(
		levelFacts.GetClassLevels(),
		intelligenceModifier,
	)
	if !ok {
		return characterSkillRankBudgetFacts{}, false
	}

	totalSkillRanks := classSkillRanks + racialBonusRanks
	if totalSkillRanks < 0 {
		return characterSkillRankBudgetFacts{}, false
	}

	return characterSkillRankBudgetFacts{
		totalSkillRanks:      totalSkillRanks,
		classSkillRanks:      classSkillRanks,
		racialBonusRanks:     racialBonusRanks,
		intelligenceModifier: intelligenceModifier,
		totalCharacterLevel:  levelFacts.GetTotalCharacterLevel(),
	}, true
}

func (f characterSkillRankBudgetFacts) GetTotalSkillRanks() int {
	return f.totalSkillRanks
}

func (f characterSkillRankBudgetFacts) GetClassSkillRanks() int {
	return f.classSkillRanks
}

func (f characterSkillRankBudgetFacts) GetRacialBonusRanks() int {
	return f.racialBonusRanks
}

func (f characterSkillRankBudgetFacts) GetIntelligenceModifier() int {
	return f.intelligenceModifier
}

func (f characterSkillRankBudgetFacts) GetTotalCharacterLevel() int {
	return f.totalCharacterLevel
}

func characterClassSkillRankBudget(
	classLevels []CharacterClassLevel,
	intelligenceModifier int,
) (int, bool) {
	total := 0

	for _, classLevel := range classLevels {
		if _, ok := NewCharacterClassLevel(classLevel.GetClassID(), classLevel.GetLevel()); !ok {
			return 0, false
		}

		class, ok := characterclass.GetClassByID(classLevel.GetClassID())
		if !ok {
			return 0, false
		}

		ranksPerLevel := class.GetSkillRanksPerLevel() + intelligenceModifier
		if ranksPerLevel < minimumCoreSkillRanksPerLevel {
			ranksPerLevel = minimumCoreSkillRanksPerLevel
		}

		total += ranksPerLevel * classLevel.GetLevel()
	}

	return total, total >= 0
}

func characterIntelligenceModifierForSkillRanks(
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

	return score.GetModifier()
}

func characterRacialSkillRankBonus(
	selectedRace CharacterRace,
	totalCharacterLevel int,
) (int, bool) {
	if totalCharacterLevel <= 0 {
		return 0, false
	}

	race, ok := selectedRace.GetRace()
	if !ok {
		return 0, false
	}

	if race.HasFeature(characterrace.SkilledFeatureID) {
		return totalCharacterLevel, true
	}

	return 0, true
}
