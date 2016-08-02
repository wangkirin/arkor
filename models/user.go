package models

// Attention:
// Arkor did not support User Manamgement & Access Control yet

// import (
// 	"github.com/jinzhu/gorm"
// )

type Owner struct {
	BucketName  string `json:"-"`
	ContentKey  string `json:"-"`
	ID          string `json:"id,omitempty" gorm:"column:ID"`
	DisplayName string `json:"display_name,omitempty"`
}
