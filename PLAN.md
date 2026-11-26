# Next task: select entries in a terminal list

**Status:** Ready to start.

## Goal

Show the existing size summaries in a minimal Bubble Tea interface with a
keyboard-controlled selection. This establishes the interaction needed for
directory navigation.

## Scope

- Keep accepting one directory argument and reuse the existing summary function.
- Finish scanning before starting the interface; report scan errors as before.
- Display the root path, names, kinds, raw byte sizes, and combined total.
- Preserve the existing sort order and mark the selected row with `>`.
- Move selection with the up/down arrow keys.
- Quit with `q` or Ctrl+C and restore the terminal.
- Show a short controls hint.

No opening folders, scrolling, styling framework, background scans, progress,
snapshots, or comparisons. Use a small fixture that fits in the terminal.

## Flow

Directory argument → collect sorted summaries → initialize model → render list
→ handle keys → update selection → render again.

## Ownership

You write the code. Codex provides guidance and reviews unless asked to implement.

## Steps

1. Add Bubble Tea and define a model holding the root, summaries, and cursor.
2. Implement rendering and up/down/quit handling without changing scan logic.
3. Test selection behavior and manually check the terminal interaction.

## Done when

- A fixture with `sub` (5 bytes), `a.txt` (3), and `empty` (0) displays those
  rows in that order, their kinds, the root path, and a total of 8 bytes.
- The first row is selected initially; down selects the next row and up returns.
- Up at the first row and down at the last row leave the selection in place.
- An empty root displays an empty-state message and total 0; arrow keys are safe.
- `q` and Ctrl+C exit cleanly.
- A scan failure exits nonzero before the interface starts.
- Focused tests cover cursor movement, boundaries, and an empty list; all tests
  pass with `go test ./...`.
- `go run . <fixture-directory>` supports the interactions above in a terminal.

Stop here; entering and leaving folders is a separate task.
