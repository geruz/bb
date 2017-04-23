package transport

import (
	"fmt"
	"github.com/geruz/bb/codec"
	"github.com/valyala/fasthttp"
)

type HttpQueryProvider struct {
	ctx   *fasthttp.RequestCtx
	codec codec.Codec
}

func (this HttpQueryProvider) In(obj interface{}) error {
	data, err := this.bytes()
	if err != nil {
		return err
	}
	err = this.codec.Decode(data, obj)
	if err != nil {
		return err
	}
	return nil
	// TODO remoteFromHeader(this.req.Header, this.ioConverter, obj)
}
func (this HttpQueryProvider) bytes() ([]byte, error) {
	method := this.ctx.Method()
	if string(method) == "GET" {
		val := this.ctx.FormValue("query")
		if len(val) == 0 {
			return val, fmt.Errorf("нет объекта запроса")
		}
		return val, nil
	} else {
		fmt.Println("Body:", string(this.ctx.Request.Body()))
		return this.ctx.Request.Body(), nil
	}
	return nil, fmt.Errorf("method %n not imlemented", method)
}
