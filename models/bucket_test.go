package models

import (
	"fmt"
	"testing"
	"time"

	. "github.com/containerops/arkor/setting"
	"github.com/containerops/arkor/utils/db"
)

func init() {
	if err := InitConf("../conf/global.yaml", "../conf/runtime.yaml"); err != nil {
		fmt.Errorf("Read config error: %v", err.Error())
		return
	}
	if err := db.SelectSQLDriver(RunTime.Sqldatabase.Driver); err != nil {
		fmt.Printf("Register database driver error: %s\n", err.Error())
	} else {
		DBuri := fmt.Sprintf("%s:%s", RunTime.Sqldatabase.Host, RunTime.Sqldatabase.Port)
		err := db.SQLDB.InitDB(RunTime.Sqldatabase.Driver, RunTime.Sqldatabase.Username, RunTime.Sqldatabase.Password, DBuri, RunTime.Sqldatabase.Schema, 0)
		if err != nil {
			fmt.Printf("Connect database error: %s\n", err.Error())
		}
		db.SQLDB.RegisterModel(&DataServer{}, &Bucket{}, &Content{}, &Owner{})
	}

	if err := db.SelectKVDriver(RunTime.Kvdatabase.Driver); err != nil {
		fmt.Printf("Register database driver error: %s\n", err.Error())
	} else {
		DBuri := fmt.Sprintf("%s:%s", RunTime.Kvdatabase.Host, RunTime.Kvdatabase.Port)
		err := db.KVDB.InitDB(RunTime.Kvdatabase.Driver, RunTime.Kvdatabase.Username, RunTime.Kvdatabase.Password, DBuri, RunTime.Kvdatabase.Schema, 0)
		if err != nil {
			fmt.Printf("Connect database error: %s\n", err.Error())
		}
		db.KVDB.RegisterModel(&DataServer{})
	}
}

func Test_BucketCreate(t *testing.T) {

	owner := Owner{
		ID:          "arkor",
		DisplayName: "containerops-arkor",
	}
	ct1 := Content{
		Key:          "jpeg1",
		ETag:         "abcdefg",
		Type:         "pic",
		LastModified: time.Now(),
		Owner:        owner,
	}

	ct2 := Content{
		Key:          "jpeg2",
		ETag:         "abcdefg",
		LastModified: time.Now(),
		Type:         "video",
		Owner:        owner,
	}

	contents := []Content{ct1, ct2}
	buckettest := &Bucket{
		Name:         "bucket7",
		CreationDate: time.Now(),
		Contents:     contents,
		Owner:        owner,
	}

	if err := db.SQLDB.Create(buckettest); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

}
