package codec

type Codec interface {
	Decode([]byte, interface{}) error
	Encode(interface{}) ([]byte, error)
	Name() string
}

var StdCodecs []Codec = []Codec{
	Json{},
	MsgPack{},
}