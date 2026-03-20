package controller

import "goplate/pkg/server"

func Ping(ctx *server.Context) {
	ctx.OK(map[string]string{
		"message": "pong",
	})
}
