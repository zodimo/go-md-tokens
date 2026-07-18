## Why

The current M3 generator heavily relies on imperative string building (`strings.Builder`) to generate the Go theme files and Sass bridge code. This approach intertwines the logic for gathering tokens with the formatting of the output, making the code harder to read, maintain, and modify. Moving to a declarative `text/template` approach with external `.tmpl` files will separate these concerns, providing clear syntax highlighting and making the structure of the generated files immediately obvious.

## What Changes

- Replace the manual `strings.Builder` logic in `internal/m3generator/cmd/main.go` with Go's `text/template` engine.
- Create a new `internal/m3generator/templates/` directory to house `.tmpl` files.
- Move the Go theme generation logic into `go_theme.tmpl`.
- Move the Sass bridge generation logic into `sass_bridge.tmpl`.
- Utilize `go:embed` to package the templates into the binary.
- Refactor the data gathering logic in `main.go` to construct structured contexts (`ThemeFileContext`, `SassBridgeContext`) that are passed to the templates.

## Capabilities

### New Capabilities
- `template-engine`: Utilizing declarative text templates for generating the Sass bridge and Go theme code.

### Modified Capabilities

## Impact

- **Affected code:** `internal/m3generator/cmd/main.go`
- **Dependencies:** Standard library `text/template` and `embed` packages will be utilized. No new external dependencies required.
- **Maintainability:** Significantly improved readability and easier modification of generated code structures.
