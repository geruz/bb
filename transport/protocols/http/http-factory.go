package transport

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/transport/configuration"
	"github.com/geruz/bb/transport/protocols"
	"github.com/valyala/fasthttp"
)

type HttpFactory struct {
	Port       int
	Host       string
	Extensions []Extension
}

func (this HttpFactory) Create(configuration configuration.Configuration) protocols.Transport {
	return &HttpTransport{
		Port:          this.Port,
		Host:          this.Host,
		DefCodec:      codec.Json{},
		Configuration: configuration,
		Extension:     this.Extensions,
		handlers:      map[string]func(ctx *fasthttp.RequestCtx, resProvider HttpResultProvider){},
	}
}
