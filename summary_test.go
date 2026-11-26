package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestSummary(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create the sub: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "empty"), 0o700); err != nil {
		t.Fatalf("cant create the sub: %v", err)
	}

	if err := save(filepath.Join(dir, "a.txt"), 3); err != nil {
		t.Fatalf("cant create the file: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 5); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	want := []Result{
		{entityName: "sub", entityType: "folder", size: 5},
		{entityName: "a.txt", entityType: "file", size: 3},
		{entityName: "empty", entityType: "folder", size: 0},
	}
	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant get the summary: %v", err)
	}

	if !slices.Equal(want, result) {
		t.Fatalf("want: %v, got: %v", want, result)
	}

}

func TestSummarySameSize(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create the sub: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "empty"), 0o700); err != nil {
		t.Fatalf("cant create the sub: %v", err)
	}

	if err := save(filepath.Join(dir, "b.txt"), 3); err != nil {
		t.Fatalf("cant create the file: %v", err)
	}

	if err := save(filepath.Join(dir, "a.txt"), 3); err != nil {
		t.Fatalf("cant create the file: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 5); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	want := []Result{
		{entityName: "sub", entityType: "folder", size: 5},
		{entityName: "a.txt", entityType: "file", size: 3},
		{entityName: "b.txt", entityType: "file", size: 3},
		{entityName: "empty", entityType: "folder", size: 0},
	}
	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant get the summary: %v", err)
	}

	if !slices.Equal(want, result) {
		t.Fatalf("want: %v, got: %v", want, result)
	}

}
