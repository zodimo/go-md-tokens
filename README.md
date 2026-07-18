# Go Gen M3

`go-gen-m3` is a code generator that parses Material 3 SCSS token schemas and generates a strongly-typed Go token resolution engine.

## Typed Token System

Instead of relying on fragile string-based map lookups, `go-gen-m3` generates a fully type-safe API for resolving Material 3 component tokens, complete with state fallback semantics.

### Features

- **Type-safe Component Tokens:** Every component (e.g., `FilledButton`) has a strongly-typed struct.
- **State Resolution Engine:** Automatically handles state-specific overlays (e.g., Hover, Focus) and compound states (e.g., `StateSelected | StateHover`).
- **Base Fallbacks:** Gracefully falls back to base token values when state-specific tokens are nil or unavailable.
- **Custom Components:** Supports registering your own custom components with the same resolver API.
- **FlatMap Compatibility:** Provides a `FlatMap()` method to export resolved tokens back to a standard string map if needed by legacy systems.

### Usage

```go
package main

import (
	"fmt"
	"github.com/zodimo/go-gen-m3/pkg/m3tokens"
)

func main() {
	// Initialize a theme
	theme := m3tokens.NewTheme("light")

	// Resolve a token with a specific state
	color := theme.Resolver.FilledButton().ContainerColor(m3tokens.StateHover)
	fmt.Println("Hover Container Color:", color)

	// Resolve compound states
	compoundColor := theme.Resolver.Checkbox().StateLayerColor(m3tokens.StateSelected | m3tokens.StateHover)
	fmt.Println("Selected & Hover State Layer:", compoundColor)

	// Register a custom component
	theme.RegisterCustomComponent("my-widget", map[string]string{
		"container-color": "custom-color",
	})
	
	customColor := theme.Resolver.Custom("my-widget", "container-color", m3tokens.StateDefault)
	fmt.Println("Custom Widget Color:", customColor)
}
```

# Links
- https://m3.material.io/foundations/design-tokens/overview
- https://styledictionary.com/
- https://tr.designtokens.org/format/
- https://www.npmjs.com/package/wireit
- https://www.npmjs.com/package/lit
- https://github.com/tdewolff/parse
