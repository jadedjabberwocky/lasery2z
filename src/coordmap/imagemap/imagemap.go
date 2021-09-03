package imagemap

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/jadedjabberwocky/lasery2z/coordmap/deltaymap"
)

var (
	ErrImageWidthZero = errors.New("image width is 0")
)

func New(filename string, options *Options) (*deltaymap.DeltaYMap, error) {
	if options == nil {
		options = DefaultOptions()
	}

	imgFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	x1 := img.Bounds().Min.X
	x2 := img.Bounds().Max.X
	dx := x2 - x1
	y1 := img.Bounds().Min.Y
	y2 := img.Bounds().Max.Y
	dy := y2 - y1

	w := options.imageWidth
	if w == 0 {
		return nil, ErrImageWidthZero
	}
	h := options.imageHeight
	if h == 0 {
		h = float64(y2-y1) * w / float64(x2-x1)
	}

	z := make([]float64, dx)
	for i := 0; i < dx; i++ {
		x := i + x1
		y := depthAt(img, y1, y2, x)

		mapz := float64(y) * h / float64(dy)
		z[i] = mapz
	}

	return deltaymap.New(0, w/float64(dx), z), nil
}

func depthAt(img image.Image, y1, y2 int, x int) int {
	for y := y1; y < y2; y++ {
		c := img.At(x, y)
		r, g, b, a := c.RGBA()
		fr := float64(r) / float64(65535)
		fg := float64(g) / float64(65535)
		fb := float64(b) / float64(65535)
		fa := float64(a) / float64(65535)
		v := (fr + fg + fb) * fa / 3
		if v > 0.5 {
			return y
		}
	}
	return y2 - 1
}
