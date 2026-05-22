package main

import (
	"context"
	"fmt"
	"os"

	"capture-board-selector/internal/captureboard"
)

func main() {
	ctx := context.Background()
	selector := captureboard.NewSelector()

	selection, err := selector.Select(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}

	if err := selector.RunPreview(ctx, selection); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
}
