package transport

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
	h "github.com/geruz/bb/transport/http"
)

type Transport interface {
	Start()
}

type Factory interface {
	Create(version Version, handlers []resource.Handler, codecs []codec.Codec) Transport
}

type Version struct {
	Major int
	Minor int
	Patch int
}

type HttpFactory struct {
	Port int
	Host string
}

func (this HttpFactory) Create(version Version, handlers []resource.Handler, codecs []codec.Codec) Transport {
	return &h.Server{
		Port:     this.Port,
		Host:     this.Host,
		DefCodec: codec.Json{},
		Major:    version.Major,
		Handlers: handlers,
		Codecs:   codecs,
	}
}
