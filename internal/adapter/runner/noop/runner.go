package noop

import (
	"context"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
)

type Runner struct{}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Execute(_ context.Context, _ entity.Step) error {
	return nil
}
