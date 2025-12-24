package pluck

type Plucker interface {
	Pluck(code string, name string, kind Kind) (snippet string, err error)
}
