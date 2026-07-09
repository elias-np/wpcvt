package webp

import (
	"os"
	"path/filepath"
	"testing"

	"webpcvt/internal/testutil"
)

func TestEncodeFileWritesWebP(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "image.png")
	output := filepath.Join(dir, "image.webp")
	testutil.WritePNG(t, input)

	if err := EncodeFile(input, output, 80); err != nil {
		t.Fatalf("EncodeFile returned error: %v", err)
	}

	testutil.AssertWebPFile(t, output)
}

func TestEncodeFileRejectsUnknownInput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "image.txt")
	output := filepath.Join(dir, "image.webp")
	if err := os.WriteFile(input, []byte("not an image"), 0644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if err := EncodeFile(input, output, 80); err == nil {
		t.Fatal("EncodeFile returned nil error")
	}
}
