package models

import (
	"time"
)

type ObjectMeta struct {
	ID        string     `json:"objectID,omitempty"`
	Key       string     `json:"objectKey,omitempty"`
	Md5Key    string     `json:"md5Key,omitempty"`
	Fragments []Fragment `json:"fragments,omitempty"`
}

type Fragment struct {
	Index   int       `json:"index,omitempty"`
	Start   int64     `json:"start,omitempty"`
	End     int64     `json:"end,omitempty"`
	GroupID string    `json:"groupID,omitempty"`
	FileID  string    `json:"fileId,omitempty"`
	IsLast  bool      `json:"isLast,omitempty"`
	ModTime time.Time `json:"modTime,omitempty"`
}
