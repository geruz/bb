package main

import (
	"fmt"
	"github.com/geruz/bb/client"
	"github.com/geruz/bb/discovery"
)

func main() {
	address := "localhost"
	port := 8088
	cl := client.Client{
		Provider:  discovery.UseSimpleProvider(&address, &port),
		Transport: client.HttpTransport{Resource: `"test"`, Action: "echo", Major: 0, Minor: 0},
	}
	rs, err := cl.Call([]byte("echo"))
	select {
	case e := <-err:
		fmt.Println("e", e)
		return
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
}
