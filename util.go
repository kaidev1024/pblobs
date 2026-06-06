package pblobs

import (
	"github.com/h2non/bimg"
)

const THUMBNAIL_SIZE = 360

func resize(image []byte, size int) ([]byte, error) {
	meta, err := bimg.NewImage(image).Metadata()
	if err != nil {
		return nil, err
	}
	w, h := meta.Size.Width, meta.Size.Height
	if w <= size && h <= size {
		return image, nil
	}
	if w >= h {
		return bimg.NewImage(image).Resize(size, 0)
	}
	return bimg.NewImage(image).Resize(0, size)
}

func resizeSquare(image []byte) ([]byte, error) {
	// Auto-rotate based on EXIF before reading dimensions; otherwise the crop
	// operates on raw stored pixels and ForceResize later re-rotates the content.
	image, err := bimg.NewImage(image).AutoRotate()
	if err != nil {
		return nil, err
	}
	meta, err := bimg.NewImage(image).Metadata()
	if err != nil {
		return nil, err
	}
	w, h := meta.Size.Width, meta.Size.Height
	side := min(w, h)
	if w != h {
		image, err = bimg.NewImage(image).Process(bimg.Options{
			Width:   side,
			Height:  side,
			Crop:    true,
			Gravity: bimg.GravityCentre,
		})
		if err != nil {
			return nil, err
		}
	}
	if side <= THUMBNAIL_SIZE {
		return image, nil
	}
	return bimg.NewImage(image).Resize(THUMBNAIL_SIZE, THUMBNAIL_SIZE)
}
