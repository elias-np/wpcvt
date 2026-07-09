package testutil

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// WritePNG writes a small valid PNG image to path.
func WritePNG(t *testing.T, path string) {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, A: 255})

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
}

// AssertWebPFile checks that path contains a RIFF WebP file.
func AssertWebPFile(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if len(data) < 12 {
		t.Fatalf("WebP file is too small: %d bytes", len(data))
	}
	if string(data[:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		t.Fatalf("file does not have a WebP RIFF header")
	}
}
