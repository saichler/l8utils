package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"testing"
)

func TestEdgeConfig(t *testing.T) {
	node := interfaces.EdgeConfig()
	node.Local_Uuid = "12345"
	sw := interfaces.SwitchConfig()
	sw.Local_Uuid = "54321"
	if node.Local_Uuid == sw.Local_Uuid {
		t.Fail()
		interfaces.Error("Expected uuid to be different")
		return
	}
	swEdge := interfaces.EdgeSwitchConfig()
	if swEdge.Local_Uuid == node.Local_Uuid {
		t.Fail()
		interfaces.Error("Expected uuid to be different")
		return
	}
}
