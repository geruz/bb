package transport

import (
	"fmt"
	"strings"

	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/transport/configuration"
	"github.com/valyala/fasthttp"
)

type HttpTransport struct {
	Port          int
	Host          string
	DefCodec      codec.Codec
	Configuration configuration.Configuration
	Extension     []Extension
	handlers      map[string]func(ctx *fasthttp.RequestCtx, resProvider HttpResultProvider)
}

func findCodec(name []byte, defCodec codec.Codec) codec.Codec{
	n := string(name)
	for _,c := range codec.StdCodecs{
		if c.Name() == n{
			return c
		}
	}
	return defCodec
}
func (this *HttpTransport) Start() {
	if this.DefCodec == nil {
		this.DefCodec = codec.Json{}
	}
	for _, resource := range this.Configuration.Handlers {
		for _, action := range resource.Actions {
			p := this.getPath(resource.Name, action.Name)
			//TODO log
			fmt.Println("Register path: ", p)
			this.handlers[p] = func(ctx *fasthttp.RequestCtx, resProvider HttpResultProvider) {
				defer resProvider.Recover()
				codec := findCodec(ctx.Request.Header.ContentType(), this.DefCodec)
				answ, err := action.Exec(HttpQueryProvider{
					ctx, codec,
				})
				if err == nil {
					resProvider.Success(answ)
					return
				}
				resProvider.AutoErr(err)

			}
		}
	}
	address := fmt.Sprintf("%v:%v", this.Host, this.Port)
	for _, ext := range this.Extension {
		ext.Configure(this)
	}
	err := fasthttp.ListenAndServe(address, func(ctx *fasthttp.RequestCtx) {
		resProvider := HttpResultProvider{ctx, this.DefCodec}
		path := ctx.Path()
		for p, f := range this.handlers {
			if p == string(path) {
				f(ctx, resProvider)
				return
			}
		}
		resProvider.NotImplemented()

	})
	if err != nil {
		panic(err)
	}
}
func (this *HttpTransport) AddHandler(name string, handler func(ctx *fasthttp.RequestCtx)) error {
	this.handlers[name] = func(ctx *fasthttp.RequestCtx, resProvider HttpResultProvider) {
		handler(ctx)
	}
	return nil
}

func (this *HttpTransport) getPath(resource string, action string) string {
	path := fmt.Sprintf("/v%v/%v/%v/", this.Configuration.Version.Major, resource, action)
	path = strings.ToLower(path)
	return path
}
