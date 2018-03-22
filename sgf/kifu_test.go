package sgf

import (
	"testing"
	"fmt"
)

func TestKifu_ToSGF(t *testing.T) {
	kifu := ParseSgf("(;GM[1]PB[张键光]BR[7.2D]PW[李刚]WR[5.4D]RE[暂无结果]KM[7.5]HA[0]TM[600]DT[2017-12-28]GN[标准区即时]SO[弈客围棋]OT[3/0.5]FF[4]CA[UTF-8]RU[zh]AP[WGo.js:2]BL[600]WL[600]SZ[19]BL[600]WL[600];B[dd];W[pp];B[cp];W[pc];B[pe];W[qe])")
	kifu.Last()
	sAll := kifu.ToCurSgf()
	fmt.Println(sAll)
	if sAll == "(;SZ[19]KM[7.5]HA[0];B[dd];W[pp];B[cp];W[pc];B[pe];W[qe])" {
		t.Log(sAll)
	} else {
		t.Error("kifu to sgf all move fail")
	}

}

