# State Resolution

## Purpose
TBD

## Requirements

### Requirement: Define all token states
The system SHALL define a `TokenState` type with at minimum the following states: `Default`, `Hover`, `Focus`, `Pressed`, `Disabled`, `Dragged`, `Selected`, `Error`, `Active`.

#### Scenario: State enumeration completeness
- **WHEN** a developer inspects the generated `TokenState` type
- **THEN** it SHALL contain all listed states as distinct constants
- **AND** it SHALL be usable as a bitmask for compound state combinations

### Requirement: Base fallback semantics
For any token queried with a non-default state, the resolver SHALL first check the state overlay for that token; if the overlay is nil or the token field is nil, it SHALL fall back to the base token value.

#### Scenario: Hover token with overlay
- **WHEN** a consumer requests `FilledButton.ContainerElevation(StateHover)`
- **AND** the hover overlay exists and contains `ContainerElevation = "1"`
- **THEN** the resolver SHALL return `"1"`

#### Scenario: Hover token without overlay
- **WHEN** a consumer requests `FilledButton.ContainerColor(StateHover)`
- **AND** the hover overlay is nil or does not contain `ContainerColor`
- **THEN** the resolver SHALL return the base `Container.Color`

### Requirement: Compound state resolution
The resolver SHALL support compound states (e.g., `Selected | Hover`) by checking single-state overlays in a defined priority order.

#### Scenario: Selected + Hover compound state
- **WHEN** a consumer requests `DataTable.RowColor(StateSelected|StateHover)`
- **AND** a `SelectedHover` overlay exists with `RowColor = "#6750A4"`
- **THEN** the resolver SHALL return `"#6750A4"`
- **AND** if no compound overlay exists, it SHALL check `Selected` then `Hover` then base

### Requirement: Nil-safe overlay access
All generated resolver methods SHALL be nil-safe: if the state overlay struct pointer is nil, resolution SHALL immediately fall back to the base value without panicking.

#### Scenario: Nil overlay for disabled state
- **WHEN** a consumer requests `ElevatedCard.ContainerColor(StateDisabled)`
- **AND** the `ElevatedCard` schema has no disabled tokens (overlay struct is nil)
- **THEN** the resolver SHALL return the base `Container.Color` without panic
