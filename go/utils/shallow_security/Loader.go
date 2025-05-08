package main

import "github.com/saichler/l8types/go/ifs"

var Loader ifs.ISecurityProviderLoader = &ShallowLoader{}

type ShallowLoader struct {
}

func (this *ShallowLoader) LoadSecurityProvider() (ifs.ISecurityProvider, error) {
	return NewShallowSecurityProvider(), nil
}
