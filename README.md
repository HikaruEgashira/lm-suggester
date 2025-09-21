# lm-suggester

[![Go Reference](https://pkg.go.dev/badge/github.com/HikaruEgashira/lm-suggester.svg)](https://pkg.go.dev/github.com/HikaruEgashira/lm-suggester)
[![Test](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/test.yml/badge.svg)](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/test.yml)
[![Release](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/release.yml/badge.svg)](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/release.yml)

Convert LLM suggestions and external tool outputs to [reviewdog](https://github.com/reviewdog/reviewdog) JSON format for seamless code review automation.

## Installation

### Go Install

```bash
go install github.com/HikaruEgashira/lm-suggester/cmd/lm-suggester@latest
```

### Binary Download

Download the latest binary from the [releases page](https://github.com/HikaruEgashira/lm-suggester/releases).

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

### Input Format

The CLI expects JSON input with the following structure:

```json
{
  "file_path": "path/to/file.go",
  "base_text": "original file content",
  "lm_before": "old code",
  "lm_after": "new code",
  "message": "Description of the change",
  "severity": "INFO",
  "source_name": "my-linter"
}
```

#### Fields

- `file_path` (required): Path to the target file
- `lm_after` (required): The suggested replacement text
- `message` (required): Description of the suggestion
- `base_text` (optional): Original file content. If not provided, reads from `file_path`
- `lm_before` (optional): Text to be replaced. If not provided, computes minimal diff automatically
- `severity` (optional): Severity level (INFO, WARNING, ERROR). Default: INFO
- `source_name` (optional): Name of the tool. Default: lm-suggester

### Integration with reviewdog

#### Manual Integration

Pipe the output directly to reviewdog:

```bash
# Generate suggestion and review
your-llm-tool | lm-suggester | reviewdog -f=rdjson -reporter=github-pr-review

# With GitHub Actions
lm-suggester -i suggestion.json | reviewdog -f=rdjson -reporter=github-pr-check
```

#### Automatic reviewdog Execution

Use the `--reviewdog` flag to automatically run reviewdog:

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

For available reviewdog options, see [reviewdog documentation](https://github.com/reviewdog/reviewdog).

### Examples

#### Simple Replacement

```bash
cat <<EOF | lm-suggester
{
  "file_path": "main.go",
  "lm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print"
}
EOF
```


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

### Custom Tool Integration

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
      --reporter string      reviewdog reporter (default: local)
      --filter-mode string   reviewdog filter mode (default: nofilter)
      --fail-on-error        Exit with non-zero code when reviewdog finds errors
  -h, --help                 Help for lm-suggester
      --version              Version information
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [reviewdog](https://github.com/reviewdog/reviewdog) - Automated code review tool
- [reviewdog/action-suggester](https://github.com/reviewdog/action-suggester) - GitHub Action for suggestions