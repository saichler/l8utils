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

package tests

import (
	"testing"

	"github.com/saichler/l8utils/go/utils/queues"
)

func TestQueue(t *testing.T) {
	q := queues.NewQueue("test", 3)
	go addToQueue(q)
	popFromQueue(q, t)
	q.Add("g")
	if q.Size() != 1 {
		Log.Fail(t, "Expected queue size to be 1")
		return
	}
	q.Clear()
	if q.Size() != 0 {
		Log.Fail(t, "Expected queue size to be 0")
		return
	}
	if q.Active() {
		q.Shutdown()
	}
	q.Add("s")
	s := q.Next()
	if s != nil {
		Log.Fail(t, "Expected nil")
		return
	}
}

func addToQueue(q *queues.Queue) {
	q.Add("a")
	q.Add("b")
	q.Add("c")
	q.Add("d")
	q.Add("e")
}

func popFromQueue(q *queues.Queue, t *testing.T) {
	for q.Size() < 3 {

	}

	if q.Size() != 3 {
		Log.Fail(t, "Expected queue size to be 3 per the limit")
		return
	}

	for i := 0; i < 5; i++ {
		nxt := q.Next()
		Log.Debug(nxt)
	}
}
