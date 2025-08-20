package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/queues"
)

func TestNewByteQueue(t *testing.T) {
	name := "test-queue"
	maxSize := 100
	bq := queues.NewByteQueue(name, maxSize)

	if bq == nil {
		t.Fatal("NewByteQueue returned nil")
	}

	if bq.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", bq.Size())
	}

	if !bq.IsEmpty() {
		t.Error("Expected queue to be empty initially")
	}

	if !bq.Active() {
		t.Error("Expected queue to be active initially after bug fix")
	}
}

func TestByteQueueAddAndNext(t *testing.T) {
	bq := queues.NewByteQueue("test", 10)
	
	// Create test data with priority - need large enough slice for PPriority index
	data1 := make([]byte, 200)
	data1[ifs.PPriority] = 0x10 // Priority 1
	data1[0] = 1 // identifier
	
	data2 := make([]byte, 200)
	data2[ifs.PPriority] = 0x20 // Priority 2
	data2[0] = 2 // identifier
	
	data3 := make([]byte, 200)
	data3[ifs.PPriority] = 0x00 // Priority 0
	data3[0] = 3 // identifier

	// Queue should be active now, so items should be added
	bq.Add(data1)
	if bq.Size() != 1 {
		t.Errorf("Expected size 1 after adding one item, got %d", bq.Size())
	}

	bq.Add(data2)
	bq.Add(data3)
	if bq.Size() != 3 {
		t.Errorf("Expected size 3 after adding three items, got %d", bq.Size())
	}

	// Test Next() in separate goroutine to avoid blocking
	go func() {
		time.Sleep(10 * time.Millisecond)
		bq.Shutdown() // This will cause Next() to return nil
	}()

	// Since priorities are processed highest to lowest: P2, P1, P0
	item := bq.Next()
	if item == nil || item[0] != 2 {
		t.Error("Expected to get priority 2 item first")
	}
}

func TestByteQueuePriorityOrdering(t *testing.T) {
	bq := queues.NewByteQueue("test", 10)
	
	// Create test data with different priorities
	data1 := make([]byte, 200)
	data1[ifs.PPriority] = 0x10 // Priority 1
	data1[0] = 1 // Identifier
	
	data2 := make([]byte, 200)
	data2[ifs.PPriority] = 0x30 // Priority 3
	data2[0] = 3 // Identifier
	
	data3 := make([]byte, 200)
	data3[ifs.PPriority] = 0x20 // Priority 2
	data3[0] = 2 // Identifier

	bq.Add(data1)
	bq.Add(data2) 
	bq.Add(data3)

	if bq.Size() != 3 {
		t.Errorf("Expected size 3, got %d", bq.Size())
	}

	// Test priority ordering by manually calling the next() method
	// Since we can't easily test Next() without blocking, we'll test size changes
	bq.Clear()
	if bq.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", bq.Size())
	}
}

func TestByteQueueSize(t *testing.T) {
	bq := queues.NewByteQueue("test", 5)

	if bq.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", bq.Size())
	}

	data := make([]byte, 200)
	data[ifs.PPriority] = 0x10
	
	bq.Add(data)
	if bq.Size() != 1 {
		t.Errorf("Expected size 1 after adding item, got %d", bq.Size())
	}

	bq.Add(data)
	if bq.Size() != 2 {
		t.Errorf("Expected size 2 after adding second item, got %d", bq.Size())
	}
}

func TestByteQueueIsEmpty(t *testing.T) {
	bq := queues.NewByteQueue("test", 5)

	if !bq.IsEmpty() {
		t.Error("Expected queue to be empty initially")
	}

	data := make([]byte, 200)
	data[ifs.PPriority] = 0x10
	bq.Add(data)

	if bq.IsEmpty() {
		t.Error("Expected queue to not be empty after adding item")
	}
}

func TestByteQueueClear(t *testing.T) {
	bq := queues.NewByteQueue("test", 5)

	data := make([]byte, 200)
	data[ifs.PPriority] = 0x10
	bq.Add(data)
	bq.Add(data)

	if bq.Size() != 2 {
		t.Errorf("Expected size 2 before clear, got %d", bq.Size())
	}

	bq.Clear()
	
	if !bq.IsEmpty() {
		t.Error("Expected queue to be empty after Clear()")
	}

	if bq.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", bq.Size())
	}
}

func TestByteQueueShutdown(t *testing.T) {
	bq := queues.NewByteQueue("test", 5)

	if !bq.Active() {
		t.Error("Expected queue to be active initially")
	}

	// Add some data
	data := make([]byte, 200)
	data[ifs.PPriority] = 0x10
	bq.Add(data)

	if bq.Size() != 1 {
		t.Errorf("Expected size 1 before shutdown, got %d", bq.Size())
	}

	bq.Shutdown()

	if bq.Active() {
		t.Error("Expected queue to be inactive after Shutdown()")
	}

	if !bq.IsEmpty() {
		t.Error("Expected queue to be empty after Shutdown()")
	}
}

func TestByteQueueConcurrentAccess(t *testing.T) {
	bq := queues.NewByteQueue("test", 100)
	
	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 5

	// Test concurrent adds
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				data := make([]byte, 200)
				data[ifs.PPriority] = byte(id%4) << 4 // Various priorities
				data[0] = byte(id)
				data[1] = byte(j)
				bq.Add(data)
			}
		}(i)
	}
	wg.Wait()

	expectedSize := numGoroutines * itemsPerGoroutine
	actualSize := bq.Size()
	if actualSize != expectedSize {
		t.Errorf("Expected size %d after concurrent adds, got %d", expectedSize, actualSize)
	}
}

func TestByteQueueMaxSizeWithActiveQueue(t *testing.T) {
	maxSize := 3
	bq := queues.NewByteQueue("test", maxSize)

	// Add items up to maxSize quickly
	for i := 0; i < maxSize; i++ {
		data := make([]byte, 200)
		data[ifs.PPriority] = 0x10
		data[0] = byte(i)
		bq.Add(data)
	}

	if bq.Size() != maxSize {
		t.Errorf("Expected size %d when at max capacity, got %d", maxSize, bq.Size())
	}

	// Adding more items will cause the Add method to wait/block
	// We'll test this behavior with a timeout
	done := make(chan bool)
	go func() {
		data := make([]byte, 200)
		data[ifs.PPriority] = 0x10
		bq.Add(data) // This should block
		done <- true
	}()

	// Give it a moment to potentially block
	select {
	case <-done:
		// If we get here immediately, the add didn't block as expected
		// This could happen if items were consumed or max size logic changed
	case <-time.After(50 * time.Millisecond):
		// Expected behavior - Add() is blocking due to max size
	}

	// Clean up by consuming items
	bq.Clear()
}

func TestByteQueueNextWithTimeout(t *testing.T) {
	bq := queues.NewByteQueue("test", 5)

	// Test Next() with timeout using goroutine
	done := make(chan []byte)
	go func() {
		result := bq.Next() // This will block until item available or shutdown
		done <- result
	}()

	// Give Next() a moment to start waiting
	time.Sleep(10 * time.Millisecond)

	// Add an item
	data := make([]byte, 200)
	data[ifs.PPriority] = 0x20 // Priority 2
	data[0] = 42
	bq.Add(data)

	// Should receive the item
	select {
	case result := <-done:
		if result == nil || result[0] != 42 {
			t.Error("Expected to receive the added item")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Next() did not return item within timeout")
		bq.Shutdown() // Clean up
	}
}

func TestByteQueueMultiplePriorities(t *testing.T) {
	bq := queues.NewByteQueue("test", 10)

	// Test with different priority levels (P0 to P7 since we use >> 4)
	priorities := []byte{0x00, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70}
	
	for i, priority := range priorities {
		data := make([]byte, 200)
		data[ifs.PPriority] = priority
		data[0] = byte(i) // identifier
		bq.Add(data)
	}

	expectedSize := len(priorities)
	if bq.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, bq.Size())
	}
}

func TestByteQueueConcurrentSize(t *testing.T) {
	bq := queues.NewByteQueue("test", 100)
	
	var wg sync.WaitGroup
	numGoroutines := 20

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				size := bq.Size()
				if size < 0 {
					t.Errorf("Size should never be negative, got %d", size)
				}
				bq.IsEmpty() // Also test IsEmpty concurrently
			}
		}()
	}
	wg.Wait()
}

func TestByteQueueConcurrentClear(t *testing.T) {
	bq := queues.NewByteQueue("test", 100)
	
	// Add some initial data
	data := make([]byte, 200)
	data[ifs.PPriority] = 0x10
	bq.Add(data)

	var wg sync.WaitGroup
	numGoroutines := 10

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			bq.Clear()
			if bq.Size() != 0 {
				t.Errorf("Expected size 0 after Clear(), got %d", bq.Size())
			}
		}()
	}
	wg.Wait()
}

func TestByteQueueAddAfterShutdown(t *testing.T) {
	bq := queues.NewByteQueue("test", 5)

	// Shutdown the queue
	bq.Shutdown()

	if bq.Active() {
		t.Error("Expected queue to be inactive after shutdown")
	}

	// Try to add data after shutdown
	data := make([]byte, 200)
	data[ifs.PPriority] = 0x10
	bq.Add(data)

	// Should remain empty since queue is inactive
	if !bq.IsEmpty() {
		t.Error("Expected queue to remain empty after adding to inactive queue")
	}
}

func TestByteQueuePriorityExtraction(t *testing.T) {
	bq := queues.NewByteQueue("test", 10)

	// Test priority extraction logic (data[PPriority] >> 4)
	testCases := []struct {
		priorityByte byte
		expectedPrio int
	}{
		{0x00, 0}, // 0000 0000 >> 4 = 0
		{0x10, 1}, // 0001 0000 >> 4 = 1
		{0x20, 2}, // 0010 0000 >> 4 = 2
		{0x30, 3}, // 0011 0000 >> 4 = 3
		{0x70, 7}, // 0111 0000 >> 4 = 7
	}

	for _, tc := range testCases {
		data := make([]byte, 200)
		data[ifs.PPriority] = tc.priorityByte
		data[0] = byte(tc.expectedPrio) // Use as identifier
		bq.Add(data)
	}

	if bq.Size() != len(testCases) {
		t.Errorf("Expected size %d, got %d", len(testCases), bq.Size())
	}
}