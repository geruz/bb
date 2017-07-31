package main

import (
	"github.com/geruz/bb"
	"github.com/geruz/bb/example/server/exts"
	"github.com/geruz/bb/resource"
	"time"
)

type Contr struct {
	resource.Context
}

func (this Contr) Echo(query string) (string, error) {
	return query, nil
}

type Query struct {
	Id int
}

func main() {
	server := bb.NewWsServer(
		"echo-server",
		bb.Version{1, 1, 1},
		8088,
		"localhost",
		exts.PingExtension{},
		exts.RamlExtension{},
		//exts.LogExtension{},
	)
	server.AddResource(`test`, func() interface{} {
		return Contr{}
	})
	server.Loop()
	for {
		time.Sleep(time.Second)
	}
}
