package job

import (
	"context"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
	"github.com/natanchagas/go-workflow/internal/usecase/port"
)

type Runner interface {
	Run(ctx context.Context, job entity.Job) error
}

type runner struct {
	executor port.StepExecutor
}

func NewRunner(executor port.StepExecutor) Runner {
	return &runner{executor: executor}
}

func (r *runner) Run(ctx context.Context, job entity.Job) error {
	for _, step := range job.Steps {
		if err := r.executor.Execute(ctx, step); err != nil {
			return err
		}
	}
	return nil
}
