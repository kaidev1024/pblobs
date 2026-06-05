package pblobs

func ProcessImage(image []byte) ([][]byte, error) {
	small, err := resizeSquare(image)
	if err != nil {
		return nil, err
	}
	medium, err := resize(image, 820)
	if err != nil {
		return nil, err
	}
	large, err := resize(image, 1280)
	if err != nil {
		return nil, err
	}
	return [][]byte{small, medium, large}, nil
}

func ProcessAvatar(image []byte, x, y, size int) ([][]byte, error) {
	croppedImage, err := cropAvatar(image, x, y, size)
	if err != nil {
		return nil, err
	}
	small, err := resize(croppedImage, 128)
	if err != nil {
		return nil, err
	}
	large, err := resize(croppedImage, 360)
	if err != nil {
		return nil, err
	}
	return [][]byte{small, large}, nil
}
