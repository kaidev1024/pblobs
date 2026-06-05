package pblobs

import (
	"os"
	"testing"
)

func TestGenerateOutputs(t *testing.T) {
	data, err := os.ReadFile("test.png")
	if err != nil {
		t.Fatalf("failed to read test.png: %v", err)
	}

	t.Run("ProcessImage", func(t *testing.T) {
		results, err := ProcessImage(data)
		if err != nil {
			t.Fatalf("ProcessImage failed: %v", err)
		}
		outputs := []string{"output_small.webp", "output_medium.webp", "output_large.webp"}
		for i, name := range outputs {
			if err := os.WriteFile(name, results[i], 0644); err != nil {
				t.Fatalf("failed to write %s: %v", name, err)
			}
			t.Logf("wrote %s (%d bytes)", name, len(results[i]))
		}
	})

	t.Run("ProcessAvatar", func(t *testing.T) {
		// center-crop a square from the 3840x2160 image
		x, y, size := (3840-2160)/2, 0, 2160
		results, err := ProcessAvatar(data, x, y, size)
		if err != nil {
			t.Fatalf("ProcessAvatar failed: %v", err)
		}
		outputs := []string{"output_avatar_small.webp", "output_avatar_large.webp"}
		for i, name := range outputs {
			if err := os.WriteFile(name, results[i], 0644); err != nil {
				t.Fatalf("failed to write %s: %v", name, err)
			}
			t.Logf("wrote %s (%d bytes)", name, len(results[i]))
		}
	})
}
