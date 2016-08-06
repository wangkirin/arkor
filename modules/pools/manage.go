package pools

import (
	"fmt"

	"github.com/containerops/arkor/models"
)

const (
	CONNECTION_POOL_CAPACITY = 200
)

func SyncDataServerConnectionPools(groups []models.Group) error {
	if ConnectionPools == nil {
		ConnectionPools = NewDataServerConnectionPools()
	}
	serverMap := make(map[string]models.DataServer)
	// Add new pool
	for _, group := range groups {
		for _, server := range group.Servers {
			key := fmt.Sprintf("%s:%s", server.IP, server.Port)
			serverMap[key] = server
			_, ok := ConnectionPools.Pools[key]
			if !ok {
				ConnectionPools.AddPool(&server, CONNECTION_POOL_CAPACITY)
			}
		}
	}
	// Remove deleted pool
	for k, _ := range ConnectionPools.Pools {
		server2remove, ok := serverMap[k]
		if !ok {
			ConnectionPools.RemovePool(&server2remove)
		}
	}
	return nil
}
