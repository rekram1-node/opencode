# OpenCode Development Guide

## Commands
- Build: `go build`
- Run: `go run main.go`
- Test: `go test ./...` (all tests) or `go test ./internal/llm/tools/...` (specific package)
- Single test: `go test -run TestName ./path/to/package`
- Format: `go fmt ./...`
- Lint: `go vet ./...`
- Generate SQL code: `sqlc generate`
- Build snapshot: `./scripts/snapshot`
- Release: `./scripts/release` or `./scripts/release --minor`

## Code Style
- **Naming**: PascalCase for exported, camelCase for unexported
- **Constructors**: Use `New` prefix (e.g., `NewService`)
- **Errors**: Use `fmt.Errorf("context: %w", err)` for wrapping
- **Imports**: Group standard library, external, and internal imports
- **Testing**: Use testify for assertions (`assert`, `require`)
- **Comments**: Document exported functions, types, and constants
- **Error handling**: Check errors immediately, avoid nested if statements
- **File organization**: One primary type per file, related helpers together
- **Interfaces**: Define interfaces where they're used, not implemented