package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func save(path string, size int) error {
	err := os.WriteFile(path, bytes.Repeat([]byte{0xFF}, size), 0o600)
	if err != nil {
		return fmt.Errorf("cant write byte: %v", err)
	}

	return nil
}

func TestScan(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub", "deep"), 0700); err != nil {
		t.Fatalf("cant create sub dir: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "sibling"), 0700); err != nil {
		t.Fatalf("cant create sibling: %v", err)
	}

	if err := save(filepath.Join(dir, "a.txt"), 3); err != nil {
		t.Fatalf("cant save txt: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "b.txt"), 5); err != nil {
		t.Fatalf("cant save txt: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "deep", "c.txt"), 7); err != nil {
		t.Fatalf("cant save txt: %v", err)
	}

	if err := save(filepath.Join(dir, "sibling", "d.txt"), 11); err != nil {
		t.Fatalf("cant save txt: %v", err)
	}

	total, err := Scan(dir)
	if err != nil {
		t.Fatalf("cant scan dir: %v", err)
	}

	var want int64 = 26

	if total != want {
		t.Fatalf("total number of bytes under given folder should be %d but it's %d", want, total)
	}
}

func TestMetaDataNotCounted(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub", "deeper"), 0o700); err != nil {
		t.Fatalf("cant crete files: %v", err)
	}

	total, err := Scan(dir)
	if err != nil {
		t.Errorf("cant scan the files: %v", err)
	}

	if total != 0 {
		t.Errorf("total should not have the metadata but has some bytes: %d", total)
	}
}

func TestSymbolinkFilesNotCounted(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant creaste the folders: %v", err)
	}

	if err := save(filepath.Join(dir, "a.txt"), 3); err != nil {
		t.Fatalf("cant save a file: %v", err)
	}

	if err := os.Symlink(filepath.Join(dir, "a.txt"), filepath.Join(dir, "sub", "a.txt")); err != nil {
		t.Fatalf("cant create the symbolink: %v", err)
	}

	total, err := Scan(filepath.Join(dir, "sub"))
	if err != nil {
		t.Fatalf("total bytes cant be calculated: %v", err)
	}

	var want int64 = 0

	if total != want {
		t.Fatalf("total bytes: %d, want: %d", total, want)
	}
}

func TestMissingPath(t *testing.T) {
	_, err := Scan(filepath.Join(t.TempDir(), "sub"))
	if err == nil {
		t.Fatalf("should have error, due to lacking a path")
	}
}

func TestFilePathGiven(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant creaste the folders: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 3); err != nil {
		t.Fatalf("cant save a file: %v", err)
	}

	_, err := Scan(filepath.Join(dir, "sub", "a.txt"))
	if err == nil {
		t.Fatalf("error should be thrown")
	}
}

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
	result, err := Summirize(dir)
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
	result, err := Summirize(dir)
	if err != nil {
		t.Fatalf("cant get the summary: %v", err)
	}

	if !slices.Equal(want, result) {
		t.Fatalf("want: %v, got: %v", want, result)
	}

}
