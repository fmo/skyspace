# SkySpace

A terminal disk-space explorer that shows what occupies your storage—and what
changed since your last scan.

**Status:** Planning. The behavior below describes the intended first version.

## Why SkySpace?

Knowing that a folder uses 40 GB is helpful. Knowing that it grew by 12 GB
since yesterday tells you where to investigate.

SkySpace combines directory navigation, size summaries, and saved scan
comparisons in one terminal interface.

## How it works

1. Choose a directory to scan.
2. Watch progress while SkySpace measures its contents.
3. Explore folders and files, sorted by size.
4. Save the completed scan locally.
5. Scan the same directory later to see what grew, shrank, appeared, or disappeared.

The first scan establishes a baseline. Each subsequent completed scan compares
against the previous saved scan of that same directory.

## Example

    ~/Projects                         Current      Change

    spreadlane/node_modules             12 GB        +3 GB
    old-project/dist                     8 GB        +5 GB
    skytui/dist                          3 GB        +1 GB
    test-data                            2 GB        unchanged

    Compared with: yesterday, 18:30

Select a directory to explore its contents and find the files responsible
for its growth.

## First-version scope

- Scan a selected directory and its descendants.
- Display file sizes and cumulative directory sizes.
- Navigate directories and sort entries by size or growth.
- Show scan progress and allow cancellation.
- Save scan snapshots locally.
- Compare snapshots for the same scan root.
- Highlight added, removed, larger, and smaller entries.
- Reveal a selected location in the system file manager.

## Comparison rules

- Compare entries by their path relative to the scan root.
- Show absolute size changes, such as `+250 MB` or `-1.2 GB`.
- Treat the first scan as a baseline, without inventing a change.
- Treat moved or renamed files as removed and added in the first version.
- Compare sizes and presence, not file contents: equal size does not mean
  identical content.
- Label unreadable or skipped locations so they are not mistaken for empty
  or deleted directories.
- Cancelled or incomplete scans do not replace the last completed baseline.

## Filesystem behavior

SkySpace reads your files and writes only its own scan snapshots.

Symbolic links are listed without following their targets. Reported sizes
represent file lengths; they may differ from physical disk usage because of
compression, sparse files, hard links, or filesystem features.

Files can change during scanning, so a scan is a best-effort observation,
not an atomic filesystem snapshot.

## Outside the first version

- Deleting files or emptying caches.
- Automatic cleanup.
- Duplicate-file detection.
- Content hashing.
- Continuous filesystem monitoring.
- Scheduled scans or cloud synchronization.

## Success criteria

SkySpace is useful when you can answer:

- Which directories consume the most space?
- What grew since my previous scan?
- Which new files account for that growth?
- Where should I investigate next?

All without modifying the files being inspected.

## Development approach

- Use Go and Bubble Tea.
- Start by defining a focused PLAN.md: goal, success criteria, flow,
  ownership, milestones, and verification.
- Work one small milestone at a time.
- Discuss design decisions before implementing substantial changes.
- The owner wants to write code and understand the flow; provide guidance
  and reviews rather than complete implementations unless requested.
- Start with one directory scan and size navigation. Add saved snapshots
  and comparisons afterward.
- Keep the first version read-only toward scanned files.
- Treat this README as a proposal, not evidence that features exist.
