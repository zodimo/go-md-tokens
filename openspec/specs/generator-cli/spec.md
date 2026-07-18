# Generator CLI

## Purpose
TBD

## Requirements

### Requirement: CLI Execution
The system SHALL provide a CLI entrypoint at `internal/m3generator/cmd/main.go` that executes the generation pipeline.

#### Scenario: Running generation
- **WHEN** the user runs `go run ./internal/m3generator/cmd/main.go`
- **THEN** the system generates the token structures in `pkg/m3tokens` and the theme structures in `theme/tokens.go`.

### Requirement: Go Generate Hook
The system SHALL expose a root `generate.go` script to invoke the tool via `go generate ./...`.

#### Scenario: Generating via Go tools
- **WHEN** the user runs `go generate ./...`
- **THEN** the token system is regenerated using the internal CLI.
