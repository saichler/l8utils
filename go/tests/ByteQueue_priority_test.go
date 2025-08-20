package tests

import (
	"testing"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/queues"
)

func TestPriorityMaskFunctionality(t *testing.T) {
	bq := queues.NewByteQueue("priority-test", 20)

	// Test empty queue - priority mask should be 0
	if !bq.IsEmpty() {
		t.Error("Queue should be empty initially")
	}

	// Add items with different priorities
	priorities := []byte{0x10, 0x30, 0x50, 0x20} // Priorities 1, 3, 5, 2
	expectedOrder := []byte{5, 3, 2, 1}          // Should come out highest first

	for i, priority := range priorities {
		data := make([]byte, 200)
		data[ifs.PPriority] = priority
		data[0] = priority >> 4 // Store priority in first byte for verification
		bq.Add(data)
		
		if bq.Size() != i+1 {
			t.Errorf("Expected size %d after adding item %d, got %d", i+1, i, bq.Size())
		}
	}

	// Verify items come out in priority order (highest first)
	for i, expectedPriority := range expectedOrder {
		item := bq.Next()
		if item == nil {
			t.Fatalf("Expected item %d, got nil", i)
		}
		
		actualPriority := item[0]
		if actualPriority != expectedPriority {
			t.Errorf("Item %d: expected priority %d, got %d", i, expectedPriority, actualPriority)
		}
		
		expectedSize := len(expectedOrder) - i - 1
		if bq.Size() != expectedSize {
			t.Errorf("After removing item %d: expected size %d, got %d", i, expectedSize, bq.Size())
		}
	}

	// Queue should be empty now
	if !bq.IsEmpty() {
		t.Error("Queue should be empty after removing all items")
	}
}

func TestPriorityOverflow(t *testing.T) {
	bq := queues.NewByteQueue("overflow-test", 10)

	// Test priority capping at 7
	data := make([]byte, 200)
	data[ifs.PPriority] = 0xFF // Priority would be 15, should be capped at 7
	data[0] = 99               // Identifier

	bq.Add(data)

	// Add lower priority item
	data2 := make([]byte, 200)
	data2[ifs.PPriority] = 0x60 // Priority 6
	data2[0] = 66               // Identifier

	bq.Add(data2)

	// First item should be the capped priority (7), which is higher than 6
	item1 := bq.Next()
	if item1 == nil || item1[0] != 99 {
		t.Error("Expected overflow priority item to come first")
	}

	item2 := bq.Next()
	if item2 == nil || item2[0] != 66 {
		t.Error("Expected priority 6 item to come second")
	}
}

func TestSamePriorityFIFO(t *testing.T) {
	bq := queues.NewByteQueue("fifo-test", 10)

	// Add multiple items with same priority
	priority := byte(0x30) // Priority 3
	for i := 0; i < 5; i++ {
		data := make([]byte, 200)
		data[ifs.PPriority] = priority
		data[0] = byte(i) // Use as sequence identifier
		bq.Add(data)
	}

	// Should come out in FIFO order within same priority
	for i := 0; i < 5; i++ {
		item := bq.Next()
		if item == nil {
			t.Fatalf("Expected item %d, got nil", i)
		}
		
		if item[0] != byte(i) {
			t.Errorf("Expected sequence %d, got %d", i, item[0])
		}
	}
}

func TestMixedPriorityAndFIFO(t *testing.T) {
	bq := queues.NewByteQueue("mixed-test", 20)

	// Add items: P1(a), P3(x), P1(b), P3(y), P2(m)
	items := []struct {
		priority byte
		id       byte
	}{
		{0x10, 'a'}, // P1
		{0x30, 'x'}, // P3
		{0x10, 'b'}, // P1
		{0x30, 'y'}, // P3
		{0x20, 'm'}, // P2
	}

	for _, item := range items {
		data := make([]byte, 200)
		data[ifs.PPriority] = item.priority
		data[0] = item.id
		bq.Add(data)
	}

	// Expected order: P3(x), P3(y), P2(m), P1(a), P1(b)
	expectedOrder := []byte{'x', 'y', 'm', 'a', 'b'}

	for i, expected := range expectedOrder {
		item := bq.Next()
		if item == nil {
			t.Fatalf("Expected item %d, got nil", i)
		}
		
		if item[0] != expected {
			t.Errorf("Position %d: expected %c, got %c", i, expected, item[0])
		}
	}
}

func TestClearResetsPriorityMask(t *testing.T) {
	bq := queues.NewByteQueue("clear-test", 10)

	// Add items with various priorities
	priorities := []byte{0x10, 0x30, 0x50}
	for _, priority := range priorities {
		data := make([]byte, 200)
		data[ifs.PPriority] = priority
		bq.Add(data)
	}

	if bq.Size() != 3 {
		t.Errorf("Expected size 3, got %d", bq.Size())
	}

	// Clear should reset everything
	bq.Clear()

	if !bq.IsEmpty() {
		t.Error("Queue should be empty after Clear()")
	}

	if bq.Size() != 0 {
		t.Errorf("Expected size 0 after Clear(), got %d", bq.Size())
	}

	// Should be able to add and retrieve normally after clear
	data := make([]byte, 200)
	data[ifs.PPriority] = 0x40
	data[0] = 42
	bq.Add(data)

	item := bq.Next()
	if item == nil || item[0] != 42 {
		t.Error("Failed to add and retrieve after Clear()")
	}
}