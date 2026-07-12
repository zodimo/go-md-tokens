## ADDED Requirements

### Requirement: Generate component enum types
The code generator SHALL produce a Go file containing enum types for `Component`, `TokenCategory`, `TokenProperty`, and `TokenState` derived from the union of all component schemas.

#### Scenario: Component enum includes all components
- **WHEN** the generator processes the schema for 49 components
- **THEN** it SHALL produce a `Component` enum with values including `ComponentFilledButton`, `ComponentElevatedCard`, `ComponentSwitch`, `ComponentFilledTextField`
- **AND** each component SHALL have a generated `String()` method returning the kebab-case name

### Requirement: Generate base token structs
For every component, the generator SHALL produce a struct representing the base (non-state) tokens, with fields grouped by category.

#### Scenario: FilledButton base struct
- **WHEN** the generator processes the `filled-button` schema
- **THEN** it SHALL produce a `FilledButtonTokens` struct containing:
  - `Container ContainerTokens`
  - `Icon IconTokens`
  - `LabelText TextTokens`
  - `Leading SpaceTokens`
  - `Trailing SpaceTokens`

### Requirement: Generate state overlay structs
For every component, the generator SHALL produce overlay structs for each state that has at least one state-specific token, using `*string` fields for tokens that exist and omitting fields for tokens that do not.

#### Scenario: FilledButton hover overlay
- **WHEN** the generator processes the `filled-button` schema
- **THEN** it SHALL produce a `FilledButtonHoverOverlay` struct containing:
  - `ContainerElevation *string`
  - `IconColor *string`
  - `LabelTextColor *string`
  - `StateLayer *StateLayerTokens`
- **AND** the struct SHALL NOT contain fields for tokens without a hover variant

### Requirement: Generate resolver methods
The generator SHALL produce a `Resolver` type with methods for each component that accept a `TokenState` parameter and return the resolved string value, implementing base fallback semantics.

#### Scenario: State-aware resolution
- **WHEN** the generator produces the resolver for `FilledButton`
- **THEN** the method `LabelTextColor(state TokenState) string` SHALL:
  - Return `FilledButtonTokens.Hover.LabelTextColor` when state is `StateHover` and the overlay is non-nil
  - Return `FilledButtonTokens.Disabled.LabelTextColor` when state is `StateDisabled` and the overlay is non-nil
  - Fall back to `FilledButtonTokens.LabelText.Color` for all other cases

### Requirement: Generate system token maps
The generator SHALL produce static Go maps for system modules (color, elevation, shape, state, typescale) containing all known token names and their default values.

#### Scenario: System elevation map
- **WHEN** the generator processes the `md-sys-elevation` schema
- **THEN** it SHALL produce a map `sysElevation = map[string]string{"level0": "0", "level1": "1", ...}`

### Requirement: Generate a theme constructor
The generator SHALL produce a `Theme` struct and a `NewTheme(mode string)` constructor that initializes the resolver with the appropriate system token maps for the requested theme mode.

#### Scenario: Light theme initialization
- **WHEN** a consumer calls `m3.NewTheme(m3.Light)`
- **THEN** the returned `Theme` SHALL have `SysColor` populated from the light scheme in `m3/material-theme.json`
- **AND** `SysElevation` SHALL contain the static elevation values
