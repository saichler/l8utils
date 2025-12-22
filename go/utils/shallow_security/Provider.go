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

// Package main provides a basic security provider implementation for Layer 8 services.
// ShallowSecurityProvider implements ISecurityProvider with minimal security features
// using AES encryption with an MD5-derived key. This is intended for development and
// testing purposes rather than production use.
//
// Key features:
//   - Basic connection validation with shared secret
//   - AES encryption/decryption for data
//   - Permissive authentication (always allows access)
//   - Placeholder implementations for TFA and registration
package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/saichler/l8types/go/aes"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/nets"
	"github.com/saichler/l8types/go/types/l8sysconfig"

	"net"
	"strconv"
	"strings"
)

// ShallowSecurityProvider implements ISecurityProvider with basic AES encryption.
// Uses a hardcoded secret for key derivation - suitable for testing only.
type ShallowSecurityProvider struct {
	secret string
	key    string
}

// NewShallowSecurityProvider creates a new provider with a hardcoded secret.
func NewShallowSecurityProvider() *ShallowSecurityProvider {
	sp := &ShallowSecurityProvider{}
	hash := md5.New()
	secret := "Shallow Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	sp.key = base64.StdEncoding.EncodeToString(kHash)
	sp.secret = secret
	return sp
}

// CanDial establishes a TCP connection to the specified host and port.
func (this *ShallowSecurityProvider) CanDial(host string, port uint32) (net.Conn, error) {
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
}

// CanAccept always allows incoming connections (permissive).
func (this *ShallowSecurityProvider) CanAccept(conn net.Conn) error {
	return nil
}

// ValidateConnection verifies the connection by exchanging encrypted secrets.
func (this *ShallowSecurityProvider) ValidateConnection(conn net.Conn, config *l8sysconfig.L8SysConfig) error {
	err := nets.WriteEncrypted(conn, []byte(this.secret), config, this)
	if err != nil {
		conn.Close()
		return err
	}

	secret, err := nets.ReadEncrypted(conn, config, this)
	if err != nil {
		conn.Close()
		return err
	}

	if this.secret != secret {
		conn.Close()
		return errors.New("incorrect Secret/Key, aborting connection")
	}

	return nets.ExecuteProtocol(conn, config, this)
}

// Encrypt encrypts data using AES with the derived key.
func (this *ShallowSecurityProvider) Encrypt(data []byte) (string, error) {
	return aes.Encrypt(data, this.key)
}

// Decrypt decrypts AES-encrypted data using the derived key.
func (this *ShallowSecurityProvider) Decrypt(data string) ([]byte, error) {
	return aes.Decrypt(data, this.key)
}

// CanDoAction always permits any action (permissive authorization).
func (this *ShallowSecurityProvider) CanDoAction(action ifs.Action, o ifs.IElements, uuid string, token string, salts ...string) error {
	return nil
}
// ScopeView returns the original elements without filtering (permissive).
func (this *ShallowSecurityProvider) ScopeView(o ifs.IElements, uuid string, token string, salts ...string) ifs.IElements {
	return o
}
// Authenticate always succeeds with a dummy bearer token (testing only).
func (this *ShallowSecurityProvider) Authenticate(user string, pass string) (string, bool, bool, error) {
	return "bearer token", false, false, nil
}
// ValidateToken always validates tokens successfully (permissive).
func (this *ShallowSecurityProvider) ValidateToken(token string) (string, bool) {
	return ifs.NewUuid(), true
}
func (this *ShallowSecurityProvider) Message(string) (*ifs.Message, error) {
	return &ifs.Message{}, nil
}

func (this *ShallowSecurityProvider) TFASetup(userid string, nic ifs.IVNic) (string, []byte, error) {
	return "", nil, nil
}
func (this *ShallowSecurityProvider) TFAVerify(userid string, code string, bearer string, nic ifs.IVNic) error {
	return nil
}
func (this *ShallowSecurityProvider) Captcha() []byte {
	return nil
}
func (this *ShallowSecurityProvider) Register(userId, password, captcha string, vnic ifs.IVNic) error {
	return nil
}

func (this *ShallowSecurityProvider) Credential(crId, cId string, r ifs.IResources) (string, string, string, string, error) {
	return "", "", "", "", nil
}
