package sgf

type KNode struct {
	X int32
	Y int32
	C int32
}

func (k KNode) GetXSizeCoor(size int32) int32 {
	return k.X*size + k.Y
}
func (k KNode) GetYSizeCoor(size int32) int32 {
	return k.Y*size + k.X
}
