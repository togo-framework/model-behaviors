# model-behaviors — usage

Pure-Go, dependency-free. Embed the structs or call the helpers.

## Timestamps
```go
type Post struct { modelbehaviors.Timestamps }
p.Touch() // CreatedAt set once; UpdatedAt every call
```

## SoftDeletes
```go
type Post struct { modelbehaviors.SoftDeletes }
p.Delete(); p.IsTrashed(); p.Restore()
```

## Slugs
```go
modelbehaviors.Slugify("Héllo, Wörld!")  // "hello-world"
modelbehaviors.UniqueSlug("Post", exists) // appends -2,-3,… until exists() is false
```

## Tags
```go
idx := modelbehaviors.NewTagIndex()
idx.Tag(id, "go", "web"); idx.Untag(id, "web")
idx.TaggedWith("go"); idx.TagsOf(id)
```

## Sortable
```go
s := modelbehaviors.NewSortable("a","b","c")
s.MoveUp("c"); s.MoveDown("a"); s.Reorder([]string{"b","a"})
s.Position("b"); s.Ordered()
```

## Tree[T]
```go
tr := modelbehaviors.NewTree[string]()
tr.AddChild("root","a"); tr.AddChild("a","a1")
tr.Ancestors("a1"); tr.Descendants("root"); tr.Depth("a1"); tr.Path("a1"); tr.Roots()
tr.AddChild("b","a1") // re-parents a1 under b
```
