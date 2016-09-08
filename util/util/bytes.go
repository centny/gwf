package util

// SplitTwo split bytes slice to two bytes slice by index
func SplitTwo(bys []byte, idx int) ([]byte, []byte) {
	return bys[:idx], bys[idx:]
}

// SplitThree split bytes slice to thrd bytes slice by index
func SplitThree(bys []byte, idxa, idxb int) ([]byte, []byte, []byte) {
	return bys[:idxa], bys[idxa:idxb], bys[idxb:]
}
