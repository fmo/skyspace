package main

import (
	"cmp"
	"os"
	"path/filepath"
	"slices"
)

type Result struct {
	entityName string
	size       int64
	entityType string
}

func Summarize(rootDir string) ([]Result, error) {
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
		if order := cmp.Compare(y.size, x.size); order != 0 {
			return order
		}
		return cmp.Compare(x.entityName, y.entityName)
	})

	return result, nil
}
