package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/setting"
	"github.com/containerops/arkor/utils/db"
	"github.com/containerops/arkor/utils/db/mysql"
)

func GetDataGroups() ([]models.Group, err) {
	var DataGroups []models.Group
	// Get DataGroups Information from
	rcURI := fmt.Sprintf("http://%s:%s/internal/v1/groups", setting.ObjectServer.RegistrationCenter.Address, setting.ObjectServer.RegistrationCenter.Port)
	resp, err := http.Get(rcURI)
	if err != nil {
		return nil, err
	}
	err := json.Unmarshal(resp, &DataGroups)
	if err != nil {
		return nil, err
	}
	return DataGroups, nil
}

func SelectDataGroup(groups models.Group, size int64) (models.Group, error) {
	for index, dg := range groups {
		var find bool = true
		for _, server := range groups.Servers {
			if server.MaxFreeSpace < size {
				find = false
				break
			}
		}
		if find {
			resultGroupId = groupId
		}
	}

	if resultGroupId != "" {
		return groups.GroupMap[resultGroupId], nil
	}
	return nil, fmt.Errorf("can not find an available chunkserver")
}
