package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestArrowKeys(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create dir: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "sub2"), 0o700); err != nil {
		t.Fatalf("cant create dir: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant get the summary: %v", err)
	}

	m := model{
		mainDir: dir,
		list:    result,
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	got := updated.(model)

	if got.selected != 1 {
		t.Fatalf("expected: 1, got: %d", got.selected)
	}

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	got = updated.(model)

	if got.selected != 0 {
		t.Fatalf("expected: 0, got: %d", got.selected)
	}

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	got = updated.(model)

	if got.selected != 0 {
		t.Fatalf("expected: 0, got: %d", got.selected)
	}

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	got = updated.(model)

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	got = updated.(model)

	if got.selected != 1 {
		t.Fatalf("expected: 1, got: %d", got.selected)
	}
}

func TestEmptyFolder(t *testing.T) {
	dir := t.TempDir()

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant get the summary: %v", err)
	}

	m := model{
		mainDir: dir,
		list:    result,
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	got := updated.(model)

	if got.selected != 0 {
		t.Fatalf("expected: 0, got: %d", got.selected)
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	got = updated.(model)

	if got.selected != 0 {
		t.Fatalf("expected: 0, got: %d", got.selected)
	}
}

func TestEnteringFoldersInTheList(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "b.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant summirize: %v", err)
	}

	m := model{
		mainDir: dir,
		nav:     dir,
		list:    result,
	}

	if m.list[m.selected].entityName != "another" {
		t.Fatalf("sub directory should be selected, selected: %s", m.list[m.selected].entityName)
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	want := []Result{{entityName: "b.txt", size: 3, entityType: "file"}}

	if !reflect.DeepEqual(want, got.list) {
		t.Fatalf("result and list are not matching, want: %v, got: %v", want, got.list)
	}
}

func TestEnteringSecondLevelFoldersInTheList(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another", "deep"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "b.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "deep", "ab.txt"), 33); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant summirize: %v", err)
	}

	m := model{
		mainDir: dir,
		nav:     dir,
		list:    result,
	}

	if m.list[m.selected].entityName != "another" {
		t.Fatalf("sub directory should be selected, selected: %s", m.list[m.selected].entityName)
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got = updated.(model)

	if got.list[got.selected].entityName != "ab.txt" {
		t.Fatalf("expected name ab, got: %s", got.list[m.selected].entityName)
	}
}

func TestNavigation(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another", "deep"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "b.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "deep", "ab.txt"), 33); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant summirize: %v", err)
	}

	m := model{
		mainDir: dir,
		nav:     dir,
		list:    result,
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	if got.nav != filepath.Join(dir, "another") {
		t.Fatalf("got: %s, expected: %s", got.nav, filepath.Join(dir, "another"))
	}

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got = updated.(model)

	want := filepath.Join(dir, "another", "deep")

	if got.nav != want {
		t.Fatalf("got: %s, expected: %s", got.nav, want)
	}
}

func TestEnterEmptyFolder(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant summirize: %v", err)
	}

	m := model{
		mainDir: dir,
		nav:     dir,
		list:    result,
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got = updated.(model)

	if got.nav != filepath.Join(dir, "sub") {
		t.Fatalf("nav should not change")
	}
}

func TestEnterOnFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "another", "deep"), 0o700); err != nil {
		t.Fatalf("cant create folders: %v", err)
	}

	if err := save(filepath.Join(dir, "sub", "a.txt"), 3); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "b.txt"), 333); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	if err := save(filepath.Join(dir, "another", "deep", "ab.txt"), 33); err != nil {
		t.Fatalf("cant save the file: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant summirize: %v", err)
	}

	m := model{
		mainDir: dir,
		nav:     dir,
		list:    result,
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	if got.list[got.selected].entityName != "b.txt" {
		t.Fatalf("expected: b.txt, got: %s", got.list[got.selected].entityName)
	}

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got = updated.(model)

	if got.nav != filepath.Join(dir, "another") {
		t.Fatalf("nav should not change")
	}
}

func TestBackspace(t *testing.T) {
	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("cant create the sub dir: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "deep"), 0o700); err != nil {
		t.Fatalf("cant create the sub dir: %v", err)
	}

	if err := save(filepath.Join(dir, "deep", "a.txt"), 4); err != nil {
		t.Fatalf("cant create a.txt: %v", err)
	}

	result, err := Summarize(dir)
	if err != nil {
		t.Fatalf("cant summirize: %v", err)
	}

	m := model{
		mainDir: dir,
		nav:     dir,
		list:    result,
	}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	if got.list[got.selected].entityName != "a.txt" {
		t.Fatalf("expected: %v, got: %v", "a.txt", got.list[got.selected].entityName)
	}

	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})

	got = updated.(model)

	if got.list[got.selected].entityName != "deep" {
		t.Fatalf("expected: %v, got: %v", "deep", got.list[got.selected].entityName)
	}
}

func TestNavigationRoundTrip(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "b", "deep")
	if err := os.MkdirAll(child, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "a"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := save(filepath.Join(child, "file.txt"), 5); err != nil {
		t.Fatal(err)
	}
	entries, err := Summarize(root)
	if err != nil {
		t.Fatal(err)
	}
	// Select the empty folder, which sorts after b, to check cursor reset.
	m := model{mainDir: root, nav: root, list: entries, selected: 1}
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)
	if got.nav != filepath.Join(root, "a") || got.selected != 0 || len(got.list) != 0 {
		t.Fatalf("enter empty folder: %+v", got)
	}
	if !strings.Contains(got.View().Content, "Total: 0") {
		t.Fatal("empty folder must display total 0")
	}
	updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	got = updated.(model)
	if got.nav != root || got.selected != 0 || !slices.Equal(got.list, entries) {
		t.Fatalf("return from empty folder: %+v", got)
	}
	for _, want := range []string{filepath.Join(root, "b"), child} {
		updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		got = updated.(model)
		if got.nav != want || got.selected != 0 || !strings.Contains(got.View().Content, "Total: 5") {
			t.Fatalf("enter %s: %+v", want, got)
		}
	}
	for _, want := range []string{filepath.Join(root, "b"), root} {
		updated, _ = got.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
		got = updated.(model)
		if got.nav != want || got.selected != 0 {
			t.Fatalf("backspace to %s: %+v", want, got)
		}
	}
}

func TestNavigationRootBoundary(t *testing.T) {
	t.Chdir(t.TempDir())
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	m := model{mainDir: root, nav: root}
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if got := updated.(model); !reflect.DeepEqual(got, m) {
		t.Fatalf("backspace changed the original root: %+v", got)
	}
}

func TestNavigationFailureAndRecovery(t *testing.T) {
	for _, key := range []rune{tea.KeyEnter, tea.KeyBackspace} {
		t.Run(fmt.Sprint(key), func(t *testing.T) {
			root := t.TempDir()
			child := filepath.Join(root, "child")
			if err := os.Mkdir(child, 0o700); err != nil {
				t.Fatal(err)
			}
			m := model{mainDir: root, nav: root, selected: 1,
				list: []Result{{entityName: "file.txt", entityType: "file", size: 9},
					{entityName: "child", entityType: "folder"}}}
			missing := child
			if key == tea.KeyBackspace {
				m.nav = filepath.Join(child, "deep")
			}
			if err := os.Remove(missing); err != nil {
				t.Fatal(err)
			}
			updated, _ := m.Update(tea.KeyPressMsg{Code: key})
			got := updated.(model)
			if got.error == nil {
				t.Fatal("expected a navigation error")
			}
			if got.nav != m.nav || got.selected != m.selected || !slices.Equal(got.list, m.list) {
				t.Fatalf("failed navigation changed current state: %+v", got)
			}
			view := got.View().Content
			if !strings.Contains(view, "Error:") || !strings.Contains(view, "file.txt (9)") || !strings.Contains(view, m.nav) {
				t.Fatalf("error must appear alongside current view: %s", view)
			}
			if err := os.Mkdir(missing, 0o700); err != nil {
				t.Fatal(err)
			}
			updated, _ = got.Update(tea.KeyPressMsg{Code: key})
			got = updated.(model)
			if got.error != nil || got.nav != child || got.selected != 0 || len(got.list) != 0 {
				t.Fatalf("successful recovery did not replace state and clear error: %+v", got)
			}
		})
	}
}
