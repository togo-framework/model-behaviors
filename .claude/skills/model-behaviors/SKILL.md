---
name: model-behaviors
description: Add common model behaviors to a togo app with the model-behaviors plugin — soft-deletes, timestamps, URL slugs (unique), tags, sortable positions, and tree/hierarchy helpers.
---

# togo model-behaviors

Embed or call these pure-Go mixins.

- **Timestamps:** embed `modelbehaviors.Timestamps`; call `.Touch()` on save.
- **SoftDeletes:** embed `modelbehaviors.SoftDeletes`; `.Delete()/.Restore()/.IsTrashed()`; filter out trashed rows in queries.
- **Slugs:** `Slugify(title)` for URL slugs; `UniqueSlug(title, existsFn)` to guarantee uniqueness (appends -2,-3,…).
- **Tags:** `NewTagIndex()` → `Tag/Untag/TaggedWith/TagsOf` (case-insensitive).
- **Sortable:** `NewSortable(ids...)` → `Reorder/MoveUp/MoveDown/Position/Ordered`.
- **Tree[T]:** `NewTree[T]()` → `AddChild/Ancestors/Descendants/Depth/Path/Roots` for categories/comments/menus.

Persist the underlying fields (DeletedAt, slug, position, parent_id) in your schema; these helpers are the in-memory/logic layer.
