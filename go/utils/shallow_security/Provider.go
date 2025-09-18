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

type ShallowSecurityProvider struct {
	secret string
	key    string
}

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

func (this *ShallowSecurityProvider) CanDial(host string, port uint32) (net.Conn, error) {
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
}

func (this *ShallowSecurityProvider) CanAccept(conn net.Conn) error {
	return nil
}

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

func (this *ShallowSecurityProvider) Encrypt(data []byte) (string, error) {
	return aes.Encrypt(data, this.key)
}

func (this *ShallowSecurityProvider) Decrypt(data string) ([]byte, error) {
	return aes.Decrypt(data, this.key)
}

func (this *ShallowSecurityProvider) CanDoAction(action ifs.Action, o ifs.IElements, uuid string, token string, salts ...string) error {
	return nil
}
func (this *ShallowSecurityProvider) ScopeView(o ifs.IElements, uuid string, token string, salts ...string) ifs.IElements {
	return o
}
func (this *ShallowSecurityProvider) Authenticate(user string, pass string) (string, error) {
	return "token", nil
}
func (this *ShallowSecurityProvider) Message(string) (*ifs.Message, error) {
	return &ifs.Message{}, nil
}
