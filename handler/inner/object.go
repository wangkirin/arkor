package inner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/utils"
	"github.com/containerops/arkor/utils/db"
)

func AllocateFileID(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	m := make(map[string]string)
	m["file_id"] = utils.MD5ID()

	result, err := json.Marshal(m)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, result
}

func PutObjectInfoHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	data, _ := ctx.Req.Body().Bytes()
	reqBody := make(map[string]interface{})
	json.Unmarshal(data, &reqBody)
	fragments := reqBody["fragments"].([]interface{})

	if len(fragments) == 0 {
		return http.StatusBadRequest, []byte("Invalid Parameters")
	}

	dbInstance := db.SQLDB.GetDB().(*gorm.DB)

	object_id := reqBody["object_id"].(string)
	object_key := reqBody["object_key"].(string)
	md5_key := reqBody["md5_key"].(string)

	// Remove the old relations
	if result := dbInstance.Exec("DELETE FROM object WHERE object_id='" + object_id + "'"); result.Error != nil {
		log.Errorln(result.Error.Error())
		return http.StatusInternalServerError, []byte("Internal Server Error")
	}

	// Save ObjectMeta
	insertOrUpdateSQL := fmt.Sprintf("REPLACE INTO object_meta(id, object_key, md5_key) VALUES(%q, %q, %q)", object_id, object_key, md5_key)
	if result := dbInstance.Exec(insertOrUpdateSQL); result.Error != nil {
		log.Errorln(result.Error.Error())
		return http.StatusInternalServerError, []byte("Internal Server Error")
	}

	// Save fragments and the relation
	insertFragmentsSQL := "INSERT INTO fragment(id, `index`, start, end, group_id, file_id, is_last, mod_time) VALUES "
	insertRelationsSQL := "REPLACE INTO object(object_id, fragment_id) VALUES "

	for index := range fragments {
		f := fragments[index].(map[string]interface{})

		fragmentID := utils.MD5ID()
		index := int(f["index"].(float64))
		start := int64(f["start"].(float64))
		end := int64(f["end"].(float64))
		group_id := f["group_id"].(string)
		file_id := f["file_id"].(string)
		is_last := f["is_last"].(bool)

		modTimeStr := f["mod_time"].(string)
		t, err := time.Parse("2006-01-02T15:04:05Z", modTimeStr)
		if err != nil {
			log.Errorln(err.Error())
			return http.StatusBadRequest, []byte("Invalid Parameters")
		}
		mod_time := t.Format("2006-01-02 15:04:05")

		insertFragmentSQL := fmt.Sprintf("(%q,%d,%d,%d,%q,%q,%t,%q),", fragmentID, index, start, end, group_id, file_id, is_last, mod_time)
		insertFragmentsSQL += insertFragmentSQL

		insertRelationSQL := fmt.Sprintf("(%q,%q),", reqBody["object_id"], fragmentID)
		insertRelationsSQL += insertRelationSQL
	}

	insertFragmentsSQL = insertFragmentsSQL[:len(insertFragmentsSQL)-1]
	if result := db.SQLDB.GetDB().(*gorm.DB).Exec(insertFragmentsSQL); result.Error != nil {
		log.Errorln(result.Error.Error())
		return http.StatusInternalServerError, []byte("Internal Server Error")
	}

	insertRelationsSQL = insertRelationsSQL[:len(insertRelationsSQL)-1]
	if result := db.SQLDB.GetDB().(*gorm.DB).Exec(insertRelationsSQL); result.Error != nil {
		log.Errorln(result.Error.Error())
		return http.StatusInternalServerError, []byte("Internal Server Error")
	}

	return http.StatusOK, nil
}

func GetObjectInfoHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	objectID := ctx.Params(":object")

	sql := fmt.Sprintf(
		"SELECT object.object_id, object.fragment_id, object_key, md5_key, `index`, start, end, group_id, file_id, is_last, mod_time FROM object, object_meta, fragment WHERE object_meta.id=%q AND object.object_id=%q  AND fragment.id=object.fragment_id",
		objectID, objectID)

	dbInstace := db.SQLDB.GetDB().(*gorm.DB)

	rows, err := dbInstace.Raw(sql).Rows()
	if err != nil {
		return http.StatusInternalServerError, []byte("Internal Server Error")
	}

	defer rows.Close()

	object := make(map[string]interface{})
	fragments := []interface{}{}
	for rows.Next() {
		var object_id string
		var fragment_id string
		var object_key string
		var md5_key string
		var index int
		var start int64
		var end int64
		var group_id string
		var file_id string
		var is_last bool
		var mod_time string

		rows.Scan(&object_id, &fragment_id, &object_key, &md5_key, &index, &start, &end, &group_id, &file_id, &is_last, &mod_time)

		object["object_id"] = object_id
		object["md5_key"] = md5_key
		object["object_key"] = object_key

		fragment := make(map[string]interface{})
		fragment["index"] = index
		fragment["start"] = start
		fragment["end"] = end
		fragment["group_id"] = group_id
		fragment["file_id"] = file_id
		fragment["is_last"] = is_last
		fragment["mod_time"] = mod_time

		fragments = append(fragments, fragment)
	}
	object["fragments"] = fragments

	result, err := json.Marshal(object)

	if err != nil {
		log.Errorln(err.Error())
		return http.StatusInternalServerError, []byte("Internal Server Error")

	}
	ctx.Resp.Header().Set("Content-Type", "application/json")
	return http.StatusOK, result
}
