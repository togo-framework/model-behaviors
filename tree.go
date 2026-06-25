package modelbehaviors

// Tree is a generic adjacency tree (parent/children) over comparable node ids —
// the togo answer to nested-set / closure-table helpers (django-mptt-style).
type Tree[T comparable] struct {
	parent   map[T]T
	children map[T][]T
	hasNode  map[T]bool
}

// NewTree creates an empty tree.
func NewTree[T comparable]() *Tree[T] {
	return &Tree[T]{parent: map[T]T{}, children: map[T][]T{}, hasNode: map[T]bool{}}
}

// AddRoot registers a node with no parent.
func (t *Tree[T]) AddRoot(node T) { t.hasNode[node] = true }

// AddChild attaches child under parent (registering both). A node keeps a single
// parent; re-parenting moves it.
func (t *Tree[T]) AddChild(parent, child T) {
	t.hasNode[parent] = true
	t.hasNode[child] = true
	// detach from a previous parent, if any
	if old, ok := t.parent[child]; ok {
		t.children[old] = remove(t.children[old], child)
	}
	t.parent[child] = parent
	t.children[parent] = append(t.children[parent], child)
}

// Children returns a node's direct children (in insertion order).
func (t *Tree[T]) Children(node T) []T { return append([]T{}, t.children[node]...) }

// Parent returns a node's parent and whether it has one.
func (t *Tree[T]) Parent(node T) (T, bool) {
	p, ok := t.parent[node]
	return p, ok
}

// Ancestors returns a node's ancestors from immediate parent up to the root.
func (t *Tree[T]) Ancestors(node T) []T {
	var out []T
	cur := node
	for {
		p, ok := t.parent[cur]
		if !ok {
			break
		}
		out = append(out, p)
		cur = p
	}
	return out
}

// Path returns root → … → node (the node included).
func (t *Tree[T]) Path(node T) []T {
	anc := t.Ancestors(node)
	// reverse ancestors (root first) then append node
	out := make([]T, 0, len(anc)+1)
	for i := len(anc) - 1; i >= 0; i-- {
		out = append(out, anc[i])
	}
	return append(out, node)
}

// Descendants returns all nodes beneath a node (depth-first, pre-order).
func (t *Tree[T]) Descendants(node T) []T {
	var out []T
	var walk func(n T)
	walk = func(n T) {
		for _, c := range t.children[n] {
			out = append(out, c)
			walk(c)
		}
	}
	walk(node)
	return out
}

// Depth returns a node's depth (0 for a root).
func (t *Tree[T]) Depth(node T) int { return len(t.Ancestors(node)) }

// Roots returns nodes with no parent.
func (t *Tree[T]) Roots() []T {
	var out []T
	for n := range t.hasNode {
		if _, ok := t.parent[n]; !ok {
			out = append(out, n)
		}
	}
	return out
}

func remove[T comparable](s []T, v T) []T {
	for i, x := range s {
		if x == v {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
