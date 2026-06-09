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
