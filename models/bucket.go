package models

import (
	"time"
)

type Bucket struct {
	Type        string    `json:"Type,omitempty"`
	Name        string    `json:"Name,omitempty"`
	MaxKeys     string    `json:"MaxKeys,omitempty"`
	IsTruncated bool      `json:"IsTruncated,omitempty"`
	Contents    []Content `json:"Contents,omitempty"`
}

type Content struct {
	Key          string    `json:"Key,omitempty"`
	LastModified time.Time `json:"LastModified,omitempty"`
	ETag         string    `json:"ETag,omitempty"`
	Type         string    `json:"Type,omitempty"`
	Size         int64     `json:"Size,omitempty"`
	StorageClass string    `json:"StorageClass,omitempty"`
	Owner        Owner     `json:"Owner,omitempty"`
}

type Owner struct {
	ID          string `json:"ID,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
}
