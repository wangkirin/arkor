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
	Index   int       `json:"index,omitempty" gorm:"column:index"`
	Start   int64     `json:"start,omitempty"`
	End     int64     `json:"end,omitempty"`
	GroupID string    `json:"group_id,omitempty"`
	IsLast  bool      `json:"is_last,omitempty"`
	ModTime time.Time `json:"mod_time,omitempty"`
}

// Object :The relation of ObjectMeta and Fragment
type Object struct {
	ObjectID   string `json:"object_id" gorm:"primary_key"`
	FragmentID string `json:"fragment_id" gorm:"primary_key"`
}

type FragIDConvert struct {
	FragIDstr string `json:"fragIDstr" gorm:"unique"`
	FragIDint int64  `json:"fragIDstr" gorm:"AUTO_INCREMENT"`
}
