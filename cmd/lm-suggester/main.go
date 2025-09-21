package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/HikaruEgashira/lm-suggester/suggester"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	var (
		inputFile   string
		outputFile  string
		pretty      bool
		reviewdog   bool
		reporter    string
		filterMode  string
		failOnError bool
	)

	rootCmd := &cobra.Command{
		Use:   "lm-suggester",
		Short: "Convert LLM suggestions to reviewdog JSON format",
		Long: `lm-suggester transforms suggestions from LLMs and other external tools
into reviewdog-compatible JSON format for code review automation.

Automatically detects and supports both single JSON and JSONL (JSON Lines) formats.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var input io.Reader = os.Stdin
			if inputFile != "" {
				f, err := os.Open(inputFile)
				if err != nil {
					return fmt.Errorf("failed to open input file: %w", err)
				}
				defer f.Close()
				input = f
			}

			data, err := io.ReadAll(input)
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			// Convert (auto-detects JSON/JSONL)
			rdJSON, err := suggester.Convert(data, "reviewdog")
			if err != nil {
				return fmt.Errorf("failed to convert to reviewdog JSON: %w", err)
			}

			var output []byte
			if pretty {
				var jsonObj interface{}
				if err := json.Unmarshal(rdJSON, &jsonObj); err == nil {
					output, _ = json.MarshalIndent(jsonObj, "", "  ")
				} else {
					output = rdJSON
				}
			} else {
				output = rdJSON
			}

			if reviewdog {
				if err := runReviewdog(output, reporter, filterMode, failOnError); err != nil {
					return fmt.Errorf("failed to run reviewdog: %w", err)
				}
			} else if outputFile != "" {
				dir := filepath.Dir(outputFile)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
				if err := os.WriteFile(outputFile, output, 0644); err != nil {
					return fmt.Errorf("failed to write output file: %w", err)
				}
			} else {
				fmt.Print(string(output))
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input JSON file (default: stdin)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Pretty-print JSON output")
	rootCmd.Flags().BoolVar(&reviewdog, "reviewdog", false, "Run reviewdog with the output")
	rootCmd.Flags().StringVar(&reporter, "reporter", "local", "reviewdog reporter (local, github-pr-review, github-pr-check, etc.)")
	rootCmd.Flags().StringVar(&filterMode, "filter-mode", "nofilter", "reviewdog filter mode (added, diff_context, file, nofilter)")
	rootCmd.Flags().BoolVar(&failOnError, "fail-on-error", false, "Exit with non-zero code when reviewdog finds errors")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("lm-suggester %s\n", version)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built at: %s\n", date)
			fmt.Printf("  built by: %s\n", builtBy)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runReviewdog(jsonData []byte, reporter, filterMode string, failOnError bool) error {
	if _, err := exec.LookPath("reviewdog"); err != nil {
		return fmt.Errorf("reviewdog is not installed. Please install it from https://github.com/reviewdog/reviewdog")
	}

	args := []string{
		"-f=rdjson",
		fmt.Sprintf("-reporter=%s", reporter),
		fmt.Sprintf("-filter-mode=%s", filterMode),
	}

	if failOnError {
		args = append(args, "-fail-on-error=true")
	}

	cmd := exec.Command("reviewdog", args...)
	cmd.Stdin = bytes.NewReader(jsonData)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Only exit with error code if failOnError is true
			// Otherwise, ignore the error (same as pipe behavior)
			if failOnError {
				os.Exit(exitErr.ExitCode())
			}
			// Return nil to match the behavior of direct pipe
			// (reviewdog outputs to stdout even when it exits with code 1)
			return nil
		}
		return err
	}

	return nil
}