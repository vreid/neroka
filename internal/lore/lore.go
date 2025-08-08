package lore

type LoreBook interface {
	AddEntry(guildId int64, userId int64, lore string) error
	GetEntry(guildId int64, userId int64) string
	RemoveEntry(guildId int64, userId int64) error
}
