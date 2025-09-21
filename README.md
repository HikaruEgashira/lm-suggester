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
echo '{"file_path":"main.go","lm_after":"fixed code","message":"Fix typo"}' | lm-suggester
```

### Input/Output Format

#### Input

```json
{
  "file_path": "main.go",
  "lm_before": "print(\"Hello, World!\")",
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

### Pass-through Format Support

lm-suggester supports pass-through of additional fields. Any extra fields in the input JSON are preserved alongside the computed diagnostics:

#### Input with additional fields

```json
{
  "file_path": "main.go",
  "base_text": "package main\n\nfunc main() {\n\tprint(\"Hello, World!\")\n}",
  "lm_before": "print(\"Hello, World!\")",
  "lm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print",
  "ruleId": "go/print-style",
  "level": "warning",
  "properties": {
    "tags": ["style", "best-practice"],
    "category": "code-quality"
  }
}
```

#### Output (merged fields with diagnostics)

```json
{
  "file_path": "main.go",
  "base_text": "package main\n\nfunc main() {\n\tprint(\"Hello, World!\")\n}",
  "lm_before": "print(\"Hello, World!\")",
  "lm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print",
  "ruleId": "go/print-style",
  "level": "warning",
  "properties": {
    "tags": ["style", "best-practice"],
    "category": "code-quality"
  },
  "diagnostics": [
    {
      "message": "Replace code with suggestion\n```suggestion\nfmt.Println(\"Hello, World!\")\n```",
      "location": {
        "path": "main.go",
        "range": {
          "start": {"line": 4, "column": 2},
          "end": {"line": 4, "column": 24}
        }
      }
    }
  ],
  "source": {
    "name": "reviewdog-converter"
  }
}
```

This pass-through behavior preserves all input fields, allowing integration with various tools while maintaining reviewdog compatibility.

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

## LLM Code Review Prompt

Minimal prompt for automated PR review with lm-suggester:

```
Analyze the code and execute the following pipeline:

1. Generate suggestions in this format (one per line, not an array):
{"file_path":"path/to/file","base_text":"<full file content>","lm_before":"<exact match>","lm_after":"<replacement>","message":"<reason>"}

Requirements:
- lm_before must match exactly (including whitespace)
- Include complete base_text for line number calculation
- Each suggestion as separate JSON object

2. Save suggestions and execute:
cat suggestions.json | lm-suggester --reviewdog --reporter=github-pr-review
```

## Related Projects

- [reviewdog](https://github.com/reviewdog/reviewdog) - Automated code review tool
- [reviewdog/action-suggester](https://github.com/reviewdog/action-suggester) - GitHub Action for suggestions