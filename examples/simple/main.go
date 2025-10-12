package main

import (
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
	out, err := suggester.Convert(b, "reviewdog")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Stdout.Write(out)
}