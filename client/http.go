package client

import (
	"fmt"

	"github.com/geruz/bb/discovery"
	"github.com/valyala/fasthttp"
)

var client = &fasthttp.Client{}

type HttpTransport struct {
	Resource string
	Action   string
	Major    int
	Minor    int
}

func (this HttpTransport) Send(address discovery.Address, data []byte, meta Meta, chans AnswerChans) {
	url := fmt.Sprintf("http://%v:%v/v%v/%v/%v/", address.Host, address.Port, this.Major, this.Resource, this.Action)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	for key, value := range meta {
		req.Header.Set(key, value)
	}
	req.SetRequestURI(url)
	req.AppendBody(data)
	resp := fasthttp.AcquireResponse()

	err := client.Do(req, resp)
	if err != nil {
		chans.Error <- []byte(err.Error())
	}
	body := resp.Body()
	statusCode := resp.StatusCode()
	switch statusCode {
	case fasthttp.StatusOK:
		chans.Success <- body
		return
	case fasthttp.StatusNotFound:
		chans.NotFound <- body
		return
	case fasthttp.StatusInternalServerError:
		chans.Error <- body
		return
	case fasthttp.StatusNotImplemented:
		chans.NotIplemented <- body
		return
	case fasthttp.StatusBadRequest:
		chans.ValidateError <- body
		return
	default:
		chans.Error <- []byte(fmt.Sprintf("Bad status code :v", statusCode))
	}
}
