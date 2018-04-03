package sgf

import (
	"fmt"
	"math"
	"strings"
)

const (
	N_MODE = iota //正常形势判断
	C_MODE        //颜色偏好
	E_MODE        //终局数字
	ANFLAG        //AN结果标记
	GOFLAG        //正常标记
)

type ScoreConfig struct {
	Normal     string `toml:"n"`
	Color      string `toml:"c"` //乐观
	End        string `toml:"e"`
	HandsCount int    `toml:"hands"`
	Empty      int    `toml:"empty"`
}
type Region struct {
	Limit                  int
	MMax, MMin, SMax, SMin float64
}

func (self ScoreConfig) GetRegion(mode int, hands int) Region {
	region := Region{
		Limit: 0,
		MMax:  1,
		MMin:  1,
		SMax:  1,
		SMin:  1,
	}
	ss := self.End
	if mode == N_MODE {
		ss = self.Normal
	} else if mode == C_MODE {
		ss = self.Color
	} else if mode == E_MODE {
		ss = self.End
	}
	regions := ParseRegion(ss)
	for _, v := range regions {
		if hands >= v.Limit {
			region = v
		} else {
			break
		}
	}
	return region
}
func NewScoreConfig() ScoreConfig {
	return ScoreConfig{
		Normal:     "0:0.8:1",
		Color:      "0:1:1:0.8:1",
		End:        "0:0.95:0.5:0.95:0.5",
		HandsCount: 150,
		Empty:      4,
	}
}

type AreaResult struct {
	ANResult []float64
	CurPos   Position
	JPResult []float64
	CNResult []float64
	Flag     int
	Kifu     Kifu
}

type ScoreUtil struct {
	AreaResult *AreaResult
	Kifu       Kifu
	Config     ScoreConfig
}
type ScoreResult struct {
	Pos            interface{} `json:"pos"`
	Komi           float32     `json:"km"`
	Handicap       int         `json:"ha"`
	JP             string      `json:"jp"`
	CN             string      `json:"cn"`
	BlackArea      int         `json:"ba"`
	WhiteArea      int         `json:"wa"`
	BlackTerritory string      `json:"bt"`
	WhiteTerritory string      `json:"wt"`
	BlackPrisoners int         `json:"bp"`
	WhitePrisoners int         `json:"wp"`
	BlackCap       int         `json:"bc"`
	WhiteCap       int         `json:"wc"`
	Empty          int         `json:"dame"`
	Black          string      `json:"b"`
	White          string      `json:"w"`
}

// Score 计算
func (self *ScoreUtil) Score(mode int, color int32) ScoreResult {
	s := ScoreResult{}
	if self.AreaResult.Flag == GOFLAG {
		self.ZHUScore()
	} else {
		self.ANScore(mode, color)
	}
	s.Komi = self.Kifu.Komi
	s.Handicap = self.Kifu.Handicap
	//提子数
	s.BlackCap = self.Kifu.CurPos.BlackCap
	s.WhiteCap = self.Kifu.CurPos.WhiteCap
	//目数
	black_territory := float64(0)
	white_territory := float64(0)
	ww_territory := 0
	bb_territory := 0
	//死子
	white_prisoners := 0
	black_prisoners := 0
	//单官
	empty := 0
	//计算系数
	black_area := make([]string, 0) //all
	white_area := make([]string, 0) //all
	for i := int32(0); i < self.Kifu.Size; i++ {
		for j := int32(0); j < self.Kifu.Size; j++ {
			coor := self.Kifu.CurPos.GetCoor(i, j)
			o := self.Kifu.CurPos.GetPosition(i, j)
			cn := self.AreaResult.CNResult[coor]
			jp := self.AreaResult.JPResult[coor]
			coon := fmt.Sprintf("%s%s", IntToChar(i), IntToChar(j))
			if cn > 0 {
				black_area = append(black_area, coon)
				if o == W {
					white_prisoners = white_prisoners + 1
					bb_territory += 1
				}
				if o == Empty {
					bb_territory += 1
				}

			} else if cn < 0 {
				white_area = append(white_area, coon)
				if o == B {
					black_prisoners = black_prisoners + 1
					ww_territory += 1
				}
				if o == Empty {
					ww_territory += 1
				}

			} else {
				empty += 1
			}
			if jp > 0 {
				black_territory += jp
			} else if jp < 0 {
				white_territory += math.Abs(jp)
			}
		}
	}
	s.BlackArea = len(black_area)
	s.WhiteArea = len(white_area)
	s.Black = strings.Join(black_area, ":")
	s.White = strings.Join(white_area, ":")
	s.BlackTerritory = fmt.Sprintf("%.2f", black_territory)
	s.WhiteTerritory = fmt.Sprintf("%.2f", white_territory)
	s.WhitePrisoners = white_prisoners
	s.BlackPrisoners = black_prisoners
	s.Empty = empty
	sear := s.Komi
	if s.Handicap > 0 {
		if s.Handicap == 1 {
			sear = 0
		} else {
			sear = float32(s.Handicap)
		}
	}
	if mode == E_MODE || (mode == N_MODE && empty <= self.Config.Empty && self.Kifu.NodeCount >= self.Config.HandsCount) || self.Kifu.Size < 19 {
		board := float32(self.Kifu.Size*self.Kifu.Size) / 2
		br := float32(s.BlackArea) + float32(s.Empty)/2 - sear/2
		s.CN = fmt.Sprintf("%.2f", br-board)
		bs := float32(bb_territory+white_prisoners+s.BlackCap) - s.Komi
		ws := float32(ww_territory + black_prisoners + s.WhiteCap)

		s.BlackTerritory = fmt.Sprintf("%v", bs)
		s.WhiteTerritory = fmt.Sprintf("%v", ws)
		s.JP = fmt.Sprintf("%.2f", bs-ws)
		s.Pos = self.AreaResult.CNResult
	} else {
		s.CN = fmt.Sprintf("%.2f", (black_territory-white_territory-float64(s.Komi))/2)
		s.JP = fmt.Sprintf("%.2f", black_territory-white_territory-float64(s.Komi))
		s.Pos = self.AreaResult.JPResult
	}
	return s
}

func (self *ScoreUtil) ScoreCustom(max, min float64, limit int) ScoreResult {
	fmt.Print(max, min)
	s := ScoreResult{}
	if self.AreaResult.Flag == GOFLAG {
		self.ZHUScore()
	} else {
		self.ANScoreCustom(max, min, max, min)
	}
	s.Komi = self.Kifu.Komi
	s.Handicap = self.Kifu.Handicap
	//提子数
	s.BlackCap = self.Kifu.CurPos.BlackCap
	s.WhiteCap = self.Kifu.CurPos.WhiteCap
	//目数
	black_territory := float64(0)
	white_territory := float64(0)
	ww_territory := 0
	bb_territory := 0
	//死子
	white_prisoners := 0
	black_prisoners := 0
	//单官
	empty := 0
	//计算系数
	black_area := make([]string, 0) //all
	white_area := make([]string, 0) //all
	for i := int32(0); i < self.Kifu.Size; i++ {
		for j := int32(0); j < self.Kifu.Size; j++ {
			coor := self.Kifu.CurPos.GetCoor(i, j)
			o := self.Kifu.CurPos.GetPosition(i, j)
			cn := self.AreaResult.CNResult[coor]
			jp := self.AreaResult.JPResult[coor]
			coon := fmt.Sprintf("%s%s", IntToChar(i), IntToChar(j))
			if cn > 0 {
				black_area = append(black_area, coon)
				if o == W {
					white_prisoners = white_prisoners + 1
					bb_territory += 1
				}
				if o == Empty {
					bb_territory += 1
				}

			} else if cn < 0 {
				white_area = append(white_area, coon)
				if o == B {
					black_prisoners = black_prisoners + 1
					ww_territory += 1
				}
				if o == Empty {
					ww_territory += 1
				}

			} else {
				empty += 1
			}
			if jp > 0 {
				black_territory += jp
			} else if jp < 0 {
				white_territory += math.Abs(jp)
			}
		}
	}
	s.BlackArea = len(black_area)
	s.WhiteArea = len(white_area)
	s.Black = strings.Join(black_area, ":")
	s.White = strings.Join(white_area, ":")
	s.BlackTerritory = fmt.Sprintf("%.2f", black_territory)
	s.WhiteTerritory = fmt.Sprintf("%.2f", white_territory)
	s.WhitePrisoners = white_prisoners
	s.BlackPrisoners = black_prisoners
	s.Empty = empty
	sear := s.Komi
	if s.Handicap > 0 {
		if s.Handicap == 1 {
			sear = 0
		} else {
			sear = float32(s.Handicap)
		}
	}

	if empty <= limit && self.Kifu.NodeCount >= self.Config.HandsCount || self.Kifu.Size < 19 {
		board := float32(self.Kifu.Size*self.Kifu.Size) / 2
		br := float32(s.BlackArea) + float32(s.Empty)/2 - sear/2
		s.CN = fmt.Sprintf("%.2f", br-board)
		bs := float32(bb_territory+white_prisoners+s.BlackCap) - s.Komi
		ws := float32(ww_territory + black_prisoners + s.WhiteCap)
		s.BlackTerritory = fmt.Sprintf("%v", bs)
		s.WhiteTerritory = fmt.Sprintf("%v", ws)
		s.JP = fmt.Sprintf("%.2f", bs-ws)
		s.Pos = self.AreaResult.CNResult
	} else {
		s.CN = fmt.Sprintf("%.2f", (black_territory-white_territory-float64(s.Komi))/2)
		s.JP = fmt.Sprintf("%.2f", black_territory-white_territory-float64(s.Komi))
		s.Pos = self.AreaResult.JPResult
	}
	return s
}

func (self *ScoreUtil) ZHUScore() {
	score_r := make([]float64, self.Kifu.Size*self.Kifu.Size)
	for i, v := range self.AreaResult.CurPos.Schema {
		score_r[i] = float64(v)
	}
	self.AreaResult.JPResult = score_r
	self.AreaResult.CNResult = score_r
	return
}

func (self *ScoreUtil) ANScore(mode int, color int32) {
	r := self.Config.GetRegion(mode, self.Kifu.NodeCount)
	if color == 0 {
		color = self.Kifu.CurColor
	}
	bm := r.MMax
	bs := r.MMin
	wm := r.SMax
	ws := r.SMin
	if mode == C_MODE && color == B {
		bm = r.SMax
		bs = r.SMin
		wm = r.MMax
		ws = r.MMin
	}
	er := self.Config.GetRegion(E_MODE, self.Kifu.NodeCount)
	cnResult := self.clacANResult(er.MMax, er.MMin, er.SMax, er.SMin, true)
	self.AreaResult.CNResult = cnResult
	jpResult := self.clacANResult(bm, bs, wm, ws, false)
	self.AreaResult.JPResult = jpResult

}
func (self *ScoreUtil) ANScoreCustom(bm, bs, wm, ws float64) {
	cnResult := self.clacANResult(bm, bs, wm, ws, true)
	self.AreaResult.CNResult = cnResult
	jpResult := self.clacANResult(bm, bs, wm, ws, false)
	self.AreaResult.JPResult = jpResult
}

func (self ScoreUtil) clacANResult(bm, bs, wm, ws float64, isCalc bool) []float64 {
	pos := self.AreaResult.ANResult
	myPos := NewPosition(self.Kifu.Size)
	result := make([]float64, self.Kifu.Size*self.Kifu.Size)
	black_can_reach := make([]bool, self.Kifu.Size*self.Kifu.Size)
	white_can_reach := make([]bool, self.Kifu.Size*self.Kifu.Size)
	for x := int32(0); x < self.Kifu.Size; x++ {
		for y := int32(0); y < self.Kifu.Size; y++ {
			result[myPos.GetCoor(x, y)] = NOTHING
			coor := y*self.Kifu.Size + x
			coorValue := pos[coor]
			if coorValue >= bm || coorValue >= bs && self.Kifu.CurPos.GetPosition(x, y) != Empty {
				myPos.SetPosition(x, y, B)
				result[myPos.GetCoor(x, y)] = float64(B)
			} else if coorValue <= (-wm) || coorValue <= (-ws) && self.Kifu.CurPos.GetPosition(x, y) != Empty {
				myPos.SetPosition(x, y, W)
				result[myPos.GetCoor(x, y)] = float64(W)
			}
			//如果为BEGIN不计算
			black_can_reach[coor] = myPos.GetPosition(x, y) == B
			white_can_reach[coor] = myPos.GetPosition(x, y) == W
		}
	}
	find_new_point := isCalc
	for find_new_point {
		find_new_point = false
		for x := int32(0); x < self.Kifu.Size; x++ {
			for y := int32(0); y < self.Kifu.Size; y++ {
				if self.Kifu.CurPos.GetPosition(x, y) == Empty {

					for _, v := range self.Kifu.CurPos.Neighbar4(x, y) {
						if !black_can_reach[y*self.Kifu.Size+x] && black_can_reach[v.GetYSizeCoor(self.Kifu.Size)] {
							black_can_reach[y*self.Kifu.Size+x] = true
							find_new_point = true
						}
						if !white_can_reach[y*self.Kifu.Size+x] && white_can_reach[v.GetYSizeCoor(self.Kifu.Size)] {
							white_can_reach[y*self.Kifu.Size+x] = true
							find_new_point = true
						}
					}
				}
			}
		}
	}
	for x := int32(0); x < self.Kifu.Size; x++ {
		for y := int32(0); y < self.Kifu.Size; y++ {
			if result[myPos.GetCoor(x, y)] == NOTHING {
				coor := y*self.Kifu.Size + x
				if black_can_reach[coor] && white_can_reach[coor] {
					result[myPos.GetCoor(x, y)] = 0
				} else {
					coorValue := pos[coor]
					if coorValue > 0 {
						if isCalc {
							result[myPos.GetCoor(x, y)] = B
						} else {
							result[myPos.GetCoor(x, y)] = coorValue
						}
					} else {
						if isCalc {
							result[myPos.GetCoor(x, y)] = W
						} else {
							result[myPos.GetCoor(x, y)] = coorValue
						}

					}
				}
			}
		}
	}
	return result
}
