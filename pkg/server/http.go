package server

import "qqbot-go/config"

//Http服务器

type HttpServer struct {
	host string
	port uint8
}

func (h *HttpServer) Init() {

}

func (h *HttpServer) Start() {

}

func (h *HttpServer) Shutdown() {

}

func newHttpServer(conf *config.Config) (*Server, error) {
	return interface{}(&HttpServer{}).(*Server), nil
}
