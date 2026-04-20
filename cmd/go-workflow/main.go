package main

import (
	"context"
	"log"
	"os"

	"github.com/natanchagas/go-workflow/internal/adapter/runner/local"
	"github.com/natanchagas/go-workflow/internal/domain/entity"
	"github.com/natanchagas/go-workflow/internal/domain/step"
	jobUseCase "github.com/natanchagas/go-workflow/internal/usecase/job"
	"github.com/natanchagas/go-workflow/internal/usecase/port"
	workflowUseCase "github.com/natanchagas/go-workflow/internal/usecase/workflow"
)

func main() {
	executors := map[string]port.StepExecutor{
		"local": local.New(os.Stdout),
	}

	jobRunner := jobUseCase.NewRunner(executors)
	workflowRunner := workflowUseCase.NewRunner(jobRunner)

	workflow := entity.Workflow{
		Name: "example",
		Jobs: []entity.Job{
			{
				ID:     "job-1",
				Name:   "build",
				Runner: "local",
				Steps: []entity.Step{
					step.Shell{Command: "echo 'Hello, World!'"},
				},
			},
		},
	}

	if err := workflowRunner.Run(context.Background(), workflow); err != nil {
		log.Fatal(err)
	}
}
