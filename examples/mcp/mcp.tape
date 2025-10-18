# VHS documentation: https://github.com/charmbracelet/vhs
Output examples/mcp/mcp.gif

Set FontSize 14
Set Width 1600
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 30

Type "# lm-suggester MCP server - pipe JSONRPC messages"
Enter
Sleep 1s

Type 'echo -e '''
Sleep 0.5
Type '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}\n'
Sleep 0.3
Type '{"jsonrpc":"2.0","method":"notifications/initialized"}\n'
Sleep 0.3
Type '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"suggest","arguments":{"file_path":"/tmp/test.go","base_text":"package main\\n\\nfunc main() {\\n\\tprint(\\"Hello\\")\\n}","lm_before":"print(\\"Hello\\")","lm_after":"fmt.Println(\\"Hello\\")","message":"Use fmt.Println instead of print","reporter":"local"}}}'
Sleep 0.5
Type "' | lm-suggester mcp"
Enter
Sleep 4s
