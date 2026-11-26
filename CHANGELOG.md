# Changelog

## Unreleased

### Added

- Enter selected folders and return with Backspace, staying within the original
  scan root even when launched with a relative path.
- Navigation errors preserve the current view and appear alongside the list;
  successful navigation clears errors and resets selection.
- Navigation controls hint and tests for nested and empty folders, file entries,
  root boundaries, selection reset, and failure recovery.
- Bubble Tea interface with arrow-key selection and `q`/Ctrl+C to quit.
- Directory heading, entry kinds and byte sizes, combined total, controls hint,
  and an empty-state message.
- Tests for selection movement, first/last-row boundaries, and arrow keys on
  an empty list.
- List immediate files and folders by size descending, with names ascending
  for equal sizes, and print their combined total in bytes.
- Structured entry summaries containing name, kind, and byte count, with tests
  for size sorting and name tie-breaking.
- Go command to recursively total regular-file bytes beneath a directory.
- Skip symbolic links, special files, and directory metadata.
- Report scan failures without printing a partial total.
- Tests for nested and sibling folders, empty directories, external symlinks,
  missing paths, and a file supplied as the scan root.
