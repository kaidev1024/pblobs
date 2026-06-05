package pblobs

import (
	"os"
	"testing"

	"github.com/h2non/bimg"
)

func logResults(t *testing.T, results [][]byte, names []string) {
	t.Helper()
	for i, name := range names {
		if err := os.WriteFile(name, results[i], 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
		m, _ := bimg.NewImage(results[i]).Metadata()
		t.Logf("wrote %s: %dx%d type=%s", name, m.Size.Width, m.Size.Height, m.Type)
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
			"output_landscape_medium.png",
			"output_landscape_large.png",
		})
	})


	t.Run("portrait/ProcessImage", func(t *testing.T) {
		data, err := os.ReadFile("portrait.JPG")
		if err != nil {
			t.Fatalf("failed to read portrait.JPG: %v", err)
		}
		meta, _ := bimg.NewImage(data).Metadata()
		t.Logf("portrait.JPG raw: %dx%d orientation=%d", meta.Size.Width, meta.Size.Height, meta.Orientation)
		results, err := ProcessImage(data)
		if err != nil {
			t.Fatalf("ProcessImage failed: %v", err)
		}
		logResults(t, results, []string{
			"output_portrait_small.webp",
			"output_portrait_medium.jpeg",
			"output_portrait_large.jpeg",
		})
	})

}
