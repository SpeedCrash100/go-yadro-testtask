package pkg

import "testing"

func TestQueueFIFOCheck(t *testing.T) {
	count := 100
	q := NewQueue[int](count)
	for i := 0; i < count; i++ {
		if err := q.Push(i); err != nil {
			t.Errorf("Failed to insert in queue: %v", err)
		}
	}

	for i := 0; i < count; i++ {
		v, err := q.Pop()
		if err != nil {
			t.Errorf("Failed to dequeue: %v", err)
		}

		if v != i {
			t.Errorf("Queue is not FIFO")
		}
	}
}

func TestQueueOverflow(t *testing.T) {
	count := 100
	q := NewQueue[int](count)
	for i := 0; i < count; i++ {
		if err := q.Push(i); err != nil {
			t.Errorf("Failed to insert in queue: %v", err)
		}
	}

	if !q.IsFull() {
		t.Errorf("Expected queue to be full")
	}

	if err := q.Push(0); err == nil {
		t.Errorf("Expected overflow error")
	}

}

func TestQueueUnderflow(t *testing.T) {
	count := 100
	q := NewQueue[int](count)

	if _, err := q.Pop(); err == nil {
		t.Errorf("Expected underrun error")
	}

	for i := 0; i < count; i++ {
		if err := q.Push(i); err != nil {
			t.Errorf("Failed to insert in queue: %v", err)
		}
	}

	for i := 0; i < count; i++ {
		v, err := q.Pop()
		if err != nil {
			t.Errorf("Failed to dequeue: %v", err)
		}

		if v != i {
			t.Errorf("Queue is not FIFO")
		}
	}

	if _, err := q.Pop(); err == nil {
		t.Errorf("Expected underrun error")
	}

}
