## Context

The current `go-gen-m3` pipeline uses `godartsass/v2` to compile Sass bridge code and extract token values as flat CSS properties. The output is a single `theme/tokens.go` with `map[string]string` for each theme mode. While functional, this approach loses all structural information encoded in the SCSS token names (state, category, property hierarchy) and requires consumers to use stringly-typed keys with no compile-time validation.

This change introduces a typed token system that:
1. Parses the SCSS schema to understand token hierarchy
2. Generates Go structs, enums, and resolvers from that schema
3. Supports state-aware resolution with base fallback
4. Extends cleanly to custom components

## Goals / Non-Goals

**Goals:**
- Generate typed Go code from the 49 M3 component SCSS files
- Decompose flat token names into structured categories, properties, and states
- Provide state-aware token resolution with automatic base fallback
- Support custom components with the same typed generation pipeline
- Maintain backward compatibility via a flat-map accessor

**Non-Goals:**
- Full SCSS parser or general Sass evaluator in Go
- Runtime theme switching (out of scope; theme is resolved at initialization)
- CSS custom property generation (we resolve values, not output CSS)
- Support for `@mixin`, `@include`, `@media`, `@extend` (token files don't use these)

## Decisions

### 1. Schema extraction via Sass compilation (not parsing)
**Decision:** We will continue using `godartsass/v2` to compile a bridge that dumps token metadata, rather than writing a full SCSS parser.
**Rationale:** The M3 token subset is small and stable. Using the official Sass engine guarantees correctness against upstream changes. A custom parser would require handling Sass maps, functions, `@use`, and interpolation — a multi-thousand-line project.
**Alternative considered:** Hand-written SCSS tokenizer + AST. Rejected due to implementation cost and maintenance burden.

### 2. Code generation produces `pkg/m3tokens/` package
**Decision:** Generated code lives in a dedicated package, not inline in `main.go` or `theme/`.
**Rationale:** Separation of concerns: the generator is a build tool, the generated code is a library. Consumers import `pkg/m3tokens` without needing the generator at runtime.
**Alternative considered:** Generate into `theme/tokens.go` directly. Rejected because the flat map format is fundamentally different from typed structs and mixing both would be confusing.

### 3. State overlays modeled as nilable struct pointers
**Decision:** State-specific tokens are stored in overlay structs with `*string` fields. The overlay struct itself is a pointer on the base struct.
**Rationale:** 
- `nil` overlay = component doesn't support this state (e.g., `ElevatedCard` has no `Disabled` tokens)
- `nil` field = state doesn't override this specific token (e.g., hover doesn't change `ContainerColor`)
- This is identical to CSS custom property semantics: if not defined, inherit from base.
**Alternative considered:** Use zero-value strings (`""`) and special markers. Rejected because it can't distinguish "not defined" from "intentionally empty."

### 4. Component accessors are methods on `Theme`, not free functions
**Decision:** Usage pattern is `theme.FilledButton().LabelTextColor(StateHover)`, not `m3.GetFilledButtonLabelTextColor(theme, StateHover)`.
**Rationale:** Fluent API chains naturally, is discoverable via IDE autocomplete, and groups component-related methods under the component accessor. This mirrors how Material Web's Sass API works (`md-comp-filled-button.values($deps)`).

### 5. Flat map backward compatibility via reflection
**Decision:** The `Theme.FlatMap()` method uses struct tag reflection to produce the legacy `map[string]string` format.
**Rationale:** Allows gradual migration. Reflection cost is acceptable because `FlatMap()` is called once at initialization or for serialization, not in hot paths.
**Alternative considered:** Generate both typed and flat accessors side-by-side. Rejected because it doubles the generated code size for limited benefit.

## Risks / Trade-offs

- **[Risk] Schema drift when Google updates M3 tokens** → **Mitigation:** The generator is re-runnable. Update `m3web/` and re-run `go generate`. The typed system will catch missing tokens at compile time.
- **[Risk] Generated code size for 49 components** → **Mitigation:** Code is generated and committed; consumers compile only the Go output, not the generator. Estimated ~3,000 lines of generated Go.
- **[Risk] Nil pointer panic in consumer code** → **Mitigation:** All resolver methods are nil-safe. The only nil risk is direct struct field access (e.g., `tokens.Hover.LabelTextColor` without checking `Hover != nil`), which is a standard Go pattern.
- **[Trade-off] Custom components require SCSS syntax** → Custom components must write token definitions in the same Sass dialect as M3. This is intentional: it enforces consistency with the upstream design system.

## Migration Plan

1. **Phase 1 (This change):** Build the generator, produce `pkg/m3tokens/`, verify output matches `theme/tokens.go` for all 6 modes.
2. **Phase 2:** Update component builder code to use typed accessors (`theme.FilledButton()` instead of flat map lookups).
3. **Phase 3:** Deprecate `theme/tokens.go` once all consumers migrate. Flat map compatibility layer remains indefinitely.

## Open Questions

1. Should the generator also produce JSON Schema for the extracted schema format, to enable validation?
2. How should compound states beyond `Selected|Hover` be represented? Bitmask? Slice of states?
3. Should custom components support `$exclude-hardcoded-values` semantics, or always resolve hardcoded values?
