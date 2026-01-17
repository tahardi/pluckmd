package pluck

type Lang string

const (
	Go   Lang = "go"
	YAML Lang = "yaml"
)

func (l Lang) Valid() bool {
	switch l {
	case Go, YAML:
		return true
	default:
		return false
	}
}
