package handler

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

func Ping(ctx *macaron.Context, log *logrus.Logger) (int, string) {
	log.Info("Pong")
	return http.StatusOK, "Pong! \n"
}
