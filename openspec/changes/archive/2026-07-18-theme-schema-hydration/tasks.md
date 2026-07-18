## 1. Update Generator Logic

- [x] 1.1 Modify `generateTheme` function signature string in `internal/m3generator/generator/generator.go` to generate `func NewTheme(tokens map[string]string) (*Theme, error)`.
- [x] 1.2 Implement hydration code generation in `generateTheme`: Iterate over `schema.Components` and `SupportedTokens` to generate explicit assignment blocks mapping the flat string key to the component struct field.
- [x] 1.3 Add validation code generation: The generated code should track missing keys and return an error if the map is missing any expected token.

## 2. Regenerate Code

- [x] 2.1 Run `go run internal/m3generator/cmd/main.go` to generate the updated token structures and functions.
- [x] 2.2 Verify that the generated `pkg/m3tokens/theme.go` file contains the explicit O(1) assignments.

## 3. Update Callers

- [x] 3.1 Update tests or any initialization code calling `m3tokens.NewTheme` to pass a token map (e.g., `theme.GetM3LightTokens()`) and properly handle the returned error.
