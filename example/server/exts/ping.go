package exts

import (
	"fmt"
	h "github.com/geruz/bb/transport/protocols/http"
	"github.com/valyala/fasthttp"
)

type PingExtension struct{}

func (this PingExtension) Configure(tr *h.HttpTransport) {

	name := tr.Configuration.Name
	version := tr.Configuration.Version
	tr.AddHandler("/_/ping/", func(ctx *fasthttp.RequestCtx) {
		answer := fmt.Sprintf("%v %v:%v", "OK", name, version.String())
		ctx.SetBody([]byte(answer))
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
}
