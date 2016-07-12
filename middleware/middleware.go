package middleware

import (
	"gopkg.in/macaron.v1"
)

func SetMiddlewares(m *macaron.Macaron) {
	//Set global Logger
	m.Map(Log)
	//Set logger handler function, deal with all the Request log output
	m.Use(logger("dev"))
	//Set recovery handler to returns a middleware that recovers from any panics
	m.Use(macaron.Recovery())
}
