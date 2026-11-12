# First task: scan one directory

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
