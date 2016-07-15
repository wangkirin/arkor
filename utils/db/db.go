package db

import (
	"fmt"

	"github.com/containerops/arkor/utils/db/factory"
	_ "github.com/containerops/arkor/utils/db/mysql"
	_ "github.com/containerops/arkor/utils/db/redis"
)

var SQLDB factory.DriverFactory
var KVDB factory.DriverFactory

func SelectSQLDriver(name string) error {
	if SQLDB != nil {
		return fmt.Errorf("Only support one SQL DB driver at one time")
	}

	for k, v := range factory.SQLDrivers {
		if k == name && v != nil {
			SQLDB = factory.SQLDrivers[k]
			return nil
		}
	}
	return fmt.Errorf("Not support driver %v", name)
}

func SelectKVDriver(name string) error {
	if KVDB != nil {
		return fmt.Errorf("Only support one KV DB driver at one time")
	}

	for k, v := range factory.KVDrivers {
		if k == name && v != nil {
			KVDB = factory.KVDrivers[k]
			return nil
		}
	}
	return fmt.Errorf("Not support driver %v", name)
}
