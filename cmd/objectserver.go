package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/codegangsta/cli"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/setting"
	"github.com/containerops/arkor/web"
)

var CmdObjectServer = cli.Command{
	Name:        "objectserver",
	Usage:       "run arkor as object server",
	Description: "arkor is the object storage service of containerops",
	Action:      runObjectServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Value: "0.0.0.0",
			Usage: "web service listen ip, default is 0.0.0.0",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "8990",
			Usage: "web service listen port;default is 8990",
		},
	},
}

func runObjectServer(c *cli.Context) {
	m := macaron.New()

	//Set Macaron Web Middleware And Routers
	web.SetObjectServerMacaron(m)

	switch setting.RunTime.Http.ListenMode {
	case "http":
		listenaddr := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			fmt.Printf("start arkor Object Server http service error: %v \n", err.Error())
		}
		break
	case "https":
		listenaddr := fmt.Sprintf("%s:%s", c.String("address"), "443")
		server := &http.Server{Addr: listenaddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		if err := server.ListenAndServeTLS(setting.RunTime.Http.HttpsCertFile, setting.RunTime.Http.HttpsKeyFile); err != nil {
			fmt.Printf("start arkor Object Server https service error: %v \n", err.Error())
		}
		break
	default:
		break
	}
}
