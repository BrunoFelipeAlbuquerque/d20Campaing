package language

type languageID string
type LanguageID = languageID

type language struct {
	id     languageID
	secret bool
}
type Language = language

const (
	AbyssalLanguageID     LanguageID = "Abyssal"
	AkloLanguageID        LanguageID = "Aklo"
	AquanLanguageID       LanguageID = "Aquan"
	AuranLanguageID       LanguageID = "Auran"
	CelestialLanguageID   LanguageID = "Celestial"
	CommonLanguageID      LanguageID = "Common"
	DraconicLanguageID    LanguageID = "Draconic"
	DruidicLanguageID     LanguageID = "Druidic"
	DwarvenLanguageID     LanguageID = "Dwarven"
	ElvenLanguageID       LanguageID = "Elven"
	GiantLanguageID       LanguageID = "Giant"
	GnomeLanguageID       LanguageID = "Gnome"
	GoblinLanguageID      LanguageID = "Goblin"
	GnollLanguageID       LanguageID = "Gnoll"
	HalflingLanguageID    LanguageID = "Halfling"
	IgnanLanguageID       LanguageID = "Ignan"
	InfernalLanguageID    LanguageID = "Infernal"
	OrcLanguageID         LanguageID = "Orc"
	SylvanLanguageID      LanguageID = "Sylvan"
	TerranLanguageID      LanguageID = "Terran"
	UndercommonLanguageID LanguageID = "Undercommon"
)

var coreLanguages = mustBuildCoreLanguages()

var coreLanguageOrder = []LanguageID{
	AbyssalLanguageID,
	AkloLanguageID,
	AquanLanguageID,
	AuranLanguageID,
	CelestialLanguageID,
	CommonLanguageID,
	DraconicLanguageID,
	DruidicLanguageID,
	DwarvenLanguageID,
	ElvenLanguageID,
	GiantLanguageID,
	GnomeLanguageID,
	GoblinLanguageID,
	GnollLanguageID,
	HalflingLanguageID,
	IgnanLanguageID,
	InfernalLanguageID,
	OrcLanguageID,
	SylvanLanguageID,
	TerranLanguageID,
	UndercommonLanguageID,
}

func GetLanguageByID(id LanguageID) (Language, bool) {
	value, ok := coreLanguages[id]
	if !ok {
		return language{}, false
	}

	return value, true
}

func GetLanguages() []Language {
	languages := make([]Language, 0, len(coreLanguageOrder))

	for _, id := range coreLanguageOrder {
		languages = append(languages, coreLanguages[id])
	}

	return languages
}

func IsKnownLanguageID(id LanguageID) bool {
	_, ok := coreLanguages[id]
	return ok
}

func IsSecretLanguageID(id LanguageID) bool {
	value, ok := coreLanguages[id]
	return ok && value.secret
}

func IsNonSecretLanguageID(id LanguageID) bool {
	value, ok := coreLanguages[id]
	return ok && !value.secret
}

func (l language) GetID() LanguageID {
	return l.id
}

func (l language) IsSecret() bool {
	return l.secret
}

func mustBuildCoreLanguages() map[LanguageID]Language {
	result := make(map[LanguageID]Language, len(coreLanguageOrder))

	for _, id := range coreLanguageOrder {
		result[id] = language{id: id, secret: id == DruidicLanguageID}
	}

	return result
}
