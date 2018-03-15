package sgf


type Node struct {
	X         int32
	Y         int32
	C         int32
	Parent    *Node
	Childrens []*Node
	Steup     []*Node
	LastSelect int
}

func NewNode() *Node {
	return &Node{
		Childrens: make([]*Node, 0),
		Steup:     make([]*Node, 0),
	}
}

func (n Node) GetChild(i int) *Node {
	if len(n.Childrens) < i {
		return nil
	} else {
		return n.Childrens[i]
	}
}
func (n *Node) AppendChild()*Node  {
	node:=NewNode()
	node.Parent=n
	n.Childrens=append(n.Childrens, node)
	return node
}
func (n *Node) AddSetup(obj *Node)  {
	n.Steup=append(n.Steup, obj)
}

func (k Node) GetXSizeCoor(size int32) int32 {
	return k.X*size + k.Y
}
func (k Node) GetYSizeCoor(size int32) int32 {
	return k.Y*size + k.X
}
