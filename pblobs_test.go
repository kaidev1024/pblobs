package pblobs

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	"github.com/h2non/bimg"
)

func createTestJPEG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func TestResize(t *testing.T) {
	img := createTestJPEG(800, 600)

	t.Run("resizes to target width preserving aspect ratio", func(t *testing.T) {
		result, err := resize(img, 400)
		if err != nil {
			t.Fatalf("resize failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 400 {
			t.Errorf("expected width 400, got %d", meta.Size.Width)
		}
	})

	t.Run("returns error for invalid image data", func(t *testing.T) {
		_, err := resize([]byte("not an image"), 400)
		if err == nil {
			t.Error("expected error for invalid image data")
		}
	})
}

func TestResizeSquare(t *testing.T) {
	t.Run("square image is resized to 360x360", func(t *testing.T) {
		img := createTestJPEG(600, 600)
		result, err := resizeSquare(img)
		if err != nil {
			t.Fatalf("resizeSquare failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 360 || meta.Size.Height != 360 {
			t.Errorf("expected 360x360, got %dx%d", meta.Size.Width, meta.Size.Height)
		}
	})

	t.Run("landscape image is cropped and resized to 360x360", func(t *testing.T) {
		img := createTestJPEG(800, 400)
		result, err := resizeSquare(img)
		if err != nil {
			t.Fatalf("resizeSquare failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 360 || meta.Size.Height != 360 {
			t.Errorf("expected 360x360, got %dx%d", meta.Size.Width, meta.Size.Height)
		}
	})

	t.Run("portrait image is cropped and resized to 360x360", func(t *testing.T) {
		img := createTestJPEG(400, 800)
		result, err := resizeSquare(img)
		if err != nil {
			t.Fatalf("resizeSquare failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 360 || meta.Size.Height != 360 {
			t.Errorf("expected 360x360, got %dx%d", meta.Size.Width, meta.Size.Height)
		}
	})

	t.Run("small image is not upscaled", func(t *testing.T) {
		img := createTestJPEG(200, 300)
		result, err := resizeSquare(img)
		if err != nil {
			t.Fatalf("resizeSquare failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 200 || meta.Size.Height != 200 {
			t.Errorf("expected 200x200, got %dx%d", meta.Size.Width, meta.Size.Height)
		}
	})

	t.Run("returns error for invalid image data", func(t *testing.T) {
		_, err := resizeSquare([]byte("not an image"))
		if err == nil {
			t.Error("expected error for invalid image data")
		}
	})
}

func TestCropAvatar(t *testing.T) {
	t.Run("crops region", func(t *testing.T) {
		img := createTestJPEG(800, 800)
		result, err := cropAvatar(img, 100, 100, 200)
		if err != nil {
			t.Fatalf("cropAvatar failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 200 || meta.Size.Height != 200 {
			t.Errorf("expected 200x200, got %dx%d", meta.Size.Width, meta.Size.Height)
		}
	})

	t.Run("crops from top-left origin", func(t *testing.T) {
		img := createTestJPEG(800, 800)
		result, err := cropAvatar(img, 0, 0, 300)
		if err != nil {
			t.Fatalf("cropAvatar failed: %v", err)
		}
		meta, err := bimg.NewImage(result).Metadata()
		if err != nil {
			t.Fatalf("failed to get metadata: %v", err)
		}
		if meta.Size.Width != 300 || meta.Size.Height != 300 {
			t.Errorf("expected 300x300, got %dx%d", meta.Size.Width, meta.Size.Height)
		}
	})

	t.Run("returns error for invalid image data", func(t *testing.T) {
		_, err := cropAvatar([]byte("not an image"), 0, 0, 100)
		if err == nil {
			t.Error("expected error for invalid image data")
		}
	})
}

func TestProcessImage(t *testing.T) {
	t.Run("returns three variants: small square, medium, large", func(t *testing.T) {
		img := createTestJPEG(1920, 1080)
		results, err := ProcessImage(img)
		if err != nil {
			t.Fatalf("ProcessImage failed: %v", err)
		}
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}

		smallMeta, err := bimg.NewImage(results[0]).Metadata()
		if err != nil {
			t.Fatalf("failed to get small metadata: %v", err)
		}
		if smallMeta.Size.Width != 360 || smallMeta.Size.Height != 360 {
			t.Errorf("small: expected 360x360, got %dx%d", smallMeta.Size.Width, smallMeta.Size.Height)
		}

		medMeta, err := bimg.NewImage(results[1]).Metadata()
		if err != nil {
			t.Fatalf("failed to get medium metadata: %v", err)
		}
		if medMeta.Size.Width != 820 {
			t.Errorf("medium: expected width 820, got %d", medMeta.Size.Width)
		}

		largeMeta, err := bimg.NewImage(results[2]).Metadata()
		if err != nil {
			t.Fatalf("failed to get large metadata: %v", err)
		}
		if largeMeta.Size.Width != 1280 {
			t.Errorf("large: expected width 1280, got %d", largeMeta.Size.Width)
		}
	})

	t.Run("returns error for invalid image data", func(t *testing.T) {
		_, err := ProcessImage([]byte("not an image"))
		if err == nil {
			t.Error("expected error for invalid image data")
		}
	})
}

func TestProcessAvatar(t *testing.T) {
	t.Run("returns two sizes: 128px and 360px wide", func(t *testing.T) {
		img := createTestJPEG(800, 800)
		results, err := ProcessAvatar(img, 100, 100, 400)
		if err != nil {
			t.Fatalf("ProcessAvatar failed: %v", err)
		}
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		smallMeta, err := bimg.NewImage(results[0]).Metadata()
		if err != nil {
			t.Fatalf("failed to get small metadata: %v", err)
		}
		if smallMeta.Size.Width != 128 {
			t.Errorf("small: expected width 128, got %d", smallMeta.Size.Width)
		}

		largeMeta, err := bimg.NewImage(results[1]).Metadata()
		if err != nil {
			t.Fatalf("failed to get large metadata: %v", err)
		}
		if largeMeta.Size.Width != 360 {
			t.Errorf("large: expected width 360, got %d", largeMeta.Size.Width)
		}
	})

	t.Run("returns error for invalid image data", func(t *testing.T) {
		_, err := ProcessAvatar([]byte("not an image"), 0, 0, 100)
		if err == nil {
			t.Error("expected error for invalid image data")
		}
	})
}
