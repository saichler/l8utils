package tests

import (
	"bytes"
	"github.com/saichler/types/go/nets"
	"net"
	"sync"
	"testing"
	"time"
)

type MockConn struct {
	data []byte
	mtx  sync.RWMutex
}
type MockAddr struct{}

func TestNets(t *testing.T) {
	conn := &MockConn{}
	writeData := []byte("Testing Read/Write data to socket")
	config := globals.Config()
	config.LocalUuid = "abcde"

	err := nets.Write(nil, nil, nil)
	if err == nil {
		log.Fail(t, "Error is nil")
		return
	}

	err = nets.Write(nil, conn, nil)
	if err == nil {
		log.Fail(t, "Error is nil")
		return
	}

	err = nets.Write(nil, conn, config)
	if err == nil {
		log.Fail(t, "Error is nil")
		return
	}

	err = nets.Write(make([]byte, config.MaxDataSize+1), conn, config)
	if err == nil {
		log.Fail(t, "Error is nil")
		return
	}

	go func() {
		time.Sleep(time.Millisecond * 500)
		err = nets.Write(writeData, conn, config)
		if err != nil {
			log.Fail(t, err)
			return
		}
	}()

	_, err = nets.Read(nil, nil)
	if err == nil {
		log.Fail(t, "Error is nil")
		return
	}

	_, err = nets.Read(conn, nil)
	if err == nil {
		log.Fail(t, "Error is nil")
		return
	}

	readData, err := nets.Read(conn, config)
	if err != nil {
		log.Fail(t, err)
		return
	}
	if bytes.Compare(writeData, readData) != 0 {
		log.Fail(t, "Write Data & Read Date do not match")
		return
	}

	err = nets.WriteEncrypted(conn, writeData, config, globals.Security())
	if err != nil {
		log.Fail(t, err)
		return
	}
	readStr, err := nets.ReadEncrypted(conn, config, globals.Security())
	if err != nil {
		log.Fail(t, err)
		return
	}
	if readStr != string(writeData) {
		log.Fail(t, "Write Data do not match read data")
		return
	}
}

func (c *MockConn) Read(b []byte) (n int, err error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	index := 0
	for index = 0; index < len(b) && index < len(c.data); index++ {
		b[index] = c.data[index]
	}
	c.data = c.data[index:]
	return index, nil
}

func (c *MockConn) Write(b []byte) (n int, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.data == nil {
		c.data = make([]byte, 0)
	}
	c.data = append(c.data, b...)
	return len(b), nil
}

func (c *MockConn) Close() error {
	return nil
}

func (c *MockConn) LocalAddr() net.Addr {
	return &MockAddr{}
}
func (c *MockConn) RemoteAddr() net.Addr {
	return &MockAddr{}
}
func (c *MockConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (a *MockAddr) Network() string {
	return "tcp"
}

func (a *MockAddr) String() string {
	return "127.0.0.1"
}
