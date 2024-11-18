package tests

import (
	. "github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"strings"
	"testing"
)

func TestShalowSecurity(t *testing.T) {
	sp := SecurityProvider()
	conn, err := sp.CanDial("127.0.0.1", 8910)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {

		Fail(t, err)
		return
	}
	err = sp.CanAccept(conn)
	if err != nil {
		Fail(t, err)
		return
	}
	conn = &MockConn{}
	config := EdgeConfig()
	config.Uuid = "Test Validate Connection"

	_, err = sp.ValidateConnection(conn, config)
	if err != nil {
		Fail(t, err)
		return
	}
	if config.IsAdjacentASwitch {
		Fail(t, "This connection is adjucent.")
		return
	}
	
	sp.CanDo(types.Action_GET, "", "")
	sp.CanView("", "", "")
}
