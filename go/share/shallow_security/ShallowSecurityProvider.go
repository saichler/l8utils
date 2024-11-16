package shallow_security

import (
	"errors"
	"github.com/saichler/shared/go/share/aes"
	"github.com/saichler/shared/go/share/nets"
	"github.com/saichler/shared/go/types"
	"net"
	"strconv"
)

type ShallowSecurityProvider struct {
	secret string
	key    string
}

func NewShallowSecurityProvider(key, secret string) *ShallowSecurityProvider {
	sp := &ShallowSecurityProvider{}
	sp.key = key
	sp.secret = secret
	return sp
}

func (sp *ShallowSecurityProvider) CanDial(host string, port uint32, salts ...interface{}) (net.Conn, error) {
	return net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
}

func (sp *ShallowSecurityProvider) CanAccept(conn net.Conn, salts ...interface{}) error {
	return nil
}

func (sp *ShallowSecurityProvider) ValidateConnection(conn net.Conn, uuid string, config *types.MessagingConfig, salts ...interface{}) (string, error) {
	err := nets.WriteEncrypted(conn, []byte(sp.secret), config, salts...)
	if err != nil {
		conn.Close()
		return "", err
	}

	secret, err := nets.ReadEncrypted(conn, config, salts...)
	if err != nil {
		conn.Close()
		return "", err
	}

	if sp.secret != secret {
		conn.Close()
		return "", errors.New("incorrect Secret/Key, aborting connection")
	}

	err = nets.WriteEncrypted(conn, []byte(uuid), config, salts...)
	if err != nil {
		conn.Close()
		return "", err
	}

	zside, err := nets.ReadEncrypted(conn, config, salts...)
	if err != nil {
		conn.Close()
		return "", err
	}

	return zside, nil
}

func (sp *ShallowSecurityProvider) Encrypt(data []byte, salts ...interface{}) (string, error) {
	return aes.Encrypt(data, sp.key)
}

func (sp *ShallowSecurityProvider) Decrypt(data string, salts ...interface{}) ([]byte, error) {
	return aes.Decrypt(data, sp.key)
}

func (sp *ShallowSecurityProvider) CanDo(action types.Action, endpoint string, token string, salts ...interface{}) error {
	return nil
}
func (sp *ShallowSecurityProvider) CanView(typ string, attrName string, token string, salts ...interface{}) error {
	return nil
}
