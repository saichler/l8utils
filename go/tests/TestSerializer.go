package tests

import "github.com/saichler/l8types/go/ifs"

type TestSerializer struct {
}

func (ts *TestSerializer) Mode() ifs.SerializerMode {
	return ifs.BINARY
}
func (ts *TestSerializer) Marshal(interface{}, ifs.IResources) ([]byte, error) {
	return nil, nil
}
func (ts *TestSerializer) Unmarshal([]byte, ifs.IResources) (interface{}, error) {
	return nil, nil
}
