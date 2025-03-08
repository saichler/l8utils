package workers

import "sync"

type Workers struct {
	limit   int
	running int
	cond    *sync.Cond
}

type IWorker interface {
	Run()
}

type Worker struct {
	worker  IWorker
	workers *Workers
}

func NewWorkers(limit int) *Workers {
	return &Workers{limit: limit, cond: sync.NewCond(&sync.Mutex{})}
}

func (this *Workers) canStart() {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	for this.running >= this.limit {
		this.cond.Wait()
	}
	this.running++
}

func (this *Workers) Run(worker IWorker) {
	this.canStart()
	w := &Worker{worker: worker, workers: this}
	go w.run()
}

func (this Worker) run() {
	this.worker.Run()
	this.workers.cond.L.Lock()
	defer this.workers.cond.L.Unlock()
	this.workers.running--
	this.workers.cond.Broadcast()
}
