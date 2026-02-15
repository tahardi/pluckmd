# Test File

This is a test file containing pluck directives and empty code blocks.

Let's test grabbing a go type definition:
<!-- pluck("go", "type", "Processor", "./processor.go", 0, 0) -->
```go
type Processor struct {
	cacher   cache.Cacher
	fetchers []fetch.Fetcher
	pluckers map[pluck.Lang]pluck.Plucker
}
```

Let's test grabbing a go function definition:
<!-- pluck("go", "function", "Processor.ProcessMarkdown", "./processor.go", 0, 0) -->
```go
func (p *Processor) ProcessMarkdown(
	ctx context.Context,
	md []byte,
) ([]byte, error) {
	// Split markdown into lines. If the markdown ends with a newline, Split
	// will return an empty string as the last element. This will cause us
	// to write an extra newline to the output. Remove the ending newline
	// (if any) so that we don't write an extra newline.
	lines := strings.Split(strings.TrimSuffix(string(md), "\n"), "\n")

	var processed bytes.Buffer
	for i := 0; i < len(lines); i++ {
		processed.WriteString(lines[i] + "\n")
		if !ContainsPluckDirective(lines[i]) {
			continue
		}

		directiveLine := lines[i]
		directive, err := NewDirective(directiveLine)
		if err != nil {
			return nil, fmt.Errorf("%w: creating directive: %w", ErrProcessor, err)
		}

		snippet, err := p.GetCodeSnippet(ctx, directive)
		if err != nil {
			return nil, fmt.Errorf("%w: getting snippet: %w", ErrProcessor, err)
		}

		codeBlockStartLine := ""
		switch directive.Lang() {
		case pluck.Go:
			codeBlockStartLine = GoCodeBlockStartLine
		case pluck.YAML:
			codeBlockStartLine = YAMLCodeBlockStartLine
		}
		err = WriteCodeBlock(&processed, directiveLine, codeBlockStartLine, snippet)
		if err != nil {
			return nil, fmt.Errorf("%w: writing code block: %w", ErrProcessor, err)
		}

		end, err := FindCodeBlockEnd(codeBlockStartLine, lines, i)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: %w: snippet uri: %s",
				ErrProcessor,
				err,
				directive.CodeSnippetURI(),
			)
		}
		i = end
	}
	return processed.Bytes(), nil
}
```

Let's test grabbing a YAML file:
<!-- pluck("yaml", "file", "nonclave-sev.yaml", "./testdata/nonclave-sev.yaml", 0, 0) -->
```yaml
platform: "sev"
nonclave:
  measurement: |
    {
      "version": 5,
      "guest_svn": 0,
      "policy": 196608,
      "family_id": "AAAAAAAAAAAAAAAAAAAAAA==",
      "image_id": "AAAAAAAAAAAAAAAAAAAAAA==",
      "vmpl": 0,
      "current_tcb": 16004667175767900164,
      "platform_info": 37,
      "signer_info": 0,
      "measurement": "FBlf3jaFK2nyrBcWbr8yIzUbHDSBTxEDOUEmUoGQc+Bh7XxH5uANAqgXrrSZLOIN",
      "host_data": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
      "id_key_digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
      "author_key_digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
      "report_id": "",
      "report_id_ma": "//////////////////////////////////////////8=",
      "reported_tcb": 16004667175767900164,
      "chip_id": "",
      "committed_tcb": 16004667175767900164,
      "current_build": 0,
      "current_minor": 58,
      "current_major": 1,
      "committed_build": 0,
      "committed_minor": 58,
      "committed_major": 1,
      "launch_tcb": 16004667175767900164,
      "cpuid_1eax_fms": 10489617
    }
```

Let's test grabbing a component of a YAML file:
<!-- pluck("yaml", "node", "nonclave.measurement", "./testdata/nonclave-sev.yaml", 0, 0) -->
```yaml
measurement: |
  {
    "version": 5,
    "guest_svn": 0,
    "policy": 196608,
    "family_id": "AAAAAAAAAAAAAAAAAAAAAA==",
    "image_id": "AAAAAAAAAAAAAAAAAAAAAA==",
    "vmpl": 0,
    "current_tcb": 16004667175767900164,
    "platform_info": 37,
    "signer_info": 0,
    "measurement": "FBlf3jaFK2nyrBcWbr8yIzUbHDSBTxEDOUEmUoGQc+Bh7XxH5uANAqgXrrSZLOIN",
    "host_data": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
    "id_key_digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "author_key_digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "report_id": "",
    "report_id_ma": "//////////////////////////////////////////8=",
    "reported_tcb": 16004667175767900164,
    "chip_id": "",
    "committed_tcb": 16004667175767900164,
    "current_build": 0,
    "current_minor": 58,
    "current_major": 1,
    "committed_build": 0,
    "committed_minor": 58,
    "committed_major": 1,
    "launch_tcb": 16004667175767900164,
    "cpuid_1eax_fms": 10489617
  }
```
