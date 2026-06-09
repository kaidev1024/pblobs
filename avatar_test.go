package pblobs

import (
	"os"
	"testing"
)

func checkAvatarResults(t *testing.T, results [][]byte, smallFile, largeFile string) {
	t.Helper()

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	cases := []struct {
		data          []byte
		file          string
		expectedWidth int
	}{
		{results[0], smallFile, 128},
		{results[1], largeFile, 360},
	}

	for _, c := range cases {
		m, err := getImageInfo(c.data)
		if err != nil {
			t.Fatalf("%s: metadata failed: %v", c.file, err)
		}
		t.Logf("%s: %dx%d orientation=%d type=%s", c.file, m.Width, m.Height, m.Orientation, m.Type)

		if m.Width != c.expectedWidth || m.Height != c.expectedWidth {
			t.Errorf("%s: expected %dx%d, got %dx%d", c.file, c.expectedWidth, c.expectedWidth, m.Width, m.Height)
		}
		if m.Orientation > 1 {
			t.Errorf("%s: unexpected orientation %d (image may be rotated)", c.file, m.Orientation)
		}

		if err := os.WriteFile(c.file, c.data, 0644); err != nil {
			t.Fatalf("failed to write %s: %v", c.file, err)
		}
	}
}

func TestAvatarLandscape(t *testing.T) {
	data, err := os.ReadFile("landscape.png")
	if err != nil {
		t.Fatalf("failed to read landscape.png: %v", err)
	}

	meta, _ := getImageInfo(data)
	t.Logf("landscape.png: %dx%d orientation=%d", meta.Width, meta.Height, meta.Orientation)

	size := meta.Height
	x := (meta.Width - size) / 2
	results, err := ProcessAvatar(data, x, 0, size)
	if err != nil {
		t.Fatalf("ProcessAvatar failed: %v", err)
	}

	checkAvatarResults(t, results, "output_landscape_avatar_small.webp", "output_landscape_avatar_large.webp")
}

func TestAvatarPortrait(t *testing.T) {
	data, err := os.ReadFile("portrait.JPG")
	if err != nil {
		t.Fatalf("failed to read portrait.JPG: %v", err)
	}

	// Auto-rotate to get display dimensions for crop coordinate calculation.
	rotated, err := autoRotate(data)
	if err != nil {
		t.Fatalf("autoRotate failed: %v", err)
	}
	meta, _ := getImageInfo(rotated)
	t.Logf("portrait.JPG display: %dx%d orientation=%d", meta.Width, meta.Height, meta.Orientation)

	size := meta.Width // portrait: width is the short side
	x := 0
	y := (meta.Height - size) / 2
	t.Logf("crop: x=%d y=%d size=%d", x, y, size)

	results, err := ProcessAvatar(data, x, y, size)
	if err != nil {
		t.Fatalf("ProcessAvatar failed: %v", err)
	}

	checkAvatarResults(t, results, "output_portrait_avatar_small.webp", "output_portrait_avatar_large.webp")
}
