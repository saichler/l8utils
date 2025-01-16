package shallow_security

import (
	"errors"
	"github.com/saichler/shared/go/share/aes"
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

func (sp *ShallowSecurityProvider) CanDial(host string, port uint32) (net.Conn, error) {
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
}

func (sp *ShallowSecurityProvider) CanAccept(conn net.Conn) error {
	return nil
}

func (sp *ShallowSecurityProvider) ValidateConnection(conn net.Conn, config *types.VNicConfig) error {
	err := nets.WriteEncrypted(conn, []byte(sp.secret), config, sp)
	if err != nil {
		conn.Close()
		return err
	}

	secret, err := nets.ReadEncrypted(conn, config, sp)
	if err != nil {
		conn.Close()
		return err
	}

	if sp.secret != secret {
		conn.Close()
		return errors.New("incorrect Secret/Key, aborting connection")
	}

	err = nets.WriteEncrypted(conn, []byte(config.Local_Uuid), config, sp)
	if err != nil {
		conn.Close()
		return err
	}

	config.RemoteUuid, err = nets.ReadEncrypted(conn, config, sp)
	if err != nil {
		conn.Close()
		return err
	}

	forceExternal := "false"
	if config.ForceExternal {
		forceExternal = "true"
	}

	err = nets.WriteEncrypted(conn, []byte(forceExternal), config, sp)
	if err != nil {
		conn.Close()
		return err
	}

	forceExternal, err = nets.ReadEncrypted(conn, config, sp)
	if err != nil {
		conn.Close()
		return err
	}
	if forceExternal == "true" {
		config.ForceExternal = true
	}

	return nil
}

func (sp *ShallowSecurityProvider) Encrypt(data []byte) (string, error) {
	return aes.Encrypt(data, sp.key)
}

func (sp *ShallowSecurityProvider) Decrypt(data string) ([]byte, error) {
	return aes.Decrypt(data, sp.key)
}

func (sp *ShallowSecurityProvider) CanDo(action types.Action, endpoint string, token string) error {
	return nil
}
func (sp *ShallowSecurityProvider) CanView(typ string, attrName string, token string) error {
	return nil
}
