package goban

import (
	"fmt"
	"strings"
	"strconv"
)

// ToNum 根据字符串转为ASCII码
func StrToASCII(value string, index int) int32 {
	temp := string(value[index])
	xr := []rune(temp)
	return xr[0] - 97
}

//坐标转为sgf中的坐标字符串
func CoorToSgfNode(x, y int32) string {
	xChar := fmt.Sprintf("%s", string(x+ACSII))
	yChar := fmt.Sprintf("%s", string(y+ACSII))
	if x == -1 {
		xChar = "t"
	}

	if y == -1 {
		yChar = "t"
	}
	return fmt.Sprintf("%s%s", xChar, yChar)
}

//坐标转为棋盘中的坐标字符串
func CoorToBoardNode(x, y, size int32) string {
	if x >= 8 {
		x++
	}
	return strings.ToUpper(fmt.Sprintf("%s%d", string(x+ACSII), size-y))
}

// 解析leelazero 数据
func ParseBranch(log string) []map[string]interface{} {
	lines := strings.Split(log, "\n")
	result := make([]map[string]interface{},0)
	for _, v := range lines {
		if strings.Contains(v, "->") {
			item:=make(map[string]interface{})
			first := strings.Split(v, "->")
			//选点
			item["select"] = strings.TrimSpace(first[0])
			second := strings.Split(strings.TrimSpace(first[1]), "(")
			// 模拟次数
			item["times"] = strings.TrimSpace(second[0])
			// 胜利
			item["wine_rate"] = strings.TrimSpace(strings.Replace(strings.Replace(second[1], "V:", "", -1), "%)", "", -1))
			three := strings.Split(strings.Replace(second[2], "N:", "", -1), "%)")
			// 策略网络概率
			item["playout"] = strings.TrimSpace(three[0])
			four := strings.Fields(strings.TrimSpace(three[1]))
			if len(four) > 0 && four[0] == "PV:" {
				item["branch"] = four[1:]
			}
			result=append(result, item)
		}
	}
	return result
}

//解析leelazero heatmap
func ParseHeatMap(log string) ([362]float64, float64) {
	position := [362]float64{}
	wineRate := 0.0
	for x, v := range strings.Split(log, "\n") {
		lines := strings.Fields(v)
		switch len(lines) {
		case 19:
			for y, p := range lines {
				pp, _ := strconv.ParseFloat(p, 64)
				position[x+y*19] = pp
			}
		case 2:
			if lines[0] == "pass:" {
				pp, _ := strconv.ParseFloat(lines[1], 64)
				position[361] = pp
			} else if lines[0] == "winrate:" {
				rate, _ := strconv.ParseFloat(lines[1], 64)
				wineRate = rate
			}
		}
	}
	return position, wineRate
}
