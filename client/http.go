package client

import (
	"fmt"
	"github.com/geruz/bb/discovery"
	"github.com/valyala/fasthttp"
)


type HttpTransport struct {
	Resource string
	Action   string
	Major    int
	Minor    int
}

func (this HttpTransport) Call(address discovery.Address, data []byte, chans AnswerChans) {
	url := fmt.Sprintf("http://%v:%v/v%v/%v/%v/", address.Host, address.Port, this.Major, this.Resource, this.Action)
	fmt.Println(string(data))
	statusCode, body, err := fasthttp.Post(data, url, nil)
	if err != nil {
		chans.Error <- []byte(err.Error())
		return
	}
	switch statusCode {
	case 200:
		chans.Success <- body
		return
	case 404:
		chans.NotFound <- body
		return
	case 500:
		chans.Error <- body
		return
	case 501:
		chans.NotIplemented <- body
		return
	default:
		chans.Error <- []byte(fmt.Sprintf("Bad status code :v", statusCode))
	}
}