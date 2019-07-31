package pdf417

import "math"

const (
	minCols         = 2
	maxCols         = 30
	maxRows         = 30
	minRows         = 2
	moduleHeight    = 2
	preferred_ratio = 3.0
)

func calculateNumberOfRows(m, k, c int) int {
	r := ((m + 1 + k) / c) + 1
	if c*r >= (m + 1 + k + c) {
		r--
	}
	return r
}

func calcDimensions(dataWords, eccWords int) (cols, rows int) {
	ratio := 0.0
	cols = 0
	rows = 0

	for c := minCols; c <= maxCols; c++ {
		r := calculateNumberOfRows(dataWords, eccWords, c)

		if r < minRows {
			break
		}

		if r > maxRows {
			continue
		}

		newRatio := float64(17*cols+69) / float64(rows*moduleHeight)
		if rows != 0 && math.Abs(newRatio-preferred_ratio) > math.Abs(ratio-preferred_ratio) {
			continue
		}

		ratio = newRatio
		cols = c
		rows = r
	}

	if rows == 0 {
		r := calculateNumberOfRows(dataWords, eccWords, minCols)
		if r < minRows {
			rows = minRows
			cols = minCols
		}
	}

	return
}
