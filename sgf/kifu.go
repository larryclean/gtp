package sgf

import "fmt"

type Kifu struct {
	Nodes     []KNode
	Steps     []KNode
	Size      int32
	Handicap  int
	Komi      float32
	CurColor  int32
	SgfStr    string
	NodeCount int
	CurPos    Position
}

func (k *Kifu) GoTo(move int) Position {
	pos := NewPosition(k.Size)
	if move > k.NodeCount {
		move = k.NodeCount
	}
	for i := 0; i < move; i++ {
		node := k.Nodes[i]
		k.CurColor = 0 - node.C
		if node.X == -1 && node.Y == -1 {
			continue
		}
		pos.Move(node.X, node.Y, node.C)
		pos.ResetKO(node.X, node.Y, node.C)
	}
	k.CurPos = pos
	return pos
}

func (k *Kifu) Last() Position {
	return k.GoTo(100000)
}
func(k *Kifu) AppendNode(node KNode)  {
	k.Nodes=append(k.Nodes,node)
}
func (k *Kifu) ToSGF(move int) (sgf string) {
	sgfHeader := fmt.Sprintf(";SZ[%d]KM[%d]HA[%d]", k.Size,k.Komi,k.Handicap)
	if move > len(k.Nodes) || move == -1 {
		move = len(k.Nodes)
	}
	for i := 0; i < move; i++ {
		sgf += NodeToString(k.Nodes[i])
	}
	sgf = fmt.Sprintf("(%s%s)", sgfHeader, sgf)
	return
}
