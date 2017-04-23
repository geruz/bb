package exts

import (
	"fmt"
	"github.com/alecthomas/rapid"
	"github.com/geruz/bb/resource"
	h "github.com/geruz/bb/transport/protocols/http"
	"github.com/valyala/fasthttp"
	"reflect"
	"strings"
)

type RamlExtension struct {
	builder RamlBuilder
}

func (this RamlExtension) Configure(tr *h.HttpTransport) {
	builder := newRamlBuilder(tr.Configuration.Name, tr.Configuration.Handlers)
	tr.AddHandler("/_/raml/", func(ctx *fasthttp.RequestCtx) {
		rapid.SchemaToRAML(fmt.Sprintf("http://%v:%v/v%v/", tr.Host, tr.Port, tr.Configuration.Version.Major), builder.Build(), ctx)
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
}

type RamlBuilder interface {
	Build() *rapid.Schema
}

func newRamlBuilder(name string, handlers []resource.Handler) RamlBuilder {
	builder := rapid.Define(name)
	for _, handler := range handlers {
		resourcePath := strings.ToLower(fmt.Sprintf("%v/", handler.Name))

		ramlResource := builder.Resource(handler.Name, resourcePath)
		for _, action := range handler.Actions {
			actionPath := strings.ToLower(resourcePath + action.Name + "/")
			route := ramlResource.Route(handler.Name+" "+action.Name, actionPath).Post()
			if action.Out != nil {
				route.Response(200, reflect.Zero(action.Out).Interface())
			}
			if action.In != nil {
				q := reflect.Zero(action.In).Interface()
				route.Request(q)
			}
		}
	}
	return builder
}
