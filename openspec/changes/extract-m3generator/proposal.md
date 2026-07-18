## Why

Currently, the generator logic and generated output are entangled in both `pkg/m3tokens` and `main.go`. This creates a dependency tangle, bloats the package, and obscures the extraction logic. Extracting the generator into `internal/m3generator` separates the build-time generation tools from the runtime library, improving maintainability, modularity, and compilation efficiency.

## What Changes

- Move extraction logic (`extractor.go`, `generator.go`, and related `schema-def.json`) from `pkg/m3tokens` to `internal/m3generator/generator/`.
- Create a new CLI entrypoint `internal/m3generator/cmd/main.go` to handle parsing, schema extraction, and theme generation.
- Remove hardcoded configurations in `main.go` and introduce an `m3gen.yaml` file for the parameters.
- Expose the pipeline generation via `go generate`.

## Capabilities

### New Capabilities
- `generator-cli`: A standalone CLI tool to parse, extract, and build the token architecture via a configuration file.
- `generator-config`: Implementation of the `m3gen.yaml` configuration format to orchestrate pipeline behavior.

### Modified Capabilities

- None

## Impact

- `pkg/m3tokens`: Will become a clean runtime package containing only `types.go`, `schema.json`, and the generated output, completely devoid of sass/transpiler dependencies.
- `main.go` (Root): Will be deleted and replaced with a `//go:generate` hook.
- The build process will require developers to provide/update `m3gen.yaml` for generation paths.
