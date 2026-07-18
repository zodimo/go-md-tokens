## Why

The current M3 token pipeline generates flat `map[string]string` with dot-separated keys like `md.comp.filled.button.disabled.label.text.opacity`. This stringly-typed representation has zero compile-time safety, no IDE autocomplete, and makes state-dependent resolution manual and error-prone. Developers building Material 3 components in Go need a typed token system that mirrors the hierarchical structure defined in the SCSS source, supports state overlays, and extends cleanly to custom components.

## What Changes

- **Replace flat token maps with typed Go structs**: Component tokens decomposed into categories (Container, Icon, LabelText), with state overlays (Hover, Focus, Disabled, Pressed) modeled as `*string` fields that override base values.
- **Schema extractor**: Parse the 49 `m3web/v2.4.1/tokens/_md-comp-*.scss` files to extract `$supported-tokens` and generate the type system automatically.
- **State-aware resolver**: Runtime resolution of tokens by `TokenState` with automatic fallback to base values when state overlays are nil.
- **Cross-component aliasing**: Custom components can reference resolved tokens from other components (e.g., a DataTable using Checkbox tokens) through the same typed resolver.
- **Code generator**: `main.go` becomes the generator that produces `pkg/m3tokens/` with enums, structs, and resolver methods.

## Capabilities

### New Capabilities
- `schema-extraction`: Parse SCSS token files to extract component schema, supported tokens, state variants, and dependency graph.
- `token-codegen`: Generate typed Go structs, state overlay types, and resolver methods from the extracted schema.
- `state-resolution`: Resolve component tokens by state (Default, Hover, Focus, Pressed, Disabled, Dragged, Selected, Error, Active) with base fallback semantics.
- `custom-component-extension`: Allow user-defined components to declare their own `$supported-tokens` and participate in the same typed generation pipeline.

### Modified Capabilities
<!-- No existing capabilities to modify -->

## Impact

- **API surface**: `theme/tokens.go` (flat maps) will be replaced by `pkg/m3tokens/` (typed system). Existing consumers of the flat map can continue using a compatibility layer.
- **Build process**: `go run .` will now generate Go source files instead of just a single token map file.
- **Dependencies**: No new runtime dependencies; code generation uses the existing `godartsass/v2` transpiler for schema extraction.
