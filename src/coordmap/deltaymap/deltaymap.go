package deltaymap

import (
	"encoding/json"
	"math"
)

type DeltaYMap struct {
	y0   float64
	ymax float64
	dy   float64
	z    []float64
}

type deltaYMapJSON struct {
	Type string    `json:"type"`
	Y0   float64   `json:"y0"`
	DY   float64   `json:"dy"`
	Z    []float64 `json:"z"`
}

func New(y0 float64, dy float64, z []float64) *DeltaYMap {
	return &DeltaYMap{
		y0:   y0,
		ymax: y0 + dy*float64(len(z)-1),
		dy:   dy,
		z:    z,
	}
}

func (m *DeltaYMap) Map(y float64) (float64, error) {
	if len(m.z) == 0 {
		return 0, nil
	}

	if m.dy == 0 {
		return m.z[0], nil
	} else if m.dy > 0 {
		if y <= m.y0 {
			return m.z[0], nil
		}

		if y >= m.ymax {
			return m.z[len(m.z)-1], nil
		}
	} else /* if m.dy < 0 */ {
		if y >= m.y0 {
			return m.z[0], nil
		}

		if y <= m.ymax {
			return m.z[len(m.z)-1], nil
		}
	}

	yi := (y - m.y0) / m.dy
	yi1 := int(math.Floor(yi))
	yi2 := yi1 + 1
	frac := yi - float64(yi1)

	z1 := m.z[yi1]
	z2 := m.z[yi2]

	return z1*(1-frac) + z2*(frac), nil
}

func (m *DeltaYMap) String() string {
	jm := deltaYMapJSON{
		Type: "deltaYMap",
		Y0:   m.y0,
		DY:   m.dy,
		Z:    m.z,
	}
	b, _ := json.Marshal(jm)
	return string(b)
}
