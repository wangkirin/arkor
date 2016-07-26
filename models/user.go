package models

// Attention:
// Arkor did not support User Manamgement & Access Control yet

// import (
// 	"github.com/jinzhu/gorm"
// )

type Owner struct {
	ID          string `json:"ID,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
