package client

import (
	"fmt"
	"github.com/geruz/bb/discovery"
	"github.com/valyala/fasthttp"
)

type Client struct {
	Provider  discovery.AdddressProvider
	Transport Transport
}
type AnswerChans struct {
	Success       chan []byte
	Error         chan []byte
	NotFound      chan []byte
	NotIplemented chan []byte
}

func NewAnswerChans() AnswerChans {
	return AnswerChans{
		Success:       make(chan []byte),
		Error:         make(chan []byte),
		NotFound:      make(chan []byte),
		NotIplemented: make(chan []byte),
	}
}

type Transport interface {
	Call(address discovery.Address, data []byte, chans AnswerChans, errCh chan error)
}
type HttpTransport struct {
	Resource string
	Action   string
	Major    int
	Minor    int
}

func (this HttpTransport) Call(address discovery.Address, data []byte, chans AnswerChans, errCh chan error) {
	url := fmt.Sprintf("http://%v:%v/v%v/%v/%v/", address.Host, address.Port, this.Major, this.Resource, this.Action)
	fmt.Println(string(data))
	statusCode, body, err := fasthttp.Post(data, url, nil)
	if err != nil {
		errCh <- err
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
		errCh <- fmt.Errorf("Bad ststau code :v", statusCode)
	}
}

func (this *Client) Call(data []byte) (AnswerChans, chan error) {
	errCh := make(chan error)
	chans := NewAnswerChans()
	go func() {
		address, err := this.Provider()
		if err != nil {
			errCh <- err
		}
		this.Transport.Call(address, data, chans, errCh)
	}()
	return chans, errCh
}
