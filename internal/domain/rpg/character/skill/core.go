package skill

import ability "d20campaigngenerator/internal/domain/rpg/character/ability"

var coreSkills = mustBuildCoreSkills()

var coreSkillOrder = []SkillID{
	AcrobaticsSkillID,
	AppraiseSkillID,
	BluffSkillID,
	ClimbSkillID,
	CraftSkillID,
	DiplomacySkillID,
	DisableDeviceSkillID,
	DisguiseSkillID,
	EscapeArtistSkillID,
	FlySkillID,
	HandleAnimalSkillID,
	HealSkillID,
	IntimidateSkillID,
	KnowledgeSkillID,
	LinguisticsSkillID,
	PerceptionSkillID,
	PerformSkillID,
	ProfessionSkillID,
	RideSkillID,
	SenseMotiveSkillID,
	SleightOfHandSkillID,
	SpellcraftSkillID,
	StealthSkillID,
	SurvivalSkillID,
	SwimSkillID,
	UseMagicDeviceSkillID,
}

func GetSkillByID(id SkillID) (Skill, bool) {
	value, ok := coreSkills[id]
	if !ok {
		return skill{}, false
	}

	return value, true
}

func GetSkills() []Skill {
	skills := make([]Skill, 0, len(coreSkillOrder))

	for _, id := range coreSkillOrder {
		skills = append(skills, coreSkills[id])
	}

	return skills
}

func NewSpecializedGroupedSkill(id SkillID) (Skill, bool) {
	familyID, specialization, ok := parseSkillID(id)
	if !ok || specialization == "" {
		return skill{}, false
	}

	familySkill, ok := GetSkillByID(familyID)
	if !ok || !familySkill.IsGrouped() {
		return skill{}, false
	}

	return NewSkill(
		id,
		familySkill.GetAbilityScoreID(),
		familySkill.IsTrainedOnly(),
		familySkill.AppliesArmorCheckPenalty(),
		true,
	)
}

func mustBuildCoreSkills() map[SkillID]Skill {
	return map[SkillID]Skill{
		AcrobaticsSkillID:     mustNewSkill(AcrobaticsSkillID, ability.DexterityScore, false, true, false),
		AppraiseSkillID:       mustNewSkill(AppraiseSkillID, ability.IntelligenceScore, false, false, false),
		BluffSkillID:          mustNewSkill(BluffSkillID, ability.CharismaScore, false, false, false),
		ClimbSkillID:          mustNewSkill(ClimbSkillID, ability.StrengthScore, false, true, false),
		CraftSkillID:          mustNewSkill(CraftSkillID, ability.IntelligenceScore, false, false, true),
		DiplomacySkillID:      mustNewSkill(DiplomacySkillID, ability.CharismaScore, false, false, false),
		DisableDeviceSkillID:  mustNewSkill(DisableDeviceSkillID, ability.DexterityScore, true, true, false),
		DisguiseSkillID:       mustNewSkill(DisguiseSkillID, ability.CharismaScore, false, false, false),
		EscapeArtistSkillID:   mustNewSkill(EscapeArtistSkillID, ability.DexterityScore, false, true, false),
		FlySkillID:            mustNewSkill(FlySkillID, ability.DexterityScore, false, true, false),
		HandleAnimalSkillID:   mustNewSkill(HandleAnimalSkillID, ability.CharismaScore, true, false, false),
		HealSkillID:           mustNewSkill(HealSkillID, ability.WisdomScore, false, false, false),
		IntimidateSkillID:     mustNewSkill(IntimidateSkillID, ability.CharismaScore, false, false, false),
		KnowledgeSkillID:      mustNewSkill(KnowledgeSkillID, ability.IntelligenceScore, true, false, true),
		LinguisticsSkillID:    mustNewSkill(LinguisticsSkillID, ability.IntelligenceScore, true, false, false),
		PerceptionSkillID:     mustNewSkill(PerceptionSkillID, ability.WisdomScore, false, false, false),
		PerformSkillID:        mustNewSkill(PerformSkillID, ability.CharismaScore, false, false, true),
		ProfessionSkillID:     mustNewSkill(ProfessionSkillID, ability.WisdomScore, true, false, true),
		RideSkillID:           mustNewSkill(RideSkillID, ability.DexterityScore, false, true, false),
		SenseMotiveSkillID:    mustNewSkill(SenseMotiveSkillID, ability.WisdomScore, false, false, false),
		SleightOfHandSkillID:  mustNewSkill(SleightOfHandSkillID, ability.DexterityScore, true, true, false),
		SpellcraftSkillID:     mustNewSkill(SpellcraftSkillID, ability.IntelligenceScore, true, false, false),
		StealthSkillID:        mustNewSkill(StealthSkillID, ability.DexterityScore, false, true, false),
		SurvivalSkillID:       mustNewSkill(SurvivalSkillID, ability.WisdomScore, false, false, false),
		SwimSkillID:           mustNewSkill(SwimSkillID, ability.StrengthScore, false, true, false),
		UseMagicDeviceSkillID: mustNewSkill(UseMagicDeviceSkillID, ability.CharismaScore, true, false, false),
	}
}

func mustNewSkill(
	id SkillID,
	abilityScoreID ability.AbilityScoreID,
	trainedOnly bool,
	armorCheckPenaltyApplies bool,
	grouped bool,
) Skill {
	skill, ok := NewSkill(id, abilityScoreID, trainedOnly, armorCheckPenaltyApplies, grouped)
	if !ok {
		panic("invalid core skill seed")
	}

	return skill
}

func coreSkillAbilityScoreID(id SkillID) (ability.AbilityScoreID, bool) {
	switch id {
	case AcrobaticsSkillID,
		DisableDeviceSkillID,
		EscapeArtistSkillID,
		FlySkillID,
		RideSkillID,
		SleightOfHandSkillID,
		StealthSkillID:
		return ability.DexterityScore, true
	case AppraiseSkillID,
		CraftSkillID,
		KnowledgeSkillID,
		LinguisticsSkillID,
		SpellcraftSkillID:
		return ability.IntelligenceScore, true
	case BluffSkillID,
		DiplomacySkillID,
		DisguiseSkillID,
		HandleAnimalSkillID,
		IntimidateSkillID,
		PerformSkillID,
		UseMagicDeviceSkillID:
		return ability.CharismaScore, true
	case ClimbSkillID,
		SwimSkillID:
		return ability.StrengthScore, true
	case HealSkillID,
		PerceptionSkillID,
		ProfessionSkillID,
		SenseMotiveSkillID,
		SurvivalSkillID:
		return ability.WisdomScore, true
	default:
		return "", false
	}
}
