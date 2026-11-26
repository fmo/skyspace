package main

import (
	"bytes"
	"fmt"
	"os"
)

func save(path string, size int) error {
	err := os.WriteFile(path, bytes.Repeat([]byte{0xFF}, size), 0o600)
	if err != nil {
		return fmt.Errorf("cant write byte: %v", err)
	}

	return nil
}
