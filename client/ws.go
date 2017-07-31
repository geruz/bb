package client

import (
	"fmt"
	"github.com/geruz/bb/discovery"
	"github.com/geruz/bb/transport/protocols/ws"
)

type WsTransport struct {
	Pool WsPool
	Resource string
	Action   string
	Major    int
	Minor    int
}

func (this WsTransport) Send (address discovery.Address, data []byte, ans AnswerChans )  {
	wss, err := this.Pool.get(address.Host, address.Port)
	if err != nil{
		ans.Error <- []byte(err.Error())
		return
	}
	meta := ws.Meta {
		this.Major, this.Resource, this.Action,
	}

	answer := wss.Send(meta, data, ans)
	switch answer.Code {
	case 200:
		ans.Success <- answer.Body
	case 404:
		ans.NotFound <- answer.Body
	case 501:
		ans.NotIplemented <- answer.Body
	case 500:
		ans.Error <- answer.Body
	default:
		ans.Error <- []byte(fmt.Sprintf("Bad status code :%v", answer.Code))
	}
}

