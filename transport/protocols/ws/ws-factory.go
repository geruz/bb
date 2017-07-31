package ws

import (
	"github.com/geruz/bb/transport/configuration"
	"github.com/geruz/bb/transport/protocols"

	h "github.com/geruz/bb/transport/protocols/http"
)

type WsFactory struct {
	HttpFactory h.HttpFactory
}

func (this WsFactory) Create(configuration configuration.Configuration) protocols.Transport {
	httpTransport := this.HttpFactory.Create(configuration).(*h.HttpTransport)
	return &WsTransport {
		httpTransport: httpTransport,
	}
}
