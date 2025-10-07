package tests

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/saichler/l8utils/go/utils/workers"
)

type TestWorker struct {
	counter *int32
	delay   time.Duration
}

func (tw *TestWorker) Run() {
	if tw.delay > 0 {
		time.Sleep(tw.delay)
	}
	atomic.AddInt32(tw.counter, 1)
}

func TestWorkers(t *testing.T) {
	limit := 2
	w := workers.NewWorkers(limit)
	var counter int32

	// Run 5 workers with limit of 2, so they should be throttled
	for i := 0; i < 5; i++ {
		tw := &TestWorker{counter: &counter, delay: 100 * time.Millisecond}
		w.Run(tw)
	}

	// Wait for all workers to complete
	time.Sleep(500 * time.Millisecond)

	if counter != 5 {
		t.Errorf("Expected 5 workers to complete, got %d", counter)
	}
}

type TestTask struct {
	value int
}

func (tt *TestTask) Run() interface{} {
	return tt.value * 2
}

func TestMultiTask(t *testing.T) {
	mt := workers.NewMultiTask()

	// Add multiple tasks
	for i := 0; i < 5; i++ {
		mt.AddTask(&TestTask{value: i})
	}

	// Run all tasks
	results := mt.RunAll()

	// Verify results
	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	for i := 0; i < 5; i++ {
		expected := i * 2
		actual := results[i].(int)
		if actual != expected {
			t.Errorf("Expected result[%d] = %d, got %d", i, expected, actual)
		}
	}
}
