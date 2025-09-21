# reviewdog-converter

[![Go Reference](https://pkg.go.dev/badge/github.com/HikaruEgashira/reviewdog-converter.svg)](https://pkg.go.dev/github.com/HikaruEgashira/reviewdog-converter)
[![Test](https://github.com/HikaruEgashira/reviewdog-converter/actions/workflows/test.yml/badge.svg)](https://github.com/HikaruEgashira/reviewdog-converter/actions/workflows/test.yml)
[![Release](https://github.com/HikaruEgashira/reviewdog-converter/actions/workflows/release.yml/badge.svg)](https://github.com/HikaruEgashira/reviewdog-converter/actions/workflows/release.yml)

Convert LLM suggestions and external tool outputs to [reviewdog](https://github.com/reviewdog/reviewdog) JSON format for seamless code review automation.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap HikaruEgashira/tap
brew install reviewdog-converter
```

### Go Install

```bash
go install github.com/HikaruEgashira/reviewdog-converter/cmd/reviewdog-converter@latest
```

### Binary Download

Download the latest binary from the [releases page](https://github.com/HikaruEgashira/reviewdog-converter/releases).

## Usage

### Basic Usage

The CLI reads JSON input and converts it to reviewdog format:

```bash
# From stdin
echo '{"file_path":"main.go","llm_after":"fixed code","message":"Fix typo"}' | reviewdog-converter

# From file
reviewdog-converter -i suggestion.json

# To file
reviewdog-converter -i suggestion.json -o reviewdog.json

# Pretty print
reviewdog-converter -i suggestion.json -p
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
- `source_name` (optional): Name of the tool. Default: reviewdog-converter

### Integration with reviewdog

Pipe the output directly to reviewdog:

```bash
# Generate suggestion and review
your-llm-tool | reviewdog-converter | reviewdog -f=rdjson -reporter=github-pr-review

# With GitHub Actions
reviewdog-converter -i suggestion.json | reviewdog -f=rdjson -reporter=github-pr-check
```

### Examples

#### Simple Replacement

```bash
cat <<EOF | reviewdog-converter
{
  "file_path": "main.go",
  "llm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print"
}
EOF
```

#### With Context

```bash
cat <<EOF | reviewdog-converter -p
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
  reviewdog-converter -i "$file" >> combined.json
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
    reviewdog-converter -i suggestions.json | \
      reviewdog -f=rdjson -reporter=github-pr-review
```

### Local Development

```bash
# Run local code review
git diff HEAD^ | llm-tool suggest | \
  reviewdog-converter | \
  reviewdog -f=rdjson -reporter=local -diff="git diff HEAD^"
```

### Custom Tool Integration

```go
// Use as a library
import "github.com/HikaruEgashira/reviewdog-converter/suggester"

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
  reviewdog-converter [flags]

Flags:
  -i, --input string    Input JSON file (default: stdin)
  -o, --output string   Output file (default: stdout)
  -p, --pretty          Pretty-print JSON output
  -h, --help            Help for reviewdog-converter
      --version         Version information
```

## Development

### Build from Source

```bash
git clone https://github.com/HikaruEgashira/reviewdog-converter.git
cd reviewdog-converter
go build ./cmd/reviewdog-converter
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