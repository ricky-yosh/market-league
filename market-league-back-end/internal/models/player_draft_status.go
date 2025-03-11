package models

type DraftStatus string

const (
	DraftReady    DraftStatus = "ready"
	DraftNotReady DraftStatus = "not_ready"
)

type PlayerDraftStatus struct {
	PlayerID uint        `json:"player_id"`
	Status   DraftStatus `json:"status"`
}
