package coordmap

import "errors"

type CoordMap interface {
	Map(y float64) (float64, error)
	String() string
}

// Errors
var (
	ErrOutOfRange = errors.New("coordinate is out of range")
)
