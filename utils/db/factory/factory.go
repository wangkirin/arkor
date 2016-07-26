package factory

import (
	"fmt"
)

type DriverFactory interface {
	RegisterModel(models ...interface{}) error
	InitDB(driver, user, passwd, uri, name string, partition int64) error
	Create(obj interface{}) error
	Delete(obj interface{}) error
	Query(obj interface{}) (bool, error)
	QueryMulti(condition interface{}, value interface{}) (bool, error)
	Save(obj interface{}) error
}

var SQLDrivers = make(map[string]DriverFactory)
var KVDrivers = make(map[string]DriverFactory)

func RegisterSQL(name string, instance DriverFactory) error {
	if _, existed := SQLDrivers[name]; existed {
		return fmt.Errorf("%v has already been registered", name)
	}

	SQLDrivers[name] = instance
	return nil
}

func RegisterKV(name string, instance DriverFactory) error {
	if _, existed := KVDrivers[name]; existed {
		return fmt.Errorf("%v has already been registered", name)
	}

	KVDrivers[name] = instance
	return nil
}
