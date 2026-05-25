package skill

import (
	"strings"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
)

type skillID string
type SkillID = skillID

const (
	AcrobaticsSkillID     SkillID = "Acrobatics"
	AppraiseSkillID       SkillID = "Appraise"
	BluffSkillID          SkillID = "Bluff"
	ClimbSkillID          SkillID = "Climb"
	CraftSkillID          SkillID = "Craft"
	DiplomacySkillID      SkillID = "Diplomacy"
	DisableDeviceSkillID  SkillID = "Disable Device"
	DisguiseSkillID       SkillID = "Disguise"
	EscapeArtistSkillID   SkillID = "Escape Artist"
	FlySkillID            SkillID = "Fly"
	HandleAnimalSkillID   SkillID = "Handle Animal"
	HealSkillID           SkillID = "Heal"
	IntimidateSkillID     SkillID = "Intimidate"
	KnowledgeSkillID      SkillID = "Knowledge"
	LinguisticsSkillID    SkillID = "Linguistics"
	PerceptionSkillID     SkillID = "Perception"
	PerformSkillID        SkillID = "Perform"
	ProfessionSkillID     SkillID = "Profession"
	RideSkillID           SkillID = "Ride"
	SenseMotiveSkillID    SkillID = "Sense Motive"
	SleightOfHandSkillID  SkillID = "Sleight of Hand"
	SpellcraftSkillID     SkillID = "Spellcraft"
	StealthSkillID        SkillID = "Stealth"
	SurvivalSkillID       SkillID = "Survival"
	SwimSkillID           SkillID = "Swim"
	UseMagicDeviceSkillID SkillID = "Use Magic Device"
)

type skill struct {
	id                       skillID
	familyID                 skillID
	specialization           string
	abilityScoreID           ability.AbilityScoreID
	trainedOnly              bool
	armorCheckPenaltyApplies bool
	grouped                  bool
}
type Skill = skill

func NewSkill(
	id SkillID,
	abilityScoreID ability.AbilityScoreID,
	trainedOnly bool,
	armorCheckPenaltyApplies bool,
	grouped bool,
) (Skill, bool) {
	familyID, specialization, ok := parseSkillID(id)
	if !ok || grouped != isGroupedSkillID(familyID) {
		return skill{}, false
	}

	expectedAbilityScoreID, ok := coreSkillAbilityScoreID(familyID)
	if !ok || abilityScoreID != expectedAbilityScoreID {
		return skill{}, false
	}

	return skill{
		id:                       id,
		familyID:                 familyID,
		specialization:           specialization,
		abilityScoreID:           abilityScoreID,
		trainedOnly:              trainedOnly,
		armorCheckPenaltyApplies: armorCheckPenaltyApplies,
		grouped:                  grouped,
	}, true
}

func (id skillID) GetName() string {
	familyID, specialization, ok := parseSkillID(SkillID(id))
	if !ok {
		return ""
	}

	if specialization != "" {
		return string(familyID) + " (" + specialization + ")"
	}

	return familyID.GetBaseName()
}

func (id skillID) GetBaseName() string {
	name := strings.TrimSpace(string(id))
	if name == "" || name != string(id) {
		return ""
	}

	switch SkillID(name) {
	case AcrobaticsSkillID,
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
		UseMagicDeviceSkillID:
		return name
	default:
		return ""
	}
}

func (s skill) GetID() SkillID {
	return s.id
}

func (s skill) GetFamilyID() SkillID {
	return s.familyID
}

func (s skill) GetSpecialization() (string, bool) {
	if s.specialization == "" {
		return "", false
	}

	return s.specialization, true
}

func (s skill) GetAbilityScoreID() ability.AbilityScoreID {
	return s.abilityScoreID
}

func (s skill) IsTrainedOnly() bool {
	return s.trainedOnly
}

func (s skill) AppliesArmorCheckPenalty() bool {
	return s.armorCheckPenaltyApplies
}

func (s skill) IsGrouped() bool {
	return s.grouped
}

func isValidSkillID(id SkillID) bool {
	_, _, ok := parseSkillID(id)
	return ok
}

func IsValidSkillID(id SkillID) bool {
	return isValidSkillID(id)
}

func isGroupedSkillID(id SkillID) bool {
	switch id {
	case CraftSkillID, KnowledgeSkillID, PerformSkillID, ProfessionSkillID:
		return true
	default:
		return false
	}
}

func parseSkillID(id SkillID) (SkillID, string, bool) {
	if baseName := id.GetBaseName(); baseName != "" {
		return SkillID(baseName), "", true
	}

	value := string(id)
	openIndex := strings.Index(value, " (")
	if openIndex <= 0 || !strings.HasSuffix(value, ")") {
		return "", "", false
	}

	familyID := SkillID(value[:openIndex])
	if !isGroupedSkillID(familyID) {
		return "", "", false
	}

	specialization := value[openIndex+2 : len(value)-1]
	if specialization == "" || strings.TrimSpace(specialization) != specialization {
		return "", "", false
	}

	if strings.ContainsRune(specialization, '(') || strings.ContainsRune(specialization, ')') {
		return "", "", false
	}

	return familyID, specialization, true
}
