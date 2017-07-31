package client

import (
	"sync"
	"github.com/geruz/bb/transport/protocols/ws"
)


type ChanMap struct {
	chans map[uint64] chan<- ws.Answer
	mutex *sync.Mutex
}

func NewChanMap ()ChanMap {
	return ChanMap{
		chans: map[uint64] chan<- ws.Answer{},
		mutex: &sync.Mutex{},
	}
}
func (this ChanMap)GetAndRemove(id uint64)(chan<- ws.Answer, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ch, ok := this.chans[id]
	delete(this.chans, id)
	return ch, ok
}
func (this ChanMap) Set (id uint64, ch chan<- ws.Answer) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.chans[id] = ch
}

type AnswerChans struct {
	Success       chan []byte
	NotFound      chan []byte
	NotIplemented chan []byte
	ValidateError chan []byte
	Error         chan []byte
}

func NewAnswerChans() AnswerChans {
	return AnswerChans{
		Success:       make(chan []byte),
		Error:         make(chan []byte),
		NotFound:      make(chan []byte),
		NotIplemented: make(chan []byte),
		ValidateError: make(chan []byte),
	}
}

func(this AnswerChans)Close(){
	close(this.Error)
	close(this.Success)
	close(this.NotFound)
	close(this.NotIplemented)
	close(this.ValidateError)
}