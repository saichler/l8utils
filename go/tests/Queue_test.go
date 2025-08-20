package tests

import (
	"testing"

	"github.com/saichler/l8utils/go/utils/queues"
)

func TestQueue(t *testing.T) {
	q := queues.NewQueue("test", 3)
	go addToQueue(q)
	popFromQueue(q, t)
	q.Add("g")
	if q.Size() != 1 {
		Log.Fail(t, "Expected queue size to be 1")
		return
	}
	q.Clear()
	if q.Size() != 0 {
		Log.Fail(t, "Expected queue size to be 0")
		return
	}
	if q.Active() {
		q.Shutdown()
	}
	q.Add("s")
	s := q.Next()
	if s != nil {
		Log.Fail(t, "Expected nil")
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
		Log.Fail(t, "Expected queue size to be 3 per the limit")
		return
	}

	for i := 0; i < 5; i++ {
		nxt := q.Next()
		Log.Debug(nxt)
	}
}
