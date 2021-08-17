package goal

const (
	HttpServer ServerType = iota
	WSServer
)

type ServerType int8
