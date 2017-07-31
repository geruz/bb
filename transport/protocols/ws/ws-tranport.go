package ws

import (
	h "github.com/geruz/bb/transport/protocols/http"
	"strings"
	"fmt"
	"github.com/geruz/bb/resource"
	"sync"
)

type WsTransport struct {
	httpTransport *h.HttpTransport
}

func (this WsTransport)Start(){
	actions := map[string]resource.Action{}
	for _, resource := range this.httpTransport.Configuration.Handlers {
		for _, action := range resource.Actions {
			p := strings.ToLower(resource.Name + "/" + action.Name)
			fmt.Println("Register ws path: ", p)
			actions[p] = action
		}
	}
	defCodec := this.httpTransport.DefCodec
	locker := &sync.Mutex{}
	this.httpTransport.AddHandler("/ws/proxy/", Upgrade(func(ws WsContext){
		defer ws.Close()
		for {
			select {
			case data, ok := <-ws.ReadNext:
				if !ok {
					return
				}
				parts, err := ParseParts(data)
				if err != nil{
					//TODO  log
					fmt.Println(err)
					return
				}
				locker.Lock()
				action, ok := actions[parts.Meta.Resource + "/" + parts.Meta.Action]
				locker.Unlock()
				resProvider := WsResultProvider{ws, defCodec, parts.Id}
				if !ok {
					resProvider.NotImplemented()
					continue
				}
				func () {
					defer resProvider.Recover()
					answ, err := action.Exec(simpleProvider{parts.Data, defCodec})
					if err == nil {
						resProvider.Success(answ)
						return
					}
					resProvider.AutoErr(err)
				}()
			}
		}
	}))
	this.httpTransport.Start()
}

type Answer struct {
	Code int
	Body []byte
}

func ParseAnswer(code int64 , data []byte) Answer{
	return Answer{int(code), data}
}