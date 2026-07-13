package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"webpcvt/internal/testutil"
)

// stubPrompter returns queued responses in order and fails the test's
// expectations by erroring if it is asked more questions than expected.
type stubPrompter struct {
	choices []string
	calls   int
}

func (s *stubPrompter) Choose(question string, _ []string) (string, error) {
	if s.calls >= len(s.choices) {
		return "", fmt.Errorf("unexpected prompt: %q", question)
	}
	choice := s.choices[s.calls]
	s.calls++
	return choice, nil
}

func TestRunNonRecursiveIgnoresSubdirectories(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	mustMkdir(t, sub)
	testutil.WritePNG(t, filepath.Join(dir, "top.png"))
	testutil.WritePNG(t, filepath.Join(sub, "nested.png"))

	if err := Run(Options{Root: dir, Quality: 80, Prompter: &stubPrompter{}}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, filepath.Join(dir, "top.webp"))
	if _, err := os.Stat(filepath.Join(sub, "nested.webp")); !os.IsNotExist(err) {
		t.Fatalf("nested.webp should not have been created in non-recursive mode")
	}
}

func TestRunRecursiveConvertsSubdirectories(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	mustMkdir(t, sub)
	testutil.WritePNG(t, filepath.Join(dir, "top.png"))
	testutil.WritePNG(t, filepath.Join(sub, "nested.png"))

	options := Options{Root: dir, Quality: 80, Recursive: true, Prompter: &stubPrompter{}}
	if err := Run(options); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, filepath.Join(dir, "top.webp"))
	testutil.AssertWebPFile(t, filepath.Join(sub, "nested.webp"))
}

func TestRunFlattenWithoutSubdirectoriesSkipsPrompt(t *testing.T) {
	dir := t.TempDir()
	testutil.WritePNG(t, filepath.Join(dir, "a.png"))
	testutil.WritePNG(t, filepath.Join(dir, "b.png"))
	outDir := filepath.Join(t.TempDir(), "out")

	prompter := &stubPrompter{}
	options := Options{Root: dir, OutputDir: outDir, Quality: 80, Prompter: prompter}
	if err := Run(options); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if prompter.calls != 0 {
		t.Fatalf("prompter called %d times, want 0", prompter.calls)
	}

	testutil.AssertWebPFile(t, filepath.Join(outDir, "a.webp"))
	testutil.AssertWebPFile(t, filepath.Join(outDir, "b.webp"))
}

func TestRunAsksLayoutAndMirrors(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	mustMkdir(t, sub)
	testutil.WritePNG(t, filepath.Join(dir, "top.png"))
	testutil.WritePNG(t, filepath.Join(sub, "nested.png"))
	outDir := filepath.Join(t.TempDir(), "out")

	prompter := &stubPrompter{choices: []string{"mirror"}}
	options := Options{Root: dir, OutputDir: outDir, Quality: 80, Recursive: true, Prompter: prompter}
	if err := Run(options); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if prompter.calls != 1 {
		t.Fatalf("prompter called %d times, want 1", prompter.calls)
	}

	testutil.AssertWebPFile(t, filepath.Join(outDir, "top.webp"))
	testutil.AssertWebPFile(t, filepath.Join(outDir, "sub", "nested.webp"))
}

func TestRunAsksLayoutAndFlattens(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	mustMkdir(t, sub)
	testutil.WritePNG(t, filepath.Join(dir, "top.png"))
	testutil.WritePNG(t, filepath.Join(sub, "nested.png"))
	outDir := filepath.Join(t.TempDir(), "out")

	prompter := &stubPrompter{choices: []string{"flatten"}}
	options := Options{Root: dir, OutputDir: outDir, Quality: 80, Recursive: true, Prompter: prompter}
	if err := Run(options); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, filepath.Join(outDir, "top.webp"))
	testutil.AssertWebPFile(t, filepath.Join(outDir, "nested.webp"))
	if _, err := os.Stat(filepath.Join(outDir, "sub")); !os.IsNotExist(err) {
		t.Fatalf("sub directory should not have been created when flattening")
	}
}

func TestRunAsksOverwriteAndSkips(t *testing.T) {
	dir := t.TempDir()
	testutil.WritePNG(t, filepath.Join(dir, "a.png"))
	testutil.WritePNG(t, filepath.Join(dir, "b.png"))

	sentinel := []byte("existing")
	if err := os.WriteFile(filepath.Join(dir, "a.webp"), sentinel, 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	prompter := &stubPrompter{choices: []string{"skip"}}
	if err := Run(Options{Root: dir, Quality: 80, Prompter: prompter}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if prompter.calls != 1 {
		t.Fatalf("prompter called %d times, want 1", prompter.calls)
	}

	data, err := os.ReadFile(filepath.Join(dir, "a.webp"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "existing" {
		t.Fatalf("a.webp was overwritten despite skip choice")
	}

	testutil.AssertWebPFile(t, filepath.Join(dir, "b.webp"))
}

func TestRunAsksOverwriteAndReconverts(t *testing.T) {
	dir := t.TempDir()
	testutil.WritePNG(t, filepath.Join(dir, "a.png"))

	if err := os.WriteFile(filepath.Join(dir, "a.webp"), []byte("existing"), 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	prompter := &stubPrompter{choices: []string{"overwrite"}}
	if err := Run(Options{Root: dir, Quality: 80, Prompter: prompter}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, filepath.Join(dir, "a.webp"))
}

func TestRunReturnsErrorWhenNoImagesFound(t *testing.T) {
	dir := t.TempDir()

	if err := Run(Options{Root: dir, Quality: 80, Prompter: &stubPrompter{}}); err == nil {
		t.Fatal("Run returned nil error")
	}
}

func TestRunConvertsManyFilesConcurrently(t *testing.T) {
	dir := t.TempDir()
	const count = 12
	for i := 0; i < count; i++ {
		testutil.WritePNG(t, filepath.Join(dir, fmt.Sprintf("img%02d.png", i)))
	}

	if err := Run(Options{Root: dir, Quality: 80, Prompter: &stubPrompter{}}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	for i := 0; i < count; i++ {
		testutil.AssertWebPFile(t, filepath.Join(dir, fmt.Sprintf("img%02d.webp", i)))
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
}
