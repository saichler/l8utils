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

package queues

import (
	"math/bits"
	"sync"

	"github.com/saichler/l8types/go/ifs"
)

type ByteQueue struct {
	name         string
	queues       [][][]byte
	priorityMask uint8
	rwMtx        *sync.RWMutex
	cond         *sync.Cond
	maxSize      int
	active       bool
	size         int
}

func NewByteQueue(name string, maxSize int) *ByteQueue {
	bq := &ByteQueue{
		name:         name,
		maxSize:      maxSize,
		active:       true,
		priorityMask: 0,
		size:         0,
	}
	bq.rwMtx = &sync.RWMutex{}
	bq.cond = sync.NewCond(bq.rwMtx)
	bq.queues = make([][][]byte, ifs.P1+1)
	for i := range bq.queues {
		bq.queues[i] = make([][]byte, 0, 16)
	}
	return bq
}

func (this *ByteQueue) Add(data []byte) {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	
	// Wait if queue is full using proper condition variable
	for this.size >= this.maxSize && this.active {
		this.cond.Wait()
	}
	
	if !this.active {
		return
	}
	
	priority := data[ifs.PPriority] >> 4
	if priority > 7 {
		priority = 7 // Cap at maximum priority
	}
	
	this.queues[priority] = append(this.queues[priority], data)
	this.priorityMask |= (1 << priority) // Set bit for this priority
	this.size++
	
	this.cond.Broadcast()
}

func (this *ByteQueue) Next() []byte {
	for this.active {
		var item []byte
		this.rwMtx.Lock()
		if this.size == 0 {
			this.cond.Wait()
		} else {
			item = this.next()
		}
		this.rwMtx.Unlock()
		if item != nil {
			return item
		}
	}
	return nil
}

func (this *ByteQueue) Active() bool {
	return this.active
}

func (this *ByteQueue) Shutdown() {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	
	this.active = false
	this.clear()
	this.cond.Broadcast()
}

func (this *ByteQueue) next() []byte {
	if this.priorityMask == 0 {
		return nil // No items in any queue
	}
	
	// Find highest priority with items - O(1) bit operation
	priority := 7 - bits.LeadingZeros8(this.priorityMask)
	
	// Dequeue from highest priority - O(1)
	queue := &this.queues[priority]
	item := (*queue)[0]
	*queue = (*queue)[1:]
	this.size--
	
	// Clear bit if this priority queue becomes empty - O(1)
	if len(*queue) == 0 {
		this.priorityMask &^= (1 << priority)
	}
	
	this.cond.Broadcast() // Signal waiting producers
	return item
}

func (this *ByteQueue) clear() {
	for i := range this.queues {
		this.queues[i] = this.queues[i][:0] // Keep capacity, reset length
	}
	this.priorityMask = 0
	this.size = 0
}

func (this *ByteQueue) Size() int {
	this.rwMtx.RLock()
	defer this.rwMtx.RUnlock()
	return this.size
}

func (this *ByteQueue) Clear() {
	this.rwMtx.Lock()
	defer this.rwMtx.Unlock()
	this.clear()
}

func (this *ByteQueue) IsEmpty() bool {
	this.rwMtx.RLock()
	defer this.rwMtx.RUnlock()
	return this.size == 0
}
