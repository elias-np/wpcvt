package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"webpcvt/internal/batch"
	"webpcvt/internal/convert"
)

const defaultQuality = 80

// ErrVersion is returned when the user requests the version information.
var ErrVersion = errors.New("version requested")

// Options holds the validated command-line arguments before they are
// resolved into a single-file or a directory conversion job.
type Options struct {
	Input     string
	Output    string
	Quality   int
	Recursive bool
}

// Run parses arguments and starts the conversion workflow.
func Run(args []string, version string) error {
	options, err := Parse(args, version)
	if err != nil {
		return err
	}

	return dispatch(options)
}

// dispatch decides, based on what Input actually is on disk, whether to
// run a single-file conversion or a directory batch conversion.
func dispatch(options Options) error {
	info, err := os.Stat(options.Input)
	if err != nil {
		return fmt.Errorf("stat input: %w", err)
	}

	if info.IsDir() {
		return batch.Run(batch.Options{
			Root:      options.Input,
			OutputDir: options.Output,
			Quality:   options.Quality,
			Recursive: options.Recursive,
			Prompter:  batch.NewStdinPrompter(),
		})
	}

	if options.Recursive {
		return errors.New("flag -r/--recursive requires a directory input")
	}

	output := options.Output
	if output == "" {
		output = convert.DefaultOutput(options.Input)
	}

	return convert.Run(convert.Options{
		Input:   options.Input,
		Output:  output,
		Quality: options.Quality,
	})
}

// Parse validates command line arguments and returns the requested options.
func Parse(args []string, version string) (Options, error) {
	quality := defaultQuality
	recursive := false
	paths := make([]string, 0, 2)

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "-v" || arg == "--version" {
			fmt.Println("webpcvt", version)
			return Options{}, ErrVersion
		}
		if arg == "-r" || arg == "--recursive" {
			recursive = true
			continue
		}
		if arg == "-q" {
			parsedQuality, nextIndex, err := parseQuality(args, index)
			if err != nil {
				return Options{}, err
			}
			quality = parsedQuality
			index = nextIndex
			continue
		}
		if strings.HasPrefix(arg, "-q=") {
			parsedQuality, err := parseQualityValue(strings.TrimPrefix(arg, "-q="))
			if err != nil {
				return Options{}, err
			}
			quality = parsedQuality
			continue
		}
		if strings.HasPrefix(arg, "-") {
			return Options{}, fmt.Errorf("unknown flag %q", arg)
		}

		paths = append(paths, arg)
	}

	if len(paths) == 0 {
		return Options{}, errors.New("input file or directory is required")
	}
	if len(paths) > 2 {
		return Options{}, errors.New("too many arguments")
	}

	input := paths[0]
	if strings.TrimSpace(input) == "" {
		return Options{}, errors.New("input file or directory is required")
	}

	output := ""
	if len(paths) == 2 {
		output = paths[1]
	}

	return Options{
		Input:     input,
		Output:    output,
		Quality:   quality,
		Recursive: recursive,
	}, nil
}

func parseQuality(args []string, index int) (int, int, error) {
	nextIndex := index + 1
	if nextIndex >= len(args) {
		return 0, index, errors.New("quality value is required")
	}

	quality, err := parseQualityValue(args[nextIndex])
	if err != nil {
		return 0, index, err
	}

	return quality, nextIndex, nil
}

func parseQualityValue(value string) (int, error) {
	quality, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse quality: %w", err)
	}
	if quality < 0 || quality > 100 {
		return 0, errors.New("quality must be between 0 and 100")
	}

	return quality, nil
}
