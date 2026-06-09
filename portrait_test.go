package pblobs

import (
	"testing"
)

func TestPortraitResize(t *testing.T) {
	// Portrait image: width < height
	img := createTestJPEG(600, 1200)

	result, err := resize(img, 820)
	if err != nil {
		t.Fatalf("resize failed: %v", err)
	}

	meta, err := getImageInfo(result)
	if err != nil {
		t.Fatalf("metadata failed: %v", err)
	}
	t.Logf("Input: 600x1200, Output: %dx%d", meta.Width, meta.Height)

	if meta.Width > meta.Height {
		t.Errorf("portrait image got rotated: expected height > width, got %dx%d", meta.Width, meta.Height)
	}

	// Larger portrait image that actually needs downscaling
	img2 := createTestJPEG(1200, 2400)
	result2, err := resize(img2, 820)
	if err != nil {
		t.Fatalf("resize failed: %v", err)
	}
	meta2, err := getImageInfo(result2)
	if err != nil {
		t.Fatalf("metadata failed: %v", err)
	}
	t.Logf("Input: 1200x2400, Output: %dx%d", meta2.Width, meta2.Height)

	if meta2.Width > meta2.Height {
		t.Errorf("portrait image got rotated: expected height > width, got %dx%d", meta2.Width, meta2.Height)
	}
	if meta2.Width != 820 {
		t.Errorf("expected width 820, got %d", meta2.Width)
	}
}
