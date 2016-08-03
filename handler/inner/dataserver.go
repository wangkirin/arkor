package inner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/utils"
	"github.com/containerops/arkor/utils/db"
	"github.com/jinzhu/gorm"
)

func PutDataserverHandler(ctx *macaron.Context, req models.DataServer, log *logrus.Logger) (int, []byte) {
	ds := &models.DataServer{
		IP:   req.IP,
		Port: req.Port,
	}

	// Query DataServer ID from SQL Database
	if exist, err := db.SQLDB.Query(ds); !exist && err == nil {
		return http.StatusNotFound, []byte("Data server is NOT registered")
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

func AddDataserverHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	data, _ := ctx.Req.Body().Bytes()
	dataServers := []models.DataServer{}
	json.Unmarshal(data, &dataServers)

	if len(dataServers) == 0 {
		return http.StatusBadRequest, []byte("Invalid Parameters or Incorrect json content")
	}

	insertDataServerSql := "INSERT INTO data_server (id, group_id, ip, port, create_time, update_time) VALUES "
	insertGroupServerSql := "INSERT INTO group_server (group_id, server_id) VALUES "

	resultAry := []interface{}{}
	for _, dataServer := range dataServers {

		now := time.Now()
		nowStr := now.Format("2006-01-02 15:04:05") // I don't know what the writer of Go is thinking of!
		serverID := utils.MD5ID()

		insertDataServer := fmt.Sprintf("(%q, %q, %q, %d, %q, %q),", serverID, dataServer.GroupID, dataServer.IP, dataServer.Port, nowStr, nowStr)
		insertGroupServer := fmt.Sprintf("(%q, %q),", dataServer.GroupID, serverID)

		insertDataServerSql += insertDataServer
		insertGroupServerSql += insertGroupServer

		ds := &models.DataServer{
			ID:         serverID,
			GroupID:    dataServer.GroupID,
			IP:         dataServer.IP,
			Port:       dataServer.Port,
			CreateTime: now,
			UpdateTime: now,
		}

		log.Println(ds)
		// Save dataserver info to K/V Database as cache
		if err := db.KVDB.Create(ds); err != nil {
			log.Println(err.Error())
			return http.StatusInternalServerError, []byte(err.Error())
		}

		dsObj := make(map[string]interface{})
		dsObj["ip"] = dataServer.IP
		dsObj["port"] = dataServer.Port
		dsObj["group_id"] = dataServer.GroupID
		dsObj["data_server_id"] = serverID

		resultAry = append(resultAry, dsObj)
	}

	dbInstance := db.SQLDB.GetDB().(*gorm.DB)
	// Remove the last ','
	insertDataServerSql = insertDataServerSql[:len(insertDataServerSql)-1]
	insertGroupServerSql = insertGroupServerSql[:len(insertGroupServerSql)-1]

	if result := dbInstance.Exec(insertDataServerSql); result.Error != nil {
		log.Println(result.Error)
		return http.StatusInternalServerError, []byte(result.Error.Error())
	}

	if result := dbInstance.Exec(insertGroupServerSql); result.Error != nil {
		log.Println(result.Error)
		return http.StatusInternalServerError, []byte(result.Error.Error())
	}

	ctx.Resp.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(resultAry)
	return http.StatusOK, result
}

func DeleteDataserverHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	dataserverID := ctx.Params(":dataserver")
	ds := &models.DataServer{
		ID: dataserverID,
	}

	if err := db.KVDB.Delete(ds); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	if err := db.SQLDB.Delete(ds); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	whereClause := fmt.Sprintf("server_id=%q", dataserverID)

	if result := db.SQLDB.GetDB().(*gorm.DB).Where(whereClause).Delete(models.GroupServer{}); result.Error != nil {
		return http.StatusInternalServerError, []byte(result.Error.Error())
	}

	return http.StatusOK, nil
}

func GetDataserverHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	dataserverID := ctx.Params(":dataserver")
	ds := &models.DataServer{
		ID: dataserverID,
	}

	exist, err := db.KVDB.Query(ds)
	if exist && err == nil { // Got the server from cache
		result, _ := json.Marshal(ds)
		ctx.Resp.Header().Set("Content-Type", "application/json")
		return http.StatusOK, result
	}

	// If there is no info in cache, try to fetch it from SQLDB and rebuild the cache
	if exist, err := db.SQLDB.Query(ds); !exist && err == nil {
		return http.StatusNotFound, []byte("Data server is NOT registered")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	db.KVDB.Create(ds)

	ctx.Resp.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(ds)
	return http.StatusOK, result
}

func GetGroupsHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	dbInstance := db.SQLDB.GetDB().(*gorm.DB)
	rows, err := dbInstance.Raw("SELECT * FROM data_server, group_server WHERE data_server.id=group_server.server_id").Rows()
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	defer rows.Close()

	groupMap := make(map[string]interface{})

	for rows.Next() {
		var gsInfo models.GroupServerInfo
		dbInstance.ScanRows(rows, &gsInfo)

		server := make(map[string]interface{})
		server["data_server_id"] = gsInfo.ServerID
		server["ip"] = gsInfo.IP
		server["port"] = gsInfo.Port
		server["status"] = gsInfo.Status
		server["group_status"] = gsInfo.GroupStatus
		server["total_chunks"] = gsInfo.TotalChunks
		server["total_free_space"] = gsInfo.TotalFreeSpace
		server["max_free_space"] = gsInfo.MaxFreeSpace
		server["pending_writes"] = gsInfo.PendingWrites
		server["data_path"] = gsInfo.DataPath
		server["reading_count"] = gsInfo.ReadingCount
		server["conn_counts"] = gsInfo.ConnCounts
		server["create_time"] = gsInfo.CreateTime
		server["update_time"] = gsInfo.UpdateTime
		server["group_id"] = gsInfo.GroupID

		if _, exists := groupMap[gsInfo.GroupID]; exists == true {
			group := groupMap[gsInfo.GroupID].(map[string]interface{})
			servers := group["servers"].([]interface{})
			servers = append(servers, server)
			group["servers"] = servers
		} else {
			g := make(map[string]interface{})
			g["id"] = gsInfo.GroupID
			g["group_status"] = gsInfo.GroupStatus
			g["servers"] = []interface{}{server}
			groupMap[gsInfo.GroupID] = g
		}
	}

	groups := []interface{}{}

	for _, group := range groupMap {
		groups = append(groups, group)
	}

	ctx.Resp.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(groups)

	return http.StatusOK, result
}

func GetGroupHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	groupID := ctx.Params(":group")
	dbInstance := db.SQLDB.GetDB().(*gorm.DB)
	sqlFormat := "SELECT  * from group_server, data_server WHERE data_server.group_id=%q AND group_server.group_id=%q AND group_server.server_id = data_server.id"
	sql := fmt.Sprintf(sqlFormat, groupID, groupID)

	rows, err := dbInstance.Raw(sql).Rows()
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	defer rows.Close()

	groupMap := make(map[string]interface{})
	groupMap["id"] = groupID
	groupMap["servers"] = []interface{}{}

	for rows.Next() {
		var gsInfo models.GroupServerInfo
		dbInstance.ScanRows(rows, &gsInfo)

		server := make(map[string]interface{})
		server["data_server_id"] = gsInfo.ServerID
		server["ip"] = gsInfo.IP
		server["port"] = gsInfo.Port
		server["status"] = gsInfo.Status
		server["group_status"] = gsInfo.GroupStatus
		server["total_chunks"] = gsInfo.TotalChunks
		server["total_free_space"] = gsInfo.TotalFreeSpace
		server["max_free_space"] = gsInfo.MaxFreeSpace
		server["pending_writes"] = gsInfo.PendingWrites
		server["data_path"] = gsInfo.DataPath
		server["reading_count"] = gsInfo.ReadingCount
		server["conn_counts"] = gsInfo.ConnCounts
		server["create_time"] = gsInfo.CreateTime
		server["update_time"] = gsInfo.UpdateTime
		server["group_id"] = gsInfo.GroupID

		servers := groupMap["servers"].([]interface{})
		servers = append(servers, server)
		groupMap["servers"] = servers
	}

	numServers := len(groupMap["servers"].([]interface{}))
	if numServers == 0 {
		return http.StatusNotFound, []byte("Not Found(Group not found)")
	}

	result, _ := json.Marshal(groupMap)

	ctx.Resp.Header().Set("Content-Type", "application/json")
	return http.StatusOK, result
}
