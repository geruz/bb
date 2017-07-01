package exts

import (
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
	h "github.com/geruz/bb/transport/protocols/http"
	"fmt"
)

func a(tr *h.HttpTransport) {


}
func aa(ctx *fasthttp.RequestCtx) {
	websocket.Upgrade(ctx.Response, ctx.Request,)
	answer := fmt.Sprintf("%v %v:%v", "OK", name, version.String())
	ctx.SetBody([]byte(answer))
	ctx.SetStatusCode(fasthttp.StatusOK)
}