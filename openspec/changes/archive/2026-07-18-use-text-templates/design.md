## Context

Currently, the M3 token code generation (both Go code for themes and intermediate Sass bridge code) is constructed imperatively using `strings.Builder`. This tightly couples data gathering with code formatting, making the generator difficult to read and modify. The proposal introduces Go's `text/template` engine to separate these concerns.

## Goals / Non-Goals

**Goals:**
- Separate the code generation logic from the data extraction logic.
- Use `text/template` to generate the final Go code and the intermediate Sass bridge.
- Package the templates within the binary using `go:embed`.
- Ensure the output of the generator remains exactly the same as before.

**Non-Goals:**
- Changing the structure of the generated Go output.
- Modifying the underlying godartsass integration or JSON schema validation logic.

## Decisions

**1. Template File Location and Embedding**
Templates will be stored in `internal/m3generator/templates/` as `.tmpl` files (e.g., `go_theme.tmpl`, `sass_bridge.tmpl`). A `templates.go` file in this package will use `//go:embed *.tmpl` to load these templates at compile time and expose parsed `*template.Template` instances to `cmd/main.go`.

*Rationale:* External `.tmpl` files allow IDEs to provide HTML/Go Template syntax highlighting. `go:embed` ensures the CLI remains a single, easily distributable binary.

**2. Data Context Structs**
`cmd/main.go` will build strongly-typed struct contexts to pass into `Template.Execute`.

*   `ThemeFileContext`: Represents the entire output Go file, containing a slice of `ThemeModeContext` (for light, dark, etc.).
*   `ThemeModeContext`: Contains the function name (`Light`, `Dark`), mode name, and the map of sorted tokens.
*   `SassBridgeContext`: Contains the variables needed for the Sass execution block (`FileURL`, `CustomColorsSass`, `Component`).

*Rationale:* This explicitly defines the API boundary between the Go logic and the templates. It eliminates the need for inline string formatting and logic within the generator loop.

## Risks / Trade-offs

- **[Risk]** Whitespace drift in templates causing differences in the generated Go output.
  → *Mitigation*: The generator already runs `gofmt` on the output directory (line 260 of `main.go`). This will normalize any minor whitespace issues from the template execution.
- **[Risk]** Slightly slower execution time due to `text/template` parsing/execution overhead.
  → *Mitigation*: The overhead is negligible for a CLI tool that runs infrequently and processes a finite set of tokens. Maintainability outweighs raw speed here.
