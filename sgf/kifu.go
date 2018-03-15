package sgf

import (
	"fmt"
	"strings"
)

type Kifu struct {
	Root      *Node
	Size      int32
	Handicap  int
	Komi      float32
	CurColor  int32
	SgfStr    string
	NodeCount int
	CurPos    Position
	CurPath   int
}

func (k *Kifu) GoTo(move int) Position {
	pos := NewPosition(k.Size)
	if move > k.NodeCount || move == -1 {
		move = k.NodeCount
	}
	node := k.Root
	for i := 0; i < move; i++ {
		if len(node.Steup) > 0 {
			for _, v := range node.Steup {
				pos.SetPosition(v.X, v.Y, v.C)
			}
		}
		if len(node.Childrens) > node.LastSelect {
			temp := node.GetChild(node.LastSelect)
			pos.Move(temp.X, temp.Y, temp.C)
			pos.ResetKO(temp.X, temp.Y, temp.C)
			node = temp
		} else {
			break;
		}
	}
	k.CurPath = move
	k.CurPos = pos
	return pos
}

func (k *Kifu) Last() Position {
	return k.GoTo(-1)
}
func (k Kifu) ToSgf() string {
	sss := fmt.Sprintf("(;SZ[%v]KM[%v]HA[%v]", k.Size, k.Komi, k.Handicap)
	node := k.Root
	//sss = getSetup(node.Steup, sss)
	sss += k.WriteNode(node, "")
	sss += ")"
	return sss
}
func (k Kifu) WriteNode(node *Node, s string) string {
	if (node.C != Empty) {
		if node.C == B {
			s += fmt.Sprintf(";B%s", CoorToStr(node.X, node.Y))
		} else if node.C == W {
			s += fmt.Sprintf(";W%s", CoorToStr(node.X, node.Y))
		}
	}
	s += getSetup(node.Steup, "")
	cnt := len(node.Childrens)
	if cnt == 1 {
		s = k.WriteNode(node.GetChild(0), s)
	} else if cnt > 1 {
		for _, v := range node.Childrens {
			s = k.WriteVarian(v, s)
		}
	}
	return s
}
func (k Kifu) WriteVarian(node *Node, s string) string {
	s += "("
	s = k.WriteNode(node, s)
	s += ")"
	return s
}
func (k Kifu) GetCleanSgf() []string {
	result := make([]string, 0)
	sss := fmt.Sprintf("(;SZ[%v]KM[%v]HA[%v]", k.Size, k.Komi, k.Handicap)
	node := k.Root
	sss = getSetup(node.Steup, sss)
	result = k.ForChildrens(node, sss, result)
	return result
}
func (k Kifu) ForChildrens(node *Node, s string, result []string) []string {
	cnt := len(node.Childrens)
	if cnt == 0 {
		s += ")"
		result = append(result, s)
		return result
	}
	for i := 0; i < cnt; i++ {
		tn := node.GetChild(i)
		temp := s
		if tn.C == B {
			temp += fmt.Sprintf(";B%s", CoorToStr(tn.X, tn.Y))
		} else if tn.C == W {
			temp += fmt.Sprintf(";W%s", CoorToStr(tn.X, tn.Y))
		}
		result = k.ForChildrens(tn, temp, result)
	}
	return result
}
func getSetup(setups []*Node, s string) string {
	ab := make([]string, 0)
	aw := make([]string, 0)
	for _, v := range setups {
		if v.C == B {
			ab = append(ab, CoorToStr(v.X, v.Y))
		} else if v.C == W {
			aw = append(aw, CoorToStr(v.X, v.Y))
		}
	}
	if len(ab) > 0 {
		s = fmt.Sprintf("%sAB%s", s, strings.Join(ab, ""))
	}
	if len(aw) > 0 {
		s = fmt.Sprintf("%sAW%s", s, strings.Join(aw, ""))
	}
	return s
}
