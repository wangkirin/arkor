package pools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/containerops/arkor/models"
)

var (
	PING uint8 = 0x0A
)

var ConnectionPools *DataServerConnectionPools

type DataServerConnectionPools struct {
	mu    sync.Mutex
	Pools map[string]*ConnectionPool // <ip:port>:connectionpool
}

func NewDataServerConnectionPools() *DataServerConnectionPools {
	return &DataServerConnectionPools{
		mu:    sync.Mutex{},
		Pools: make(map[string]*ConnectionPool),
	}
}

func (d *DataServerConnectionPools) GetConn(dataserver *models.DataServer) (PoolConnection, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := fmt.Sprintf("%s:%d", dataserver.IP, dataserver.Port)
	pool, ok := d.Pools[key]
	if !ok {
		return nil, fmt.Errorf("pool %s not exist", key)
	}

	return pool.Get()
}

//chunkserver closed, the state of connection in pool is close_wait, need to close those connection
func (d *DataServerConnectionPools) CheckConnPool(dataserver *models.DataServer) error {
	for {
		conn, err := d.GetConn(dataserver)
		if err != nil {
			return err
		}
		err = Ping(dataserver, conn.(*PooledConn))
		if err != nil {
			conn.Close()
			d.ReleaseConn(conn)
			continue
		}

		return nil
	}
}

func (d *DataServerConnectionPools) ReleaseConn(pc PoolConnection) {
	pc.Recycle()
}

func (d *DataServerConnectionPools) AddPool(dataserver *models.DataServer, capacity int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := fmt.Sprintf("%s:%d", dataserver.IP, dataserver.Port)
	log.Debugf("Add Connection Pool of DataServer: %s", key)

	pool, ok := d.Pools[key]
	if ok {
		log.Debugf("Connection Pool of DataServer: %s already exist", key)
		return nil
	}

	pool = NewConnectionPool("DataServer Connection Pool", capacity, 3600*time.Second)

	log.Debugf("Tring to Open DataServer Connection Pool")
	pool.Open(ConnectionCreator(key))
	log.Debugf("Open DataServer Connection Pool success")

	d.Pools[key] = pool
	return nil
}

func (d *DataServerConnectionPools) AddExistPool(key string, pool *ConnectionPool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Debugf("AddExistPool, key: %v, pool: %v", key, pool)

	_, ok := d.Pools[key]
	if ok {
		log.Infof("AddExistPool key: %s already exist", key)
		return
	}

	d.Pools[key] = pool
	log.Debugf("AddExistPool, key: %v, pool: %v", key, pool)
	return
}

func (d *DataServerConnectionPools) RemovePool(dataserver *models.DataServer) {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := fmt.Sprintf("%s:%d", dataserver.IP, dataserver.Port)
	log.Debugf("RemovePool, key: %v", key)

	delete(d.Pools, key)
}

func (d *DataServerConnectionPools) RemoveAndClosePool(dataserver *models.DataServer) error {
	d.mu.Lock()

	key := fmt.Sprintf("%s:%d", dataserver.IP, dataserver.Port)
	pool, ok := d.Pools[key]
	if !ok {
		d.mu.Unlock()
		return fmt.Errorf("pool %s not exist", key)
	}

	delete(d.Pools, key)

	d.mu.Unlock()

	pool.Close()
	return nil
}

func Ping(server *models.DataServer, conn *PooledConn) error {
	output := new(bytes.Buffer)
	header := make([]byte, 6)
	binary.Write(output, binary.BigEndian, PING)
	binary.Write(output, binary.BigEndian, uint32(0))

	_, err := conn.Write(output.Bytes())
	if err != nil {
		return err
	}

	_, err = io.ReadFull(conn.Br, header)
	if err != nil {
		return err
	}

	if header[0] != PING || header[1] != 0 {
		return fmt.Errorf("ping %s:%d error, header[0]:%d, header[1]:%d", server.IP, server.Port, uint8(header[0]), uint8(header[1]))
	}

	return nil
}
