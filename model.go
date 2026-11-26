package main

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	mainDir  string
	nav      string
	list     []Result
	error    error
	selected int
}

func (m model) Init() tea.Cmd {
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
			m.navigate(filepath.Dir(m.nav))
			return m, nil
		case "enter":
			if len(m.list) == 0 {
				return m, nil
			}
			selected := m.list[m.selected]
			if selected.entityType != "folder" {
				return m, nil
			}
			m.navigate(filepath.Join(m.nav, selected.entityName))
			return m, nil
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		}
	}

	return m, nil
}

// navigate replaces the current directory only after it has been read successfully.
func (m *model) navigate(path string) {
	entries, err := Summarize(path)
	if err != nil {
		m.error = err
		return
	}
	m.nav = path
	m.list = entries
	m.selected = 0
	m.error = nil
}

func (m model) View() tea.View {
	var view strings.Builder
	fmt.Fprintf(&view, "Directory: %s\n\n", m.nav)
	if m.error != nil {
		fmt.Fprintf(&view, "Error: %v\n\n", m.error)
	}
	var total int64
	for index, entity := range m.list {
		total += entity.size
		if index == m.selected {
			view.WriteString("> ")
		}
		fmt.Fprintf(&view, "%s (%d) - %s\n", entity.entityName, entity.size, entity.entityType)
	}
	if len(m.list) == 0 {
		view.WriteString("No files or folders to display\n")
	}
	fmt.Fprintf(&view, "Total: %d\n\n", total)

	view.WriteString("↑/↓ select • Enter open • Backspace parent • q/Ctrl+C quit\n")

	return tea.NewView(view.String())
}
