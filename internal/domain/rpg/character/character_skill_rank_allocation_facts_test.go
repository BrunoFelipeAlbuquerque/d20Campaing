package character

import (
	"testing"

	"d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
	"d20campaigngenerator/internal/domain/rpg/character/skill"
)

func TestNewCharacterSkillRankAllocationFacts_ComposesWithinBudgetAndLevelCap(t *testing.T) {
	craftAlchemySkillID := skill.SkillID("Craft (alchemy)")
	knowledgeArcanaSkillID := skill.SkillID("Knowledge (arcana)")

	budget := mustNewCharacterSkillRankBudgetFactsForValidationTest(
		t,
		characterclass.FighterClassID,
		2,
		10,
		characterrace.HumanRaceID,
	)
	rideRanks := mustNewCharacterSkillRankAllocationForValidationTest(t, skill.RideSkillID, 2)
	craftAlchemyRanks := mustNewCharacterSkillRankAllocationForValidationTest(t, craftAlchemySkillID, 2)
	knowledgeArcanaRanks := mustNewCharacterSkillRankAllocationForValidationTest(t, knowledgeArcanaSkillID, 1)

	facts, ok := NewCharacterSkillRankAllocationFacts(
		budget,
		[]CharacterSkillRankAllocation{rideRanks, craftAlchemyRanks, knowledgeArcanaRanks},
	)
	if !ok {
		t.Fatal("expected skill-rank allocation facts to compose")
	}

	if facts.GetTotalAllocatedRanks() != 5 {
		t.Fatalf("expected 5 allocated ranks, got %d", facts.GetTotalAllocatedRanks())
	}

	assertSkillRankAllocation(t, facts, skill.RideSkillID, 2)
	assertSkillRankAllocation(t, facts, craftAlchemySkillID, 2)
	assertSkillRankAllocation(t, facts, knowledgeArcanaSkillID, 1)

	if facts.HasSkillRankAllocation(skill.KnowledgeSkillID) {
		t.Fatal("expected grouped family skill not to resolve as a concrete allocation")
	}

	allocations := facts.GetAllocations()
	allocations[0] = mustNewCharacterSkillRankAllocationForValidationTest(t, skill.AcrobaticsSkillID, 1)

	allocationsAgain := facts.GetAllocations()
	if allocationsAgain[0].GetSkillID() != skill.RideSkillID {
		t.Fatalf("expected defensive allocation copy to preserve %q", skill.RideSkillID)
	}
}

func TestNewCharacterSkillRankAllocationFacts_AllowsEmptySelectionWithinBudget(t *testing.T) {
	budget := mustNewCharacterSkillRankBudgetFactsForValidationTest(
		t,
		characterclass.WizardClassID,
		1,
		14,
		characterrace.ElfRaceID,
	)

	facts, ok := NewCharacterSkillRankAllocationFacts(budget, nil)
	if !ok {
		t.Fatal("expected empty skill-rank allocation facts to compose")
	}

	if facts.GetTotalAllocatedRanks() != 0 {
		t.Fatalf("expected no allocated ranks, got %d", facts.GetTotalAllocatedRanks())
	}

	if len(facts.GetAllocations()) != 0 {
		t.Fatal("expected no allocations")
	}
}

func TestNewCharacterSkillRankAllocationFacts_RejectsBudgetAndLevelCapViolations(t *testing.T) {
	budget := mustNewCharacterSkillRankBudgetFactsForValidationTest(
		t,
		characterclass.FighterClassID,
		1,
		10,
		characterrace.ElfRaceID,
	)

	testCases := []struct {
		name        string
		allocations []CharacterSkillRankAllocation
	}{
		{
			name: "skill ranks above character level",
			allocations: []CharacterSkillRankAllocation{
				mustNewCharacterSkillRankAllocationForValidationTest(t, skill.RideSkillID, 2),
			},
		},
		{
			name: "total ranks above budget",
			allocations: []CharacterSkillRankAllocation{
				mustNewCharacterSkillRankAllocationForValidationTest(t, skill.RideSkillID, 1),
				mustNewCharacterSkillRankAllocationForValidationTest(t, skill.AcrobaticsSkillID, 1),
				mustNewCharacterSkillRankAllocationForValidationTest(t, skill.ClimbSkillID, 1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterSkillRankAllocationFacts(budget, tc.allocations); ok {
				t.Fatal("expected skill-rank allocation facts to fail")
			}
		})
	}
}

func TestNewCharacterSkillRankAllocationFacts_RejectsDuplicateMalformedOrUnknownAllocations(t *testing.T) {
	budget := mustNewCharacterSkillRankBudgetFactsForValidationTest(
		t,
		characterclass.FighterClassID,
		2,
		10,
		characterrace.HumanRaceID,
	)
	rideRanks := mustNewCharacterSkillRankAllocationForValidationTest(t, skill.RideSkillID, 1)

	testCases := []struct {
		name        string
		allocations []CharacterSkillRankAllocation
	}{
		{"duplicate skill", []CharacterSkillRankAllocation{rideRanks, rideRanks}},
		{"zero allocation", []CharacterSkillRankAllocation{{}}},
		{"unknown skill", []CharacterSkillRankAllocation{{skillID: skill.SkillID("Jump"), ranks: 1}}},
		{"grouped family skill", []CharacterSkillRankAllocation{{skillID: skill.KnowledgeSkillID, ranks: 1}}},
		{"zero ranks", []CharacterSkillRankAllocation{{skillID: skill.RideSkillID, ranks: 0}}},
		{"negative ranks", []CharacterSkillRankAllocation{{skillID: skill.RideSkillID, ranks: -1}}},
		{"over cap grouped skill", []CharacterSkillRankAllocation{{skillID: skill.SkillID("Knowledge (arcana)"), ranks: 3}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterSkillRankAllocationFacts(budget, tc.allocations); ok {
				t.Fatal("expected skill-rank allocation facts to fail")
			}
		})
	}
}

func TestNewCharacterSkillRankAllocationFacts_RejectsMalformedBudgetFacts(t *testing.T) {
	rideRanks := mustNewCharacterSkillRankAllocationForValidationTest(t, skill.RideSkillID, 1)

	testCases := []struct {
		name   string
		budget CharacterSkillRankBudgetFacts
	}{
		{"zero budget", CharacterSkillRankBudgetFacts{}},
		{"negative budget", CharacterSkillRankBudgetFacts{totalSkillRanks: -1, classSkillRanks: -1, totalCharacterLevel: 1}},
		{"inconsistent total", CharacterSkillRankBudgetFacts{totalSkillRanks: 2, classSkillRanks: 1, totalCharacterLevel: 1}},
		{"zero character level", CharacterSkillRankBudgetFacts{totalSkillRanks: 1, classSkillRanks: 1}},
		{"character level above core range", CharacterSkillRankBudgetFacts{totalSkillRanks: 21, classSkillRanks: 21, totalCharacterLevel: 21}},
		{"missing class ranks", CharacterSkillRankBudgetFacts{totalSkillRanks: 1, racialBonusRanks: 1, totalCharacterLevel: 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterSkillRankAllocationFacts(tc.budget, []CharacterSkillRankAllocation{rideRanks}); ok {
				t.Fatal("expected skill-rank allocation facts to fail")
			}
		})
	}
}

func assertSkillRankAllocation(
	t *testing.T,
	facts CharacterSkillRankAllocationFacts,
	id skill.SkillID,
	expectedRanks int,
) {
	t.Helper()

	ranks, ok := facts.GetSkillRanks(id)
	if !ok || ranks != expectedRanks {
		t.Fatalf("expected %q skill ranks (%d, true), got (%d, %t)", id, expectedRanks, ranks, ok)
	}
}

func mustNewCharacterSkillRankBudgetFactsForValidationTest(
	t *testing.T,
	classID characterclass.ClassID,
	level int,
	intelligenceScore int,
	raceID characterrace.RaceID,
) CharacterSkillRankBudgetFacts {
	t.Helper()

	facts, ok := NewCharacterSkillRankBudgetFacts(
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, classID, level)},
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForTest(t, ability.IntelligenceScore, intelligenceScore)},
		mustNewCharacterRaceForSkillRankBudgetTest(t, raceID),
	)
	if !ok {
		t.Fatalf("expected skill-rank budget facts to compose")
	}

	return facts
}

func mustNewCharacterSkillRankAllocationForValidationTest(
	t *testing.T,
	id skill.SkillID,
	ranks int,
) CharacterSkillRankAllocation {
	t.Helper()

	allocation, ok := NewCharacterSkillRankAllocation(id, ranks)
	if !ok {
		t.Fatalf("expected skill-rank allocation %q %d to compose", id, ranks)
	}

	return allocation
}
