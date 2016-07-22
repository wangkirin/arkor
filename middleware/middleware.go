package middleware

import (
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/setting"
)

func SetMiddlewares(m *macaron.Macaron) {
	// init log settings
	InitLog()
	//Set global Logger
	m.Map(Log)
	//Set logger handler function, deal with all the Request log output
	m.Use(logger(setting.RunTime.Run.RunMode))
	//Set recovery handler to returns a middleware that recovers from any panics
	m.Use(macaron.Recovery())
}
