package job

import (
	"context"
	"fmt"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
	"github.com/natanchagas/go-workflow/internal/usecase/port"
)

type Runner interface {
	Run(ctx context.Context, job entity.Job) error
}

type runner struct {
	executors map[string]port.StepExecutor
}

func NewRunner(executors map[string]port.StepExecutor) Runner {
	return &runner{executors: executors}
}

func (r *runner) Run(ctx context.Context, job entity.Job) error {
	executor, ok := r.executors[job.Runner]
	if !ok {
		return fmt.Errorf("unknown runner %q", job.Runner)
	}
	for _, step := range job.Steps {
		if err := executor.Execute(ctx, step); err != nil {
			return err
		}
	}
	return nil
}
