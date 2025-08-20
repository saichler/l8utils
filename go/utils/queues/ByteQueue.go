package queues

import (
	"sync"
	"time"

	"github.com/saichler/l8types/go/ifs"
)

type ByteQueue struct {
	name    string
	queues  [][][]byte
	rwMtx   *sync.RWMutex
	cond    *sync.Cond
	maxSize int
	active  bool
	size    int
}

func NewByteQueue(name string, maxSize int) *ByteQueue {
	bq := &ByteQueue{}
	bq.active = true
	bq.rwMtx = &sync.RWMutex{}
	bq.cond = sync.NewCond(bq.rwMtx)
	bq.name = name
	bq.maxSize = maxSize
	bq.queues = make([][][]byte, ifs.P1+1)
	for i := range bq.queues {
		bq.queues[i] = [][]byte{}
	}
	return bq
}

func (this *ByteQueue) Add(data []byte) {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	if this.size >= this.maxSize && this.active {
		this.rwMtx.Unlock()
		for this.Size() >= this.maxSize && this.active {
			this.cond.Broadcast()
			time.Sleep(time.Millisecond * 100)
		}
		this.rwMtx.Lock()
	}
	if this.active {
		priority := data[ifs.PPriority] >> 4
		this.queues[priority] = append(this.queues[priority], data)
		this.size++
	} else {
		this.clear()
	}
	this.cond.Broadcast()
}

func (this *ByteQueue) Next() []byte {
	for this.active {
		var item []byte
		this.rwMtx.Lock()
		if this.size == 0 {
			this.cond.Wait()
		} else {
			item = this.next()
		}
		this.rwMtx.Unlock()
		if item != nil {
			return item
		}
	}
	return nil
}

func (this *ByteQueue) Active() bool {
	return this.active
}

func (this *ByteQueue) Shutdown() {
	this.active = false
	this.Clear()
	for i := 0; i < this.Size(); i++ {
		this.cond.Broadcast()
	}
}

func (this *ByteQueue) next() []byte {
	for i := len(this.queues) - 1; i >= 0; i-- {
		if len(this.queues[i]) > 0 {
			data := this.queues[i][0]
			this.queues[i] = this.queues[i][1:]
			this.size--
			return data
		}
	}
	return nil
}

func (this *ByteQueue) clear() {
	for i := range this.queues {
		this.queues[i] = [][]byte{}
	}
	this.size = 0
}

func (this *ByteQueue) Size() int {
	this.rwMtx.RLock()
	defer this.rwMtx.RUnlock()
	return this.size
}

func (this *ByteQueue) Clear() {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	this.clear()
}

func (this *ByteQueue) IsEmpty() bool {
	this.rwMtx.RLock()
	defer this.rwMtx.RUnlock()
	return this.size == 0
}
