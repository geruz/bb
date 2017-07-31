package bb

import (
	"github.com/geruz/bb/codec"
	h "github.com/geruz/bb/transport/protocols/http"
	"github.com/geruz/bb/transport/protocols/ws"
)

func NewWsServer(name string, version Version, port int, host string, exts ...h.Extension) *BBServer {
	server := NewBBServer(name, version)
	server.AddCodec(codec.MsgPack{})
	server.AddCodec(codec.Json{})
	httpFactory := h.HttpFactory{
		Port:       port,
		Host:       host,
		Extensions: exts,
	}
	server.AddTransport(ws.WsFactory{httpFactory})
	return server
}
