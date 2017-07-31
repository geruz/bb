package client

import "github.com/geruz/bb/discovery"

type Client struct {
	Provider  discovery.AdddressProvider
	Transport Transport
}

type Transport interface {
	Send(address discovery.Address, data []byte, meta Meta, chans AnswerChans)
}

func (this *Client) Call(data []byte, meta Meta) (AnswerChans) {
	chans := NewAnswerChans()
	go func() {
		address, err := this.Provider()
		if err != nil {
			chans.Error <- []byte(err.Error())
			return
		}
		this.Transport.Send(address, data, meta, chans)
	}()
	return chans
}

type Meta map[string]string