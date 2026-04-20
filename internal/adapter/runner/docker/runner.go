package docker

import (
	"context"
	"fmt"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
)

type Runner struct{}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Execute(_ context.Context, s entity.Step) error {
	args, err := s.Args()
	if err != nil {
		return err
	}

	// TODO: spin up a container and run: docker exec <container> args...
	_ = args
	return fmt.Errorf("docker runner: not yet implemented")
}
