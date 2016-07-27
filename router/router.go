package router

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/handler"
	"github.com/containerops/arkor/handler/internal"
	"github.com/containerops/arkor/models"
)

func SetRouters(m *macaron.Macaron) {
	m.Group("/v1", func() {
		m.Get("/", handler.GetServiceHandler)
		m.Put("/:bucket", handler.PutBucketHandler)
		m.Head("/:bucket", handler.HeadBucketHandler)
		m.Delete("/:bucket", handler.DeleteBucketHandler)
		m.Get("/:bucket", handler.GetBucketHandler)
	})
	// internal APIS
	m.Group("/internal", func() {
		m.Group("/v1", func() {
			m.Put("/dataserver", binding.Bind(models.DataServer{}), internal.PutDataserverHandler)
		})
	})
	// interface to test whether the arkor is working
	m.Get("/ping", handler.Ping)
}
