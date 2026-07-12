## ADDED Requirements

### Requirement: Parse SCSS component token files
The schema extractor SHALL parse each `m3web/v2.4.1/tokens/_md-comp-*.scss` file to extract the component name, `$supported-tokens`, `$unsupported-tokens`, `$renamed-tokens`, system dependencies (`$_default`), and any state-specific token prefixes.

#### Scenario: Extracting filled-button schema
- **WHEN** the extractor reads `_md-comp-filled-button.scss`
- **THEN** it SHALL produce a schema containing:
  - component name: `filled-button`
  - supported tokens list including `container-color`, `disabled-label-text-opacity`, `hover-state-layer-opacity`
  - dependencies: `md-sys-color`, `md-sys-elevation`, `md-sys-shape`, `md-sys-state`, `md-sys-typescale`
  - state prefixes identified: `disabled-`, `hover-`, `focus-`, `pressed-`

### Requirement: Identify token hierarchy from names
The extractor SHALL decompose every kebab-case token name into a hierarchy of State, Category, and Property by analyzing the prefix patterns against a known state vocabulary.

#### Scenario: Decomposing disabled-label-text-opacity
- **WHEN** the extractor processes the token `disabled-label-text-opacity`
- **THEN** it SHALL classify:
  - State: `disabled`
  - Category: `label-text`
  - Property: `opacity`

#### Scenario: Decomposing container-shape-end-end
- **WHEN** the extractor processes the token `container-shape-end-end`
- **THEN** it SHALL classify:
  - State: `default` (no state prefix)
  - Category: `container`
  - Sub-category: `shape`
  - Property: `end-end`

### Requirement: Extract system module schemas
The extractor SHALL parse the `md-sys-*` module files to produce a static map of system token names and their resolution rules (e.g., `md-sys-elevation` levels, `md-sys-shape` corners, `md-sys-state` opacities).

#### Scenario: Extracting elevation system tokens
- **WHEN** the extractor reads `md-sys-elevation` files
- **THEN** it SHALL produce a map containing `level0` → `0`, `level1` → `1`, `level2` → `2`, `level3` → `3`, `level4` → `4`, `level5` → `5`

### Requirement: Output a consolidated JSON schema
The extractor SHALL write a consolidated schema file (`schema.json`) containing all component schemas, system module schemas, and the dependency graph, to be consumed by the code generator.

#### Scenario: Valid schema generation
- **WHEN** the extractor runs against the full `m3web/v2.4.1/tokens/` directory
- **THEN** it SHALL produce a valid `schema.json` containing schemas for all 49 components and all 6 system modules
- **AND** the JSON SHALL validate against a predefined JSON Schema for the extractor output
