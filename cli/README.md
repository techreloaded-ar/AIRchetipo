# archetipo CLI

Deterministic Go implementation of the ARchetipo connector contracts. Replaces the markdown connector files (`file.md`, `github.md`) with one binary that performs every operation defined in `.archetipo/contracts.md`.

## Build

```bash
cd cli
go build ./cmd/archetipo
```

The output binary `archetipo` reads `.archetipo/config.yaml` from the project root (or any ancestor) to choose the connector (`file` or `github`) and execute the requested sub-command.

## Layout

```
cmd/archetipo/        # entry point
internal/
  cli/                # cobra sub-commands (one file per operation)
  connector/          # interface, registry, two implementations + inmemory ref
    filefs/           # markdown + HTML-comment markers
    github/           # gh CLI shell-out + GraphQL aliased mutations
    inmemory/         # reference impl used by the conformance suite
    conformance/      # behavioural test suite shared by every implementation
  config/             # .archetipo/config.yaml loader
  domain/             # canonical data types
  iox/                # JSON envelope on stdin/stdout/stderr + typed errors
  version/            # injected via -ldflags at release time
```

## Tests

```bash
go test ./...
```

The conformance suite runs against `filefs` and `inmemory`. The `github` connector is exercised with a mock `gh` runner; live smoke tests are gated behind `ARCHETIPO_E2E_GH=1` and need a sandbox repo with `gh` authenticated.

## Distribution

Tags `vX.Y.Z` produce a single bare binary per platform via goreleaser. Release assets are named `archetipo-<os>-<arch>` and downloaded by `install.sh` / `install.ps1` into `.archetipo/bin/` of the target project.
