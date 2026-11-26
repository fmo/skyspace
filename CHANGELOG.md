# Changelog

## Unreleased

### Added

- List immediate files and folders by size descending, with names ascending
  for equal sizes, and print their combined total in bytes.
- Structured entry summaries containing name, kind, and byte count, with tests
  for size sorting and name tie-breaking.
- Go command to recursively total regular-file bytes beneath a directory.
- Skip symbolic links, special files, and directory metadata.
- Report scan failures without printing a partial total.
- Tests for nested and sibling folders, empty directories, external symlinks,
  missing paths, and a file supplied as the scan root.
