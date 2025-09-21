package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/HikaruEgashira/lm-suggester/suggester"
)

func main() {
	// Read input from stdin
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Parse JSON input
	var in suggester.Input
	json.Unmarshal(b, &in)

	// Build reviewdog format
	out, _ := suggester.BuildRDJSON(in)

	// Output result
	fmt.Print(string(out))
}