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
	mtx *sync.Cond
	// Maximum size for the queue, in which the queue will block the input go routine
	maxSize int
	// Is the queue active, e.g. shutdown was not called
	active bool
}

// NewQueue Constructs a new queue
func NewQueue(queueName string, maxSize int) *Queue {
	queue := &Queue{}
	queue.mtx = sync.NewCond(&sync.Mutex{})
	queue.queue = make([]interface{}, 0)
	queue.maxSize = maxSize
	queue.active = true
	queue.queueName = queueName
	return queue
}

// Add an element to the queue and broadcast notification
func (queue *Queue) Add(any interface{}) {
	queue.mtx.L.Lock()
	defer queue.mtx.L.Unlock()

	if len(queue.queue) >= queue.maxSize && queue.active {
		queue.mtx.L.Unlock()
		for len(queue.queue) >= queue.maxSize && queue.active {
			queue.mtx.Broadcast()
			time.Sleep(time.Millisecond * 100)
		}
		queue.mtx.L.Lock()
	}
	if queue.active {
		queue.queue = append(queue.queue, any)
	} else {
		queue.queue = queue.queue[0:0]
	}
	queue.mtx.Broadcast()
}

// Next retrieve the next element in the queue, if the queue is empty this is a blocking queue
func (queue *Queue) Next() interface{} {
	for queue.active {
		var any interface{}
		queue.mtx.L.Lock()
		if len(queue.queue) == 0 {
			queue.mtx.Wait()
		} else {
			any = queue.queue[0]
			queue.queue = queue.queue[1:]

		}
		queue.mtx.L.Unlock()
		if any != nil {
			return any
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
	queue.mtx.Broadcast()
}

// Clear all the content of the queue and return it
func (queue *Queue) Clear() []interface{} {
	queue.mtx.L.Lock()
	defer queue.mtx.L.Unlock()
	result := queue.queue
	queue.queue = make([]interface{}, 0)
	return result
}

// Size of the queue
func (queue *Queue) Size() int {
	queue.mtx.L.Lock()
	defer queue.mtx.L.Unlock()
	return len(queue.queue)
}
