reviewdog-converter は LLM などの外部ツールから受け取った提案を reviewdog JSON 形式へ変換するGoライブラリです。

```bash
# 依存関係の管理
go mod tidy

# テスト実行
go test ./...
go test -race ./...  # レース検出付き

# 静的検査
go vet ./...

# カバレッジ確認
go test -cover ./...

# ベンチマーク（性能計測）
go test -run none -bench . ./suggester

# サンプル実行
cat _examples/testdata/simple_replacement.json | go run _examples/simple/main.go
```

## アーキテクチャ

### コアライブラリ構成

`suggester/` パッケージが変換ロジックを提供。

- public.go: 公開 API (`Input` 構造体、`BuildRDJSON` 関数)
  - LLM の提案を受け取り、reviewdog JSON 形式に変換
  - `LLMBefore` が空の場合は最小差分を自動計算

- align.go: 差分位置の整列処理
  - ベーステキストと LLM 提案の位置合わせ

- diff.go: 差分検出
  - 最小差分範囲の特定

- normalize.go: テキスト正規化
  - UTF-8 マルチバイト文字と改行の正規化

- errors.go: エラー定義
  - カスタムエラー型の定義

### 入出力形式

入力 (`suggester.Input`):
```go
type Input struct {
    FilePath   string  // 対象ファイルパス
    BaseText   string  // 元のファイル内容
    LLMBefore  string  // 変更前テキスト（optional）
    LLMAfter   string  // 変更後テキスト
    Message    string  // サジェストメッセージ
    Severity   string  // 重要度
    SourceName string  // ツール名
}
```

出力: reviewdog JSON 形式のバイト列

## コーディングスタイルと命名
Go の標準スタイル (タブインデント) を厳守し、`gofmt` / `goimports` 後にコミットします。公開シンボルは `Suggester` 接頭辞で役割を示し、パッケージ内限定のものは小文字始まりに統一します。公開 API には `// Name ...` 形式のコメントを付け、エラーは `errors.go` の型を利用して英語で簡潔に表現します。補助スクリプトは `_examples/` または `tools/` に整理し、`mise tasks edit` で共有タスクを登録します。

## テスト指針
テーブル駆動の `testing` を基本とし、ファイル名は `*_test.go` へ統一します。サブテストは `t.Run` を使い、バグ修正は `Test_issue123_description` のように Issue 番号を含めます。`go test -cover ./suggester` の実測カバレッジは 52.1% なので、新規コードはカバレッジ向上を意識してください。性能が疑わしいときは `go test -run none -bench . ./suggester` を活用します。

## コミットとプルリクエスト
履歴は Conventional Commits (`feat: ...`, `fix: ...`) に沿っています。各コミット前に `go test ./...` を再度実行し、変更の粒度を最小限に分割します。PR は `gh pr create --fill` を基本に、概要・動機・テスト結果 (`go test ./...` ログ)・関連 Issue (`Fixes #123`) を本文へ記載します。提出後は `gh pr checks --watch` で CI 成功を確認し、必要に応じて `gh pr view --web` で差分を共有します。
