package batch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// StdinPrompter is a Prompter that reads the user's choice from an input
// stream and writes the question to an output stream. In production it is
// wired to os.Stdin and os.Stderr, keeping the prompt off of stdout so
// scripts piping webpcvt's output are not affected.
type StdinPrompter struct {
	In  io.Reader
	Out io.Writer
}

// NewStdinPrompter returns a StdinPrompter wired to the process's stdin
// and stderr.
func NewStdinPrompter() StdinPrompter {
	return StdinPrompter{In: os.Stdin, Out: os.Stderr}
}

// Choose implements Prompter by printing question with the available
// choices, then re-reading a line until it matches one of them.
func (p StdinPrompter) Choose(question string, choices []string) (string, error) {
	scanner := bufio.NewScanner(p.In)

	for {
		fmt.Fprintf(p.Out, "%s [%s]: ", question, strings.Join(choices, "/"))

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("read choice: %w", err)
			}
			return "", errors.New("no input available to read choice")
		}

		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		for _, choice := range choices {
			if answer == choice {
				return choice, nil
			}
		}

		fmt.Fprintf(p.Out, "please answer one of: %s\n", strings.Join(choices, ", "))
	}
}
