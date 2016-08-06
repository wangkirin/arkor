package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/setting"
	"github.com/containerops/arkor/utils/db"
)

func GetObjectInfo(objectid string) (*models.ObjectMeta, error) {
	var objectmeta models.ObjectMeta
	// Read object server configs
	if err := setting.InitObjectServerConf("conf/objectserver.yaml"); err != nil {
		log.Warningf("Read config error: %v", err.Error())
		return nil, err
	}
	// Get ObjectMeta Information from Registration Center
	rcURI := fmt.Sprintf("http://%s:%s/internal/v1/object/%s", setting.ObjectServerConf.RegistrationCenter.Address, setting.ObjectServerConf.RegistrationCenter.Port, objectid)
	resp, err := http.Get(rcURI)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	err = json.Unmarshal(result, &objectmeta)
	if err != nil {
		return nil, err
	}
	return &objectmeta, nil
}

func SaveObjectInfo(objectmetainfo models.ObjectMeta) error {
	// Read object server configs
	if err := setting.InitObjectServerConf("conf/objectserver.yaml"); err != nil {
		log.Warningf("Read config error: %v", err.Error())
		return err
	}
	objectjson, err := json.Marshal(objectmetainfo)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer([]byte(objectjson))
	// Get ObjectMeta Information from Registration Center
	rcURI := fmt.Sprintf("http://%s:%s/internal/v1/object/info", setting.ObjectServerConf.RegistrationCenter.Address, setting.ObjectServerConf.RegistrationCenter.Port)
	resp, err := http.Post(rcURI, "application/json", body)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Save Object Info Error, Status Code: %d, Info: %s", resp.StatusCode, result)
	}
	return nil
}

func FragIDStr2Int(idstr string) (int64, error) {
	convertstruct := models.FragIDConvert{
		FragIDstr: idstr,
	}
	if exist, err := db.SQLDB.Query(&convertstruct); !exist && err == nil {
		if err := db.SQLDB.Create(&convertstruct); err != nil {
			return -1, err
		}
		if exist, err := db.SQLDB.Query(&convertstruct); exist && err == nil {
			return convertstruct.FragIDint, nil
		}
	} else if err != nil {
		return -1, err
	}
	return convertstruct.FragIDint, nil
}
