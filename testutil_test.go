package pblobs

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

type imageInfo struct {
	Width       int
	Height      int
	Orientation int
	Type        string
}

func getImageInfo(data []byte) (imageInfo, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return imageInfo{}, err
	}
	return imageInfo{
		Width:       cfg.Width,
		Height:      cfg.Height,
		Orientation: readJPEGOrientation(data),
		Type:        format,
	}, nil
}
