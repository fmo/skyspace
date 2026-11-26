package main

import (
	tea "charm.land/bubbletea/v2"

	"cmp"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

type model struct {
	mainDir  string
	total    int64
	nav      string
	list     []Result
	error    error
	selected int
}

func (m model) Init() tea.Cmd {
	m.selected = 0

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down":
			if m.selected < len(m.list)-1 {
				m.selected++
			}
		case "backspace":
			if m.mainDir == m.nav {
				return m, nil
			}
			parent := filepath.Dir(m.nav)
			result, err := Summirize(parent)
			if err != nil {
				m.error = err
				return m, nil
			}
			m.nav = parent
			m.list = result
			m.selected = 0
			m.error = nil
			return m, nil
		case "enter":
			if len(m.list) == 0 {
				return m, nil
			}
			selected := m.list[m.selected]
			if selected.entityType != "folder" {
				return m, nil
			}
			nav := filepath.Join(m.nav, selected.entityName)
			result, err := Summirize(nav)
			if err != nil {
				m.error = err
				return m, nil
			}
			m.nav = nav
			m.selected = 0
			m.list = result
			m.error = nil
			return m, nil
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		}
	case summary:
		m.list = msg
		var total int64
		for _, entry := range msg {
			total += entry.size
		}
		m.total = total
	}

	return m, nil
}

func (m model) View() tea.View {
	summary := fmt.Sprintf("Directory: %s\n\n", m.nav)
	if m.error != nil {
		summary += fmt.Sprintf("Error: %v\n\n", m.error)
	}
	var total int64
	for ind, entity := range m.list {
		total += entity.size
		if ind == m.selected {
			summary += "> "
		}
		summary += fmt.Sprintf("%s (%d) - %s\n", entity.entityName, entity.size, entity.entityType)
	}
	if len(m.list) == 0 {
		summary += fmt.Sprintf("No files or folders to display\n")
	}
	summary += fmt.Sprintf("Total: %d\n\n", total)

	summary += "↑/↓ select • Enter open • Backspace parent • q/Ctrl+C quit\n"

	return tea.NewView(summary)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("cant get arguments")
	}

	rootDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalf("cant resolve directory path: %v", err)
	}

	results, err := Summirize(rootDir)
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

type summary []Result
type errorType error

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
