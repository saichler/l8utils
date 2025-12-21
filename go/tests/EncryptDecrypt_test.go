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
	"github.com/saichler/l8types/go/aes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	base := "Test Data To Encrypt"
	key := aes.GenerateAES256Key()
	encData, err := aes.Encrypt([]byte(base), key)
	if err != nil {
		Log.Fail("Failed to encrypt data:", err)
		return
	}
	decData, err := aes.Decrypt(encData, key)
	if err != nil {
		Log.Fail("Failed to decrypt data:", err)
		return
	}
	out := string(decData)
	if out != base {
		Log.Fail("Decrypted data is not equal to base:", out)
		return
	}
}
