package codec

import (
	"fmt"
	"reflect"

	"github.com/ugorji/go/codec"
)

var jsonHandle = &codec.JsonHandle{}

func init() {
	jsonHandle.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

type Json struct {}

func (this Json) Decode(b []byte, m interface{}) error {
	err := codec.NewDecoderBytes(b, jsonHandle).Decode(m)
	if err != nil {
		err = fmt.Errorf("tried to decode into %T struct %+v error: %v", m, m, err)
	}
	return err
}
func (this Json) Encode(m interface{}) ([]byte, error) {
	var bytes []byte
	err := codec.NewEncoderBytes(&bytes, jsonHandle).Encode(m)
	if err == nil {
		return bytes, nil
	}
	return nil, err
}

func (this Json) Name() string {
	return "application/json"
}
