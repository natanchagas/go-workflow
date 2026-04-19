# ADR-0001: Project Layout and Clean Architecture

**Status:** Accepted  
**Date:** 2026-04-19  
**Deciders:** Natan das Chagas Silva  

---

## Context

`go-workflow` is a workflow runner and manager whose core value proposition is platform independence: the same workflow definition should execute on Docker, Kubernetes, GitHub Actions, AWS, Azure, or any future target without changing the business logic.

The initial prototype had all three objects (`Workflow`, `Job`, `Step`) as concrete structs with execution logic baked directly into methods (`Dispatch()`, `Run()`, `Execute()`). Port interfaces existed in `internal/domain/ports/` but the model types did not implement them — `context.Context` was absent from the method signatures. The `Workflow` type also had a `Steps []Step` field, skipping the `Job` level entirely.

The design needed to:

- Make it cheap to add new execution platforms without touching business logic.
- Establish a dependency rule that prevents infrastructure concerns from bleeding into domain objects.
- Follow a community-standard layout so contributors can orient themselves quickly.

---

## Decision

Adopt **Clean Architecture** (concentric layers, dependency rule pointing inward) for internal structure, and the **golang-standards/project-layout** convention for top-level directories.

---

## Options Considered

### Option A: Clean Architecture + golang-standards layout *(chosen)*

Organise `internal/` into three concentric layers. Dependencies flow inward only: adapters depend on use cases, use cases depend on domain entities, entities depend on nothing.

```
cmd/go-workflow/          ← composition root
internal/
  domain/entity/          ← Layer 1 – entities (no deps)
  usecase/                ← Layer 2 – application logic
    port/                 ← output port interfaces owned by use cases
    job/
    workflow/
  adapter/runner/         ← Layer 3 – infrastructure adapters
    noop/
    docker/
```

| Dimension | Assessment |
|-----------|------------|
| Complexity | Medium — more directories than a flat layout |
| Extensibility | High — new platform = new package under `adapter/runner/` |
| Testability | High — use cases testable with a `noop` executor, no real infra needed |
| Team familiarity | Medium — Clean Architecture is widely documented; golang-standards is the de-facto Go convention |

**Pros:**
- Adding a Kubernetes, GitHub Actions, or AWS runner requires only a new `adapter/runner/<platform>/` package implementing one interface.
- Domain entities carry zero external dependencies and can be reasoned about in isolation.
- Use-case logic is independently testable without Docker, a cluster, or cloud credentials.
- `cmd/go-workflow/main.go` is the single place that knows about concrete types; everything else talks to interfaces.

**Cons:**
- More directories than a flat or layered-by-type approach.
- New contributors must understand the concentric-layer model.

---

### Option B: Layered architecture (packages by type)

```
internal/
  model/       ← all structs
  service/     ← all business logic
  runner/      ← all runners
```

| Dimension | Assessment |
|-----------|------------|
| Complexity | Low — fewer directories |
| Extensibility | Medium — adding a runner is straightforward, but service layer may accumulate cross-cutting platform concerns |
| Testability | Medium — services depend directly on runner implementations without an interface boundary |
| Team familiarity | High — familiar MVC/layered pattern |

**Pros:** Simpler to navigate for small teams unfamiliar with Clean Architecture.

**Cons:**
- No enforced boundary between business logic and infrastructure; platform-specific concerns tend to leak into services over time.
- Harder to test business logic without a real execution environment.

---

## Trade-off Analysis

The defining requirement is **platform independence**: the same `Workflow`/`Job`/`Step` model must execute on arbitrarily different backends. This is precisely the problem the `StepExecutor` output port solves — it is the seam between application logic and infrastructure. Option B can achieve this with discipline, but Clean Architecture makes the boundary structural and therefore enforceable at compile time via Go's package import rules.

The cost is directory depth. Given that the number of platform adapters is expected to grow, this cost is paid once and amortised across every new backend.

---

## Consequences

### What becomes easier
- Adding a new execution platform: implement `port.StepExecutor`, wire it in `main.go`, done.
- Unit-testing workflow orchestration logic: inject `noop.Runner` — no external dependencies.
- Onboarding contributors: the golang-standards layout makes top-level navigation predictable.

### What becomes harder
- Simple one-off changes now require awareness of which layer a change belongs to.
- The `Workflow → Job → Step` hierarchy is now fixed in the entity model; flattening it later would be a breaking change.

### What we will need to revisit
- If a step executor needs to return output (e.g. logs, artifacts) back to the use case, `port.StepExecutor` will need to be extended.
- If workflows need to be persisted or loaded from YAML/JSON, a `WorkflowRepository` output port and a corresponding adapter will be needed in a future ADR.
- Concurrent job execution within a workflow is not addressed by the current `WorkflowRunner` implementation; this is an intentional deferral.

---

## Key Structural Rules

1. `internal/domain/entity` imports nothing outside the standard library.
2. `internal/usecase` imports `domain/entity` and `usecase/port` only — never `adapter`.
3. `internal/adapter` imports `usecase/port` and `domain/entity` — never other adapters.
4. `cmd/` is the only package that imports concrete adapter types and wires the dependency graph.

---

## Action Items

1. [x] Restructure project to match this layout.
2. [ ] Implement `docker.Runner.Execute` using the Docker SDK.
3. [ ] Add a `WorkflowRepository` port + YAML adapter so workflows can be loaded from files.
4. [ ] Revisit `WorkflowRunner` for concurrent job execution.
