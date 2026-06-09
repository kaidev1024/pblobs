package pblobs

func ProcessAvatar(image []byte, x, y, size int) ([][]byte, error) {
	cropped, err := cropAvatar(image, x, y, size)
	if err != nil {
		return nil, err
	}
	small, err := resize(cropped, 128)
	if err != nil {
		return nil, err
	}
	medium, err := resize(cropped, 360)
	if err != nil {
		return nil, err
	}
	return [][]byte{small, medium}, nil
}

func cropAvatar(imageData []byte, x, y, size int) ([]byte, error) {
	if rotated, err := autoRotate(imageData); err == nil {
		imageData = rotated
	}
	img, _, err := decodeImage(imageData)
	if err != nil {
		return nil, err
	}
	cropped := cropRect(img, x, y, size, size)
	return encodeJPEG(cropped)
}
