package bb

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/transport"
)

func NewHttpServer(port int, host string) *BBServer {
	server := BBServer{}
	server.AddFormat(codec.MsgPack{})
	server.AddFormat(codec.Json{})
	httpTransport := transport.HttpFactory{
		Port: port,
		Host: host,
	}
	server.AddTransport(&httpTransport)
	return &server
}
