package codec

import (
	"fmt"
	"reflect"

	"github.com/ugorji/go/codec"
)

var msgpackHandle = &codec.MsgpackHandle{}

func init() {
	msgpackHandle.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

type MsgPack struct {}

func (this MsgPack) Decode(b []byte, m interface{}) error {
	return codec.NewDecoderBytes(b, msgpackHandle).Decode(m)
}
func (this MsgPack) Encode(m interface{}) ([]byte, error) {
	var bytes []byte
	err := codec.NewEncoderBytes(&bytes, msgpackHandle).Encode(m)
	if err == nil {
		return bytes, nil
	}
	err = fmt.Errorf("tried to decode into %T struct %+v error: %v", m, m, err)
	return nil, err
}
func (this MsgPack) Name() string {
	return "application/msgpack"
}
