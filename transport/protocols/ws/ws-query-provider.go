package ws

import "github.com/geruz/bb/codec"

type simpleProvider struct {
	data []byte
	codec codec.Codec
}
func (this simpleProvider) In(i interface{}) error{
	return this.codec.Decode(this.data, i)
}