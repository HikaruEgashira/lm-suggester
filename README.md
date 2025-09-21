# llm-suggester

[![Go Reference](https://pkg.go.dev/badge/github.com/HikaruEgashira/llm-suggester.svg)](https://pkg.go.dev/github.com/HikaruEgashira/llm-suggester)
[![Test](https://github.com/HikaruEgashira/llm-suggester/actions/workflows/test.yml/badge.svg)](https://github.com/HikaruEgashira/llm-suggester/actions/workflows/test.yml)
[![Release](https://github.com/HikaruEgashira/llm-suggester/actions/workflows/release.yml/badge.svg)](https://github.com/HikaruEgashira/llm-suggester/actions/workflows/release.yml)

Convert LLM suggestions and external tool outputs to [reviewdog](https://github.com/reviewdog/reviewdog) JSON format for seamless code review automation.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap HikaruEgashira/tap
brew install llm-suggester
```

### Go Install

```bash
go install github.com/HikaruEgashira/llm-suggester/cmd/llm-suggester@latest
```

### Binary Download

Download the latest binary from the [releases page](https://github.com/HikaruEgashira/llm-suggester/releases).

## Usage

### Basic Usage

The CLI reads JSON input and converts it to reviewdog format:

```bash
# From stdin
echo '{"file_path":"main.go","llm_after":"fixed code","message":"Fix typo"}' | llm-suggester

# From file
llm-suggester -i suggestion.json

# To file
llm-suggester -i suggestion.json -o reviewdog.json

# Pretty print
llm-suggester -i suggestion.json -p
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
- `source_name` (optional): Name of the tool. Default: llm-suggester

### Integration with reviewdog

Pipe the output directly to reviewdog:

```bash
# Generate suggestion and review
your-llm-tool | llm-suggester | reviewdog -f=rdjson -reporter=github-pr-review

# With GitHub Actions
llm-suggester -i suggestion.json | reviewdog -f=rdjson -reporter=github-pr-check
```

### Examples

#### Simple Replacement

```bash
cat <<EOF | llm-suggester
{
  "file_path": "main.go",
  "llm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print"
}
EOF
```

#### With Context

```bash
cat <<EOF | llm-suggester -p
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
  llm-suggester -i "$file" >> combined.json
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
    llm-suggester -i suggestions.json | \
      reviewdog -f=rdjson -reporter=github-pr-review
```

### Local Development

```bash
# Run local code review
git diff HEAD^ | llm-tool suggest | \
  llm-suggester | \
  reviewdog -f=rdjson -reporter=local -diff="git diff HEAD^"
```

### Custom Tool Integration

```go
// Use as a library
import "github.com/HikaruEgashira/llm-suggester/suggester"

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
  llm-suggester [flags]

Flags:
  -i, --input string    Input JSON file (default: stdin)
  -o, --output string   Output file (default: stdout)
  -p, --pretty          Pretty-print JSON output
  -h, --help            Help for llm-suggester
      --version         Version information
```

## Development

### Build from Source

```bash
git clone https://github.com/HikaruEgashira/llm-suggester.git
cd llm-suggester
go build ./cmd/llm-suggester
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