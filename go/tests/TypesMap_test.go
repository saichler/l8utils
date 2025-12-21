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
	"reflect"
	"testing"

	"github.com/saichler/l8utils/go/utils/registry"
)

type TestStruct struct {
	Name string
}

func TestTypesMapDelAndContains(t *testing.T) {
	tm := registry.NewTypesMap()

	// Test Contains on empty map
	if tm.Contains("TestStruct") {
		t.Error("Should not contain TestStruct initially")
	}

	// Add a value
	testType := reflect.TypeOf(TestStruct{})
	tm.Put("TestStruct", testType)

	// Test Contains
	if !tm.Contains("TestStruct") {
		t.Error("Should contain TestStruct after adding")
	}

	// Test Del
	tm.Del("TestStruct")

	// Test Contains after delete
	if tm.Contains("TestStruct") {
		t.Error("Should not contain TestStruct after deletion")
	}
}
