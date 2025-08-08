package dm

type DirectMessageManager interface {
	SetDmToggle(userId int64, enabled bool) error
	DmToggleEnabled(userId int64) bool
	GetUsersNeedingCheckUp() []int64
	MarkCheckUpSent(userId int64) error
	SetDmFullHistory(userId int64, enabled bool) error
	FullHistoryEnabled(userId int64) bool
}
