package interfaces

type Serializer interface {
	Add(interface{}, IStructRegistry) ([]byte, int)
	Get([]byte, int, IStructRegistry) (interface{}, int)
}
