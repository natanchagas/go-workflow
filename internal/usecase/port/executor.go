package port

import (
	"context"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
)

type StepExecutor interface {
	Execute(ctx context.Context, step entity.Step) error
}
