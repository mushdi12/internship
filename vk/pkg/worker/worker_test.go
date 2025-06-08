package worker

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool_Processing(t *testing.T) {
	var mu sync.Mutex
	results := make([]string, 0)

	task := func(id, msg string) {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, fmt.Sprintf("worker-%s: %s", id, msg))
	}

	pool := NewWorkerPool(task)
	pool.AddWorkers(2)
	pool.Start()

	input := []string{"A", "B", "C"}
	pool.Submit(input)
	pool.WaitAndStop()

	if len(results) != len(input) {
		t.Errorf("Expected %d results, got %d", len(input), len(results))
	}
}

func TestWorkerPool_RemoveWorker(t *testing.T) {
	var mu sync.Mutex
	calledIDs := make(map[string]bool)

	task := func(id, msg string) {
		mu.Lock()
		defer mu.Unlock()
		calledIDs[id] = true
	}

	pool := NewWorkerPool(task)
	pool.AddWorkers(1) // добавляет 0 и 1
	pool.Start()
	pool.Remove(1)

	pool.Submit([]string{"X"})
	pool.WaitAndStop()

	if calledIDs["1"] {
		t.Errorf("Worker 1 should have been removed")
	}
}

func TestWorkerPool_LargeSubmit(t *testing.T) {
	counter := 0
	mu := sync.Mutex{}

	task := func(id, msg string) {
		mu.Lock()
		counter++
		mu.Unlock()
	}

	pool := NewWorkerPool(task)
	pool.AddWorkers(4)
	pool.Start()

	var data []string
	for i := 0; i < 100; i++ {
		data = append(data, fmt.Sprintf("msg-%d", i))
	}
	pool.Submit(data)
	pool.WaitAndStop()

	if counter != 100 {
		t.Errorf("Expected 100 processed messages, got %d", counter)
	}
}

func TestWorkerPool_ManualStop(t *testing.T) {
	task := func(id, msg string) {
		time.Sleep(10 * time.Millisecond) // имитируем работу
	}

	pool := NewWorkerPool(task)
	pool.AddWorkers(1)
	pool.Start()

	pool.Submit([]string{"1", "2", "3", "4", "5"})
	pool.Stop()
}
