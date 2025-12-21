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

package main

import (
	"bytes"
	"os"
)

func SeekResource(path string, filename string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if fileInfo.Name() == filename {
		return path
	}
	if fileInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return ""
		}
		for _, file := range files {
			found := SeekResource(pathOf(path, file), filename)
			if found != "" {
				return found
			}
		}
	}
	return ""
}

func pathOf(path string, file os.DirEntry) string {
	buff := bytes.Buffer{}
	buff.WriteString(path)
	buff.WriteString("/")
	buff.WriteString(file.Name())
	return buff.String()
}
