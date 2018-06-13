package goban

import (
	"strings"
	"fmt"
	"bytes"
	"encoding/gob"
)

type Position struct {
	Schema   []int32
	Size     int32
	BlackCap int
	WhiteCap int
	HisNode  Node
}

// 创建position对象
func NewPosition(size int32) Position {
	position := Position{}
	position.Size = size
	position.Schema = position.CreateSchema()
	return position
}

// position坐标规则x*size+y
func (p Position) GetPos(x, y int32) int32 {
	return x*p.Size + y
}

// 克隆POSITION对象
func (p *Position) Clone() (*Position) {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(p);

	pos:=&Position{}
	gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(pos)
	return pos
}

//创建schema
func (p Position) CreateSchema() []int32 {
	poss := make([]int32, p.Size*p.Size)
	return poss
}

// 设置POSITION坐标对应颜色 Position规则为x*size+y
func (p *Position) SetColor(x, y, c int32) {
	if x >= 0 && y >= 0 && x <= p.Size-1 && y <= p.Size-1 {
		p.Schema[p.GetPos(x, y)] = c
	}
}

// 获取position坐标对应的颜色 Position规则为x*size+y
func (p Position) GetColor(x, y int32) int32 {
	if x >= 0 && y >= 0 {
		return p.Schema[p.GetPos(x, y)]
	}
	return 0
}

// 遍历x,y
func (p Position) ForeachXY(cb func(x, y int32)) {
	for i := int32(0); i < p.Size; i++ {
		for j := int32(0); j < p.Size; j++ {
			cb(i, j)
		}
	}
}

//获取对应坐标的四领域,并回调触发
func (p Position) Neighbor4(x, y int32, cb func(x, y int32)) {
	// up
	if y > 0 {
		cb(x, y-1)
	}
	//left
	if x > 0 {
		cb(x-1, y)
	}
	//down
	if y < p.Size-1 {
		cb(x, y+1)
	}
	//right
	if x < p.Size-1 {
		cb(x+1, y)
	}
}

// 根据坐标和颜色获取是否可提子
func (p *Position) Cap(nodes []Node) {
	for _, v := range nodes {
		p.SetColor(v.X, v.Y, Empty)
	}
}

// 根据坐标和颜色获取是否可提子
func (p *Position) GetDeadByPointColor(x, y, c int32) []Node {
	nodes := make([]Node, 0)
	p.Neighbor4(x, y, func(x, y int32) {
		if p.GetColor(x, y) == c {
			nodes = append(nodes, p.CalcDeadNode(x, y, c)...)
		}
	})
	return nodes
}

//计算死子但不提子
func (p *Position) CalcDeadNode(x, y, c int32) []Node {
	//新建一个计算的POSITION用于判断死子
	calcPos := NewPosition(p.Size)
	isDead := true
	nodes := make([]Node, 0)
	calcPos = p.FindAreaByC(calcPos, x, y, c)
	//判断是否可提子
	p.ForeachXY(func(x, y int32) {
		if calcPos.GetColor(x, y) == 3 {
			isDead = false
		}
	})
	//如果可提子进入提子，但是不动原始数据
	if isDead {
		p.ForeachXY(func(i, j int32) {
			if calcPos.GetColor(i, j) == c {
				p.SetColor(i, j, Empty)
				nodes = append(nodes, Node{
					X: i,
					Y: j,
					C: c,
				})
			}
		})
	}
	return nodes
}

// FindAreaByC 查找区域连块逻辑
func (p Position) FindAreaByC(pos Position, x, y, c int32) Position {
	if pos.GetColor(x, y) != c && p.GetColor(x, y) == c {
		pos.SetColor(x, y, c)
		//上区域连块
		if y > 0 && p.GetColor(x, y-1) == c {
			pos = p.FindAreaByC(pos, x, y-1, c)
		} else if y > 0 && p.GetColor(x, y-1) == Empty {
			pos.SetColor(x, y-1, 3)
			return pos
		}
		//左区域连块
		if x > 0 && p.GetColor(x-1, y) == c {
			pos = p.FindAreaByC(pos, x-1, y, c)
		} else if x > 0 && p.GetColor(x-1, y) == Empty {
			pos.SetColor(x-1, y, 3)
			return pos
		}
		//下区域连块
		if y < p.Size-1 && p.GetColor(x, y+1) == c {
			pos = p.FindAreaByC(pos, x, y+1, c)
		} else if y < p.Size-1 && p.GetColor(x, y+1) == 0 {
			pos.SetColor(x, y+1, 3)
			return pos
		}
		//右区域连块
		if x < p.Size-1 && p.GetColor(x+1, y) == c {
			pos = p.FindAreaByC(pos, x+1, y, c)
		} else if x < p.Size-1 && p.GetColor(x+1, y) == 0 {
			pos.SetColor(x+1, y, 3)
			return pos
		}
	}
	return pos
}
func (p Position) ShowBoard() string {
	boards:=make([]string,p.Size)
	p.ForeachXY(func(x, y int32) {
		color:=p.GetColor(x,y)
		str:="."
		switch color {
		case B:
			str="X"
		case W:
			str="O"
		}
		boards[y]=fmt.Sprintf("%s%+3v",boards[y],str)
	})
	return strings.Join(boards,"\n")
}
