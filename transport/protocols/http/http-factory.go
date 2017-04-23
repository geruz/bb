package transport

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/transport/configuration"
	"github.com/geruz/bb/transport/protocols"
)

type HttpFactory struct {
	Port int
	Host string
}

func (this HttpFactory) Create(configuration configuration.Configuration) protocols.Transport {
	return &HttpTransport{
		Port:          this.Port,
		Host:          this.Host,
		DefCodec:      codec.Json{},
		Configuration: configuration,
	}
}
