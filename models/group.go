package models

// status of a Group
const (
	GROUP_STATUS_NORMAL   = 1
	GROUP_STATUS_UNNORMAL = 2
)

type Group struct {
	ID          string       `json:"groupID,omitempty"`
	GroupStatus int          `json:"globalStatus,omitempty"`
	Servers     []DataServer `json:"servers,omitempty"`
}
