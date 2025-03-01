package nets

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
	"net"
)

func ExecuteProtocol(conn net.Conn, config *types.VNicConfig, security interfaces.ISecurityProvider) error {
	err := WriteEncrypted(conn, []byte(config.LocalUuid), config, security)
	if err != nil {
		conn.Close()
		return err
	}

	config.RemoteUuid, err = ReadEncrypted(conn, config, security)
	if err != nil {
		conn.Close()
		return err
	}

	forceExternal := "false"
	if config.ForceExternal {
		forceExternal = "true"
	}

	err = WriteEncrypted(conn, []byte(forceExternal), config, security)
	if err != nil {
		conn.Close()
		return err
	}

	forceExternal, err = ReadEncrypted(conn, config, security)
	if err != nil {
		conn.Close()
		return err
	}
	if forceExternal == "true" {
		config.ForceExternal = true
	}

	err = WriteEncrypted(conn, []byte(config.LocalAlias), config, security)
	if err != nil {
		conn.Close()
		return err
	}

	remoteAlias, err := ReadEncrypted(conn, config, security)
	if err != nil {
		conn.Close()
		return err
	}
	config.RemoteAlias = remoteAlias

	err = WriteEncrypted(conn, VlansToBytes(config.Vlans), config, security)
	if err != nil {
		conn.Close()
		return err
	}

	vlans, err := ReadEncryptedBytes(conn, config, security)
	if err != nil {
		conn.Close()
		return err
	}
	config.Vlans = BytesToVlans(vlans)

	return nil
}

func VlansToBytes(areas *types.Vlans) []byte {
	data, err := proto.Marshal(areas)
	if err != nil {
		return []byte{}
	}
	return data
}

func BytesToVlans(data []byte) *types.Vlans {
	areas := &types.Vlans{}
	err := proto.Unmarshal(data, areas)
	if err != nil {
		return nil
	}
	return areas
}
