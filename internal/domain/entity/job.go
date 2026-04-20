package entity

type Job struct {
	ID     string
	Name   string
	Runner string
	Steps  []Step
}
