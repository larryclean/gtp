package goban

type Node struct {
	X          int32
	Y          int32
	C          int32
	Comment    string
	Parent     *Node               `json:"-"`
	Childrens  []*Node             `json:"-"`
	Steup      []*Node             `json:"-"`
	LastSelect int                 `json:"-"`
	Info       map[string][]string `json:"-"`
}

//创建空对象
func NewNode() *Node {
	return &Node{
		Childrens: make([]*Node, 0),
		Steup:     make([]*Node, 0),
		Parent:    nil,
		Info:      make(map[string][]string),
	}
}

//获取子节点
func (n Node) GetChild(i int) *Node {
	if len(n.Childrens)==0{
		return nil
	}
	if len(n.Childrens) <=i {
		i=0
	}
	n.LastSelect=i
	return n.Childrens[n.LastSelect]
}

//追加子节点
func (n *Node) AppendChild() *Node {
	node := NewNode()
	node.Parent = n
	n.Childrens = append(n.Childrens, node)
	return node
}

// 添加AB/AW标签
func (n *Node) AddSetup(obj *Node) {
	n.Steup = append(n.Steup, obj)
}
