package interfaces

type SerializerMode int

const (
	BINARY SerializerMode = 1
	JSON   SerializerMode = 2
)

type ISerializer interface {
	Mode() SerializerMode
	Marshal(interface{}, IRegistry) ([]byte, error)
	Unmarshal([]byte, string, IRegistry) (interface{}, error)
}
