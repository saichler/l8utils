package tests

import (
	"testing"

	"github.com/saichler/l8utils/go/utils/queues"
)

func TestQueueIsEmpty(t *testing.T) {
	q := queues.NewQueue("test-queue", 100)

	// Should be empty initially
	if !q.IsEmpty() {
		t.Error("Queue should be empty initially")
	}

	// Add an item
	q.Add("test-item")

	// Should not be empty
	if q.IsEmpty() {
		t.Error("Queue should not be empty after adding item")
	}

	// Remove item
	q.Next()

	// Should be empty again
	if !q.IsEmpty() {
		t.Error("Queue should be empty after removing all items")
	}
}
