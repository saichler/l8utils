package interfaces

type SerializerMode int

const (
	BINARY SerializerMode = 1
	JSON   SerializerMode = 2
)

type Serializer interface {
	Mode() SerializerMode
	Marshal(interface{}, ITypeRegistry) ([]byte, error)
	Unmarshal([]byte, string, ITypeRegistry) (interface{}, error)
}
