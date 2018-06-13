package goban

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
	CurNode   *Node
	IsKo      bool
	liberty   Node //劫
}

func NewKifu(sgf string) Kifu {
	return Kifu{
		Komi:     7.5,
		Handicap: 0,
		Size:     19,
		CurColor: B,
		SgfStr:   sgf,
		Root:     NewNode(),
		IsKo:     false,
		liberty:  Node{},
	}
}
func (k *Kifu) GoTo(step int) {
	k.CurPos = NewPosition(k.Size)
	if step > k.NodeCount || step == -1 {
		step = k.NodeCount
	}
	k.CurNode=k.Root
	node := k.Root
	for i := 0; i < step; i++ {
		if len(node.Steup) > 0 {
			for _, v := range node.Steup {
				k.CurPos.SetColor(v.X, v.Y, v.C)
			}
		}

		temp := node.GetChild(node.LastSelect)
		if temp!=nil{
			k.Move(temp.X, temp.Y, temp.C)
			k.CurNode = temp
			k.CurColor = -node.C
			node = temp
		}
	}
}
func (k *Kifu) Last() {
	k.GoTo(-1)
}
// 落子算法逻辑
func (k *Kifu) Move(x, y, c int32) bool {
	if k.IsBoard(x, y) {
		return false
	}
	pos := k.CurPos.Clone()
	pos.SetColor(x, y, c)
	//判断是否为打劫
	deads := pos.GetDeadByPointColor(x, y, 0-c)
	if len(deads) == 0 {
		//如果无法提对方子判断自己是否为死子则阻止落子
		deads = pos.CalcDeadNode(x, y, c)
		if len(deads) > 0 {
			return false
		}
	}
	if len(deads) == 1 {
		if k.CheckKo(x,y,c) {
			return false
		}
		k.IsKo = true
		k.liberty = deads[0]
	}else{
		k.IsKo = false
	}
	k.CurPos.SetColor(x, y, c)
	k.CurPos.Cap(deads)
	k.CurPath++
	return true
}
func (k *Kifu) Play(x, y, c int32) bool {
	if k.Move(x, y, c) {
		k.CurNode = k.CurNode.AppendChild()
		k.CurNode.X = x
		k.CurNode.Y = y
		k.CurNode.C = c
		k.NodeCount++
		k.CurColor = 0 - c
		return true
	}
	return false

}
func (k Kifu) IsBoard(x, y int32) bool {
	return k.CurPos.GetColor(x, y) != Empty
}
func (k Kifu) CheckKo(x,y,c int32) bool {
	if k.IsKo {
		if k.liberty.X == x && k.liberty.Y == y && k.liberty.C == c{
			return true
		}
	}
	return false
}

//生成SGF
func (k Kifu) ToSgf() string {
	sgf := fmt.Sprintf("(;SZ[%v]KM[%v]HA[%v]", k.Size, k.Komi, k.Handicap)
	node := k.Root
	sgf += k.SgfWriteNode(node)
	sgf += ")"
	return sgf
}

// 解析节点信息
func (k Kifu) toNodeInfo(node *Node) string {
	sgf := ""
	if node.Comment != "" {
		sgf += fmt.Sprintf("C[%s]", node.Comment)
	}
	for k, v := range node.Info {
		sgf += fmt.Sprintf("%s", k)
		for _, v1 := range v {
			sgf += fmt.Sprintf("[%s]", v1)
		}
	}
	return sgf
}

// 解析AB,AW点
func (k Kifu) toSetup(setups []*Node) string {
	sgf := ""
	ab := make([]string, 0)
	aw := make([]string, 0)
	for _, v := range setups {
		if v.C == B {
			ab = append(ab, fmt.Sprintf("[%s]", CoorToSgfNode(v.X, v.Y)))
		} else if v.C == W {
			aw = append(aw, fmt.Sprintf("[%s]", CoorToSgfNode(v.X, v.Y)))
		}
	}
	if len(ab) > 0 {
		sgf += fmt.Sprintf("AB%s", strings.Join(ab, ""))
	}
	if len(aw) > 0 {
		sgf += fmt.Sprintf("AW%s", strings.Join(aw, ""))
	}
	return sgf
}

// 正序解析节点
func (k Kifu) SgfWriteNode(node *Node) string {
	sgf := ""
	if node.C != Empty {
		if node.C == B {
			sgf += fmt.Sprintf(";B[%s]", CoorToSgfNode(node.X, node.Y))
		} else if node.C == W {
			sgf += fmt.Sprintf(";W[%s]", CoorToSgfNode(node.X, node.Y))
		}
	}
	if len(node.Steup) > 0 {
		sgf += k.toSetup(node.Steup)
	}
	sgf+=k.toNodeInfo(node)
	if len(node.Childrens) == 1 {
		sgf += k.SgfWriteNode(node.Childrens[0])
	} else if len(node.Childrens) > 1 {
		for _, v := range node.Childrens {
			sgf += k.SgfWriteVariantion(v)
		}
	}
	return sgf
}

// 带有变化分支生成
func (k Kifu) SgfWriteVariantion(node *Node) string {
	sgf := "("
	sgf += k.SgfWriteNode(node)
	sgf += ")"
	return sgf
}

// 根据当前节点生成SGF获取一条分支
func (k Kifu) ToCurSgf() string {
	sgf := fmt.Sprintf("(;SZ[%v]KM[%v]HA[%v]", k.Size, k.Komi, k.Handicap)
	sgf+=k.toNodeInfo(k.Root)+k.RefletSgfWriteNode(k.CurNode)
	sgf += ")"
	return sgf
}

// 方向解析节点
func (k Kifu) RefletSgfWriteNode(node *Node) string {
	sgf := ""
	if node.C != Empty {
		if node.C == B {
			sgf += fmt.Sprintf(";B[%s]", CoorToSgfNode(node.X, node.Y))
		} else if node.C == W {
			sgf += fmt.Sprintf(";W[%s]", CoorToSgfNode(node.X, node.Y))
		}
	}
	if len(node.Steup) > 0 {
		sgf += k.toSetup(node.Steup)
	}

	if node.Parent != nil {
		sgf = k.RefletSgfWriteNode(node.Parent) + sgf
	}
	return sgf
}
