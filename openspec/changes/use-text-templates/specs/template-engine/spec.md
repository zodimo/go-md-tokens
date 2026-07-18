## ADDED Requirements

### Requirement: Template-driven generation
The system SHALL use external text templates rather than imperative string building to generate both the intermediate Sass bridge code and the final Go theme code.

#### Scenario: Code structure generation
- **WHEN** the generator parses SCSS tokens and custom colors
- **THEN** it generates the required output code using `text/template` by passing structured context structs to the embedded templates

### Requirement: Embedded templates
The system SHALL embed the templates into the generator binary to maintain a single distributable artifact.

#### Scenario: Running the generator
- **WHEN** the user executes the compiled `go-gen-m3` binary
- **THEN** the generator executes without requiring the external `.tmpl` files on the filesystem at runtime
