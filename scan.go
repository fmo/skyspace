package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func Scan(currentDir string) (int64, error) {
	var totalSize int64

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if !entry.Type().IsRegular() && !entry.IsDir() {
			continue
		}
		if entry.IsDir() {
			childDir := filepath.Join(currentDir, entry.Name())
			subTotalSize, err := Scan(childDir)
			if err != nil {
				return 0, err
			}
			totalSize += subTotalSize
			continue
		}
		fileInfo, err := entry.Info()
		if err != nil {
			return 0, fmt.Errorf("cant get the file info: %w", err)
		}
		totalSize += fileInfo.Size()
	}

	return totalSize, nil
}
