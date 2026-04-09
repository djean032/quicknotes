# Implementation Plan: QuickNotes Zettelkasten

## Phase 1: Foundation (Start Here)
**Goal:** Core data structures and configuration

### Step 1.1: pkg/config/config.go
Create the configuration package:
- `QNDir()` - returns `$QN_DIR` or `~/quicknotes`
- `Editor()` - returns `$EDITOR` or `nvim`
- `IndexPath()` - returns `QNDir/.qn_index.json`

### Step 1.2: pkg/note/note.go
Define the Note struct:
```go
type Note struct {
    ID      string    `json:"id"`      // 20060102150405
    Title   string    `json:"title"`
    Created time.Time `json:"created"`
    Tags    []string  `json:"tags"`
    Links   []string  `json:"links"`   // IDs this note links TO
}
```

### Step 1.3: pkg/note/parser.go
Create the markdown parser:
- `ParseNote(path string) (*Note, error)` - reads file, extracts frontmatter
- `ExtractLinks(content string) []string` - regex `\[\[(\d{14})\]\]`
- Return Note struct populated from file

**Test:** Write a sample note with tags and wiki-links, verify parsing works.

---

## Phase 2: Index System
**Goal:** Fast lookups and auto-rebuild

### Step 2.1: pkg/index/index.go
Create the index structure:
```go
type Index struct {
    Version      int                  `json:"version"`
    LastRebuild  time.Time            `json:"last_rebuild"`
    Notes        map[string]*Note     `json:"notes"`      // ID → Note
    TagIndex     map[string][]string  `json:"tag_index"`   // tag → []IDs
}
```

Implement:
- `NewIndex(dir string) *Index` - constructor
- `Rebuild() error` - walk QN_DIR, parse all .md files, populate maps
- `Save() error` - JSON marshal to IndexPath
- `Load() (*Index, error)` - JSON unmarshal from IndexPath
- `NoteExists(id string) bool` - check if note ID exists

### Step 2.2: pkg/index/freshness.go
Add the freshness check:
- `EnsureFresh(dir string) (*Index, error)` - check mtimes, rebuild if stale
- `isStale() bool` - compare index modtime vs newest note modtime

**Test:** Create 2-3 notes, build index manually, verify JSON output, test freshness detection.

---

## Phase 3: Commands
**Goal:** Working CLI tools

### Step 3.1: Refactor cmd/qn-new/main.go
Move existing `qn_new.go` to `cmd/qn-new/main.go`:
- Use `pkg/config` for paths
- Use `pkg/note` if needed for structure
- Add index rebuild call after creating note (trigger `EnsureFresh`)

Build: `go build -o qn-new ./cmd/qn-new`

### Step 3.2: Create cmd/qn-link/main.go
Usage: `qn-link <source_id> <target_id>`

Logic:
1. Load index via `EnsureFresh()`
2. Verify both notes exist via `NoteExists()`
3. Read source file content
4. Add `[[target_id]]` to source content (if not exists)
5. Read target file content
6. Add `[[source_id]]` to target content (if not exists)
7. Write both files back
8. Rebuild index (or update in-memory and save)

**Edge cases:** Handle duplicate links, missing IDs, same ID linking.

### Step 3.3: Create cmd/qn-suggest/main.go
Usage: `qn-suggest <note_id> [--limit N]`

Logic:
1. Load index via `EnsureFresh()`
2. Get source note from index
3. Score all other notes:
   - Shared tags: +10 per match
   - Linked neighbors: +5 for notes linked to already-linked notes
   - Keyword similarity: +1 per common word (simplified TF-IDF)
   - Recency: +1 if created within last 30 days
4. Sort by score descending
5. Print top N: `ID    Title    Score: XX [tags: N, links: M neighbors]`

Build both: `go build -o qn-link ./cmd/qn-link && go build -o qn-suggest ./cmd/qn-suggest`

---

## Phase 4: Polish & Testing
**Goal:** Robust, usable system

### Step 4.1: Error Handling
Add proper error messages for:
- Missing note files
- Invalid note ID format (must be 14 digits)
- Empty QN_DIR
- No matches for suggestions

### Step 4.2: Manual Testing Checklist
Create 5-10 test notes with various tags and links, verify:
- [ ] `qn-new "Test Note" tag1 tag2` creates file with correct format
- [ ] Index auto-rebuilds after creating notes
- [ ] `qn-link 20240101120000 20240101130000` adds `[[...]]` to both notes
- [ ] `qn-suggest 20240101120000` returns sensible rankings
- [ ] Wiki-links appear correctly in note bodies
- [ ] Index freshness check works (modify note manually, run command)

### Step 4.3: Build & Install
Add to `Makefile` or just:
```bash
go build -o qn-new ./cmd/qn-new
go build -o qn-link ./cmd/qn-link
go build -o qn-suggest ./cmd/qn-suggest
```

Then: `mv qn-* ~/bin/` or wherever your `$PATH` includes.

---

## Suggested Order of Implementation

**Day 1:** Phase 1 (config, note struct, parser) - foundation  
**Day 2:** Phase 2 (index, freshness) - the "brain"  
**Day 3:** Phase 3 (all three commands) - the interface  
**Day 4:** Phase 4 (testing, edge cases, polish) - reliability

---

## Package Structure

```
quicknotes/
├── cmd/
│   ├── qn-new/main.go      (refactored from root)
│   ├── qn-link/main.go
│   └── qn-suggest/main.go
├── pkg/
│   ├── config/config.go    (QN_DIR, EDITOR, paths)
│   ├── note/
│   │   ├── note.go         (Note struct with ID, Title, Created, Tags, Links, Content)
│   │   └── parser.go       (Parse markdown: extract frontmatter + wiki-links)
│   └── index/
│       ├── index.go        (Index struct, Load/Save, EnsureFresh, Rebuild)
│       └── suggest.go      (Scoring: tags→neighbors→keywords→recency)
└── go.mod
```

---

## Key Design Decisions

| Aspect | Decision |
|--------|----------|
| Links | Wiki-link style `[[YYYYMMDDHHMMSS]]` embedded in note body |
| Index format | JSON (`.qn_index.json` in QN_DIR root) |
| Freshness check | Sequential `EnsureFresh()` on every command startup |
| Suggestion count | Default 10, override with `--limit N` |
| qn-link lookup | Uses index to verify target note exists |

---

## Note Format

```markdown
# My Note Title

Created: 2025-01-01T12:00:00Z
Tags: programming, go, zettelkasten

---

This is the note body with a [[20250101130000]] link to another note.
And [[20250101140000]] is another link.
```

---

## Index Schema (`.qn_index.json`)

```json
{
  "version": 1,
  "last_rebuild": "2025-01-01T12:30:00Z",
  "notes": {
    "20250101120000": {
      "id": "20250101120000",
      "title": "My Note Title",
      "created": "2025-01-01T12:00:00Z",
      "tags": ["programming", "go"],
      "links": ["20250101130000", "20250101140000"]
    }
  },
  "tag_index": {
    "programming": ["20250101120000", "20250101130000"],
    "go": ["20250101120000"]
  }
}
```

---

## Command Interfaces

```
qn-new "Note Title" [tag1] [tag2]...     → Creates note, opens in EDITOR
qn-link <source_id> <target_id>          → Bidirectional wiki-link, updates index
qn-suggest <note_id> [--limit N]          → Top N related notes with scores
```

---

## Suggestion Engine Specification (Refined)

### Weighted Category System

**Default Weights (Hardcoded):**
```go
TagWeight      = 0.4  // TF-IDF weighted tag similarity
NeighborWeight = 0.3  // 2nd-degree link neighbors (Jaccard similarity)
ContentWeight  = 0.2  // TF-IDF content similarity
RecencyWeight  = 0.1  // Exponential decay
```

### Category Scoring Details

**1. Tag Similarity (0-100)**
- TF-IDF weighted cosine similarity between tag vectors
- Rare tags contribute more than common tags
- Similar tags between notes increase score

**2. Linked Neighbors (0-100)**
- Calculate 2nd-degree neighbors for both source and candidate
- Use Jaccard similarity: `shared_2nd_degree / total_unique_2nd_degree * 100`
- Cutoff at 2nd degree (no 3rd degree scoring)

**3. Content TF-IDF (0-100)**
- Pre-computed `word_index` in `.qn_index.json` (word → [note IDs])
- Calculate TF-IDF vectors for note content during indexing
- Cosine similarity between source and candidate content

**4. Recency (0-100)**
- Exponential decay: `100 * e^(-days/30)`
- Today: 100, 30 days: ~37, 90 days: ~5
- Notes > 90 days approach 0

### Final Score Calculation
```go
FinalScore = (TagScore * TagWeight) + 
             (NeighborScore * NeighborWeight) + 
             (ContentScore * ContentWeight) + 
             (RecencyScore * RecencyWeight)
```

### Exclusions
- Source note itself
- Notes already linked to source

### Output Format
```
20250101130000  Advanced Go Patterns          82.4
20250101140000  Programming Best Practices    67.1
20250101150000  Recent Project Ideas          45.3
```

### Updated Index Schema

```json
{
  "version": 1,
  "last_rebuild": "2025-01-01T12:30:00Z",
  "notes": {
    "20250101120000": {
      "id": "20250101120000",
      "title": "My Note Title",
      "created": "2025-01-01T12:00:00Z",
      "tags": ["programming", "go"],
      "links": ["20250101130000"]
    }
  },
  "tag_index": {
    "programming": ["20250101120000", "20250101130000"],
    "go": ["20250101120000"]
  },
  "word_index": {
    "function": ["20250101120000", "20250101140000"],
    "variable": ["20250101120000"],
    "struct": ["20250101130000"]
  }
}
```

### Implementation Notes for Phase 3.3

**Files to modify:**
- `pkg/index/index.go`: Add `word_index` map to Index struct
- `pkg/index/freshness.go`: Build `word_index` during `Rebuild()`
- `pkg/index/suggest.go`: New file with all scoring functions

**Algorithm:**
1. Load index via `EnsureFresh()`
2. Pre-calculate source note's 2nd-degree neighbor set once
3. For each candidate note (excluding self and already-linked):
   - Calculate TagScore using TF-IDF on tag_index
   - Calculate NeighborScore using pre-computed sets
   - Calculate ContentScore using word_index
   - Calculate RecencyScore using time decay
   - Apply weights and sum
4. Sort all candidates by final score descending
5. Print top N (default 10)
