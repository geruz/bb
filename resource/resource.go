package resource

import (
	"reflect"
	"strings"
	"unicode"
)

type Handler struct {
	Name    string
	Actions []Action
}

type Context struct{}

func NewResourceRunner(name string, factory func() interface{}) Handler {
	name = strings.ToLower(name)
	actions := []Action{}
	controller := reflect.TypeOf(factory())
	for i := 0; i < controller.NumMethod(); i++ {
		method := controller.Method(i)
		if strings.IndexFunc(method.Name[0:1], unicode.IsUpper) != 0 {
			continue
		}
		actions = append(actions, NewAction(factory, method))
	}
	return Handler{
		Name:    name,
		Actions: actions,
	}
}
