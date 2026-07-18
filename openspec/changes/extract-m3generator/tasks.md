## 1. Setup

- [ ] 1.1 Create `internal/m3generator` directories (`cmd`, `generator`, `schema`)
- [ ] 1.2 Move `schema-def.json` to `internal/m3generator/schema/`
- [ ] 1.3 Move `extractor.go` and `generator.go` to `internal/m3generator/generator/` and update package references to `generator`

## 2. CLI Implementation

- [ ] 2.1 Create `m3gen.yaml` config struct in `internal/m3generator/cmd/main.go`
- [ ] 2.2 Migrate logic from root `main.go` into `internal/m3generator/cmd/main.go`
- [ ] 2.3 Implement yaml parsing to replace hardcoded strings
- [ ] 2.4 Delete old root `main.go`

## 3. Tooling and Hookup

- [ ] 3.1 Create root `generate.go` file with `//go:generate go run ./internal/m3generator/cmd/main.go`
- [ ] 3.2 Create sample `m3gen.yaml` config in root directory
- [ ] 3.3 Test generation pipeline (`go generate ./...`) and verify parity with old output
