# Comic Archiver

A Go CLI tool that scrapes and locally archives webcomics, plus a built-in web server to browse archived comics.

## Build & Run

```bash
go build          # produces ./comic-archiver binary
go build -v ./... # verbose, used in CI
```

If you modify files under `server/html/`, rebuild CSS first:

```bash
./server/html/build.sh
go build
```

There are no tests. CI only verifies the build succeeds.

## Architecture

```
main.go           → imports cmd/, calls cmd.Execute()
cmd/              → Cobra CLI commands (archive, list, serve)
archivers/        → all scraping logic + the Comics registry
server/           → stdlib HTTP server for browsing local archives
```

**Entry point flow for `comic-archiver archive xkcd`:**

1. `cmd/root.go` initializes config (viper), reads `~/.config/comic-archiver.yaml`, and **`chdir`s into `ComicsDir`** (default: `./comics`). All subsequent file paths are relative to that directory.
2. `cmd/archive.go` spawns one goroutine per comic via `archivers.Archive()`, coordinated with a `sync.WaitGroup`.
3. `archivers/main.go: Archive()` dispatches to the appropriate archiver function based on `comic.Archiver` (a plain string, not an interface — just a chain of `if` statements).

## Adding a Comic

All comics are defined in `archivers/main.go` in the `Comics` map (key: lowercase, hyphen-separated slug). Choose an archiver type:

| Archiver | Use when |
|---|---|
| `Generic` | Single image per page, navigate via "previous" link |
| `GenericCustomStart` | Like Generic, but the start page redirects/links to actual first comic |
| `MultiImageGeneric` | Multiple images per page (e.g., xkcd's comic + alt-text image) |
| `Sequential` | Direct URL pattern with sequential numbering, no HTML scraping needed |
| Custom (in `custom.go`) | Site-specific quirks that can't be handled generically |

**Generic/MultiImageGeneric fields:**
- `StartURL`: most recent comic page (archivers walk **backwards** from newest to oldest)
- `FileMatch`: regexp with one capture group for the path/filename after `FilePrefix`
- `FilePrefix`: URL prefix prepended to the capture group to form the full image URL
- `PrevLinkMatch`: regexp with one capture group for the URL of the previous page

**Sequential fields:** `FilePrefix` (URL base) + `SeqPattern` (Go `fmt.Sprintf` pattern, e.g. `"IMG%04d.png"`) + `SeqStart`/`SeqEnd`.

## Key Conventions

- **Backwards traversal**: Generic archivers start at the newest page and follow "previous" links. `StartURL` should point to the current/latest page, not the first.
- **Resume support**: Each comic directory gets a `.last_url` file written after each successful page download. On the next run with `--continue`, the archiver skips files until it reaches that URL, then resumes.
- **`skipExisting` behavior**: Without `--continue`, hitting an existing file immediately stops (assumes everything older is already downloaded). With `--continue`, it skips existing files and keeps going.
- **Rate limiting**: 500ms sleep between page fetches and between file downloads is hardcoded in all archivers.
- **Regexp capture groups**: `FileMatch` and `PrevLinkMatch` always use capture group 1 (`[1]`) for the relevant path. `FilePrefix + match[1]` forms the full URL.
- **Server templates**: HTML templates in `server/html/` are embedded at compile time via `//go:embed`. The `server.Page` struct is the data model passed to all templates.
- **Config**: `viper` reads `~/.config/comic-archiver.yaml`. The only meaningful setting is `ComicsDir` (default: `"comics"`).
- **`hiveworks.md`**: A personal tracking list of comics to potentially add support for — not code, just a reference.
