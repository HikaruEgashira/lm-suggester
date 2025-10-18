# Development Guide

This guide covers development workflows for lm-suggester, including Nix-based development environment setup.

## Development Environment

### Using Nix Flake (Recommended)

This project uses Nix flakes for reproducible development environments.

#### Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled
- [direnv](https://direnv.net/) (optional but recommended)

#### Setup

1. **With direnv (automatic shell activation)**:

```bash
# Allow direnv in this directory
direnv allow

# The development shell will automatically activate when you cd into the directory
```

2. **Without direnv (manual activation)**:

```bash
# Enter the Nix development shell
nix develop

# Your shell now has all development tools available
```

#### Available Tools in Nix Shell

The development environment includes:

- **Go 1.25**: Go compiler and toolchain
- **gotools**: Additional Go tools (goimports, etc.)
- **golangci-lint**: Linter for Go code
- **gopls**: Go language server for IDE support
- **git**: Version control
- **syft**: SBOM generation tool
- **bash & coreutils**: Basic shell utilities

#### Environment Variables

The Nix shell sets up:

- `GOROOT`: Points to the Nix-managed Go installation
- `GOPATH`: Set to `$PWD/.go` (project-local Go cache)
- `PATH`: Contains only Nix-managed binaries (no system interference)

### Traditional Setup (without Nix)

If you prefer not to use Nix:

1. Install [Go 1.25+](https://go.dev/dl/)
2. Install development tools:
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

## Project Structure

### Multi-Module Setup

This project uses a multi-module structure:

- **Root `go.mod`**: Manages library (`suggester/`) dependencies
- **`cmd/lm-suggester/go.mod`**: Manages CLI-specific dependencies (cobra, etc.)

This keeps the library dependencies clean and minimal.

### Core Library

The `suggester/` package provides the core conversion logic:

- **public.go**: Public API (`Input` struct, `BuildRDJSON` function)
- **align.go**: Diff position alignment
- **diff.go**: Minimal diff detection
- **normalize.go**: Text normalization (UTF-8, newlines)
- **errors.go**: Custom error types

## Common Tasks

### Dependency Management

```bash
# Update dependencies for the library
go mod tidy

# Update dependencies for the CLI
cd cmd/lm-suggester
go mod tidy
cd ../..
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./suggester

# Check test coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Current coverage: ~52% (aim to improve with new code)

### Benchmarking

```bash
# Run benchmarks for suggester package
go test -run none -bench . ./suggester

# Run benchmarks with memory profiling
go test -run none -bench . -benchmem ./suggester
```

### Static Analysis

```bash
# Run go vet
go vet ./...

# Run golangci-lint (in Nix shell)
golangci-lint run

# Run goimports
goimports -w .
```

### Building

```bash
# Build the library (check compilation)
go build ./suggester

# Build the CLI tool
cd cmd/lm-suggester
go build -v -o ../../lm-suggester .
cd ../..

# Test the built binary
./lm-suggester version
```

### Running Examples

```bash
# Run simple example
cat _examples/testdata/simple_replacement.json | go run _examples/simple/main.go

# Run with CLI
cat examples/testdata/simple_replacement.json | ./lm-suggester -p
```

## Coding Standards

### Style Guidelines

- **Indentation**: Use tabs (Go standard)
- **Formatting**: Run `gofmt` / `goimports` before committing
- **Naming**:
  - Public symbols: Use descriptive names with `Suggester` prefix where appropriate
  - Package-private: Start with lowercase
- **Documentation**: All public symbols must have doc comments (`// Name ...`)
- **Error messages**: Use `errors.go` types, write in English, keep concise

### Testing Guidelines

- Use table-driven tests (Go idiom)
- Test files: `*_test.go`
- Subtests: Use `t.Run` for grouping
- Bug fix tests: Name as `Test_issue123_description`
- Aim to improve coverage (currently 52.1%)

### Scripts and Tools

- Examples: `_examples/` directory
- Helper scripts: `tools/` directory
- Task definitions: Use `mise tasks edit` for shared tasks

## Commit Guidelines

### Conventional Commits

Follow [Conventional Commits](https://www.conventionalcommits.org/) format:

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Test additions/changes
- `chore:` Maintenance tasks

### Pre-commit Checklist

1. Run tests: `go test ./...`
2. Run linters: `go vet ./...`
3. Format code: `goimports -w .`
4. Keep changes atomic (one logical change per commit)

## Pull Request Process

### Creating a PR

1. **Prepare your changes**:
   ```bash
   # Ensure all tests pass
   go test ./...

   # Check for lint issues
   go vet ./...
   ```

2. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

3. **Create PR**:
   ```bash
   gh pr create --fill
   ```

### PR Content

Include in PR description:

- **Summary**: What changes were made
- **Motivation**: Why these changes are needed
- **Test results**: Output of `go test ./...`
- **Related issues**: `Fixes #123` or `Relates to #456`

### After PR Creation

```bash
# Watch CI checks
gh pr checks --watch

# View PR in browser
gh pr view --web
```

## Troubleshooting

### Nix Shell Issues

**Problem**: `nix develop` fails or hangs

**Solution**:
```bash
# Update flake lock
nix flake update

# Try with verbose output
nix develop --show-trace
```

**Problem**: direnv not loading automatically

**Solution**:
```bash
# Ensure direnv is hooked in your shell
# For bash:
eval "$(direnv hook bash)"

# For zsh:
eval "$(direnv hook zsh)"

# Then allow again
direnv allow
```

### Go Module Issues

**Problem**: Module dependencies not resolving

**Solution**:
```bash
# Clear module cache
go clean -modcache

# Re-download modules
go mod download

# For CLI module
cd cmd/lm-suggester
go mod download
cd ../..
```

### Build Issues

**Problem**: `GOPATH` or `GOROOT` conflicts

**Solution**: Ensure you're in the Nix shell, which sets these correctly:
```bash
nix develop
echo $GOROOT  # Should point to /nix/store/...
echo $GOPATH  # Should be $PWD/.go
```

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Nix Flakes Documentation](https://nixos.wiki/wiki/Flakes)
- [reviewdog Documentation](https://github.com/reviewdog/reviewdog)
