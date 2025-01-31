package interfaces

import (
	"github.com/saichler/shared/go/types"
	"reflect"
)

type IIntrospector interface {
	Inspect(interface{}) (*types.RNode, error)
	Node(string) (*types.RNode, bool)
	NodeByType(p reflect.Type) (*types.RNode, bool)
	NodeByTypeName(string) (*types.RNode, bool)
	NodeByValue(interface{}) (*types.RNode, bool)
	Nodes(bool, bool) []*types.RNode
	Registry() IRegistry
	Kind(*types.RNode) reflect.Kind
	Clone(interface{}) interface{}
	AddDecorator(types.DecoratorType, interface{}, *types.RNode)
	DecoratorOf(types.DecoratorType, *types.RNode) interface{}
	TableView(string) (*types.TableView, bool)
	TableViews() []*types.TableView
}
