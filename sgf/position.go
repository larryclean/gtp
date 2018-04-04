package sgf

import (
	"fmt"
	"strings"
)

type Position struct {
	Schema   []int32
	Size     int32
	BlackCap int
	WhiteCap int
	HisNode  Node
}

func NewPosition(size int32) Position {
	position := Position{}
	position.Size = size
	position.Schema = position.CreateSchema()
	return position
}

func (p Position) Clone() (Position, error) {
	newPos := Position{}
	err := DeepCopy(&newPos, p)
	if err != nil {
		return newPos, err
	}
	return newPos, nil
}

func (p Position) CreateSchema() []int32 {
	poss := make([]int32, p.Size*p.Size)
	return poss
}
func (p Position) Neighbar4(x, y int32) []Node {
	result := make([]Node, 0)
	// up
	if y > 0 {
		result = append(result, Node{
			X: x,
			Y: y - 1,
		})
	}
	//left
	if x > 0 {
		result = append(result, Node{
			X: x - 1,
			Y: y,
		})
	}
	//down
	if y < p.Size-1 {
		result = append(result, Node{
			X: x,
			Y: y + 1,
		})
	}
	//right
	if x < p.Size-1 {
		result = append(result, Node{
			X: x + 1,
			Y: y,
		})
	}
	return result
}

func (p *Position) SetPosition(x, y, c int32) {
	if x >= 0 && y >= 0 {
		p.Schema[x*p.Size+y] = c
	}

}
func (p Position) GetPosition(x, y int32) int32 {
	if x >= 0 && y >= 0 {
		return p.Schema[x*p.Size+y]
	}
	return 0
}
func (p Position) GetCoor(x, y int32) int32 {
	return x*p.Size + y
}

// Move 落子
func (p *Position) Move(x, y, c int32) (bool, int) {
	newPos, _ := p.Clone()
	newPos.SetPosition(x, y, c)
	//return p.CheckDead(x, y, c)
	nodes := newPos.CheckDead(x, y, c)
	cnt := len(nodes)
	p.SetPosition(x, y, c)
	if cnt > 0 {
		p.CapStones(nodes)
	}
	p.ResetKO(x, y, c)
	if cnt == 1 {
		p.HisNode = nodes[0]
	}
	return cnt > 0, cnt
}

// Move 落子
func (p *Position) Play(x, y, c int32) (bool, int) {
	newPos, _ := p.Clone()
	newPos.SetPosition(x, y, c)
	nodes := newPos.CheckDead(x, y, c)
	cnt := len(nodes)
	if cnt > 0 {
		p.CapStones(nodes)
	} else {
		newPos.CalcDeadNotCap(x, y, c, nodes)
		if len(nodes) > 0 {
			return false, len(nodes)
		}
	}
	p.SetPosition(x, y, c)
	p.ResetKO(x, y, c)
	if cnt == 1 {
		p.HisNode = nodes[0]
	}
	p.SetPosition(x, y, c)

	return true, cnt
}
func (p *Position) CapStones(nodes []Node) {

	black := 0
	white := 0
	for _, v := range nodes {
		if v.C == B {
			black++
		} else if v.C == W {
			white++
		}
		p.SetPosition(v.X, v.Y, Empty)
	}
	p.BlackCap = p.BlackCap + black
	p.WhiteCap = p.WhiteCap + white
}

// ResetKO 重置打劫
func (p *Position) ResetKO(x, y, c int32) {
	if p.HisNode.C == c && (p.HisNode.X != x || p.HisNode.Y != y) {
		p.HisNode = Node{}
	}
}

// 是否打劫
func (p *Position) CheckKO(x, y, c int32, deadcount int) bool {
	if p.HisNode.X == x && p.HisNode.Y == y && p.HisNode.C == c && deadcount == 1 {
		return false
	} else {
		return true
	}
}

//校验死子
func (p *Position) CheckDead(x, y, c int32) []Node {
	otherColor := int32(Empty)
	otherColor = int32(B)
	if c == 1 {
		otherColor = int32(W)
	}
	nodes := make([]Node, 0)
	//up
	if y > 0 && p.GetPosition(x, y-1) == int32(otherColor) {
		nodes=p.CalcDeadNotCap(x, y-1, otherColor, nodes)
	}
	//left
	if x > 0 && p.GetPosition(x-1, y) == int32(otherColor) {
		nodes=p.CalcDeadNotCap(x-1, y, otherColor, nodes)
	}
	//down
	if y < p.Size-1 && p.GetPosition(x, y+1) == int32(otherColor) {
		nodes=p.CalcDeadNotCap(x, y+1, otherColor, nodes)
	}
	//right
	if x < p.Size-1 && p.GetPosition(x+1, y) == int32(otherColor) {
		nodes=p.CalcDeadNotCap(x+1, y, otherColor, nodes)
	}

	return nodes
}

//计算死子但不提子
func (p *Position) CalcDeadNotCap(x, y, c int32, nodes []Node)[]Node {
	temp_pos := NewPosition(p.Size)
	isDead := true
	temp_pos = p.FindAreaByC(temp_pos, x, y, c)
	for i := int32(0); i < p.Size; i++ {
		for j := int32(0); j < p.Size; j++ {
			if temp_pos.GetPosition(i, j) == 3 {
				isDead = false
			}
		}
	}
	if isDead {
		for i := int32(0); i < p.Size; i++ {
			for j := int32(0); j < p.Size; j++ {
				if temp_pos.GetPosition(i, j) == c {
					p.SetPosition(i, j, Empty)
					nodes = append(nodes, Node{
						X: i,
						Y: j,
						C: c,
					})
				}
			}
		}
	}
	return nodes
}

// FindAreaByC 查找区域连块逻辑
func (p Position) FindAreaByC(pos Position, x, y, c int32) Position {
	if pos.GetPosition(x, y) != c && p.GetPosition(x, y) == c {
		pos.SetPosition(x, y, c)
		//上区域连块
		if y > 0 && p.GetPosition(x, y-1) == c {
			pos = p.FindAreaByC(pos, x, y-1, c)
		} else if y > 0 && p.GetPosition(x, y-1) == Empty {
			pos.SetPosition(x, y-1, 3)
			return pos
		}
		//左区域连块
		if x > 0 && p.GetPosition(x-1, y) == c {
			pos = p.FindAreaByC(pos, x-1, y, c)
		} else if x > 0 && p.GetPosition(x-1, y) == Empty {
			pos.SetPosition(x-1, y, 3)
			return pos
		}
		//下区域连块
		if y < p.Size-1 && p.GetPosition(x, y+1) == c {
			pos = p.FindAreaByC(pos, x, y+1, c)
		} else if y < p.Size-1 && p.GetPosition(x, y+1) == 0 {
			pos.SetPosition(x, y+1, 3)
			return pos
		}
		//右区域连块
		if x < p.Size-1 && p.GetPosition(x+1, y) == c {
			pos = p.FindAreaByC(pos, x+1, y, c)
		} else if x < p.Size-1 && p.GetPosition(x+1, y) == 0 {
			pos.SetPosition(x+1, y, 3)
			return pos
		}
	}
	return pos
}

// SetDead 设置死子
func (p *Position) SetDead(nodes []Node) {
	for _, v := range nodes {
		color := p.GetPosition(v.X, v.Y)
		if color == B {
			p.SetPosition(v.X, v.Y, W)
			p.WhiteCap = p.WhiteCap + 1
		} else if color == W {
			p.SetPosition(v.X, v.Y, B)
			p.BlackCap = p.BlackCap + 1
		}
	}
}

func (p *Position) SetColor(nodes []Node, color int32) {
	for _, v := range nodes {
		if p.GetPosition(v.X, v.Y) != 0 && color != p.GetPosition(v.X, v.Y) {
			if color == B {
				p.WhiteCap = p.WhiteCap + 1
			} else if color == W {
				p.BlackCap = p.BlackCap + 1
			}
		}
		p.SetPosition(v.X, v.Y, color)
	}
}

// SetBoardPos 设置棋盘着点
func (p *Position) SetBoardPos(b string, w string) {
	blacks := strings.Split(b, ":")
	for _, v := range blacks {
		v = strings.TrimSpace(v)
		if len(v) == 2 {
			x := ToNum(v, 0)
			y := ToNum(v, 1)
			p.SetPosition(x, y, B)
		}
	}
	whites := strings.Split(w, ":")
	for _, v := range whites {
		v = strings.TrimSpace(v)
		if len(v) == 2 {
			x := ToNum(v, 0)
			y := ToNum(v, 1)
			p.SetPosition(x, y, W)
		}
	}
}

// ResetEmptyAreaToC 重置空白区域所属
func (p *Position) ResetEmptyAreaToC() {
	for i := int32(0); i < p.Size; i++ {
		for j := int32(0); j < p.Size; j++ {
			if p.GetPosition(i, j) == 0 {
				p.checkEmptyArea(i, j)
			}
		}
	}
}

//checkEmptyArea 计算空白区域所以方
func (p *Position) checkEmptyArea(x, y int32) {
	isAttach := int32(0)
	temp_pos := NewPosition(p.Size)
	temp_pos, isAttach = p.calcMyArea(temp_pos, x, y, 0)
	if isAttach == 1 || isAttach == -1 {
		for i := int32(0); i < p.Size; i++ {
			for j := int32(0); j < p.Size; j++ {
				if temp_pos.GetPosition(i, j) == 4 {
					p.SetPosition(i, j, isAttach)
				}
			}
		}
	}
}

// 计算相邻空白区域以及终止颜色
func (p Position) calcMyArea(pos Position, x, y, c int32) (Position, int32) {
	if c == 5 {
		return pos, 5
	}

	if pos.GetPosition(x, y) != 4 && p.GetPosition(x, y) == 0 {
		pos.SetPosition(x, y, 4)
		//上区域连块
		if y > 0 && p.GetPosition(x, y-1) == 0 {
			pos, c = p.calcMyArea(pos, x, y-1, c)
		} else if y > 0 && p.GetPosition(x, y-1) != 0 {
			if c == 0 {
				c = p.GetPosition(x, y-1)
			} else if c != p.GetPosition(x, y-1) {
				return pos, 5
			}
		}
		//左区域连块
		if x > 0 && p.GetPosition(x-1, y) == 0 {
			pos, c = p.calcMyArea(pos, x-1, y, c)
		} else if x > 0 && p.GetPosition(x-1, y) != 0 {
			if c == 0 {
				c = p.GetPosition(x-1, y)
			} else if c != p.GetPosition(x-1, y) {
				return pos, 5
			}
		}
		//下区域连块
		if y < p.Size-1 && p.GetPosition(x, y+1) == 0 {
			pos, c = p.calcMyArea(pos, x, y+1, c)
		} else if y < p.Size-1 && p.GetPosition(x, y+1) != 0 {
			if c == 0 {
				c = p.GetPosition(x, y+1)
			} else if c != p.GetPosition(x, y+1) {
				return pos, 5
			}
		}
		//右区域连块
		if x < p.Size-1 && p.GetPosition(x+1, y) == 0 {
			pos, c = p.calcMyArea(pos, x+1, y, c)
		} else if x < p.Size-1 && p.GetPosition(x+1, y) != 0 {
			if c == 0 {
				c = p.GetPosition(x+1, y)
			} else if c != p.GetPosition(x+1, y) {
				return pos, 5
			}
		}
	}
	return pos, c
}

// GetStones 获取各个颜色的棋子列表
func (p Position) GetStones() (blackList []string, whiteList []string) {
	blackList = make([]string, 0)
	whiteList = make([]string, 0)
	for i := int32(0); i < p.Size; i++ {
		for j := int32(0); j < p.Size; j++ {
			coon := fmt.Sprintf("%s%s", IntToChar(i), IntToChar(j))
			if p.GetPosition(i, j) == B {
				blackList = append(blackList, coon)
			} else if p.GetPosition(i, j) == W {
				whiteList = append(whiteList, coon)
			}
		}
	}
	return
}

func (p Position) PrintBoard() []string {
	result := make([]string, 0)
	for i := int32(0); i < p.Size; i++ {
		arr := make([]string, 0)
		for j := int32(0); j < p.Size; j++ {
			if p.GetPosition(j, i) == 1 {
				arr = append(arr, "X")
			} else if p.GetPosition(j, i) == -1 {
				arr = append(arr, "O")
			} else {
				arr = append(arr, "-")
			}
		}
		result = append(result, strings.Join(arr, " "))
	}
	return result
}
func (p Position) getNextMove(x, y, c int32, deadCount int, hisNode Node) (*Node, int) {
	kill, cnt := p.Move(x, y, c)
	if kill {
		if p.CheckKO(x, y, c, deadCount) && cnt > deadCount {
			return &Node{
				X: x,
				Y: y,
				C: c,
			}, cnt
		}
	}
	p.SetPosition(x, y, Empty)
	return nil, 0
}
func (op Position) CalcCap(color int32) *Node {
	p, _ := op.Clone()
	deadCount := 0
	result := &Node{}
	for i := int32(0); i < p.Size; i++ {
		for j := int32(0); j < p.Size; j++ {
			n := p.GetPosition(i, j)
			if n == Empty {
				// 上
				if j > 0 && p.GetPosition(i, j-1) != Empty {
					c := p.GetPosition(i, j-1)
					if 0-c == color {
						node, cnt := p.getNextMove(i, j, color, deadCount, op.HisNode)
						if node != nil {
							deadCount = cnt
							result = node
						}
					}

				}
				// 左
				if i > 0 && p.GetPosition(i-1, j) != Empty {
					c := p.GetPosition(i-1, j)
					if 0-c == color {
						node, cnt := p.getNextMove(i, j, color, deadCount, op.HisNode)
						if node != nil {
							deadCount = cnt
							result = node
						}
					}
				}
				// 下
				if j < p.Size-1 && p.GetPosition(i, j+1) != Empty {
					c := p.GetPosition(i, j+1)
					if 0-c == color {
						node, cnt := p.getNextMove(i, j, color, deadCount, op.HisNode)
						if node != nil {
							deadCount = cnt
							result = node
						}
					}

				}
				// 右
				if i < p.Size-1 && p.GetPosition(i+1, j) != Empty {
					c := p.GetPosition(i+1, j)
					if 0-c == color {
						node, cnt := p.getNextMove(i, j, color, deadCount, op.HisNode)
						if node != nil {
							deadCount = cnt
							result = node
						}
					}
				}
			}
		}
	}
	if deadCount > 0 {
		return result
	}
	return nil
}
