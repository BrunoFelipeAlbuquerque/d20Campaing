package character

import (
	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	characterclass "d20campaigngenerator/internal/domain/rpg/character/class"
	characterfeat "d20campaigngenerator/internal/domain/rpg/character/feat"
	"d20campaigngenerator/internal/domain/rpg/character/skill"
	characterspell "d20campaigngenerator/internal/domain/rpg/character/spell"
)

type characterAbilityScore struct {
	id    ability.AbilityScoreID
	score int
}
type CharacterAbilityScore = characterAbilityScore

type characterClassLevel struct {
	classID characterclass.ClassID
	level   int
}
type CharacterClassLevel = characterClassLevel

type characterCasterLevel struct {
	source ability.CasterSource
	level  int
}
type CharacterCasterLevel = characterCasterLevel

type characterSkillRanks struct {
	skillID skill.SkillID
	ranks   int
}
type CharacterSkillRanks = characterSkillRanks

type characterFeatPrerequisiteState struct {
	valid                       bool
	abilityScores               map[ability.AbilityScoreID]int
	baseAttackBonus             int
	casterLevels                map[ability.CasterSource]int
	classLevels                 map[characterclass.ClassID]int
	classFeatures               map[characterclass.ClassFeatureID]struct{}
	skillRanks                  map[skill.SkillID]int
	selectedWeapon              characterSelectedWeapon
	selectedWeaponFeats         map[characterfeat.FeatID]characterSelectedWeapon
	selectedSpellSchool         characterSelectedSpellSchool
	selectedSpellSchoolFeats    map[characterfeat.FeatID]characterSelectedSpellSchool
	selectedFamiliarEligibility characterSelectedFamiliarEligibility
	feats                       map[characterfeat.FeatID]struct{}
}
type CharacterFeatPrerequisiteState = characterFeatPrerequisiteState

type characterFeat struct {
	id characterfeat.FeatID
}
type CharacterFeat = characterFeat

func NewCharacterAbilityScore(id ability.AbilityScoreID, score int) (CharacterAbilityScore, bool) {
	value, ok := ability.NewAbilityScoreValue(score, true)
	if !ok {
		return characterAbilityScore{}, false
	}

	if _, ok := ability.NewAbilityScore(id, value); !ok {
		return characterAbilityScore{}, false
	}

	return characterAbilityScore{
		id:    id,
		score: score,
	}, true
}

func NewCharacterClassLevel(
	id characterclass.ClassID,
	level int,
) (CharacterClassLevel, bool) {
	if _, ok := characterclass.GetClassByID(id); !ok || level <= 0 || level > maxCoreCharacterLevel {
		return characterClassLevel{}, false
	}

	return characterClassLevel{
		classID: id,
		level:   level,
	}, true
}

func NewCharacterCasterLevel(source ability.CasterSource, level int) (CharacterCasterLevel, bool) {
	if level <= 0 {
		return characterCasterLevel{}, false
	}

	casterLevel := ability.NewImpossibleCasterLevel()
	if !casterLevel.SetSourceLevel(source, level) {
		return characterCasterLevel{}, false
	}

	return characterCasterLevel{
		source: source,
		level:  level,
	}, true
}

func NewCharacterSkillRanks(id skill.SkillID, ranks int) (CharacterSkillRanks, bool) {
	if ranks <= 0 || !isValidCharacterSkillID(id) {
		return characterSkillRanks{}, false
	}

	return characterSkillRanks{
		skillID: id,
		ranks:   ranks,
	}, true
}

func NewCharacterFeatPrerequisiteState(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	return newCharacterFeatPrerequisiteState(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		characterSelectedWeapon{},
		nil,
		characterSelectedSpellSchool{},
		nil,
		characterSelectedFamiliarEligibility{},
		feats,
	)
}

func NewCharacterFeatPrerequisiteStateWithSkillRankAllocationFacts(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRankAllocations CharacterSkillRankAllocationFacts,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	skillRanks, ok := characterSkillRanksFromAllocationFacts(skillRankAllocations)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	return newCharacterFeatPrerequisiteState(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		characterSelectedWeapon{},
		nil,
		characterSelectedSpellSchool{},
		nil,
		characterSelectedFamiliarEligibility{},
		feats,
	)
}

func NewCharacterFeatPrerequisiteStateWithSelectedWeapon(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedWeapon CharacterSelectedWeapon,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	return newCharacterFeatPrerequisiteState(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		selectedWeapon,
		nil,
		characterSelectedSpellSchool{},
		nil,
		characterSelectedFamiliarEligibility{},
		feats,
	)
}

func NewCharacterFeatPrerequisiteStateWithSelectedWeaponFeats(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedWeapon CharacterSelectedWeapon,
	selectedWeaponFeats []CharacterSelectedWeaponFeat,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	return newCharacterFeatPrerequisiteState(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		selectedWeapon,
		selectedWeaponFeats,
		characterSelectedSpellSchool{},
		nil,
		characterSelectedFamiliarEligibility{},
		feats,
	)
}

func NewCharacterFeatPrerequisiteStateWithSelectedSpellSchoolFeats(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedSpellSchool CharacterSelectedSpellSchool,
	selectedSpellSchoolFeats []CharacterSelectedSpellSchoolFeat,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	return newCharacterFeatPrerequisiteState(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		characterSelectedWeapon{},
		nil,
		selectedSpellSchool,
		selectedSpellSchoolFeats,
		characterSelectedFamiliarEligibility{},
		feats,
	)
}

func NewCharacterFeatPrerequisiteStateWithSelectedFamiliarEligibility(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedFamiliarEligibility CharacterSelectedFamiliarEligibility,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	return newCharacterFeatPrerequisiteState(
		abilityScores,
		baseAttackBonus,
		casterLevels,
		classLevels,
		classFeatures,
		skillRanks,
		characterSelectedWeapon{},
		nil,
		characterSelectedSpellSchool{},
		nil,
		selectedFamiliarEligibility,
		feats,
	)
}

func newCharacterFeatPrerequisiteState(
	abilityScores []CharacterAbilityScore,
	baseAttackBonus int,
	casterLevels []CharacterCasterLevel,
	classLevels []CharacterClassLevel,
	classFeatures []characterclass.ClassFeatureID,
	skillRanks []CharacterSkillRanks,
	selectedWeapon CharacterSelectedWeapon,
	selectedWeaponFeats []CharacterSelectedWeaponFeat,
	selectedSpellSchool CharacterSelectedSpellSchool,
	selectedSpellSchoolFeats []CharacterSelectedSpellSchoolFeat,
	selectedFamiliarEligibility CharacterSelectedFamiliarEligibility,
	feats []characterfeat.FeatID,
) (CharacterFeatPrerequisiteState, bool) {
	if baseAttackBonus < 0 {
		return characterFeatPrerequisiteState{}, false
	}

	abilityScoreMap, ok := buildCharacterAbilityScoreMap(abilityScores)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	casterLevelMap, ok := buildCharacterCasterLevelMap(casterLevels)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	classLevelMap, ok := buildCharacterClassLevelMap(classLevels)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	classFeatureSet, ok := buildCharacterClassFeatureSet(classFeatures)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	skillRankMap, ok := buildCharacterSkillRankMap(skillRanks)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	selectedWeaponValue, ok := buildCharacterSelectedWeapon(selectedWeapon)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	selectedWeaponFeatMap, ok := buildCharacterSelectedWeaponFeatMap(selectedWeaponFeats)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	selectedSpellSchoolValue, ok := buildCharacterSelectedSpellSchool(selectedSpellSchool)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	spellSchoolFeatMap, ok := buildCharacterSelectedSpellSchoolFeatMap(selectedSpellSchoolFeats)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	selectedFamiliarEligibilityValue, ok := buildCharacterSelectedFamiliarEligibility(selectedFamiliarEligibility)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	featSet, ok := buildCharacterFeatSet(feats)
	if !ok {
		return characterFeatPrerequisiteState{}, false
	}

	if hasOverlappingFeatOwnership(featSet, selectedWeaponFeatMap, spellSchoolFeatMap) {
		return characterFeatPrerequisiteState{}, false
	}

	return characterFeatPrerequisiteState{
		valid:                       true,
		abilityScores:               abilityScoreMap,
		baseAttackBonus:             baseAttackBonus,
		casterLevels:                casterLevelMap,
		classLevels:                 classLevelMap,
		classFeatures:               classFeatureSet,
		skillRanks:                  skillRankMap,
		selectedWeapon:              selectedWeaponValue,
		selectedWeaponFeats:         selectedWeaponFeatMap,
		selectedSpellSchool:         selectedSpellSchoolValue,
		selectedSpellSchoolFeats:    spellSchoolFeatMap,
		selectedFamiliarEligibility: selectedFamiliarEligibilityValue,
		feats:                       featSet,
	}, true
}

func NewCharacterFeat(
	id characterfeat.FeatID,
	prerequisites CharacterFeatPrerequisiteState,
) (CharacterFeat, bool) {
	feat, ok := characterfeat.GetFeatByID(id)
	if !ok || !prerequisites.SatisfiesFeat(feat) {
		return characterFeat{}, false
	}

	return characterFeat{id: id}, true
}

func (a characterAbilityScore) GetAbilityScoreID() ability.AbilityScoreID {
	return a.id
}

func (a characterAbilityScore) GetScore() int {
	return a.score
}

func (l characterClassLevel) GetClassID() characterclass.ClassID {
	return l.classID
}

func (l characterClassLevel) GetLevel() int {
	return l.level
}

func (l characterCasterLevel) GetSource() ability.CasterSource {
	return l.source
}

func (l characterCasterLevel) GetLevel() int {
	return l.level
}

func (r characterSkillRanks) GetSkillID() skill.SkillID {
	return r.skillID
}

func (r characterSkillRanks) GetRanks() int {
	return r.ranks
}

func (s characterFeatPrerequisiteState) SatisfiesFeat(feat characterfeat.Feat) bool {
	if !s.valid {
		return false
	}

	if _, ok := characterfeat.GetFeatByID(feat.GetID()); !ok {
		return false
	}

	for _, prerequisite := range feat.GetPrerequisites() {
		if !s.SatisfiesPrerequisite(prerequisite) {
			return false
		}
	}

	return true
}

func (s characterFeatPrerequisiteState) SatisfiesPrerequisite(
	prerequisite characterfeat.Prerequisite,
) bool {
	if !s.valid || prerequisite == nil {
		return false
	}

	switch value := prerequisite.(type) {
	case characterfeat.AbilityScorePrerequisite:
		score, ok := s.abilityScores[value.GetAbilityScoreID()]
		return ok && score >= value.GetMinimumScore()
	case characterfeat.BaseAttackBonusPrerequisite:
		return s.baseAttackBonus >= value.GetMinimumBonus()
	case characterfeat.SkillRanksPrerequisite:
		return s.skillRanks[value.GetSkillID()] >= value.GetMinimumRanks()
	case characterfeat.AnySkillRanksPrerequisite:
		for _, skillID := range value.GetSkillIDs() {
			if s.skillRanks[skillID] >= value.GetMinimumRanks() {
				return true
			}
		}
		return false
	case characterfeat.SpellcastingPrerequisite:
		return s.satisfiesSpellcastingAccess(value.GetAccess())
	case characterfeat.CasterLevelPrerequisite:
		return s.satisfiesCasterLevel(value.GetMinimumLevel())
	case characterfeat.CharacterLevelPrerequisite:
		return s.totalCharacterLevel() >= value.GetMinimumLevel()
	case characterfeat.ClassLevelPrerequisite:
		return s.classLevels[value.GetClassID()] >= value.GetMinimumLevel()
	case characterfeat.ClassFeaturePrerequisite:
		_, ok := s.classFeatures[value.GetFeatureID()]
		return ok
	case characterfeat.SelectedWeaponProficiencyPrerequisite:
		return s.satisfiesSelectedWeaponProficiency()
	case characterfeat.SelectedFamiliarEligibilityPrerequisite:
		return s.selectedFamiliarEligibility.IsEligible()
	case characterfeat.SameSelectionFeatPrerequisite:
		return s.satisfiesSameSelectedWeaponFeat(value.GetFeatID()) ||
			s.satisfiesSameSelectedSpellSchoolFeat(value.GetFeatID())
	case characterfeat.SpellSchoolFeatPrerequisite:
		return s.satisfiesSpellSchoolFeat(value.GetFeatID(), value.GetSchoolID())
	case characterfeat.FeatPrerequisite:
		return s.hasFeat(value.GetFeatID())
	case characterfeat.AnyFeatPrerequisite:
		for _, featID := range value.GetFeatIDs() {
			if s.hasFeat(featID) {
				return true
			}
		}
		return false
	case characterfeat.FeatCategoryCountPrerequisite:
		return s.featCategoryCount(value.GetCategory()) >= value.GetMinimumCount()
	default:
		return false
	}
}

func (s characterFeatPrerequisiteState) hasFeat(id characterfeat.FeatID) bool {
	if _, ok := s.feats[id]; ok {
		return true
	}

	if _, ok := s.selectedWeaponFeats[id]; ok {
		return true
	}

	_, ok := s.selectedSpellSchoolFeats[id]
	return ok
}

func (f characterFeat) GetFeatID() characterfeat.FeatID {
	return f.id
}

func (f characterFeat) GetFeat() (characterfeat.Feat, bool) {
	return characterfeat.GetFeatByID(f.id)
}

func (s characterFeatPrerequisiteState) totalCharacterLevel() int {
	total := 0
	for _, level := range s.classLevels {
		total += level
	}

	return total
}

func (s characterFeatPrerequisiteState) satisfiesSpellcastingAccess(
	access characterfeat.SpellcastingAccess,
) bool {
	for classID := range s.classLevels {
		class, ok := characterclass.GetClassByID(classID)
		if !ok {
			return false
		}

		spellcasting := class.GetSpellcasting()
		if !spellcasting.HasSpellcasting() {
			continue
		}

		switch access {
		case characterfeat.AnySpellcastingAccess:
			return true
		case characterfeat.ArcaneSpellcastingAccess:
			if isArcaneSpellcastingKind(spellcasting.GetKind()) {
				return true
			}
		case characterfeat.DivineSpellcastingAccess:
			if isDivineSpellcastingKind(spellcasting.GetKind()) {
				return true
			}
		default:
			return false
		}
	}

	return false
}

func (s characterFeatPrerequisiteState) satisfiesCasterLevel(minimumLevel int) bool {
	if minimumLevel <= 0 {
		return false
	}

	for _, level := range s.casterLevels {
		if level >= minimumLevel {
			return true
		}
	}

	return false
}

func (s characterFeatPrerequisiteState) satisfiesSelectedWeaponProficiency() bool {
	if !s.selectedWeapon.valid {
		return false
	}

	for classID := range s.classLevels {
		class, ok := characterclass.GetClassByID(classID)
		if !ok {
			return false
		}

		if s.selectedWeapon.IsProficientWith(class.GetWeaponProficiencies()) {
			return true
		}
	}

	return false
}

func (s characterFeatPrerequisiteState) satisfiesSameSelectedWeaponFeat(
	featID characterfeat.FeatID,
) bool {
	if !s.selectedWeapon.valid {
		return false
	}

	previousSelection, ok := s.selectedWeaponFeats[featID]
	if !ok {
		return false
	}

	return previousSelection.valid && previousSelection.id == s.selectedWeapon.id
}

func (s characterFeatPrerequisiteState) satisfiesSameSelectedSpellSchoolFeat(
	featID characterfeat.FeatID,
) bool {
	if !s.selectedSpellSchool.valid {
		return false
	}

	previousSelection, ok := s.selectedSpellSchoolFeats[featID]
	if !ok {
		return false
	}

	return previousSelection.valid && previousSelection.id == s.selectedSpellSchool.id
}

func (s characterFeatPrerequisiteState) satisfiesSpellSchoolFeat(
	featID characterfeat.FeatID,
	schoolID characterspell.SchoolID,
) bool {
	previousSelection, ok := s.selectedSpellSchoolFeats[featID]
	if !ok {
		return false
	}

	return previousSelection.valid && previousSelection.id == schoolID
}

func (s characterFeatPrerequisiteState) featCategoryCount(
	category characterfeat.FeatCategory,
) int {
	count := 0
	for featID := range s.feats {
		feat, ok := characterfeat.GetFeatByID(featID)
		if !ok {
			return 0
		}

		if feat.GetCategory() == category {
			count++
		}
	}

	for featID := range s.selectedWeaponFeats {
		feat, ok := characterfeat.GetFeatByID(featID)
		if !ok {
			return 0
		}

		if feat.GetCategory() == category {
			count++
		}
	}

	for featID := range s.selectedSpellSchoolFeats {
		feat, ok := characterfeat.GetFeatByID(featID)
		if !ok {
			return 0
		}

		if feat.GetCategory() == category {
			count++
		}
	}

	return count
}

func buildCharacterAbilityScoreMap(
	values []CharacterAbilityScore,
) (map[ability.AbilityScoreID]int, bool) {
	result := make(map[ability.AbilityScoreID]int, len(values))

	for _, value := range values {
		if _, ok := NewCharacterAbilityScore(value.id, value.score); !ok {
			return nil, false
		}

		if _, ok := result[value.id]; ok {
			return nil, false
		}

		result[value.id] = value.score
	}

	return result, true
}

func buildCharacterCasterLevelMap(
	values []CharacterCasterLevel,
) (map[ability.CasterSource]int, bool) {
	result := make(map[ability.CasterSource]int, len(values))

	for _, value := range values {
		if _, ok := NewCharacterCasterLevel(value.source, value.level); !ok {
			return nil, false
		}

		if _, ok := result[value.source]; ok {
			return nil, false
		}

		result[value.source] = value.level
	}

	return result, true
}

func buildCharacterClassLevelMap(
	values []CharacterClassLevel,
) (map[characterclass.ClassID]int, bool) {
	result := make(map[characterclass.ClassID]int, len(values))

	for _, value := range values {
		if _, ok := NewCharacterClassLevel(value.classID, value.level); !ok {
			return nil, false
		}

		if _, ok := result[value.classID]; ok {
			return nil, false
		}

		result[value.classID] = value.level
	}

	return result, true
}

func buildCharacterClassFeatureSet(
	values []characterclass.ClassFeatureID,
) (map[characterclass.ClassFeatureID]struct{}, bool) {
	result := make(map[characterclass.ClassFeatureID]struct{}, len(values))

	for _, value := range values {
		if value.GetName() == "" {
			return nil, false
		}

		if _, ok := result[value]; ok {
			return nil, false
		}

		result[value] = struct{}{}
	}

	return result, true
}

func buildCharacterSkillRankMap(values []CharacterSkillRanks) (map[skill.SkillID]int, bool) {
	result := make(map[skill.SkillID]int, len(values))

	for _, value := range values {
		if _, ok := NewCharacterSkillRanks(value.skillID, value.ranks); !ok {
			return nil, false
		}

		if _, ok := result[value.skillID]; ok {
			return nil, false
		}

		result[value.skillID] = value.ranks
	}

	return result, true
}

func characterSkillRanksFromAllocationFacts(
	facts CharacterSkillRankAllocationFacts,
) ([]CharacterSkillRanks, bool) {
	if !isValidCharacterSkillRankAllocationFacts(facts) {
		return nil, false
	}

	allocations := facts.GetAllocations()
	skillRanks := make([]CharacterSkillRanks, 0, len(allocations))
	groupedFamilyRanks := make(map[skill.SkillID]int)

	for _, allocation := range allocations {
		allocatedRanks, ok := NewCharacterSkillRanks(allocation.GetSkillID(), allocation.GetRanks())
		if !ok {
			return nil, false
		}

		skillRanks = append(skillRanks, allocatedRanks)

		allocatedSkill, ok := allocation.GetSkill()
		if !ok {
			return nil, false
		}

		if _, ok := allocatedSkill.GetSpecialization(); !ok {
			continue
		}

		familyID := allocatedSkill.GetFamilyID()
		if allocation.GetRanks() > groupedFamilyRanks[familyID] {
			groupedFamilyRanks[familyID] = allocation.GetRanks()
		}
	}

	for _, coreSkill := range skill.GetSkills() {
		ranks, ok := groupedFamilyRanks[coreSkill.GetID()]
		if !ok {
			continue
		}

		familyRanks, ok := NewCharacterSkillRanks(coreSkill.GetID(), ranks)
		if !ok {
			return nil, false
		}

		skillRanks = append(skillRanks, familyRanks)
	}

	return skillRanks, true
}

func buildCharacterSelectedWeapon(
	value CharacterSelectedWeapon,
) (characterSelectedWeapon, bool) {
	if isEmptyCharacterSelectedWeapon(value) {
		return characterSelectedWeapon{}, true
	}

	if !value.valid {
		return characterSelectedWeapon{}, false
	}

	selectedWeapon, ok := NewCharacterSelectedWeapon(value.id)
	if !ok {
		return characterSelectedWeapon{}, false
	}

	if value.proficiencyCategory != selectedWeapon.proficiencyCategory {
		return characterSelectedWeapon{}, false
	}

	return selectedWeapon, true
}

func buildCharacterSelectedWeaponFeatMap(
	values []CharacterSelectedWeaponFeat,
) (map[characterfeat.FeatID]characterSelectedWeapon, bool) {
	result := make(map[characterfeat.FeatID]characterSelectedWeapon, len(values))

	for _, value := range values {
		selectedFeat, ok := NewCharacterSelectedWeaponFeat(value.featID, value.selectedWeapon)
		if !ok {
			return nil, false
		}

		if _, ok := result[selectedFeat.featID]; ok {
			return nil, false
		}

		result[selectedFeat.featID] = selectedFeat.selectedWeapon
	}

	return result, true
}

func buildCharacterSelectedSpellSchoolFeatMap(
	values []CharacterSelectedSpellSchoolFeat,
) (map[characterfeat.FeatID]characterSelectedSpellSchool, bool) {
	result := make(map[characterfeat.FeatID]characterSelectedSpellSchool, len(values))

	for _, value := range values {
		selectedFeat, ok := NewCharacterSelectedSpellSchoolFeat(value.featID, value.selectedSpellSchool)
		if !ok {
			return nil, false
		}

		if _, ok := result[selectedFeat.featID]; ok {
			return nil, false
		}

		result[selectedFeat.featID] = selectedFeat.selectedSpellSchool
	}

	return result, true
}

func hasOverlappingFeatOwnership(
	feats map[characterfeat.FeatID]struct{},
	selectedWeaponFeats map[characterfeat.FeatID]characterSelectedWeapon,
	selectedSpellSchoolFeats map[characterfeat.FeatID]characterSelectedSpellSchool,
) bool {
	selectedFeatIDs := make(map[characterfeat.FeatID]struct{}, len(selectedWeaponFeats))
	for featID := range selectedWeaponFeats {
		if _, ok := feats[featID]; ok {
			return true
		}

		selectedFeatIDs[featID] = struct{}{}
	}

	for featID := range selectedSpellSchoolFeats {
		if _, ok := feats[featID]; ok {
			return true
		}

		if _, ok := selectedFeatIDs[featID]; ok {
			return true
		}
	}

	return false
}

func isEmptyCharacterSelectedWeapon(value CharacterSelectedWeapon) bool {
	return !value.valid && value.id == "" && value.proficiencyCategory == ""
}

func buildCharacterFeatSet(
	values []characterfeat.FeatID,
) (map[characterfeat.FeatID]struct{}, bool) {
	result := make(map[characterfeat.FeatID]struct{}, len(values))

	for _, value := range values {
		if _, ok := characterfeat.GetFeatByID(value); !ok {
			return nil, false
		}

		if _, ok := result[value]; ok {
			return nil, false
		}

		result[value] = struct{}{}
	}

	return result, true
}

func isValidCharacterSkillID(id skill.SkillID) bool {
	if _, ok := skill.GetSkillByID(id); ok {
		return true
	}

	return skill.IsValidSkillID(id)
}

func isArcaneSpellcastingKind(kind characterclass.SpellcastingKind) bool {
	switch kind {
	case characterclass.ArcanePreparedSpellcastingKind,
		characterclass.ArcaneSpontaneousSpellcastingKind:
		return true
	default:
		return false
	}
}

func isDivineSpellcastingKind(kind characterclass.SpellcastingKind) bool {
	return kind == characterclass.DivinePreparedSpellcastingKind
}
