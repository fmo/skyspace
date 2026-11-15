# Changelog

## Unreleased

### Added

- Go command to recursively total regular-file bytes beneath a directory.
- Skip symbolic links, special files, and directory metadata.
- Report scan failures without printing a partial total.
- Tests for nested and sibling folders, empty directories, external symlinks,
  missing paths, and a file supplied as the scan root.
