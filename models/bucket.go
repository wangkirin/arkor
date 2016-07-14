package models

import (
	"time"
)

type Bucket struct {
	Type        string    `json:"type,omitempty"`
	Name        string    `json:"name,omitempty"`
	MaxKeys     string    `json:"maxKeys,omitempty"`
	KeyCount    string    `json:"keyCount,omitempty"`
	IsTruncated bool      `json:"isTruncated,omitempty"`
	Contents    []Content `json:"contents,omitempty"`
}

type Content struct {
	Key          string    `json:"key,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	ETag         string    `json:"eTag,omitempty"`
	Type         string    `json:"type,omitempty"`
	Size         int64     `json:"size,omitempty"`
	StorageClass string    `json:"storageClass,omitempty"`
	Owner        Owner     `json:"owner,omitempty"`
}

type Owner struct {
	ID          string `json:"ID,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
