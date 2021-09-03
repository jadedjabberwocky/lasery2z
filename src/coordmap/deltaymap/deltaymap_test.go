package deltaymap

import "testing"

func isNearly(v1, v2 float64) bool {
	delta := 0.0001
	dv := v1 - v2
	return dv > -delta && dv < delta
}

func TestNew(t *testing.T) {
	m := New(0.5, 0.1, []float64{1, 2, 3})
	if !isNearly(m.ymax, 0.5+0.1*2) {
		t.Fatalf("ymax is wrong")
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		description string
		y0          float64
		dy          float64
		z           []float64
		tests       map[float64]float64
	}{
		{
			"0",
			0,
			1,
			[]float64{1, 2, 3, 4},
			map[float64]float64{0: 1, 0.5: 1.5, 1: 2, -1: 1, 5: 4},
		},
		{
			"1",
			10,
			-10,
			[]float64{5, 0},
			map[float64]float64{10: 5, 0: 0, 1: 0.5, 20: 5, -20: 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			m := New(testCase.y0, testCase.dy, testCase.z)

			for y, wantZ := range testCase.tests {
				z, _ := m.Map(y)

				if !isNearly(z, wantZ) {
					t.Fatalf("expected %v, but received %v", wantZ, z)
					return
				}
			}
		})
	}
}
