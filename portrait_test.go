package pblobs

import (
	"testing"
	"github.com/h2non/bimg"
)

func TestPortraitResize(t *testing.T) {
	// Portrait image: width < height
	img := createTestJPEG(600, 1200)

	result, err := resize(img, 820)
	if err != nil {
		t.Fatalf("resize failed: %v", err)
	}

	meta, err := bimg.NewImage(result).Metadata()
	if err != nil {
		t.Fatalf("metadata failed: %v", err)
	}
	t.Logf("Input: 600x1200, Output: %dx%d", meta.Size.Width, meta.Size.Height)

	if meta.Size.Width > meta.Size.Height {
		t.Errorf("portrait image got rotated: expected height > width, got %dx%d", meta.Size.Width, meta.Size.Height)
	}

	// Larger portrait image that actually needs downscaling
	img2 := createTestJPEG(1200, 2400)
	result2, err := resize(img2, 820)
	if err != nil {
		t.Fatalf("resize failed: %v", err)
	}
	meta2, err := bimg.NewImage(result2).Metadata()
	if err != nil {
		t.Fatalf("metadata failed: %v", err)
	}
	t.Logf("Input: 1200x2400, Output: %dx%d", meta2.Size.Width, meta2.Size.Height)

	if meta2.Size.Width > meta2.Size.Height {
		t.Errorf("portrait image got rotated: expected height > width, got %dx%d", meta2.Size.Width, meta2.Size.Height)
	}
	if meta2.Size.Width != 820 {
		t.Errorf("expected width 820, got %d", meta2.Size.Width)
	}
}
