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
	"github.com/saichler/l8types/go/sec"
	"github.com/saichler/l8types/go/types/l8sysconfig"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
)

var globals ifs.IResources
var Log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func init() {
	_log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	_log.SetLogLevel(ifs.Trace_Level)
	_resources := resources.NewResources(_log)
	_resources.Set(registry.NewRegistry())
	_security, err := sec.LoadSecurityProvider(nil)
	if err != nil {
		panic("Failed to load security provider " + err.Error())
	}
	_resources.Set(_security)
	_config := &l8sysconfig.L8SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		VnetPort:    50000}
	_resources.Set(_config)
	globals = _resources

}
