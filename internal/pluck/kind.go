package plucker

type Kind string

const (
	Type Kind = "type"
	Func Kind = "function"
)

func (k Kind) Valid() bool {
	switch k {
	case Type, Func:
		return true
	default:
		return false
	}
}
