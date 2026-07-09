package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"webpcvt/internal/convert"
)

const defaultQuality = 80

// Run parses arguments and starts the conversion workflow.
func Run(args []string) error {
	options, err := Parse(args)
	if err != nil {
		return err
	}

	return convert.Run(options)
}

// Parse validates command line arguments and returns conversion options.
func Parse(args []string) (convert.Options, error) {
	quality := defaultQuality
	paths := make([]string, 0, 2)

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "-q" {
			parsedQuality, nextIndex, err := parseQuality(args, index)
			if err != nil {
				return convert.Options{}, err
			}
			quality = parsedQuality
			index = nextIndex
			continue
		}
		if strings.HasPrefix(arg, "-q=") {
			parsedQuality, err := parseQualityValue(strings.TrimPrefix(arg, "-q="))
			if err != nil {
				return convert.Options{}, err
			}
			quality = parsedQuality
			continue
		}
		if strings.HasPrefix(arg, "-") {
			return convert.Options{}, fmt.Errorf("unknown flag %q", arg)
		}

		paths = append(paths, arg)
	}

	if len(paths) == 0 {
		return convert.Options{}, errors.New("input file is required")
	}
	if len(paths) > 2 {
		return convert.Options{}, errors.New("too many arguments")
	}

	input := paths[0]
	if strings.TrimSpace(input) == "" {
		return convert.Options{}, errors.New("input file is required")
	}

	output := ""
	if len(paths) == 2 {
		output = paths[1]
	}
	if strings.TrimSpace(output) == "" {
		output = defaultOutput(input)
	}

	return convert.Options{
		Input:   input,
		Output:  output,
		Quality: quality,
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

func defaultOutput(input string) string {
	ext := filepath.Ext(input)
	if ext == "" {
		return input + ".webp"
	}

	return strings.TrimSuffix(input, ext) + ".webp"
}
