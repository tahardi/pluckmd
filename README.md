# PluckMD

Pluck for Markdown (PluckMD) is a CLI tool built on the Blocky
[Pluck](https://github.com/blocky/pluck) tool, which allows you to "pluck"
Golang type and function definitions from source code files. PluckMD uses this
functionality to programmatically replace code blocks in Markdown files to
help ensure your code documentation stays up-to-date. Along with Go, PluckMD
also supports plucking YAML code as well.

- [**How It Works**](#how-it-works)
- [**Installation**](#installation)
- [**Reference**](#reference)
  - [**CLI Usage**](#cli-usage)
  - [**Directives**](#directives)
  - [**YAML**](#yaml)
- [**Assumptions & Limitations**](#assumptions--limitations)

## How It Works

PluckMD recursively searches a given directory for Markdown files.
For each file, it scans the contents looking for Markdown comments that contain
a PluckMD "directive", where a directive looks something like:

```
pluck("lang", "kind", "name", "source", start, end)
```

Let's look at a concrete example. There is a file in our repository called
`goplucker.go` that contains a `GoPlucker` struct with a `Pluck` method. To
extract the `Pluck` function and include it here in our README, we 
define a Markdown comment containing the following directive:

```
pluck("go", "function", "GoPlucker.Pluck", "https://github.com/tahardi/pluckmd/blob/main/internal/pluck/goplucker.go", -1, -1)
```

This directive tells PluckMD to pluck a Go function called `GoPlucker.Pluck`
from a file located at the given URL and to hide the function body (the pair 
`-1,-1` is used to indicate that we don't want to display the body).

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

Check for yourself. Delete the contents of the code block but leave
the opening ticks, language identifier, and closing ticks. Then run
`pluckmd --dir .` and the code block will once again be populated with the
`GoPlucker.Pluck` function.

The `start` and `end` fields can be used to display only a portion of
the struct or function body. This is useful when you want to walk a user through
the logical sections of a function or struct. For example, let's look at the
first part of the `GoPlucker.Pluck` function:

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

Here we might describe what this first section of the function is doing, before
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
## Installation

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

## Reference

Here we detail the various options and operating modes supported by the PluckMD
CLI tool.

### CLI Usage
Below is an example illustrating how to use the `pluckmd` tool:

```bash
pluckmd --dir . \
  --ignore-dir testdata/ \
  --ignore-dir .github/ \
  --timeout 120
```

Note you must pass a `--dir` argument so `pluckmd` knows which directory to
process, whereas the other arguments are optional. You can use `pluckmd --help`
for a full list and description of supported arguments.

### Directives

Directives are used to tell `pluckmd` what code to fetch, where to fetch it from,
and how to render it. The format for a directive is:

```
pluck("lang", "kind", "name", "source", start, end)
```

#### Lang

Currently, PluckMD supports plucking code for the following languages:

- `go`
- `yaml`

#### Kind

This field indicates what kind of code block is being plucked.  

- `file` used to read an entire file. Can be used with both `go` and `yaml`.
- `function` used to read a function. Only used with `go`.
- `node` used to read a node component. Only used with `yaml`.
- `type` used to read a type definition. Only used with `go`.

#### Name

The name of the code or file to be plucked.

- `GoPlucker` the name of a standalone function or type
- `GoPlucker.Pluck` functions defined on structs are named `<struct>.<func>`
- `enclave-sev.yaml` when readings files the entire filename must be specified
- `enclave.args.domain` for nested YAML nodes you must specify the node path

#### Source

Currently, PluckMD supports fetching source code from:

- GitHub
- Local Files

The local fetcher reads local files given an absolute or relative path. Note
that relative paths are assumed to be relative _to the directory in which 
`pluckmd` is run._

For example, this repository has a makefile target for running `pluckmd` to
(re-)generate code blocks in our README.md. Since `pluckmd` is run from the 
top-level of this repository, we must use a path relative to the top-level
directory of our repository:

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
tracked in a remote repository, such as when you are working on a feature
branch that introduces a new function.

#### Start & End

The pair `[start, end)` is used to display a portion of the plucked type or
function body. Currently, this feature is only supported for Golang code and
not YAML. Depending on whether the displayed code is at the beginning, middle, or
end of the body, PluckMD will add lines containing `// ...` to indicate to the
reader that part of the body is hidden.

There are two special pairs to be aware of:
- `-1, -1` indicates that the entire code block should be excluded from the output
- `0, 0` indicates that the entire code block should be included in the output

### YAML

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

## Assumptions & Limitations

- pluck directives are contained within a single-line Markdown comment (i.e., `<!-- directive -->`)
- pluck directives are on the line directly preceding the code block
- if code blocks are indented, the pluck directive has the same indentation
- the code block is marked as Golang or YAML code
- the YAML snipper does not currently support returning partial YAML components
- the YAML plucker may not support all YAML features
