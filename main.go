package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/containerops/arkor/cmd"
	"github.com/containerops/arkor/setting"
)

func main() {
	app := cli.NewApp()

	app.Name = setting.Global.AppName
	app.Usage = setting.Global.Usage
	app.Version = setting.Global.Version
	app.Author = setting.Global.Author
	app.Email = setting.Global.Email

	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
