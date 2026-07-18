## ADDED Requirements

### Requirement: YAML Configuration Support
The system SHALL read a YAML configuration file named `m3gen.yaml` at the root directory instead of using hardcoded variables.

#### Scenario: Using config for paths
- **WHEN** the CLI tool executes
- **THEN** it reads `theme_source`, `sass_tokens_dir`, `schema_output_dir`, and `theme_output` from `m3gen.yaml`.
