package web

import (
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/middleware"
	"github.com/containerops/arkor/router"
	"github.com/containerops/arkor/utils/setting"
)

func SetArkorMacaron(m *macaron.Macaron) {
	//Setting Middleware
	middleware.SetMiddlewares(m)
	//Setting Router
	router.SetRouters(m)
	//static
	if setting.RunTime.Run.RunMode == "dev" {
		m.Use(macaron.Static("external"))
	}

}
