package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/saichler/shared/go/share/aes"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/nets"
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
	"strings"
)

var SecurityProvider interfaces.ISecurityProvider = createShallowSecurityProvider()

type ShallowSecurityProvider struct {
	secret    string
	key       string
	salts     []string
	resources interfaces.IResources
}

func NewShallowSecurityProvider(key, secret string, salts ...string) *ShallowSecurityProvider {
	sp := &ShallowSecurityProvider{}
	sp.key = key
	sp.secret = secret
	sp.salts = salts
	return sp
}

func (this *ShallowSecurityProvider) Init(resources interfaces.IResources) {
	if this.resources == nil {
		this.resources = resources
	}
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

func (this *ShallowSecurityProvider) ValidateConnection(conn net.Conn) error {
	err := nets.WriteEncrypted(conn, []byte(this.secret), this.resources.Config(), this)
	if err != nil {
		conn.Close()
		return err
	}

	secret, err := nets.ReadEncrypted(conn, this.resources.Config(), this)
	if err != nil {
		conn.Close()
		return err
	}

	if this.secret != secret {
		conn.Close()
		return errors.New("incorrect Secret/Key, aborting connection")
	}

	return nets.ExecuteProtocol(conn, this.resources.Config(), this)
}

func (this *ShallowSecurityProvider) Encrypt(data []byte) (string, error) {
	return aes.Encrypt(data, this.key)
}

func (this *ShallowSecurityProvider) Decrypt(data string) ([]byte, error) {
	return aes.Decrypt(data, this.key)
}

func (this *ShallowSecurityProvider) CanDoAction(action types.Action, pb proto.Message, uuid string, token string, salts ...string) error {
	return nil
}
func (this *ShallowSecurityProvider) ScopeView(pb proto.Message, uuid string, token string, salts ...string) (proto.Message, error) {
	return pb, nil
}
func (this *ShallowSecurityProvider) Authenticate(user string, pass string, salts ...string) string {
	return "token"
}

func createShallowSecurityProvider() interfaces.ISecurityProvider {
	hash := md5.New()
	secret := "Shallow Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	k := base64.StdEncoding.EncodeToString(kHash)
	return NewShallowSecurityProvider(k, secret)
}
