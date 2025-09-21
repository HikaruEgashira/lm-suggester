package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/HikaruEgashira/reviewdog-converter/suggester"
)

func main() {
	payload, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var in suggester.Input
	if err := json.Unmarshal(payload, &in); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out, err := suggester.BuildRDJSON(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
