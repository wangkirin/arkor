package models

import (
	"fmt"
	"testing"

	"github.com/containerops/arkor/utils/db"
	_ "github.com/containerops/arkor/utils/db/mysql"
	_ "github.com/containerops/arkor/utils/db/redis"
	. "github.com/containerops/arkor/utils/setting"
)

func Test_DataserverSQLCreate(t *testing.T) {
	ds := &DataServer{
		ID:      "testID02",
		GroupID: "1",
		IP:      "10.67.147.80",
	}

	if err := db.SelectSQLDriver(RunTime.Sqldatabase.Driver); err != nil {
		fmt.Printf("Register database driver error: %s\n", err.Error())
	} else {
		DBuri := fmt.Sprintf("%s:%s", RunTime.Sqldatabase.Host, RunTime.Sqldatabase.Port)
		err := db.SQLDB.InitDB(RunTime.Sqldatabase.Driver, RunTime.Sqldatabase.Username, RunTime.Sqldatabase.Password, DBuri, RunTime.Sqldatabase.Schema, 0)
		if err != nil {
			fmt.Printf("Connect database error: %s\n", err.Error())
		}
		db.SQLDB.RegisterModel(&DataServer{})
	}

	if err := db.SQLDB.Create(ds); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

	if err := db.SQLDB.Delete(ds); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

}

func Test_DataserverKVCreate(t *testing.T) {
	ds := &DataServer{
		ID:      "testkv",
		GroupID: "1",
		IP:      "10.67.147.80",
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

	if err := db.KVDB.Create(ds); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

	if err := db.KVDB.Delete(ds); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

}
