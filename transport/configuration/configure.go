package configuration

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

type Configuration struct {
	Name     string
	Version  Version
	Handlers []resource.Handler
	Codecs   []codec.Codec
	Control  Control
}
