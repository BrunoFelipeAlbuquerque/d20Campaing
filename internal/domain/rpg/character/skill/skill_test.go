package skill

import (
	"testing"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
)

func TestNewSkill_ConstructsValidatedSkillChassis(t *testing.T) {
	skill, ok := NewSkill(SkillID("Acrobatics"), ability.DexterityScore, false, true, false)
	if !ok {
		t.Fatal("expected skill chassis to be constructed")
	}

	if skill.GetID() != SkillID("Acrobatics") {
		t.Fatalf("expected skill id %q, got %q", SkillID("Acrobatics"), skill.GetID())
	}

	if skill.GetID().GetName() != "Acrobatics" {
		t.Fatalf("expected skill name %q, got %q", "Acrobatics", skill.GetID().GetName())
	}

	if skill.GetAbilityScoreID() != ability.DexterityScore {
		t.Fatalf("expected ability score %q, got %q", ability.DexterityScore, skill.GetAbilityScoreID())
	}

	if skill.IsTrainedOnly() {
		t.Fatal("expected Acrobatics not to be trained only")
	}

	if !skill.AppliesArmorCheckPenalty() {
		t.Fatal("expected Acrobatics to apply armor check penalty metadata")
	}

	if skill.IsGrouped() {
		t.Fatal("expected Acrobatics not to be grouped")
	}
}

func TestNewSkill_AllowsCoreMultiwordSkillIDs(t *testing.T) {
	skill, ok := NewSkill(SkillID("Sleight of Hand"), ability.DexterityScore, true, true, false)
	if !ok {
		t.Fatal("expected multiword skill id to be accepted")
	}

	if skill.GetAbilityScoreID() != ability.DexterityScore {
		t.Fatalf("expected Sleight of Hand ability score %q, got %q", ability.DexterityScore, skill.GetAbilityScoreID())
	}

	if !skill.IsTrainedOnly() {
		t.Fatal("expected Sleight of Hand to preserve trained-only metadata")
	}

	if !skill.AppliesArmorCheckPenalty() {
		t.Fatal("expected Sleight of Hand to preserve armor check penalty metadata")
	}
}

func TestNewSkill_ModelsGroupedSkillFamilies(t *testing.T) {
	tests := []struct {
		id             SkillID
		abilityScoreID ability.AbilityScoreID
		trainedOnly    bool
	}{
		{CraftSkillID, ability.IntelligenceScore, false},
		{KnowledgeSkillID, ability.IntelligenceScore, true},
		{PerformSkillID, ability.CharismaScore, false},
		{ProfessionSkillID, ability.WisdomScore, true},
	}

	for _, tc := range tests {
		skill, ok := NewSkill(tc.id, tc.abilityScoreID, tc.trainedOnly, false, true)
		if !ok {
			t.Fatalf("expected grouped skill %q to be constructed", tc.id)
		}

		if !skill.IsGrouped() {
			t.Fatalf("expected skill %q to be marked grouped", tc.id)
		}

		if skill.GetAbilityScoreID() != tc.abilityScoreID {
			t.Fatalf("expected skill %q ability score %q, got %q", tc.id, tc.abilityScoreID, skill.GetAbilityScoreID())
		}
	}
}

func TestNewSkill_ModelsSpecializedGroupedSkillEntries(t *testing.T) {
	skill, ok := NewSkill(SkillID("Knowledge (arcana)"), ability.IntelligenceScore, true, false, true)
	if !ok {
		t.Fatal("expected specialized grouped skill entry to be constructed")
	}

	if skill.GetID() != SkillID("Knowledge (arcana)") {
		t.Fatalf("expected specialized skill id %q, got %q", SkillID("Knowledge (arcana)"), skill.GetID())
	}

	if skill.GetFamilyID() != KnowledgeSkillID {
		t.Fatalf("expected grouped family id %q, got %q", KnowledgeSkillID, skill.GetFamilyID())
	}

	if skill.GetAbilityScoreID() != ability.IntelligenceScore {
		t.Fatalf("expected grouped ability score %q, got %q", ability.IntelligenceScore, skill.GetAbilityScoreID())
	}

	specialization, ok := skill.GetSpecialization()
	if !ok || specialization != "arcana" {
		t.Fatalf("expected specialization (%q, true), got (%q, %t)", "arcana", specialization, ok)
	}

	if skill.GetID().GetName() != "Knowledge (arcana)" {
		t.Fatalf("expected specialized grouped skill name %q, got %q", "Knowledge (arcana)", skill.GetID().GetName())
	}
}

func TestNewSkill_RejectsInvalidInputs(t *testing.T) {
	if _, ok := NewSkill("", ability.DexterityScore, false, false, false); ok {
		t.Fatal("expected empty skill id to be rejected")
	}

	if _, ok := NewSkill("   ", ability.DexterityScore, false, false, false); ok {
		t.Fatal("expected blank skill id to be rejected")
	}

	if _, ok := NewSkill(" Acrobatics", ability.DexterityScore, false, true, false); ok {
		t.Fatal("expected skill id with surrounding whitespace to be rejected")
	}

	if _, ok := NewSkill(CraftSkillID, ability.IntelligenceScore, false, false, false); ok {
		t.Fatal("expected grouped skill id without grouped metadata to be rejected")
	}

	if _, ok := NewSkill("Acrobatics", ability.DexterityScore, false, true, true); ok {
		t.Fatal("expected ungrouped skill id with grouped metadata to be rejected")
	}

	if _, ok := NewSkill("Knowledge (arcana)", ability.IntelligenceScore, true, false, false); ok {
		t.Fatal("expected specialized grouped skill id without grouped metadata to be rejected")
	}

	if _, ok := NewSkill("Knowledge()", ability.IntelligenceScore, true, false, true); ok {
		t.Fatal("expected malformed grouped specialization to be rejected")
	}

	if _, ok := NewSkill("Acrobatics (urban)", ability.DexterityScore, false, true, false); ok {
		t.Fatal("expected specialization on ungrouped skill family to be rejected")
	}

	if _, ok := NewSkill(AcrobaticsSkillID, ability.AbilityScoreID("LUCK"), false, true, false); ok {
		t.Fatal("expected unknown skill ability score to be rejected")
	}

	if _, ok := NewSkill(AcrobaticsSkillID, ability.IntelligenceScore, false, true, false); ok {
		t.Fatal("expected mismatched skill ability score to be rejected")
	}

	if _, ok := NewSkill(AcrobaticsSkillID, ability.ConstitutionScore, false, true, false); ok {
		t.Fatal("expected unsupported skill ability score to be rejected")
	}

	if SkillID(" ").GetName() != "" {
		t.Fatal("expected invalid skill id name lookup to be empty")
	}
}
