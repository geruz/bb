package configuration

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
	"strconv"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (this Version) String() string {
	return strconv.Itoa(this.Major) + "." + strconv.Itoa(this.Minor) + "." + strconv.Itoa(this.Patch)
}

type Configuration struct {
	Name     string
	Version  Version
	Handlers []resource.Handler
	Codecs   []codec.Codec
	Control  Control
}
