// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tasks provides a dual-queue task scheduling system with configurable worker limits.
// Tasks can be categorized into T1 and T2 tiers with independent tracking and execution.
// This is a work-in-progress implementation.
//
// Key features:
//   - Dual task queues (T1 and T2) for priority/categorization
//   - Configurable maximum concurrent workers
//   - Task deduplication via T1Id and T2Id tracking
package tasks

import "sync"

// Tasks manages dual task queues with worker pool coordination.
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

// ITask defines the interface for executable tasks with tier identification.
type ITask interface {
	Execute()
	T1Id() string
	T2Id() string
	Timeout() int64
}

// NewTasks creates a new dual-queue task manager with the specified worker limit.
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

// AddTask queues a task for execution (work in progress).
func (this *Tasks) AddTask(task ITask) {
	/*
		this.mtx.Lock()
		defer this.mtx.Unlock()
		if
		this.queue = append(this.queue, task)
	*/
}
