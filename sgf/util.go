package sgf

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
	"regexp"
)

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// ToNum 根据sgf字母转为数字
func ToNum(value string, index int) int32 {
	temp := string(value[index])
	xr := []rune(temp)
	return xr[0] - 97
}

// IntToChar 数字转为sgf识别字母
func IntToChar(x int32) string {
	x = x + 97
	return fmt.Sprintf("%s", string(x))
}
func IntToAIChar(x int32) string {
	x = x + 97
	if x >= 105 {
		x++
	}
	return fmt.Sprintf("%s", string(x))
}
func CoorToXY(index, size int32) (int32, int32) {
	y := int32(index / size)
	x := int32(index % size)
	return x, y
}
func NodeToString(node Node) string {
	xChar := IntToChar(node.X)
	yChar := IntToChar(node.Y)
	if node.X == -1 {
		xChar = "t"
	}

	if node.Y == -1 {
		yChar = "t"
	}
	if node.C == B {
		return fmt.Sprintf(";B[%s%s]", xChar, yChar)
	} else if node.C == W {
		return fmt.Sprintf(";W[%s%s]", xChar, yChar)
	}
	return ""
}
func CoorToStr(x, y int32) string {
	xChar := IntToChar(x)
	yChar := IntToChar(y)
	if x == -1 {
		xChar = "t"
	}

	if y == -1 {
		yChar = "t"
	}
	return fmt.Sprintf("[%s%s]", xChar, yChar)
}
func PointToStr(x, y int32) string {
	xChar := IntToChar(x)
	yChar := IntToChar(y)
	if x == -1 {
		xChar = "t"
	}

	if y == -1 {
		yChar = "t"
	}
	return fmt.Sprintf("%s%s", xChar, yChar)
}

// SaveStringToPath 保存字符串到文件
func SaveStringToPath(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err1 := f.WriteString(content)
	if err1 != nil {
		return err1
	}
	w := bufio.NewWriter(f)
	_, err2 := w.WriteString(content)
	if err2 != nil {
		return err2
	}
	w.Flush()
	return nil
}

// RemoveByPath 更加路径删除文件
func RemoveByPath(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func ParseRegion(m string) []Region {
	r1 := strings.Split(m, ";")
	regions := make([]Region, 0)
	for _, v := range r1 {
		r2 := strings.Split(v, ":")
		region := Region{
			Limit: 1,
			MMax:  1,
			MMin:  1,
			SMax:  1,
			SMin:  1,
		}
		for i, v1 := range r2 {
			if i == 0 {
				hands, _ := strconv.Atoi(v1)
				region.Limit = hands
			} else if i == 1 {
				max, _ := strconv.ParseFloat(v1, 64)
				region.MMax = max
				region.SMax = max
			} else if i == 2 {
				min, _ := strconv.ParseFloat(v1, 64)
				region.MMin = min
				region.SMin = min
			} else if i == 3 {
				max, _ := strconv.ParseFloat(v1, 64)
				region.SMax = max
			} else if i == 4 {
				min, _ := strconv.ParseFloat(v1, 64)
				region.SMin = min
			}

		}
		regions = append(regions, region)
	}
	return regions
}

// parseMove 解析AI节点
func ParseMove(ss string, size int32) ([]Node, error) {
	result := make([]Node, 0, 5)
	pat := `\s+`
	reg := regexp.MustCompile(pat)
	list := reg.Split(ss, -1)
	for _, v := range list {
		if len(v) > 0 {
			temp := strings.ToLower(string(v[0]))
			if v == "PASS" {
				pos := Node{
					X: -1,
					Y: -1,
				}
				result = append(result, pos)
				continue
			}
			xr := []rune(temp)
			xInt := xr[0]
			if xInt > 105 {
				xInt = xInt - 1
			}

			y, err := strconv.Atoi(string(v[1:]))
			//utils.CheckError(err)
			if err != nil {
				return nil, err
			}
			yInt := size - int32(y)
			pos := Node{
				X: xInt - 97,
				Y: yInt,
			}
			result = append(result, pos)
		}
	}
	return result, nil
}

// ParseScore 解析AI Result
func ParseResult(list []string) []string {
	value := make([]string, 0, 5)
	temp := ""
	for i, v := range list {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			if strings.Contains(v, "=") {
				if i > 0 {
					value = append(value, temp)
					temp = ""
				}
				temp = strings.TrimSpace(strings.Replace(v, "=", "", -1))
			} else {
				temp = fmt.Sprintf("%s %s", temp, strings.TrimSpace(v))
			}
		}
	}
	value = append(value, temp)
	return value
}

// 解析leelazero 数据
func ParseBranch(log string) map[string]interface{} {
	lines := strings.Split(log, "\n")
	result := make(map[string]interface{})
	for _, v := range lines {
		if strings.Contains(v, "->") {
			first := strings.Split(v, "->")
			//选点
			result["select"] = strings.TrimSpace(first[0])
			second := strings.Split(strings.TrimSpace(first[1]), "(")
			// 模拟次数
			result["times"] = strings.TrimSpace(second[0])
			// 胜利
			result["wine_rate"] = strings.TrimSpace(strings.Replace(strings.Replace(second[1], "V:", "", -1), "%)", "", -1))
			three := strings.Split(strings.Replace(second[2], "N:", "", -1), "%)")
			// 策略网络概率
			result["playout"] = strings.TrimSpace(three[0])
			four := strings.Fields(strings.TrimSpace(three[1]))
			if len(four) > 0 && four[0] == "PV:" {
				result["branch"] = four[1:]
			}
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
