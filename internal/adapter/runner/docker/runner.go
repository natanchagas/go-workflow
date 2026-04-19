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

func (r *Runner) Execute(_ context.Context, step entity.Step) error {
	// TODO: implement Docker-based step execution using step.Type and step.Parameters
	_ = fmt.Sprintf("docker run step type=%s", step.Type)
	return nil
}
