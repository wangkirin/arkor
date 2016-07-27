package models

import (
	"time"

	"github.com/containerops/arkor/utils/db/mysql"
	// "github.com/jinzhu/gorm"
)

type Bucket struct {
	Name         string    `json:"name,omitempty"`
	MaxKeys      string    `json:"maxKeys,omitempty"`
	KeyCount     string    `json:"keyCount,omitempty"`
	IsTruncated  bool      `json:"isTruncated,omitempty"`
	CreationDate time.Time `json:"creationDate,omitempty"`
	Contents     []Content `json:"contents,omitempty" gorm:"ForeignKey:Key;AssociationForeignKey:Name"`
	Owner        Owner     `json:"-" gorm:"ForeignKey:ID;AssociationForeignKey:Name"`
}

type Content struct {
	Key          string    `json:"key,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	ETag         string    `json:"eTag,omitempty"`
	Type         string    `json:"type,omitempty"`
	Size         int64     `json:"size,omitempty"`
	StorageClass string    `json:"storageClass,omitempty"`
	Owner        Owner     `json:"owner,omitempty" gorm:"ForeignKey:ID;AssociationForeignKey:Key"`
}

type BucketListResponse struct {
	Type    string         `json:"type,omitempty"`
	Owner   Owner          `json:"owner,omitempty"`
	Buckets []BucketSimple `json:"buckets,omitempty"`
}

type BucketSimple struct {
	Name         string    `json:"name,omitempty"`
	CreationDate time.Time `json:"creationDate,omitempty"`
}

func (b *Bucket) Associate() {
	var bucket Bucket
	var content Content
	var owner Owner
	mysqldb := mysql.MySQLInstance()
	mysqldb.Model(&bucket).Related(&content, "Contents")
	mysqldb.Model(&content).Related(&owner, "Owner")
}
