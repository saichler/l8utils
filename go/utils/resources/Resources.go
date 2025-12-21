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

package resources

import (
	"reflect"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8sysconfig"
)

var DEFAULT_MAX_DATA_SIZE uint64 = 1024 * 1024 * 50
var DEFAULT_QUEUE_SIZE uint64 = 100000

type Resources struct {
	logger       ifs.ILogger
	registry     ifs.IRegistry
	services     ifs.IServices
	security     ifs.ISecurityProvider
	dataListener ifs.IDatatListener
	serializers  map[ifs.SerializerMode]ifs.ISerializer
	config       *l8sysconfig.L8SysConfig
	introspector ifs.IIntrospector
}

func NewResources(logger ifs.ILogger) ifs.IResources {
	r := &Resources{}
	r.logger = logger
	r.serializers = make(map[ifs.SerializerMode]ifs.ISerializer)
	return r
}

func (this *Resources) AddService(serviceName string, serviceArea int32) {
	ifs.AddService(this.config, serviceName, serviceArea)
}

func (this *Resources) Set(any interface{}) {
	if any == nil {
		return
	}
	registry, ok := any.(ifs.IRegistry)
	if ok {
		this.registry = registry
		return
	}

	services, ok := any.(ifs.IServices)
	if ok {
		this.services = services
		return
	}

	security, ok := any.(ifs.ISecurityProvider)
	if ok {
		this.security = security
		return
	}

	dataListener, ok := any.(ifs.IDatatListener)
	if ok {
		this.dataListener = dataListener
		return
	}

	serializer, ok := any.(ifs.ISerializer)
	if ok {
		this.serializers[serializer.Mode()] = serializer
	}

	config, ok := any.(*l8sysconfig.L8SysConfig)
	if ok {
		this.config = config
		return
	}

	introspector, ok := any.(ifs.IIntrospector)
	if ok {
		this.introspector = introspector
		return
	}
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	this.logger.Error("Unknown Set type ", v.Type().Name())
}

func (this *Resources) Copy(other ifs.IResources) {
	this.registry = other.Registry()
	this.security = other.Security()
	this.services = other.Services()
	this.serializers[ifs.BINARY] = other.Serializer(ifs.BINARY)
	this.introspector = other.Introspector()
	this.dataListener = other.DataListener()
	this.config = other.SysConfig()
}

func (this *Resources) Registry() ifs.IRegistry {
	return this.registry
}
func (this *Resources) Services() ifs.IServices {
	return this.services
}
func (this *Resources) Security() ifs.ISecurityProvider {
	return this.security
}
func (this *Resources) DataListener() ifs.IDatatListener {
	return this.dataListener
}
func (this *Resources) Serializer(mode ifs.SerializerMode) ifs.ISerializer {
	return this.serializers[mode]
}
func (this *Resources) Logger() ifs.ILogger {
	return this.logger
}
func (this *Resources) SysConfig() *l8sysconfig.L8SysConfig {
	return this.config
}
func (this *Resources) Introspector() ifs.IIntrospector {
	return this.introspector
}
