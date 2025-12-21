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

	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/resources"
)

func TestResources(t *testing.T) {
	// Create logger
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

	// Create resources
	r := resources.NewResources(log)

	// Test logger getter
	if r.Logger() == nil {
		t.Error("Logger should not be nil")
	}

	// Test other getters returning nil initially
	if r.Registry() != nil {
		t.Error("Registry should be nil initially")
	}
	if r.Services() != nil {
		t.Error("Services should be nil initially")
	}
	if r.Security() != nil {
		t.Error("Security should be nil initially")
	}
	if r.DataListener() != nil {
		t.Error("DataListener should be nil initially")
	}
	if r.Introspector() != nil {
		t.Error("Introspector should be nil initially")
	}
	if r.SysConfig() != nil {
		t.Error("SysConfig should be nil initially")
	}
}

func TestResourcesSet(t *testing.T) {
	// Create logger
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

	// Create resources
	r := resources.NewResources(log)

	// Test Set with nil (should do nothing)
	r.Set(nil)

	// Test Set with unknown type (should log error)
	r.Set("unknown type")
}
