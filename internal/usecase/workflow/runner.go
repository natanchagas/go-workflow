package workflow

import (
	"context"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
	"github.com/natanchagas/go-workflow/internal/usecase/job"
)

type Runner interface {
	Run(ctx context.Context, workflow entity.Workflow) error
}

type runner struct {
	jobRunner job.Runner
}

func NewRunner(jobRunner job.Runner) Runner {
	return &runner{jobRunner: jobRunner}
}

func (r *runner) Run(ctx context.Context, workflow entity.Workflow) error {
	for _, j := range workflow.Jobs {
		if err := r.jobRunner.Run(ctx, j); err != nil {
			return err
		}
	}
	return nil
}
