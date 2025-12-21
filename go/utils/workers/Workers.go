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
