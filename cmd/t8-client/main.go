package main

import (
	"fmt"
	"os"
)

func main() {
	// El comando raíz se configura en root.go
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
