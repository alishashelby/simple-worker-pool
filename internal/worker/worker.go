package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Worker struct {
	id     int64
	work   <-chan string
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorker(id int64, work <-chan string, wg *sync.WaitGroup) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		id:     id,
		work:   work,
		wg:     wg,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *Worker) DoWork() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			select {
			case <-w.ctx.Done():
				log.Printf("worker %d finished working\n", w.id)
				return

			case work, ok := <-w.work:
				if !ok {
					log.Printf("worker %d finished working because channel is closed\n", w.id)
					return
				}

				time.Sleep(3 * time.Millisecond)
				fmt.Printf("worker: %d doing work: %s\n", w.id, work)
			}
		}
	}()
}

func (w *Worker) Finish() {
	w.cancel()
}
