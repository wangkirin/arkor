package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/containerops/arkor/cmd"
	"github.com/containerops/arkor/setting"
)

func main() {

	if err := setting.InitConf("conf/global.yaml", "conf/runtime.yaml"); err != nil {
		fmt.Printf("Read config error: %v", err.Error())
		return
	}

	app := cli.NewApp()

	app.Name = setting.Global.AppName
	app.Usage = setting.Global.Usage
	app.Version = setting.Global.Version
	app.Author = setting.Global.Author
	app.Email = setting.Global.Email

	app.Commands = []cli.Command{
		cmd.CmdWeb,
		cmd.CmdObjectServer,
		cmd.CmdRegistrationCenter,
		cmd.CmdDataServer,
		cmd.CmdAllInOne,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
