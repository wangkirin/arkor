package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/utils/db"
	"github.com/containerops/arkor/utils/db/mysql"
)

const (
	FRAGEMENT_SIZE_MB = 4
)

var owner = models.Owner{
	ID:          "arkor",
	DisplayName: "containerops-arkor",
}

func PutObjectHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	// Recive Object Content And Check
	objectLengthStr := ctx.Req.Header.Get("Content-Length")
	content, err := ctx.Req.Body().Bytes()
	if err != nil {
		log.Errorf("[Upload Object Error] Can't NOT Get Object Content")
		return http.StatusBadRequest, []byte("Can't NOT Get Object Content")
	}
	objectLength, err := strconv.ParseInt(objectLengthStr, 10, 64)
	if err != nil {
		log.Errorf("[Upload Object Error] Convert Content-Length Header Error")
		return http.StatusBadRequest, []byte("Convert Content-Length Header Error")
	}
	if objectLength != len(content) {
		return http.StatusBadRequest, []byte("Incorrect Content-Length")
	}
	if objectLength == 0 {
		return http.StatusBadRequest, []byte("Upload empty object")
	}

	// Divide into fragments
	fragSize := FRAGEMENT_SIZE_MB * 1024 * 1024
	fragCount := int(objectLength / int64(fragSize))
	partial := int(objectLength % int64(fragSize))

	// The object contains only one fragment
	if fragCount == 0 && partial != 0 {

	}

	// Return the upload result
	return http.StatusOK, nil
}
