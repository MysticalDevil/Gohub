package tests

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func CreatePNGFile(t *testing.T) string {
	t.Helper()

	file, err := os.CreateTemp("", "avatar-*.png")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("encode png failed: %v", err)
	}

	return file.Name()
}
