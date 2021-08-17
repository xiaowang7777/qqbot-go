package server

import (
	"qqbot-go/config"
	"qqbot-go/pkg/goal"
)

type Server interface {
	Init()
	Start()
	Shutdown()
}

func NewServer(conf *config.Config) (*Server, error) {
	switch conf.Server.Type {
	case goal.HttpServer:
		return newHttpServer(conf)
	case goal.WSServer:
		return newWSServer(conf)
	}
	return nil, nil
}
