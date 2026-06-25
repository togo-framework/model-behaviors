<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/model-behaviors</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/model-behaviors"><img src="https://pkg.go.dev/badge/github.com/togo-framework/model-behaviors.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Composable model mixins for <a href="https://to-go.dev">togo</a> — soft-deletes, slugs, tags, sortable & trees.</strong></p>
</div>

## Install

```bash
togo install togo-framework/model-behaviors
```

The togo equivalent of Rails `acts_as_*` / `friendly_id` / `django-mptt` — small, pure-Go, dependency-free building blocks you embed or call directly.

## Behaviors

### Timestamps
```go
type Post struct { modelbehaviors.Timestamps; Title string }
p.Touch() // sets CreatedAt once, UpdatedAt every time
```

### SoftDeletes
```go
type Post struct { modelbehaviors.SoftDeletes }
p.Delete(); p.IsTrashed() // true
p.Restore()
```

### Slugs
```go
modelbehaviors.Slugify("Héllo, Wörld!")     // "hello-world"  (ASCII-folds accents)
slug := modelbehaviors.UniqueSlug("Post", func(s string) bool {
    return db.SlugExists(s)                  // appends -2, -3, … until free
})
```

### Tags
```go
tags := modelbehaviors.NewTagIndex()
tags.Tag("post-1", "go", "web")
tags.TaggedWith("go")   // ["post-1", …]
tags.TagsOf("post-1")   // ["go", "web"]
tags.Untag("post-1", "web")
```

### Sortable (positions)
```go
order := modelbehaviors.NewSortable("a", "b", "c")
order.MoveUp("c"); order.Reorder([]string{"b","a"})
order.Position("b"); order.Ordered()
```

### Tree (generic adjacency)
```go
tr := modelbehaviors.NewTree[string]()
tr.AddChild("root", "a"); tr.AddChild("a", "a1")
tr.Ancestors("a1")   // ["a", "root"]
tr.Descendants("root")
tr.Depth("a1")       // 1
tr.Path("a1")        // ["root", "a", "a1"]
```

---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
