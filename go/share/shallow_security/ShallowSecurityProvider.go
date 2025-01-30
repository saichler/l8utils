package shallow_security

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/saichler/shared/go/share/aes"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/nets"
	"github.com/saichler/shared/go/types"
	"net"
	"strconv"
	"strings"
)

type ShallowSecurityProvider struct {
	secret string
	key    string
	salts  []string
}

func NewShallowSecurityProvider(key, secret string, salts ...string) *ShallowSecurityProvider {
	sp := &ShallowSecurityProvider{}
	sp.key = key
	sp.secret = secret
	sp.salts = salts
	return sp
}

func (this *ShallowSecurityProvider) CanDial(host string, port uint32) (net.Conn, error) {
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
}

func (this *ShallowSecurityProvider) CanAccept(conn net.Conn) error {
	return nil
}

func (this *ShallowSecurityProvider) ValidateConnection(conn net.Conn, config *types.VNicConfig) error {
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

func (this *ShallowSecurityProvider) Encrypt(data []byte) (string, error) {
	return aes.Encrypt(data, this.key)
}

func (this *ShallowSecurityProvider) Decrypt(data string) ([]byte, error) {
	return aes.Decrypt(data, this.key)
}

func (this *ShallowSecurityProvider) CanDo(action types.Action, endpoint string, token string) error {
	return nil
}
func (this *ShallowSecurityProvider) CanView(typ string, attrName string, token string) error {
	return nil
}

func CreateShallowSecurityProvider() interfaces.ISecurityProvider {
	hash := md5.New()
	secret := "Default Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	k := base64.StdEncoding.EncodeToString(kHash)
	return NewShallowSecurityProvider(k, secret)
}
