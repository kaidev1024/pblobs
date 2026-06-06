package pblobs

import (
	"os"
	"testing"

	"github.com/h2non/bimg"
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
		m, err := bimg.NewImage(c.data).Metadata()
		if err != nil {
			t.Fatalf("%s: metadata failed: %v", c.file, err)
		}
		t.Logf("%s: %dx%d orientation=%d type=%s", c.file, m.Size.Width, m.Size.Height, m.Orientation, m.Type)

		if m.Size.Width != c.expectedWidth || m.Size.Height != c.expectedWidth {
			t.Errorf("%s: expected %dx%d, got %dx%d", c.file, c.expectedWidth, c.expectedWidth, m.Size.Width, m.Size.Height)
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

	meta, _ := bimg.NewImage(data).Metadata()
	t.Logf("landscape.png: %dx%d orientation=%d", meta.Size.Width, meta.Size.Height, meta.Orientation)

	size := meta.Size.Height
	x := (meta.Size.Width - size) / 2
	results, err := ProcessAvatar(data, x, 0, size)
	if err != nil {
		t.Fatalf("ProcessAvatar failed: %v", err)
	}

	checkAvatarResults(t, results, "output_landscape_avatar_small.png", "output_landscape_avatar_large.png")
}

func TestAvatarPortrait(t *testing.T) {
	data, err := os.ReadFile("portrait.JPG")
	if err != nil {
		t.Fatalf("failed to read portrait.JPG: %v", err)
	}

	// Auto-rotate to get display dimensions for crop coordinate calculation.
	rotated, err := bimg.NewImage(data).AutoRotate()
	if err != nil {
		t.Fatalf("AutoRotate failed: %v", err)
	}
	meta, _ := bimg.NewImage(rotated).Metadata()
	t.Logf("portrait.JPG display: %dx%d orientation=%d", meta.Size.Width, meta.Size.Height, meta.Orientation)

	size := meta.Size.Width // portrait: width is the short side
	x := 0
	y := (meta.Size.Height - size) / 2
	t.Logf("crop: x=%d y=%d size=%d", x, y, size)

	results, err := ProcessAvatar(data, x, y, size)
	if err != nil {
		t.Fatalf("ProcessAvatar failed: %v", err)
	}

	checkAvatarResults(t, results, "output_portrait_avatar_small.jpeg", "output_portrait_avatar_large.jpeg")
}
