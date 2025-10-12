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
	Message string `json:"message" jsonschema_description:"Result message or error details"`
}

// suggest handles the MCP tool call to suggest code changes via reviewdog
func suggest(ctx context.Context, req *mcp.CallToolRequest, input SuggestInput) (*mcp.CallToolResult, SuggestOutput, error) {
	// Marshal the input back to JSON for the suggester
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to marshal input: %v", err)}},
			IsError: true,
		}, SuggestOutput{Success: false, Message: fmt.Sprintf("failed to marshal input: %v", err)}, nil
	}

	// Convert using the suggester library
	rdJSON, err := suggester.Convert(inputJSON, "reviewdog")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to convert: %v", err)}},
			IsError: true,
		}, SuggestOutput{Success: false, Message: fmt.Sprintf("failed to convert: %v", err)}, nil
	}

	// Determine reporter
	reporter := "local"
	if input.Reporter != "" {
		reporter = input.Reporter
	}

	// Run reviewdog with the converted JSON
	if err := runReviewdog(rdJSON, reporter, "nofilter", false); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("reviewdog failed: %v", err)}},
			IsError: true,
		}, SuggestOutput{Success: false, Message: fmt.Sprintf("reviewdog failed: %v", err)}, nil
	}

	return nil, SuggestOutput{
		Success: true,
		Message: fmt.Sprintf("Successfully posted suggestion to %s for %s", reporter, input.FilePath),
	}, nil
}

// newMCPCommand creates the MCP subcommand
func newMCPCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Run lm-suggester as an MCP server",
		Long: `Run lm-suggester as a Model Context Protocol (MCP) server.

This allows AI assistants and LLMs to post code review suggestions directly
through the MCP protocol over stdin/stdout.

The server provides the following tool:
  - suggest: Post code review suggestions via reviewdog

Example usage with an MCP client:
  lm-suggester mcp

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
			// Create MCP server
			server := mcp.NewServer(&mcp.Implementation{
				Name:    "lm-suggester",
				Version: version,
			}, nil)

			// Add the suggest tool
			mcp.AddTool(server, &mcp.Tool{
				Name:        "suggest",
				Description: "Post code review suggestions via reviewdog. Converts LLM suggestions to reviewdog format and automatically runs reviewdog to post them. Takes file path, optional base text, optional exact match text (lm_before), replacement text (lm_after), explanation message, and optional reporter (local, github-pr-review, etc.). Returns success status and message.",
			}, suggest)

			// Run the server over stdin/stdout
			log.Printf("Starting MCP server (version %s)...", version)
			if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
				return fmt.Errorf("MCP server failed: %w", err)
			}

			return nil
		},
	}
}
