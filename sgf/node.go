package sgf

import (
	"fmt"
)

type Node struct {
	X          int32
	Y          int32
	C          int32
	Parent     *Node `json:"-"`
	Childrens  []*Node `json:"-"`
	Steup      []*Node `json:"-"`
	LastSelect int `json:"-"`
	Info       map[string][]string `json:"-"`
}

func NewNode() *Node {
	return &Node{
		Childrens: make([]*Node, 0),
		Steup:     make([]*Node, 0),
		Info:      make(map[string][]string),
	}
}

func (n Node) GetChild(i int) *Node {
	if len(n.Childrens) < i {
		return nil
	} else {
		return n.Childrens[i]
	}
}
func (n *Node) AppendChild() *Node {
	node := NewNode()
	node.Parent = n
	n.Childrens = append(n.Childrens, node)
	return node
}
func (n *Node) AddSetup(obj *Node) {
	n.Steup = append(n.Steup, obj)
}

func (n *Node) AddTR() {
	n.AddInfo("TR", PointToStr(n.X, n.Y))
}
func (n *Node) AddInfo(k, v string) {
	aa := n.Info[k]
	aa = append(aa, v)
	n.Info[k] = aa
}
func (n *Node) AddComment(v string) {
	aa := n.Info["C"]
	if len(aa)==0{
		aa = append(aa, v)
	}else{
		aa[0]=fmt.Sprintf("%s %s",aa[0],v)
	}
	n.Info["C"] = aa
}

func (k Node) GetXSizeCoor(size int32) int32 {
	return k.X*size + k.Y
}
func (k Node) GetYSizeCoor(size int32) int32 {
	return k.Y*size + k.X
}
func (k Node) GetSgfMove() string {
	if k.C == B {
		return fmt.Sprintf(";B%s", CoorToStr(k.X, k.Y))
	} else if k.C == W {
		return fmt.Sprintf(";W%s", CoorToStr(k.X, k.Y))
	} else {
		return ""
	}
}
func (k Node) GetColor()string {
	if k.C==1{
		return "B"
	}else if k.C==-1{
		return "W"
	}
	return ""
}
