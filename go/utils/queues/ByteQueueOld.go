package queues

import (
	"sync"
	"time"

	"github.com/saichler/l8types/go/ifs"
)

// ByteQueueOld represents the original implementation before optimizations
// This version has:
// - O(P) Next() operation (scans all priority levels)
// - Race condition in Add() method with double locking
// - Polling with 100ms sleep delays
// - Inefficient size calculation (before cached size was added)
type ByteQueueOld struct {
	name    string
	queues  [][][]byte
	rwMtx   *sync.RWMutex
	cond    *sync.Cond
	maxSize int
	active  bool
	size    int
}

func NewByteQueueOld(name string, maxSize int) *ByteQueueOld {
	bq := &ByteQueueOld{}
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

// Add has the original race condition with double locking
func (this *ByteQueueOld) Add(data []byte) {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	
	// RACE CONDITION: This releases and reacquires the lock
	if this.size >= this.maxSize && this.active {
		this.rwMtx.Unlock()
		for this.Size() >= this.maxSize && this.active {
			this.cond.Broadcast()
			time.Sleep(time.Millisecond * 100) // 100ms polling delay
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

func (this *ByteQueueOld) Next() []byte {
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

func (this *ByteQueueOld) Active() bool {
	return this.active
}

func (this *ByteQueueOld) Shutdown() {
	this.active = false
	this.Clear()
	for i := 0; i < this.Size(); i++ {
		this.cond.Broadcast()
	}
}

// Original O(P) next implementation - scans all priority levels
func (this *ByteQueueOld) next() []byte {
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

// Original size calculation (before cached size optimization)
func (this *ByteQueueOld) sizeCalculation() int {
	sum := 0
	for _, b := range this.queues {
		sum += len(b)
	}
	return sum
}

func (this *ByteQueueOld) clear() {
	for i := range this.queues {
		this.queues[i] = [][]byte{}
	}
	this.size = 0
}

func (this *ByteQueueOld) Size() int {
	this.rwMtx.RLock()
	defer this.rwMtx.RUnlock()
	return this.size
}

func (this *ByteQueueOld) Clear() {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	this.clear()
}

func (this *ByteQueueOld) IsEmpty() bool {
	this.rwMtx.RLock()
	defer this.rwMtx.RUnlock()
	return this.size == 0
}

// Performance characteristics of this old implementation:
//
// Time Complexity:
// - Add(): O(1) best case, O(∞) worst case (due to polling)
// - Next(): O(P) where P = number of priority levels (8)
// - Size(): O(1) (after size caching was added)
// - Clear(): O(P)
// - IsEmpty(): O(1)
//
// Issues:
// 1. Race condition in Add() method (lines 43-49)
// 2. 100ms polling delays causing poor responsiveness
// 3. O(P) scanning in next() method
// 4. Thread safety violations
// 5. Inefficient priority queue scanning
//
// This implementation is preserved for comparison and reference purposes.