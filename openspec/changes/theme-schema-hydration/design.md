## Context

Currently, the `m3-token-gen` code generator extracts token schema (Components and SupportedTokens) from Material Design SASS files. It generates typed Go structs for components (e.g. `FilledButtonTokens`) but leaves them empty during `NewTheme` initialization. In a separate pipeline step, it evaluates the exact hex values and exports them as flat maps (e.g., `GetM3LightTokens()`). Bridging the flat map to the structured API is necessary.

## Goals / Non-Goals

**Goals:**
- Provide a fully hydrated `Theme` struct from `NewTheme` so users can access tokens via structured API without nil pointers.
- Ensure theme instantiation is highly performant.
- Leverage the existing schema data (`schema.json`) available during code generation.

**Non-Goals:**
- Removing the flat map export (`theme/tokens.go`).
- Modifying how SASS tokens are parsed (this relies on the existing extractor).

## Decisions

**1. Generate explicit assignments in NewTheme**
We will update `generator.go` to iterate through the parsed `schema.Components` and output explicit Go assignment statements mapping the string key to the struct field path.
*Alternative considered:* Use Go reflection (`reflect` package) at runtime to iterate over the struct fields and look up the `m3` tags. Rejected due to reflection overhead and heap allocations. Generating code moves this cost to compile time.

**2. Change NewTheme signature to include error return**
We will change `func NewTheme(mode string) *Theme` to `func NewTheme(tokens map[string]string) (*Theme, error)`.
*Rationale:* This decouples the `m3tokens` package from knowing about specific modes (light/dark/custom). It allows consumers to inject any token map, including custom dynamically generated themes. Returning an error ensures we can validate completeness and block the creation of a fundamentally broken theme.

**3. Validate Map Completeness**
The generated code will track whether any expected token key was missing from the input map. If missing keys are found, it returns an aggregated error outlining which tokens were missing.
*Rationale:* Material Design relies heavily on all tokens being present. Silently allowing nil or empty tokens leads to subtle visual bugs or runtime panics. Failing fast at initialization is the safest approach.

## Risks / Trade-offs

- **[Risk]** Massive file size: Generating thousands of assignments will drastically increase the size of `theme.go`.
  - **Mitigation:** Go compiler is highly optimized for flat assignment blocks. The source file may be large, but the binary size increase is negligible and runtime performance is optimal.
