package main

import (
	"github.com/geruz/bb"
	"github.com/geruz/bb/resource"
	"time"
)

type Contr struct {
	resource.Context
}

func (this Contr) Echo(query string) (string, error) {
	return query, nil
}

type Contr1 struct {
	resource.Context
}

func (this Contr1) Echo(query string) string {
	return query
}

type Query struct {
	Id int
}

func main() {
	server := bb.NewHttpServer("echo-server", bb.Version{1, 1, 1}, 8088, "localhost")
	server.AddResource(`test`, func() interface{} {
		return Contr{}
	})
	server.Loop()
	for {
		time.Sleep(time.Second)
	}
}
