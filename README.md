# lm-suggester

[![Go Reference](https://pkg.go.dev/badge/github.com/HikaruEgashira/lm-suggester.svg)](https://pkg.go.dev/github.com/HikaruEgashira/lm-suggester)
[![Test](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/test.yml/badge.svg)](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/test.yml)
[![Release](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/release.yml/badge.svg)](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/release.yml)

Convert LLM suggestions and external tool outputs to [reviewdog](https://github.com/reviewdog/reviewdog) JSON format for seamless code review automation.

## Installation

### Using lm-suggester cli

```bash
# Install globally with mise
mise use -g github:HikaruEgashira/lm-suggester

# Or install locally in your project
mise use github:HikaruEgashira/lm-suggester
```

For other installation methods, see the [releases page](https://github.com/HikaruEgashira/lm-suggester/releases).

## Usage

### Basic Usage

The CLI reads JSON input and converts it to reviewdog format:

```bash
# From stdin
echo '{"file_path":"main.go","lm_after":"fixed code","message":"Fix typo"}' | lm-suggester

# From file
lm-suggester -i suggestion.json

# To file
lm-suggester -i suggestion.json -o reviewdog.json

# Pretty print
lm-suggester -i suggestion.json -p
```

### Input/Output Format

#### Input (required fields only)

```json
{
  "file_path": "main.go",
  "lm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print"
}
```

#### Output (reviewdog JSON)

```json
{
  "diagnostics": [
    {
      "message": "Use fmt.Println instead of print",
      "location": {
        "path": "main.go",
        "range": {
          "start": {"line": 10, "column": 5},
          "end": {"line": 10, "column": 30}
        }
      },
      "severity": "INFO",
      "source": {"name": "lm-suggester"},
      "suggestions": [
        {
          "range": {
            "start": {"line": 10, "column": 5},
            "end": {"line": 10, "column": 30}
          },
          "text": "fmt.Println(\"Hello, World!\")"
        }
      ]
    }
  ]
}
```

### SARIF Format Conversion

Convert to SARIF format via reviewdog:

```bash
lm-suggester -i suggestion.json | reviewdog -f=rdjson -reporter=sarif
```

#### SARIF Output Example

```json
{
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "lm-suggester"
        }
      },
      "results": [
        {
          "message": {"text": "Use fmt.Println instead of print"},
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {"uri": "main.go"},
                "region": {
                  "startLine": 10,
                  "startColumn": 5,
                  "endLine": 10,
                  "endColumn": 30
                }
              }
            }
          ],
          "fixes": [
            {
              "artifactChanges": [
                {
                  "artifactLocation": {"uri": "main.go"},
                  "replacements": [
                    {
                      "deletedRegion": {
                        "startLine": 10,
                        "startColumn": 5,
                        "endLine": 10,
                        "endColumn": 30
                      },
                      "insertedContent": {"text": "fmt.Println(\"Hello, World!\")"}
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### Integration with reviewdog

#### Manual Integration

Pipe the output directly to reviewdog:

```bash
# Generate suggestion and review
your-llm-tool | lm-suggester | reviewdog -f=rdjson -reporter=github-pr-review

# With GitHub Actions
lm-suggester -i suggestion.json | reviewdog -f=rdjson -reporter=github-pr-check
```

#### Automatic Execution

Automatically run reviewdog with the `--reviewdog` flag:

```bash
# Run reviewdog with local reporter
lm-suggester -i suggestion.json --reviewdog

# Specify reporter for CI/CD
lm-suggester -i suggestion.json --reviewdog --reporter=github-pr-review

# With custom options
lm-suggester -i suggestion.json --reviewdog \
  --reporter=github-pr-check \
  --filter-mode=diff_context \
  --fail-on-error

# From stdin
echo '{"file_path":"main.go","llm_after":"fixed","message":"Fix"}' | \
  lm-suggester --reviewdog --reporter=local
```

For available reviewdog options, see the [reviewdog documentation](https://github.com/reviewdog/reviewdog).


## Use Cases

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Run LLM Review
  run: |
    llm-tool analyze --output suggestions.json
    lm-suggester -i suggestions.json | \
      reviewdog -f=rdjson -reporter=github-pr-review
```

### Local Development

```bash
# Run local code review
git diff HEAD^ | llm-tool suggest | \
  lm-suggester | \
  reviewdog -f=rdjson -reporter=local -diff="git diff HEAD^"
```

### As a Library

```go
// Use as a library
import "github.com/HikaruEgashira/lm-suggester/suggester"

input := suggester.Input{
    FilePath:  "main.go",
    LMAfter:  "improved code",
    Message:   "Optimization suggestion",
}

rdJSON, err := suggester.BuildRDJSON(input)
```

## Command Options

```
Usage:
  lm-suggester [flags]

Flags:
  -i, --input string         Input JSON file (default: stdin)
  -o, --output string        Output file (default: stdout)
  -p, --pretty               Pretty-print JSON output
      --reviewdog            Run reviewdog with the output
      --reporter string      reviewdog reporter (local, github-pr-review, github-pr-check, etc.) (default: local)
      --filter-mode string   reviewdog filter mode (added, diff_context, file, nofilter) (default: nofilter)
      --fail-on-error        Exit with non-zero code when reviewdog finds errors
  -h, --help                 Help for lm-suggester
      --version              Print version information
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [reviewdog](https://github.com/reviewdog/reviewdog) - Automated code review tool
- [reviewdog/action-suggester](https://github.com/reviewdog/action-suggester) - GitHub Action for suggestions