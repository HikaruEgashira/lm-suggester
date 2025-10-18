# SBOM (Software Bill of Materials)

This document describes how lm-suggester generates and manages Software Bill of Materials (SBOM) for supply chain security.

## Overview

lm-suggester automatically generates SBOM files to provide transparency about the software components and dependencies used in the project. These SBOMs are submitted to GitHub's Dependency Graph, enabling Dependabot to monitor for security vulnerabilities.

## Formats

We generate SBOM in two industry-standard formats:

- **SPDX 2.3 JSON**: ISO/IEC standard format for software component information
- **CycloneDX JSON**: OWASP standard format designed for application security

## Automation

### GitHub Actions Workflow

The SBOM generation is automated through GitHub Actions (`.github/workflows/sbom.yml`):

**Triggers:**
- Push to `main` branch
- Pull requests to `main`
- Weekly schedule (Monday at 00:00 UTC)
- Manual workflow dispatch

**Steps:**
1. Install [Syft](https://github.com/anchore/syft) - a CLI tool for SBOM generation
2. Scan the entire project directory
3. Generate SBOM in both SPDX and CycloneDX formats
4. Submit SPDX SBOM to GitHub Dependency Graph (main branch only)
5. Upload SBOM files as workflow artifacts (90-day retention)

### Dependency Graph Submission

On push to `main` branch, the SPDX SBOM is automatically submitted to GitHub's Dependency Graph using the `advanced-security/spdx-dependency-submission-action`. This enables:

- Automatic dependency tracking in the repository's Insights tab
- Dependabot alerts for vulnerable dependencies
- Dependency review for pull requests

## Local Generation

You can generate SBOM files locally using Nix:

```bash
# Enter the Nix development shell (syft is pre-installed)
nix develop

# Generate SPDX format
syft dir:. -o spdx-json=sbom.spdx.json

# Generate CycloneDX format
syft dir:. -o cyclonedx-json=sbom.cyclonedx.json
```

Or using a one-liner:

```bash
nix develop --command syft dir:. -o spdx-json=sbom.spdx.json
nix develop --command syft dir:. -o cyclonedx-json=sbom.cyclonedx.json
```

## What's Included

The generated SBOM includes:

- **Go modules**: All direct and transitive dependencies from `go.mod` files
- **GitHub Actions**: Actions used in workflows (e.g., `actions/checkout`, `actions/setup-go`)
- **Package metadata**:
  - Package names and versions
  - License information
  - Package URLs (purl)
  - CPE identifiers for security tracking
  - SHA256 checksums where available

### Example Components

From a typical SBOM generation, you'll find:

- Core dependencies: `github.com/spf13/cobra`, `github.com/sergi/go-diff`
- MCP SDK: `github.com/modelcontextprotocol/go-sdk`
- Testing utilities: `github.com/stretchr/testify` dependencies
- CI/CD actions: `actions/checkout@v4.1.7`, `actions/setup-go@v5.1.0`

## SBOM Files

Generated SBOM files are:
- **Gitignored**: Pattern `sbom.*.json` in `.gitignore`
- **Size**: Typically 70-75 KB (JSON format, single line)
- **Location**: Project root directory when generated locally

## Security Benefits

1. **Transparency**: Complete visibility into all software components
2. **Vulnerability Management**: Automatic detection of known vulnerabilities through Dependabot
3. **Compliance**: Meet supply chain security requirements (e.g., Executive Order 14028)
4. **Audit Trail**: Historical SBOM artifacts available for 90 days

## Tools

### Syft

[Syft](https://github.com/anchore/syft) is used for SBOM generation because:

- Fast and accurate Go module scanning
- Supports multiple SBOM formats
- Actively maintained by Anchore
- Works well in CI/CD environments
- Available in nixpkgs

### Alternative: sbomnix

Initially, we attempted to use `sbomnix` for Nix-native SBOM generation, but encountered build issues with the `ruby3.3-nokogiri` dependency. Syft was chosen as a stable alternative with excellent Go support.

## Troubleshooting

### SBOM Generation Fails

If SBOM generation fails locally:

```bash
# Ensure you're in the Nix shell
nix develop

# Verify syft is available
syft version

# Check Go modules are downloaded
go mod download
cd cmd/lm-suggester && go mod download && cd ../..

# Try generating again with verbose output
syft dir:. -o spdx-json -vv
```

### Dependency Graph Submission Fails

Common issues:

- **Permissions**: The workflow needs `contents: write` permission
- **Branch protection**: Ensure the workflow has permission to access the repository
- **Fork PRs**: Dependency submission is skipped for PRs from forks (security limitation)

Check the workflow run logs in the Actions tab for detailed error messages.

## References

- [SPDX Specification](https://spdx.dev/specifications/)
- [CycloneDX Specification](https://cyclonedx.org/specification/overview/)
- [Syft Documentation](https://github.com/anchore/syft)
- [GitHub Dependency Submission API](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/using-the-dependency-submission-api)
- [Executive Order on Cybersecurity](https://www.nist.gov/itl/executive-order-14028-improving-nations-cybersecurity)
