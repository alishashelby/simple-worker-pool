package service

import (
	"errors"
	"github.com/alishashelby/simple-worker-pool/internal/worker"
	"log"
	"sync"
)

const (
	bufferedChanSize = 100
)

type WorkerPool struct {
	work      chan string
	workers   map[int64]*worker.Worker
	wg        *sync.WaitGroup
	mu        *sync.Mutex
	isRunning bool
	workerID  int64
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		work:    make(chan string, bufferedChanSize),
		workers: make(map[int64]*worker.Worker),
		wg:      &sync.WaitGroup{},
		mu:      &sync.Mutex{},
	}
}

func (pool *WorkerPool) checkRunning() error {
	if !pool.isRunning {
		return errors.New("worker-pool is not running")
	}

	return nil
}

func (pool *WorkerPool) Start() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.isRunning {
		return
	}
	pool.isRunning = true
	log.Println("starting worker-pool...")
}

func (pool *WorkerPool) Stop() {
	pool.mu.Lock()

	if err := pool.checkRunning(); err != nil {
		pool.mu.Unlock()
		log.Println(err)
		return
	}

	pool.isRunning = false
	close(pool.work)

	workers := pool.workers
	pool.workers = nil
	pool.mu.Unlock()

	for _, poolWorker := range workers {
		poolWorker.Finish()
	}

	pool.wg.Wait()
	log.Println("stopped entire worker-pool")
}

func (pool *WorkerPool) AddWork(work string) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if err := pool.checkRunning(); err != nil {
		log.Println(err)
		return
	}

	select {
	case pool.work <- work:
	default:
		log.Println("worker-pool is full")
	}
}

func (pool *WorkerPool) AddWorker() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if err := pool.checkRunning(); err != nil {
		log.Println(err)
		return
	}

	id := pool.workerID
	pool.workerID++

	poolWorker := worker.NewWorker(id, pool.work, pool.wg)
	poolWorker.DoWork()
	pool.workers[id] = poolWorker

	log.Printf("added worker: %d\n", id)
}

func (pool *WorkerPool) RemoveWorker() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if err := pool.checkRunning(); err != nil {
		log.Println(err)
		return
	}

	if len(pool.workers) == 0 {
		log.Println("worker-pool is empty")
		return
	}

	var id int64
	var poolWorker *worker.Worker
	for i, w := range pool.workers {
		id = i
		poolWorker = w
		break
	}

	delete(pool.workers, id)
	go poolWorker.Finish()
	log.Printf("removed worker: %d\n", id)
}
