package main

import (
	"log"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("cant get arguments")
	}

	rootDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalf("cant resolve directory path: %v", err)
	}

	results, err := Summarize(rootDir)
	if err != nil {
		log.Fatalf("cant get the summary: %v", err)
	}

	m := model{
		mainDir: rootDir,
		nav:     rootDir,
		list:    results,
	}

	program := tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		log.Fatalf("cant run the bubbletea")
	}
}
