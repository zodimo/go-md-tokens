## 1. Setup Templates

- [ ] 1.1 Create `internal/m3generator/templates/` directory
- [ ] 1.2 Create `sass_bridge.tmpl` containing the SCSS interpolation block
- [ ] 1.3 Create `go_theme.tmpl` containing the Go file structure
- [ ] 1.4 Create `templates.go` to use `//go:embed *.tmpl` and export a parsed template object

## 2. Refactor Generator Logic

- [ ] 2.1 Define `SassBridgeContext` struct in `main.go`
- [ ] 2.2 Refactor the Sass bridge generation loop to execute `sass_bridge.tmpl` instead of `fmt.Sprintf`
- [ ] 2.3 Define `ThemeFileContext` and `ThemeModeContext` structs in `main.go`
- [ ] 2.4 Refactor the Go code generation loop to populate a `ThemeFileContext` instead of calling `goCode.WriteString`
- [ ] 2.5 Execute `go_theme.tmpl` with the populated `ThemeFileContext` and write the result to the output file

## 3. Verification

- [ ] 3.1 Run `go run internal/m3generator/cmd/main.go`
- [ ] 3.2 Verify the generated `theme/tokens.go` and intermediate SCSS compile exactly the same as before
