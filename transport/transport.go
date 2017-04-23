package transport

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
	"github.com/geruz/bb/transport/configuration"
	"github.com/geruz/bb/transport/protocols"
)

func NewConfiguration(name string, major, minor, patch int, handlers []resource.Handler, codecs []codec.Codec) configuration.Configuration {
	return configuration.Configuration{
		Name:     name,
		Version:  configuration.Version{major, minor, patch},
		Codecs:   codecs,
		Handlers: handlers,
		Control:  configuration.NewControl(),
	}
}

type Factory interface {
	Create(configuration configuration.Configuration) protocols.Transport
}
