package bb

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
	"github.com/geruz/bb/transport"
)

type Version struct {
	Major int
	Minor int
	Patch int
}
type BBServer struct {
	Name      string
	Version   Version
	codecs    []codec.Codec
	factories []transport.Factory
	handlers  []resource.Handler
}

func NewBBServer(name string, version Version) *BBServer {
	server := BBServer{
		Name:    name,
		Version: version,
	}
	return &server
}

func (this *BBServer) AddTransport(tr transport.Factory) {
	this.factories = append(this.factories, tr)
}

func (this *BBServer) AddResource(name string, factory func() interface{}) {
	this.handlers = append(this.handlers, resource.NewResourceRunner(name, factory))
}

func (this *BBServer) AddCodec(cnv codec.Codec) {
	this.codecs = append(this.codecs, cnv)
}

func (this *BBServer) Loop() {
	conf := transport.NewConfiguration(
		this.Name,
		this.Version.Major, this.Version.Major, this.Version.Patch,
		this.handlers,
		this.codecs,
	)
	for _, factory := range this.factories {
		tr := factory.Create(conf)
		go tr.Start()
	}
}
