package main

import (
	"context"
	"log"

	"github.com/natanchagas/go-workflow/internal/adapter/runner/noop"
	"github.com/natanchagas/go-workflow/internal/domain/entity"
	jobUseCase "github.com/natanchagas/go-workflow/internal/usecase/job"
	workflowUseCase "github.com/natanchagas/go-workflow/internal/usecase/workflow"
)

func main() {
	executor := noop.New()
	jobRunner := jobUseCase.NewRunner(executor)
	workflowRunner := workflowUseCase.NewRunner(jobRunner)

	workflow := entity.Workflow{
		Name: "example",
		Jobs: []entity.Job{
			{
				ID:   "job-1",
				Name: "build",
				Steps: []entity.Step{
					{Type: "shell", Parameters: map[string]string{"command": "echo hello"}},
				},
			},
		},
	}

	if err := workflowRunner.Run(context.Background(), workflow); err != nil {
		log.Fatal(err)
	}
}
