package lore

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/vreid/neroka/internal/common"
)

type filesystemLoreBook struct {
	fs       afero.Fs
	filename string
	entries  map[int64]map[int64]string
}

func NewFilesystemLoreBook(fs afero.Fs, filename string) (LoreBook, error) {
	entries, err := common.LoadData[map[int64]map[int64]string](fs, filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't load entry data: %s", err.Error())
	}

	if entries == nil {
		entries = &map[int64]map[int64]string{}
	}

	result := &filesystemLoreBook{
		fs,
		filename,
		*entries,
	}

	return result, nil
}

func (l *filesystemLoreBook) AddEntry(guildId int64, userId int64, lore string) error {
	_, ok := l.entries[guildId]
	if !ok {
		l.entries[guildId] = map[int64]string{}
	}

	l.entries[guildId][userId] = lore

	err := common.SaveData(l.fs, l.filename, &l.entries)
	if err != nil {
		return err
	}

	return nil
}

func (l *filesystemLoreBook) GetEntry(guildId int64, userId int64) string {
	guildEntries, ok := l.entries[guildId]
	if !ok {
		return ""
	}

	userEntry, ok := guildEntries[userId]
	if !ok {
		return ""
	}

	return userEntry
}

func (l *filesystemLoreBook) RemoveEntry(guildId int64, userId int64) error {
	guildEntries, ok := l.entries[guildId]
	if !ok {
		return nil
	}

	delete(guildEntries, userId)

	err := common.SaveData(l.fs, l.filename, &l.entries)
	if err != nil {
		return err
	}

	return nil
}
