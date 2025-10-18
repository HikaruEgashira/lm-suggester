# Changelog

All notable changes to lm-suggester will be documented in this file.

## [v0.6.0] - 2025-10-18

![v0.6.0 Demo](examples/mcp/mcp.gif)

### ✨ New Features
- MCP Server Support: Add Model Context Protocol server subcommand ([5dfa50c](https://github.com/HikaruEgashira/lm-suggester/commit/5dfa50c2f0b8c0ec0c7ce8e8bdd8c36b776ae1aa))
  - Integrate with MCP Inspector for interactive suggestion conversion
  - Implement JSONRPC protocol for suggest tool
  - Auto-exit on stdin close for better process management
- Development Environment Migration: Migrate from Nix Flakes to devenv ([fa17c14](https://github.com/HikaruEgashira/lm-suggester/commit/fa17c14d79ed7ad29e47f3b71d0b29b6f41f65ab))
  - Improved developer experience with devenv shell scripts
  - Support for test, lint, coverage, bench, and example commands

### 🔧 Others
- SBOM Generation: Add Software Bill of Materials generation ([0c0ac3f](https://github.com/HikaruEgashira/lm-suggester/commit/0c0ac3fe7e30ea48dc42e8aefe7eb8e1e5f29c9a))
- Documentation updates and cleanups ([c494850](https://github.com/HikaruEgashira/lm-suggester/commit/c49485077cfd1f3a2e9af69a05cdb4c9f57e6a21), [7f24081](https://github.com/HikaruEgashira/lm-suggester/commit/7f24081c82bc63e0c1dc7c7b7a2c7f3a11e8fe96))

---

## [v0.5.1] - 2025-09-21

![v0.5.1 Demo](examples/e2e/v0.5.1_unified_api.gif)

### ✨ New Features
- Unified Convert API with Auto-detection: FilePath is now required while BaseText is optional ([004250c](https://github.com/HikaruEgashira/lm-suggester/commit/004250ce26dbf4a70782cb17a7f206ff123731eb))
- Single function API: Refactored to use a single `Convert` function with automatic format detection ([6d3d078](https://github.com/HikaruEgashira/lm-suggester/commit/6d3d0789c2e3365c64fc61f903d30f013170558d))

### 🐛 Bug Fixes
- Rename JSONL tests ([7a896d2](https://github.com/HikaruEgashira/lm-suggester/commit/7a896d28b64415fcebb29e547ba7485fc8b8d119))
- Support lowercase 'message' field in passthrough conversion ([6ac45a6](https://github.com/HikaruEgashira/lm-suggester/commit/6ac45a63a57083c7f09dad236e0815381ec9f609))

### 🔧 Others
- Cleanup and refactoring ([af882d1](https://github.com/HikaruEgashira/lm-suggester/commit/af882d163e93c14f29a48e0d478721871a918fe9))

---

## [v0.4.0] - 2025-09-21

![v0.4.0 Demo](examples/e2e/v0.4.0_jsonl_support.gif)

### ✨ New Features
- JSONL Format Support: Add support for processing multiple suggestions via JSONL format ([b9e1829](https://github.com/HikaruEgashira/lm-suggester/commit/b9e1829682f02670f5afdadbc7f2297b7ff878ce))
  - Process line-by-line JSON input
  - Handle multiple conversion results in a single operation
  - Support both standard and pretty-print output

### 🔧 Others
- Documentation updates ([75eae4d](https://github.com/HikaruEgashira/lm-suggester/commit/75eae4d651a7c7da8d5677dde8ed698c86bd763d))

---

## [v0.3.1] - 2025-09-21

### 🐛 Bug Fixes
- Configure goreleaser for separate CLI module build ([782d6dd](https://github.com/HikaruEgashira/lm-suggester/commit/782d6dd5141f6885439f3240f2d6e66f5aa410e6))

---

## [v0.2.4] - 2025-09-21

### 🐛 Bug Fixes
- Add source archives to goreleaser config ([8ce1450](https://github.com/HikaruEgashira/lm-suggester/commit/8ce1450cf274b0d76b375c34b764757bd03fefd0))

---

## [v0.2.3] - 2025-09-21

![v0.2.3 Demo](examples/e2e/v0.2.3_passthrough.gif)

### ✨ New Features
- Generic JSON Transformation System: Add support for multiple output formats ([f54c992](https://github.com/HikaruEgashira/lm-suggester/commit/f54c992183afb0b62f2db0047154383dac0b9b7f))
- Pure Passthrough JSON Transformation: Implement passthrough transformation that preserves custom fields ([f93da78](https://github.com/HikaruEgashira/lm-suggester/commit/f93da78de21fc75746e1721b3e4df7cbbe2771a0))

### 🐛 Bug Fixes
- Address code review comments for passthrough implementation ([95ad312](https://github.com/HikaruEgashira/lm-suggester/commit/95ad312fb910496234569f49a6f382c25f29206c))
- Remove binary builds from goreleaser config for library-only project ([a12dc67](https://github.com/HikaruEgashira/lm-suggester/commit/a12dc67187d715e1b115ec85f0f1e4980409895a))

### 🔧 Others
- Create SECURITY.md for security policy ([22ae54f](https://github.com/HikaruEgashira/lm-suggester/commit/22ae54f0517b1f7e26ba06c14d3a193b5d02a571))
- Refactor to pure passthrough architecture ([42d9d4e](https://github.com/HikaruEgashira/lm-suggester/commit/42d9d4e8059692f75246e68167a4aac44f0e8a0f))
- Use passthrough API everywhere ([7b76ce0](https://github.com/HikaruEgashira/lm-suggester/commit/7b76ce0c541179f41358d76ec456fa3ed0517fdc))

---

## Installation

### Go Module
```bash
go get github.com/HikaruEgashira/lm-suggester@latest
```

### CLI Tool
```bash
go install github.com/HikaruEgashira/lm-suggester/cmd/lm-suggester@latest
```

---

## Legend

- ✨ New Features
- 🐛 Bug Fixes
- 🔧 Other Changes (refactoring, documentation, etc.)
