package queues

// ByteSliceQueue wrapper type over queue to make it easy to add and get []byte elements
type ByteSliceQueue struct {
	queue *Queue
}

// NewByteSliceQueue constructs a new byte slice queue
func NewByteSliceQueue(queueName string, maxSize int) *ByteSliceQueue {
	bsQueue := &ByteSliceQueue{}
	bsQueue.queue = NewQueue(queueName, maxSize)
	return bsQueue
}

// Add a []byte array element to the queue
func (bsQueue *ByteSliceQueue) Add(packet []byte) {
	bsQueue.queue.Add(packet)
}

// Next fetch next available []byte element, if empty block
func (bsQueue *ByteSliceQueue) Next() []byte {
	any := bsQueue.queue.Next()
	data, _ := any.([]byte)
	return data
}

// Active is this queue still active or already been shutdown
func (bsQueue *ByteSliceQueue) Active() bool {
	return bsQueue.queue.Active()
}

// Shutdown this queue
func (bsQueue *ByteSliceQueue) Shutdown() {
	bsQueue.queue.Shutdown()
}
