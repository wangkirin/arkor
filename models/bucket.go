package models

import (
	"time"

	"github.com/containerops/arkor/utils/db/mysql"
	// "github.com/jinzhu/gorm"
)

type Bucket struct {
	Name         string    `json:"name,omitempty" gorm:"unique"`
	MaxKeys      string    `json:"maxKeys,omitempty" gorm:"-"`
	KeyCount     string    `json:"keyCount,omitempty" gorm:"-"`
	IsTruncated  bool      `json:"isTruncated,omitempty" gorm:"-"`
	CreationDate time.Time `json:"creationDate,omitempty"`
	Contents     []Content `json:"contents,omitempty" gorm:"ForeignKey:BucketName;AssociationForeignKey:Name"`
	Owner        Owner     `json:"-" gorm:"ForeignKey:BucketName;AssociationForeignKey:Name"`
}

type Content struct {
	BucketName   string    `json:"-"`
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"eTag"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
	StorageClass string    `json:"storageClass"`
	Owner        Owner     `json:"owner" gorm:"ForeignKey:ContentKey;AssociationForeignKey:Key"`
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
	mysqldb.Model(&bucket).Related(&owner, "Owner")
	mysqldb.Model(&content).Related(&owner, "Owner")
}
