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

// TestTypesMapPut_PreservesExistingInfo asserts that re-Put of the same type
// name preserves the existing Info instance — and anything attached to it —
// instead of overwriting with a fresh, empty one. This is the contract that
// the Registry.RegisterType doc-comment promises and that the introspector
// relies on when it walks already-registered nested types.
func TestTypesMapPut_PreservesExistingInfo(t *testing.T) {
	tm := registry.NewTypesMap()
	testType := reflect.TypeOf(TestStruct{})

	newReg, err := tm.Put("TestStruct", testType)
	if err != nil {
		t.Fatal("Initial Put failed:", err)
	}
	if !newReg {
		t.Fatal("Initial Put should report a new registration")
	}

	info1, ok := tm.Get("TestStruct")
	if !ok || info1 == nil {
		t.Fatal("Initial Info missing after Put")
	}

	// Re-Put the same type. Must report not-new and preserve the same Info.
	newReg, err = tm.Put("TestStruct", testType)
	if err != nil {
		t.Fatal("Re-Put returned error:", err)
	}
	if newReg {
		t.Fatal("Re-Put should report not-new (false) when key already exists")
	}

	info2, ok := tm.Get("TestStruct")
	if !ok || info2 == nil {
		t.Fatal("Info missing after re-Put")
	}
	if info1 != info2 {
		t.Fatal("Re-Put must preserve the existing *Info pointer; got a fresh Info instead")
	}
}
