package queues

import (
	"sync"
	"time"
)

// Queue is a simple blocking/thread safe "abstract" queue.
// The queue entry is an interface
type Queue struct {
	// The name of the queue for purposes of reporting and loggin
	queueName string
	// The queue itself
	queue []interface{}
	// The cond for waking up go routined
	cond  *sync.Cond
	rwMtx *sync.RWMutex
	// Maximum size for the queue, in which the queue will block the input go routine
	maxSize int
	// Is the queue active, e.g. shutdown was not called
	active bool
}

// NewQueue Constructs a new queue
func NewQueue(queueName string, maxSize int) *Queue {
	queue := &Queue{}
	queue.rwMtx = &sync.RWMutex{}
	queue.cond = sync.NewCond(queue.rwMtx)
	queue.queue = make([]interface{}, 0)
	queue.maxSize = maxSize
	queue.active = true
	queue.queueName = queueName
	return queue
}

// Add an element to the queue and broadcast notification
func (queue *Queue) Add(any interface{}) {
	queue.rwMtx.Lock()
	defer queue.rwMtx.Unlock()
	if len(queue.queue) >= queue.maxSize && queue.active {
		queue.rwMtx.Unlock()
		for len(queue.queue) >= queue.maxSize && queue.active {
			queue.cond.Broadcast()
			time.Sleep(time.Millisecond * 100)
		}
		queue.rwMtx.Lock()
	}
	if queue.active {
		queue.queue = append(queue.queue, any)
	} else {
		queue.queue = queue.queue[0:0]
	}
	queue.cond.Broadcast()
}

// Next retrieve the next element in the queue, if the queue is empty this is a blocking queue
func (queue *Queue) Next() interface{} {
	for queue.active {
		var item interface{}
		queue.rwMtx.Lock()
		if len(queue.queue) == 0 {
			queue.cond.Wait()
		} else {
			item = queue.queue[0]
			queue.queue = queue.queue[1:]

		}
		queue.rwMtx.Unlock()
		if item != nil {
			return item
		}
	}
	return nil
}

// Active is the shutdown was not called
func (queue *Queue) Active() bool {
	return queue.active
}

// Shutdown the queue should unblock and shutdown
func (queue *Queue) Shutdown() {
	queue.active = false
	queue.Clear()
	for i := 0; i < 100; i++ {
		queue.cond.Broadcast()
	}
}

// Clear all the content of the queue and return it
func (queue *Queue) Clear() []interface{} {
	queue.rwMtx.Lock()
	defer queue.rwMtx.Unlock()
	result := queue.queue
	queue.queue = make([]interface{}, 0)
	return result
}

// Size of the queue
func (queue *Queue) Size() int {
	queue.rwMtx.RLock()
	defer queue.rwMtx.RUnlock()
	return len(queue.queue)
}

func (queue *Queue) IsEmpty() bool {
	queue.rwMtx.RLock()
	defer queue.rwMtx.RUnlock()
	return len(queue.queue) == 0
}
