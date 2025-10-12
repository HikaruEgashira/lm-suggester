package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/HikaruEgashira/lm-suggester/suggester"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var in suggester.Input
	if err := json.Unmarshal(b, &in); err != nil {
		panic(err)
	}
	out, err := suggester.BuildRDJSON(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Stdout.Write(out)
}