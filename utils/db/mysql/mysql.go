package mysql

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/containerops/arkor/utils/db/factory"
)

type mysql struct{}

func init() {
	factory.RegisterSQL("mysql", &mysql{})
}

var db *gorm.DB

func (my *mysql) InitDB(driver, user, passwd, uri, name string, partition int64) error {
	var err error
	databaseuri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, uri, name)
	if db, err = gorm.Open(driver, databaseuri); err != nil {
		log.Errorln(err.Error())
		log.Fatal("Initlization database connection error.")
		os.Exit(1)
	} else {
		db.DB()
		db.DB().Ping()
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		db.SingularTable(true)
	}
	return err
}

func (my *mysql) RegisterModel(datastructs ...interface{}) error {
	for _, datastruct := range datastructs {
		if !db.HasTable(datastruct) {
			if result := db.CreateTable(datastruct); result.Error != nil {
				log.Infoln("Create Table Error")
				return result.Error
			}
		}
		db.AutoMigrate(datastruct)
	}
	return nil
}

// Create insert records to database
func (my *mysql) Create(value interface{}) error {
	if result := db.Create(value); result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete a record
func (my *mysql) Delete(value interface{}) error {
	if result := db.Delete(value); result.Error != nil {
		return result.Error
	}
	return nil
}

// Save update the record, and if the record does not exist ,insert it
func (my *mysql) Save(value interface{}) error {
	if result := db.Save(value); result.Error != nil {
		return result.Error
	}
	return nil
}

// Query one record
func (my *mysql) Query(value interface{}) (bool, error) {
	if result := db.Where(value).Find(value); result.Error != nil && strings.EqualFold(result.Error.Error(), "record not found") {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	} else if result.RowsAffected > 1 {
		return true, fmt.Errorf("query records more than one")
	}
	return true, nil
}
