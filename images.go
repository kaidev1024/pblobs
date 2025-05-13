package pblobs

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"

	"github.com/h2non/bimg"
)

func GetImageDataAndFormat(file io.Reader, contentType string) ([]byte, string) {
	imageFormat := contentType[strings.IndexByte(contentType, '/')+1:]
	buf := &bytes.Buffer{}
	buf.ReadFrom(file)
	data := buf.Bytes()
	return data, imageFormat
}

func GetImageData(file io.ReadSeeker) []byte {
	buf := &bytes.Buffer{}
	buf.ReadFrom(file)
	return buf.Bytes()
}

func GetProfileImageTypeFromDataUrl(url string) string {
	profileImageTypeHeader := url[0:strings.IndexByte(url, ';')]
	return profileImageTypeHeader[strings.IndexByte(profileImageTypeHeader, '/')+1:]
}

func GetImageDataFromDataUrl(url string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(url[strings.IndexByte(url, ',')+1:])
}

// dataUrl is png/jpeg encoded in base64
func Resize(imageData []byte, width, height int) ([]byte, error) {
	var err error
	if width != 0 {
		imageData, err = bimg.NewImage(imageData).ForceResize(width, height)
		if err != nil {
			return nil, err
		}
	}

	return imageData, nil
}

func ConvertToPng(image []byte) ([]byte, error) {
	return bimg.NewImage(image).Convert(bimg.PNG)
}
