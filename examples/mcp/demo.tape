# VHS documentation: https://github.com/charmbracelet/vhs
Output examples/mcp/demo.gif

Set FontSize 13
Set Width 1400
Set Height 600
Set Theme "Catppuccin Mocha"
Set Padding 20

Type "# MCP oneshot mode demo - pipe JSONRPC messages directly"
Enter
Sleep 1s

Type "echo -e '"
Sleep 0.3
Type '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}'
Enter
Type '{"jsonrpc":"2.0","method":"notifications/initialized"}'
Enter
Type '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"suggest","arguments":{"file_path":"/tmp/test.go","base_text":"package main\\n\\nfunc main() {\\n\\tprint(\\"Hello\\")\\n}","lm_before":"print(\\"Hello\\")","lm_after":"fmt.Println(\\"Hello\\")","message":"Use fmt.Println","reporter":"local"}}}'
Type "' | lm-suggester mcp --oneshot"
Enter
Sleep 3s
