package sgf

import (
	"regexp"
	"strconv"
)

func abaw(kifu Kifu, node *Node, value string, indent string) {
	c := W
	if indent == "AB" {
		c = B
	} else {
		c = W
	}
	node1 := &Node{}
	if len(value) == 0 || (kifu.Size <= 19 && value == "tt") {
		node1.C = int32(c)
		node1.X = int32(-1)
		node1.Y = int32(-1)
	} else {
		node1.C = int32(c)
		node1.X = ToNum(value, 0)
		node1.Y = ToNum(value, 1)
	}
	node.AddSetup(node1)
}
func bw(kifu Kifu, node *Node, value string, indent string) {
	c := W
	if indent == "B" {
		c = B
	} else {
		c = W
	}
	if len(value) == 0 || (kifu.Size <= 19 && value == "tt") {
		node.C = int32(c)
		node.X = int32(-1)
		node.Y = int32(-1)
	} else {
		node.C = int32(c)
		node.X = ToNum(value, 0)
		node.Y = ToNum(value, 1)
	}
}

var properties = map[string]func(kifu Kifu, node *Node, value string, indent string){
	"AB": abaw,
	"AW": abaw,
	"B":  bw,
	"W":  bw,
}

//解析SGF 生成树状结构数据
func ParseSgf(sgf string) Kifu {
	stack := make([]*Node, 0)
	var node *Node
	kifu := Kifu{
		Komi:     7.5,
		Handicap: 0,
		Size:     19,
		CurColor: B,
		SgfStr:   sgf,
		Root:     NewNode(),
	}
	//解析SGF文件
	reg_seq := regexp.MustCompile(pat_seq)
	reg_node := regexp.MustCompile(pat_node)
	reg_indent := regexp.MustCompile(pat_ident)
	reg_props := regexp.MustCompile(pat_props)
	reg_re := regexp.MustCompile(`\\(!\\)`)
	sequence := reg_seq.FindAllString(sgf, -1)
	for _, v := range sequence {
		if v == "(" {
			stack = append(stack, node)
			continue
		} else if v == ")" {
			ll:=len(stack)
			node = stack[ll-1]
			if ll>1{
				stack=stack[:ll-1]
			}
			continue
		}
 		if (node==nil){
 			node=kifu.Root
		}else{
			node=node.AppendChild()
		}
		props := reg_node.FindAllString(v, -1)
		for _, v1 := range props {
			indent := reg_indent.FindString(v1)
			vals := reg_props.FindAllString(v1, -1)
			for i, v2 := range vals {
				v2 = reg_re.ReplaceAllString(v2[1:len(v2)-1], "")
				vals[i] = v2
			}
			if indent == "B" || indent == "W"||indent == "AW" || indent == "AB" {
				for _, v10 := range vals {
					properties[indent](kifu,node,v10,indent)
					kifu.NodeCount += 1
				}
			}
			if indent == "SZ" {
				size, err := strconv.Atoi(vals[0])
				if err != nil {
					size = 19
				}
				kifu.Size = int32(size)
			}
			if indent == "KM" {
				km, err := strconv.ParseFloat(vals[0], 32)
				if err != nil {
					km = 7.5
				}
				kifu.Komi = float32(km)
			}
			if indent == "HA" {
				ha, err := strconv.Atoi(vals[0])
				if err != nil {
					ha = 0
				}
				kifu.Handicap = ha
			}
		}

	}
	if kifu.Size == 0 {
		kifu.Size = 19
	}
	return kifu
}

