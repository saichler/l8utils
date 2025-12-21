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
	"errors"
	"strings"
	"testing"
)

func TestFmtLogger(t *testing.T) {
	err := errors.New("Sample Error")
	Log.Trace("my trace message: ", err)
	Log.Debug("my debug message: ", err)
	Log.Info("my info message: ", err)
	Log.Warning("my warning message: ", err)
	err = Log.Error("my error message: ", err)
	if !strings.Contains(err.Error(), "(Error) - my error message: Sample Error") {
		t.Fail()
		Log.Error("Expected a formatted error message:", err.Error())
		return
	}
	Log.Empty()

	tt := &testing.T{}
	Log.Fail(tt, "my fail message", err)
}
