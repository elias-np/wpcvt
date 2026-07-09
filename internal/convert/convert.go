package convert

import (
	"fmt"

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
