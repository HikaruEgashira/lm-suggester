# VHS documentation: https://github.com/charmbracelet/vhs
Output examples/mcp/demo.gif

Set FontSize 14
Set Width 1200
Set Height 1000
Set Theme "Catppuccin Mocha"
Set Padding 20

# Show MCP command help
Type "# lm-suggester MCP Demo"
Enter
Sleep 1s
Type "# Step 1: Show MCP command help"
Enter
Sleep 500ms
Type "lm-suggester mcp --help"
Enter
Sleep 3s

Type "clear"
Enter
Sleep 500ms

# Create a test Go file with an issue
Type "# Step 2: Create test file with a code issue"
Enter
Sleep 500ms
Type "cat /tmp/example.go"
Enter
Sleep 2s

Type "clear"
Enter
Sleep 500ms

# Run the MCP client that calls the suggest tool
Type "# Step 3: Run MCP client to call suggest tool"
Enter
Sleep 500ms
Type "cd /tmp && go run mcp_demo_client.go"
Enter
Sleep 5s
