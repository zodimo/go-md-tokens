## Why

The current token code generator creates component token structs and a `NewTheme` function that does not actually initialize the component values. Meanwhile, it generates a perfectly flat, pre-evaluated map of tokens (like `GetM3LightTokens()`). To bridge the gap, we need to populate the structured fields from the map. Doing this with reflection at runtime is slow and generates garbage. Since we already parse the `schema.json` during code generation, the optimal approach is to use the schema to generate explicit, statically-typed hydration assignments inside `NewTheme` directly from the flat token map.

## What Changes

- Modify `NewTheme` signature in `m3tokens` package to accept `tokens map[string]string` and return `(*Theme, error)` instead of `mode string` (**BREAKING**).
- Add validation to `NewTheme` to return an error and block theme creation if any required token from the schema is missing in the provided map.
- Update `generateTheme` in `internal/m3generator/generator/generator.go` to iterate over `schema.Components` and `SupportedTokens`.
- During code generation, write explicit O(1) assignments for every token mapping the flat string key (e.g. `"md.comp.filled.button.container.color"`) to its struct field (e.g. `theme.FilledButtonTokens.Container.Color`).
- Update `main.go` and any dependent code to instantiate `NewTheme` by passing a token map (e.g., `theme.GetM3LightTokens()`).

## Capabilities

### New Capabilities
- `theme-hydration`: Generating statically typed struct hydration from the flat token map using the token schema.

### Modified Capabilities

## Impact

- **API Break**: `m3tokens.NewTheme(mode string)` becomes `m3tokens.NewTheme(tokens map[string]string) (*Theme, error)`.
- **Safety**: Prevents runtime panics or visual bugs by strictly validating that all required tokens are present during initialization.
- **Performance**: Near-instant O(N) theme instantiation with no runtime reflection or heap allocations.
- **Code Size**: The generated `theme.go` will grow significantly due to thousands of explicit assignment blocks.
