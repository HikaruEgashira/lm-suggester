# VHS documentation: https://github.com/charmbracelet/vhs
Output examples/mcp/inspector.gif

Set FontSize 14
Set Width 1600
Set Height 900
Set Theme "Catppuccin Mocha"
Set Padding 30

Type "# lm-suggester MCP server with Inspector"
Enter
Sleep 1s

Type "# Start MCP Inspector - automatically handles initialize handshake"
Enter
Sleep 0.5

Type "npx @modelcontextprotocol/inspector --cli cmd/lm-suggester/lm-suggester mcp --method tools/list"
Enter
Sleep 8s

Type "# Call the suggest tool via Inspector"
Enter
Sleep 0.5

Type 'npx @modelcontextprotocol/inspector --cli cmd/lm-suggester/lm-suggester mcp --method tools/call --tool-name suggest --tool-arg file_path=/tmp/test.go --tool-arg lm_after="fmt.Println(\"Hello\")" --tool-arg message="Use fmt.Println instead of print"'
Enter
Sleep 8s
