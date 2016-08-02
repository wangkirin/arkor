package modules

import (
	"fmt"
	"testing"

	"github.com/containerops/arkor/models"
)

func Test_SelectDataGroup(t *testing.T) {

	server1 := models.DataServer{
		ID:           "testServer1",
		GroupID:      "1",
		IP:           "10.229.40.140",
		Port:         7456,
		Status:       models.RW_STATUS,
		MaxFreeSpace: 20,
	}
	server2 := models.DataServer{
		ID:           "testServer2",
		GroupID:      "1",
		IP:           "10.229.40.140",
		Port:         7457,
		Status:       models.RW_STATUS,
		MaxFreeSpace: 1222200,
	}
	server3 := models.DataServer{
		ID:           "testServer3",
		GroupID:      "1",
		IP:           "10.229.40.140",
		Port:         7458,
		Status:       models.RW_STATUS,
		MaxFreeSpace: 1222200,
	}
	group1 := models.Group{
		ID:          "testgroup1",
		GroupStatus: models.GROUP_STATUS_NORMAL,
		Servers:     []models.DataServer{server1, server2, server3},
	}
	server4 := models.DataServer{
		ID:           "testServer4",
		GroupID:      "2",
		IP:           "10.229.40.120",
		Port:         7456,
		Status:       models.RW_STATUS,
		MaxFreeSpace: 1222200,
	}
	server5 := models.DataServer{
		ID:           "testServer5",
		GroupID:      "2",
		IP:           "10.229.40.120",
		Port:         7457,
		Status:       models.RW_STATUS,
		MaxFreeSpace: 20,
	}
	server6 := models.DataServer{
		ID:           "testServer6",
		GroupID:      "2",
		IP:           "10.229.40.120",
		Port:         7458,
		Status:       models.RW_STATUS,
		MaxFreeSpace: 1222200,
	}
	group2 := models.Group{
		ID:          "testgroup2",
		GroupStatus: models.GROUP_STATUS_NORMAL,
		Servers:     []models.DataServer{server4, server5, server6},
	}
	groups := []models.Group{group1, group2}
	groupselect, err := SelectDataGroup(groups, 1200)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf(err.Error())
	}
	fmt.Println(groupselect)
}
