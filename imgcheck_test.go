package pblobs

import (
	"os"
	"testing"
)

func TestRealPortraitJPEG(t *testing.T) {
	data, err := os.ReadFile("portrait.JPG")
	if err != nil {
		t.Fatalf("failed to read portrait.JPG: %v", err)
	}

	meta, err := getImageInfo(data)
	if err != nil {
		t.Fatalf("failed to get metadata: %v", err)
	}
	t.Logf("Input: %dx%d orientation=%d type=%s", meta.Width, meta.Height, meta.Orientation, meta.Type)

	results, err := ProcessImage(data)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	for i, r := range results {
		m, err := getImageInfo(r)
		if err != nil {
			t.Fatalf("result[%d] metadata failed: %v", i, err)
		}
		t.Logf("result[%d]: %dx%d type=%s", i, m.Width, m.Height, m.Type)
	}

	medMeta, _ := getImageInfo(results[1])
	if medMeta.Width != 820 {
		t.Errorf("medium: expected width 820, got %d", medMeta.Width)
	}
	if medMeta.Height < medMeta.Width {
		t.Errorf("portrait orientation lost: medium is %dx%d (width > height)", medMeta.Width, medMeta.Height)
	}
}

func TestCheckDimensions(t *testing.T) {
	data, _ := os.ReadFile("landscape.png")
	meta, _ := getImageInfo(data)
	t.Logf("landscape.png: %dx%d orientation=%d type=%s", meta.Width, meta.Height, meta.Orientation, meta.Type)

	portrait := createTestJPEG(400, 800)
	result, err := resize(portrait, 820)
	if err != nil {
		t.Fatalf("%v", err)
	}
	m, _ := getImageInfo(result)
	t.Logf("Portrait 400x800 -> resize(820) -> %dx%d", m.Width, m.Height)

	portrait2 := createTestJPEG(1200, 2400)
	result2, err := resize(portrait2, 820)
	if err != nil {
		t.Fatalf("%v", err)
	}
	m2, _ := getImageInfo(result2)
	t.Logf("Portrait 1200x2400 -> resize(820) -> %dx%d", m2.Width, m2.Height)

	results, err := ProcessImage(portrait2)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	for i, r := range results {
		m, _ := getImageInfo(r)
		t.Logf("ProcessImage result[%d]: %dx%d", i, m.Width, m.Height)
	}
}
