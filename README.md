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
echo '{"file_path":"main.go","llm_after":"fixed code","message":"Fix typo"}' | lm-suggester

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
  "llm_before": "old code",
  "llm_after": "new code",
  "message": "Description of the change",
  "severity": "INFO",
  "source_name": "my-linter"
}
```

#### Fields

- `file_path` (required): Path to the target file
- `llm_after` (required): The suggested replacement text
- `message` (required): Description of the suggestion
- `base_text` (optional): Original file content. If not provided, reads from `file_path`
- `llm_before` (optional): Text to be replaced. If not provided, computes minimal diff automatically
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

Available reviewdog options:
- `--reviewdog`: Enable automatic reviewdog execution
- `--reporter`: Set reviewdog reporter (default: local)
  - `local`: Show results in terminal
  - `github-pr-review`: GitHub PR review comments
  - `github-pr-check`: GitHub PR checks
  - `gitlab-mr-discussion`: GitLab MR discussions
- `--filter-mode`: Set filter mode (default: added)
  - `added`: Only new issues
  - `diff_context`: Issues in diff context
  - `file`: All issues in changed files
  - `nofilter`: All issues
- `--fail-on-error`: Exit with non-zero code if issues found

### Examples

#### Simple Replacement

```bash
cat <<EOF | lm-suggester
{
  "file_path": "main.go",
  "llm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print"
}
EOF
```

#### With Context

```bash
cat <<EOF | lm-suggester -p
{
  "file_path": "utils/helper.go",
  "llm_before": "// TODO: implement this",
  "llm_after": "// parseConfig reads the configuration file",
  "message": "Add proper documentation",
  "severity": "WARNING",
  "source_name": "doc-linter"
}
EOF
```

#### Batch Processing

```bash
# Process multiple suggestions
for file in suggestions/*.json; do
  lm-suggester -i "$file" >> combined.json
done

# Review all at once
cat combined.json | reviewdog -f=rdjson -reporter=local
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
    LLMAfter:  "improved code",
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
      --filter-mode string   reviewdog filter mode (default: added)
      --fail-on-error        Exit with non-zero code when reviewdog finds errors
  -h, --help                 Help for lm-suggester
      --version              Version information
```

## Development

### Build from Source

```bash
git clone https://github.com/HikaruEgashira/lm-suggester.git
cd lm-suggester
go build ./cmd/lm-suggester
```

### Run Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [reviewdog](https://github.com/reviewdog/reviewdog) - Automated code review tool
- [reviewdog/action-suggester](https://github.com/reviewdog/action-suggester) - GitHub Action for suggestions