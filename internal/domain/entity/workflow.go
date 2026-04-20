package entity

type Workflow struct {
	Name   string
	Runner string
	Jobs   []Job
}
