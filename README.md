# PluckMD

Pluck for Markdown (PluckMD) is a CLI tool built on the Blocky
[Pluck](https://github.com/blocky/pluck) tool, which allows you to "pluck"
Golang type and function definitions from source code files. PluckMD uses this
functionality to programmatically replace code blocks in Markdown files to
help ensure your code documentation stays up-to-date. Along with Go, PluckMD
also supports plucking YAML code as well.

## How It Works

The `pluckmd` tool recursively searches a given directory for Markdown files.
For each file, it scans the contents looking for Markdown comments that contain
a PluckMD "directive". The format of a directive is defined as:

```
pluck("lang", "kind", "name", "source", start, end)
```

- `lang` the language of the source code to pluck (e.g., "go", "yaml")
- `kind` the kind of code block to pluck (e.g., "func", "node", "type")
- `name` the name of the code block to pluck
- `source` the file path or GitHub URL for the file containing the code to pluck
- `start` and `end` are line numbers indicating the range of lines within the
code block to include in the output

Let's demonstrate this with an example. There is a file in our repository called
`internal/pluck/goplucker.go` that contains a `GoPlucker` struct with a `Pluck` 
method. To extract the `Pluck` function and include it here in our README, we
will define a Markdown comment containing the following directive:

```
pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", -1, -1)
```

This directive tells PluckMD to pluck a Go function called `GoPlucker.Pluck`
from a file located at the given URL and to not include the function body (the 
pair `-1,-1` is used to indicate that we want to "hide" the body).

If you view the "raw" version of our README.md, you will see a comment immediately
following this text that contains our directive. Initially, the code block
below was empty, but after running `pluckmd --dir .` it was populated using the
information contained in the directive.

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

To see this in action, delete the contents of the code block but leave
the opening ticks, language identifier, and closing ticks. Then run
`pluckmd --dir .` and the code block will once again be populated with the
`GoPlucker.Pluck` function.

As we mentioned earlier, the pair `(-1, -1)` is a special case that it tells
PluckMD to exclude the entire struct or function body contents. It does add a
small comment, however, to indicate that there is hidden code. This feature is
useful when you want to walk a user through a function in logical chunks. For
example, let's look at the first "chunk" of the `GoPlucker.Pluck` function:

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

Here we might describe what the first chunk of the function is doing, before
moving on to the next bit...

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

...and the next one...

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

...until we reach the end. Finally, we might use the pair `(0, 0)` to tell
PluckMD to display the entire function body:

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

2. Install the Bearclave PluckMD CLI tool.

```bash
go install github.com/tahardi/pluckmd/cmd/pluckmd@v0.1.1
```

3. That's it! Add some pluck directives to your Markdown files and try it out!
Simply define an empty code block with a pluck directive on the line that
immediately precedes it and then run PluckMD.

## Assumptions & Limitations

- pluck directives are contained within a single-line Markdown comment (i.e., `<!-- directive -->`)
- pluck directives are on the line directly preceding the code block
- if code blocks are indented, the pluck directive has the same indentation
- the code block is marked as Golang or YAML code
- the YAML snipper does not currently support returning partial YAML components
- the YAML plucker may not support all YAML features

## Usage

- \[start, end) indicates the range of lines within the code block to include in the output
- `-1, -1` indicates that the entire code block should be excluded from the output
- `0, 0` indicates that the entire code block should be included in the output
- `"file"` can be used to print the entire contents of a file

### Local File Fetcher

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

### YAML Usage

The YAML plucker can be used to extract specific YAML components from a file.
Note that the YAML plucker does not currently support returning partial YAML
components. Thus, the `start` and `end` parameters of the pluck directive are
ignored. Let's use our `enclave-sev.yaml` file to demonstrate how to pluck YAML. 
Use the following directive to print the entire contents of the file:

```
pluck("yaml", "file", "enclave-sev.yaml", "internal/pluck/testdata/enclave-sev.yaml", -1, -1)
```

<!-- pluck("yaml", "file", "enclave-sev.yaml", "internal/pluck/testdata/enclave-sev.yaml", -1, -1) -->
```yaml
platform: "sev"
enclave:
  addr: "http://127.0.0.1:8083"
  addr_tls: "https://127.0.0.1:8444"
  args:
    domain: "bearclave.tee"
proxy:
  addr_tls: "http://127.0.0.1:8084"
  rev_addr: "http://0.0.0.0:8080"
  rev_addr_tls: "https://0.0.0.0:8443"
```

Note that the `type` and `func` kinds only apply to programming languages such
as Go. To print a component of a configuration language such as YAML, use the
`node` kind and the path to the component within the file:

```
pluck("yaml", "node", "enclave", "internal/pluck/testdata/enclave-sev.yaml", -1, -1)
```

<!-- pluck("yaml", "node", "enclave", "internal/pluck/testdata/enclave-sev.yaml", -1, -1) -->
```yaml
enclave:
  addr: "http://127.0.0.1:8083"
  addr_tls: "https://127.0.0.1:8444"
  args:
    domain: "bearclave.tee"
```

You must provide the full path to the component within the file:

```
pluck("yaml", "node", "enclave.args", "internal/pluck/testdata/enclave-sev.yaml", -1, -1)
```

<!-- pluck("yaml", "node", "enclave.args", "internal/pluck/testdata/enclave-sev.yaml", -1, -1) -->
```yaml
args:
  domain: "bearclave.tee"
```
