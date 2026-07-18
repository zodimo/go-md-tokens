# Custom Component Extension

## Purpose
TBD

## Requirements

### Requirement: Parse custom component token files
The extractor SHALL support an additional directory (e.g., `custom/tokens/`) containing user-defined `_md-comp-*.scss` files with the same structure as M3 components.

#### Scenario: Custom data-table component
- **WHEN** the extractor finds `custom/tokens/_md-comp-data-table.scss`
- **THEN** it SHALL parse `$supported-tokens` and produce a schema for `data-table`
- **AND** it SHALL be indistinguishable from a Google M3 component in the generated output

### Requirement: Resolve cross-component aliases
Custom components SHALL be able to reference tokens from other components by calling their resolver methods during code generation.

#### Scenario: DataTable referencing Checkbox
- **WHEN** the generator processes a custom component `data-table` with a token `checkbox-selected-container-color`
- **THEN** the generated `DataTableTokens.Checkbox` field SHALL be typed as `CheckboxTokens`
- **AND** the resolver SHALL delegate to `r.Checkbox().Selected.ContainerColor`

### Requirement: Extend system token maps
Custom components SHALL be able to define new system-level tokens (e.g., `cell-padding-start`) that are not part of the standard M3 system but are shared across multiple custom components.

#### Scenario: Shared custom spacing tokens
- **WHEN** two custom components (`data-table` and `list`) both define `cell-padding-start`
- **THEN** the generator SHALL produce a `CustomSysSpacing` map containing `cell-padding-start`
- **AND** both component resolvers SHALL reference the same `CustomSysSpacing` value

### Requirement: Custom state definitions
Custom components SHALL be able to define their own state vocabularies beyond the standard M3 states.

#### Scenario: DataTable sort state
- **WHEN** a custom component defines tokens with prefix `sorted-` and `unsorted-`
- **THEN** the generator SHALL produce `StateSorted` and `StateUnsorted` constants
- **AND** the component resolver SHALL handle these states with the same fallback semantics as standard states

### Requirement: Backward compatibility for flat maps
The generated code SHALL include a compatibility function that exposes all component tokens as a flat `map[string]string`, mapping dot-separated keys to the typed values.

#### Scenario: Flat map compatibility
- **WHEN** a legacy consumer calls `theme.FlatMap()`
- **THEN** it SHALL return a map containing `"md.comp.filled.button.disabled.label.text.opacity" → "0.38"`
- **AND** the flat map SHALL contain all tokens from all components for the current theme mode
