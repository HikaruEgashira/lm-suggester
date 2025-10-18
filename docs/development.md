# Development Guide

## Development Environment

This project uses [devenv](https://devenv.sh/) for development environment management.

### Setup

```bash
# With direnv
direnv allow

# Or manually
devenv shell
```

### Available Commands

Run devenv scripts:

```bash
devenv shell test          # Run tests
devenv shell test-race     # Run tests with race detection
devenv shell lint          # Run linter
devenv shell coverage      # Check test coverage
devenv shell bench         # Run benchmarks
devenv shell example       # Run example
```

Or enter the shell and run commands directly:

```bash
devenv shell
# Now inside devenv shell:
test
test-race
lint
```

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [devenv Documentation](https://devenv.sh/)
- [reviewdog Documentation](https://github.com/reviewdog/reviewdog)
