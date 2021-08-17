package server

import "qqbot-go/config"

type WSServer struct {
}

func (w *WSServer) Init() {

}

func (w *WSServer) Start() {

}

func (w *WSServer) Shutdown() {

}

func newWSServer(conf *config.Config) (*Server, error) {
	return interface{}(&WSServer{}).(*Server), nil
}
