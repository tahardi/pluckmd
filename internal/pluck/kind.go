package pluck

type Kind string

const (
	Type Kind = "type"
	Func Kind = "function"
	Node Kind = "node"
	File Kind = "file"
)

func (k Kind) Valid() bool {
	switch k {
	case Type, Func, Node, File:
		return true
	default:
		return false
	}
}
