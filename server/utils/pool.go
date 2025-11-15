package utils

import (
	"sync"
)

type WorkerPool struct {
	workers int
	jobs    chan func()
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		workers: workers,
		jobs:    make(chan func(), workers),
	}
	pool.startWorkers()
	return pool
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.workers; i++ {
		go func() {
			for job := range p.jobs {
				job()
				p.wg.Done()
			}
		}()
	}
}

func (p *WorkerPool) Do(job func()) {
	p.wg.Add(1)
	p.jobs <- job

}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.jobs)
}
