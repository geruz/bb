package ws

import (
	"bytes"
	"fmt"
	"encoding/binary"
	"sync"
)

type requestParts struct {
	Id   uint64
	Meta Meta
	Data []byte
}

func ParseParts(body []byte )(requestParts, error){
	id, _ := binary.Uvarint(body[0:8])
	parts := bytes.SplitN(body[8:], selector, 2)
	meta, err := ParseMeta(parts[0])
	if err != nil {
		return requestParts{},  fmt.Errorf("not parsed meta: %v", err.Error())
	}
	return requestParts{id, meta, parts[1]}, nil
}


var generater = NewGenId(0)
var selector = []byte("\n")

func (this requestParts) ToBytes() []byte {
	idBuff := make([]byte,8)
	binary.PutUvarint(idBuff, this.Id)
	s := [][]byte{idBuff, this.Meta.ToBytes(), this.Data}
	return bytes.Join(s,[]byte{})
}
func NewRequestParts(meta Meta, data []byte) requestParts {
	id := generater.NextId()
	return requestParts {
		id,meta, data,
	}
}

func NewGenId(start uint64) genId {
	return genId{
		id: start,
		mutex: &sync.Mutex{},
	}
}
func (this *genId) NextId() uint64{
	this.mutex.Lock()
	this.id++;
	i := this.id
	this.mutex.Unlock()
	return i;
}

type genId struct{
	id uint64
	mutex *sync.Mutex
}
