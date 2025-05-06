package main

import "github.com/saichler/types/go/common"

var Loader common.ISecurityProviderLoader = &ShallowLoader{}

type ShallowLoader struct {
}

func (this *ShallowLoader) LoadSecurityProvider() (common.ISecurityProvider, error) {
	return NewShallowSecurityProvider(), nil
}
