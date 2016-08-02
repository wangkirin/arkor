package models

import "time"

// status of a Group
const (
	GROUP_STATUS_NORMAL   = 1
	GROUP_STATUS_UNNORMAL = 2
)

type Group struct {
	ID          string       `json:"group_id,omitempty"`
	GroupStatus int          `json:"group_status,omitempty"`
	Servers     []DataServer `json:"servers,omitempty"`
}

type GroupServer struct {
	ID          int    `json:"id,omitempty" gorm:"primary_key:true;AUTO_INCREMENT"`
	GroupStatus int    `json:"group_status,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	ServerID    string `json:"server_id,omitempty"`
}

/*
A struct to receve the group-dataserver relations info in SQL
*/
type GroupServerInfo struct {
	GroupID        string
	GropuStatus    int
	ServerID       string
	IP             string    `json:"ip"`
	Status         int       `json:"status"`
	Port           int       `json:"port"`
	GroupStatus    int       `json:"status"`
	Deleted        int       `json:"deleted"`
	TotalChunks    int       `json:"total_chunks"`
	TotalFreeSpace int64     `json:"total_free_space"`
	MaxFreeSpace   int64     `json:"max_free_space"`
	DataPath       string    `json:"data_path"`
	PendingWrites  int       `json:"pend_writes"`
	ReadingCount   int64     `json:"reading_count"`
	ConnCounts     int       `json:"conn_counts"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
}
