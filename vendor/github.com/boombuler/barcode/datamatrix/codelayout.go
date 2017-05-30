package datamatrix

import (
	"github.com/boombuler/barcode/utils"
	"strconv"
)

type setValFunc func(byte)

type codeLayout struct {
	matrix *utils.BitList
	occupy *utils.BitList
	size   *dmCodeSize
}

func newCodeLayout(size *dmCodeSize) *codeLayout {
	result := new(codeLayout)
	result.matrix = utils.NewBitList(size.MatrixColumns() * size.MatrixRows())
	result.occupy = utils.NewBitList(size.MatrixColumns() * size.MatrixRows())
	result.size = size
	return result
}

func (l *codeLayout) Occupied(row, col int) bool {
	return l.occupy.GetBit(col + row*l.size.MatrixColumns())
}

func (l *codeLayout) Set(row, col int, value, bitNum byte) {
	val := ((value >> (7 - bitNum)) & 1) == 1
	if row < 0 {
		row += l.size.MatrixRows()
		col += 4 - ((l.size.MatrixRows() + 4) % 8)
	}
	if col < 0 {
		col += l.size.MatrixColumns()
		row += 4 - ((l.size.MatrixColumns() + 4) % 8)
	}
	if l.Occupied(row, col) {
		panic("Field already occupied row: " + strconv.Itoa(row) + " col: " + strconv.Itoa(col))
	}

	l.occupy.SetBit(col+row*l.size.MatrixColumns(), true)

	l.matrix.SetBit(col+row*l.size.MatrixColumns(), val)
}

func (l *codeLayout) SetSimple(row, col int, value byte) {
	l.Set(row-2, col-2, value, 0)
	l.Set(row-2, col-1, value, 1)
	l.Set(row-1, col-2, value, 2)
	l.Set(row-1, col-1, value, 3)
	l.Set(row-1, col-0, value, 4)
	l.Set(row-0, col-2, value, 5)
	l.Set(row-0, col-1, value, 6)
	l.Set(row-0, col-0, value, 7)
}

func (l *codeLayout) Corner1(value byte) {
	l.Set(l.size.MatrixRows()-1, 0, value, 0)
	l.Set(l.size.MatrixRows()-1, 1, value, 1)
	l.Set(l.size.MatrixRows()-1, 2, value, 2)
	l.Set(0, l.size.MatrixColumns()-2, value, 3)
	l.Set(0, l.size.MatrixColumns()-1, value, 4)
	l.Set(1, l.size.MatrixColumns()-1, value, 5)
	l.Set(2, l.size.MatrixColumns()-1, value, 6)
	l.Set(3, l.size.MatrixColumns()-1, value, 7)
}

func (l *codeLayout) Corner2(value byte) {
	l.Set(l.size.MatrixRows()-3, 0, value, 0)
	l.Set(l.size.MatrixRows()-2, 0, value, 1)
	l.Set(l.size.MatrixRows()-1, 0, value, 2)
	l.Set(0, l.size.MatrixColumns()-4, value, 3)
	l.Set(0, l.size.MatrixColumns()-3, value, 4)
	l.Set(0, l.size.MatrixColumns()-2, value, 5)
	l.Set(0, l.size.MatrixColumns()-1, value, 6)
	l.Set(1, l.size.MatrixColumns()-1, value, 7)
}

func (l *codeLayout) Corner3(value byte) {
	l.Set(l.size.MatrixRows()-3, 0, value, 0)
	l.Set(l.size.MatrixRows()-2, 0, value, 1)
	l.Set(l.size.MatrixRows()-1, 0, value, 2)
	l.Set(0, l.size.MatrixColumns()-2, value, 3)
	l.Set(0, l.size.MatrixColumns()-1, value, 4)
	l.Set(1, l.size.MatrixColumns()-1, value, 5)
	l.Set(2, l.size.MatrixColumns()-1, value, 6)
	l.Set(3, l.size.MatrixColumns()-1, value, 7)
}

func (l *codeLayout) Corner4(value byte) {
	l.Set(l.size.MatrixRows()-1, 0, value, 0)
	l.Set(l.size.MatrixRows()-1, l.size.MatrixColumns()-1, value, 1)
	l.Set(0, l.size.MatrixColumns()-3, value, 2)
	l.Set(0, l.size.MatrixColumns()-2, value, 3)
	l.Set(0, l.size.MatrixColumns()-1, value, 4)
	l.Set(1, l.size.MatrixColumns()-3, value, 5)
	l.Set(1, l.size.MatrixColumns()-2, value, 6)
	l.Set(1, l.size.MatrixColumns()-1, value, 7)
}

func (l *codeLayout) SetValues(data []byte) {
	idx := 0
	row := 4
	col := 0

	for (row < l.size.MatrixRows()) || (col < l.size.MatrixColumns()) {
		if (row == l.size.MatrixRows()) && (col == 0) {
			l.Corner1(data[idx])
			idx++
		}
		if (row == l.size.MatrixRows()-2) && (col == 0) && (l.size.MatrixColumns()%4 != 0) {
			l.Corner2(data[idx])
			idx++
		}
		if (row == l.size.MatrixRows()-2) && (col == 0) && (l.size.MatrixColumns()%8 == 4) {
			l.Corner3(data[idx])
			idx++
		}

		if (row == l.size.MatrixRows()+4) && (col == 2) && (l.size.MatrixColumns()%8 == 0) {
			l.Corner4(data[idx])
			idx++
		}

		for true {
			if (row < l.size.MatrixRows()) && (col >= 0) && !l.Occupied(row, col) {
				l.SetSimple(row, col, data[idx])
				idx++
			}
			row -= 2
			col += 2
			if (row < 0) || (col >= l.size.MatrixColumns()) {
				break
			}
		}
		row += 1
		col += 3

		for true {
			if (row >= 0) && (col < l.size.MatrixColumns()) && !l.Occupied(row, col) {
				l.SetSimple(row, col, data[idx])
				idx++
			}
			row += 2
			col -= 2
			if (row >= l.size.MatrixRows()) || (col < 0) {
				break
			}
		}
		row += 3
		col += 1
	}

	if !l.Occupied(l.size.MatrixRows()-1, l.size.MatrixColumns()-1) {
		l.Set(l.size.MatrixRows()-1, l.size.MatrixColumns()-1, 255, 0)
		l.Set(l.size.MatrixRows()-2, l.size.MatrixColumns()-2, 255, 0)
	}
}

func (l *codeLayout) Merge() *datamatrixCode {
	result := newDataMatrixCode(l.size)

	//dotted horizontal lines
	for r := 0; r < l.size.Rows; r += (l.size.RegionRows() + 2) {
		for c := 0; c < l.size.Columns; c += 2 {
			result.set(c, r, true)
		}
	}

	//solid horizontal line
	for r := l.size.RegionRows() + 1; r < l.size.Rows; r += (l.size.RegionRows() + 2) {
		for c := 0; c < l.size.Columns; c++ {
			result.set(c, r, true)
		}
	}

	//dotted vertical lines
	for c := l.size.RegionColumns() + 1; c < l.size.Columns; c += (l.size.RegionColumns() + 2) {
		for r := 1; r < l.size.Rows; r += 2 {
			result.set(c, r, true)
		}
	}

	//solid vertical line
	for c := 0; c < l.size.Columns; c += (l.size.RegionColumns() + 2) {
		for r := 0; r < l.size.Rows; r++ {
			result.set(c, r, true)
		}
	}
	count := 0
	for hRegion := 0; hRegion < l.size.RegionCountHorizontal; hRegion++ {
		for vRegion := 0; vRegion < l.size.RegionCountVertical; vRegion++ {
			for x := 0; x < l.size.RegionColumns(); x++ {
				colMatrix := (l.size.RegionColumns() * hRegion) + x
				colResult := ((2 + l.size.RegionColumns()) * hRegion) + x + 1

				for y := 0; y < l.size.RegionRows(); y++ {
					rowMatrix := (l.size.RegionRows() * vRegion) + y
					rowResult := ((2 + l.size.RegionRows()) * vRegion) + y + 1
					val := l.matrix.GetBit(colMatrix + rowMatrix*l.size.MatrixColumns())
					if val {
						count++
					}

					result.set(colResult, rowResult, val)
				}
			}
		}
	}

	return result
}
