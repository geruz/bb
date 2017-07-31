package client

import (
	"fmt"
	"github.com/geruz/bb/transport/protocols/ws"
	"sync"
	"net"
	"bufio"
)

type WsPool struct {
	opened map[string] *WsReadWrite
}
func NewWsPool()WsPool{
	return WsPool{
		opened: map[string] *WsReadWrite{},
	}
}
func (this WsPool) get(host string, port int)(*WsReadWrite, error) {
	//TODO lock
	address := fmt.Sprintf("%v:%v", host, port)
	if wss, ok := this.opened[address]; ok {
		return wss, nil
	}
	wss, err := this.open(address)
	if err != nil {
		return nil, err
	}
	this.opened[address] = wss
	return wss, nil
}

func (this WsPool) open(address string)(*WsReadWrite, error) {

	//TODO go на ping pong
	//TODO defer wsContext.Close()
	conn, err := this.getWsConn(address)
	if err != nil {
		return nil, err
	}
	wsContext := ws.NewWsContext(conn);
	return NewWsReadWrite(wsContext), nil
}

func (this WsPool)getWsConn(address string) (net.Conn, error){
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(conn, `GET ws://%v/ws/proxy/ HTTP/1.1
Host: %v
Connection: Upgrade
Upgrade: websocket
Origin: http://%v
Sec-WebSocket-Version: 13
Sec-WebSocket-Key: njOpWRHDp4AxY+EPDuvSvg==

`, address, address, address)
	//TODO добавить случайный ключ
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if message == "\r\n" {
			break;
		}
	}
	//TODO проверка на заголовки UPGRADE
	return conn, nil
}

type WsReadWrite struct {
	chans ChanMap
	mutext *sync.Mutex
	context ws.WsContext
}
func NewWsReadWrite(wsContext ws.WsContext) *WsReadWrite {

	chans := NewChanMap()

	go func() {
		for {
			next := <-wsContext.ReadNext
			res, err := ws.ParseResponce(next)
			if err != nil {
				return
			}
			if ch, ok := chans.GetAndRemove(res.Id); ok {
				ch <- ws.ParseAnswer(res.Code, res.Data)
				close(ch)
				continue
			}
			//TODO log
			fmt.Printf("Неожиданный id %v", res)
		}
	}()
	return &WsReadWrite{
		chans:chans,
		mutext: &sync.Mutex{},
		context: wsContext,
	}
}

func (this *WsReadWrite) Send(meta ws.Meta, data []byte, chans AnswerChans) ws.Answer{
	result := make(chan ws.Answer)
	func() {
		this.mutext.Lock()
		defer this.mutext.Unlock()
		parts := ws.NewRequestParts(meta, data)
		this.chans.Set(parts.Id, result)
		this.context.Write(parts.ToBytes())
	}()
	return <-result
}


