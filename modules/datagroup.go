package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/modules/pools"
	"github.com/containerops/arkor/setting"
)

func GetDataGroups() ([]models.Group, error) {
	var DataGroups []models.Group
	// Read object server configs
	if err := setting.InitObjectServerConf("conf/objectserver.yaml"); err != nil {
		log.Warningf("Read config error: %v", err.Error())
		return nil, err
	}
	// Get DataGroups Information from Registration Center
	rcURI := fmt.Sprintf("http://%s:%s/internal/v1/groups", setting.ObjectServerConf.RegistrationCenter.Address, setting.ObjectServerConf.RegistrationCenter.Port)
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
	err = json.Unmarshal(result, &DataGroups)
	if err != nil {
		return nil, err
	}
	if err := pools.SyncDataServerConnectionPools(DataGroups); err != nil {
		return nil, err
	}
	return DataGroups, nil
}

func SelectDataGroup(groups []models.Group, size int64) (*models.Group, error) {
	// Check if all servers in the Group have enought free space
	indexlist := []int{}
	for index, dg := range groups {
		var find bool = true
		if dg.GroupStatus != models.GROUP_STATUS_NORMAL {
			find = false
		}
		for _, server := range dg.Servers {
			if server.MaxFreeSpace < size {
				find = false
				break
			}
		}
		if find {
			indexlist = append(indexlist, index)
		}
	}
	// Select a available group and return
	if len(indexlist) == 0 {
		return nil, fmt.Errorf("Can not find an available Data Server Group")
	}
	randindex := rand.Int() % len(indexlist)
	outgroup := groups[indexlist[randindex]]
	return &outgroup, nil
}

func GetServersByGroupID(groupID string) ([]models.DataServer, error) {
	var dataservers []models.DataServer
	groups, err := GetDataGroups()
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		if group.ID == groupID {
			dataservers = group.Servers
		}
	}
	return dataservers, nil
}
