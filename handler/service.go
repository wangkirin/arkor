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
)

func GetServiceHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	// Get All Buckets
	querycondition := &models.Bucket{}
	buckets := []models.Bucket{}
	arkorOwner := models.Owner{
		ID:          "arkor",
		DisplayName: "containerops-arkor",
	}
	if exist, err := db.SQLDB.QueryMulti(querycondition, &buckets); !exist && err == nil {
		return http.StatusNotFound, []byte("NOT Found Any Buckets")
	} else if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// Convert Bucket list to ouptut format
	responseBuckets := []models.BucketSimple{}
	for _, bucket := range buckets {
		bucketSimpleTmp := models.BucketSimple{
			Name:         bucket.Name,
			CreationDate: bucket.CreationDate,
		}
		responseBuckets = append(responseBuckets, bucketSimpleTmp)
	}

	//  Set up response json and output
	bucketlist := models.BucketListResponse{
		Type:    "ListAllMyBucketsResult",
		Owner:   arkorOwner,
		Buckets: responseBuckets,
	}

	result, err := json.Marshal(bucketlist)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	ctx.Resp.Header().Set("Date", time.Now().Format("2006-01-02 15:04:05 -0700"))
	ctx.Resp.Header().Set("Content-Type", "application/json")
	ctx.Resp.Header().Set("Content-Length", fmt.Sprint(len(result)))
	return http.StatusOK, result
}
