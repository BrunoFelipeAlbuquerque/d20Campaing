package character

import (
	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	"d20campaigngenerator/internal/domain/rpg/character/skill"
)

const coreClassSkillCheckBonus = 3

type characterSkillCheckTotal struct {
	skillID                  skill.SkillID
	abilityScoreID           ability.AbilityScoreID
	abilityModifier          int
	ranks                    int
	classSkill               bool
	classSkillBonus          int
	armorCheckPenaltyApplies bool
	total                    int
}
type CharacterSkillCheckTotal = characterSkillCheckTotal

type characterSkillCheckTotalFacts struct {
	valid            bool
	abilityModifiers map[ability.AbilityScoreID]int
	skillRanks       map[skill.SkillID]int
	classSkillIDs    []skill.SkillID
}
type CharacterSkillCheckTotalFacts = characterSkillCheckTotalFacts

func NewCharacterSkillCheckTotalFacts(
	classLevels []CharacterClassLevel,
	abilityScores []CharacterAbilityScore,
	rankAllocations CharacterSkillRankAllocationFacts,
) (CharacterSkillCheckTotalFacts, bool) {
	levelFacts, ok := NewCharacterLevelFacts(classLevels)
	if !ok {
		return characterSkillCheckTotalFacts{}, false
	}

	if !isValidCharacterSkillRankAllocationFacts(rankAllocations) ||
		rankAllocations.rankCap != levelFacts.GetTotalCharacterLevel() {
		return characterSkillCheckTotalFacts{}, false
	}

	abilityModifiers, ok := characterAbilityModifierMap(abilityScores)
	if !ok {
		return characterSkillCheckTotalFacts{}, false
	}

	skillRanks, ok := characterSkillRankMap(rankAllocations)
	if !ok {
		return characterSkillCheckTotalFacts{}, false
	}

	classSkillIDs, ok := characterClassSkillIDs(levelFacts.GetClassLevels())
	if !ok {
		return characterSkillCheckTotalFacts{}, false
	}

	return characterSkillCheckTotalFacts{
		valid:            true,
		abilityModifiers: abilityModifiers,
		skillRanks:       skillRanks,
		classSkillIDs:    classSkillIDs,
	}, true
}

func (f characterSkillCheckTotalFacts) GetSkillCheckTotal(
	id skill.SkillID,
) (CharacterSkillCheckTotal, bool) {
	if !f.valid {
		return characterSkillCheckTotal{}, false
	}

	selectedSkill, ok := concreteCharacterSkillForRankAllocation(id)
	if !ok {
		return characterSkillCheckTotal{}, false
	}

	ranks := f.skillRanks[selectedSkill.GetID()]
	if selectedSkill.IsTrainedOnly() && ranks == 0 {
		return characterSkillCheckTotal{}, false
	}

	abilityScoreID := selectedSkill.GetAbilityScoreID()
	abilityModifier, ok := f.abilityModifiers[abilityScoreID]
	if !ok {
		return characterSkillCheckTotal{}, false
	}

	classSkill := characterHasClassSkill(f.classSkillIDs, selectedSkill)
	classSkillBonus := 0
	if classSkill && ranks > 0 {
		classSkillBonus = coreClassSkillCheckBonus
	}

	return characterSkillCheckTotal{
		skillID:                  selectedSkill.GetID(),
		abilityScoreID:           abilityScoreID,
		abilityModifier:          abilityModifier,
		ranks:                    ranks,
		classSkill:               classSkill,
		classSkillBonus:          classSkillBonus,
		armorCheckPenaltyApplies: selectedSkill.AppliesArmorCheckPenalty(),
		total:                    ranks + abilityModifier + classSkillBonus,
	}, true
}

func (t characterSkillCheckTotal) GetSkillID() skill.SkillID {
	return t.skillID
}

func (t characterSkillCheckTotal) GetAbilityScoreID() ability.AbilityScoreID {
	return t.abilityScoreID
}

func (t characterSkillCheckTotal) GetAbilityModifier() int {
	return t.abilityModifier
}

func (t characterSkillCheckTotal) GetRanks() int {
	return t.ranks
}

func (t characterSkillCheckTotal) IsClassSkill() bool {
	return t.classSkill
}

func (t characterSkillCheckTotal) GetClassSkillBonus() int {
	return t.classSkillBonus
}

func (t characterSkillCheckTotal) AppliesArmorCheckPenalty() bool {
	return t.armorCheckPenaltyApplies
}

func (t characterSkillCheckTotal) GetTotal() int {
	return t.total
}

func characterAbilityModifierMap(
	abilityScores []CharacterAbilityScore,
) (map[ability.AbilityScoreID]int, bool) {
	scoreMap, ok := buildCompleteCharacterAbilityScoreMap(abilityScores)
	if !ok {
		return nil, false
	}

	result := make(map[ability.AbilityScoreID]int, len(scoreMap))
	for id, scoreValue := range scoreMap {
		value, ok := ability.NewAbilityScoreValue(scoreValue, true)
		if !ok {
			return nil, false
		}

		score, ok := ability.NewAbilityScore(id, value)
		if !ok {
			return nil, false
		}

		modifier, ok := score.GetModifier()
		if !ok {
			return nil, false
		}

		result[id] = modifier
	}

	return result, true
}

func characterSkillRankMap(
	facts CharacterSkillRankAllocationFacts,
) (map[skill.SkillID]int, bool) {
	if !isValidCharacterSkillRankAllocationFacts(facts) {
		return nil, false
	}

	result := make(map[skill.SkillID]int, len(facts.allocations))
	for _, allocation := range facts.allocations {
		selectedSkill, ok := allocation.GetSkill()
		if !ok {
			return nil, false
		}

		result[selectedSkill.GetID()] = allocation.GetRanks()
	}

	return result, true
}

func characterClassSkillIDs(
	classLevels []CharacterClassLevel,
) ([]skill.SkillID, bool) {
	result := make([]skill.SkillID, 0)
	seen := make(map[skill.SkillID]struct{})

	for _, classLevel := range classLevels {
		class, ok := characterclass.GetClassByID(classLevel.GetClassID())
		if !ok {
			return nil, false
		}

		for _, classSkillID := range class.GetClassSkills() {
			if _, ok := seen[classSkillID]; ok {
				continue
			}

			seen[classSkillID] = struct{}{}
			result = append(result, classSkillID)
		}
	}

	return result, true
}

func characterHasClassSkill(
	classSkillIDs []skill.SkillID,
	selectedSkill skill.Skill,
) bool {
	for _, classSkillID := range classSkillIDs {
		if classSkillID == selectedSkill.GetID() {
			return true
		}

		if selectedSkill.IsGrouped() && classSkillID == selectedSkill.GetFamilyID() {
			return true
		}
	}

	return false
}
