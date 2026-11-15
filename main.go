package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("cant get arguments")
	}

	totalSize, err := Scan(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("directory:", os.Args[1], "size:", totalSize)
}

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
			currentDir := filepath.Join(currentDir, entry.Name())
			subTotalSize, err := Scan(currentDir)
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
