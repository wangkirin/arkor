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
		return http.StatusBadRequest, []byte("Upload empty object")
	}

	// Divide into fragments
	fragSize := FRAGEMENT_SIZE_MB * 1024 * 1024
	fragCount := int64(objectLength / int64(fragSize))
	partial := int64(objectLength % int64(fragSize))

	// The object divided into one fragment
	if fragCount == 0 && partial != 0 {
		log.Infoln("enter partial")
		content, err := ctx.Req.Body().Bytes()
		log.Infoln(content)
		if err != nil {
			log.Errorf("[Upload Object Error] Can't NOT Get Object Content")
			return http.StatusBadRequest, []byte("Can't NOT Get Object Content")
		}
		// Set up fragmentInfo Metaesrver
		fragmentInfo := models.Fragment{
			Index:  0,
			Start:  0,
			End:    partial,
			IsLast: true,
		}
		if err := modules.Upload(content, &fragmentInfo); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		fragmentInfo.ModTime = time.Now()
		log.Infoln("fragmentInfo===")
		log.Infoln(fragmentInfo)
		objectMetadata.Fragments = append(objectMetadata.Fragments, fragmentInfo)
	}

	// The object divided into more than one fragment and have partial left
	if fragCount != 0 && partial != 0 {
		log.Infoln("enter multi-part")
		// Read and Upload all fragments
		for k := 0; k < int(fragCount); k++ {
			log.Infof("k=%d", k)
			fragmentInfo := models.Fragment{
				Index:  k,
				Start:  int64(k * fragSize),
				End:    int64((k + 1) * fragSize),
				IsLast: false,
			}
			buf := make([]byte, fragSize)
			n, err := ctx.Req.Request.Body.Read(buf)
			if err != nil {
				return http.StatusInternalServerError, []byte(err.Error())
			}
			log.Infof("read buffer number:= %d", n)
			if n != fragSize {
				return http.StatusInternalServerError, []byte("Read Body Error")
			}
			if err := modules.Upload(buf, &fragmentInfo); err != nil {
				return http.StatusInternalServerError, []byte(err.Error())
			}
			objectMetadata.Fragments = append(objectMetadata.Fragments, fragmentInfo)
		}
		// Read and Upload partial data
		fragmentInfo := models.Fragment{
			Index:  int(fragCount + 1),
			Start:  (fragCount + 1) * int64(fragSize),
			End:    (fragCount+1)*int64(fragSize) + partial,
			IsLast: true,
		}
		buf := make([]byte, partial)
		n, err := ctx.Req.Request.Body.Read(buf)
		if err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
		if n != fragSize {
			return http.StatusInternalServerError, []byte("Read Body Error")
		}
		if err := modules.Upload(buf, &fragmentInfo); err != nil {
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

			buf := make([]byte, fragSize)
			n, err := ctx.Req.Request.Body.Read(buf)
			if err != nil {
				return http.StatusInternalServerError, []byte(err.Error())
			}
			if n != fragSize {
				return http.StatusInternalServerError, []byte("Read Body Error")
			}
			if err := modules.Upload(buf, &fragmentInfo); err != nil {
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
	log.Infoln("enter get obj handler------")
	outputBuf := bytes.NewBuffer([]byte{})
	// Handle Input Para
	objectName := ctx.Params(":object")
	objectMetadata, err := modules.GetObjectInfo(objectName)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	log.Infoln("objectmeta======")
	log.Println(objectMetadata)
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
		log.Println("outputdata*******")
		log.Println(string(data))
		if _, err := outputBuf.Write(data); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
	}
	return http.StatusOK, outputBuf.Bytes()
}
