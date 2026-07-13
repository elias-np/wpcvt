package convert

import (
	"fmt"
	"path/filepath"
	"strings"

	"webpcvt/internal/webp"
)

// Options contains the validated inputs for one image conversion.
type Options struct {
	Input   string
	Output  string
	Quality int
}

// Run converts one image to WebP using the configured options.
func Run(options Options) error {
	if err := webp.EncodeFile(options.Input, options.Output, options.Quality); err != nil {
		return fmt.Errorf("convert %q to %q: %w", options.Input, options.Output, err)
	}

	return nil
}

// DefaultOutput derives an output path from input by swapping its
// extension for .webp, or appending .webp when input has no extension.
func DefaultOutput(input string) string {
	ext := filepath.Ext(input)
	if ext == "" {
		return input + ".webp"
	}

	return strings.TrimSuffix(input, ext) + ".webp"
}
