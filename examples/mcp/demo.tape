# VHS documentation: https://github.com/charmbracelet/vhs
Output examples/mcp/demo.gif

Set FontSize 14
Set Width 1200
Set Height 600

Type "# MCP Server Demo - lm-suggester"
Enter
Sleep 1s

Type "# Start the MCP server"
Enter
Type "lm-suggester mcp --help"
Enter
Sleep 2s

Type "clear"
Enter
Sleep 500ms

Type "# The MCP server provides a 'suggest' tool for code review"
Enter
Type "# Let's test it with a Go client"
Enter
Sleep 1s

Type "cat << 'EOF' > /tmp/test_suggest.go"
Enter
Type "package main"
Enter
Enter
Type "import ("
Enter
Type '	"context"'
Enter
Type '	"fmt"'
Enter
Type '	"log"'
Enter
Type '	"os/exec"'
Enter
Enter
Type '	"github.com/modelcontextprotocol/go-sdk/mcp"'
Enter
Type ")"
Enter
Enter
Type "func main() {"
Enter
Type '	ctx := context.Background()'
Enter
Type '	client := mcp.NewClient(&mcp.Implementation{Name: "demo"}, nil)'
Enter
Type '	transport := &mcp.CommandTransport{'
Enter
Type '		Command: exec.Command("lm-suggester", "mcp"),'
Enter
Type '	}'
Enter
Type '	session, _ := client.Connect(ctx, transport, nil)'
Enter
Type '	defer session.Close()'
Enter
Enter
Type '	// List available tools'
Enter
Type '	tools, _ := session.ListTools(ctx, &mcp.ListToolsParams{})'
Enter
Type '	fmt.Println("Available tools:")'
Enter
Type '	for _, tool := range tools.Tools {'
Enter
Type '		fmt.Printf("  - %s\\n", tool.Name)'
Enter
Type '	}'
Enter
Type "}"
Enter
Type "EOF"
Enter
Sleep 500ms

Type "cd /tmp && go mod init demo && go get github.com/modelcontextprotocol/go-sdk/mcp"
Enter
Sleep 2s

Type "go run test_suggest.go"
Enter
Sleep 2s

Type "# MCP server successfully lists the 'suggest' tool!"
Enter
Sleep 2s
