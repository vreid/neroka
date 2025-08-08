package dm

import (
	"fmt"
	"time"

	"github.com/spf13/afero"
	"github.com/vreid/neroka/internal/common"
)

type _settings struct {
	DmToggleSettings map[int64]bool      `json:"dm_toggle_settings"`
	LastInteractions map[int64]time.Time `json:"last_interactions"`
	PendingCheckUps  map[int64]bool      `json:"pending_check_ups"`
	DmPersonalities  map[int64]any       `json:"dm_personalities"`
	DmFullHistory    map[int64]bool      `json:"dm_full_history"`
	CheckUpSent      map[int64]bool      `json:"check_up_sent"`
}

type filesystemDirectMessageManager struct {
	fs       afero.Fs
	filename string

	settings _settings
}

func NewFilesystemDirectMessageManager(fs afero.Fs, filename string) (DirectMessageManager, error) {
	settings, err := common.LoadData[_settings](fs, filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't load entry data: %s", err.Error())
	}

	if settings == nil {
		settings = &_settings{
			DmToggleSettings: map[int64]bool{},
			LastInteractions: map[int64]time.Time{},
			PendingCheckUps:  map[int64]bool{},
			DmPersonalities:  map[int64]any{},
			DmFullHistory:    map[int64]bool{},
			CheckUpSent:      map[int64]bool{},
		}
	}

	result := &filesystemDirectMessageManager{
		fs,
		filename,
		*settings,
	}

	return result, nil
}

func (d *filesystemDirectMessageManager) updateLastInteractions(userId int64) {
	d.settings.LastInteractions[userId] = time.Now().UTC()
	delete(d.settings.PendingCheckUps, userId)
	d.settings.CheckUpSent[userId] = false
}

func (d *filesystemDirectMessageManager) SetDmToggle(userId int64, enabled bool) error {
	d.settings.DmToggleSettings[userId] = enabled
	if enabled {
		d.updateLastInteractions(userId)
	} else {
		delete(d.settings.PendingCheckUps, userId)
	}

	err := common.SaveData(d.fs, d.filename, &d.settings)
	if err != nil {
		return err
	}

	return nil
}

func (d *filesystemDirectMessageManager) DmToggleEnabled(userId int64) bool {
	if enabled, ok := d.settings.DmToggleSettings[userId]; ok {
		return enabled
	}

	return false
}

func (d *filesystemDirectMessageManager) GetUsersNeedingCheckUp() []int64 {
	usersNeedingCheckUp := []int64{}

	sixHoursAgo := time.Now().UTC().Add(time.Duration(-6) * time.Hour)
	for userId, enabled := range d.settings.DmToggleSettings {
		if !enabled {
			continue
		}

		lastInteraction, ok := d.settings.LastInteractions[userId]
		if !ok {
			continue
		}

		if pending, ok := d.settings.PendingCheckUps[userId]; ok && pending {
			continue
		}

		if sent, ok := d.settings.CheckUpSent[userId]; ok && sent {
			continue
		}

		if lastInteraction.After(sixHoursAgo) {
			continue
		}

		usersNeedingCheckUp = append(usersNeedingCheckUp, userId)
	}

	return usersNeedingCheckUp
}

func (d *filesystemDirectMessageManager) MarkCheckUpSent(userId int64) error {
	d.settings.PendingCheckUps[userId] = true
	d.settings.CheckUpSent[userId] = true

	err := common.SaveData(d.fs, d.filename, &d.settings)
	if err != nil {
		return err
	}

	return nil
}

func (d *filesystemDirectMessageManager) SetDmFullHistory(userId int64, enabled bool) error {
	d.settings.DmFullHistory[userId] = enabled

	err := common.SaveData(d.fs, d.filename, &d.settings)
	if err != nil {
		return err
	}

	return nil
}

func (d *filesystemDirectMessageManager) FullHistoryEnabled(userId int64) bool {
	if enabled, ok := d.settings.DmFullHistory[userId]; ok {
		return enabled
	}

	return false
}
