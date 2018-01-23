package sgf

import (
	"regexp"
	"strconv"
	"strings"
)

// 解析每个子的坐标
func parseNode(kifu Kifu, value string, indent string) KNode {
	node := KNode{}
	c := W
	if indent == "B" || indent == "AB" {
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
	return node
}
func ParseSgf(sgf string) Kifu {
	stack := make([]KNode, 0, 5)
	kifu := Kifu{
		Komi:     7.5,
		Handicap: 0,
		Size:     19,
		CurColor: B,
		SgfStr:   sgf,
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
			continue
		} else if v == ")" {
			continue
		}

		props := reg_node.FindAllString(v, -1)
		for _, v1 := range props {
			indent := reg_indent.FindString(v1)
			vals := reg_props.FindAllString(v1, -1)
			for i, v2 := range vals {
				v2 = reg_re.ReplaceAllString(v2[1:len(v2)-1], "")
				vals[i] = v2
			}
			if indent == "B" || indent == "W" {
				for _, v10 := range vals {
					node := parseNode(kifu, v10, indent)
					stack = append(stack, node)
					kifu.NodeCount += 1
				}
			}
			if indent == "AW" || indent == "AB" {
				for _, v10 := range vals {
					node := parseNode(kifu, v10, indent)
					stack = append(stack, node)
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
	kifu.Nodes = stack
	if kifu.Size == 0 {
		kifu.Size = 19
	}
	return kifu
}

// ParseABAW 转化形势判断
func ParseABAW(ab, aw string, komi float32, size int32, handicap int) Kifu {
	stack := make([]KNode, 0, 5)
	kifu := Kifu{
		Komi:     komi,
		Handicap: handicap,
		Size:     size,
		CurColor: B,
		SgfStr:   "",
	}
	bList := strings.Split(ab, ":")
	wList := strings.Split(aw, ":")
	step := len(bList)
	if step < len(wList) {
		step = len(wList)
	}
	for i := 0; i < step; i++ {
		if i >= len(bList) {
			node := parseNode(kifu, "tt", "B")
			stack = append(stack, node)
			kifu.NodeCount += 1
		} else {
			node := parseNode(kifu, bList[i], "B")
			stack = append(stack, node)
			kifu.NodeCount += 1
		}
		if i >= len(wList) {
			node := parseNode(kifu, "tt", "W")
			stack = append(stack, node)
			kifu.NodeCount += 1
		} else {
			node := parseNode(kifu, wList[i], "W")
			stack = append(stack, node)
			kifu.NodeCount += 1
		}

	}
	kifu.Nodes = stack
	return kifu
}
