package bb

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
	"github.com/geruz/bb/transport"
)

type Version struct {
	transport.Version
}
type BBServer struct {
	Name      string
	Version   Version
	codecs    []codec.Codec
	factories []transport.Factory
	handlers  []resource.Handler
}

func (this *BBServer) AddTransport(tr transport.Factory) {
	this.factories = append(this.factories, tr)
}

func (this *BBServer) AddResource(name string, factory func() interface{}) {
	this.handlers = append(this.handlers, resource.NewResourceRunner(name, factory))
}

func (this *BBServer) AddFormat(cnv codec.Codec) {
	this.codecs = append(this.codecs, cnv)
}

func (this *BBServer) Loop() {

	for _, factory := range this.factories {
		tr := factory.Create(transport.Version{}, this.handlers, this.codecs)
		go tr.Start()
	}
}
func (this *BBServer) Stop() {
	//	this.stop <- true
}
