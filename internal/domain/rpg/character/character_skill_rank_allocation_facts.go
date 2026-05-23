package character

import "d20campaigngenerator/internal/domain/rpg/character/skill"

type characterSkillRankAllocationFacts struct {
	allocations         []CharacterSkillRankAllocation
	totalAllocatedRanks int
}
type CharacterSkillRankAllocationFacts = characterSkillRankAllocationFacts

func NewCharacterSkillRankAllocationFacts(
	budget CharacterSkillRankBudgetFacts,
	allocations []CharacterSkillRankAllocation,
) (CharacterSkillRankAllocationFacts, bool) {
	if !isValidCharacterSkillRankBudgetFacts(budget) {
		return characterSkillRankAllocationFacts{}, false
	}

	validatedAllocations, totalAllocatedRanks, ok := characterSkillRankAllocationsFromList(
		budget,
		allocations,
	)
	if !ok {
		return characterSkillRankAllocationFacts{}, false
	}

	return characterSkillRankAllocationFacts{
		allocations:         validatedAllocations,
		totalAllocatedRanks: totalAllocatedRanks,
	}, true
}

func (f characterSkillRankAllocationFacts) GetAllocations() []CharacterSkillRankAllocation {
	return append([]CharacterSkillRankAllocation(nil), f.allocations...)
}

func (f characterSkillRankAllocationFacts) GetTotalAllocatedRanks() int {
	return f.totalAllocatedRanks
}

func (f characterSkillRankAllocationFacts) GetSkillRanks(id skill.SkillID) (int, bool) {
	if _, ok := concreteCharacterSkillForRankAllocation(id); !ok {
		return 0, false
	}

	for _, allocation := range f.allocations {
		if allocation.GetSkillID() == id {
			return allocation.GetRanks(), true
		}
	}

	return 0, false
}

func (f characterSkillRankAllocationFacts) HasSkillRankAllocation(id skill.SkillID) bool {
	_, ok := f.GetSkillRanks(id)
	return ok
}

func isValidCharacterSkillRankBudgetFacts(budget CharacterSkillRankBudgetFacts) bool {
	if budget.totalCharacterLevel <= 0 || budget.totalCharacterLevel > maxCoreCharacterLevel {
		return false
	}

	if budget.totalSkillRanks <= 0 || budget.classSkillRanks <= 0 || budget.racialBonusRanks < 0 {
		return false
	}

	return budget.classSkillRanks+budget.racialBonusRanks == budget.totalSkillRanks
}

func characterSkillRankAllocationsFromList(
	budget CharacterSkillRankBudgetFacts,
	allocations []CharacterSkillRankAllocation,
) ([]CharacterSkillRankAllocation, int, bool) {
	seen := make(map[skill.SkillID]struct{}, len(allocations))
	validatedAllocations := make([]CharacterSkillRankAllocation, 0, len(allocations))
	totalAllocatedRanks := 0

	for _, allocation := range allocations {
		validatedAllocation, ok := NewCharacterSkillRankAllocation(
			allocation.GetSkillID(),
			allocation.GetRanks(),
		)
		if !ok {
			return nil, 0, false
		}

		if _, ok := seen[validatedAllocation.GetSkillID()]; ok {
			return nil, 0, false
		}

		if validatedAllocation.GetRanks() > budget.GetTotalCharacterLevel() {
			return nil, 0, false
		}

		totalAllocatedRanks += validatedAllocation.GetRanks()
		if totalAllocatedRanks > budget.GetTotalSkillRanks() {
			return nil, 0, false
		}

		seen[validatedAllocation.GetSkillID()] = struct{}{}
		validatedAllocations = append(validatedAllocations, validatedAllocation)
	}

	return validatedAllocations, totalAllocatedRanks, true
}
