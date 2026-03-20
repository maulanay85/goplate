package controller

import (
	"goplate/internal/model"
	"goplate/pkg/server"
)

func Ping(ctx *server.Context) {
	ctx.OK(map[string]string{
		"message": "pong",
	})
}

func Short(ctx *server.Context) {
	var data model.Short

	if err := ctx.BindJSON(&data); err != nil {
		ctx.InternalServerError(map[string]string{
			"message": "error",
		})
		return
	}

	ctx.OK(map[string]string{
		"data": data.Url,
	})
}
