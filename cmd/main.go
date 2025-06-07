package main

import (
	"github.com/alishashelby/simple-worker-pool/internal/service"
	"time"
)

func main() {
	pool := service.NewWorkerPool()
	pool.Start()

	input := []string{"brown0", "green1", "red2", "blue3", "yellow4", "auburn5",
		"purple6", "cyan7", "maroon8", "black9", "white10", "grey11",
		"olive12", "coral13", "beige14"}

	go func() {
		for _, str := range input {
			pool.AddWork(str)
		}
	}()

	pool.AddWorker()
	pool.AddWorker()
	pool.AddWorker()
	time.Sleep(10 * time.Millisecond)
	pool.RemoveWorker()
	time.Sleep(10 * time.Millisecond)
	pool.RemoveWorker()
	pool.AddWorker()
	time.Sleep(10 * time.Millisecond)
	pool.RemoveWorker()

	time.Sleep(15 * time.Millisecond)
	pool.Stop()
}
