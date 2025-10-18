# SBOM (Software Bill of Materials)

lm-suggester automatically generates SBOM files to provide transparency about the software components and dependencies used in the project. These SBOMs are submitted to GitHub's Dependency Graph, enabling Dependabot to monitor for security vulnerabilities.

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

## References

- [SPDX Specification](https://spdx.dev/specifications/)
- [CycloneDX Specification](https://cyclonedx.org/specification/overview/)
- [Syft Documentation](https://github.com/anchore/syft)
- [GitHub Dependency Submission API](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/using-the-dependency-submission-api)
- [Executive Order on Cybersecurity](https://www.nist.gov/itl/executive-order-14028-improving-nations-cybersecurity)
