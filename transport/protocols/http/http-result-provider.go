package transport

import (
	"github.com/geruz/bb/codec"
	"github.com/geruz/bb/results"
	"github.com/valyala/fasthttp"
)

type HttpResultProvider struct {
	Cntx  *fasthttp.RequestCtx
	Codec codec.Codec
}

func (this HttpResultProvider) send(code int, data interface{}) {
	bytes, err := this.Codec.Encode(data)
	if err != nil {
		bytes, _ = this.Codec.Encode(data)
		this.Cntx.SetBody(bytes)
		this.Cntx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	this.Cntx.SetStatusCode(code)
	this.Cntx.SetBody(bytes)
}
func (this HttpResultProvider) Success(data interface{}) {
	this.send(fasthttp.StatusOK, data)
}
func (this HttpResultProvider) NotImplemented() {
	this.send(fasthttp.StatusNotImplemented, "Not implemented")
}
func (this HttpResultProvider) NotFound(err interface{}) {
	this.send(fasthttp.StatusNotFound, err)
}
func (this HttpResultProvider) Validation(err interface{}) {
	this.send(fasthttp.StatusBadRequest, err)
}
func (this HttpResultProvider) ServerError(err interface{}) {
	this.send(fasthttp.StatusInternalServerError, err)
}
func (this HttpResultProvider) Recover() {
	if err := recover(); err != nil {
		this.AutoErr(err)
	}
}
func (this HttpResultProvider) AutoErr(err interface{}) {
	switch e := err.(type) {
	case results.NotFound, *results.NotFound:
		this.NotFound(err)
		break
	case results.ValidationError, *results.ValidationError:
		this.Validation(err)
		break
	case results.Error, *results.Error:
		this.ServerError(err)
		break
	case error:
		this.ServerError(e.Error())
		break
	default:
		this.ServerError(err)
		break
	}
}
