# Test File

This is a test file containing pluck directives and empty code blocks.

Let's test grabbing a go type definition:
<!-- pluck("go", "type", "Processor", "./processor.go", 0, 0) -->
```go

```

Let's test grabbing a go function definition:
<!-- pluck("go", "function", "Processor.ProcessMarkdown", "./processor.go", 0, 0) -->
```go

```

Let's test grabbing a YAML file:
<!-- pluck("yaml", "file", "nonclave-sev.yaml", "./testdata/nonclave-sev.yaml", 0, 0) -->
```yaml

```

Let's test grabbing a component of a YAML file:
<!-- pluck("yaml", "node", "nonclave.measurement", "./testdata/nonclave-sev.yaml", 0, 0) -->
```yaml

```
