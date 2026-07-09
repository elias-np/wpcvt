package convert

import (
	"path/filepath"
	"testing"

	"webpcvt/internal/testutil"
)

func TestRunUsesWebPEncoder(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "image.png")
	output := filepath.Join(dir, "image.webp")
	testutil.WritePNG(t, input)

	err := Run(Options{
		Input:   input,
		Output:  output,
		Quality: 80,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	testutil.AssertWebPFile(t, output)
}
