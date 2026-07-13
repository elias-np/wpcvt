package cli

import (
	"path/filepath"
	"testing"

	"webpcvt/internal/testutil"
)

const testVersion = "0.1.0-test"

func TestParseWithExplicitOutput(t *testing.T) {
	options, err := Parse([]string{"image.jpg", "-q", "85", "out.webp"}, testVersion)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	assertOptions(t, options, Options{
		Input:   "image.jpg",
		Output:  "out.webp",
		Quality: 85,
	})
}

func TestParseLeavesOutputEmptyWhenOmitted(t *testing.T) {
	options, err := Parse([]string{"assets/image.png", "-q", "70"}, testVersion)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	assertOptions(t, options, Options{
		Input:   "assets/image.png",
		Output:  "",
		Quality: 70,
	})
}

func TestParseUsesDefaultQuality(t *testing.T) {
	options, err := Parse([]string{"image.jpg"}, testVersion)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if options.Quality != defaultQuality {
		t.Fatalf("Quality = %d, want %d", options.Quality, defaultQuality)
	}
}

func TestParseVersionFlag(t *testing.T) {
	for _, flag := range []string{"-v", "--version"} {
		_, err := Parse([]string{flag}, testVersion)
		if err != ErrVersion {
			t.Fatalf("Parse(%q) error = %v, want ErrVersion", flag, err)
		}
	}
}

func TestParseRecursiveFlag(t *testing.T) {
	for _, flag := range []string{"-r", "--recursive"} {
		options, err := Parse([]string{"images", flag}, testVersion)
		if err != nil {
			t.Fatalf("Parse(%q) returned error: %v", flag, err)
		}
		if !options.Recursive {
			t.Fatalf("Parse(%q).Recursive = false, want true", flag)
		}
	}
}

func TestParseDefaultsRecursiveToFalse(t *testing.T) {
	options, err := Parse([]string{"image.jpg"}, testVersion)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if options.Recursive {
		t.Fatal("Recursive = true, want false")
	}
}

func TestParseRejectsMissingInput(t *testing.T) {
	_, err := Parse([]string{"-q", "85"}, testVersion)
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
}

func TestParseRejectsInvalidQuality(t *testing.T) {
	_, err := Parse([]string{"image.jpg", "-q", "101"}, testVersion)
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
}

func TestParseRejectsTooManyArguments(t *testing.T) {
	_, err := Parse([]string{"image.jpg", "one.webp", "two.webp"}, testVersion)
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
}

func TestRunCallsConversionWorkflow(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "image.png")
	output := filepath.Join(dir, "image.webp")
	testutil.WritePNG(t, input)

	if err := Run([]string{input, "-q", "80", output}, testVersion); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, output)
}

func TestRunUsesDefaultOutputForFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "image.png")
	testutil.WritePNG(t, input)

	if err := Run([]string{input}, testVersion); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, filepath.Join(dir, "image.webp"))
}

func TestRunRejectsRecursiveOnFileInput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "image.png")
	testutil.WritePNG(t, input)

	if err := Run([]string{input, "-r"}, testVersion); err == nil {
		t.Fatal("Run returned nil error")
	}
}

func TestRunConvertsDirectory(t *testing.T) {
	dir := t.TempDir()
	testutil.WritePNG(t, filepath.Join(dir, "a.png"))
	testutil.WritePNG(t, filepath.Join(dir, "b.png"))

	if err := Run([]string{dir}, testVersion); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, filepath.Join(dir, "a.webp"))
	testutil.AssertWebPFile(t, filepath.Join(dir, "b.webp"))
}

func assertOptions(t *testing.T, got Options, want Options) {
	t.Helper()

	if got != want {
		t.Fatalf("Options = %+v, want %+v", got, want)
	}
}
