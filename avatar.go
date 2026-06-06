package pblobs

import "github.com/h2non/bimg"

func ProcessAvatar(image []byte, x, y, size int) ([][]byte, error) {
	croppedImage, err := cropAvatar(image, x, y, size)
	if err != nil {
		return nil, err
	}
	small, err := resize(croppedImage, 128)
	if err != nil {
		return nil, err
	}
	medium, err := resize(croppedImage, 360)
	if err != nil {
		return nil, err
	}
	return [][]byte{small, medium}, nil
}

func cropAvatar(imageData []byte, x, y, size int) ([]byte, error) {
	var err error
	imageData, err = bimg.NewImage(imageData).AutoRotate()
	if err != nil {
		return nil, err
	}
	return bimg.NewImage(imageData).Process(bimg.Options{
		AreaWidth:  size,
		AreaHeight: size,
		Left:       x,
		Top:        y,
	})
}
