package main

import (
	prop "goplate/config"
	"goplate/internal/router"
	"goplate/pkg/config"
	"goplate/pkg/server"
	"log"
)

func main() {
	cfg, err := config.Load[prop.Config](
		config.WithEnvFile(".env"),
	)
	if err != nil {
		log.Fatalf("error running: %v", err)
	}

	s := server.New(server.WithAppName(cfg.AppName), server.WithEnv(cfg.Env), server.WithPort(cfg.Port))
	router.Route(s)
	s.Run()
}
