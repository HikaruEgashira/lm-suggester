# lm-suggester

[![Go Reference](https://pkg.go.dev/badge/github.com/HikaruEgashira/lm-suggester.svg)](https://pkg.go.dev/github.com/HikaruEgashira/lm-suggester)
[![Test](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/test.yml/badge.svg)](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/test.yml)
[![Release](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/release.yml/badge.svg)](https://github.com/HikaruEgashira/lm-suggester/actions/workflows/release.yml)

LLMや外部ツールの提案を[reviewdog](https://github.com/reviewdog/reviewdog) JSON形式に変換し、コードレビューを自動化するツールです。

## Installation

### Go Install

```bash
go install github.com/HikaruEgashira/lm-suggester/cmd/lm-suggester@latest
```

### Using lm-suggester cli

```bash
# Install globally with mise
mise use -g github:HikaruEgashira/lm-suggester

# Or install locally in your project
mise use github:HikaruEgashira/lm-suggester
```

その他のインストール方法は[リリースページ](https://github.com/HikaruEgashira/lm-suggester/releases)を参照してください。

## Usage

### 基本的な使い方

CLIはJSON入力を読み込み、reviewdog形式に変換します：

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

### 入出力形式

#### 入力 (必須フィールドのみ)

```json
{
  "file_path": "main.go",
  "lm_after": "fmt.Println(\"Hello, World!\")",
  "message": "Use fmt.Println instead of print"
}
```

#### 出力 (reviewdog JSON)

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

### SARIF形式への変換

reviewdog経由でSARIF形式に変換できます：

```bash
lm-suggester -i suggestion.json | reviewdog -f=rdjson -reporter=sarif
```

#### SARIF出力例

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

### reviewdogとの連携

#### 手動連携

出力を直接reviewdogにパイプ：

```bash
# Generate suggestion and review
your-llm-tool | lm-suggester | reviewdog -f=rdjson -reporter=github-pr-review

# With GitHub Actions
lm-suggester -i suggestion.json | reviewdog -f=rdjson -reporter=github-pr-check
```

#### 自動実行

`--reviewdog`フラグで自動的にreviewdogを実行：

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

利用可能なreviewdogオプションは[reviewdogドキュメント](https://github.com/reviewdog/reviewdog)を参照してください。


## 使用例

### CI/CD連携

```yaml
# GitHub Actions example
- name: Run LLM Review
  run: |
    llm-tool analyze --output suggestions.json
    lm-suggester -i suggestions.json | \
      reviewdog -f=rdjson -reporter=github-pr-review
```

### ローカル開発

```bash
# Run local code review
git diff HEAD^ | llm-tool suggest | \
  lm-suggester | \
  reviewdog -f=rdjson -reporter=local -diff="git diff HEAD^"
```

### ライブラリとして使用

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

## コマンドオプション

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

## 関連プロジェクト

- [reviewdog](https://github.com/reviewdog/reviewdog) - Automated code review tool
- [reviewdog/action-suggester](https://github.com/reviewdog/action-suggester) - GitHub Action for suggestions