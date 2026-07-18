# Capability: Theme Hydration

## Purpose
TBD

## Requirements

### Requirement: Explicit Theme Hydration
The code generator SHALL use the parsed schema to produce explicit assignment statements inside `NewTheme` for every supported component token.

#### Scenario: Generating NewTheme function
- **WHEN** the generator runs the `generateTheme` function
- **THEN** it SHALL write a mapping assignment from the flat token map key (e.g., `"md.comp.filled.button.container.color"`) to the struct field (e.g., `theme.FilledButtonTokens.Container.Color`).

### Requirement: Token Map Injection
The `NewTheme` initialization function SHALL accept a generic map of tokens rather than a hardcoded mode string.

#### Scenario: Instantiating a Theme
- **WHEN** a developer creates a new theme using `NewTheme`
- **THEN** they MUST pass a pre-evaluated token map (e.g., `theme.GetM3LightTokens()`) as the argument.

### Requirement: Token Map Completeness Validation
The `NewTheme` function SHALL validate that the provided map contains all tokens defined in the schema.

#### Scenario: Missing tokens provided
- **WHEN** a developer passes a token map missing one or more required tokens
- **THEN** `NewTheme` SHALL return a non-nil error listing the missing tokens and return a nil `*Theme`.

#### Scenario: Complete tokens provided
- **WHEN** a developer passes a complete token map
- **THEN** `NewTheme` SHALL return a populated `*Theme` and a nil error.
