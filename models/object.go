package models

import (
	"time"
)

type ObjectMeta struct {
	ID        string     `json:"object_id,omitempty"`
	Key       string     `json:"object_key,omitempty"`
	Md5Key    string     `json:"md5Key,omitempty"`
	Fragments []Fragment `json:"fragments,omitempty"`
}

type Fragment struct {
	Index   int       `json:"index,omitempty"`
	Start   int64     `json:"start,omitempty"`
	End     int64     `json:"end,omitempty"`
	GroupID string    `json:"group_id,omitempty"`
	FileID  string    `json:"file_id,omitempty"`
	IsLast  bool      `json:"is_last,omitempty"`
	ModTime time.Time `json:"mod_time,omitempty"`
}
