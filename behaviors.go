// Package modelbehaviors provides composable model mixins for togo apps —
// timestamps, soft-deletes, slugs, tags, sortable positions, and trees
// (Rails acts_as_* / friendly_id / django-mptt for Go). Embed a struct or use
// the helpers; everything is pure-Go and dependency-free.
package modelbehaviors

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// ---- Timestamps ----

// Timestamps is an embeddable created/updated stamp.
type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Touch sets UpdatedAt to now (and CreatedAt if unset).
func (t *Timestamps) Touch() {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
}

// ---- SoftDeletes ----

// SoftDeletes is an embeddable soft-delete marker.
type SoftDeletes struct {
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Delete marks the record trashed (no-op if already trashed).
func (s *SoftDeletes) Delete() {
	if s.DeletedAt == nil {
		now := time.Now()
		s.DeletedAt = &now
	}
}

// Restore un-trashes the record.
func (s *SoftDeletes) Restore() { s.DeletedAt = nil }

// IsTrashed reports whether the record is soft-deleted.
func (s *SoftDeletes) IsTrashed() bool { return s.DeletedAt != nil }

// ---- Sluggable ----

// accentFold maps common accented letters to their ASCII base.
var accentFold = map[rune]rune{
	'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a', 'ā': 'a',
	'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e', 'ē': 'e',
	'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i', 'ī': 'i',
	'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o', 'ø': 'o', 'ō': 'o',
	'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u', 'ū': 'u',
	'ñ': 'n', 'ç': 'c', 'ß': 's', 'ý': 'y', 'ÿ': 'y',
}

// Slugify turns a string into a URL-safe slug: lowercase, ASCII-folded,
// hyphen-separated, collapsed, trimmed. e.g. "Héllo, Wörld!" → "hello-world".
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastHyphen := false
	for _, r := range s {
		if f, ok := accentFold[r]; ok {
			r = f
		}
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			lastHyphen = false
		default:
			// any separator / punctuation / unmappable letter → a single hyphen
			if !lastHyphen && b.Len() > 0 {
				b.WriteByte('-')
				lastHyphen = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// UniqueSlug returns Slugify(base), appending -2, -3, … until exists() is false.
func UniqueSlug(base string, exists func(string) bool) string {
	slug := Slugify(base)
	if slug == "" {
		slug = "n"
	}
	if !exists(slug) {
		return slug
	}
	for i := 2; ; i++ {
		cand := fmt.Sprintf("%s-%d", slug, i)
		if !exists(cand) {
			return cand
		}
	}
}

// ---- Taggable ----

// TagIndex is a concurrency-safe two-way tag index (id ↔ tags).
type TagIndex struct {
	mu    sync.RWMutex
	byID  map[string]map[string]struct{}
	byTag map[string]map[string]struct{}
}

// NewTagIndex creates an empty tag index.
func NewTagIndex() *TagIndex {
	return &TagIndex{byID: map[string]map[string]struct{}{}, byTag: map[string]map[string]struct{}{}}
}

// Tag adds one or more tags to an id.
func (t *TagIndex) Tag(id string, tags ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.byID[id] == nil {
		t.byID[id] = map[string]struct{}{}
	}
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		t.byID[id][tag] = struct{}{}
		if t.byTag[tag] == nil {
			t.byTag[tag] = map[string]struct{}{}
		}
		t.byTag[tag][id] = struct{}{}
	}
}

// Untag removes tags from an id.
func (t *TagIndex) Untag(id string, tags ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		delete(t.byID[id], tag)
		delete(t.byTag[tag], id)
		if len(t.byTag[tag]) == 0 {
			delete(t.byTag, tag)
		}
	}
}

// TagsOf returns an id's tags, sorted.
func (t *TagIndex) TagsOf(id string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return sortedKeys(t.byID[id])
}

// TaggedWith returns the ids carrying a tag, sorted.
func (t *TagIndex) TaggedWith(tag string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return sortedKeys(t.byTag[strings.ToLower(strings.TrimSpace(tag))])
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ---- Sortable ----

// Sortable maintains an explicit ordering of ids (positions).
type Sortable struct {
	mu    sync.Mutex
	order []string
}

// NewSortable creates a Sortable with an optional initial order.
func NewSortable(ids ...string) *Sortable { return &Sortable{order: append([]string{}, ids...)} }

// Add appends an id if not present.
func (s *Sortable) Add(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.indexOf(id) < 0 {
		s.order = append(s.order, id)
	}
}

// Reorder sets the order to the given ids (unknown ids are appended in order).
func (s *Sortable) Reorder(ids []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	seen := map[string]bool{}
	next := make([]string, 0, len(s.order))
	for _, id := range ids {
		if s.indexOf(id) >= 0 && !seen[id] {
			next = append(next, id)
			seen[id] = true
		}
	}
	for _, id := range s.order {
		if !seen[id] {
			next = append(next, id)
		}
	}
	s.order = next
}

// Position returns the 0-based position of an id, or -1.
func (s *Sortable) Position(id string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.indexOf(id)
}

// MoveUp swaps an id one slot earlier.
func (s *Sortable) MoveUp(id string) { s.move(id, -1) }

// MoveDown swaps an id one slot later.
func (s *Sortable) MoveDown(id string) { s.move(id, 1) }

func (s *Sortable) move(id string, delta int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.indexOf(id)
	j := i + delta
	if i < 0 || j < 0 || j >= len(s.order) {
		return
	}
	s.order[i], s.order[j] = s.order[j], s.order[i]
}

// Ordered returns the ids in their current order.
func (s *Sortable) Ordered() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string{}, s.order...)
}

func (s *Sortable) indexOf(id string) int {
	for i, x := range s.order {
		if x == id {
			return i
		}
	}
	return -1
}
