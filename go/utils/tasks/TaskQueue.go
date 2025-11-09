package tasks

import "sync"

type Tasks struct {
	name       string
	maxWorkers int

	mtx  *sync.RWMutex
	cond *sync.Cond

	queueT1   []ITask
	queueT2   []ITask
	runningT1 map[string]bool
	runningT2 map[string]bool
}

type ITask interface {
	Execute()
	T1Id() string
	T2Id() string
	Timeout() int64
}

func NewTasks(name string, maxWorkers int) *Tasks {
	this := &Tasks{}
	this.name = name
	this.maxWorkers = maxWorkers

	this.mtx = &sync.RWMutex{}
	this.cond = sync.NewCond(this.mtx)

	this.queueT1 = []ITask{}
	this.queueT2 = []ITask{}
	this.runningT1 = map[string]bool{}
	this.runningT2 = map[string]bool{}

	return this
}

func (this *Tasks) AddTask(task ITask) {
	/*
		this.mtx.Lock()
		defer this.mtx.Unlock()
		if
		this.queue = append(this.queue, task)
	*/
}
