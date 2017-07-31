package main

import (
	"fmt"
	"github.com/geruz/bb/client"
	"github.com/geruz/bb/discovery"
)
const n = 100000

func main() {
	cl := client.Client{
		Provider:  func() (discovery.Address, error){
			return discovery.Address{
				"localhost", 8088,
			}, nil
		},
		Transport: client.WsTransport{
			Pool: client.NewWsPool(),
			Resource: `test`,
			Action: "echo",
			Major: 1,
			Minor: 0,
		},
	}
	for i := 0; i < n; i++ {
			rs := cl.Call([]byte(fmt.Sprintf("\"echo: %v\"", i)))
			select {
			case r := <-rs.Success:
				fmt.Println("Success", string(r))
				break
			case r := <-rs.Error:
				fmt.Println("Error", string(r))
				break
			case r := <-rs.NotFound:
				fmt.Println("NotFound", string(r))
				break
			case r := <-rs.NotIplemented:
				fmt.Println("NotIplemented", string(r))
				break
			}
			//time.Sleep(time.Second)
	}
}
