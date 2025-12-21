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
	"github.com/saichler/l8utils/go/utils/maps"
	"reflect"

	"testing"
)

func TestSyncMap(t *testing.T) {
	key := "key"
	val := "val"
	m := maps.NewSyncMap()
	m.Put(key, val)
	v, ok := m.Get(key)
	if !ok {
		Log.Fail(t, "Expected key to exist")
		return
	}
	if v != val {
		Log.Fail(t, "Expected value to be '"+val+"'")
		return
	}

	if m.Size() != 1 {
		Log.Fail(t, "Expected size to be 1")
		return
	}

	m.Clean()
	if m.Size() != 0 {
		Log.Fail(t, "Expected size to be 0")
		return
	}

	m.Put(key, val)

	v, ok = m.Delete(key)
	if !ok {
		Log.Fail(t, "Expected key to exist")
		return
	}
	if v != val {
		Log.Fail(t, "Expected value to be '"+val+"'")
		return
	}

	if m.Contains(key) {
		Log.Fail(t, "Expected key '"+key+" to NOT exist")
	}

	m.Put("a", "b")
	m.Put("c", "d")
	m.Put("e", "f")

	vFilter := func(filter interface{}) bool {
		k := filter.(string)
		if k == "d" {
			return false
		}
		return true
	}

	l := m.ValuesAsList(reflect.ValueOf(val).Type(), vFilter)
	list := l.([]string)

	if len(list) != 2 {
		Log.Fail(t, "Expected length of list to be 2, but it is:", len(list))
		return
	}

	if !m.Contains("a") || !m.Contains("e") {
		Log.Fail(t, "Expected 'a' & 'e' keys to exist")
		return
	}

	l = m.ValuesAsList(reflect.ValueOf(val).Type(), nil)
	list = l.([]string)

	if len(list) != 3 {
		Log.Fail(t, "Expected length of list to be 3, but it is:", len(list))
		return
	}

	if !m.Contains("a") || !m.Contains("e") || !m.Contains("c") {
		Log.Fail(t, "Expected 'a', 'c' & 'e' keys to exist")
		return
	}

	l = m.KeysAsList(reflect.ValueOf(val).Type(), nil)
	list = l.([]string)
	if len(list) != 3 {
		Log.Fail(t, "Expected length of list to be 3, but it is:", len(list))
		return
	}

	kFilter := func(filter interface{}) bool {
		k := filter.(string)
		if k == "c" {
			return false
		}
		return true
	}

	l = m.KeysAsList(reflect.ValueOf(val).Type(), kFilter)
	list = l.([]string)
	if len(list) != 2 {
		Log.Fail(t, "Expected length of list to be 2, but it is:", len(list))
		return
	}

	itf := func(key interface{}, val interface{}) {
		Log.Debug("key:", key, " val:", val)
	}

	m.Iterate(itf)
}
