package entity

type Step interface {
	Type() string
	Args() ([]string, error)
}
