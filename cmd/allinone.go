package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/models"
	. "github.com/containerops/arkor/setting"
	// "github.com/containerops/arkor/utils/db"
	"github.com/containerops/arkor/web"
)

var CmdAllInOne = cli.Command{
	Name:        "allinone",
	Usage:       "run all arkor object server, registartion center and data servers in local machine",
	Description: "arkor is the object storage service of containerops",
	Action:      runAllInOne,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Value: "0.0.0.0",
			Usage: "web service listen ip, default is 0.0.0.0",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "8990",
			Usage: "web service listen at port 80; if run with https will be 443.",
		},
	},
}

var dataservers []models.DataServer

func runAllInOne(c *cli.Context) {
	m := macaron.New()

	//Set Macaron Web Middleware And Routers
	web.SetArkorMacaron(m)

	//Init Object Server Config
	rc := &RegistrationCenter{
		Address: "127.0.0.1",
		Port:    c.String("port"),
	}
	ObjectServerConf = &ObjectServer{
		RegistrationCenter: rc,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		listenaddr := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			fmt.Printf("start generator http service error: %v", err.Error())
		}
	}()
	runtime.Gosched()

	// Init & Register Dataservers
	InitDataServer()
	RegisterDataServer()

	// Check if chunkserver binary exsist
	_, err := os.Stat(DATASERVER_BINARY_PATH)
	if err != nil && os.IsNotExist(err) {
		log.Fatalln("Cannot find binary file of  DataServer")
	}

	//Start Data Servers
	for _, ds := range dataservers {
		go func(ds models.DataServer) {
			log.Infof("Data Server IP: %v, Port:%v", ds.IP, ds.Port)
			// Check if errlog folder exsist , if not ,create it
			datadir := fmt.Sprintf("./data/data_%v_%v", ds.IP, ds.Port)
			_, err = os.Stat(datadir)
			if err != nil || os.IsNotExist(err) {
				os.MkdirAll(datadir, 0777)
			}
			// Check if errlog folder exsist , if not ,create it
			errlogpath := fmt.Sprintf("./data/error/errlog_%v_%v.log", ds.IP, ds.Port)
			_, err = os.Stat("./data/error")
			if err != nil || os.IsNotExist(err) {
				os.MkdirAll("./data/error", 0777)
			}
			_, err = os.Stat(errlogpath)
			if err != nil || os.IsNotExist(err) {
				os.Create(errlogpath)
			}

			var stdout, stderr bytes.Buffer
			cmd := exec.Command(DATASERVER_BINARY_PATH, "--ip", ds.IP, "--port", fmt.Sprintf("%v", ds.Port), "--master_ip", ObjectServerConf.RegistrationCenter.Address,
				"--master_port", ObjectServerConf.RegistrationCenter.Port, "--chunks", "1",
				"--group_id", ds.GroupID, "--data_dir", datadir, "--error_log", errlogpath)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Infof("Start local DataServer error, STDOUT: %s", stdout.Bytes())
				log.Infof("Start local DataServer error, STDERR: %s", stderr.Bytes())
				log.Fatalf("Start local DataServer error, error INFO: %v", err)
			}
		}(ds)
		runtime.Gosched()
	}
	wg.Wait()
}

func InitDataServer() {
	ds1 := models.DataServer{
		ID:      "ALLinONEdataserver1",
		GroupID: "2",
		IP:      "127.0.0.1",
		Port:    8125,
	}
	ds2 := models.DataServer{
		ID:      "ALLinONEdataserver2",
		GroupID: "2",
		IP:      "127.0.0.1",
		Port:    8126,
	}
	ds3 := models.DataServer{
		ID:      "ALLinONEdataserver3",
		GroupID: "2",
		IP:      "127.0.0.1",
		Port:    8127,
	}
	dataservers = append(dataservers, ds1, ds2, ds3)
}

func RegisterDataServer() error {
	// Sent POST Request to Register
	dataServerJson, _ := json.Marshal(dataservers)
	registerURI := fmt.Sprintf("http://%s:%s/internal/v1/dataserver", ObjectServerConf.RegistrationCenter.Address, ObjectServerConf.RegistrationCenter.Port)
	body := bytes.NewBuffer([]byte(dataServerJson))
	resp, err := http.Post(registerURI, "application/json", body)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Register Dataserver Error, Status Code: %d, Info: %s", resp.StatusCode, result)
	}
	return nil
}
