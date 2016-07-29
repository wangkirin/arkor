package inner

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/utils"
	"github.com/containerops/arkor/utils/db"
)

func PutDataserverHandler(ctx *macaron.Context, req models.DataServer, log *logrus.Logger) (int, []byte) {
	ds := &models.DataServer{
		IP:   req.IP,
		Port: req.Port,
	}

	// Query DataServer ID from SQL Database
	if exist, err := db.SQLDB.Query(ds); !exist && err == nil {
		return http.StatusNotFound, []byte("Data server have NOT registered")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	req.ID = ds.ID
	req.UpdateTime = time.Now()
	log.Println(req)
	// Save dataserver status to K/V Database
	if err := db.KVDB.Save(&req); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// Save dataserver status to SQL Database
	if err := db.SQLDB.Save(&req); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, nil
}

func AddDataserverHandler(ctx *macaron.Context, req models.DataServer, log *logrus.Logger) (int, []byte) {
	now := time.Now()
	ServerID := utils.MD5ID()

	ds := &models.DataServer{
		ID:         ServerID,
		GroupID:    req.GroupID,
		IP:         req.IP,
		Port:       req.Port,
		CreateTime: now,
		UpdateTime: now,
	}

	// Create dataserver in SQL Database
	if err := db.SQLDB.Create(ds); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	log.Println(ds)
	// Save dataserver info to K/V Database as cache
	if err := db.KVDB.Create(ds); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	gs := &models.GroupServer{
		GroupID:  req.GroupID,
		ServerID: ServerID,
	}

	if err := db.SQLDB.Create(gs); err != nil {
		return http.StatusInsufficientStorage, []byte(err.Error())
	}

	return http.StatusOK, nil
}

func GetGroupsHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {

	relations := []models.GroupServer{}

	if _, err := db.SQLDB.QueryMulti(&models.GroupServer{}, &relations); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	groups := make(map[string][]string)

	for _, relation := range relations {
		if _, exists := groups[relation.GroupID]; exists == true {
			groups[relation.GroupID] = append(groups[relation.GroupID], relation.ServerID)
		} else {
			groups[relation.GroupID] = []string{relation.ServerID}
		}
	}

	result, err := json.Marshal(groups)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	ctx.Resp.Header().Set("Content-Type", "application/json")

	return http.StatusOK, result
}
