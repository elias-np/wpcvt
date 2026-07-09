package webp

/*
#cgo CFLAGS: -I${SRCDIR}/../../vendor/libwebp/src
#cgo windows LDFLAGS: -L${SRCDIR}/../../vendor/libwebp/src -L${SRCDIR}/../../vendor/libwebp/sharpyuv -lwebp -lsharpyuv
#include <stdlib.h>
#include <webp/encode.h>
#include <webp/types.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"
)

// EncodeFile encodes input as a WebP file at output with the given quality.
func EncodeFile(input string, output string, quality int) error {
	img, err := readImage(input)
	if err != nil {
		return err
	}

	data, err := encodeRGBA(img, quality)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, data, 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func readImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode input: %w", err)
	}

	return img, nil
}

func encodeRGBA(img image.Image, quality int) ([]byte, error) {
	rgba := toRGBA(img)
	bounds := rgba.Bounds()
	if bounds.Empty() {
		return nil, errors.New("image is empty")
	}

	var output *C.uint8_t
	size := C.WebPEncodeRGBA(
		(*C.uint8_t)(unsafe.Pointer(&rgba.Pix[0])),
		C.int(bounds.Dx()),
		C.int(bounds.Dy()),
		C.int(rgba.Stride),
		C.float(quality),
		&output,
	)
	if size == 0 || output == nil {
		return nil, errors.New("encode webp")
	}
	defer C.WebPFree(unsafe.Pointer(output))

	return C.GoBytes(unsafe.Pointer(output), C.int(size)), nil
}

func toRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)
	return rgba
}
