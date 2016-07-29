package inner

import (
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

	ds := &models.DataServer{
		ID:         utils.MD5ID(),
		GroupID:    req.GroupID,
		IP:         req.IP,
		Port:       req.Port,
		CreateTime: now,
		UpdateTime: now,
	}

	// Query DataServer ID from SQL Database
	if err := db.SQLDB.Create(ds); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	log.Println(ds)
	// Save dataserver status to K/V Database
	if err := db.KVDB.Create(ds); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// Save dataserver status to SQL Database

	return http.StatusOK, nil
}
