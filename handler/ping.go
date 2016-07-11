package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

func Ping(ctx *macaron.Context) (int, string) {
	return http.StatusOK, "Pong! \n"
}
