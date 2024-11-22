package nets

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"net"
	"time"
)

// Write data to socket
func Write(data []byte, conn net.Conn, config *types.MessagingConfig) error {
	// If the connection is nil, return an error
	if conn == nil {
		return errors.New("no Connection Available")
	}
	// If the config is nil, error
	if config == nil {
		return errors.New("no Config Available")
	}
	if data == nil {
		return errors.New("no Data Available")
	}
	// Error is the data is too big
	if len(data) > int(config.MaxDataSize) {
		return errors.New("data is larger than MAX size allowed")
	}
	// Write the size of the data
	_, e := conn.Write(Long2Bytes(int64(len(data))))
	if e != nil {
		return e
	}
	// Write the actual data
	_, e = conn.Write(data)
	time.Sleep(time.Millisecond * 1)
	return e
}

func WriteEncrypted(conn net.Conn, data []byte, config *types.MessagingConfig, salts ...interface{}) error {
	encData, err := interfaces.SecurityProvider().Encrypt(data, salts...)
	if err != nil {
		return err
	}
	err = Write([]byte(encData), conn, config)
	if err != nil {
		return err
	}
	return nil
}
