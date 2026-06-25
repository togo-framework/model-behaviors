package modelbehaviors

import (
	"reflect"
	"testing"
	"time"
)

func TestTimestamps(t *testing.T) {
	var ts Timestamps
	ts.Touch()
	if ts.CreatedAt.IsZero() || ts.UpdatedAt.IsZero() {
		t.Fatal("Touch should set both stamps on first call")
	}
	created := ts.CreatedAt
	time.Sleep(time.Millisecond)
	ts.Touch()
	if !ts.CreatedAt.Equal(created) {
		t.Error("CreatedAt should not change on subsequent Touch")
	}
	if !ts.UpdatedAt.After(created) {
		t.Error("UpdatedAt should advance")
	}
}

func TestSoftDeletes(t *testing.T) {
	var s SoftDeletes
	if s.IsTrashed() {
		t.Fatal("new record should not be trashed")
	}
	s.Delete()
	if !s.IsTrashed() || s.DeletedAt == nil {
		t.Fatal("Delete should set DeletedAt")
	}
	first := *s.DeletedAt
	s.Delete() // idempotent
	if !s.DeletedAt.Equal(first) {
		t.Error("repeat Delete should not change DeletedAt")
	}
	s.Restore()
	if s.IsTrashed() {
		t.Error("Restore should clear DeletedAt")
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Hello World":      "hello-world",
		"  Trim  Me  ":     "trim-me",
		"Héllo, Wörld!":    "hello-world",
		"Café & Crème":     "cafe-creme",
		"already-a-slug":   "already-a-slug",
		"Multiple   Spaces": "multiple-spaces",
		"Ñoño_test":        "nono-test",
		"!!!":              "",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUniqueSlug(t *testing.T) {
	taken := map[string]bool{"post": true, "post-2": true}
	exists := func(s string) bool { return taken[s] }
	if got := UniqueSlug("Post", exists); got != "post-3" {
		t.Errorf("UniqueSlug = %q, want post-3", got)
	}
	if got := UniqueSlug("Brand New", exists); got != "brand-new" {
		t.Errorf("UniqueSlug free = %q, want brand-new", got)
	}
}

func TestTaggable(t *testing.T) {
	idx := NewTagIndex()
	idx.Tag("a", "go", "Web", "api")
	idx.Tag("b", "go", "cli")
	if got := idx.TagsOf("a"); !reflect.DeepEqual(got, []string{"api", "go", "web"}) {
		t.Fatalf("TagsOf(a) = %v", got)
	}
	if got := idx.TaggedWith("GO"); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("TaggedWith(go) = %v", got)
	}
	idx.Untag("a", "go")
	if got := idx.TaggedWith("go"); !reflect.DeepEqual(got, []string{"b"}) {
		t.Fatalf("after untag = %v", got)
	}
	if got := idx.TagsOf("a"); !reflect.DeepEqual(got, []string{"api", "web"}) {
		t.Fatalf("TagsOf(a) after untag = %v", got)
	}
}

func TestSortable(t *testing.T) {
	s := NewSortable("a", "b", "c")
	if s.Position("b") != 1 {
		t.Fatalf("pos b = %d", s.Position("b"))
	}
	s.MoveUp("c") // a, c, b
	if got := s.Ordered(); !reflect.DeepEqual(got, []string{"a", "c", "b"}) {
		t.Fatalf("after MoveUp = %v", got)
	}
	s.MoveDown("a") // c, a, b
	if got := s.Ordered(); !reflect.DeepEqual(got, []string{"c", "a", "b"}) {
		t.Fatalf("after MoveDown = %v", got)
	}
	s.Reorder([]string{"b", "a"}) // b, a, then leftover c
	if got := s.Ordered(); !reflect.DeepEqual(got, []string{"b", "a", "c"}) {
		t.Fatalf("after Reorder = %v", got)
	}
	s.Add("d")
	if s.Position("d") != 3 {
		t.Fatalf("Add d pos = %d", s.Position("d"))
	}
	// MoveUp on the first element is a no-op.
	s.MoveUp("b")
	if s.Position("b") != 0 {
		t.Error("MoveUp on first should be a no-op")
	}
}

func TestTree(t *testing.T) {
	tr := NewTree[string]()
	tr.AddChild("root", "a")
	tr.AddChild("root", "b")
	tr.AddChild("a", "a1")
	tr.AddChild("a1", "a1x")

	if got := tr.Children("root"); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("Children(root) = %v", got)
	}
	if got := tr.Ancestors("a1x"); !reflect.DeepEqual(got, []string{"a1", "a", "root"}) {
		t.Fatalf("Ancestors(a1x) = %v", got)
	}
	if got := tr.Path("a1x"); !reflect.DeepEqual(got, []string{"root", "a", "a1", "a1x"}) {
		t.Fatalf("Path(a1x) = %v", got)
	}
	if got := tr.Depth("a1x"); got != 3 {
		t.Fatalf("Depth(a1x) = %d, want 3", got)
	}
	if got := tr.Descendants("a"); !reflect.DeepEqual(got, []string{"a1", "a1x"}) {
		t.Fatalf("Descendants(a) = %v", got)
	}
	if got := tr.Roots(); !reflect.DeepEqual(got, []string{"root"}) {
		t.Fatalf("Roots = %v", got)
	}
	// Re-parent a1 under b → it leaves a's subtree.
	tr.AddChild("b", "a1")
	if got := tr.Descendants("a"); len(got) != 0 {
		t.Fatalf("after re-parent, Descendants(a) = %v, want none", got)
	}
	if p, _ := tr.Parent("a1"); p != "b" {
		t.Fatalf("Parent(a1) = %v, want b", p)
	}
}
