package handler

import (
	// "encoding/json"
	// "fmt"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/modules"
	"github.com/containerops/arkor/modules/pools"
	// "github.com/containerops/arkor/utils/db"
	// "github.com/containerops/arkor/utils/db/mysql"
)

const (
	FRAGEMENT_SIZE_MB = 4
)

func PutObjectHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	// Recive ObjectLength
	objectLengthStr := ctx.Req.Header.Get("Content-Length")
	objectLength, err := strconv.ParseInt(objectLengthStr, 10, 64)
	if err != nil {
		log.Errorf("[Upload Object Error] Convert Content-Length Header Error")
		return http.StatusBadRequest, []byte("Convert Content-Length Header Error")
	}
	if objectLength == 0 {
		return http.StatusBadRequest, []byte("Upload empty object")
	}

	// Divide into fragments
	fragSize := FRAGEMENT_SIZE_MB * 1024 * 1024
	fragCount := int64(objectLength / int64(fragSize))
	partial := int64(objectLength % int64(fragSize))

	// The object contains only one fragment
	if fragCount == 0 && partial != 0 {
		content, err := ctx.Req.Body().Bytes()
		if err != nil {
			log.Errorf("[Upload Object Error] Can't NOT Get Object Content")
			return http.StatusBadRequest, []byte("Can't NOT Get Object Content")
		}
		datagroups, err := modules.GetDataGroups()
		if err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		datagroup, err := modules.SelectDataGroup(datagroups, partial)
		if err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		if err := pools.SyncDataServerConnectionPools(datagroups); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		// Generate fileID
		h := md5.New()
		h.Write(content)
		fileID := hex.EncodeToString(h.Sum(nil))
		// Set up fragmentInfo Metaesrver
		fragmentInfo := models.Fragment{
			Index:   1,
			Start:   0,
			End:     partial,
			FileID:  fileID,
			GroupID: datagroup.ID,
			IsLast:  true,
		}
		if err := modules.Upload(content, datagroup, &fragmentInfo); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		fragmentInfo.ModTime = time.Now()
	}

	// Return the upload result
	return http.StatusOK, nil
}
