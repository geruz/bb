package ws

import (
	"github.com/valyala/fasthttp"
	"net"
	"crypto/sha1"
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"time"
	"fmt"
	"sync"
)

type WsContext struct{
	ReadNext <-chan []byte
	Close func()
	conn net.Conn
}

const(
	ENDED_FRAME = byte(0x80)
	NOT_ENDED_FRAME = byte(0x7F)
	TEXT_FRAME = byte(0x1)
	BIN_FARME = byte(0x2)
	CLOSE_FARME = byte(0x8)
	PING_FRAME = byte(0x9)
	PONG_FRAME = byte(0xA)
	CONTINUE_FRAME = byte(0x0)
	USE_MASK = byte(0x80)
)
var timeout = time.Second * 100
func length(l int) ( []byte) {
	if l <= 125 {
		return []byte{byte(l)}
	}
	if l < 0xFFFF {
		return []byte{byte(126), byte(l / 0xFF), byte(l % 0xFF)}
	}
	bs := make([]byte, 9)
	bs[0] = 127
	binary.LittleEndian.PutUint64(bs[1:], uint64(l))
	return bs

}
//TODO убрать
var ll = sync.Mutex{}
func (this WsContext) WriteType(typ byte, data []byte) {
	ll.Lock()
	defer ll.Unlock()
	firstByte := ENDED_FRAME | typ
	this.conn.Write([]byte{firstByte})
	this.conn.Write(length(len(data)))
	fmt.Println(data)
	this.conn.Write(data)
	time.Sleep(time.Second)
}
func (this WsContext) Write(data []byte) {
	this.WriteType(BIN_FARME, data)
}
func (this WsContext) WriteString(str string){
	data := []byte(str)
	this.WriteType(TEXT_FRAME, data)
}


func ReadBody(br *bufio.Reader, useMask bool, length uint64) ([]byte, error){
	if length == 126{
		p, err := br.Peek(2)
		if err != nil{
			return nil, err
		}
		br.Discard(2)
		length = uint64(binary.BigEndian.Uint16(p))
	}
	if length == 127{
		p, _ := br.Peek(8)
		br.Discard(8)
		length = uint64(binary.BigEndian.Uint64(p))
	}
	var mask []byte
	if useMask {
		var err error
		mask, err = br.Peek(4)
		br.Discard(4)
		if err != nil{
			return nil, err
		}
	}

	data, err := br.Peek(int(length))
	if err != nil{
		return nil, err
	}
	br.Discard(int(length))
	if useMask {
		for i := 0; i < int(length); i++ {
			data[i] = data[i] ^ mask[i%4]
		}
	}
	return data, nil
}
func NewWsContext(c net.Conn) WsContext{
	ch := make(chan []byte)
	clearTimeout := make(chan struct{})
	cl := make(chan struct{})

	br := bufio.NewReaderSize(c, 2048)
	isClose := false
	closeAll:= func(){
		if isClose{
			return
		}
		isClose = true
		c.Close()
		close(ch)
		close(cl)
		close(clearTimeout)
	}

	ws := WsContext{
		conn: c,
		ReadNext: ch,
		Close: closeAll,

	}
	go func () {
		for {
			select {
				case <-time.After(timeout):
					ws.Close()
				case <-clearTimeout:
					continue;
				case <-cl:
					return
			}
		}
	}()

	go func() {
		for {
			select {
				case <-time.After(timeout / 2000):
					ws.WriteType(PING_FRAME, []byte{0x1})
				case <-cl:
					return
			}
		}
	}()

	go func (){
		defer closeAll()
		for {
			arr, err := br.Peek(2)
			if err != nil{
				return
			}
			br.Discard(2)
			firstByte := arr[0]
			//TODO добавить чтение если фрейм не конечный
			typ := firstByte & 0xF
			second := arr[1]
			useMask := second & USE_MASK > 0
			length := uint64(second & (USE_MASK - 1))
			clearTimeout <- struct{}{}
			switch typ {
			case CLOSE_FARME:
				return;
			case PING_FRAME:
				fmt.Println("PING")
				_, err := ReadBody(br, useMask, length)
				if err != nil{
					return
				}
				/*ws.WriteType(PONG_FRAME, data)*/
				continue
			case PONG_FRAME:
				fmt.Println("PONG")
				_, err := ReadBody(br, useMask, length)
				if err != nil{
					return
				}
				continue
			default:
				data, err := ReadBody(br, useMask, length)
				if err != nil {
					return
				}
				ch <- data
			}
		}
	}()
	return ws
}



func Upgrade(handleF func (context WsContext)) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {

		challengeKey := ctx.Request.Header.Peek("Sec-Websocket-Key")
		if len(challengeKey) == 0 {
			ctx.Error("websocket: key missing or blank", fasthttp.StatusBadRequest)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusSwitchingProtocols)
		ctx.Response.Header.Set("Upgrade", "websocket")
		ctx.Response.Header.Set("Connection", "Upgrade")
		ctx.Response.Header.Set("Sec-WebSocket-Accept", computeAcceptKeyByte(challengeKey))
		ctx.Hijack(func(c net.Conn) {
			ws := NewWsContext(c)
			handleF(ws)
		})
	}
}

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
func computeAcceptKeyByte(challengeKey []byte) string {
	h := sha1.New()
	h.Write(challengeKey)
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}