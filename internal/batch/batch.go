// Package batch converts every image in a directory to WebP, discovering
// files, resolving where converted files should go, and running the
// conversions concurrently.
package batch

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"webpcvt/internal/convert"
)

// imageExtensions are the input extensions batch discovery recognizes.
var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

const (
	layoutMirror  = "mirror"
	layoutFlatten = "flatten"

	conflictSkip      = "skip"
	conflictOverwrite = "overwrite"
)

// Prompter asks the user to pick one of a fixed set of choices, letting Run
// resolve ambiguous situations (output layout, overwrite conflicts)
// without hard-coding a particular UI.
type Prompter interface {
	Choose(question string, choices []string) (string, error)
}

// Options configures a directory conversion. Prompter must be non-nil: Run
// only calls it when a choice is actually ambiguous, but it must be ready
// to answer when that happens.
type Options struct {
	Root      string
	OutputDir string
	Quality   int
	Recursive bool
	Prompter  Prompter
}

// Run discovers images under options.Root and converts them to WebP.
func Run(options Options) error {
	root := filepath.Clean(options.Root)

	files, err := discover(root, options.Recursive)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .jpg, .jpeg or .png images found in %q", root)
	}

	jobs, err := planJobs(root, options.OutputDir, options.Prompter, files)
	if err != nil {
		return err
	}

	jobs, err = resolveConflicts(options.Prompter, jobs)
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		fmt.Fprintln(os.Stderr, "nothing to convert")
		return nil
	}

	return convertAll(jobs, options.Quality)
}

type job struct {
	input  string
	output string
}

// discover collects image files under root. Non-recursive mode only looks
// at root's direct entries; recursive mode walks the whole subtree.
func discover(root string, recursive bool) ([]string, error) {
	if !recursive {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, fmt.Errorf("read directory %q: %w", root, err)
		}

		var files []string
		for _, entry := range entries {
			if entry.IsDir() || !isImage(entry.Name()) {
				continue
			}
			files = append(files, filepath.Join(root, entry.Name()))
		}
		sort.Strings(files)
		return files, nil
	}

	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isImage(d.Name()) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory %q: %w", root, err)
	}
	sort.Strings(files)
	return files, nil
}

func isImage(name string) bool {
	return imageExtensions[strings.ToLower(filepath.Ext(name))]
}

// planJobs decides, for each discovered file, where its converted copy
// should be written. With no output directory, conversion happens in
// place. With an output directory, the layout (mirrored subfolders or a
// single flat folder) is only ambiguous when files come from more than
// one source directory, in which case the user is asked once.
func planJobs(root, outputDir string, prompter Prompter, files []string) ([]job, error) {
	if outputDir == "" {
		return inPlaceJobs(files), nil
	}

	layout := layoutFlatten
	if hasSubdirectoryFiles(root, files) {
		choice, err := prompter.Choose(
			"Images were found in subdirectories. Mirror the source folder structure inside the output directory, or flatten everything into one folder?",
			[]string{layoutMirror, layoutFlatten},
		)
		if err != nil {
			return nil, fmt.Errorf("ask output layout: %w", err)
		}
		layout = choice
	}

	if layout == layoutMirror {
		return mirrorJobs(root, outputDir, files)
	}
	return flattenJobs(outputDir, files), nil
}

func inPlaceJobs(files []string) []job {
	jobs := make([]job, len(files))
	for i, f := range files {
		jobs[i] = job{input: f, output: convert.DefaultOutput(f)}
	}
	return jobs
}

func hasSubdirectoryFiles(root string, files []string) bool {
	for _, f := range files {
		if rel, err := filepath.Rel(root, filepath.Dir(f)); err == nil && rel != "." {
			return true
		}
	}
	return false
}

func mirrorJobs(root, outputDir string, files []string) ([]job, error) {
	jobs := make([]job, len(files))
	for i, f := range files {
		relDir, err := filepath.Rel(root, filepath.Dir(f))
		if err != nil {
			return nil, fmt.Errorf("resolve relative path for %q: %w", f, err)
		}

		targetDir := outputDir
		if relDir != "." {
			targetDir = filepath.Join(outputDir, relDir)
		}

		name := filepath.Base(convert.DefaultOutput(f))
		jobs[i] = job{input: f, output: filepath.Join(targetDir, name)}
	}
	return jobs, nil
}

func flattenJobs(outputDir string, files []string) []job {
	jobs := make([]job, len(files))
	seenBy := make(map[string]string, len(files))
	for i, f := range files {
		name := filepath.Base(convert.DefaultOutput(f))
		output := filepath.Join(outputDir, name)

		if prior, ok := seenBy[output]; ok {
			fmt.Fprintf(os.Stderr, "warning: %q and %q both flatten to %q; the later conversion will overwrite the earlier one\n", prior, f, output)
		}
		seenBy[output] = f

		jobs[i] = job{input: f, output: output}
	}
	return jobs
}

// resolveConflicts asks once, only when at least one planned output
// already exists, whether to skip those files or reconvert them.
func resolveConflicts(prompter Prompter, jobs []job) ([]job, error) {
	existing := make(map[string]bool)
	for _, j := range jobs {
		if _, err := os.Stat(j.output); err == nil {
			existing[j.output] = true
		}
	}
	if len(existing) == 0 {
		return jobs, nil
	}

	choice, err := prompter.Choose(
		fmt.Sprintf("%d output file(s) already exist. Skip them or overwrite (reconvert)?", len(existing)),
		[]string{conflictSkip, conflictOverwrite},
	)
	if err != nil {
		return nil, fmt.Errorf("ask overwrite conflicts: %w", err)
	}
	if choice == conflictOverwrite {
		return jobs, nil
	}

	remaining := make([]job, 0, len(jobs))
	for _, j := range jobs {
		if !existing[j.output] {
			remaining = append(remaining, j)
		}
	}
	return remaining, nil
}

// convertAll runs jobs concurrently across a small worker pool: there is
// real parallel work to do (independent image encodes), so a pool of
// goroutines bounded by CPU count keeps machines with many cores busy
// without spawning one goroutine per file.
func convertAll(jobs []job, quality int) error {
	workerCount := runtime.NumCPU()
	if workerCount > len(jobs) {
		workerCount = len(jobs)
	}

	jobCh := make(chan job)
	errCh := make(chan error, len(jobs))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				errCh <- convertOne(j, quality)
			}
		}()
	}

	go func() {
		defer close(jobCh)
		for _, j := range jobs {
			jobCh <- j
		}
	}()

	wg.Wait()
	close(errCh)

	var errs []error
	converted := 0
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		converted++
	}

	fmt.Fprintf(os.Stderr, "converted %d file(s), %d failed\n", converted, len(errs))
	return errors.Join(errs...)
}

func convertOne(j job, quality int) error {
	if err := os.MkdirAll(filepath.Dir(j.output), 0755); err != nil {
		return fmt.Errorf("create output directory for %q: %w", j.output, err)
	}
	if err := convert.Run(convert.Options{Input: j.input, Output: j.output, Quality: quality}); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "converted %s -> %s\n", j.input, j.output)
	return nil
}
