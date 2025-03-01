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

	err = WriteEncrypted(conn, TopicsToBytes(config.Topics), config, security)
	if err != nil {
		conn.Close()
		return err
	}

	topics, err := ReadEncryptedBytes(conn, config, security)
	if err != nil {
		conn.Close()
		return err
	}
	config.Topics = BytesToTopics(topics)

	return nil
}

func TopicsToBytes(topics *types.Topics) []byte {
	data, err := proto.Marshal(topics)
	if err != nil {
		return []byte{}
	}
	return data
}

func BytesToTopics(data []byte) *types.Topics {
	topics := &types.Topics{}
	err := proto.Unmarshal(data, topics)
	if err != nil {
		return nil
	}
	return topics
}
