package local

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/natanchagas/go-workflow/internal/domain/entity"
)

type Runner struct {
	output io.Writer
}

func New(output io.Writer) *Runner {
	return &Runner{output: output}
}

func (r *Runner) Execute(ctx context.Context, s entity.Step) error {
	args, err := s.Args()
	if err != nil {
		return err
	}

	program, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("%q not found: %w", args[0], err)
	}

	cmd := exec.CommandContext(ctx, program, args[1:]...)
	cmd.Stdout = r.output
	cmd.Stderr = r.output

	return cmd.Run()
}
