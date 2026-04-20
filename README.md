# go-workflow

A platform-independent workflow runner written in Go. Define your workflows once and execute them on any supported backend — a local machine, a Docker container, a Kubernetes cluster, or any future target.

## Concepts

Workflows are built from three nested objects:

```
Workflow
└── Job(s)       ← declares where it runs (runner)
    └── Step(s)  ← declares what to execute and how
```

| Object | Description |
|--------|-------------|
| **Workflow** | A named, ordered collection of jobs. |
| **Job** | A unit of work that declares a **runner** (where it executes) and an ordered list of steps. |
| **Step** | A single executable unit. Defines what to run and expresses itself as a command via `Args()`. |

### Runners

A runner defines the execution environment for a job. Each job declares its runner by name:

| Runner | Description |
|--------|-------------|
| `local` | Runs steps directly on the local machine. |
| `docker` | Runs steps inside a Docker container. *(in progress)* |

### Step types

A step type defines what gets executed. Every step type implements `Args() ([]string, error)`, which returns the command to run — e.g. `["zsh", "-c", "echo hello"]`. Runners consume this without any knowledge of the step type.

| Type | Parameters | Description |
|------|------------|-------------|
| `shell` | `Shell` (optional, defaults to `sh`), `Command` (required) | Runs a shell command. |

Runners and step types are independent: any step type can run on any runner.

## Architecture

The project follows [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) convention.

```
cmd/go-workflow/          entry point — wires the dependency graph
internal/
  domain/
    entity/               Layer 1: Workflow, Job, and the Step interface
    step/                 Layer 1: concrete step types (Shell, …)
  usecase/                Layer 2: application orchestration
    port/                 output port interfaces (StepExecutor)
    workflow/             WorkflowRunner use case
    job/                  JobRunner use case
  adapter/runner/         Layer 3: platform implementations
    local/                local machine executor
    docker/               Docker executor (in progress)
    noop/                 no-op executor (testing)
docs/adr/                 Architecture Decision Records
```

Dependencies flow inward only — adapters depend on use cases, use cases depend on the domain, the domain depends on nothing. The full rationale is in [docs/adr/0001-project-layout-and-clean-architecture.md](docs/adr/0001-project-layout-and-clean-architecture.md).

## Getting Started

**Requirements:** Go 1.21+

```bash
# Run the example workflow
make run

# Build the binary
go build -o bin/go-workflow ./cmd/go-workflow
```

## Extending with a New Runner

1. Create `internal/adapter/runner/<platform>/runner.go`.
2. Implement `port.StepExecutor`:

```go
type StepExecutor interface {
    Execute(ctx context.Context, step entity.Step) error
}
```

Call `step.Args()` to get the command — no step-type logic belongs in the runner.

3. Register it in `cmd/go-workflow/main.go`:

```go
executors := map[string]port.StepExecutor{
    "local":      local.New(os.Stdout),
    "<platform>": myrunner.New(...),
}
```

## Extending with a New Step Type

1. Create `internal/domain/step/<type>.go`.
2. Implement `entity.Step`:

```go
type Step interface {
    Type() string
    Args() ([]string, error) // full command: e.g. ["bash", "-c", "echo hi"]
}
```

No runner changes are needed — runners call `Args()` and execute whatever comes back.

## Roadmap

- [ ] Docker executor
- [ ] Workflow definition from YAML files
- [ ] Concurrent job execution
- [ ] Kubernetes executor
