# Next task: enter and leave folders

**Status:** Complete. Recorded under Unreleased in CHANGELOG.md.

## Goal

Explore a selected folder and return to its parent within the original scan root.

## Scope

- Enter opens the selected directory; Enter on a file or an empty list does nothing.
- Backspace opens the parent, stopping at the original scan root.
- Track the original root separately from the currently displayed directory.
  Normalize the starting path to an absolute, clean path for boundary checks.
- Reuse the existing summary function to load each destination synchronously.
- After successful navigation, update the path, sorted entries, and total, reset
  selection to the first row, and clear any previous navigation error.
- If loading fails, preserve the current directory, entries, total, and selection;
  show an error alongside the current list and keep the interface usable.
- Update the controls hint to include Enter and Backspace.

Preserve skip rules, sorting, empty-state behavior, arrow-key selection, and quitting.
No scrolling, background scans, caching, saved cursor history, or snapshots.

## Flow

Navigation key → determine destination → load its summaries → replace current view
on success, or preserve it and show an error on failure.

## Ownership

You write the code. Codex provides guidance and reviews unless asked to implement.

## Steps

1. Separate original-root and current-directory state and implement Enter/Backspace.
2. Update the view, controls hint, and error display.
3. Test navigation and manually walk a small fixture in the terminal.

## Done when

- Enter on a selected folder shows only that folder's immediate entries, sorted
  as before, with its path and correct total; selection resets to 0.
- Enter works through two nested directories; Backspace returns one level at a time.
- Backspace at the original root does nothing, including when launched with `.`.
- Enter on a regular file or an empty list leaves the view unchanged.
- Entering an empty folder shows the empty state and total 0; Backspace still works.
- A failed navigation preserves the current view and selection and shows an error.
  A later successful navigation clears that error.
- Focused tests cover the cases above and `go test ./...` passes.
- `go run . <fixture-directory>` supports entering, returning, selection, and quitting.

Stop here; scrolling and scan responsiveness are later tasks.
