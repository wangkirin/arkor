package models

import (
	"time"
)

// status of Data Server
const (
	INIT_STATUS = 0
	RW_STATUS   = 1
	RO_STATUS   = 2
	ERR_STATUS  = 3
)

// struct of DataServer
type DataServer struct {
	ID             string    `json:"data_server_id,omitempty" gorm:"unique;column:id"`
	GroupID        string    `json:"group_id,omitempty" gorm:"column:group_id" binding:"Required"`
	IP             string    `json:"ip,omitempty" gorm:"column:ip" binding:"Required"`
	Port           int       `json:"port,omitempty" binding:"Required"`
	Status         int       `json:"status,omitempty"`
	Deleted        int       `json:"deleted,omitempty"`
	TotalChunks    int       `json:"total_chunks,omitempty"`
	TotalFreeSpace int64     `json:"total_free_space,omitempty"`
	MaxFreeSpace   int64     `json:"max_free_space,omitempty"`
	DataPath       string    `json:"data_path,omitempty"`
	PendingWrites  int       `json:"pend_writes,omitempty"`
	ReadingCount   int64     `json:"reading_count,omitempty"`
	ConnCounts     int       `json:"conn_counts,omitempty"`
	CreateTime     time.Time `json:"create_time,omitempty"`
	UpdateTime     time.Time `json:"update_time,omitempty"`
}
