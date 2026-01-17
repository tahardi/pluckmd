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
pluck("lang", "kind", "name", "source", start, end)
```

- `lang` can be "go" or "yaml"
- `kind` can be "file", "func", "node", or "type"
- `name` is the name of the type or function to pluck
- `source` is the GitHub URL of the source file containing the type or function
- `start` and `end` are line numbers indicating the range of lines within the
type or function body to include in the output


Let's demonstrate this with an example. There is a file in our repository called
`goplucker.go` that contains a `GoPlucker` type with a `Pluck` function. We
are going to use PluckMD to extract the `Pluck` function and include it in our
documentation.

```
pluck("function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", -1, -1)
```

<!-- pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", -1, -1) -->
```go
func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
}
```

The pair `(-1, -1)` is a special case, in that it tells `pluckmd` to exclude
the entire type or function body contents. As we can see above, it does just
that! It does add a small comment, however, to indicate that there is hidden
code. This can be very useful when you want to walk a user through a specific
function or type, but don't want to include the entire body all at once.

<!-- pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", 0, 10) -->
```go
func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	switch kind {
	case File:
		return code, nil
	case Func, Type:
		break
	case Node:
		return "", fmt.Errorf("%w: node kind not supported", ErrGoPlucker)
	default:
		return "", fmt.Errorf("%w: unrecognized kind: %v", ErrGoPlucker, kind)
	}
	// ...
}
```

Instead, we can selectively include the relevant lines for any given step...

<!-- pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", 11, 19) -->
```go
func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
	var out bytes.Buffer
	var stderr bytes.Buffer
	pick := fmt.Sprintf("%s=%s:%s", PickArg, kind, name)

	cmd := exec.CommandContext(ctx, GoPluckCmd, pick)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// ...
}
```

...as we work our way through the function...

<!-- pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", 20, 30) -->
```go
func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(
			"%w: running %s: %s",
			ErrGoPlucker,
			GoPluckCmd,
			stderr.String(),
		)
	}
	return out.String(), nil
}
```

...until we reach the end.

<!-- pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", 0, 0) -->
```go
func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	switch kind {
	case File:
		return code, nil
	case Func, Type:
		break
	case Node:
		return "", fmt.Errorf("%w: node kind not supported", ErrGoPlucker)
	default:
		return "", fmt.Errorf("%w: unrecognized kind: %v", ErrGoPlucker, kind)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	pick := fmt.Sprintf("%s=%s:%s", PickArg, kind, name)

	cmd := exec.CommandContext(ctx, GoPluckCmd, pick)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(
			"%w: running %s: %s",
			ErrGoPlucker,
			GoPluckCmd,
			stderr.String(),
		)
	}
	return out.String(), nil
}
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

#### Local File Fetcher

There are two ways that `pluckmd` tries to fetch source code files. The first is
with `GitHubFetcher`, which downloads and reads a source code file for a given
URL. The second uses `LocalFetcher`, which opens and reads a source code file
on the local system given an absolute or relative path. The important thing to
note about`LocalFetcher` is that relative paths are assumed to be relative _to
the directory in which `pluckmd` is being run._

For example, this repository has a makefile target for running `pluckmd` to
(re-)generate code blocks in our README.md. Thus, let's assume that `pluckmd`
is always run from the top-level of this repository. We would update the
`GoPluck.Pluck` URI's from the earlier example to specify the file using a
path relative to the top-level directory of our repository:

```
pluck("go", "function", "GoPlucker.Pluck", "internal/pluck/goplucker.go", -1, -1)
```

<!-- pluck("go", "function", "GoPlucker.Pluck", "internal/pluck/goplucker.go", -1, -1) -->
```go
func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
}
```

The local fetcher is useful when you want to include documentation not yet
tracked in a remote repository. For example, you may be working on a feature
branch that introduces a new function. Using remote URLs, you would have to
first commit and push the function to the remote respository before running
`pluckmd`to update your documentation. Otherwise, it would fail to find the
function, or it might pull an out-of-date version.
