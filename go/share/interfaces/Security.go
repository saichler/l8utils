package interfaces

import (
	"errors"
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
	"net"
	"plugin"
)

type ISecurityProvider interface {
	Init(IResources)
	CanDial(string, uint32) (net.Conn, error)
	CanAccept(net.Conn) error
	ValidateConnection(net.Conn) error

	Encrypt([]byte) (string, error)
	Decrypt(string) ([]byte, error)

	CanDoAction(types.Action, proto.Message, string, string, ...string) error
	ScopeView(proto.Message, string, string, ...string) (proto.Message, error)
	Authenticate(string, string, ...string) string
}

func LoadSecurityProvider(path string) (ISecurityProvider, error) {
	securityProviderPlugin, err := plugin.Open(path)
	if err != nil {
		return nil, errors.New("failed to load security provider plugin #1")
	}
	securityProvider, err := securityProviderPlugin.Lookup("SecurityProvider")
	if err != nil {
		return nil, errors.New("failed to load security provider plugin #2")
	}
	if securityProvider == nil {
		return nil, errors.New("failed to load security provider plugin #3")
	}
	providerInterface := *securityProvider.(*ISecurityProvider)
	return providerInterface.(ISecurityProvider), nil
}
