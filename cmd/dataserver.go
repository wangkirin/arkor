package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ghodss/yaml"
)

var CmdDataServer = cli.Command{
	Name:        "dataserver",
	Usage:       "run arkor dataserver in local machine",
	Description: "arkor is the object storage service of containerops",
	Action:      runLocalDataServer,
}

var DataServerConf *DataServer

const DATASERVER_BINARY_PATH = "./dataserver/spy_server"

type DataServer struct {
	GroupID                string "yaml:groupID"
	IP                     string "yaml:IP"
	Port                   string "yaml:port"
	RegistrationCenterIP   string "yaml:registrationCenterIP"
	RegistrationCenterPort string "yaml:registrationCenterPort"
	DataDir                string "yaml:dataDir"
	ErrorLogDir            string "yaml:errorLogDir"
	ChunkNum               string "yaml:chunkNum"
}

func initDataServerConf(Path string) error {
	conf, err := ioutil.ReadFile(Path)

	if err != nil {
		return err
	}
	DataServerConf = &DataServer{}

	err = yaml.Unmarshal([]byte(conf), &DataServerConf)
	if err != nil {
		return err
	}
	return nil
}

func runLocalDataServer(c *cli.Context) {
	// Load Configs
	if err := initDataServerConf("conf/dataserver.yaml"); err != nil {
		log.Errorf("Read Data Server config error: %v", err.Error())
		return
	}

	// Check if chunkserver binary exsist
	_, err := os.Stat(DATASERVER_BINARY_PATH)
	if err != nil && os.IsNotExist(err) {
		log.Fatalln("Cannot find binary file of  DataServer")
	}

	// Check if errlog folder exsist , if not ,create it
	_, err = os.Stat(DataServerConf.ErrorLogDir)
	if err != nil || os.IsNotExist(err) {
		os.MkdirAll(DataServerConf.ErrorLogDir, 0777)
	}
	// Check if data folder exsist , if not ,create it
	_, err = os.Stat(DataServerConf.DataDir)
	if err != nil || os.IsNotExist(err) {
		os.MkdirAll(DataServerConf.DataDir, 0777)
	}
	// Start Local DataServer
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(DATASERVER_BINARY_PATH, "--ip", DataServerConf.IP, "--port", DataServerConf.Port, "--master_ip", DataServerConf.RegistrationCenterIP,
		"--master_port", DataServerConf.RegistrationCenterPort, "--chunks", DataServerConf.ChunkNum,
		"--group_id", DataServerConf.GroupID, "--data_dir", DataServerConf.DataDir, "--error_log", DataServerConf.ErrorLogDir)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Infof("Start local DataServer error, STDOUT: %s", stdout.Bytes())
		log.Infof("Start local DataServer error, STDERR: %s", stderr.Bytes())
		log.Fatalf("Start local DataServer error, error INFO: %v", err)
	}
}
