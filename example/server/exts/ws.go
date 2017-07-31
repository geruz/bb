package exts

/*
import (
	"github.com/geruz/bb/resource"
	"strings"
	"fmt"
	"sync"
)


import (
	h "github.com/geruz/bb/transport/protocols/http"
	"fmt"
	"strings"
	"github.com/geruz/bb/resource"
	"sync"
)



type LogExtension struct {}



func (this LogExtension) Configure(tr *h.HttpTransport) {
	actions := map[string]resource.Action{}
	for _, resource := range tr.Configuration.Handlers {
		for _, action := range resource.Actions {
			p := strings.ToLower(resource.Name + "/" + action.Name)
			fmt.Println("Register ws path: ", p)
			actions[p] = action
		}
	}
	cdc := tr.DefCodec
	locker := &sync.Mutex{}

	tr.AddHandler("/_/log/", Upgrade(func(ws WsContext){
		for {
			select {
			case data := <-ws.ReadNext:
				d := string(data)
				parts := strings.SplitN(d,":", 2)
				actionName := parts[0]
				q:= parts[1]
				locker.Lock()
				action, ok := actions[actionName]
				locker.Unlock()
				if !ok {
					fmt.Println(d, actionName)
				}
				resProvider := WsResultProvider{ws, cdc}
				func () {
					defer resProvider.Recover()

					answ, err := action.Exec(simpleProvider{[]byte(q),cdc})
					if err == nil {
						resProvider.Success(answ)
						return
					}
					resProvider.AutoErr(err)
				}()
			case <- ws.CloseCh:
				return
			}
		}
	}))
}

*/