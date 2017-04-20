package resource

import (
	"fmt"
	"reflect"
)

type Action struct {
	Name string
	Exec exec
	In   reflect.Type
	Out  reflect.Type
}

type ActionRunner struct {
	controller reflect.Type
	method     reflect.Method
}

func NewAction(factory func() interface{}, method reflect.Method) Action {

	fmt.Println(method.Name)
	var in reflect.Type
	var out reflect.Type
	if method.Type.NumIn() > 2 {
		panic(fmt.Sprintf("Входной параметр должен быть один, %v", method.Type.NumIn()))
	}
	if method.Type.NumIn() == 2 {
		in = method.Type.In(1)
	} else {
		in = nil
	}
	if method.Type.NumOut() != 2 {
		panic(fmt.Sprintf("Функции должна возвращать результат и ошибку, %v", method.Type.NumIn()))
	}
	//TODO добавить проверку что второй параметр ошибка
	if method.Type.NumOut() == 1 {
		out = method.Type.Out(0)
	} else {
		out = nil
	}
	return Action{
		Name: method.Name,
		Exec: createExec(factory, in, method),
		In:   in,
		Out:  out,
	}
}

type exec func(provider InProvider) (interface{}, error)

func createExec(factory func() interface{}, in reflect.Type, method reflect.Method) exec {
	return func(provider InProvider) (interface{}, error) {
		if provider == nil {
			return nil, fmt.Errorf("Provider is nil")
		}
		inst := factory()
		sPtr := reflect.New(in).Interface()
		if err := provider.In(sPtr); err != nil {
			return nil, err
		}
		params := []reflect.Value{reflect.Indirect(reflect.ValueOf(sPtr))}
		rs := reflect.ValueOf(inst).MethodByName(method.Name).Call(params)
		if rs[1].Interface() == nil {
			return rs[0].Interface(), nil
		}
		return nil, (rs[1].Interface().(error))
	}
}
