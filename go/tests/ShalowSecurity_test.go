package tests

import (
	. "github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"strings"
	"testing"
)

func TestShalowSecurity(t *testing.T) {
	sp := globals.Security()
	conn, err := sp.CanDial("127.0.0.1", 8910)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {

		log.Fail(t, err)
		return
	}
	err = sp.CanAccept(conn)
	if err != nil {
		log.Fail(t, err)
		return
	}
	conn = &MockConn{}
	c := globals.Config(EdgeConfig)
	config := &c
	config.Local_Uuid = "Test Validate Connection"

	err = sp.ValidateConnection(conn, config)
	if err != nil {
		log.Fail(t, err)
		return
	}
	if config.IsAdjacentASwitch {
		log.Fail(t, "This connection is adjucent.")
		return
	}

	sp.CanDo(types.Action_GET, "", "")
	sp.CanView("", "", "")
}
