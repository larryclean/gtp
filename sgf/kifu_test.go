package sgf

import "testing"

func TestKifu_ToSGF(t *testing.T) {
	kifu := ParseSgf("(;GM[1]PB[张键光]BR[7.2D]PW[李刚]WR[5.4D]RE[暂无结果]KM[7.5]HA[0]TM[600]DT[2017-12-28]GN[标准区即时]SO[弈客围棋]OT[3/0.5]FF[4]CA[UTF-8]RU[zh]AP[WGo.js:2]BL[600]WL[600]SZ[19]BL[600]WL[600];B[dd];W[pp];B[cp];W[pc];B[pe];W[qe])")
	sAll := kifu.ToSGF(-1)
	if sAll == "(;SZ[19];B[dd];W[pp];B[cp];W[pc];B[pe];W[qe])" {
		t.Log(sAll)
	} else {
		t.Error("kifu to sgf all move fail")
	}
	s2 := kifu.ToSGF(2)
	if s2 == "(;SZ[19];B[dd];W[pp])" {
		t.Log(s2)
	} else {
		t.Error("kifu to sgf 2 move fail")
	}
}

