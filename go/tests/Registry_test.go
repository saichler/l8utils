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
	"github.com/saichler/l8types/go/ifs"
	. "github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/registry"
	"reflect"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	ok, err := globals.Registry().Register(nil)
	if err == nil {
		Log.Fail("Expected an error for nil type")
	}

	ok, err = globals.Registry().Register(&TestProto{})
	if !ok || err != nil {
		Log.Fail("Expected to register a proto successfully")
	}

	ok, err = globals.Registry().Register(TestProto{})
	if ok {
		Log.Fail(t, "Type should have already been registered")
		return
	}
	typ, err := globals.Registry().Info(protoName)
	if err != nil {
		Log.Fail(t, "Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		Log.Fail(t, "Wrong type by name")
		return
	}
	_, err = globals.Registry().Info(unknowProtoName)
	if err == nil {
		Log.Fail(t, "Expected an error")
		return
	}
	info, err := globals.Registry().Info(protoName)
	if err != nil {
		Log.Fail(t, "Failed to get type by name", err.Error())
		return
	}
	ins, err := info.NewInstance()
	if err != nil {
		Log.Fail(t, "Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		Log.Fail(t, "Failed to cast instance")
		return
	}
	_, err = globals.Registry().Info(unknowProtoName)
	if err == nil {
		Log.Fail(t, "Expected an error")
		return
	}

	info, err = globals.Registry().Info(protoName)
	if err != nil {
		Log.Fail(t, "Failed to get type by name", err.Error())
		return
	}

	if info.Type() == nil || info.Type().Name() != protoName {
		Log.Fail(t, "Wrong type by name")
		return
	}

	info.AddSerializer(&TestSerializer{})
	ser := info.Serializer(ifs.BINARY)

	if ser == nil {
		Log.Fail(t, "Failed to create serializer")
		return
	}

	if reflect.ValueOf(ser).Elem().Type().Name() != "TestSerializer" {
		Log.Fail(t, "Wrong type by name")
		return
	}

	pb, err := info.NewInstance()
	if err != nil {
		Log.Fail(t, "Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		Log.Fail(t, "Failed to cast protobuf instance")
		return
	}

	i, e := registry.NewInfo(nil)
	defer time.Sleep(time.Second)

	if e == nil {
		Log.Fail(t, "Expected an error")
		return
	}

	if i != nil {
		Log.Fail(t, "Expected nil instance")
		return
	}

	b, e := globals.Registry().RegisterType(nil)
	if e == nil {
		Log.Fail(t, "Expected an error")
		return
	}
	if b {
		Log.Fail(t, "Expected a false")
		return
	}

	b, e = globals.Registry().RegisterType(reflect.ValueOf(reflect.TypeOf(5)).Type().Elem())
	if e != nil {
		Log.Fail(t, "Did not expect an error")
		return
	}
	if b {
		Log.Fail(t, "Expected a false")
		return
	}
}

// stringTestSerializer is a STRING-mode serializer used by
// TestRegisterTypeIdempotent_PreservesSerializer to verify that
// re-registering a type does NOT discard previously-attached serializers.
type stringTestSerializer struct{}

func (s *stringTestSerializer) Mode() ifs.SerializerMode {
	return ifs.STRING
}
func (s *stringTestSerializer) Marshal(interface{}, ifs.IResources) ([]byte, error) {
	return nil, nil
}
func (s *stringTestSerializer) Unmarshal([]byte, ifs.IResources) (interface{}, error) {
	return nil, nil
}

// TestRegisterTypeIdempotent_PreservesSerializer is the registry-level
// regression test for the bug where Registry.RegisterType silently wiped
// any serializer attached to an already-registered type. After this fix,
// the serializer attached after first registration must still be reachable
// after a second RegisterType call for the same type — the call site that
// most often triggers the second registration is the introspector walking
// nested struct fields during service activation.
func TestRegisterTypeIdempotent_PreservesSerializer(t *testing.T) {
	r := registry.NewRegistry()

	if _, err := r.Register(&TestProto{}); err != nil {
		t.Fatal("Initial Register failed:", err)
	}

	info, err := r.Info("TestProto")
	if err != nil || info == nil {
		t.Fatal("Info missing after first Register:", err)
	}
	info.AddSerializer(&stringTestSerializer{})

	// Sanity: serializer reachable before re-register.
	if info.Serializer(ifs.STRING) == nil {
		t.Fatal("STRING serializer should be reachable immediately after AddSerializer")
	}

	// Re-register the same type — simulates what Introspector.addAttribute
	// does when it walks a nested struct field whose type is already in the
	// registry.
	newReg, err := r.RegisterType(reflect.TypeOf(TestProto{}))
	if err != nil {
		t.Fatal("Re-RegisterType returned error:", err)
	}
	if newReg {
		t.Fatal("Re-RegisterType should return false (not-new) for an existing type")
	}

	// The serializer MUST still be reachable.
	info2, err := r.Info("TestProto")
	if err != nil || info2 == nil {
		t.Fatal("Info missing after re-RegisterType:", err)
	}
	if info != info2 {
		t.Fatal("Re-RegisterType must preserve the existing *Info pointer")
	}
	ser := info2.Serializer(ifs.STRING)
	if ser == nil {
		t.Fatal("STRING serializer was wiped by re-RegisterType (root cause of READY=0/0 bug)")
	}
	if _, ok := ser.(*stringTestSerializer); !ok {
		t.Fatal("Serializer instance changed after re-RegisterType")
	}
}
