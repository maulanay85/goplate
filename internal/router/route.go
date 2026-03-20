package router

import (
	"goplate/internal/controller"
	"goplate/pkg/server"
)

func Route(s *server.Server) {
	s.GET("/ping", controller.Ping)
	s.POST("/short", controller.Short)
}
