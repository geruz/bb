package ws

import (
	"github.com/geruz/bb/results"
	"github.com/valyala/fasthttp"
	"github.com/geruz/bb/codec"
	"encoding/binary"
	"bytes"
	"fmt"
)

type WsResultProvider struct {
	Cntx  WsContext
	Codec codec.Codec
	Id uint64
}
type Res struct {
	Id uint64
	Code int64
	Data []byte
}
func ParseResponce(data []byte) (Res, error) {
	if len(data) < 16 {
		return Res{}, fmt.Errorf("answer < 16")
	}
	fmt.Println(data)

	id, _ := binary.Uvarint(data[0:8])
	code, _ := binary.Varint(data[8:16])
	return Res{
		Id: id,
		Code: code,
		Data: data[16:],
	}, nil
}

func (this Res) ToBytes() []byte{
	idBuff := make([]byte, 8)
	binary.PutUvarint(idBuff, this.Id)
	codeBuff := make([]byte, 8)
	binary.PutVarint(codeBuff, int64(this.Code))
	buffer := [][]byte{idBuff,codeBuff, this.Data}
	fmt.Println(bytes.Join(buffer, []byte{}))
	return bytes.Join(buffer, []byte{})
}
func (this WsResultProvider) send(code int64, answer interface{}) {
	//TODO error
	data, _ := this.Codec.Encode(answer)
	res := Res {
		Id: this.Id,
		Code: code,
		Data: data,
	}
	fmt.Println(res)

	this.Cntx.Write(res.ToBytes())
}


func (this WsResultProvider) Success(data interface{}) {
	this.send(fasthttp.StatusOK, data)
}
func (this WsResultProvider) NotImplemented() {
	this.send(fasthttp.StatusNotImplemented, "Not implemented")
}
func (this WsResultProvider) NotFound(err interface{}) {
	this.send(fasthttp.StatusNotFound, err)
}
func (this WsResultProvider) Validation(err interface{}) {
	this.send(fasthttp.StatusBadRequest, err)
}
func (this WsResultProvider) ServerError(err interface{}) {
	this.send(fasthttp.StatusInternalServerError, err)
}
func (this WsResultProvider) Recover() {
	if err := recover(); err != nil {
		//TODO Log
		//log.Errore(err.(error))
		this.AutoErr(err)
	}
}
func (this WsResultProvider) AutoErr(err interface{}) {
	switch e := err.(type) {
	case results.NotFound, *results.NotFound:
		this.NotFound(err)
		break
	case results.ValidationError, *results.ValidationError:
		this.Validation(err)
		break
	case results.Error, *results.Error:
		this.ServerError(err)
		break
	case error:
		this.ServerError(e.Error())
		break
	default:
		this.ServerError(err)
		break
	}
}