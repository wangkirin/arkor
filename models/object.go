package models

import (
	"time"
)

type ObjectMeta struct {
	ID        string     `json:"object_id,omitempty" gorm:"column:id;primary_key;unique"`
	Key       string     `json:"object_key,omitempty" gorm:"column:object_key"`
	Md5Key    string     `json:"md5_key,omitempty" gorm:"column:md5_key`
	Fragments []Fragment `json:"fragments,omitempty"`
}

type Fragment struct {
	ID      string    `json:"id" gorm:"column:id;primary_key;unique"`
	Index   int       `json:"index" gorm:"column:index"`
	Start   int64     `json:"start"`
	End     int64     `json:"end"`
	GroupID string    `json:"group_id"`
	FileID  string    `json:"file_id"`
	IsLast  bool      `json:"is_last"`
	ModTime time.Time `json:"mod_time"`
}

// Object :The relation of ObjectMeta and Fragment
type Object struct {
	ObjectID   string `json:"object_id" gorm:"primary_key"`
	FragmentID string `json:"fragment_id" gorm:"primary_key"`
}

type FragIDConvert struct {
	FragIDstr string `json:"fragIDstr" gorm:"unique"`
	FragIDint int64  `json:"fragIDstr" gorm:"primary_key;AUTO_INCREMENT"`
}
