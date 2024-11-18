package interfaces

import (
	"github.com/saichler/shared/go/types"
	"net"
)

type ISecurityProvider interface {
	CanDial(string, uint32, ...interface{}) (net.Conn, error)
	CanAccept(net.Conn, ...interface{}) error
	ValidateConnection(net.Conn, *types.MessagingConfig, ...interface{}) (string, error)

	Encrypt([]byte, ...interface{}) (string, error)
	Decrypt(string, ...interface{}) ([]byte, error)

	CanDo(types.Action, string, string, ...interface{}) error
	CanView(string, string, string, ...interface{}) error
}

var securityProvider ISecurityProvider

func SecurityProvider() ISecurityProvider {
	return securityProvider
}

func SetSecurityProvider(sp ISecurityProvider) {
	securityProvider = sp
}
