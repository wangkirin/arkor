package modules

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/containerops/arkor/models"
	"github.com/containerops/arkor/modules/pools"
)

func Download(fragmentInfo *models.Fragment) ([]byte, error) {
	// Get All servers of the Groups this fragment belongs to
	groupID := fragmentInfo.GroupID
	servers, err := GetServersByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	// Chose one available server
	index := rand.Int() % len(servers)
	server := servers[index]
	if server.Status != models.RW_STATUS {
		log.Infof("Server %s is NOT available, Retry", server.ID)
		for i := 0; i < len(servers); i++ {
			server = servers[i]
			if server.Status == models.RW_STATUS {
				log.Infof("Get an available dataserver: %s", server.ID)
				break
			}
			if i == (len(servers) - 1) {
				return nil, fmt.Errorf("Can not find available dataserver to download")
			}
		}
	}
	connPools := pools.ConnectionPools
	conn, err := connPools.GetConn(&server)
	if err != nil {
		log.Errorf("Download Object Get Connection error: %v", err)
		return nil, err
	}
	data, err := DownloadData(fragmentInfo, conn.(*pools.PooledConn))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DownloadData(fragmentInfo *models.Fragment, conn *pools.PooledConn) ([]byte, error) {
	output := new(bytes.Buffer)
	header := make([]byte, HEADERSIZE)

	// Temp
	groupIDnum, _ := strconv.Atoi(fragmentInfo.GroupID)
	fileIDint, _ := FragIDStr2Int(fragmentInfo.ID)
	//
	binary.Write(output, binary.BigEndian, GET)
	binary.Write(output, binary.BigEndian, uint32(2+8))
	binary.Write(output, binary.BigEndian, uint16(groupIDnum))
	binary.Write(output, binary.BigEndian, uint64(fileIDint))

	_, err := conn.Write(output.Bytes())
	if err != nil {
		fmt.Errorf("write socket error %s\n", err)
		return nil, err
	}

	_, err = io.ReadFull(conn.Br, header)
	if err != nil {
		log.Errorf("GetData read header error: %s", err)
		return nil, err
	}

	if header[0] != GET || header[1] != 0 {
		return nil, fmt.Errorf("download file failed, code = %d\n", header[1])
	}

	bodyLen := binary.BigEndian.Uint32(header[2:])
	data := make([]byte, bodyLen)
	log.Debugf("GetData len: %d, %d", bodyLen, len(data))

	if _, err := io.ReadFull(conn.Br, data); err != nil {
		return nil, fmt.Errorf("read socket error %s", err)
	}

	return data, nil
}
