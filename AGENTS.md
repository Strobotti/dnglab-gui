# Agentic Development Guidelines

This document provides guidelines for AI agents and automated development tools working on this project.

## Branching Strategy

**Always create a feature branch** before making code changes.

Exceptions (may work directly on `main`):

- Documentation-only changes (README.md, docs/, comments)
- Test-only changes (no production code modified)

For all other changes:

```sh
git checkout -b feature/descriptive-name
# or
git checkout -b fix/issue-description
```

## Project Structure

```
cmd/dnglab-gui/       # Application entry point
internal/
  assets/             # Bundled resources (icon, etc.)
  config/             # Settings management
  converter/          # dnglab wrapper and file scanning
  ui/                 # Fyne UI components
assets/               # Source assets (icon.png)
docs/                 # Documentation
```

## Build and Test

```sh
# Build
go build -ldflags "-s -w" -o bin/dnglab-gui ./cmd/dnglab-gui
# or
make build

# Run tests
go test ./...

# Check for issues
go vet ./...
```

## Code Style

- Follow standard Go conventions (gofmt, go vet)
- Tests use `_test.go` suffix alongside source files
- Internal packages go in `internal/` to prevent external imports
- Use descriptive variable and function names

## Dependencies

- **Fyne** - UI toolkit (requires CGO and C compiler)
- **dnglab** - External CLI tool for RAW-to-DNG conversion

When adding dependencies:

- Prefer well-maintained, widely-used packages
- Pin exact versions in go.mod
- Document any new system requirements

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

| Prefix      | Purpose                  | Version bump |
|-------------|--------------------------|--------------|
| `feat:`     | New feature              | Minor        |
| `fix:`      | Bug fix                  | Patch        |
| `docs:`     | Documentation only       | Patch        |
| `chore:`    | Maintenance (no release) | None         |
| `refactor:` | Code refactoring         | Patch        |

## Environment Requirements

- Go 1.27.1+
- CGO enabled (default)
- C compiler:
    - Linux: `gcc` + `libgl1-mesa-dev xorg-dev`
    - macOS: Xcode Command Line Tools

## Testing Guidelines

- Write tests for new functionality in `internal/` packages
- Use table-driven tests where appropriate
- Mock external dependencies (filesystem, dnglab CLI)
- Run `go test ./...` before committing

## UI Changes

The application uses Fyne UI toolkit. When modifying UI:

- Keep the four-section layout (source, destination, options, metadata)
- Test with `go run ./cmd/dnglab-gui` to verify visual changes
- Ensure accessibility (readable labels, logical tab order)

## Bundled Assets

To update bundled assets (e.g., application icon):

```sh
fyne bundle -o internal/assets/bundled.go assets/icon.png
```
