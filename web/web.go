package web

import (
	"fmt"

	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/middleware"
	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/router"
	. "github.com/containerops/arkor/setting"
	"github.com/containerops/arkor/utils/db"
)

func SetArkorMacaron(m *macaron.Macaron) {
	// Setting Middleware
	middleware.SetMiddlewares(m)
	// Setting Router
	router.SetRouters(m)
	// Static
	if RunTime.Run.RunMode == "dev" {
		m.Use(macaron.Static("external"))
	}
	// Init SQL DB
	InitSQLDB()
	// Init Key/Value DB
	InitKVDB()

}

func SetObjectServerMacaron(m *macaron.Macaron) {
	// Setting Middleware
	middleware.SetMiddlewares(m)
	// Setting Router
	router.SetObjectServerRouters(m)
	// Static
	if RunTime.Run.RunMode == "dev" {
		m.Use(macaron.Static("external"))
	}
	// Init SQL DB
	InitSQLDB()
	// Init Key/Value DB
	InitKVDB()
}

func SetRegistrationCenterMacaron(m *macaron.Macaron) {
	// Setting Middleware
	middleware.SetMiddlewares(m)
	// Setting Router
	router.SetRegistrationCenterRouters(m)
	// Static
	if RunTime.Run.RunMode == "dev" {
		m.Use(macaron.Static("external"))
	}
	// Init SQL DB
	InitSQLDB()
	// Init Key/Value DB
	InitKVDB()
}

func InitSQLDB() {
	if err := db.SelectSQLDriver(RunTime.Sqldatabase.Driver); err != nil {
		fmt.Printf("Register database driver error: %s\n", err.Error())
	} else {
		DBuri := fmt.Sprintf("%s:%s", RunTime.Sqldatabase.Host, RunTime.Sqldatabase.Port)
		err := db.SQLDB.InitDB(RunTime.Sqldatabase.Driver, RunTime.Sqldatabase.Username, RunTime.Sqldatabase.Password, DBuri, RunTime.Sqldatabase.Schema, 0)
		if err != nil {
			fmt.Printf("Connect database error: %s\n", err.Error())
		}
		db.SQLDB.RegisterModel(&models.DataServer{}, &models.Bucket{}, &models.Content{}, &models.Owner{}, &models.GroupServer{}, &models.Object{}, &models.ObjectMeta{}, &models.Fragment{}, &models.FragIDConvert{})
	}
}

func InitKVDB() {
	if err := db.SelectKVDriver(RunTime.Kvdatabase.Driver); err != nil {
		fmt.Printf("Register database driver error: %s\n", err.Error())
	} else {
		DBuri := fmt.Sprintf("%s:%s", RunTime.Kvdatabase.Host, RunTime.Kvdatabase.Port)
		err := db.KVDB.InitDB(RunTime.Kvdatabase.Driver, RunTime.Kvdatabase.Username, RunTime.Kvdatabase.Password, DBuri, RunTime.Kvdatabase.Schema, RunTime.Kvdatabase.Partition)
		if err != nil {
			fmt.Printf("Connect database error: %s\n", err.Error())
		}
		db.KVDB.RegisterModel(&models.DataServer{})
	}
}
