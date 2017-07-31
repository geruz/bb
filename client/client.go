package client

import "github.com/geruz/bb/discovery"

type Client struct {
	Provider  discovery.AdddressProvider
	Transport Transport
}

type Transport interface {
	Send(address discovery.Address, data []byte, chans AnswerChans)
}

func (this *Client) Call(data []byte) (AnswerChans) {
	chans := NewAnswerChans()
	go func() {
		address, err := this.Provider()
		if err != nil {
			chans.Error <- []byte(err.Error())
			return
		}
		this.Transport.Send(address, data, chans)
	}()
	return chans
}
