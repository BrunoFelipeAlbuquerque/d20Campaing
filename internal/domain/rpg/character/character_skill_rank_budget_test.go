package character

import (
	"testing"

	"d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
)

func TestNewCharacterSkillRankBudgetFacts_ComposesClassAndIntelligenceBudget(t *testing.T) {
	facts, ok := NewCharacterSkillRankBudgetFacts(
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1)},
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 14)},
		mustNewCharacterRaceForSkillRankBudgetTest(t, characterrace.ElfRaceID),
	)
	if !ok {
		t.Fatal("expected skill-rank budget facts to compose")
	}

	if facts.GetTotalSkillRanks() != 4 {
		t.Fatalf("expected total skill ranks 4, got %d", facts.GetTotalSkillRanks())
	}

	if facts.GetClassSkillRanks() != 4 {
		t.Fatalf("expected class skill ranks 4, got %d", facts.GetClassSkillRanks())
	}

	if facts.GetRacialBonusRanks() != 0 {
		t.Fatalf("expected no racial bonus ranks, got %d", facts.GetRacialBonusRanks())
	}

	if facts.GetIntelligenceModifier() != 2 {
		t.Fatalf("expected Intelligence modifier 2, got %d", facts.GetIntelligenceModifier())
	}

	if facts.GetTotalCharacterLevel() != 1 {
		t.Fatalf("expected total character level 1, got %d", facts.GetTotalCharacterLevel())
	}
}

func TestNewCharacterSkillRankBudgetFacts_ComposesMulticlassBudget(t *testing.T) {
	facts, ok := NewCharacterSkillRankBudgetFacts(
		[]CharacterClassLevel{
			mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 2),
			mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1),
		},
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 12)},
		mustNewCharacterRaceForSkillRankBudgetTest(t, characterrace.DwarfRaceID),
	)
	if !ok {
		t.Fatal("expected multiclass skill-rank budget facts to compose")
	}

	if facts.GetClassSkillRanks() != 21 {
		t.Fatalf("expected class skill ranks 21, got %d", facts.GetClassSkillRanks())
	}

	if facts.GetTotalSkillRanks() != 21 {
		t.Fatalf("expected total skill ranks 21, got %d", facts.GetTotalSkillRanks())
	}
}

func TestNewCharacterSkillRankBudgetFacts_AddsHumanSkilledRanksByCharacterLevel(t *testing.T) {
	facts, ok := NewCharacterSkillRankBudgetFacts(
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 2)},
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 10)},
		mustNewCharacterRaceForSkillRankBudgetTest(t, characterrace.HumanRaceID),
	)
	if !ok {
		t.Fatal("expected human skill-rank budget facts to compose")
	}

	if facts.GetClassSkillRanks() != 4 {
		t.Fatalf("expected class skill ranks 4, got %d", facts.GetClassSkillRanks())
	}

	if facts.GetRacialBonusRanks() != 2 {
		t.Fatalf("expected human Skilled racial bonus ranks 2, got %d", facts.GetRacialBonusRanks())
	}

	if facts.GetTotalSkillRanks() != 6 {
		t.Fatalf("expected total skill ranks 6, got %d", facts.GetTotalSkillRanks())
	}
}

func TestNewCharacterSkillRankBudgetFacts_AppliesMinimumOneRankPerLevel(t *testing.T) {
	facts, ok := NewCharacterSkillRankBudgetFacts(
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 2)},
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 1)},
		mustNewCharacterRaceForSkillRankBudgetTest(t, characterrace.ElfRaceID),
	)
	if !ok {
		t.Fatal("expected low-Intelligence skill-rank budget facts to compose")
	}

	if facts.GetIntelligenceModifier() != -5 {
		t.Fatalf("expected Intelligence modifier -5, got %d", facts.GetIntelligenceModifier())
	}

	if facts.GetClassSkillRanks() != 2 {
		t.Fatalf("expected minimum class skill ranks 2, got %d", facts.GetClassSkillRanks())
	}

	if facts.GetTotalSkillRanks() != 2 {
		t.Fatalf("expected total skill ranks 2, got %d", facts.GetTotalSkillRanks())
	}
}

func TestNewCharacterSkillRankBudgetFacts_RejectsInvalidInputs(t *testing.T) {
	fighterLevel := mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)
	validClassLevels := []CharacterClassLevel{fighterLevel}
	validAbilityScores := []CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.IntelligenceScore, 10)}
	validRace := mustNewCharacterRaceForSkillRankBudgetTest(t, characterrace.HumanRaceID)

	testCases := []struct {
		name          string
		classLevels   []CharacterClassLevel
		abilityScores []CharacterAbilityScore
		race          CharacterRace
	}{
		{"missing class levels", nil, validAbilityScores, validRace},
		{"malformed class level", []CharacterClassLevel{{classID: characterclass.FighterClassID, level: 0}}, validAbilityScores, validRace},
		{"duplicate class levels", []CharacterClassLevel{fighterLevel, fighterLevel}, validAbilityScores, validRace},
		{"missing intelligence", validClassLevels, nil, validRace},
		{"missing intelligence score", validClassLevels, []CharacterAbilityScore{mustNewCharacterAbilityScoreForAbilityTest(t, ability.StrengthScore, 10)}, validRace},
		{"malformed intelligence", validClassLevels, []CharacterAbilityScore{{id: ability.IntelligenceScore, score: -1}}, validRace},
		{"duplicate intelligence", validClassLevels, []CharacterAbilityScore{validAbilityScores[0], validAbilityScores[0]}, validRace},
		{"zero race", validClassLevels, validAbilityScores, CharacterRace{}},
		{"unknown race", validClassLevels, validAbilityScores, CharacterRace{id: characterrace.RaceID("android")}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterSkillRankBudgetFacts(tc.classLevels, tc.abilityScores, tc.race); ok {
				t.Fatal("expected skill-rank budget facts to fail")
			}
		})
	}
}

func mustNewCharacterRaceForSkillRankBudgetTest(
	t *testing.T,
	id characterrace.RaceID,
) CharacterRace {
	t.Helper()

	race, ok := NewCharacterRace(id)
	if !ok {
		t.Fatalf("expected character race %q to compose", id)
	}

	return race
}
