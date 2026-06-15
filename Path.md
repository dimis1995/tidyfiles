# TidyFiles — CLI File Organizer

A command-line tool that organizes files in directories based on configurable rules. Single binary, no database, no server.

## Tech Stack

- **Language:** Go
- **CLI Framework:** `cobra`
- **Configuration:** `viper` + TOML
- **File Watching:** `fsnotify`
- **History/Undo:** JSON log file

## Project Structure

```
tidyfiles/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go          # root command, global flags
│   ├── scan.go          # directory summary
│   ├── plan.go          # dry-run preview
│   ├── run.go           # execute moves
│   ├── undo.go          # reverse last run
│   ├── watch.go         # real-time monitoring
│   └── config.go        # config management
├── engine/
│   ├── rule.go          # rule types and matching logic
│   ├── matcher.go       # glob, size, age matching
│   ├── planner.go       # builds a plan (list of moves)
│   └── executor.go      # executes moves, writes history
├── config/
│   └── config.go        # config loading, defaults
├── history/
│   └── history.go       # read/write history log
└── output/
    └── printer.go       # formatted/colored terminal output
```

## Commands

```
tidyfiles scan <dir>            # show directory summary (file count, sizes, types)
tidyfiles plan <dir>            # dry run — show what would be moved and where
tidyfiles run <dir>             # organize files according to rules
tidyfiles undo                  # reverse the last run
tidyfiles watch <dir>           # monitor directory and organize new files in real time
tidyfiles config init           # generate default config file
tidyfiles config show           # print active configuration
```

## Config

Location: `~/.config/tidyfiles/config.toml`

Falls back to sensible built-in defaults when no config file exists (group by extension into Images/, Documents/, Videos/, etc.).

```toml
[defaults]
target_dir = "~/Documents/organized"
duplicate_strategy = "skip"  # skip, rename, overwrite

[[rules]]
name = "large-pdfs"
match = "*.pdf"
min_size = "10MB"
dest = "PDFs/large"

[[rules]]
name = "old-screenshots"
match = "Screenshot*.png"
older_than = "30d"
dest = "Archive/screenshots"

[[rules]]
name = "images"
match = "*.{jpg,jpeg,png,gif,webp}"
dest = "Images"

[[rules]]
name = "videos"
match = "*.{mp4,mkv,avi,mov}"
dest = "Videos"
```

## History

Location: `~/.local/share/tidyfiles/history.json`

Records every file move with source path, destination path, rule name, and timestamp. Keeps the last 10 runs by default.

---

## Phases

### Phase 1: CLI Scaffolding + Config + Scan

**Goal:** Working CLI that loads config and prints directory summaries.

- [ ] Initialize project (`go mod init`, install `cobra`, `viper`)
- [ ] Set up `cobra` root command with global flags (`--config`, `--verbose`, `--quiet`)
- [ ] Config loading with `viper`: read TOML from `~/.config/tidyfiles/config.toml`, fall back to built-in defaults
- [ ] `config init` subcommand: generate a default config file
- [ ] `config show` subcommand: print active config (merged defaults + file)
- [ ] `scan` command: walk directory, print summary (file count by extension, total size, oldest/newest, largest files)
- [ ] Basic formatted table output in `output/printer.go`

**Deliverable:** `tidyfiles scan ~/Downloads` prints a useful summary. `tidyfiles config init` creates a starter config.

---

### Phase 2: Rule Engine + Plan

**Goal:** Files get matched against rules and a plan of moves is generated.

- [ ] Define `Rule` struct: name, glob pattern, min/max size, older/newer than, destination
- [ ] Implement glob matching (handle `*.{jpg,png}` brace expansion — `filepath.Match` doesn't support it)
- [ ] Implement size matching: parse human-readable sizes (`10MB`, `500KB`)
- [ ] Implement age matching: parse durations (`30d`, `7d`, `1h`)
- [ ] Implement planner: walk directory, test each file against rules in order (first match wins), produce `[]PlannedMove`
- [ ] Handle unmatched files: leave them or route to catch-all folder (configurable)
- [ ] Handle duplicate filenames at destination: skip/rename/overwrite per config
- [ ] `plan` command: run planner, print proposed moves grouped by rule
- [ ] `--rule` flag to filter plan output by rule name

**Deliverable:** `tidyfiles plan ~/Downloads` shows what would happen without moving anything.

---

### Phase 3: Run + Undo

**Goal:** Actually move files and support undoing the last run.

- [ ] Implement executor: take `[]PlannedMove`, move files, create destination directories as needed
- [ ] Write history log after each run
- [ ] Handle errors mid-run: log failures, continue with remaining files, report at end
- [ ] `run` command: execute plan, print results, write history
- [ ] `--dry-run` flag as alias for `plan`
- [ ] `--yes` flag to skip confirmation prompt
- [ ] Without `--yes`: show plan and ask for confirmation before executing
- [ ] Implement history read/write: JSON format, keep last N runs (configurable, default 10)
- [ ] `undo` command: read last run, reverse all moves, clean up empty directories
- [ ] `--last N` flag to undo a specific past run
- [ ] Handle undo edge cases: file modified after move, file deleted, destination gone

**Deliverable:** `tidyfiles run ~/Downloads` organizes files. `tidyfiles undo` reverses it.

---

### Phase 4: Watch Mode

**Goal:** Monitor a directory and organize new files as they appear.

- [ ] `watch` command using `fsnotify`
- [ ] Debounce file events: wait until file is stable (no writes for N seconds) before processing
- [ ] Apply rules to each new stable file, move it, log to history
- [ ] Graceful shutdown on `SIGINT`/`SIGTERM`
- [ ] Print a log line for each file organized
- [ ] Support watching multiple directories (from config)

**Deliverable:** `tidyfiles watch ~/Downloads` runs in foreground, organizing files as they arrive.

---

### Phase 5: Polish

**Goal:** Make it feel like a production CLI tool.

- [ ] Colored output (`fatih/color` or similar): green for moves, yellow for skips, red for errors
- [ ] Progress indicator for large directories during `run`
- [ ] `--json` flag on `scan` and `plan` for machine-readable output
- [ ] Shell completion via `cobra` (`tidyfiles completion bash/zsh/fish`)
- [ ] `--log-file` flag for file-based logging
- [ ] `version` subcommand with build-time info via `-ldflags`
- [ ] README.md with install instructions and usage examples
- [ ] Bonus: `stats` command showing historical data from past runs
