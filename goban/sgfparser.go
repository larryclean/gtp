package goban

import (
	"regexp"
	"strconv"
	"strings"
)

func setup(kifu *Kifu, node *Node, value []string, indent string) {
	c := W
	if indent == "AB" {
		c = B
	} else {
		c = W
	}
	for _,v:=range value{
		node.AddSetup(&Node{
			C: int32(c),
			X: StrToASCII(v, 0),
			Y: StrToASCII(v, 1),
		})
	}
}
func node(kifu *Kifu, node *Node, value []string, indent string) {
	c := W
	if indent == "B" {
		c = B
	} else {
		c = W
	}
	if len(value) == 0 || (kifu.Size <= 19 && value[0] == "tt") {
		node.C = int32(c)
		node.X = int32(-1)
		node.Y = int32(-1)
	} else {
		node.C = int32(c)
		node.X = StrToASCII(value[0], 0)
		node.Y = StrToASCII(value[0], 1)
	}
	kifu.NodeCount += 1
}
func kifuInfo(kifu *Kifu, node *Node, value []string, indent string) {
	if indent == "SZ" {
		size, err := strconv.Atoi(value[0])
		if err != nil {
			size = 19
		}
		kifu.Size = int32(size)
	}
	if indent == "KM" {
		km, err := strconv.ParseFloat(value[0], 32)
		if err != nil {
			km = 7.5
		}
		kifu.Komi = float32(km)
	}
	if indent == "HA" {
		ha, err := strconv.Atoi(value[0])
		if err != nil {
			ha = 0
		}
		kifu.Handicap = ha
	}
}
func comment(kifu *Kifu, node *Node, value []string, indent string) {
	node.Comment = strings.Join(value, "")
}

var properties = map[string]func(kifu *Kifu, node *Node, value []string, indent string){
	"AB": setup,
	"AW": setup,
	"B":  node,
	"W":  node,
	"SZ": kifuInfo,
	"KM": kifuInfo,
	"HA": kifuInfo,
	"C":  comment,
}

//解析SGF 生成树状结构数据
func ParseSgf(sgf string) Kifu {
	stack := make([]*Node, 0)
	var node *Node
	kifu := NewKifu(sgf)
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
			ll := len(stack)
			node = stack[ll-1]
			if ll > 1 {
				stack = stack[:ll-1]
			}
			continue
		}
		if (node == nil) {
			node = kifu.Root
		} else {
			node = node.AppendChild()
		}
		props := reg_node.FindAllString(v, -1)
		for _, v1 := range props {
			indent := reg_indent.FindString(v1)
			vals := reg_props.FindAllString(v1, -1)
			for i, v2 := range vals {
				v2 = reg_re.ReplaceAllString(v2[1:len(v2)-1], "")
				vals[i] = v2
			}
			val, ok := properties[indent]
			if ok {
				val(&kifu, node, vals, indent)
			} else {
				node.Info[indent]=vals
			}

		}

	}
	return kifu
}
