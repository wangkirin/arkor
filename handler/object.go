package handler

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/modules"
	"github.com/containerops/arkor/utils"
)

const (
	FRAGEMENT_SIZE_MB = 4
)

func PutObjectHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	// TODO Handle bucket infomation
	// Recive Object Parameters
	objectName := ctx.Params(":object")
	objectMetadata := models.ObjectMeta{
		ID:     objectName,
		Key:    objectName,
		Md5Key: utils.MD5(objectName),
	}
	objectLengthStr := ctx.Req.Header.Get("Content-Length")
	objectLength, err := strconv.ParseInt(objectLengthStr, 10, 64)
	if err != nil {
		log.Errorf("[Upload Object Error] Convert Content-Length Header Error")
		return http.StatusBadRequest, []byte("Convert Content-Length Header Error")
	}
	if objectLength == 0 {
		return http.StatusBadRequest, []byte("Object Content is empty")
	}
	objectdata, err := ctx.Req.Body().Bytes()
	if err != nil {
		return http.StatusBadRequest, []byte("Recieve Object Content error")
	}
	// Divide into fragments

	fragSize := FRAGEMENT_SIZE_MB * 1024 * 1024
	fragCount := int64(objectLength / int64(fragSize))
	partial := int64(objectLength % int64(fragSize))

	// The object divided into one fragment
	if fragCount == 0 && partial != 0 {
		// Set up fragmentInfo Metaesrver
		fragmentInfo := models.Fragment{
			Index:  0,
			Start:  0,
			End:    partial,
			IsLast: true,
		}
		if err := modules.Upload(objectdata, &fragmentInfo); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		fragmentInfo.ModTime = time.Now()
		objectMetadata.Fragments = append(objectMetadata.Fragments, fragmentInfo)
	}

	// The object divided into more than one fragment and have partial left
	if fragCount != 0 && partial != 0 {
		// Read and Upload all fragments
		for k := 0; k < int(fragCount); k++ {
			fragmentInfo := models.Fragment{
				Index:  k,
				Start:  int64(k * fragSize),
				End:    int64((k + 1) * fragSize),
				IsLast: false,
			}
			fragdata := objectdata[int64(k*fragSize):int64((k+1)*fragSize)]
			if err := modules.Upload(fragdata, &fragmentInfo); err != nil {
				return http.StatusInternalServerError, []byte(err.Error())
			}
			objectMetadata.Fragments = append(objectMetadata.Fragments, fragmentInfo)
		}
		// Read and Upload partial data
		fragmentInfo := models.Fragment{
			Index:  int(fragCount),
			Start:  (fragCount) * int64(fragSize),
			End:    (fragCount)*int64(fragSize) + partial,
			IsLast: true,
		}
		fragdata := objectdata[(fragCount+1)*int64(fragSize) : (fragCount+1)*int64(fragSize)+partial]
		if err := modules.Upload(fragdata, &fragmentInfo); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		objectMetadata.Fragments = append(objectMetadata.Fragments, fragmentInfo)
	}

	// The object divided into more than one fragment and have no partial left
	if fragCount != 0 && partial == 0 {
		// Read and Upload all fragments
		for k := 0; k < int(fragCount); k++ {
			fragmentInfo := models.Fragment{
				Index: k,
				Start: int64(k * fragSize),
				End:   int64((k + 1) * fragSize),
			}
			if k != int(fragCount-1) {
				fragmentInfo.IsLast = false
			} else {
				fragmentInfo.IsLast = true
			}
			fragdata := objectdata[int64(k*fragSize):int64((k+1)*fragSize)]
			if err := modules.Upload(fragdata, &fragmentInfo); err != nil {
				return http.StatusInternalServerError, []byte(err.Error())
			}
			objectMetadata.Fragments = append(objectMetadata.Fragments, fragmentInfo)
		}
	}
	if err := modules.SaveObjectInfo(objectMetadata); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	return http.StatusOK, nil
}

func GetObjectHandler(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {

	outputBuf := bytes.NewBuffer([]byte{})
	// Handle Input Para
	objectName := ctx.Params(":object")
	objectMetadata, err := modules.GetObjectInfo(objectName)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	// Query All Fragment Info of the object
	fragmentsInfo := make([]models.Fragment, len(objectMetadata.Fragments))
	for _, fragment := range objectMetadata.Fragments {
		fragmentsInfo[fragment.Index] = fragment
	}

	// Download All Fragments
	for i := 0; i < len(fragmentsInfo); i++ {
		data, err := modules.Download(&fragmentsInfo[i])
		if err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		if _, err := outputBuf.Write(data); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
	}
	return http.StatusOK, outputBuf.Bytes()
}
