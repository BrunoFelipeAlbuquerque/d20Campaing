package character

import "d20campaigngenerator/internal/domain/rpg/character/skill"

type characterSkillRankAllocation struct {
	skillID skill.SkillID
	ranks   int
}
type CharacterSkillRankAllocation = characterSkillRankAllocation

func NewCharacterSkillRankAllocation(
	id skill.SkillID,
	ranks int,
) (CharacterSkillRankAllocation, bool) {
	if ranks <= 0 {
		return characterSkillRankAllocation{}, false
	}

	if _, ok := concreteCharacterSkillForRankAllocation(id); !ok {
		return characterSkillRankAllocation{}, false
	}

	return characterSkillRankAllocation{
		skillID: id,
		ranks:   ranks,
	}, true
}

func (a characterSkillRankAllocation) GetSkillID() skill.SkillID {
	return a.skillID
}

func (a characterSkillRankAllocation) GetRanks() int {
	return a.ranks
}

func (a characterSkillRankAllocation) GetSkill() (skill.Skill, bool) {
	return concreteCharacterSkillForRankAllocation(a.skillID)
}

func concreteCharacterSkillForRankAllocation(id skill.SkillID) (skill.Skill, bool) {
	if seededSkill, ok := skill.GetSkillByID(id); ok {
		if seededSkill.IsGrouped() {
			return skill.Skill{}, false
		}

		return seededSkill, true
	}

	specializedSkill, ok := skill.NewSkill(id, false, false, true)
	if !ok {
		return skill.Skill{}, false
	}

	if _, ok := specializedSkill.GetSpecialization(); !ok {
		return skill.Skill{}, false
	}

	familySkill, ok := skill.GetSkillByID(specializedSkill.GetFamilyID())
	if !ok || !familySkill.IsGrouped() {
		return skill.Skill{}, false
	}

	return skill.NewSkill(
		id,
		familySkill.IsTrainedOnly(),
		familySkill.AppliesArmorCheckPenalty(),
		true,
	)
}
