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
	"strings"
	"testing"
)

func TestShalowSecurity(t *testing.T) {
	sp := globals.Security()
	conn, err := sp.CanDial("127.0.0.1", 8910)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		Log.Fail(t, err)
		return
	}
	err = sp.CanAccept(conn)
	if err != nil {
		Log.Fail(t, err)
		return
	}
	conn = &MockConn{}
	config := globals.SysConfig()
	config.LocalUuid = "Test Validate Connection"

	err = sp.ValidateConnection(conn, config)
	if err != nil {
		Log.Fail(t, err)
		return
	}
	if config.ForceExternal {
		Log.Fail(t, "This connection is adjucent.")
		return
	}
}
