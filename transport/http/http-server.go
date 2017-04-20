package transport

import (
	"fmt"
	"strings"

	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/resource"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Port     int
	Host     string
	DefCodec codec.Codec
	Major    int
	Handlers []resource.Handler
	Codecs   []codec.Codec
}

func (this *Server) Start() {
	if this.DefCodec == nil {
		this.DefCodec = codec.Json{}
	}
	paths := map[string]func(ctx *fasthttp.RequestCtx, resProvider HttpResultProvider){}
	for _, resource := range this.Handlers {
		for _, action := range resource.Actions {
			p := this.getPath(resource.Name, action.Name)
			fmt.Println("Register path: ", p)
			paths[p] = func(ctx *fasthttp.RequestCtx, resProvider HttpResultProvider) {
				defer resProvider.Recover()
				answ, err := action.Exec(HttpQueryProvider{
					ctx, this.DefCodec,
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
	err := fasthttp.ListenAndServe(address, func(ctx *fasthttp.RequestCtx) {
		resProvider := HttpResultProvider{ctx, this.DefCodec}
		path := ctx.Path()
		for p, f := range paths {
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

func (this *Server) getPath(resource string, action string) string {
	path := fmt.Sprintf("/v%v/%v/%v/", this.Major, resource, action)
	path = strings.ToLower(path)
	return path
}
