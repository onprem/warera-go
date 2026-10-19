// Package console contains small standard-library helpers for interactive examples.
package console

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompt reads line-oriented input and writes prompts to out.
type Prompt struct {
	reader *bufio.Reader
	out    io.Writer
}

// New creates a prompt over the supplied streams.
func New(in io.Reader, out io.Writer) *Prompt {
	return &Prompt{reader: bufio.NewReader(in), out: out}
}

// Ask prompts for a non-empty line. Input is trimmed, so surrounding spaces do
// not accidentally become part of a handle, country name, or ID.
func (p *Prompt) Ask(ctx context.Context, question string) (string, error) {
	for {
		fmt.Fprintf(p.out, "%s\n> ", question)
		line, err := p.readLine(ctx)
		if err != nil {
			return "", err
		}
		if line != "" {
			return line, nil
		}
		fmt.Fprintln(p.out, "Please enter a value.")
	}
}

// Choose displays numbered options and returns the selected zero-based index,
// which can be used directly to index the original options or records slice.
func (p *Prompt) Choose(ctx context.Context, question string, options []string) (int, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("cannot choose from an empty list")
	}
	fmt.Fprintln(p.out, question)
	for i, option := range options {
		fmt.Fprintf(p.out, "  %d) %s\n", i+1, option)
	}
	for {
		fmt.Fprintf(p.out, "Choose 1-%d: ", len(options))
		line, err := p.readLine(ctx)
		if err != nil {
			return 0, err
		}
		choice, err := strconv.Atoi(line)
		if err == nil && choice >= 1 && choice <= len(options) {
			return choice - 1, nil
		}
		fmt.Fprintf(p.out, "Enter a number from 1 to %d.\n", len(options))
	}
}

func (p *Prompt) readLine(ctx context.Context) (string, error) {
	type readResult struct {
		line string
		err  error
	}
	result := make(chan readResult, 1)
	go func() {
		// If cancellation wins the select, this buffered channel still lets a
		// later stdin read report its result without blocking on the send.
		line, err := p.reader.ReadString('\n')
		result <- readResult{line: line, err: err}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case read := <-result:
		if read.err != nil && !(read.err == io.EOF && read.line != "") {
			return "", read.err
		}
		return strings.TrimSpace(read.line), nil
	}
}
