package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/HikaruEgashira/lm-suggester/suggester"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

// SuggestInput defines the input schema for the suggest tool
type SuggestInput struct {
	FilePath   string `json:"file_path" jsonschema:"required" jsonschema_description:"Path to the file to be reviewed"`
	BaseText   string `json:"base_text,omitempty" jsonschema_description:"Original file content for diff calculation"`
	LMBefore   string `json:"lm_before,omitempty" jsonschema_description:"Exact text to be replaced (must match exactly including whitespace)"`
	LMAfter    string `json:"lm_after" jsonschema:"required" jsonschema_description:"Replacement text or suggestion"`
	Message    string `json:"message" jsonschema:"required" jsonschema_description:"Explanation or reason for the suggestion"`
	Severity   string `json:"severity,omitempty" jsonschema_description:"Severity level (ERROR, WARNING, INFO)"`
	SourceName string `json:"source_name,omitempty" jsonschema_description:"Name of the tool or LLM that generated this suggestion"`
	Reporter   string `json:"reporter,omitempty" jsonschema_description:"Reviewdog reporter to use (local, github-pr-review, etc.)"`
}

// SuggestOutput defines the output of the tool
type SuggestOutput struct {
	Success bool   `json:"success" jsonschema_description:"Whether the suggestion was successfully posted"`
	Output  string `json:"output" jsonschema_description:"Output from reviewdog (stdout)"`
	Error   string `json:"error,omitempty" jsonschema_description:"Error output from reviewdog (stderr) if any"`
}

// jsonrpcRequest represents a JSON-RPC 2.0 request
type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// jsonrpcResponse represents a JSON-RPC 2.0 response
type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
}

// runOneshotMode processes JSONRPC messages from stdin and exits
func runOneshotMode() error {
	// Read all messages from stdin (newline-delimited JSONRPC)
	scanner := bufio.NewScanner(os.Stdin)
	var messages []json.RawMessage
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		messages = append(messages, json.RawMessage(line))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	// Find tools/call request
	var toolCallID json.RawMessage
	var toolCallParams struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	found := false
	for _, msg := range messages {
		var req jsonrpcRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			continue
		}
		if req.Method == "tools/call" {
			toolCallID = req.ID
			if err := json.Unmarshal(req.Params, &toolCallParams); err != nil {
				return fmt.Errorf("failed to unmarshal tools/call params: %w", err)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no tools/call request found in stdin")
	}

	// Only support "suggest" tool
	if toolCallParams.Name != "suggest" {
		return fmt.Errorf("unsupported tool: %s (only 'suggest' is supported)", toolCallParams.Name)
	}

	// Parse suggest arguments
	var input SuggestInput
	argsJSON, err := json.Marshal(toolCallParams.Arguments)
	if err != nil {
		return fmt.Errorf("failed to marshal arguments: %w", err)
	}
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return fmt.Errorf("failed to unmarshal suggest input: %w", err)
	}

	// Execute suggest function
	_, output, err := suggest(context.Background(), nil, input)
	if err != nil {
		return fmt.Errorf("suggest execution failed: %w", err)
	}

	// Build JSONRPC response
	response := jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      toolCallID,
		Result: map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": output.Output},
			},
			"isError": !output.Success,
		},
	}

	// Output response as JSON
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	if !output.Success {
		return fmt.Errorf("suggest failed: %s", output.Error)
	}

	return nil
}

// suggest handles the MCP tool call to suggest code changes via reviewdog
func suggest(ctx context.Context, req *mcp.CallToolRequest, input SuggestInput) (*mcp.CallToolResult, SuggestOutput, error) {
	// Marshal the input back to JSON for the suggester
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to marshal input: %v", err)}},
			IsError: true,
		}, SuggestOutput{Success: false, Error: fmt.Sprintf("marshal error: %v", err)}, nil
	}

	// Convert using the suggester library
	rdJSON, err := suggester.Convert(inputJSON, "reviewdog")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to convert: %v", err)}},
			IsError: true,
		}, SuggestOutput{Success: false, Error: fmt.Sprintf("convert error: %v", err)}, nil
	}

	// Determine reporter
	reporter := "local"
	if input.Reporter != "" {
		reporter = input.Reporter
	}

	// Run reviewdog and capture output
	stdout, stderr, err := runReviewdogCapture(rdJSON, reporter, "nofilter", false)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("reviewdog failed: %v\nStdout: %s\nStderr: %s", err, string(stdout), string(stderr))}},
			IsError: true,
		}, SuggestOutput{Success: false, Output: string(stdout), Error: fmt.Sprintf("%v\n%s", err, string(stderr))}, nil
	}

	return nil, SuggestOutput{
		Success: true,
		Output:  string(stdout),
		Error:   string(stderr),
	}, nil
}

// newMCPCommand creates the MCP subcommand
func newMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Run lm-suggester as an MCP server",
		Long: `Run lm-suggester as a Model Context Protocol (MCP) server.

This allows AI assistants and LLMs to post code review suggestions directly
through the MCP protocol over stdin/stdout.

The server provides the following tool:
  - suggest: Post code review suggestions via reviewdog

Example usage with an MCP client (interactive mode):
  lm-suggester mcp

Example usage with oneshot mode (pipe):
  echo -e 'JSONRPC_MSG1\nJSONRPC_MSG2\nJSONRPC_MSG3' | lm-suggester mcp --oneshot

Example tool call:
  {
    "file_path": "main.go",
    "base_text": "package main\n\nfunc main() {\n\tprint(\"Hello\")\n}",
    "lm_before": "print(\"Hello\")",
    "lm_after": "fmt.Println(\"Hello\")",
    "message": "Use fmt.Println instead of print",
    "reporter": "local"
  }`,
		RunE: func(cmd *cobra.Command, args []string) error {
			oneshot, _ := cmd.Flags().GetBool("oneshot")
			if oneshot {
				return runOneshotMode()
			}

			// Create MCP server
			server := mcp.NewServer(&mcp.Implementation{
				Name:    "lm-suggester",
				Version: version,
			}, nil)

			// Add the suggest tool
			mcp.AddTool(server, &mcp.Tool{
				Name:        "suggest",
				Description: "Post code review suggestions via reviewdog. Converts LLM suggestions to reviewdog format and runs reviewdog to post them. Takes file path, optional base text, optional exact match text (lm_before), replacement text (lm_after), explanation message, and optional reporter (local, github-pr-review, etc.). Returns success status, reviewdog output, and any errors. The output is captured and returned without polluting the MCP JSON-RPC protocol.",
			}, suggest)

			// Run the server over stdin/stdout
			log.Printf("Starting MCP server (version %s)...", version)
			if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
				return fmt.Errorf("MCP server failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().Bool("oneshot", false, "Process JSONRPC messages from stdin and exit (for pipe usage)")

	return cmd
}
