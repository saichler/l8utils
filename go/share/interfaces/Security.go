package interfaces

import (
	"github.com/saichler/shared/go/types"
	"net"
)

type ISecurityProvider interface {
	CanDial(string, uint32) (net.Conn, error)
	CanAccept(net.Conn) error
	ValidateConnection(net.Conn, *types.VNicConfig) error

	Encrypt([]byte) (string, error)
	Decrypt(string) ([]byte, error)

	CanDo(types.Action, string, string) error
	CanView(string, string, string) error
}
