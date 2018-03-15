package sgf

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	return fmt.Sprintf("%s", string(x+97))
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
