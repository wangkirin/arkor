package models

import (
	"fmt"
	"testing"

	"github.com/containerops/arkor/utils/db"
)

func Test_DataserverSQLCreate(t *testing.T) {
	ds := &DataServer{
		ID:      "testID05",
		GroupID: "1",
		IP:      "10.67.147.80",
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

	if err := db.KVDB.Create(ds); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

	if err := db.KVDB.Delete(ds); err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}

}
