package registry

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

/*
Type info contains only two attributes, the golang type of, usually, a struct.
And a Serializer, which indicates a which serializer should be used to transalte
this type into bytes. Serializer can be nil...
*/
type Info struct {
	/* The golang reflect type */
	typ reflect.Type
	/* The serializers */
	serializers map[interfaces.SerializerMode]interfaces.ISerializer
}

// NewInfo /* Constructs a new type info with the given attributes */
func NewInfo(typ reflect.Type) (*Info, error) {
	if typ == nil {
		return nil, errors.New("Cannot register a nil type")
	}
	return &Info{typ: typ,
		serializers: make(map[interfaces.SerializerMode]interfaces.ISerializer)}, nil
}

// Type /* Return the reflect type of this TypeInfo */
func (info *Info) Type() reflect.Type {
	return info.typ
}

// Name /* Return the name of the type/struct */
func (info *Info) Name() string {
	return info.typ.Name()
}

// Serializer /* Return the serializer to be used with this type. Can be nil... */
func (info *Info) Serializer(mode interfaces.SerializerMode) interfaces.ISerializer {
	return info.serializers[mode]
}

func (info *Info) AddSerializer(serializer interfaces.ISerializer) {
	info.serializers[serializer.Mode()] = serializer
}

// NewInstance /* Construct a new instance via reflect using the type */
func (info *Info) NewInstance() (interface{}, error) {
	n := reflect.New(info.typ)
	if !n.IsValid() {
		return nil, errors.New("was not able to create new instance of type " + info.typ.Name())
	}
	return n.Interface(), nil
}
