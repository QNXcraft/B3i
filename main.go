package main

import (
	"fmt"
	"os"

	"github.com/user/b3i/cmd/b3i"
)

func main() {
	if err := b3i.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
