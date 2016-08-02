package modules

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/modules/pools"
)

const (
	HEADERSIZE = 6
	SUCCESS    = "Success"
)

var (
	PUT    uint8 = 0x00
	GET    uint8 = 0x01
	DELETE uint8 = 0x02
	PING   uint8 = 0x0A
)

func Upload(data []byte, datagroup *models.Group, fragInfo *models.Fragment) error {
	// Count the number of normal servers
	normalCount := 0
	for _, server := range datagroup.Servers {
		if server.Status == models.RW_STATUS {
			normalCount++
		}
	}
	// Use multi-goroutine to Upload
	ch := make(chan string, normalCount)
	for _, server := range datagroup.Servers {
		if server.Status == models.RW_STATUS {
			go concurrenceUpload(server, data, ch, fragInfo.FileID)
		}
	}

	err := handleUploadResult(ch, normalCount)
	if err != nil {
		log.Errorf("upload error: %s", err)
		return err
	}
	return nil
}

func concurrenceUpload(server models.DataServer, data []byte, c chan string, fileId string) {
	connPools := pools.ConnectionPools
	if connPools == nil {
		log.Errorf("connectionPools is nil")
		c <- "connectionPools is nil"
		return
	}

	conn, err := connPools.GetConn(&server)
	if err != nil {
		log.Errorf("Can not get connection: %s", err.Error())
		c <- err.Error()
		return
	}

	log.Debugf("Begin to upload data")
	err = PutData(data, conn.(*pools.PooledConn), fileId, server.GroupID)
	if err != nil {
		log.Errorf("upload data failed: %s", err)
		conn.Close()
		pools.ConnectionPools.ReleaseConn(conn)
		c <- err.Error()
		checkErrorAndConnPool(err, &server, connPools)
		return
	}
	log.Debugf("Upload data success")
	c <- SUCCESS
	connPools.ReleaseConn(conn)
}

func PutData(data []byte, conn *pools.PooledConn, fileId string, groupID string) error {
	output := new(bytes.Buffer)
	header := make([]byte, HEADERSIZE)

	binary.Write(output, binary.BigEndian, PUT)
	binary.Write(output, binary.BigEndian, uint32(len(data)+2+8))
	binary.Write(output, binary.BigEndian, groupID)
	binary.Write(output, binary.BigEndian, fileId)

	output.Write(data)
	_, err := conn.Write(output.Bytes())
	if err != nil {
		log.Errorf("write conn error: %s", err)
		return err
	}

	if _, err := io.ReadFull(conn.Br, header); err != nil {
		log.Errorf("read header error: %s", err)
		return err
	}

	if header[0] == PUT && header[1] == 0 {
		log.Debugf("upload success")
		return nil
	}

	log.Errorf("fileId: %d, upload failed, header[0] = %d, header[1] = %d", fileId, header[0], header[1])
	return fmt.Errorf("upload error, code: %d", header[1])
}

func handleUploadResult(ch chan string, size int) error {
	var result, tempResult string
	var failed = false

	if ch == nil {
		log.Errorf("ch is nil")
		return fmt.Errorf("[handleUploadResult] channel is nil")
	}

	for k := 0; k < size; k++ {
		tempResult = <-ch
		if len(tempResult) != 0 {
			result = tempResult
			failed = true
		}
	}
	if failed {
		log.Errorf("Upload Object failed: %s", result)
		return fmt.Errorf(result)
	}

	return nil
}

func checkErrorAndConnPool(err error, server *models.DataServer, connPools *pools.DataServerConnectionPools) {
	if "EOF" == err.Error() {
		err := connPools.CheckConnPool(server)
		if err != nil {
			log.Errorf("CheckConnPool error: %v", err)
		}
	}
}
