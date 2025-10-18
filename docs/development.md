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

With direnv enabled, you can run scripts directly:

```bash
test          # Run tests
test-race     # Run tests with race detection
lint          # Run linter
coverage      # Check test coverage
bench         # Run benchmarks
example       # Run example
```

Or use `devenv shell <script>`:

```bash
devenv shell test
devenv shell test-race
devenv shell lint
```

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [devenv Documentation](https://devenv.sh/)
- [reviewdog Documentation](https://github.com/reviewdog/reviewdog)
