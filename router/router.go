package router

import (
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/handler"
)

func SetRouters(m *macaron.Macaron) {
	m.Get("/ping", handler.Ping)
}
