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

func TestQueueIsEmpty(t *testing.T) {
	q := queues.NewQueue("test-queue", 100)

	// Should be empty initially
	if !q.IsEmpty() {
		t.Error("Queue should be empty initially")
	}

	// Add an item
	q.Add("test-item")

	// Should not be empty
	if q.IsEmpty() {
		t.Error("Queue should not be empty after adding item")
	}

	// Remove item
	q.Next()

	// Should be empty again
	if !q.IsEmpty() {
		t.Error("Queue should be empty after removing all items")
	}
}
