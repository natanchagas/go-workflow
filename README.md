# go-workflow

A platform-independent workflow runner written in Go. Define your workflows once and execute them on any supported backend — Docker, Kubernetes, GitHub Actions, AWS, Azure, or any custom target.

## Concepts

Workflows are built from three nested objects:

```
Workflow
└── Job(s)
    └── Step(s)
```

| Object | Description |
|--------|-------------|
| **Workflow** | A named collection of jobs executed in order. |
| **Job** | A logical unit of work with an ID, a name, and an ordered list of steps. |
| **Step** | A single executable action identified by a `Type` and configured via `Parameters`. |

## Architecture

The project follows [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) convention.

```
cmd/go-workflow/        entry point — wires the dependency graph
internal/
  domain/entity/        Layer 1: pure data objects, no external dependencies
  usecase/              Layer 2: application logic and port interfaces
    port/               output port interfaces (StepExecutor)
    job/                JobRunner use case
    workflow/           WorkflowRunner use case
  adapter/runner/       Layer 3: platform-specific StepExecutor implementations
    noop/               no-op executor (default / testing)
    docker/             Docker executor (in progress)
docs/adr/               Architecture Decision Records
```

Dependencies flow inward only: adapters depend on use cases, use cases depend on entities, entities depend on nothing. Adding a new execution platform means implementing `port.StepExecutor` in a new `adapter/runner/<platform>/` package — no other layer is touched.

See [docs/adr/0001-project-layout-and-clean-architecture.md](docs/adr/0001-project-layout-and-clean-architecture.md) for the full rationale.

## Getting Started

**Requirements:** Go 1.21+

```bash
# Run the example workflow
make run

# Build the binary
go build -o bin/go-workflow ./cmd/go-workflow
```

## Extending with a New Platform

1. Create `internal/adapter/runner/<platform>/runner.go`.
2. Implement the `port.StepExecutor` interface:

```go
type StepExecutor interface {
    Execute(ctx context.Context, step entity.Step) error
}
```

3. Wire the new runner in `cmd/go-workflow/main.go`.

## Roadmap

- [ ] Docker executor
- [ ] Workflow definition from YAML files
- [ ] Concurrent job execution
- [ ] Kubernetes executor
- [ ] GitHub Actions executor
