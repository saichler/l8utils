package nets

import (
	"bytes"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"net"
	"strings"
)

func ExecuteProtocol(conn net.Conn, config *types.VNicConfig, security interfaces.ISecurityProvider) error {
	err := WriteEncrypted(conn, []byte(config.Local_Uuid), config, security)
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

	err = WriteEncrypted(conn, []byte(SetToString(config.Topics)), config, security)
	if err != nil {
		conn.Close()
		return err
	}

	topics, err := ReadEncrypted(conn, config, security)
	if err != nil {
		conn.Close()
		return err
	}
	config.Topics = StringToSet(topics)

	return nil
}

func SetToString(set map[string]bool) string {
	buff := bytes.Buffer{}
	first := true
	for topic, _ := range set {
		if !first {
			buff.WriteString(",")
		}
		buff.WriteString(topic)
		first = false
	}
	return buff.String()
}

func StringToSet(str string) map[string]bool {
	result := make(map[string]bool)
	tokens := strings.Split(str, ",")
	for _, token := range tokens {
		result[token] = true
	}
	return result
}
