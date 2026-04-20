# ADR-0001: Project Layout, Clean Architecture, and Core Object Model

**Status:** Accepted  
**Date:** 2026-04-20  
**Deciders:** Natan das Chagas Silva  

---

## Context

`go-workflow` is a workflow runner and manager whose core value proposition is platform independence: the same workflow definition should execute on a local machine, a Docker container, a Kubernetes cluster, or any future target without changing the business logic.

The initial prototype had all three objects (`Workflow`, `Job`, `Step`) as concrete structs with execution logic baked directly into methods. Port interfaces existed but the model types did not implement them — `context.Context` was absent from signatures and the `Workflow` type skipped the `Job` level entirely, going straight from workflow to steps.

Several rounds of review shaped the current design. The decisions below are recorded together because they are interdependent: the layer structure drives where each object lives, and the object model drives what the layer boundaries enforce.

---

## Decisions

### 1. Project layout follows golang-standards/project-layout

Top-level directories follow the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) convention:

```
cmd/go-workflow/      ← entry point and composition root
internal/             ← all private application code
docs/adr/             ← architecture decision records
```

### 2. Internal structure follows Clean Architecture

`internal/` is organised into concentric layers. The dependency rule is strict: dependencies point inward only.

```
cmd/go-workflow/          → composition root (imports everything)

internal/
  domain/                 → Layer 1: innermost, no external deps
    entity/               → core entities and the Step interface
    step/                 → concrete step type definitions
  usecase/                → Layer 2: application orchestration
    port/                 → output port interfaces owned by use cases
    workflow/             → WorkflowRunner use case
    job/                  → JobRunner use case
  adapter/                → Layer 3: platform implementations
    runner/local/         → executes steps on the local machine
    runner/docker/        → executes steps in a Docker container (stub)
    runner/noop/          → no-op executor for testing
```

**Enforced import rules:**
1. `domain/` imports nothing outside the standard library.
2. `usecase/` imports `domain/` only — never `adapter/`.
3. `adapter/` imports `domain/` and `usecase/port/` — never other adapters.
4. `cmd/` is the only package that imports concrete adapter types and wires the dependency graph.

### 3. Object model: Workflow → Job → Step

The three core objects form a strict hierarchy:

| Object | Owns | Defines |
|--------|------|---------|
| `Workflow` | `[]Job` | a named, ordered collection of jobs |
| `Job` | `Runner string`, `[]Step` | where the job runs and what steps to execute |
| `Step` | — | what a single unit of work is and how to express it as a command |

`Workflow` and `Job` are concrete structs in `domain/entity/`. `Step` is an interface in `domain/entity/`, with concrete implementations in `domain/step/`.

### 4. Runner vs. Step type are orthogonal concerns

**Runner** (defined on `Job`) answers *where* the job executes:  
`"local"` → local machine shell, `"docker"` → Docker container, `"kubernetes"` → cluster pod, etc.

**Step type** (defined by the concrete `Step` implementation) answers *what* the step does and *how to express it as a command*:  
`step.Shell{Shell: "zsh", Command: "echo hello"}` → `["zsh", "-c", "echo hello"]`

These two axes are independent. The same step type can run on any runner; the same runner can execute any step type.

### 5. Step types live in `domain/step/`, not `domain/entity/`

`domain/entity/` holds the stable, universal objects (`Workflow`, `Job`, `Step` interface). Step type implementations (`step.Shell`, and future types) are value objects that live in `domain/step/`.

This separation is required by the dependency rule: both use cases and adapters need to reference concrete step types. If step types lived in `usecase/`, adapters would be forced to import a higher layer — a violation. Placing them in `domain/step/` keeps them importable by all inner layers without introducing cycles.

### 6. Steps own their own command expression via `Args()`

The `Step` interface requires:

```go
type Step interface {
    Type() string
    Args() ([]string, error)
}
```

`Args()` returns the full command as a slice — e.g. `["zsh", "-c", "echo hello"]`. This means:

- **Step types** encapsulate their own execution semantics: which program, which flags, validation of required fields.
- **Runners** are generic executors: they call `Args()`, then run the resulting command in their environment. No runner contains step-type-specific logic.

The local runner resolves the program via `exec.LookPath` (environment concern) and invokes it. The future Docker runner will wrap the same args in a `docker exec` call. The args are identical; only the envelope changes.

### 7. Runners are resolved at runtime via a registry

The `usecase/job` runner holds a `map[string]port.StepExecutor` keyed by runner name. At execution time it looks up `job.Runner` in the map and fails fast with a clear error if the runner is not registered. New platforms are added by registering an additional entry in `cmd/go-workflow/main.go` — no use-case or domain code changes.

---

## Options Considered (summary)

| Concern | Chosen | Rejected |
|---------|--------|----------|
| Step type location | `domain/step/` — importable by all layers | `usecase/step/` — would force adapters to import use cases |
| Step execution semantics | Step owns `Args()` | Runner owns type switch — runner accumulates step-type knowledge with every new type |
| Runner selection | String registry in job use case | Single executor injected at construction — no multi-runner support |
| Layer structure | Clean Architecture (concentric) | Flat layered-by-type — no enforced boundary between logic and infrastructure |

---

## Consequences

### What becomes easier
- Adding a new execution platform: implement `port.StepExecutor`, register it in `main.go`.
- Adding a new step type: add a struct to `domain/step/` that implements `Args()`. No runner changes needed.
- Testing use-case logic: inject `noop.Runner` — no real infrastructure required.

### What becomes harder
- Understanding the layer model requires familiarity with Clean Architecture; more directories than a flat layout.
- Every new step type must express itself fully through `Args()` — step types that don't map cleanly to a single shell-style command will need a richer interface.

### What we will need to revisit
- `Workflow.Runner` field is defined but currently unused. Intent unclear — candidate for either removal or use as a default runner for jobs that don't specify one.
- If a step needs to return structured output (logs, artifacts, exit codes) back to the use case, `port.StepExecutor` and `entity.Step` will both need extension.
- Concurrent job execution within a workflow is not addressed; the `WorkflowRunner` currently runs jobs sequentially.
- `adapter/runner/docker` is a stub. Implementing it requires deciding how container lifecycle is managed (pre-existing container vs. spin-up per job).

---

## Action Items

1. [x] Restructure project to match this layout.
2. [x] Implement `step.Shell` with `Args()` for local shell execution.
3. [x] Implement `adapter/runner/local` as a generic executor.
4. [ ] Clarify or remove `Workflow.Runner` field.
5. [ ] Implement `adapter/runner/docker`.
6. [ ] Add a `WorkflowRepository` port + YAML adapter so workflows can be loaded from files.
7. [ ] Revisit `WorkflowRunner` for concurrent job execution.
