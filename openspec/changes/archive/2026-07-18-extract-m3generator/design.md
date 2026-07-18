## Context

The current `go-md-tokens` repository uses `main.go` to parse Dart Sass files, generate an intermediate schema, and build theme components. However, this same `main.go` calls generation logic in `pkg/m3tokens` and then proceeds to do local `theme/` generation, while configurations are hardcoded.

## Goals / Non-Goals

**Goals:**
- Extract the extraction tools and schema generation from `pkg/m3tokens` into a dedicated `internal/m3generator` package.
- Move configurations to an `m3gen.yaml` config file.
- Clean up `pkg/m3tokens` to only contain runtime-needed generated code and the base schema file.
- Enable `go generate` workflows via a wrapper script.

**Non-Goals:**
- Changes to the SCSS schema or mapping logic.
- Generating additional file formats outside Go.

## Decisions

- **Config File Location:** Config will be named `m3gen.yaml` and placed at the root of the project.
- **Generator Binary:** The internal CLI will be located at `internal/m3generator/cmd/main.go`. This prevents the root of the repo from being polluted and ensures it remains an internal tool.

## Risks / Trade-offs

- [Risk] Breakage in theme generation downstream → Run existing tests and validation logic after extraction to ensure output equality.
