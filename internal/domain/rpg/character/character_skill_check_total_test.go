package character

import (
	"testing"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
	"d20campaigngenerator/internal/domain/rpg/character/skill"
)

func TestNewCharacterSkillCheckTotalFacts_ComposesRankAbilityAndClassSkillBonus(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.DexterityScore:    16,
		ability.IntelligenceScore: 14,
		ability.CharismaScore:     12,
	})
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		[]CharacterSkillRankAllocation{
			mustNewCharacterSkillRankAllocationForSkillCheckTotalTest(t, skill.StealthSkillID, 1),
			mustNewCharacterSkillRankAllocationForSkillCheckTotalTest(t, skill.DisableDeviceSkillID, 1),
		},
	)

	facts, ok := NewCharacterSkillCheckTotalFacts(classLevels, abilityScores, rankFacts)
	if !ok {
		t.Fatal("expected skill check total facts to compose")
	}

	stealth, ok := facts.GetSkillCheckTotal(skill.StealthSkillID)
	if !ok {
		t.Fatal("expected Stealth total to resolve")
	}

	assertCharacterSkillCheckTotal(t, stealth, expectedCharacterSkillCheckTotal{
		skillID:                  skill.StealthSkillID,
		abilityScoreID:           ability.DexterityScore,
		abilityModifier:          3,
		ranks:                    1,
		classSkill:               true,
		classSkillBonus:          3,
		armorCheckPenaltyApplies: true,
		total:                    7,
	})

	disableDevice, ok := facts.GetSkillCheckTotal(skill.DisableDeviceSkillID)
	if !ok {
		t.Fatal("expected trained-only Disable Device total to resolve once ranked")
	}

	assertCharacterSkillCheckTotal(t, disableDevice, expectedCharacterSkillCheckTotal{
		skillID:                  skill.DisableDeviceSkillID,
		abilityScoreID:           ability.DexterityScore,
		abilityModifier:          3,
		ranks:                    1,
		classSkill:               true,
		classSkillBonus:          3,
		armorCheckPenaltyApplies: true,
		total:                    7,
	})
}

func TestCharacterSkillCheckTotalFacts_ComposesUntrainedSkillWithoutRanks(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.IntelligenceScore: 14,
		ability.CharismaScore:     12,
	})
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		nil,
	)

	facts, ok := NewCharacterSkillCheckTotalFacts(classLevels, abilityScores, rankFacts)
	if !ok {
		t.Fatal("expected skill check total facts to compose")
	}

	bluff, ok := facts.GetSkillCheckTotal(skill.BluffSkillID)
	if !ok {
		t.Fatal("expected untrained Bluff total to resolve")
	}

	assertCharacterSkillCheckTotal(t, bluff, expectedCharacterSkillCheckTotal{
		skillID:                  skill.BluffSkillID,
		abilityScoreID:           ability.CharismaScore,
		abilityModifier:          1,
		ranks:                    0,
		classSkill:               true,
		classSkillBonus:          0,
		armorCheckPenaltyApplies: false,
		total:                    1,
	})
}

func TestCharacterSkillCheckTotalFacts_RejectsTrainedOnlySkillWithoutRanks(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.IntelligenceScore: 14,
		ability.CharismaScore:     12,
	})
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		nil,
	)
	facts, ok := NewCharacterSkillCheckTotalFacts(classLevels, abilityScores, rankFacts)
	if !ok {
		t.Fatal("expected skill check total facts to compose")
	}

	if _, ok := facts.GetSkillCheckTotal(skill.UseMagicDeviceSkillID); ok {
		t.Fatal("expected trained-only Use Magic Device without ranks not to resolve")
	}
}

func TestCharacterSkillCheckTotalFacts_ComposesGroupedClassSkillTotals(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.IntelligenceScore: 14,
	})
	knowledgeArcanaSkillID := skill.SkillID("Knowledge (arcana)")
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		[]CharacterSkillRankAllocation{
			mustNewCharacterSkillRankAllocationForSkillCheckTotalTest(t, knowledgeArcanaSkillID, 1),
		},
	)

	facts, ok := NewCharacterSkillCheckTotalFacts(classLevels, abilityScores, rankFacts)
	if !ok {
		t.Fatal("expected skill check total facts to compose")
	}

	knowledgeArcana, ok := facts.GetSkillCheckTotal(knowledgeArcanaSkillID)
	if !ok {
		t.Fatal("expected Knowledge (arcana) total to resolve")
	}

	assertCharacterSkillCheckTotal(t, knowledgeArcana, expectedCharacterSkillCheckTotal{
		skillID:                  knowledgeArcanaSkillID,
		abilityScoreID:           ability.IntelligenceScore,
		abilityModifier:          2,
		ranks:                    1,
		classSkill:               true,
		classSkillBonus:          3,
		armorCheckPenaltyApplies: false,
		total:                    6,
	})
}

func TestCharacterSkillCheckTotalFacts_ComposesSpecificGroupedClassSkillTotals(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.IntelligenceScore: 14,
	})
	knowledgeEngineeringSkillID := skill.SkillID("Knowledge (engineering)")
	knowledgeArcanaSkillID := skill.SkillID("Knowledge (arcana)")
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		[]CharacterSkillRankAllocation{
			mustNewCharacterSkillRankAllocationForSkillCheckTotalTest(t, knowledgeEngineeringSkillID, 1),
			mustNewCharacterSkillRankAllocationForSkillCheckTotalTest(t, knowledgeArcanaSkillID, 1),
		},
	)

	facts, ok := NewCharacterSkillCheckTotalFacts(classLevels, abilityScores, rankFacts)
	if !ok {
		t.Fatal("expected skill check total facts to compose")
	}

	knowledgeEngineering, ok := facts.GetSkillCheckTotal(knowledgeEngineeringSkillID)
	if !ok {
		t.Fatal("expected Knowledge (engineering) total to resolve")
	}

	assertCharacterSkillCheckTotal(t, knowledgeEngineering, expectedCharacterSkillCheckTotal{
		skillID:                  knowledgeEngineeringSkillID,
		abilityScoreID:           ability.IntelligenceScore,
		abilityModifier:          2,
		ranks:                    1,
		classSkill:               true,
		classSkillBonus:          3,
		armorCheckPenaltyApplies: false,
		total:                    6,
	})

	knowledgeArcana, ok := facts.GetSkillCheckTotal(knowledgeArcanaSkillID)
	if !ok {
		t.Fatal("expected ranked Knowledge (arcana) total to resolve")
	}

	assertCharacterSkillCheckTotal(t, knowledgeArcana, expectedCharacterSkillCheckTotal{
		skillID:                  knowledgeArcanaSkillID,
		abilityScoreID:           ability.IntelligenceScore,
		abilityModifier:          2,
		ranks:                    1,
		classSkill:               false,
		classSkillBonus:          0,
		armorCheckPenaltyApplies: false,
		total:                    3,
	})
}

func TestCharacterSkillCheckTotalFacts_RejectsInvalidSkillQueries(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.IntelligenceScore: 14,
	})
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		nil,
	)
	facts, ok := NewCharacterSkillCheckTotalFacts(classLevels, abilityScores, rankFacts)
	if !ok {
		t.Fatal("expected skill check total facts to compose")
	}

	testCases := []skill.SkillID{
		skill.SkillID("Jump"),
		skill.KnowledgeSkillID,
		skill.SkillID("Acrobatics (urban)"),
	}

	for _, id := range testCases {
		if _, ok := facts.GetSkillCheckTotal(id); ok {
			t.Fatalf("expected skill %q not to resolve", id)
		}
	}
}

func TestNewCharacterSkillCheckTotalFacts_RejectsInvalidInputs(t *testing.T) {
	classLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 1),
	}
	abilityScores := completeCharacterAbilityScoresForSkillCheckTotalTest(t, map[ability.AbilityScoreID]int{
		ability.IntelligenceScore: 14,
	})
	rankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		classLevels,
		abilityScores,
		characterrace.ElfRaceID,
		nil,
	)
	mismatchedClassLevels := []CharacterClassLevel{
		mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 2),
	}
	mismatchedRankFacts := mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
		t,
		mismatchedClassLevels,
		abilityScores,
		characterrace.ElfRaceID,
		nil,
	)

	testCases := []struct {
		name          string
		classLevels   []CharacterClassLevel
		abilityScores []CharacterAbilityScore
		rankFacts     CharacterSkillRankAllocationFacts
	}{
		{"missing class levels", nil, abilityScores, rankFacts},
		{"missing ability score", classLevels, abilityScores[:5], rankFacts},
		{"duplicate ability score", classLevels, append(abilityScores, abilityScores[0]), rankFacts},
		{"zero rank facts", classLevels, abilityScores, CharacterSkillRankAllocationFacts{}},
		{"mismatched rank facts level cap", classLevels, abilityScores, mismatchedRankFacts},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterSkillCheckTotalFacts(tc.classLevels, tc.abilityScores, tc.rankFacts); ok {
				t.Fatal("expected skill check total facts to fail")
			}
		})
	}
}

type expectedCharacterSkillCheckTotal struct {
	skillID                  skill.SkillID
	abilityScoreID           ability.AbilityScoreID
	abilityModifier          int
	ranks                    int
	classSkill               bool
	classSkillBonus          int
	armorCheckPenaltyApplies bool
	total                    int
}

func assertCharacterSkillCheckTotal(
	t *testing.T,
	actual CharacterSkillCheckTotal,
	expected expectedCharacterSkillCheckTotal,
) {
	t.Helper()

	if actual.GetSkillID() != expected.skillID {
		t.Fatalf("expected skill id %q, got %q", expected.skillID, actual.GetSkillID())
	}

	if actual.GetAbilityScoreID() != expected.abilityScoreID {
		t.Fatalf("expected ability score %q, got %q", expected.abilityScoreID, actual.GetAbilityScoreID())
	}

	if actual.GetAbilityModifier() != expected.abilityModifier {
		t.Fatalf("expected ability modifier %d, got %d", expected.abilityModifier, actual.GetAbilityModifier())
	}

	if actual.GetRanks() != expected.ranks {
		t.Fatalf("expected ranks %d, got %d", expected.ranks, actual.GetRanks())
	}

	if actual.IsClassSkill() != expected.classSkill {
		t.Fatalf("expected class skill %t, got %t", expected.classSkill, actual.IsClassSkill())
	}

	if actual.GetClassSkillBonus() != expected.classSkillBonus {
		t.Fatalf("expected class skill bonus %d, got %d", expected.classSkillBonus, actual.GetClassSkillBonus())
	}

	if actual.AppliesArmorCheckPenalty() != expected.armorCheckPenaltyApplies {
		t.Fatalf("expected armor check penalty flag %t, got %t", expected.armorCheckPenaltyApplies, actual.AppliesArmorCheckPenalty())
	}

	if actual.GetTotal() != expected.total {
		t.Fatalf("expected total %d, got %d", expected.total, actual.GetTotal())
	}
}

func completeCharacterAbilityScoresForSkillCheckTotalTest(
	t *testing.T,
	overrides map[ability.AbilityScoreID]int,
) []CharacterAbilityScore {
	t.Helper()

	scores := []CharacterAbilityScore{
		mustNewCharacterAbilityScoreForTest(t, ability.StrengthScore, 10),
		mustNewCharacterAbilityScoreForTest(t, ability.DexterityScore, 10),
		mustNewCharacterAbilityScoreForTest(t, ability.ConstitutionScore, 10),
		mustNewCharacterAbilityScoreForTest(t, ability.IntelligenceScore, 10),
		mustNewCharacterAbilityScoreForTest(t, ability.WisdomScore, 10),
		mustNewCharacterAbilityScoreForTest(t, ability.CharismaScore, 10),
	}

	for i, score := range scores {
		if override, ok := overrides[score.GetAbilityScoreID()]; ok {
			scores[i] = mustNewCharacterAbilityScoreForTest(t, score.GetAbilityScoreID(), override)
		}
	}

	return scores
}

func mustNewCharacterSkillRankAllocationFactsForSkillCheckTotalTest(
	t *testing.T,
	classLevels []CharacterClassLevel,
	abilityScores []CharacterAbilityScore,
	raceID characterrace.RaceID,
	allocations []CharacterSkillRankAllocation,
) CharacterSkillRankAllocationFacts {
	t.Helper()

	rankBudget, ok := NewCharacterSkillRankBudgetFacts(
		classLevels,
		abilityScores,
		mustNewCharacterRaceForSkillCheckTotalTest(t, raceID),
	)
	if !ok {
		t.Fatal("expected skill-rank budget facts to compose")
	}

	rankFacts, ok := NewCharacterSkillRankAllocationFacts(rankBudget, allocations)
	if !ok {
		t.Fatal("expected skill-rank allocation facts to compose")
	}

	return rankFacts
}

func mustNewCharacterSkillRankAllocationForSkillCheckTotalTest(
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

func mustNewCharacterRaceForSkillCheckTotalTest(
	t *testing.T,
	id characterrace.RaceID,
) CharacterRace {
	t.Helper()

	race, ok := NewCharacterRace(id)
	if !ok {
		t.Fatalf("expected race %q to compose", id)
	}

	return race
}
