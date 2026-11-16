package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("cant get arguments")
	}

	result, err := Summirize(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	summary := ""
	var total int64 = 0
	for _, one := range result {
		summary += fmt.Sprintf("%s (%d)\n", one.entityName, one.size)
		total += one.size
	}

	fmt.Printf("%s\nTotal: %d\n", summary, total)
}

type Result struct {
	entityName string
	size       int64
	entityType string
}

func Summirize(rootDir string) ([]Result, error) {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return []Result{}, err
	}

	result := make([]Result, 0, len(entries))
	for _, entry := range entries {
		if !entry.Type().IsRegular() && !entry.Type().IsDir() {
			continue
		}
		if entry.Type().IsRegular() {
			fileInfo, err := entry.Info()
			if err != nil {
				return []Result{}, err
			}
			result = append(result, Result{entityName: entry.Name(), size: fileInfo.Size(), entityType: "file"})
			continue
		}
		total, err := Scan(filepath.Join(rootDir, entry.Name()))
		if err != nil {
			return []Result{}, err
		}
		result = append(result, Result{entityName: entry.Name(), size: total, entityType: "folder"})
	}

	slices.SortFunc(result, func(x, y Result) int {
		if cmp.Compare(y.size, x.size) == 0 {
			return cmp.Compare(x.entityName, y.entityName)
		}
		return cmp.Compare(y.size, x.size)
	})

	return result, nil

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
