package skill

import (
	"testing"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
)

func TestCoreSkills_SeedsTwentySixCoreEntries(t *testing.T) {
	testCases := []struct {
		id                       SkillID
		abilityScoreID           ability.AbilityScoreID
		trainedOnly              bool
		armorCheckPenaltyApplies bool
		grouped                  bool
	}{
		{AcrobaticsSkillID, ability.DexterityScore, false, true, false},
		{AppraiseSkillID, ability.IntelligenceScore, false, false, false},
		{BluffSkillID, ability.CharismaScore, false, false, false},
		{ClimbSkillID, ability.StrengthScore, false, true, false},
		{CraftSkillID, ability.IntelligenceScore, false, false, true},
		{DiplomacySkillID, ability.CharismaScore, false, false, false},
		{DisableDeviceSkillID, ability.DexterityScore, true, true, false},
		{DisguiseSkillID, ability.CharismaScore, false, false, false},
		{EscapeArtistSkillID, ability.DexterityScore, false, true, false},
		{FlySkillID, ability.DexterityScore, false, true, false},
		{HandleAnimalSkillID, ability.CharismaScore, true, false, false},
		{HealSkillID, ability.WisdomScore, false, false, false},
		{IntimidateSkillID, ability.CharismaScore, false, false, false},
		{KnowledgeSkillID, ability.IntelligenceScore, true, false, true},
		{LinguisticsSkillID, ability.IntelligenceScore, true, false, false},
		{PerceptionSkillID, ability.WisdomScore, false, false, false},
		{PerformSkillID, ability.CharismaScore, false, false, true},
		{ProfessionSkillID, ability.WisdomScore, true, false, true},
		{RideSkillID, ability.DexterityScore, false, true, false},
		{SenseMotiveSkillID, ability.WisdomScore, false, false, false},
		{SleightOfHandSkillID, ability.DexterityScore, true, true, false},
		{SpellcraftSkillID, ability.IntelligenceScore, true, false, false},
		{StealthSkillID, ability.DexterityScore, false, true, false},
		{SurvivalSkillID, ability.WisdomScore, false, false, false},
		{SwimSkillID, ability.StrengthScore, false, true, false},
		{UseMagicDeviceSkillID, ability.CharismaScore, true, false, false},
	}

	if len(coreSkills) != len(testCases) {
		t.Fatalf("expected %d core skills, got %d", len(testCases), len(coreSkills))
	}

	for _, tc := range testCases {
		skill, ok := coreSkills[tc.id]
		if !ok {
			t.Fatalf("expected core skill %q to be seeded", tc.id)
		}

		if skill.GetID() != tc.id {
			t.Fatalf("expected skill id %q, got %q", tc.id, skill.GetID())
		}

		if skill.GetAbilityScoreID() != tc.abilityScoreID {
			t.Fatalf("expected skill %q ability score %q, got %q", tc.id, tc.abilityScoreID, skill.GetAbilityScoreID())
		}

		if skill.IsTrainedOnly() != tc.trainedOnly {
			t.Fatalf("expected skill %q trained-only=%t, got %t", tc.id, tc.trainedOnly, skill.IsTrainedOnly())
		}

		if skill.AppliesArmorCheckPenalty() != tc.armorCheckPenaltyApplies {
			t.Fatalf("expected skill %q armor-check-penalty=%t, got %t", tc.id, tc.armorCheckPenaltyApplies, skill.AppliesArmorCheckPenalty())
		}

		if skill.IsGrouped() != tc.grouped {
			t.Fatalf("expected skill %q grouped=%t, got %t", tc.id, tc.grouped, skill.IsGrouped())
		}
	}
}

func TestNewSkill_AcceptsEverySeededCoreSkillID(t *testing.T) {
	for id, seeded := range coreSkills {
		skill, ok := NewSkill(
			id,
			seeded.GetAbilityScoreID(),
			seeded.IsTrainedOnly(),
			seeded.AppliesArmorCheckPenalty(),
			seeded.IsGrouped(),
		)
		if !ok {
			t.Fatalf("expected skill %q to be constructible from its seeded metadata", id)
		}

		if skill.GetID() != seeded.GetID() {
			t.Fatalf("expected constructed skill id %q, got %q", seeded.GetID(), skill.GetID())
		}

		if skill.GetAbilityScoreID() != seeded.GetAbilityScoreID() {
			t.Fatalf("expected constructed skill ability score %q, got %q", seeded.GetAbilityScoreID(), skill.GetAbilityScoreID())
		}
	}
}

func TestNewSkill_RejectsUnknownCoreLikeSkillIDs(t *testing.T) {
	invalidIDs := []SkillID{
		"Jump",
		"Open Lock",
		"knowledge",
		"Use magic device",
		"Knowledge(arcana)",
	}

	for _, id := range invalidIDs {
		if _, ok := NewSkill(id, ability.DexterityScore, false, false, false); ok {
			t.Fatalf("expected non-core or non-canonical skill id %q to be rejected", id)
		}
	}
}

func TestNewSpecializedGroupedSkill_UsesCoreFamilyMetadata(t *testing.T) {
	testCases := []struct {
		id                       SkillID
		familyID                 SkillID
		abilityScoreID           ability.AbilityScoreID
		trainedOnly              bool
		armorCheckPenaltyApplies bool
	}{
		{SkillID("Craft (alchemy)"), CraftSkillID, ability.IntelligenceScore, false, false},
		{SkillID("Knowledge (arcana)"), KnowledgeSkillID, ability.IntelligenceScore, true, false},
		{SkillID("Perform (sing)"), PerformSkillID, ability.CharismaScore, false, false},
		{SkillID("Profession (sailor)"), ProfessionSkillID, ability.WisdomScore, true, false},
	}

	for _, tc := range testCases {
		skill, ok := NewSpecializedGroupedSkill(tc.id)
		if !ok {
			t.Fatalf("expected specialized grouped skill %q to be constructed", tc.id)
		}

		if skill.GetFamilyID() != tc.familyID {
			t.Fatalf("expected family id %q, got %q", tc.familyID, skill.GetFamilyID())
		}

		if skill.GetAbilityScoreID() != tc.abilityScoreID {
			t.Fatalf("expected skill %q ability score %q, got %q", tc.id, tc.abilityScoreID, skill.GetAbilityScoreID())
		}

		if skill.IsTrainedOnly() != tc.trainedOnly {
			t.Fatalf("expected skill %q trained-only=%t, got %t", tc.id, tc.trainedOnly, skill.IsTrainedOnly())
		}

		if skill.AppliesArmorCheckPenalty() != tc.armorCheckPenaltyApplies {
			t.Fatalf("expected skill %q armor-check-penalty=%t, got %t", tc.id, tc.armorCheckPenaltyApplies, skill.AppliesArmorCheckPenalty())
		}
	}
}

func TestNewSpecializedGroupedSkill_RejectsFamiliesAndMalformedIDs(t *testing.T) {
	invalidIDs := []SkillID{
		PerformSkillID,
		SkillID("Acrobatics (urban)"),
		SkillID("Knowledge()"),
	}

	for _, id := range invalidIDs {
		if _, ok := NewSpecializedGroupedSkill(id); ok {
			t.Fatalf("expected specialized grouped skill %q to be rejected", id)
		}
	}
}

func TestNewSpecializedGroupedSkill_DoesNotChangeCoreCatalog(t *testing.T) {
	if _, ok := NewSpecializedGroupedSkill("Perform (sing)"); !ok {
		t.Fatal("expected grouped specialization to be accepted")
	}

	if len(coreSkills) != len(coreSkillOrder) {
		t.Fatalf("expected grouped specialization support not to change core catalog size, got %d skills and %d ordered ids", len(coreSkills), len(coreSkillOrder))
	}
}

func TestGetSkillByID_ReturnsSeededCoreSkill(t *testing.T) {
	skill, ok := GetSkillByID(PerceptionSkillID)
	if !ok {
		t.Fatal("expected perception skill lookup to succeed")
	}

	if skill.GetID() != PerceptionSkillID {
		t.Fatalf("expected looked up skill id %q, got %q", PerceptionSkillID, skill.GetID())
	}

	if skill.GetAbilityScoreID() != ability.WisdomScore {
		t.Fatalf("expected Perception ability score %q, got %q", ability.WisdomScore, skill.GetAbilityScoreID())
	}

	if skill.IsTrainedOnly() {
		t.Fatal("expected perception to not be trained-only")
	}
}

func TestGetSkillByID_RejectsUnknownSkill(t *testing.T) {
	if _, ok := GetSkillByID(SkillID("Jump")); ok {
		t.Fatal("expected unknown skill lookup to fail")
	}
}

func TestGetSkills_ReturnsSeededCatalogInCoreOrder(t *testing.T) {
	skills := GetSkills()
	if len(skills) != len(coreSkillOrder) {
		t.Fatalf("expected %d queried skills, got %d", len(coreSkillOrder), len(skills))
	}

	for i, expectedID := range coreSkillOrder {
		if skills[i].GetID() != expectedID {
			t.Fatalf("expected skill at index %d to be %q, got %q", i, expectedID, skills[i].GetID())
		}
	}
}

func TestGetSkills_ReturnsDetachedSlice(t *testing.T) {
	first := GetSkills()
	second := GetSkills()

	first[0] = skill{}

	if second[0].GetID() != AcrobaticsSkillID {
		t.Fatalf("expected detached skill slice to preserve %q, got %q", AcrobaticsSkillID, second[0].GetID())
	}
}
