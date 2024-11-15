package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"testing"
)

func TestEdgeConfig(t *testing.T) {
	node := interfaces.EdgeConfig()
	node.Uuid = "12345"
	sw := interfaces.SwitchConfig()
	sw.Uuid = "54321"
	if node.Uuid == sw.Uuid {
		t.Fail()
		interfaces.Error("Expected uuid to be different")
		return
	}
	swEdge := interfaces.EdgeSwitchConfig()
	if swEdge.Uuid == node.Uuid {
		t.Fail()
		interfaces.Error("Expected uuid to be different")
		return
	}
}
