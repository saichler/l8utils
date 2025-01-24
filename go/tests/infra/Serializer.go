package infra

import "github.com/saichler/shared/go/share/interfaces"

type TestSerializer struct {
}

func (ts *TestSerializer) Mode() interfaces.SerializerMode {
	return interfaces.BINARY
}
func (ts *TestSerializer) Marshal(interface{}, interfaces.IRegistry) ([]byte, error) {
	return nil, nil
}
func (ts *TestSerializer) Unmarshal([]byte, string, interfaces.IRegistry) (interface{}, error) {
	return nil, nil
}
