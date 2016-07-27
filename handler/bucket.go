package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/utils/db"
	"github.com/containerops/arkor/utils/db/mysql"
)

var owner = models.Owner{
	ID:          "arkor",
	DisplayName: "containerops-arkor",
}

func PutBucketHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	bucketName := ctx.Params(":bucket")

	// Check if Bucket exist
	bucketQuery := models.Bucket{
		Name: bucketName,
	}
	if exist, err := db.SQLDB.Query(&bucketQuery); exist && err == nil {
		return http.StatusConflict, []byte("Bucket exist")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// If not exist, insert it
	bucket := models.Bucket{
		Name:         bucketName,
		CreationDate: time.Now(),
		Owner:        owner,
	}
	if err := db.SQLDB.Create(&bucket); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	ctx.Resp.Header().Set("Date", time.Now().Format("2006-01-02 15:04:05 -0700"))
	ctx.Resp.Header().Set("Content-Type", "type")
	ctx.Resp.Header().Set("Content-Length", "0")
	return http.StatusOK, nil
}

func HeadBucketHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	bucketName := ctx.Params(":bucket")

	bucket := models.Bucket{
		Name: bucketName,
	}

	// Check if Bucket exist
	if exist, err := db.SQLDB.Query(&bucket); !exist && err == nil {
		return http.StatusNotFound, []byte("Bucket NOT exist")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	ctx.Resp.Header().Set("Date", time.Now().Format("2006-01-02 15:04:05 -0700"))
	ctx.Resp.Header().Set("Content-Type", "type")
	ctx.Resp.Header().Set("Content-Length", "0")
	return http.StatusOK, nil
}

func DeleteBucketHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	bucketName := ctx.Params(":bucket")

	// Check if Bucket exist
	bucket := models.Bucket{
		Name: bucketName,
	}
	if exist, err := db.SQLDB.Query(&bucket); !exist && err == nil {
		return http.StatusNotFound, []byte("Bucket NOT exist")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// If exist, Delete it
	if err := db.SQLDB.Delete(&bucket); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	ctx.Resp.Header().Set("Date", time.Now().Format("2006-01-02 15:04:05 -0700"))
	return http.StatusOK, nil
}

func GetBucketHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	bucketName := ctx.Params(":bucket")

	// Check if Bucket exist
	bucket := models.Bucket{
		Name: bucketName,
	}
	if exist, err := db.SQLDB.Query(&bucket); !exist && err == nil {
		return http.StatusNotFound, []byte("Bucket NOT exist")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// Preload Associations and Query
	mysqldb := mysql.MySQLInstance()
	mysqldb.Preload("Contents").Find(&bucket)

	// Output
	result, err := json.Marshal(bucket)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	ctx.Resp.Header().Set("Date", time.Now().Format("2006-01-02 15:04:05 -0700"))
	ctx.Resp.Header().Set("Content-Type", "application/json")
	ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(result)))
	return http.StatusOK, result
}
