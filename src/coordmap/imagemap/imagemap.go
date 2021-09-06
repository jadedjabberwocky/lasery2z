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
	ErrImageTooSmall  = errors.New("image is too small (min size 2x2)")
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

	if dx < 2 || dy < 2 {
		return nil, ErrImageTooSmall
	}

	w := options.imageWidth
	if w == 0 {
		return nil, ErrImageWidthZero
	}
	h := options.imageHeight
	if h == 0 {
		h = float64(dy-1) * w / float64(dx-1)
	}

	originZ := float64(0)
	z := make([]float64, dx)
	for i := 0; i < dx; i++ {
		x := i + x1
		y := depthAt(img, y1, y2, x)

		mapz := y * h / float64(dy-1)
		if i == 0 {
			originZ = mapz
		}
		mapz = originZ - mapz
		z[i] = mapz
	}

	return deltaymap.New(0, w/float64(dx-1), z), nil
}

func depthAt(img image.Image, y1, y2 int, x int) float64 {
	for y := y1; y < y2; y++ {
		c := img.At(x, y)
		r, g, b, a := c.RGBA()
		fr := float64(r) / float64(65535)
		fg := float64(g) / float64(65535)
		fb := float64(b) / float64(65535)
		fa := float64(a) / float64(65535)
		v := (fr + fg + fb) * fa / 3
		if v > 0 {
			return float64(y) + (1 - v)
		}
	}
	return float64(y2) - 1
}
