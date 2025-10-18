# MCP Server Examples

lm-suggester の MCP サーバーを使用する方法を示します。

## MCP Inspector を使用する方法（推奨）

MCP Inspector は公式のテスト・デバッグツールで、`initialize` ハンドシェイクを自動的に処理します。

### ツール一覧を表示

```bash
npx @modelcontextprotocol/inspector --cli \
  cmd/lm-suggester/lm-suggester mcp \
  --method tools/list
```

### suggest ツールを呼び出す

```bash
npx @modelcontextprotocol/inspector --cli \
  cmd/lm-suggester/lm-suggester mcp \
  --method tools/call \
  --tool-name suggest \
  --tool-arg file_path=/tmp/test.go \
  --tool-arg lm_after='fmt.Println("Hello")' \
  --tool-arg message="Use fmt.Println instead of print"
```

### Web UI を使用する

```bash
npx @modelcontextprotocol/inspector cmd/lm-suggester/lm-suggester mcp
```

ブラウザで `http://localhost:6274` を開いて、視覚的にツールをテストできます。

![Inspector Demo](inspector.gif)

## 生の JSON-RPC メッセージを使用する方法（上級者向け）

MCP プロトコルの詳細を理解するために、直接 JSON-RPC メッセージを送信できます。

この方法では `initialize` と `initialized` を手動で送信する必要があります。

```bash
echo -e '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}
{"jsonrpc":"2.0","method":"notifications/initialized"}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"suggest","arguments":{"file_path":"/tmp/test.go","base_text":"package main\\n\\nfunc main() {\\n\\tprint(\\"Hello\\")\\n}","lm_before":"print(\\"Hello\\")","lm_after":"fmt.Println(\\"Hello\\")","message":"Use fmt.Println instead of print","reporter":"local"}}}' | lm-suggester mcp
```

![Raw JSON-RPC Demo](mcp.gif)

## なぜ initialize が必要なのか？

MCP プロトコル仕様では、以下のライフサイクルが定義されています：

1. **Initialize**: クライアントとサーバーがプロトコルバージョンと機能をネゴシエート（必須）
2. **Initialized 通知**: クライアントが準備完了を通知（必須）
3. **Operation**: 通常の操作（tools/call など）
4. **Shutdown**: 接続の終了

### 各アプローチの比較

| 方法 | `initialize` の扱い | 難易度 | 用途 |
|------|-------------------|--------|------|
| **MCP Inspector** | ✅ 自動処理 | 簡単 | 開発・テスト・デバッグ |
| **FastMCP など高レベルライブラリ** | ✅ 自動処理 | 簡単 | 本番アプリケーション |
| **生 JSON-RPC** | ❌ 手動で必要 | 難しい | プロトコル学習・デバッグ |

## 参考資料

- [MCP Inspector 公式ドキュメント](https://modelcontextprotocol.io/docs/tools/inspector)
- [MCP プロトコル仕様](https://modelcontextprotocol.io/specification/2025-03-26/basic/lifecycle)
- [FastMCP (Python)](https://github.com/jlowin/fastmcp)
