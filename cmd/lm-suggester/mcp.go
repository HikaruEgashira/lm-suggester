package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/HikaruEgashira/lm-suggester/suggester"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

// ConvertInput defines the input schema for the convert tool
type ConvertInput struct {
	FilePath   string `json:"file_path" jsonschema:"required" jsonschema_description:"Path to the file to be reviewed"`
	BaseText   string `json:"base_text,omitempty" jsonschema_description:"Original file content for diff calculation"`
	LMBefore   string `json:"lm_before,omitempty" jsonschema_description:"Exact text to be replaced (must match exactly including whitespace)"`
	LMAfter    string `json:"lm_after" jsonschema:"required" jsonschema_description:"Replacement text or suggestion"`
	Message    string `json:"message" jsonschema:"required" jsonschema_description:"Explanation or reason for the suggestion"`
	Severity   string `json:"severity,omitempty" jsonschema_description:"Severity level (ERROR, WARNING, INFO)"`
	SourceName string `json:"source_name,omitempty" jsonschema_description:"Name of the tool or LLM that generated this suggestion"`
}

// ConvertOutput defines the output of the tool
type ConvertOutput struct {
	ReviewdogJSON string `json:"reviewdog_json" jsonschema_description:"The converted reviewdog JSON format"`
}

// convertToReviewdog handles the MCP tool call to convert LLM suggestions to reviewdog format
func convertToReviewdog(ctx context.Context, req *mcp.CallToolRequest, input ConvertInput) (*mcp.CallToolResult, ConvertOutput, error) {
	// Marshal the input back to JSON for the suggester
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to marshal input: %v", err)}},
			IsError: true,
		}, ConvertOutput{}, nil
	}

	// Convert using the suggester library
	rdJSON, err := suggester.Convert(inputJSON, "reviewdog")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to convert: %v", err)}},
			IsError: true,
		}, ConvertOutput{}, nil
	}

	return nil, ConvertOutput{ReviewdogJSON: string(rdJSON)}, nil
}

// newMCPCommand creates the MCP subcommand
func newMCPCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Run lm-suggester as an MCP server",
		Long: `Run lm-suggester as a Model Context Protocol (MCP) server.

This allows AI assistants and LLMs to convert code suggestions to reviewdog format
through the MCP protocol over stdin/stdout.

The server provides the following tool:
  - convert_to_reviewdog: Convert LLM suggestions to reviewdog JSON format

Example usage with an MCP client:
  lm-suggester mcp

Example tool call:
  {
    "file_path": "main.go",
    "base_text": "package main\n\nfunc main() {\n\tprint(\"Hello\")\n}",
    "lm_before": "print(\"Hello\")",
    "lm_after": "fmt.Println(\"Hello\")",
    "message": "Use fmt.Println instead of print"
  }`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create MCP server
			server := mcp.NewServer(&mcp.Implementation{
				Name:    "lm-suggester",
				Version: version,
			}, nil)

			// Add the convert_to_reviewdog tool
			mcp.AddTool(server, &mcp.Tool{
				Name:        "convert_to_reviewdog",
				Description: "Convert LLM code suggestions to reviewdog JSON format. Takes file path, optional base text, optional exact match text (lm_before), replacement text (lm_after), and explanation message. Returns reviewdog-compatible JSON that can be piped to reviewdog for PR comments.",
			}, convertToReviewdog)

			// Run the server over stdin/stdout
			log.Printf("Starting MCP server (version %s)...", version)
			if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
				return fmt.Errorf("MCP server failed: %w", err)
			}

			return nil
		},
	}
}
