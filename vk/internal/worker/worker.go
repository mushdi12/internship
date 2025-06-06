package worker

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

type Worker struct {
	id   int
	task func(id, msg string)

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorker(id int, task func(id, msg string), parentCtx context.Context) *Worker {
	ctx, cancel := context.WithCancel(parentCtx)
	return &Worker{id: id, task: task, ctx: ctx, cancel: cancel}
}

func (w *Worker) DoWork(input <-chan string, wg *sync.WaitGroup) {
	for {
		select {
		case <-w.ctx.Done():
			wg.Done()
			return
		case letter := <-input:
			w.task(strconv.Itoa(w.id), letter)
			wg.Done()
		}
	}
}

func (w *Worker) StopWork() {
	w.cancel()
}

type WorkerPool struct {
	workers map[int]*Worker
	input   chan string
	running bool // false
	lastID  int
	task    func(id, msg string)

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorkerPool(task func(id, msg string)) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	workers := make(map[int]*Worker)
	return &WorkerPool{
		workers: workers,
		task:    task,
		input:   make(chan string),
		ctx:     ctx,
		cancel:  cancel}
}

func (wp *WorkerPool) addWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	w := NewWorker(wp.lastID, wp.task, wp.ctx)

	wp.workers[w.id] = w
	wp.lastID++

	if wp.running {
		go w.DoWork(wp.input, &wp.wg)
	}
}

func (wp *WorkerPool) AddWorkers(count int) {
	for i := 0; i <= count; i++ {
		wp.addWorker()
	}
}

func (wp *WorkerPool) Remove(id int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	w, exist := wp.workers[id]
	if !exist {
		return
	}

	w.cancel()
	delete(wp.workers, id)
}

func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.running = true

	for _, worker := range wp.workers {
		go worker.DoWork(wp.input, &wp.wg)
	}
}

func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		worker.StopWork()
	}
	close(wp.input)
}

func (wp *WorkerPool) WaitAndStop() {
	wp.wg.Wait()
	fmt.Println("Worker pool stopped")
	wp.Stop()
}

func (wp *WorkerPool) Submit(data []string) {
	for _, d := range data {
		wp.wg.Add(1)
		wp.input <- d
	}
}
