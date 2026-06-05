package pblobs

import (
	"os"
	"testing"

	"github.com/h2non/bimg"
)

func TestRealPortraitJPEG(t *testing.T) {
	data, err := os.ReadFile("portrait.JPG")
	if err != nil {
		t.Fatalf("failed to read portrait.JPG: %v", err)
	}

	meta, err := bimg.NewImage(data).Metadata()
	if err != nil {
		t.Fatalf("failed to get metadata: %v", err)
	}
	t.Logf("Input: %dx%d orientation=%d type=%s", meta.Size.Width, meta.Size.Height, meta.Orientation, meta.Type)

	results, err := ProcessImage(data)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	for i, r := range results {
		m, err := bimg.NewImage(r).Metadata()
		if err != nil {
			t.Fatalf("result[%d] metadata failed: %v", i, err)
		}
		t.Logf("result[%d]: %dx%d type=%s", i, m.Size.Width, m.Size.Height, m.Type)
	}

	medMeta, _ := bimg.NewImage(results[1]).Metadata()
	if medMeta.Size.Width != 820 {
		t.Errorf("medium: expected width 820, got %d", medMeta.Size.Width)
	}
	if medMeta.Size.Height < medMeta.Size.Width {
		t.Errorf("portrait orientation lost: medium is %dx%d (width > height)", medMeta.Size.Width, medMeta.Size.Height)
	}
}

func TestCheckDimensions(t *testing.T) {
	data, _ := os.ReadFile("landscape.png")
	meta, _ := bimg.NewImage(data).Metadata()
	t.Logf("landscape.png: %dx%d orientation=%d type=%s", meta.Size.Width, meta.Size.Height, meta.Orientation, meta.Type)

	portrait := createTestJPEG(400, 800)
	result, err := resize(portrait, 820)
	if err != nil {
		t.Fatalf("%v", err)
	}
	m, _ := bimg.NewImage(result).Metadata()
	t.Logf("Portrait 400x800 -> resize(820) -> %dx%d", m.Size.Width, m.Size.Height)

	portrait2 := createTestJPEG(1200, 2400)
	result2, err := resize(portrait2, 820)
	if err != nil {
		t.Fatalf("%v", err)
	}
	m2, _ := bimg.NewImage(result2).Metadata()
	t.Logf("Portrait 1200x2400 -> resize(820) -> %dx%d", m2.Size.Width, m2.Size.Height)

	results, err := ProcessImage(portrait2)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	for i, r := range results {
		m, _ := bimg.NewImage(r).Metadata()
		t.Logf("ProcessImage result[%d]: %dx%d", i, m.Size.Width, m.Size.Height)
	}
}
