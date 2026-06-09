package pblobs

import (
	"os"
	"testing"
)

func logResults(t *testing.T, results [][]byte, names []string) {
	t.Helper()
	for i, name := range names {
		if err := os.WriteFile(name, results[i], 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
		m, _ := getImageInfo(results[i])
		t.Logf("wrote %s: %dx%d type=%s", name, m.Width, m.Height, m.Type)
	}
}

func TestGenerateOutputs(t *testing.T) {
	t.Run("landscape/ProcessImage", func(t *testing.T) {
		data, err := os.ReadFile("landscape.png")
		if err != nil {
			t.Fatalf("failed to read landscape.png: %v", err)
		}
		results, err := ProcessImage(data)
		if err != nil {
			t.Fatalf("ProcessImage failed: %v", err)
		}
		logResults(t, results, []string{
			"output_landscape_small.webp",
			"output_landscape_medium.webp",
			"output_landscape_large.webp",
		})
	})

	t.Run("portrait/ProcessImage", func(t *testing.T) {
		data, err := os.ReadFile("portrait.JPG")
		if err != nil {
			t.Fatalf("failed to read portrait.JPG: %v", err)
		}
		meta, _ := getImageInfo(data)
		t.Logf("portrait.JPG raw: %dx%d orientation=%d", meta.Width, meta.Height, meta.Orientation)
		results, err := ProcessImage(data)
		if err != nil {
			t.Fatalf("ProcessImage failed: %v", err)
		}
		logResults(t, results, []string{
			"output_portrait_small.webp",
			"output_portrait_medium.webp",
			"output_portrait_large.webp",
		})
	})
}
