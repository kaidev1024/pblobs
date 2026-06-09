package pblobs

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"math"
)

const THUMBNAIL_SIZE = 360

func resize(imageData []byte, size int) ([]byte, error) {
	if rotated, err := autoRotate(imageData); err == nil {
		imageData = rotated
	}
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	if w <= size {
		return imageData, nil
	}
	newW := size
	newH := int(math.Round(float64(h) * float64(size) / float64(w)))
	dst := scaleImage(img, newW, newH)
	return encodeJPEG(dst)
}

func resizeSquare(imageData []byte) ([]byte, error) {
	if rotated, err := autoRotate(imageData); err == nil {
		imageData = rotated
	}
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	side := min(w, h)
	if w != h {
		img = cropRect(img, (w-side)/2, (h-side)/2, side, side)
	}
	if side <= THUMBNAIL_SIZE {
		return encodeJPEG(img)
	}
	return encodeJPEG(scaleImage(img, THUMBNAIL_SIZE, THUMBNAIL_SIZE))
}

// scaleImage resizes src to newW×newH using bilinear interpolation.
func scaleImage(src image.Image, newW, newH int) image.Image {
	srcB := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	scaleX := float64(srcB.Dx()) / float64(newW)
	scaleY := float64(srcB.Dy()) / float64(newH)
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			fx := float64(x)*scaleX + float64(srcB.Min.X)
			fy := float64(y)*scaleY + float64(srcB.Min.Y)
			x0 := int(fx)
			y0 := int(fy)
			x1 := min(x0+1, srcB.Max.X-1)
			y1 := min(y0+1, srcB.Max.Y-1)
			dx := fx - float64(x0)
			dy := fy - float64(y0)
			c00 := toRGBAF(src.At(x0, y0))
			c10 := toRGBAF(src.At(x1, y0))
			c01 := toRGBAF(src.At(x0, y1))
			c11 := toRGBAF(src.At(x1, y1))
			r := c00[0]*(1-dx)*(1-dy) + c10[0]*dx*(1-dy) + c01[0]*(1-dx)*dy + c11[0]*dx*dy
			g := c00[1]*(1-dx)*(1-dy) + c10[1]*dx*(1-dy) + c01[1]*(1-dx)*dy + c11[1]*dx*dy
			b := c00[2]*(1-dx)*(1-dy) + c10[2]*dx*(1-dy) + c01[2]*(1-dx)*dy + c11[2]*dx*dy
			a := c00[3]*(1-dx)*(1-dy) + c10[3]*dx*(1-dy) + c01[3]*(1-dx)*dy + c11[3]*dx*dy
			dst.Set(x, y, color.RGBA{
				R: uint8(math.Round(r)),
				G: uint8(math.Round(g)),
				B: uint8(math.Round(b)),
				A: uint8(math.Round(a)),
			})
		}
	}
	return dst
}

func toRGBAF(c color.Color) [4]float64 {
	r, g, b, a := c.RGBA()
	return [4]float64{float64(r >> 8), float64(g >> 8), float64(b >> 8), float64(a >> 8)}
}

func cropRect(src image.Image, x, y, w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	srcB := src.Bounds()
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			dst.Set(dx, dy, src.At(srcB.Min.X+x+dx, srcB.Min.Y+y+dy))
		}
	}
	return dst
}

func encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
