# First task: scan one directory

**Status:** Complete. Recorded under Unreleased in CHANGELOG.md.

## Goal

Build a small Go command that accepts one directory path and prints the total
file bytes beneath it. This proves the scanning foundation before adding a UI.

## Scope

- Accept exactly one directory path as a command-line argument.
- Walk its descendants and sum regular-file lengths in bytes.
- Do not follow symbolic links or count directory metadata.
- Print the total only after a successful scan.
- Return a clear error and nonzero exit code for invalid input or a failed walk.

No Bubble Tea, navigation, progress display, snapshots, or comparisons in this task.

## Flow

Directory argument → validate directory → walk files → sum bytes → print total.
Keep scanning in a function separate from command-line handling so it can be tested.

## Ownership

You write the code. Codex provides guidance and reviews unless asked to implement.

## Steps

1. Initialize a Go module and add the scan function.
2. Add the command-line entry point and error handling.
3. Add focused tests and run the checks below.

## Done when

- A fixture with a 3-byte file and a nested 5-byte file reports exactly `8` bytes.
- An empty directory reports `0` bytes.
- A symlink to a file outside the scan root adds no bytes.
- A missing path or a file supplied as the root exits nonzero with a useful error.
- Walk errors propagate to the caller instead of producing a successful partial total.
- `go test ./...` passes and `go run . <fixture-directory>` prints the expected total.

Stop here once these checks pass; choose the next small task afterward.

---

# Next task: show what occupies the root directory

**Status:** Ready to start.

## Goal

Print the scan root's immediate files and folders, largest first, so you can
identify where its bytes are concentrated.

## Scope

- Return one summary per immediate regular file or directory: name, kind, bytes.
- For directories, use the recursive total of their regular files.
- Sort by bytes descending; break ties by name ascending.
- Print one row per summary and the overall total, using raw bytes.
- Preserve the existing skip rules and error behavior. Print results only after
  all entries have been measured successfully.

No interactive UI, navigation, human-readable units, snapshots, or comparisons.

## Flow

Read root entries → measure each child → collect summaries → sort → print rows
and total. Keep collection and sorting separate from terminal output for testing.

## Ownership

You write the code. Codex provides guidance and reviews unless asked to implement.

## Steps

1. Define a small summary type and collect the root's child summaries, reusing Scan.
2. Sort the summaries and print the rows plus the total from the command.
3. Add a focused test and run the checks below.

## Done when

- A root with a 3-byte `a.txt`, a `sub` directory containing 5 bytes, and an empty
  `empty` directory produces rows in this order: `sub` (5), `a.txt` (3), `empty` (0).
- Descendants appear only through their parent directory's total, not as extra rows.
- Equal-size entries are ordered by name ascending.
- Symlinks and special files produce no rows or bytes.
- An empty root produces no entry rows and a total of 0.
- A failed child scan returns an error and prints no successful partial report.
- Existing tests and the new summary test pass with `go test ./...`.
- `go run . <fixture-directory>` prints the expected ordered rows and total of 8.

Stop here; directory navigation is a later task.
