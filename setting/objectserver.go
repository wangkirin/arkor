package setting

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

var (
	ObjectServerConf *ObjectServer
)

type ObjectServer struct {
	RegistrationCenter *RegistrationCenter "yaml:registrationcenter"
}

type RegistrationCenter struct {
	ListenMode string "yaml:listenMode"
	Address    string "yaml:address"
	Port       string "yaml:port"
}

func InitObjectServerConf(Path string) error {
	// If Config already init , return
	if ObjectServerConf.RegistrationCenter.Address != "" {
		return nil
	}
	// else read related config file
	conf, err := ioutil.ReadFile(Path)
	if err != nil {
		return err
	}
	ObjectServerConf = &ObjectServer{}

	err = yaml.Unmarshal([]byte(conf), &ObjectServerConf)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}
