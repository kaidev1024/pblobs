package pblobs

import (
	"github.com/h2non/bimg"
)

func resize(image []byte, width int) ([]byte, error) {
	return bimg.NewImage(image).Resize(width, 0)
}

func convertToWebp(image []byte) ([]byte, error) {
	return bimg.NewImage(image).Convert(bimg.WEBP)
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
	if w != h {
		side := min(w, h)
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
	image, err = bimg.NewImage(image).ForceResize(360, 360)
	if err != nil {
		return nil, err
	}
	return convertToWebp(image)
}

func cropAvatar(imageData []byte, x, y, size int) ([]byte, error) {
	var err error
	imageData, err = bimg.NewImage(imageData).AutoRotate()
	if err != nil {
		return nil, err
	}
	croppedImage, err := bimg.NewImage(imageData).Process(bimg.Options{
		AreaWidth:  size,
		AreaHeight: size,
		Left:       x,
		Top:        y,
	})
	if err != nil {
		return nil, err
	}
	return convertToWebp(croppedImage)
}
