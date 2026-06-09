package pblobs

import (
	"bytes"
	"encoding/binary"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// decodeImage decodes image bytes and returns the image, format string, and any error.
func decodeImage(imageData []byte) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(imageData))
}

// autoRotate reads the EXIF orientation tag from JPEG data and returns the image
// rotated to its natural upright orientation. Non-JPEG or untagged images are
// returned unchanged. Only JPEG supports EXIF orientation, so we skip other formats.
func autoRotate(imageData []byte) ([]byte, error) {
	orientation := readJPEGOrientation(imageData)
	if orientation <= 1 {
		return imageData, nil
	}
	img, _, err := decodeImage(imageData)
	if err != nil {
		return nil, err
	}
	rotated := applyOrientation(img, orientation)
	return encodeJPEG(rotated)
}

// readJPEGOrientation parses JPEG APP1/EXIF to extract the Orientation tag (0x0112).
// Returns 1 (normal) if not found or not a JPEG.
func readJPEGOrientation(data []byte) int {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return 1
	}
	i := 2
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			break
		}
		marker := data[i+1]
		segLen := int(binary.BigEndian.Uint16(data[i+2 : i+4]))
		if marker == 0xE1 && i+4+segLen <= len(data) {
			seg := data[i+4 : i+4+segLen]
			if len(seg) >= 6 && string(seg[0:6]) == "Exif\x00\x00" {
				return parseExifOrientation(seg[6:])
			}
		}
		i += 2 + segLen
	}
	return 1
}

// parseExifOrientation reads the Orientation tag from raw TIFF data.
func parseExifOrientation(tiff []byte) int {
	if len(tiff) < 8 {
		return 1
	}
	var bo binary.ByteOrder
	switch string(tiff[0:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return 1
	}
	ifdOffset := int(bo.Uint32(tiff[4:8]))
	if ifdOffset+2 > len(tiff) {
		return 1
	}
	count := int(bo.Uint16(tiff[ifdOffset : ifdOffset+2]))
	for j := 0; j < count; j++ {
		entry := ifdOffset + 2 + j*12
		if entry+12 > len(tiff) {
			break
		}
		tag := bo.Uint16(tiff[entry : entry+2])
		if tag == 0x0112 {
			return int(bo.Uint16(tiff[entry+8 : entry+10]))
		}
	}
	return 1
}

// applyOrientation transforms img according to the EXIF orientation value.
func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return flipH(img)
	case 3:
		return rotate180(img)
	case 4:
		return flipV(img)
	case 5:
		return flipH(rotate90CW(img))
	case 6:
		return rotate90CW(img)
	case 7:
		return flipH(rotate90CCW(img))
	case 8:
		return rotate90CCW(img)
	default:
		return img
	}
}

func rotate90CW(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(h-1-y, x, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func rotate90CCW(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(y, w-1-x, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func rotate180(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(w-1-x, h-1-y, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func flipH(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(w-1-x, y, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func flipV(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(x, h-1-y, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}
