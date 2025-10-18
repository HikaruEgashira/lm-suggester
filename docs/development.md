# Development Guide

This guide covers development workflows for lm-suggester, including Nix-based development environment setup.

## Development Environment

This project uses Nix flakes for reproducible development environments.

- [Nix](https://nixos.org/download.html) with flakes enabled
- [direnv](https://direnv.net/) (optional but recommended)

#### Setup

1. With direnv (automatic shell activation):

```bash
# Allow direnv in this directory
direnv allow

# The development shell will automatically activate when you cd into the directory
```

2. Without direnv (manual activation):

```bash
# Enter the Nix development shell
nix develop

# Your shell now has all development tools available
```

### Build Issues

Problem: `GOPATH` or `GOROOT` conflicts

Solution: Ensure you're in the Nix shell, which sets these correctly:
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
