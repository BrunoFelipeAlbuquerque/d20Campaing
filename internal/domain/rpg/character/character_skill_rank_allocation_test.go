package character

import (
	"testing"

	"d20campaigngenerator/internal/domain/rpg/character/skill"
)

func TestNewCharacterSkillRankAllocation_ComposesCoreSkillFact(t *testing.T) {
	allocation, ok := NewCharacterSkillRankAllocation(skill.RideSkillID, 1)
	if !ok {
		t.Fatal("expected core skill rank allocation to compose")
	}

	if allocation.GetSkillID() != skill.RideSkillID {
		t.Fatalf("expected skill id %q, got %q", skill.RideSkillID, allocation.GetSkillID())
	}

	if allocation.GetRanks() != 1 {
		t.Fatalf("expected 1 rank, got %d", allocation.GetRanks())
	}

	selectedSkill, ok := allocation.GetSkill()
	if !ok {
		t.Fatal("expected allocated skill to resolve")
	}

	if selectedSkill.GetID() != skill.RideSkillID {
		t.Fatalf("expected allocated skill %q, got %q", skill.RideSkillID, selectedSkill.GetID())
	}

	if selectedSkill.IsGrouped() {
		t.Fatal("expected Ride allocation not to be grouped")
	}

	if !selectedSkill.AppliesArmorCheckPenalty() {
		t.Fatal("expected Ride metadata to preserve armor check penalty flag")
	}
}

func TestNewCharacterSkillRankAllocation_ComposesGroupedSpecializedSkillFacts(t *testing.T) {
	testCases := []struct {
		name           string
		skillID        skill.SkillID
		familyID       skill.SkillID
		specialization string
		trainedOnly    bool
		expectedRanks  int
	}{
		{"craft", skill.SkillID("Craft (alchemy)"), skill.CraftSkillID, "alchemy", false, 1},
		{"knowledge", skill.SkillID("Knowledge (arcana)"), skill.KnowledgeSkillID, "arcana", true, 2},
		{"perform", skill.SkillID("Perform (sing)"), skill.PerformSkillID, "sing", false, 3},
		{"profession", skill.SkillID("Profession (sailor)"), skill.ProfessionSkillID, "sailor", true, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allocation, ok := NewCharacterSkillRankAllocation(tc.skillID, tc.expectedRanks)
			if !ok {
				t.Fatal("expected grouped specialized skill rank allocation to compose")
			}

			if allocation.GetSkillID() != tc.skillID {
				t.Fatalf("expected skill id %q, got %q", tc.skillID, allocation.GetSkillID())
			}

			if allocation.GetRanks() != tc.expectedRanks {
				t.Fatalf("expected %d ranks, got %d", tc.expectedRanks, allocation.GetRanks())
			}

			selectedSkill, ok := allocation.GetSkill()
			if !ok {
				t.Fatal("expected allocated grouped skill to resolve")
			}

			if selectedSkill.GetFamilyID() != tc.familyID {
				t.Fatalf("expected family id %q, got %q", tc.familyID, selectedSkill.GetFamilyID())
			}

			specialization, ok := selectedSkill.GetSpecialization()
			if !ok || specialization != tc.specialization {
				t.Fatalf("expected specialization (%q, true), got (%q, %t)", tc.specialization, specialization, ok)
			}

			if !selectedSkill.IsGrouped() {
				t.Fatal("expected specialized skill to stay grouped")
			}

			if selectedSkill.IsTrainedOnly() != tc.trainedOnly {
				t.Fatalf("expected trained-only=%t, got %t", tc.trainedOnly, selectedSkill.IsTrainedOnly())
			}
		})
	}
}

func TestNewCharacterSkillRankAllocation_RejectsInvalidSkillOrRanks(t *testing.T) {
	testCases := []struct {
		name    string
		skillID skill.SkillID
		ranks   int
	}{
		{"zero ranks", skill.RideSkillID, 0},
		{"negative ranks", skill.RideSkillID, -1},
		{"unknown skill", skill.SkillID("Jump"), 1},
		{"malformed casing", skill.SkillID("knowledge"), 1},
		{"grouped family without specialization", skill.KnowledgeSkillID, 1},
		{"ungrouped skill with specialization", skill.SkillID("Acrobatics (urban)"), 1},
		{"malformed grouped specialization", skill.SkillID("Knowledge()"), 1},
		{"empty grouped specialization", skill.SkillID("Knowledge ()"), 1},
		{"trimmed grouped specialization", skill.SkillID("Knowledge ( arcana)"), 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterSkillRankAllocation(tc.skillID, tc.ranks); ok {
				t.Fatal("expected skill rank allocation to fail")
			}
		})
	}
}

func TestCharacterSkillRankAllocation_ZeroValueDoesNotResolve(t *testing.T) {
	var allocation CharacterSkillRankAllocation

	if _, ok := allocation.GetSkill(); ok {
		t.Fatal("expected zero-value skill rank allocation not to resolve")
	}
}
