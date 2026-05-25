package class

import (
	"strings"

	ability "d20campaigngenerator/internal/domain/rpg/character/ability"
	"d20campaigngenerator/internal/domain/rpg/character/skill"
)

type classID string
type ClassID = classID

type weaponProficiencyID string
type WeaponProficiencyID = weaponProficiencyID

type armorProficiencyID string
type ArmorProficiencyID = armorProficiencyID

type spellcastingKind string
type SpellcastingKind = spellcastingKind

const (
	NonSpellcastingKind               SpellcastingKind = "None"
	ArcanePreparedSpellcastingKind    SpellcastingKind = "Arcane Prepared"
	ArcaneSpontaneousSpellcastingKind SpellcastingKind = "Arcane Spontaneous"
	DivinePreparedSpellcastingKind    SpellcastingKind = "Divine Prepared"
)

type saveProgressions struct {
	fortitude ability.SavingThrowProgression
	reflex    ability.SavingThrowProgression
	will      ability.SavingThrowProgression
}
type SaveProgressions = saveProgressions

type spellcastingProfile struct {
	kind       spellcastingKind
	keyAbility ability.AbilityScoreID
}
type SpellcastingProfile = spellcastingProfile

type class struct {
	id                  classID
	hitDieType          ability.HitDieType
	baseAttackBonus     ability.BaseAttackBonusProgression
	saveProgressions    saveProgressions
	skillRanksPerLevel  int
	classSkills         []skill.SkillID
	weaponProficiencies []weaponProficiencyID
	armorProficiencies  []armorProficiencyID
	spellcasting        spellcastingProfile
}
type Class = class

func NewSaveProgressions(
	fortitude ability.SavingThrowProgression,
	reflex ability.SavingThrowProgression,
	will ability.SavingThrowProgression,
) (SaveProgressions, bool) {
	value := saveProgressions{
		fortitude: fortitude,
		reflex:    reflex,
		will:      will,
	}
	if !isValidSaveProgressions(value) {
		return saveProgressions{}, false
	}

	return value, true
}

func NewNonSpellcastingProfile() SpellcastingProfile {
	return spellcastingProfile{
		kind: NonSpellcastingKind,
	}
}

func NewSpellcastingProfile(
	kind SpellcastingKind,
	keyAbility ability.AbilityScoreID,
) (SpellcastingProfile, bool) {
	value := spellcastingProfile{
		kind:       kind,
		keyAbility: keyAbility,
	}
	if !isValidCasterSpellcastingKind(kind) || !isValidSpellcastingProfile(value) {
		return spellcastingProfile{}, false
	}

	return value, true
}

func NewClass(
	id ClassID,
	hitDieType ability.HitDieType,
	baseAttackBonus ability.BaseAttackBonusProgression,
	saveProgressions SaveProgressions,
	skillRanksPerLevel int,
	classSkills []skill.SkillID,
	weaponProficiencies []WeaponProficiencyID,
	armorProficiencies []ArmorProficiencyID,
	spellcasting SpellcastingProfile,
) (Class, bool) {
	if !isValidClassID(id) ||
		!isValidHitDieType(hitDieType) ||
		!isValidBaseAttackBonusProgression(baseAttackBonus) ||
		!isValidSaveProgressions(saveProgressions) ||
		skillRanksPerLevel <= 0 ||
		!isValidSpellcastingProfile(spellcasting) {
		return class{}, false
	}

	dedupedClassSkills, ok := dedupeClassSkills(classSkills)
	if !ok {
		return class{}, false
	}

	dedupedWeaponProficiencies, ok := dedupeWeaponProficiencies(weaponProficiencies)
	if !ok {
		return class{}, false
	}

	dedupedArmorProficiencies, ok := dedupeArmorProficiencies(armorProficiencies)
	if !ok {
		return class{}, false
	}

	return class{
		id:                  id,
		hitDieType:          hitDieType,
		baseAttackBonus:     baseAttackBonus,
		saveProgressions:    saveProgressions,
		skillRanksPerLevel:  skillRanksPerLevel,
		classSkills:         dedupedClassSkills,
		weaponProficiencies: dedupedWeaponProficiencies,
		armorProficiencies:  dedupedArmorProficiencies,
		spellcasting:        spellcasting,
	}, true
}

func (c class) GetID() ClassID {
	return c.id
}

func (c class) GetHitDieType() ability.HitDieType {
	return c.hitDieType
}

func (c class) GetBaseAttackBonusProgression() ability.BaseAttackBonusProgression {
	return c.baseAttackBonus
}

func (c class) GetSaveProgressions() SaveProgressions {
	return c.saveProgressions
}

func (c class) GetSkillRanksPerLevel() int {
	return c.skillRanksPerLevel
}

func (c class) GetClassSkills() []skill.SkillID {
	return append([]skill.SkillID(nil), c.classSkills...)
}

func (c class) GetWeaponProficiencies() []WeaponProficiencyID {
	return append([]WeaponProficiencyID(nil), c.weaponProficiencies...)
}

func (c class) GetArmorProficiencies() []ArmorProficiencyID {
	return append([]ArmorProficiencyID(nil), c.armorProficiencies...)
}

func (c class) GetSpellcasting() SpellcastingProfile {
	return c.spellcasting
}

func (p saveProgressions) GetFortitude() ability.SavingThrowProgression {
	return p.fortitude
}

func (p saveProgressions) GetReflex() ability.SavingThrowProgression {
	return p.reflex
}

func (p saveProgressions) GetWill() ability.SavingThrowProgression {
	return p.will
}

func (p saveProgressions) GetProgression(id ability.SavingThrowID) (ability.SavingThrowProgression, bool) {
	switch id {
	case ability.FortitudeSave:
		return p.fortitude, true
	case ability.ReflexSave:
		return p.reflex, true
	case ability.WillSave:
		return p.will, true
	default:
		return "", false
	}
}

func (p spellcastingProfile) GetKind() SpellcastingKind {
	return p.kind
}

func (p spellcastingProfile) HasSpellcasting() bool {
	return p.kind != NonSpellcastingKind
}

func (p spellcastingProfile) GetKeyAbility() (ability.AbilityScoreID, bool) {
	if p.kind == NonSpellcastingKind {
		return "", false
	}

	return p.keyAbility, true
}

func isValidClassID(id ClassID) bool {
	value := string(id)
	return value != "" && strings.TrimSpace(value) == value
}

func isValidHitDieType(hitDieType ability.HitDieType) bool {
	_, ok := ability.NewUniformHitDie(hitDieType, 1)
	return ok
}

func isValidBaseAttackBonusProgression(progression ability.BaseAttackBonusProgression) bool {
	_, ok := ability.NewBaseAttackBonusByClassLevel(1, progression)
	return ok
}

func isValidSaveProgressions(value SaveProgressions) bool {
	return isValidSavingThrowProgression(value.fortitude) &&
		isValidSavingThrowProgression(value.reflex) &&
		isValidSavingThrowProgression(value.will)
}

func isValidSavingThrowProgression(value ability.SavingThrowProgression) bool {
	_, ok := ability.NewSavingThrowByClassLevel(ability.FortitudeSave, 1, value)
	return ok
}

func isValidSpellcastingProfile(value SpellcastingProfile) bool {
	if !isValidSpellcastingKind(value.kind) {
		return false
	}

	if value.kind == NonSpellcastingKind {
		return value.keyAbility == ""
	}

	return value.keyAbility.GetName() != ""
}

func isValidSpellcastingKind(value SpellcastingKind) bool {
	switch value {
	case NonSpellcastingKind,
		ArcanePreparedSpellcastingKind,
		ArcaneSpontaneousSpellcastingKind,
		DivinePreparedSpellcastingKind:
		return true
	default:
		return false
	}
}

func isValidCasterSpellcastingKind(value SpellcastingKind) bool {
	return value != NonSpellcastingKind && isValidSpellcastingKind(value)
}

func dedupeClassSkills(classSkills []skill.SkillID) ([]skill.SkillID, bool) {
	if len(classSkills) == 0 {
		return nil, false
	}

	seen := make(map[skill.SkillID]struct{}, len(classSkills))
	deduped := make([]skill.SkillID, 0, len(classSkills))

	for _, classSkillID := range classSkills {
		if !isValidClassSkillID(classSkillID) {
			return nil, false
		}

		if _, ok := seen[classSkillID]; ok {
			continue
		}

		seen[classSkillID] = struct{}{}
		deduped = append(deduped, classSkillID)
	}

	return deduped, true
}

func isValidClassSkillID(id skill.SkillID) bool {
	if _, ok := skill.GetSkillByID(id); ok {
		return true
	}

	return skill.IsValidSkillID(id)
}

func dedupeWeaponProficiencies(
	weaponProficiencies []WeaponProficiencyID,
) ([]WeaponProficiencyID, bool) {
	if len(weaponProficiencies) == 0 {
		return nil, false
	}

	seen := make(map[WeaponProficiencyID]struct{}, len(weaponProficiencies))
	deduped := make([]WeaponProficiencyID, 0, len(weaponProficiencies))

	for _, proficiencyID := range weaponProficiencies {
		if !isValidWeaponProficiencyID(proficiencyID) {
			return nil, false
		}

		if _, ok := seen[proficiencyID]; ok {
			continue
		}

		seen[proficiencyID] = struct{}{}
		deduped = append(deduped, proficiencyID)
	}

	return deduped, true
}

func dedupeArmorProficiencies(
	armorProficiencies []ArmorProficiencyID,
) ([]ArmorProficiencyID, bool) {
	if len(armorProficiencies) == 0 {
		return nil, true
	}

	seen := make(map[ArmorProficiencyID]struct{}, len(armorProficiencies))
	deduped := make([]ArmorProficiencyID, 0, len(armorProficiencies))

	for _, proficiencyID := range armorProficiencies {
		if !isValidArmorProficiencyID(proficiencyID) {
			return nil, false
		}

		if _, ok := seen[proficiencyID]; ok {
			continue
		}

		seen[proficiencyID] = struct{}{}
		deduped = append(deduped, proficiencyID)
	}

	return deduped, true
}
