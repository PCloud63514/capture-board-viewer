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

	if selection.Video == "" || selection.Audio == "" {
		fmt.Fprintln(os.Stdout, "[INFO] 선택 가능한 장치가 없어 프리뷰를 실행하지 않습니다.")
		return
	}

	if err := selector.RunPreview(ctx, selection); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
}
