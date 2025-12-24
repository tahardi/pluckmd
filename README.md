# PluckMD

Pluck for Markdown (PluckMD) is a CLI tool built on the Blocky
[Pluck](https://github.com/blocky/pluck) tool, which allows you to "pluck"
Golang type and function definitions from source code files. PluckMD uses this
functionality to programmatically replace code blocks in Markdown files to
help ensure your code documentation stays up-to-date.

## How It Works

PluckMD scans Markdown files looking for comments containing "pluck directives",
where a directive is defined as:

```
pluck("kind", "name", "source", start, end)
```

- `kind` can either be "type" or "function"
- `name` is the name of the type or function to pluck
- `source` is the GitHub URL of the source file containing the type or function
- `start` and `end` are line numbers indicating the range of lines within the
type or function body to include in the output


Let's demonstrate this with an example. There is a file in our repository called
`blocky.go` that contains a `BlockyPlucker` type with a `Pluck` function. We
are going to use PluckMD to extract the `Pluck` function and include it in our
documentation.

```
pluck("function", "BlockyPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/blocky.go", -1, -1)
```

<!-- pluck("function", "BlockyPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/blocky.go", -1, -1) -->
```go

```

The pair `(-1, -1)` is a special case, in that it tells `pluckmd` to exclude
the entire type or function body contents. As we can see above, it does just
that! It does add a small comment, however, to indicate that there is hidden
code. This can be very useful when you want to walk a user through a specific
function or type, but don't want to include the entire body all at once.

<!-- pluck("function", "BlockyPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/blocky.go", 0, 3) -->
```go

```

Instead, we can selectively include the relevant lines for any given step...

<!-- pluck("function", "BlockyPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/blocky.go", 4, 12) -->
```go

```

...as we work our way through the function...

<!-- pluck("function", "BlockyPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/blocky.go", 13, 23) -->
```go

```

...until we reach the end.

<!-- pluck("function", "BlockyPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/blocky.go", 0, 0) -->
```go

```

## Getting Started

1. Install the Blocky Pluck CLI tool.

```bash
go install github.com/blocky/pluck/cmd/pluck@v0.1.1
```

2. That's it! Add some pluck directives to your Markdown files and try it out!
Simply define an empty code block with a pluck directive on the line that
immediately precedes it.

### Assumptions & Limitations

- pluck directives are contained within a single-line Markdown comment (i.e., `<!-- directive -->`)
- pluck directives are on the line directly preceding the code block
- if code blocks are indented, the pluck directive has the same indentation
- the code block contains Golang code

### Info

- \[start, end) indicates the range of lines within the code block to include in the output
- `-1, -1` indicates that the entire code block should be excluded from the output
- `0, 0` indicates that the entire code block should be included in the output
