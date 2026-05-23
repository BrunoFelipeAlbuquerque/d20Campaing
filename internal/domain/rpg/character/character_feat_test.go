package character

import (
	"testing"

	"d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterequipment "d20campaigngenerator/internal/domain/rpg/character/equipment"
	characterfeat "d20campaigngenerator/internal/domain/rpg/character/feat"
	characterrace "d20campaigngenerator/internal/domain/rpg/character/race"
	"d20campaigngenerator/internal/domain/rpg/character/skill"
	characterspell "d20campaigngenerator/internal/domain/rpg/character/spell"
)

func TestNewCharacterFeat_ComposesAbilityScoreAndBaseAttackBonusPrerequisites(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForTest(t, ability.StrengthScore, 13)},
		1,
		nil,
		nil,
		nil,
		nil,
	)

	selectedFeat, ok := NewCharacterFeat(characterfeat.PowerAttackFeatID, state)
	if !ok {
		t.Fatal("expected power attack prerequisites to compose")
	}

	if selectedFeat.GetFeatID() != characterfeat.PowerAttackFeatID {
		t.Fatalf("expected selected feat id %q, got %q", characterfeat.PowerAttackFeatID, selectedFeat.GetFeatID())
	}

	feat, ok := selectedFeat.GetFeat()
	if !ok {
		t.Fatal("expected selected feat to resolve")
	}

	if feat.GetCategory() != characterfeat.CombatFeatCategory {
		t.Fatalf("expected selected feat category %q, got %q", characterfeat.CombatFeatCategory, feat.GetCategory())
	}

	lowStrengthState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForTest(t, ability.StrengthScore, 12)},
		1,
		nil,
		nil,
		nil,
		nil,
	)
	if _, ok := NewCharacterFeat(characterfeat.PowerAttackFeatID, lowStrengthState); ok {
		t.Fatal("expected power attack to reject a low strength score")
	}

	lowBaseAttackState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		[]CharacterAbilityScore{mustNewCharacterAbilityScoreForTest(t, ability.StrengthScore, 13)},
		0,
		nil,
		nil,
		nil,
		nil,
	)
	if _, ok := NewCharacterFeat(characterfeat.PowerAttackFeatID, lowBaseAttackState); ok {
		t.Fatal("expected power attack to reject a low base attack bonus")
	}
}

func TestNewCharacterFeat_ComposesSkillRanksAndFeatPrerequisites(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		nil,
		[]CharacterSkillRanks{mustNewCharacterSkillRanksForTest(t, skill.RideSkillID, 1)},
		[]characterfeat.FeatID{characterfeat.MountedCombatFeatID},
	)

	if _, ok := NewCharacterFeat(characterfeat.MountedArcheryFeatID, state); !ok {
		t.Fatal("expected mounted archery prerequisites to compose")
	}

	missingFeatState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		nil,
		[]CharacterSkillRanks{mustNewCharacterSkillRanksForTest(t, skill.RideSkillID, 1)},
		nil,
	)
	if _, ok := NewCharacterFeat(characterfeat.MountedArcheryFeatID, missingFeatState); ok {
		t.Fatal("expected mounted archery to reject missing mounted combat")
	}

	missingSkillState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		[]characterfeat.FeatID{characterfeat.MountedCombatFeatID},
	)
	if _, ok := NewCharacterFeat(characterfeat.MountedArcheryFeatID, missingSkillState); ok {
		t.Fatal("expected mounted archery to reject missing ride ranks")
	}
}

func TestNewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts_ComposesSeededSkillRankPrerequisites(t *testing.T) {
	rideAllocations := mustNewCharacterSkillRankAllocationFactsForFeatPrerequisiteTest(
		t,
		characterclass.FighterClassID,
		1,
		10,
		characterrace.ElfRaceID,
		[]CharacterSkillRankAllocation{
			mustNewCharacterSkillRankAllocationForFeatPrerequisiteTest(t, skill.RideSkillID, 1),
		},
	)

	rideState, ok := NewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts(
		nil,
		0,
		nil,
		nil,
		nil,
		rideAllocations,
		[]characterfeat.FeatID{characterfeat.MountedCombatFeatID},
	)
	if !ok {
		t.Fatal("expected feat prerequisite state to compose from skill-rank allocations")
	}

	if _, ok := NewCharacterFeat(characterfeat.MountedArcheryFeatID, rideState); !ok {
		t.Fatal("expected mounted archery to compose from allocated Ride ranks")
	}

	craftAlchemyAllocations := mustNewCharacterSkillRankAllocationFactsForFeatPrerequisiteTest(
		t,
		characterclass.WizardClassID,
		5,
		10,
		characterrace.ElfRaceID,
		[]CharacterSkillRankAllocation{
			mustNewCharacterSkillRankAllocationForFeatPrerequisiteTest(t, skill.SkillID("Craft (alchemy)"), 5),
		},
	)

	craftState, ok := NewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts(
		nil,
		0,
		nil,
		nil,
		nil,
		craftAlchemyAllocations,
		nil,
	)
	if !ok {
		t.Fatal("expected grouped skill allocation to feed feat prerequisite state")
	}

	if _, ok := NewCharacterFeat(characterfeat.MasterCraftsmanFeatID, craftState); !ok {
		t.Fatal("expected master craftsman to compose from concrete Craft ranks")
	}

	splitCraftAllocations := mustNewCharacterSkillRankAllocationFactsForFeatPrerequisiteTest(
		t,
		characterclass.WizardClassID,
		5,
		10,
		characterrace.ElfRaceID,
		[]CharacterSkillRankAllocation{
			mustNewCharacterSkillRankAllocationForFeatPrerequisiteTest(t, skill.SkillID("Craft (alchemy)"), 3),
			mustNewCharacterSkillRankAllocationForFeatPrerequisiteTest(t, skill.SkillID("Craft (weapons)"), 2),
		},
	)

	splitCraftState, ok := NewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts(
		nil,
		0,
		nil,
		nil,
		nil,
		splitCraftAllocations,
		nil,
	)
	if !ok {
		t.Fatal("expected split grouped skill allocations to feed feat prerequisite state")
	}

	if _, ok := NewCharacterFeat(characterfeat.MasterCraftsmanFeatID, splitCraftState); ok {
		t.Fatal("expected master craftsman to reject split Craft ranks below the single-skill prerequisite")
	}
}

func TestNewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts_RejectsMissingSkillRanks(t *testing.T) {
	emptyAllocations := mustNewCharacterSkillRankAllocationFactsForFeatPrerequisiteTest(
		t,
		characterclass.FighterClassID,
		1,
		10,
		characterrace.ElfRaceID,
		nil,
	)

	state, ok := NewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts(
		nil,
		0,
		nil,
		nil,
		nil,
		emptyAllocations,
		nil,
	)
	if !ok {
		t.Fatal("expected empty allocated ranks to compose into feat prerequisite state")
	}

	if _, ok := NewCharacterFeat(characterfeat.MountedCombatFeatID, state); ok {
		t.Fatal("expected mounted combat to reject missing allocated Ride ranks")
	}
}

func TestNewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts_RejectsMalformedAllocationFacts(t *testing.T) {
	testCases := []struct {
		name        string
		allocations CharacterSkillRankAllocationFacts
	}{
		{"zero facts", CharacterSkillRankAllocationFacts{}},
		{
			name: "unknown allocation",
			allocations: CharacterSkillRankAllocationFacts{
				valid:               true,
				skillRankBudget:     1,
				rankCap:             1,
				totalAllocatedRanks: 1,
				allocations: []CharacterSkillRankAllocation{
					{skillID: skill.SkillID("Jump"), ranks: 1},
				},
			},
		},
		{
			name: "over rank cap",
			allocations: CharacterSkillRankAllocationFacts{
				valid:               true,
				skillRankBudget:     2,
				rankCap:             1,
				totalAllocatedRanks: 2,
				allocations: []CharacterSkillRankAllocation{
					{skillID: skill.RideSkillID, ranks: 2},
				},
			},
		},
		{
			name: "inconsistent total",
			allocations: CharacterSkillRankAllocationFacts{
				valid:               true,
				skillRankBudget:     2,
				rankCap:             1,
				totalAllocatedRanks: 2,
				allocations: []CharacterSkillRankAllocation{
					{skillID: skill.RideSkillID, ranks: 1},
				},
			},
		},
		{
			name: "duplicate allocation",
			allocations: CharacterSkillRankAllocationFacts{
				valid:               true,
				skillRankBudget:     2,
				rankCap:             1,
				totalAllocatedRanks: 2,
				allocations: []CharacterSkillRankAllocation{
					{skillID: skill.RideSkillID, ranks: 1},
					{skillID: skill.RideSkillID, ranks: 1},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := NewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts(
				nil,
				0,
				nil,
				nil,
				nil,
				tc.allocations,
				nil,
			); ok {
				t.Fatal("expected feat prerequisite state to reject malformed allocation facts")
			}
		})
	}
}

func TestNewCharacterFeat_ComposesClassFeaturePrerequisites(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		[]characterclass.ClassFeatureID{characterclass.RageClassFeatureID},
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ExtraRageFeatID, state); !ok {
		t.Fatal("expected extra rage prerequisites to compose")
	}

	missingFeatureState := mustNewCharacterFeatPrerequisiteStateForTest(t, nil, 0, nil, nil, nil, nil)
	if _, ok := NewCharacterFeat(characterfeat.ExtraRageFeatID, missingFeatureState); ok {
		t.Fatal("expected extra rage to reject missing rage feature")
	}
}

func TestNewCharacterFeat_ComposesClassLevelAndAnyFeatPrerequisites(t *testing.T) {
	anyFeatState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		8,
		nil,
		nil,
		nil,
		[]characterfeat.FeatID{characterfeat.CatchOffGuardFeatID},
	)

	if _, ok := NewCharacterFeat(characterfeat.ImprovisedWeaponMasteryFeatID, anyFeatState); !ok {
		t.Fatal("expected any-feat prerequisite to compose")
	}

	classLevelState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		1,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		[]characterfeat.FeatID{
			characterfeat.ShieldFocusFeatID,
			characterfeat.ShieldProficiencyFeatID,
		},
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterShieldFocusFeatID, classLevelState); !ok {
		t.Fatal("expected class-level feat prerequisites to compose")
	}

	lowClassLevelState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		1,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 7)},
		nil,
		nil,
		[]characterfeat.FeatID{
			characterfeat.ShieldFocusFeatID,
			characterfeat.ShieldProficiencyFeatID,
		},
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterShieldFocusFeatID, lowClassLevelState); ok {
		t.Fatal("expected class-level feat prerequisites to reject a low fighter level")
	}
}

func TestNewCharacterFeat_ComposesCharacterLevelAndSpellcastingPrerequisites(t *testing.T) {
	characterLevelState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		[]CharacterClassLevel{
			mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 5),
			mustNewCharacterClassLevelForTest(t, characterclass.RogueClassID, 2),
		},
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.LeadershipFeatID, characterLevelState); !ok {
		t.Fatal("expected character-level prerequisite to compose from total class levels")
	}

	spellcastingState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1)},
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ArcaneStrikeFeatID, spellcastingState); !ok {
		t.Fatal("expected arcane spellcasting prerequisite to compose from selected class levels")
	}
}

func TestNewCharacterFeat_ComposesCasterLevelPrerequisites(t *testing.T) {
	firstLevelCasterState := mustNewCharacterFeatPrerequisiteStateWithCasterLevelsForTest(
		t,
		nil,
		0,
		[]CharacterCasterLevel{mustNewCharacterCasterLevelForTest(t, ability.ArcaneCasterSource, 1)},
		nil,
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ScribeScrollFeatID, firstLevelCasterState); !ok {
		t.Fatal("expected scribe scroll caster-level prerequisite to compose")
	}

	if _, ok := NewCharacterFeat(characterfeat.BrewPotionFeatID, firstLevelCasterState); ok {
		t.Fatal("expected brew potion to reject a low caster level")
	}

	thirdLevelCasterState := mustNewCharacterFeatPrerequisiteStateWithCasterLevelsForTest(
		t,
		nil,
		0,
		[]CharacterCasterLevel{mustNewCharacterCasterLevelForTest(t, ability.DivineCasterSource, 3)},
		nil,
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.BrewPotionFeatID, thirdLevelCasterState); !ok {
		t.Fatal("expected brew potion caster-level prerequisite to compose from any caster source")
	}
}

func TestNewCharacterFeat_CasterLevelPrerequisiteRequiresExplicitCasterLevelFact(t *testing.T) {
	classOnlyState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 5)},
		nil,
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.BrewPotionFeatID, classOnlyState); ok {
		t.Fatal("expected caster-level prerequisites to reject class levels without a caster-level fact")
	}
}

func TestNewCharacterFeat_ComposesSelectedWeaponProficiencyPrerequisiteFromCategory(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
		t,
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)},
		nil,
		nil,
		mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID),
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.WeaponFocusFeatID, state); !ok {
		t.Fatal("expected weapon focus to compose from selected simple weapon and fighter category proficiency")
	}
}

func TestNewCharacterFeat_ComposesSelectedWeaponProficiencyPrerequisiteFromIndividualWeapon(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
		t,
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1)},
		nil,
		nil,
		mustNewCharacterSelectedWeaponForTest(t, characterequipment.CrossbowHeavyWeaponID),
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.WeaponFocusFeatID, state); !ok {
		t.Fatal("expected weapon focus to compose from selected heavy crossbow and wizard individual proficiency")
	}
}

func TestNewCharacterFeat_RejectsSelectedWeaponProficiencyWithoutMatchingClassProficiency(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
		t,
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.WizardClassID, 1)},
		nil,
		nil,
		mustNewCharacterSelectedWeaponForTest(t, characterequipment.SlingWeaponID),
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.WeaponFocusFeatID, state); ok {
		t.Fatal("expected weapon focus to reject selected sling without wizard proficiency")
	}
}

func TestNewCharacterFeat_RejectsSelectedWeaponProficiencyWithoutSelectedWeaponContext(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(t, nil, 1, nil, nil, nil, nil)

	if _, ok := NewCharacterFeat(characterfeat.WeaponFocusFeatID, state); ok {
		t.Fatal("expected selected weapon prerequisite to reject without selected weapon context")
	}
}

func TestNewCharacterFeat_ComposesSameSelectionWeaponPrerequisite(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)
	weaponFocus := mustNewCharacterSelectedWeaponFeatForTest(t, characterfeat.WeaponFocusFeatID, selectedWeapon)
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponFeatsForTest(
		t,
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		selectedWeapon,
		[]CharacterSelectedWeaponFeat{weaponFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterWeaponFocusFeatID, state); !ok {
		t.Fatal("expected greater weapon focus to compose from weapon focus with the same selected weapon")
	}
}

func TestNewCharacterFeat_SelectedWeaponFeatOwnershipSatisfiesPlainFeatPrerequisites(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)
	weaponFocus := mustNewCharacterSelectedWeaponFeatForTest(t, characterfeat.WeaponFocusFeatID, selectedWeapon)
	dazzlingDisplay := mustNewCharacterSelectedWeaponFeatForTest(t, characterfeat.DazzlingDisplayFeatID, selectedWeapon)
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponFeatsForTest(
		t,
		nil,
		6,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 6)},
		nil,
		nil,
		selectedWeapon,
		[]CharacterSelectedWeaponFeat{weaponFocus, dazzlingDisplay},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ShatterDefensesFeatID, state); !ok {
		t.Fatal("expected shatter defenses to compose from selected dazzling display and same-weapon focus")
	}
}

func TestNewCharacterFeat_RejectsMismatchedSameSelectionWeaponPrerequisite(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)
	weaponFocus := mustNewCharacterSelectedWeaponFeatForTest(
		t,
		characterfeat.WeaponFocusFeatID,
		mustNewCharacterSelectedWeaponForTest(t, characterequipment.SlingWeaponID),
	)
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponFeatsForTest(
		t,
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		selectedWeapon,
		[]CharacterSelectedWeaponFeat{weaponFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterWeaponFocusFeatID, state); ok {
		t.Fatal("expected greater weapon focus to reject weapon focus with a different selected weapon")
	}
}

func TestNewCharacterFeat_RejectsSameSelectionWeaponPrerequisiteWithoutSelectedFeatFact(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
		t,
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID),
		[]characterfeat.FeatID{characterfeat.WeaponFocusFeatID},
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterWeaponFocusFeatID, state); ok {
		t.Fatal("expected greater weapon focus to reject flat weapon focus without selected weapon ownership")
	}
}

func TestNewCharacterFeat_ComposesSameSelectionSpellSchoolPrerequisite(t *testing.T) {
	selectedSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID)
	spellFocus := mustNewCharacterSelectedSpellSchoolFeatForTest(t, characterfeat.SpellFocusFeatID, selectedSchool)
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		selectedSchool,
		[]CharacterSelectedSpellSchoolFeat{spellFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterSpellFocusFeatID, state); !ok {
		t.Fatal("expected greater spell focus to compose from spell focus with the same selected school")
	}
}

func TestNewCharacterFeat_ComposesSpellSchoolFeatPrerequisite(t *testing.T) {
	spellFocus := mustNewCharacterSelectedSpellSchoolFeatForTest(
		t,
		characterfeat.SpellFocusFeatID,
		mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID),
	)
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		CharacterSelectedSpellSchool{},
		[]CharacterSelectedSpellSchoolFeat{spellFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.AugmentSummoningFeatID, state); !ok {
		t.Fatal("expected augment summoning to compose from spell focus in conjuration")
	}
}

func TestNewCharacterFeat_RejectsMismatchedSpellSchoolPrerequisites(t *testing.T) {
	conjurationSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID)
	evocationSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.EvocationSchoolID)
	spellFocus := mustNewCharacterSelectedSpellSchoolFeatForTest(t, characterfeat.SpellFocusFeatID, conjurationSchool)
	sameSelectionState := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		evocationSchool,
		[]CharacterSelectedSpellSchoolFeat{spellFocus},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterSpellFocusFeatID, sameSelectionState); ok {
		t.Fatal("expected greater spell focus to reject spell focus with a different selected school")
	}

	spellSchoolState := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		CharacterSelectedSpellSchool{},
		[]CharacterSelectedSpellSchoolFeat{
			mustNewCharacterSelectedSpellSchoolFeatForTest(t, characterfeat.SpellFocusFeatID, evocationSchool),
		},
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.AugmentSummoningFeatID, spellSchoolState); ok {
		t.Fatal("expected augment summoning to reject spell focus outside conjuration")
	}
}

func TestNewCharacterFeat_RejectsSpellSchoolPrerequisitesWithoutSelectedSchoolOwnership(t *testing.T) {
	selectedSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID)
	sameSelectionState := mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		selectedSchool,
		nil,
		[]characterfeat.FeatID{characterfeat.SpellFocusFeatID},
	)

	if _, ok := NewCharacterFeat(characterfeat.GreaterSpellFocusFeatID, sameSelectionState); ok {
		t.Fatal("expected greater spell focus to reject flat spell focus without selected school ownership")
	}

	spellSchoolState := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		[]characterfeat.FeatID{characterfeat.SpellFocusFeatID},
	)
	if _, ok := NewCharacterFeat(characterfeat.AugmentSummoningFeatID, spellSchoolState); ok {
		t.Fatal("expected augment summoning to reject flat spell focus without selected school ownership")
	}
}

func TestNewCharacterFeat_RejectsSelectionPrerequisitesWithoutContext(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(t, nil, 1, nil, nil, nil, nil)

	if _, ok := NewCharacterFeat(characterfeat.GreaterSpellFocusFeatID, state); ok {
		t.Fatal("expected same-selection prerequisite to reject without selection context")
	}
}

func TestNewCharacterFeat_ComposesSelectedFamiliarEligibilityPrerequisite(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedFamiliarEligibilityForTest(
		t,
		nil,
		0,
		nil,
		nil,
		[]characterclass.ClassFeatureID{characterclass.FamiliarAccessClassFeatureID},
		nil,
		NewCharacterSelectedFamiliarEligibility(),
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ImprovedFamiliarFeatID, state); !ok {
		t.Fatal("expected improved familiar to compose from familiar access and selected familiar eligibility")
	}
}

func TestNewCharacterFeat_RejectsSelectedFamiliarEligibilityWithoutFamiliarAccess(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateWithSelectedFamiliarEligibilityForTest(
		t,
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		NewCharacterSelectedFamiliarEligibility(),
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ImprovedFamiliarFeatID, state); ok {
		t.Fatal("expected improved familiar to reject eligibility without familiar access")
	}
}

func TestNewCharacterFeat_RejectsSelectedFamiliarEligibilityWithoutEligibilityFact(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		[]characterclass.ClassFeatureID{characterclass.FamiliarAccessClassFeatureID},
		nil,
		nil,
	)

	if _, ok := NewCharacterFeat(characterfeat.ImprovedFamiliarFeatID, state); ok {
		t.Fatal("expected improved familiar to reject familiar access without selected familiar eligibility")
	}
}

func TestCharacterFeatPrerequisiteState_ZeroValueSelectedFamiliarEligibilityFailsClosed(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(
		t,
		nil,
		0,
		nil,
		[]characterclass.ClassFeatureID{characterclass.FamiliarAccessClassFeatureID},
		nil,
		nil,
	)

	if state.SatisfiesPrerequisite(characterfeat.NewSelectedFamiliarEligibilityPrerequisite()) {
		t.Fatal("expected zero-value selected familiar eligibility to fail closed")
	}
}

func TestCharacterFeatPrerequisiteState_RejectsMalformedSelectedWeaponFacts(t *testing.T) {
	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeapon(
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)},
		nil,
		nil,
		characterSelectedWeapon{
			id:                  characterequipment.WeaponID(" dagger"),
			proficiencyCategory: characterequipment.SimpleWeaponProficiencyCategory,
			valid:               true,
		},
		nil,
	); ok {
		t.Fatal("expected malformed selected weapon id to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeapon(
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)},
		nil,
		nil,
		characterSelectedWeapon{
			id:                  characterequipment.WeaponID("longsword"),
			proficiencyCategory: characterequipment.MartialWeaponProficiencyCategory,
			valid:               true,
		},
		nil,
	); ok {
		t.Fatal("expected unknown selected weapon id to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeapon(
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)},
		nil,
		nil,
		characterSelectedWeapon{
			id:                  characterequipment.DaggerWeaponID,
			proficiencyCategory: characterequipment.ExoticWeaponProficiencyCategory,
			valid:               true,
		},
		nil,
	); ok {
		t.Fatal("expected selected weapon fact with mismatched proficiency category to be rejected")
	}
}

func TestCharacterFeatPrerequisiteState_UnsupportedSelectedWeaponMappingFailsClosed(t *testing.T) {
	state := characterFeatPrerequisiteState{
		valid: true,
		classLevels: map[characterclass.ClassID]int{
			characterclass.FighterClassID: 1,
		},
		selectedWeapon: characterSelectedWeapon{
			id:                  characterequipment.WeaponID("net"),
			proficiencyCategory: characterequipment.ExoticWeaponProficiencyCategory,
			valid:               true,
		},
	}

	if state.SatisfiesPrerequisite(characterfeat.NewSelectedWeaponProficiencyPrerequisite()) {
		t.Fatal("expected unsupported selected weapon proficiency mapping to fail closed")
	}
}

func TestNewCharacterSelectedWeaponFeat_RejectsInvalidSelectedFeatFacts(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)

	if _, ok := NewCharacterSelectedWeaponFeat(characterfeat.FeatID("Weapon Focus "), selectedWeapon); ok {
		t.Fatal("expected malformed selected feat id to be rejected")
	}

	if _, ok := NewCharacterSelectedWeaponFeat(characterfeat.FeatID("Laser Focus"), selectedWeapon); ok {
		t.Fatal("expected unknown selected feat id to be rejected")
	}

	if _, ok := NewCharacterSelectedWeaponFeat(characterfeat.SpellFocusFeatID, selectedWeapon); ok {
		t.Fatal("expected non-weapon selected feat id to be rejected")
	}

	if _, ok := NewCharacterSelectedWeaponFeat(characterfeat.WeaponFocusFeatID, CharacterSelectedWeapon{}); ok {
		t.Fatal("expected zero-value selected weapon to be rejected for selected feat ownership")
	}
}

func TestNewCharacterSelectedSpellSchool_RejectsInvalidSpellSchools(t *testing.T) {
	if _, ok := NewCharacterSelectedSpellSchool(characterspell.SchoolID("Chronomancy")); ok {
		t.Fatal("expected unknown selected spell school to be rejected")
	}

	if _, ok := NewCharacterSelectedSpellSchool(characterspell.SchoolID(" conjuration")); ok {
		t.Fatal("expected malformed selected spell school to be rejected")
	}
}

func TestNewCharacterSelectedSpellSchoolFeat_RejectsInvalidSelectedFeatFacts(t *testing.T) {
	selectedSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID)

	if _, ok := NewCharacterSelectedSpellSchoolFeat(characterfeat.FeatID("Spell Focus "), selectedSchool); ok {
		t.Fatal("expected malformed selected feat id to be rejected")
	}

	if _, ok := NewCharacterSelectedSpellSchoolFeat(characterfeat.FeatID("Laser Focus"), selectedSchool); ok {
		t.Fatal("expected unknown selected feat id to be rejected")
	}

	if _, ok := NewCharacterSelectedSpellSchoolFeat(characterfeat.WeaponFocusFeatID, selectedSchool); ok {
		t.Fatal("expected non-spell-school selected feat id to be rejected")
	}

	if _, ok := NewCharacterSelectedSpellSchoolFeat(characterfeat.SpellFocusFeatID, CharacterSelectedSpellSchool{}); ok {
		t.Fatal("expected zero-value selected spell school to be rejected for selected feat ownership")
	}
}

func TestCharacterFeatPrerequisiteState_RejectsDuplicateAndOverlappingSelectedWeaponFeatFacts(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)
	weaponFocus := mustNewCharacterSelectedWeaponFeatForTest(t, characterfeat.WeaponFocusFeatID, selectedWeapon)

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeaponFeats(
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		selectedWeapon,
		[]CharacterSelectedWeaponFeat{weaponFocus, weaponFocus},
		nil,
	); ok {
		t.Fatal("expected duplicate selected weapon feat facts to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeaponFeats(
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		selectedWeapon,
		[]CharacterSelectedWeaponFeat{weaponFocus},
		[]characterfeat.FeatID{characterfeat.WeaponFocusFeatID},
	); ok {
		t.Fatal("expected overlapping flat and selected feat ownership to be rejected")
	}
}

func TestCharacterFeatPrerequisiteState_RejectsDuplicateAndOverlappingSelectedSpellSchoolFeatFacts(t *testing.T) {
	selectedSchool := mustNewCharacterSelectedSpellSchoolForTest(t, characterspell.ConjurationSchoolID)
	spellFocus := mustNewCharacterSelectedSpellSchoolFeatForTest(t, characterfeat.SpellFocusFeatID, selectedSchool)

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeats(
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		selectedSchool,
		[]CharacterSelectedSpellSchoolFeat{spellFocus, spellFocus},
		nil,
	); ok {
		t.Fatal("expected duplicate selected spell-school feat facts to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeats(
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		selectedSchool,
		[]CharacterSelectedSpellSchoolFeat{spellFocus},
		[]characterfeat.FeatID{characterfeat.SpellFocusFeatID},
	); ok {
		t.Fatal("expected overlapping flat and selected spell-school feat ownership to be rejected")
	}
}

func TestCharacterFeatPrerequisiteState_RejectsMalformedSelectedWeaponFeatFacts(t *testing.T) {
	selectedWeapon := mustNewCharacterSelectedWeaponForTest(t, characterequipment.DaggerWeaponID)

	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeaponFeats(
		nil,
		1,
		nil,
		[]CharacterClassLevel{mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 8)},
		nil,
		nil,
		selectedWeapon,
		[]CharacterSelectedWeaponFeat{
			{
				featID: characterfeat.WeaponFocusFeatID,
				selectedWeapon: characterSelectedWeapon{
					id:                  characterequipment.WeaponID(" dagger"),
					proficiencyCategory: characterequipment.SimpleWeaponProficiencyCategory,
					valid:               true,
				},
				valid: true,
			},
		},
		nil,
	); ok {
		t.Fatal("expected malformed selected weapon feat fact to be rejected")
	}
}

func TestCharacterFeatPrerequisiteState_RejectsMalformedSelectedSpellSchoolFeatFacts(t *testing.T) {
	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeats(
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		CharacterSelectedSpellSchool{},
		[]CharacterSelectedSpellSchoolFeat{
			{
				featID: characterfeat.SpellFocusFeatID,
				selectedSpellSchool: characterSelectedSpellSchool{
					id:    characterspell.SchoolID("Chronomancy"),
					valid: true,
				},
				valid: true,
			},
		},
		nil,
	); ok {
		t.Fatal("expected malformed selected spell-school feat fact to be rejected")
	}
}

func TestCharacterFeatPrerequisiteState_RejectsMalformedSelectedFamiliarEligibilityFact(t *testing.T) {
	if _, ok := NewCharacterFeatPrerequisiteStateWithSelectedFamiliarEligibility(
		nil,
		0,
		nil,
		nil,
		[]characterclass.ClassFeatureID{characterclass.FamiliarAccessClassFeatureID},
		nil,
		characterSelectedFamiliarEligibility{
			eligible: false,
			valid:    true,
		},
		nil,
	); ok {
		t.Fatal("expected malformed selected familiar eligibility fact to be rejected")
	}
}

func TestNewCharacterFeatPrerequisiteState_RejectsInvalidEntries(t *testing.T) {
	if _, ok := NewCharacterFeatPrerequisiteState(nil, -1, nil, nil, nil, nil, nil); ok {
		t.Fatal("expected negative base attack bonus to be rejected")
	}

	if _, ok := NewCharacterCasterLevel(ability.CasterSource("Mystic"), 1); ok {
		t.Fatal("expected unknown caster source to be rejected")
	}

	if _, ok := NewCharacterCasterLevel(ability.ArcaneCasterSource, 0); ok {
		t.Fatal("expected zero caster level to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		[]CharacterAbilityScore{{id: ability.AbilityScoreID("LUCK"), score: 10}},
		0,
		nil,
		nil,
		nil,
		nil,
		nil,
	); ok {
		t.Fatal("expected invalid ability score to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		[]CharacterCasterLevel{{source: ability.CasterSource("Mystic"), level: 1}},
		nil,
		nil,
		nil,
		nil,
	); ok {
		t.Fatal("expected invalid caster level to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		[]CharacterClassLevel{{classID: characterclass.ClassID("alchemist"), level: 1}},
		nil,
		nil,
		nil,
	); ok {
		t.Fatal("expected invalid class level to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		nil,
		[]characterclass.ClassFeatureID{characterclass.ClassFeatureID("alchemy")},
		nil,
		nil,
	); ok {
		t.Fatal("expected invalid class feature to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		nil,
		nil,
		[]CharacterSkillRanks{{skillID: skill.SkillID("Sailing"), ranks: 1}},
		nil,
	); ok {
		t.Fatal("expected invalid skill ranks to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		[]characterfeat.FeatID{characterfeat.FeatID("Extra Alchemy")},
	); ok {
		t.Fatal("expected invalid known-feat reference to be rejected")
	}
}

func TestNewCharacterFeatPrerequisiteState_RejectsDuplicateEntries(t *testing.T) {
	strength := mustNewCharacterAbilityScoreForTest(t, ability.StrengthScore, 13)
	if _, ok := NewCharacterFeatPrerequisiteState(
		[]CharacterAbilityScore{strength, strength},
		0,
		nil,
		nil,
		nil,
		nil,
		nil,
	); ok {
		t.Fatal("expected duplicate ability scores to be rejected")
	}

	casterLevel := mustNewCharacterCasterLevelForTest(t, ability.ArcaneCasterSource, 1)
	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		[]CharacterCasterLevel{casterLevel, casterLevel},
		nil,
		nil,
		nil,
		nil,
	); ok {
		t.Fatal("expected duplicate caster levels to be rejected")
	}

	fighterLevel := mustNewCharacterClassLevelForTest(t, characterclass.FighterClassID, 1)
	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		[]CharacterClassLevel{fighterLevel, fighterLevel},
		nil,
		nil,
		nil,
	); ok {
		t.Fatal("expected duplicate class levels to be rejected")
	}

	rideRanks := mustNewCharacterSkillRanksForTest(t, skill.RideSkillID, 1)
	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		nil,
		nil,
		[]CharacterSkillRanks{rideRanks, rideRanks},
		nil,
	); ok {
		t.Fatal("expected duplicate skill ranks to be rejected")
	}

	if _, ok := NewCharacterFeatPrerequisiteState(
		nil,
		0,
		nil,
		nil,
		nil,
		nil,
		[]characterfeat.FeatID{characterfeat.EnduranceFeatID, characterfeat.EnduranceFeatID},
	); ok {
		t.Fatal("expected duplicate feats to be rejected")
	}
}

func TestCharacterFeatPrerequisiteState_ZeroValueDoesNotSatisfyFeat(t *testing.T) {
	var state CharacterFeatPrerequisiteState

	feat, ok := characterfeat.GetFeatByID(characterfeat.AcrobaticFeatID)
	if !ok {
		t.Fatal("expected acrobatic feat seed to resolve")
	}

	if state.SatisfiesFeat(feat) {
		t.Fatal("expected zero-value prerequisite state not to satisfy feats")
	}
}

func TestCharacterFeatPrerequisiteState_DoesNotSatisfyZeroValueFeat(t *testing.T) {
	state := mustNewCharacterFeatPrerequisiteStateForTest(t, nil, 0, nil, nil, nil, nil)
	var feat characterfeat.Feat

	if state.SatisfiesFeat(feat) {
		t.Fatal("expected prerequisite state not to satisfy zero-value feat")
	}
}

func TestCharacterFeat_ZeroValueDoesNotResolve(t *testing.T) {
	var selectedFeat CharacterFeat

	if _, ok := selectedFeat.GetFeat(); ok {
		t.Fatal("expected zero-value character feat not to resolve")
	}
}

func mustNewCharacterAbilityScoreForTest(
	t *testing.T,
	id ability.AbilityScoreID,
	score int,
) CharacterAbilityScore {
	t.Helper()

	value, ok := NewCharacterAbilityScore(id, score)
	if !ok {
		t.Fatalf("expected ability score %q %d to compose", id, score)
	}

	return value
}

func mustNewCharacterClassLevelForTest(
	t *testing.T,
	id characterclass.ClassID,
	level int,
) CharacterClassLevel {
	t.Helper()

	value, ok := NewCharacterClassLevel(id, level)
	if !ok {
		t.Fatalf("expected class level %q %d to compose", id, level)
	}

	return value
}

func mustNewCharacterCasterLevelForTest(
	t *testing.T,
	source ability.CasterSource,
	level int,
) CharacterCasterLevel {
	t.Helper()

	value, ok := NewCharacterCasterLevel(source, level)
	if !ok {
		t.Fatalf("expected caster level %q %d to compose", source, level)
	}

	return value
}

func mustNewCharacterSkillRanksForTest(
	t *testing.T,
	id skill.SkillID,
	ranks int,
) CharacterSkillRanks {
	t.Helper()

	value, ok := NewCharacterSkillRanks(id, ranks)
	if !ok {
		t.Fatalf("expected skill ranks %q %d to compose", id, ranks)
	}

	return value
}

func mustNewCharacterSkillRankAllocationFactsForFeatPrerequisiteTest(
	t *testing.T,
	classID characterclass.ClassID,
	level int,
	intelligenceScore int,
	raceID characterrace.RaceID,
	allocations []CharacterSkillRankAllocation,
) CharacterSkillRankAllocationFacts {
	t.Helper()

	budget := mustNewCharacterSkillRankBudgetFactsForValidationTest(
		t,
		classID,
		level,
		intelligenceScore,
		raceID,
	)

	facts, ok := NewCharacterSkillRankAllocationFacts(budget, allocations)
	if !ok {
		t.Fatal("expected skill-rank allocation facts to compose")
	}

	return facts
}

func mustNewCharacterSkillRankAllocationForFeatPrerequisiteTest(
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

func mustNewCharacterSelectedWeaponFeatForTest(
	t *testing.T,
	id characterfeat.FeatID,
	selectedWeapon CharacterSelectedWeapon,
) CharacterSelectedWeaponFeat {
	t.Helper()

	value, ok := NewCharacterSelectedWeaponFeat(id, selectedWeapon)
	if !ok {
		t.Fatalf("expected selected weapon feat %q to compose", id)
	}

	return value
}

func mustNewCharacterSelectedSpellSchoolForTest(
	t *testing.T,
	id characterspell.SchoolID,
) CharacterSelectedSpellSchool {
	t.Helper()

	value, ok := NewCharacterSelectedSpellSchool(id)
	if !ok {
		t.Fatalf("expected selected spell school %q to compose", id)
	}

	return value
}

func mustNewCharacterSelectedSpellSchoolFeatForTest(
	t *testing.T,
	id characterfeat.FeatID,
	selectedSpellSchool CharacterSelectedSpellSchool,
) CharacterSelectedSpellSchoolFeat {
	t.Helper()

	value, ok := NewCharacterSelectedSpellSchoolFeat(id, selectedSpellSchool)
	if !ok {
		t.Fatalf("expected selected spell-school feat %q to compose", id)
	}

	return value
}

func mustNewCharacterFeatPrerequisiteStateForTest(
	t *testing.T,
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	feats []characterfeat.FeatID,
) CharacterFeatPrerequisiteState {
	t.Helper()

	return mustNewCharacterFeatPrerequisiteStateWithCasterLevelsForTest(
		t,
		abilityScores,
		baseAttackBonus,
		nil,
		classLevels,
		classFeatures,
		skillRanks,
		feats,
	)
}

func mustNewCharacterFeatPrerequisiteStateWithCasterLevelsForTest(
	t *testing.T,
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	feats []characterfeat.FeatID,
) CharacterFeatPrerequisiteState {
	t.Helper()

	return mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
		t,
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		CharacterSelectedWeapon{},
		feats,
	)
}

func mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponForTest(
	t *testing.T,
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedWeapon CharacterSelectedWeapon,
	feats []characterfeat.FeatID,
) CharacterFeatPrerequisiteState {
	t.Helper()

	return mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponFeatsForTest(
		t,
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		selectedWeapon,
		nil,
		feats,
	)
}

func mustNewCharacterFeatPrerequisiteStateWithSelectedWeaponFeatsForTest(
	t *testing.T,
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedWeapon CharacterSelectedWeapon,
	selectedWeaponFeats []CharacterSelectedWeaponFeat,
	feats []characterfeat.FeatID,
) CharacterFeatPrerequisiteState {
	t.Helper()

	state, ok := NewCharacterFeatPrerequisiteStateWithSelectedWeaponFeats(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		selectedWeapon,
		selectedWeaponFeats,
		feats,
	)
	if !ok {
		t.Fatal("expected feat prerequisite state to compose")
	}

	return state
}

func mustNewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeatsForTest(
	t *testing.T,
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedSpellSchool CharacterSelectedSpellSchool,
	selectedSpellSchoolFeats []CharacterSelectedSpellSchoolFeat,
	feats []characterfeat.FeatID,
) CharacterFeatPrerequisiteState {
	t.Helper()

	state, ok := NewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeats(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		selectedSpellSchool,
		selectedSpellSchoolFeats,
		feats,
	)
	if !ok {
		t.Fatal("expected feat prerequisite state to compose")
	}

	return state
}

func mustNewCharacterFeatPrerequisiteStateWithSelectedFamiliarEligibilityForTest(
	t *testing.T,
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedFamiliarEligibility CharacterSelectedFamiliarEligibility,
	feats []characterfeat.FeatID,
) CharacterFeatPrerequisiteState {
	t.Helper()

	state, ok := NewCharacterFeatPrerequisiteStateWithSelectedFamiliarEligibility(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		selectedFamiliarEligibility,
		feats,
	)
	if !ok {
		t.Fatal("expected feat prerequisite state to compose")
	}

	return state
}
