package tests

import (
	"github.com/saichler/shared/go/src/interfaces"
	"github.com/saichler/shared/go/src/queues"
	"testing"
)

func TestQueue(t *testing.T) {
	q := queues.NewQueue("test", 3)
	go addToQueue(q)
	popFromQueue(q, t)
	q.Add("g")
	if q.Size() != 1 {
		interfaces.Fail(t, "Expected queue size to be 1")
		return
	}
	q.Clear()
	if q.Size() != 0 {
		interfaces.Fail(t, "Expected queue size to be 0")
		return
	}
	if q.Active() {
		q.Shutdown()
	}
	q.Add("s")
	s := q.Next()
	if s != nil {
		interfaces.Fail(t, "Expected nil")
		return
	}
}

func addToQueue(q *queues.Queue) {
	q.Add("a")
	q.Add("b")
	q.Add("c")
	q.Add("d")
	q.Add("e")
}

func popFromQueue(q *queues.Queue, t *testing.T) {
	for q.Size() < 3 {

	}

	if q.Size() != 3 {
		interfaces.Fail(t, "Expected queue size to be 3 per the limit")
		return
	}

	for i := 0; i < 5; i++ {
		nxt := q.Next()
		interfaces.Debug(nxt)
	}
}

func TestByteQueue(t *testing.T) {
	bq := queues.NewByteSliceQueue("test", 3)
	bq.Add([]byte{50, 51, 51})
	b := bq.Next()
	if b[0] != 50 || b[1] != 51 || b[2] != 51 {
		interfaces.Fail(t, "Expected byte slice to be 50 51 51")
		return
	}
	bq.Shutdown()
	bq.Active()
}
